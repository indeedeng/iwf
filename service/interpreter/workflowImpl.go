package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"time"
)

func InterpreterImpl(ctx UnifiedContext, provider WorkflowProvider, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	globalVersionProvider := newGlobalVersionProvider(provider)
	if globalVersionProvider.isAfterVersionOfUsingGlobalVersioning(ctx) {
		err := globalVersionProvider.upsertGlobalVersionSearchAttribute(ctx)
		if err != nil {
			return nil, err
		}
	}

	err := provider.UpsertSearchAttributes(ctx, map[string]interface{}{
		service.SearchAttributeIwfWorkflowType: input.IwfWorkflowType,
	})
	if err != nil {
		return nil, err
	}

	execution := service.IwfWorkflowExecution{
		IwfWorkerUrl:     input.IwfWorkerUrl,
		WorkflowType:     input.IwfWorkflowType,
		WorkflowId:       provider.GetWorkflowInfo(ctx).WorkflowExecution.ID,
		RunId:            provider.GetWorkflowInfo(ctx).WorkflowExecution.RunID,
		StartedTimestamp: provider.GetWorkflowInfo(ctx).WorkflowStartTime.Unix(),
	}
	stateExeIdMgr := NewStateExecutionIdManager()
	interStateChannel := NewInterStateChannel()
	currentStates := []iwfidl.StateMovement{
		{
			StateId:      input.StartStateId,
			StateOptions: &input.StateOptions,
			StateInput:   &input.StateInput,
		},
	}
	persistenceManager := NewPersistenceManager(func(attributes map[string]interface{}) error {
		return provider.UpsertSearchAttributes(ctx, attributes)
	})

	err = provider.SetQueryHandler(ctx, service.GetDataObjectsWorkflowQueryType, func(req service.GetDataObjectsQueryRequest) (service.GetDataObjectsQueryResponse, error) {
		return persistenceManager.GetDataObjectsByKey(req), nil
	})
	if err != nil {
		return nil, err
	}

	var errToFailWf error // TODO Note that today different errors could overwrite each other, we only support last one wins. we may use multiError to improve.
	var outputsToReturnWf []iwfidl.StateCompletionOutput
	var forceCompleteWf bool
	stateExecutingMgr := newStateExecutingManager(ctx, provider)
	//inFlightExecutingStateCount := 0

	for len(currentStates) > 0 {
		// copy the whole slice(pointer)
		statesToExecute := currentStates
		err := stateExecutingMgr.startStates(currentStates)
		if err != nil {
			return nil, err
		}
		//reset to empty slice since each iteration will process all current states in the queue
		currentStates = nil

		for _, stateToExecute := range statesToExecute {
			// execute in another thread for parallelism
			// state must be passed via parameter https://stackoverflow.com/questions/67263092
			stateCtx := provider.ExtendContextWithValue(ctx, "state", stateToExecute)
			provider.GoNamed(stateCtx, stateToExecute.GetStateId(), func(ctx UnifiedContext) {
				state, ok := provider.GetContextValue(ctx, "state").(iwfidl.StateMovement)
				if !ok {
					errToFailWf = provider.NewApplicationError(
						"critical code bug when passing state via context",
						service.WorkflowErrorTypeUserInternalError,
					)
					return
				}
				defer func() {
					err := stateExecutingMgr.completeStates(state)
					if err != nil {
						errToFailWf = err
					}
				}()

				stateExeId := stateExeIdMgr.IncAndGetNextExecutionId(state.GetStateId())
				decision, err := executeState(ctx, provider, state, execution, stateExeId, persistenceManager, interStateChannel)
				if err != nil {
					errToFailWf = err
				}

				shouldClose, gracefulComplete, forceComplete, forceFail, output, err := checkClosingWorkflow(provider, decision, state.GetStateId(), stateExeId)
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
					errToFailWf = provider.NewApplicationError(
						fmt.Sprintf("user workflow decided to fail workflow execution stateId %s, stateExecutionId: %s", state.GetStateId(), stateExeId),
						service.WorkflowErrorTypeUserWorkflowDecision,
					)
				}
				if !shouldClose && decision.HasNextStates() {
					currentStates = append(currentStates, decision.GetNextStates()...)
				}
			})
		}

		awaitError := provider.Await(ctx, func() bool {
			return len(currentStates) > 0 || errToFailWf != nil || forceCompleteWf || stateExecutingMgr.getTotalExecutingStates() == 0
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
	provider WorkflowProvider, decision *iwfidl.StateDecision, currentStateId, currentStateExeId string,
) (shouldClose, gracefulComplete, forceComplete, forceFail bool, completeOutput *iwfidl.StateCompletionOutput, err error) {
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.GracefulCompletingWorkflowStateId {
			shouldClose = true
			gracefulComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceCompletingWorkflowStateId {
			shouldClose = true
			forceComplete = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceFailingWorkflowStateId {
			shouldClose = true
			forceFail = true
		}
	}
	if shouldClose && len(decision.NextStates) > 1 {
		// Illegal decision
		err = provider.NewApplicationError(
			"closing workflow decision should have only one state movement, but got more than one",
			service.WorkflowErrorTypeUserWorkflowError,
		)
		return
	}
	return
}

func executeState(
	ctx UnifiedContext,
	provider WorkflowProvider,
	state iwfidl.StateMovement,
	execution service.IwfWorkflowExecution,
	stateExeId string,
	attrMgr *PersistenceManager,
	interStateChannel *InterStateChannel,
) (*iwfidl.StateDecision, error) {
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	if state.StateOptions != nil {
		if state.StateOptions.GetStartApiTimeoutSeconds() > 0 {
			activityOptions.StartToCloseTimeout = time.Duration(state.StateOptions.GetStartApiTimeoutSeconds()) * time.Second
		}
		activityOptions.RetryPolicy = state.StateOptions.StartApiRetryPolicy
	}

	ctx = provider.WithActivityOptions(ctx, activityOptions)

	exeCtx := iwfidl.Context{
		WorkflowId:               execution.WorkflowId,
		WorkflowRunId:            execution.RunId,
		WorkflowStartedTimestamp: execution.StartedTimestamp,
		StateExecutionId:         stateExeId,
	}

	var startResponse *iwfidl.WorkflowStateStartResponse
	err := provider.ExecuteActivity(ctx, StateStart, provider.GetBackendType(), service.StateStartActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateStartRequest{
			Context:          exeCtx,
			WorkflowType:     execution.WorkflowType,
			WorkflowStateId:  state.StateId,
			StateInput:       state.StateInput,
			SearchAttributes: attrMgr.LoadSearchAttributes(state.StateOptions),
			DataObjects:      attrMgr.LoadDataObjects(state.StateOptions),
		},
	}).Get(ctx, &startResponse)
	if err != nil {
		return nil, err
	}

	err = attrMgr.ProcessUpsertSearchAttribute(startResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, err
	}
	err = attrMgr.ProcessUpsertDataObject(startResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, err
	}
	interStateChannel.ProcessPublishing(startResponse.GetPublishToInterStateChannel())

	commandReq := startResponse.GetCommandRequest()
	commandReqDone := false

	completedTimerCmds := map[int]bool{}
	if len(commandReq.GetTimerCommands()) > 0 {
		for idx, cmd := range commandReq.GetTimerCommands() {
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getThreadName("timer", cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				cmd, ok := provider.GetContextValue(ctx, "cmd").(iwfidl.TimerCommand)
				if !ok {
					panic("critical code bug")
				}
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				now := provider.Now(ctx).Unix()
				fireAt := cmd.GetFiringUnixTimestampSeconds()
				duration := time.Duration(fireAt-now) * time.Second
				future := provider.NewTimer(ctx, duration)
				_ = provider.Await(ctx, func() bool {
					return future.IsReady() || commandReqDone
				})
				if future.IsReady() {
					completedTimerCmds[idx] = true
				}
			})
		}
	}

	completedSignalCmds := map[int]*iwfidl.EncodedObject{}
	if len(commandReq.GetSignalCommands()) > 0 {
		for idx, cmd := range commandReq.GetSignalCommands() {
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getThreadName("signal", cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				cmd, ok := provider.GetContextValue(ctx, "cmd").(iwfidl.SignalCommand)
				if !ok {
					panic("critical code bug")
				}
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}
				ch := provider.GetSignalChannel(ctx, cmd.GetSignalChannelName())
				value := iwfidl.EncodedObject{}
				received := false
				_ = provider.Await(ctx, func() bool {
					received = ch.ReceiveAsync(&value)
					return received || commandReqDone
				})
				if received {
					completedSignalCmds[idx] = &value
				}
			})
		}
	}

	completedInterStateChannelCmds := map[int]*iwfidl.EncodedObject{}
	if len(commandReq.GetInterStateChannelCommands()) > 0 {
		for idx, cmd := range commandReq.GetInterStateChannelCommands() {
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getThreadName("interstate", cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				cmd, ok := provider.GetContextValue(ctx, "cmd").(iwfidl.InterStateChannelCommand)
				if !ok {
					panic("critical code bug")
				}
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				received := false
				_ = provider.Await(ctx, func() bool {
					received = interStateChannel.HasData(cmd.ChannelName)
					return received || commandReqDone
				})

				if received {
					completedInterStateChannelCmds[idx] = interStateChannel.Retrieve(cmd.ChannelName)
				}
			})
		}
	}

	// TODO process long running activity command

	if len(commandReq.GetTimerCommands())+len(commandReq.GetSignalCommands())+len(commandReq.GetInterStateChannelCommands()) > 0 {
		triggerType := commandReq.GetDeciderTriggerType()
		if triggerType == iwfidl.ALL_COMMAND_COMPLETED {
			err = provider.Await(ctx, func() bool {
				return len(completedTimerCmds) == len(commandReq.GetTimerCommands()) &&
					len(completedSignalCmds) == len(commandReq.GetSignalCommands()) &&
					len(completedInterStateChannelCmds) == len(commandReq.GetInterStateChannelCommands())
			})
		} else if triggerType == iwfidl.ANY_COMMAND_COMPLETED {
			err = provider.Await(ctx, func() bool {
				return len(completedTimerCmds)+
					len(completedSignalCmds)+
					len(completedInterStateChannelCmds) > 0
			})
		} else {
			return nil, provider.NewApplicationError("unsupported decider trigger type", "unsupported", triggerType)
		}
	}
	commandReqDone = true

	if err != nil {
		return nil, err
	}
	commandRes := &iwfidl.CommandResults{}
	if len(commandReq.GetTimerCommands()) > 0 {
		var timerResults []iwfidl.TimerResult
		for idx, cmd := range commandReq.GetTimerCommands() {
			status := iwfidl.FIRED
			if !completedTimerCmds[idx] {
				status = iwfidl.SCHEDULED
			}
			timerResults = append(timerResults, iwfidl.TimerResult{
				CommandId:   cmd.GetCommandId(),
				TimerStatus: status,
			})
		}
		commandRes.SetTimerResults(timerResults)
	}

	if len(commandReq.GetSignalCommands()) > 0 {
		var signalResults []iwfidl.SignalResult
		for idx, cmd := range commandReq.GetSignalCommands() {
			status := iwfidl.RECEIVED
			result, completed := completedSignalCmds[idx]
			if !completed {
				status = iwfidl.WAITING
			}

			signalResults = append(signalResults, iwfidl.SignalResult{
				CommandId:           cmd.GetCommandId(),
				SignalChannelName:   cmd.GetSignalChannelName(),
				SignalValue:         result,
				SignalRequestStatus: status,
			})
		}
		commandRes.SetSignalResults(signalResults)
	}

	if len(commandReq.GetInterStateChannelCommands()) > 0 {
		var interStateChannelResults []iwfidl.InterStateChannelResult
		for idx, cmd := range commandReq.GetInterStateChannelCommands() {
			status := iwfidl.RECEIVED
			result, completed := completedInterStateChannelCmds[idx]
			if !completed {
				status = iwfidl.WAITING
			}

			interStateChannelResults = append(interStateChannelResults, iwfidl.InterStateChannelResult{
				CommandId:     cmd.CommandId,
				ChannelName:   cmd.ChannelName,
				RequestStatus: status,
				Value:         result,
			})
		}
		commandRes.SetInterStateChannelResults(interStateChannelResults)
	}

	activityOptions = ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	if state.StateOptions != nil {
		if state.StateOptions.GetDecideApiTimeoutSeconds() > 0 {
			activityOptions.StartToCloseTimeout = time.Duration(state.StateOptions.GetDecideApiTimeoutSeconds()) * time.Second
		}
		activityOptions.RetryPolicy = state.StateOptions.DecideApiRetryPolicy
	}

	ctx = provider.WithActivityOptions(ctx, activityOptions)
	var decideResponse *iwfidl.WorkflowStateDecideResponse
	err = provider.ExecuteActivity(ctx, StateDecide, provider.GetBackendType(), service.StateDecideActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateDecideRequest{
			Context:          exeCtx,
			WorkflowType:     execution.WorkflowType,
			WorkflowStateId:  state.StateId,
			CommandResults:   commandRes,
			StateLocals:      startResponse.GetUpsertStateLocals(),
			SearchAttributes: attrMgr.LoadSearchAttributes(state.StateOptions),
			DataObjects:      attrMgr.LoadDataObjects(state.StateOptions),
			StateInput:       state.StateInput,
		},
	}).Get(ctx, &decideResponse)
	if err != nil {
		return nil, err
	}

	decision := decideResponse.GetStateDecision()
	err = attrMgr.ProcessUpsertSearchAttribute(decideResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, err
	}
	err = attrMgr.ProcessUpsertDataObject(decideResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, err
	}
	interStateChannel.ProcessPublishing(decideResponse.GetPublishToInterStateChannel())

	return &decision, nil
}

func getThreadName(prefix string, cmdId string, idx int) string {
	return fmt.Sprintf("%v-%v-%v", prefix, cmdId, idx)
}
