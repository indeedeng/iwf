package temporal

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	execution := service.IwfWorkflowExecution{
		IwfWorkerUrl:     input.IwfWorkerUrl,
		WorkflowType:     input.IwfWorkflowType,
		WorkflowId:       workflow.GetInfo(ctx).WorkflowExecution.ID,
		RunId:            workflow.GetInfo(ctx).WorkflowExecution.RunID,
		StartedTimestamp: int32(workflow.GetInfo(ctx).WorkflowStartTime.Unix()),
	}
	stateExeIdMgr := newStateExecutionIdManager()
	currentStates := []iwfidl.StateMovement{
		{
			StateId:          &input.StartStateId,
			NextStateOptions: &input.StateOptions,
			NextStateInput:   &input.StateInput,
		},
	}

	var err error
	for len(currentStates) > 0 {
		statesToExecute := currentStates
		//reset to empty slice since each iteration will process all current states in the queue
		currentStates = nil

		for _, state := range statesToExecute {
			decision, err := executeState(ctx, state, execution, stateExeIdMgr)
			if err != nil {
				return nil, err
			}
			// TODO process search attributes
			// TODO process query attributes

			isClosing, output, err := checkClosingWorkflow(decision)
			if isClosing {
				return output, err
			}
			if decision.HasNextStates() {
				currentStates = append(currentStates, decision.GetNextStates()...)
			}

		}
		err = workflow.Await(ctx, func() bool {
			return len(currentStates) > 0
		})
		if err != nil {
			break
		}
	}

	return nil, err
}

func checkClosingWorkflow(decision *iwfidl.StateDecision) (bool, *service.InterpreterWorkflowOutput, error) {
	hasClosingDecision := false
	var output *service.InterpreterWorkflowOutput
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.CompletingWorkflowStateId || stateId == service.FailingWorkflowStateId {
			hasClosingDecision = true
			output = &service.InterpreterWorkflowOutput{
				CompletedStateId: stateId,
				StateOutput:      movement.GetNextStateInput(),
			}
		}
	}
	if hasClosingDecision && len(decision.NextStates) > 1 {
		// Illegal decision, should fail the workflow
		return true, nil, temporal.NewApplicationError(
			"closing workflow decision shouldn't have other state movements",
			"Illegal closing workflow decision",
		)
	}
	return hasClosingDecision, output, nil
}

func executeState(
	ctx workflow.Context, state iwfidl.StateMovement, execution service.IwfWorkflowExecution, idMgr *stateExecutionIdManager,
) (*iwfidl.StateDecision, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	stateExeId := idMgr.incAndGetNextExecutionId(state.GetStateId())
	exeCtx := iwfidl.Context{
		WorkflowId:               &execution.WorkflowId,
		WorkflowRunId:            &execution.RunId,
		WorkflowStartedTimestamp: &execution.StartedTimestamp,
		StateExecutionId:         &stateExeId,
	}

	var startResponse *iwfidl.WorkflowStateStartResponse
	err := workflow.ExecuteActivity(ctx, StateStartActivity, service.StateStartActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateStartRequest{
			Context:          &exeCtx,
			WorkflowType:     &execution.WorkflowType,
			WorkflowStateId:  state.StateId,
			StateInput:       state.NextStateInput,
			SearchAttributes: nil, // TODO
			QueryAttributes:  nil, // TODO
		},
	}).Get(ctx, &startResponse)
	if err != nil {
		return nil, err
	}

	// TODO process timer command
	// TODO process signal command
	// TODO process long running activity command
	// TODO process upsert search attribute
	// TODO process upsert query attribute
	// TODO process state local attribute

	commandReq := startResponse.GetCommandRequest()
	triggerType := commandReq.GetDeciderTriggerType()
	if triggerType != "ALL_COMMAND_COMPLETED" {
		return nil, temporal.NewApplicationError("unsupported decider trigger type", "unsupported", triggerType)
	}

	var decideResponse *iwfidl.WorkflowStateDecideResponse
	err = workflow.ExecuteActivity(ctx, StateDecideActivity, service.StateDecideActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateDecideRequest{
			Context:              &exeCtx,
			WorkflowType:         &execution.WorkflowType,
			WorkflowStateId:      state.StateId,
			CommandResults:       nil, // TODO
			StateLocalAttributes: nil, // TODO
			SearchAttributes:     nil, // TODO
			QueryAttributes:      nil, // TODO
		},
	}).Get(ctx, &decideResponse)
	if err != nil {
		return nil, err
	}

	return decideResponse.StateDecision, nil
}
