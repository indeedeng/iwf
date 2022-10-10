package temporal

import (
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter"
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
		StartedTimestamp: workflow.GetInfo(ctx).WorkflowStartTime.Unix(),
	}
	stateExeIdMgr := interpreter.NewStateExecutionIdManager()
	currentStates := []iwfidl.StateMovement{
		{
			StateId:          input.StartStateId,
			NextStateOptions: &input.StateOptions,
			NextStateInput:   &input.StateInput,
		},
	}
	attrMgr := interpreter.NewAttributeManager(func(attributes map[string]interface{}) error {
		return workflow.UpsertSearchAttributes(ctx, attributes)
	})

	err := workflow.SetQueryHandler(ctx, service.AttributeQueryType, func(req service.QueryAttributeRequest) (service.QueryAttributeResponse, error) {
		return attrMgr.GetQueryAttributesByKey(req), nil
	})
	if err != nil {
		return nil, err
	}

	var errToFailWf error // TODO Note that today different errors could overwrite each other, we only support last one wins. we may use multiError to improve.
	var outputsToReturnWf []iwfidl.StateCompletionOutput
	var forceCompleteWf bool
	inFlightExecutingStateCount := 0

	for len(currentStates) > 0 {
		// copy the whole slice(pointer)
		inFlightExecutingStateCount += len(currentStates)
		statesToExecute := currentStates
		//reset to empty slice since each iteration will process all current states in the queue
		currentStates = nil

		for _, state := range statesToExecute {
			// execute in another thread for parallelism
			// state must be passed via parameter https://stackoverflow.com/questions/67263092
			stateCtx := workflow.WithValue(ctx, "state", state)
			workflow.GoNamed(stateCtx, state.GetStateId(), func(ctx workflow.Context) {
				defer func() {
					inFlightExecutingStateCount--
				}()

				thisState, ok := ctx.Value("state").(iwfidl.StateMovement)
				if !ok {
					errToFailWf = temporal.NewApplicationError(
						"critical code bug when passing state via context",
						service.WorkflowErrorTypeUserInternalError,
					)
					return
				}

				stateExeId := stateExeIdMgr.IncAndGetNextExecutionId(state.GetStateId())
				decision, err := executeState(ctx, thisState, execution, stateExeId, attrMgr)
				if err != nil {
					errToFailWf = err
				}

				shouldClose, gracefulComplete, forceComplete, forceFail, output, err := checkClosingWorkflow(decision, state.GetStateId(), stateExeId)
				if err != nil {
					errToFailWf = err
				}
				if gracefulComplete || forceComplete {
					outputsToReturnWf = append(outputsToReturnWf, *output)
				}
				if forceComplete {
					forceCompleteWf = true
				}
				if forceFail {
					errToFailWf = temporal.NewApplicationError(
						fmt.Sprintf("user workflow decided to fail workflow execution stateId %s, stateExecutionId: %s", state.GetStateId(), stateExeId),
						service.WorkflowErrorTypeUserWorkflowDecision,
					)
				}
				if !shouldClose && decision.HasNextStates() {
					currentStates = append(currentStates, decision.GetNextStates()...)
				}
			})
		}

		awaitError := workflow.Await(ctx, func() bool {
			return len(currentStates) > 0 || errToFailWf != nil || forceCompleteWf || inFlightExecutingStateCount == 0
		})
		if errToFailWf != nil || forceCompleteWf {
			return &service.InterpreterWorkflowOutput{
				StateCompletionOutputs: outputsToReturnWf,
			}, errToFailWf
		}

		if awaitError != nil {
			// this could happen for cancellation
			errToFailWf = awaitError
			break
		}
	}

	// gracefully complete workflow when all states are executed to dead ends
	return &service.InterpreterWorkflowOutput{
		StateCompletionOutputs: outputsToReturnWf,
	}, errToFailWf
}

func checkClosingWorkflow(
	decision *iwfidl.StateDecision, currentStateId, currentStateExeId string,
) (shouldClose, gracefulComplete, forceComplete, forceFail bool, completeOutput *iwfidl.StateCompletionOutput, err error) {
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.GracefulCompletingWorkflowStateId {
			shouldClose = true
			gracefulComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.NextStateInput,
			}
		}
		if stateId == service.ForceCompletingWorkflowStateId {
			shouldClose = true
			forceComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.NextStateInput,
			}
		}
		if stateId == service.ForceFailingWorkflowStateId {
			shouldClose = true
			forceFail = true
		}
	}
	if shouldClose && len(decision.NextStates) > 1 {
		// Illegal decision
		err = temporal.NewApplicationError(
			"closing workflow decision should have only one state movement, but got more than one",
			service.WorkflowErrorTypeUserWorkflowError,
		)
		return
	}
	return
}

func executeState(
	ctx workflow.Context,
	state iwfidl.StateMovement,
	execution service.IwfWorkflowExecution,
	stateExeId string,
	attrMgr *interpreter.AttributeManager,
) (*iwfidl.StateDecision, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	exeCtx := iwfidl.Context{
		WorkflowId:               execution.WorkflowId,
		WorkflowRunId:            execution.RunId,
		WorkflowStartedTimestamp: execution.StartedTimestamp,
		StateExecutionId:         stateExeId,
	}

	var startResponse *iwfidl.WorkflowStateStartResponse
	err := workflow.ExecuteActivity(ctx, interpreter.StateStart, service.BackendTypeTemporal, service.StateStartActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateStartRequest{
			Context:          exeCtx,
			WorkflowType:     execution.WorkflowType,
			WorkflowStateId:  state.StateId,
			StateInput:       state.NextStateInput,
			SearchAttributes: attrMgr.GetAllSearchAttributes(),
			QueryAttributes:  attrMgr.GetAllQueryAttributes(),
		},
	}).Get(ctx, &startResponse)
	if err != nil {
		return nil, err
	}

	err = attrMgr.ProcessUpsertSearchAttribute(startResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, err
	}
	err = attrMgr.ProcessUpsertQueryAttribute(startResponse.GetUpsertQueryAttributes())
	if err != nil {
		return nil, err
	}

	commandReq := startResponse.GetCommandRequest()

	completedTimerCmds := 0
	if len(commandReq.GetTimerCommands()) > 0 {
		for _, cmd := range commandReq.GetTimerCommands() {
			cmdCtx := workflow.WithValue(ctx, "cmd", cmd)
			workflow.Go(cmdCtx, func(ctx workflow.Context) {
				cmd, ok := ctx.Value("cmd").(iwfidl.TimerCommand)
				if !ok {
					panic("critical code bug")
				}

				now := workflow.Now(ctx).Unix()
				fireAt := cmd.GetFiringUnixTimestampSeconds()
				duration := time.Duration(fireAt-now) * time.Second
				_ = workflow.Sleep(ctx, duration)
				completedTimerCmds++
			})
		}
	}

	completedSignalCmds := map[string]*iwfidl.EncodedObject{}
	if len(commandReq.GetSignalCommands()) > 0 {
		for _, cmd := range commandReq.GetSignalCommands() {
			cmdCtx := workflow.WithValue(ctx, "cmd", cmd)
			workflow.Go(cmdCtx, func(ctx workflow.Context) {
				cmd, ok := ctx.Value("cmd").(iwfidl.SignalCommand)
				if !ok {
					panic("critical code bug")
				}
				ch := workflow.GetSignalChannel(ctx, cmd.GetSignalName())
				value := iwfidl.EncodedObject{}
				ch.Receive(ctx, &value)
				completedSignalCmds[cmd.GetCommandId()] = &value
			})
		}
	}

	// TODO process long running activity command

	triggerType := commandReq.GetDeciderTriggerType()
	if triggerType != service.DeciderTypeAllCommandCompleted {
		return nil, temporal.NewApplicationError("unsupported decider trigger type", "unsupported", triggerType)
	}

	err = workflow.Await(ctx, func() bool {
		return completedTimerCmds == len(commandReq.GetTimerCommands()) &&
			len(completedSignalCmds) == len(commandReq.GetSignalCommands())
	})
	if err != nil {
		return nil, err
	}
	commandRes := &iwfidl.CommandResults{}
	if len(commandReq.GetTimerCommands()) > 0 {
		var timerResults []iwfidl.TimerResult
		for _, cmd := range commandReq.GetTimerCommands() {
			timerResults = append(timerResults, iwfidl.TimerResult{
				CommandId:   cmd.GetCommandId(),
				TimerStatus: service.TimerStatusFired,
			})
		}
		commandRes.SetTimerResults(timerResults)
	}

	if len(commandReq.GetSignalCommands()) > 0 {
		var signalResults []iwfidl.SignalResult
		for _, cmd := range commandReq.GetSignalCommands() {
			signalResults = append(signalResults, iwfidl.SignalResult{
				CommandId:    cmd.GetCommandId(),
				SignalName:   cmd.GetSignalName(),
				SignalValue:  completedSignalCmds[cmd.GetCommandId()],
				SignalStatus: service.SignalStatusReceived,
			})
		}
		commandRes.SetSignalResults(signalResults)
	}

	var decideResponse *iwfidl.WorkflowStateDecideResponse
	err = workflow.ExecuteActivity(ctx, interpreter.StateDecide, service.BackendTypeTemporal, service.StateDecideActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateDecideRequest{
			Context:              exeCtx,
			WorkflowType:         execution.WorkflowType,
			WorkflowStateId:      state.StateId,
			CommandResults:       commandRes,
			StateLocalAttributes: startResponse.GetUpsertStateLocalAttributes(),
			SearchAttributes:     attrMgr.GetAllSearchAttributes(),
			QueryAttributes:      attrMgr.GetAllQueryAttributes(),
			StateInput:           state.NextStateInput,
		},
	}).Get(ctx, &decideResponse)
	if err != nil {
		return nil, err
	}

	decision := decideResponse.GetStateDecision()
	err = attrMgr.ProcessUpsertSearchAttribute(decision.GetUpsertSearchAttributes())
	if err != nil {
		return nil, err
	}
	err = attrMgr.ProcessUpsertQueryAttribute(decision.GetUpsertQueryAttributes())
	if err != nil {
		return nil, err
	}

	return &decision, nil
}
