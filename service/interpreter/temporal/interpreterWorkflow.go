package temporal

import (
	"fmt"
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

	var errToReturn error
	var outputToReturn *service.InterpreterWorkflowOutput
	for len(currentStates) > 0 {
		// copy the whole slice(pointer)
		statesToExecute := currentStates
		//reset to empty slice since each iteration will process all current states in the queue
		currentStates = nil

		for _, state := range statesToExecute {
			// execute in another thread for parallelism
			// state must be passed via parameter https://stackoverflow.com/questions/67263092
			stateCtx := workflow.WithValue(ctx, "state", state)
			//stateCtx := newParametrizedContext(ctx2, &state)
			workflow.GoNamed(stateCtx, state.GetStateId(), func(ctx workflow.Context) {
				thisState, ok := ctx.Value("state").(iwfidl.StateMovement)
				if !ok {
					panic("critical code bug")
				}
				fmt.Println("check stateId", thisState.GetStateId())

				decision, err := executeState(ctx, thisState, execution, stateExeIdMgr)
				if err != nil {
					errToReturn = err
				}
				// TODO process search attributes
				// TODO process query attributes

				isClosing, output, err := checkClosingWorkflow(decision)
				if isClosing {
					errToReturn = err
					outputToReturn = output
				}
				if decision.HasNextStates() {
					currentStates = append(currentStates, decision.GetNextStates()...)
				}
			})
		}

		awaitError := workflow.Await(ctx, func() bool {
			return len(currentStates) > 0 || errToReturn != nil || outputToReturn != nil
		})
		if errToReturn != nil || outputToReturn != nil {
			return outputToReturn, errToReturn
		}

		if awaitError != nil {
			errToReturn = awaitError
			break
		}
	}

	return nil, errToReturn
}

func checkClosingWorkflow(decision *iwfidl.StateDecision) (bool, *service.InterpreterWorkflowOutput, error) {
	hasClosingDecision := false
	var output *service.InterpreterWorkflowOutput
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.CompletingWorkflowStateId {
			hasClosingDecision = true
			output = &service.InterpreterWorkflowOutput{
				CompletedStateExecutionId: "TODO", // TODO get prev state execution Id
				StateOutput:               movement.GetNextStateInput(),
			}
		}
		if stateId == service.FailingWorkflowStateId {
			return true, nil, temporal.NewApplicationError(
				"failing by user workflow decision",
				"failing on request",
			)
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
	if triggerType != service.DeciderTypeAllCommandCompleted {
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
