package interpreter

import (
	"fmt"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

func InterpreterImpl(ctx UnifiedContext, provider WorkflowProvider, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	var err error
	globalVersionProvider := NewGlobalVersioner(provider, ctx)
	if globalVersionProvider.IsAfterVersionOfUsingGlobalVersioning() {
		err = globalVersionProvider.UpsertGlobalVersionSearchAttribute()
		if err != nil {
			return nil, err
		}
	}

	if input.ContinueAsNew {
		return ResumeFromPreviousRun(input)
	}

	if !input.Config.GetDisableSystemSearchAttribute() {
		if !globalVersionProvider.IsAfterVersionOfOptimizedUpsertSearchAttribute() {
			// stop upsert it here since it's done in start workflow request
			err = provider.UpsertSearchAttributes(ctx, map[string]interface{}{
				service.SearchAttributeIwfWorkflowType: input.IwfWorkflowType,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	execution := service.IwfWorkflowExecution{
		IwfWorkerUrl:     input.IwfWorkerUrl,
		WorkflowType:     input.IwfWorkflowType,
		WorkflowId:       provider.GetWorkflowInfo(ctx).WorkflowExecution.ID,
		RunId:            provider.GetWorkflowInfo(ctx).WorkflowExecution.RunID,
		StartedTimestamp: provider.GetWorkflowInfo(ctx).WorkflowStartTime.Unix(),
	}
	interStateChannel := NewInterStateChannel()
	statesToExecuteQueue := []iwfidl.StateMovement{
		{
			StateId:      input.StartStateId,
			StateOptions: &input.StateOptions,
			StateInput:   &input.StateInput,
		},
	}

	persistenceManager := NewPersistenceManager(provider, input.InitSearchAttributes)
	timerProcessor := NewTimerProcessor(ctx, provider)
	continueAsNewCounter := NewContinueAsCounter(input.Config, ctx, provider)
	signalReceiver := NewSignalReceiver(ctx, provider, timerProcessor, continueAsNewCounter)

	err = provider.SetQueryHandler(ctx, service.GetDataObjectsWorkflowQueryType, func(req service.GetDataObjectsQueryRequest) (service.GetDataObjectsQueryResponse, error) {
		return persistenceManager.GetDataObjectsByKey(req), nil
	})
	if err != nil {
		return nil, err
	}
	err = provider.SetQueryHandler(ctx, service.GetSearchAttributesWorkflowQueryType, func() ([]iwfidl.SearchAttribute, error) {
		return persistenceManager.GetAllSearchAttributes(), nil
	})
	if err != nil {
		return nil, err
	}

	var errToFailWf error // Note that today different errors could overwrite each other, we only support last one wins. we may use multiError to improve.
	var outputsToReturnWf []iwfidl.StateCompletionOutput
	var forceCompleteWf bool
	stateExecutionCounter := NewStateExecutionCounter(ctx, provider, input.Config, continueAsNewCounter)

	continueAsNewer := NewContinueAsNewer(provider, interStateChannel, signalReceiver, stateExecutionCounter, persistenceManager)
	err = continueAsNewer.SetQueryHandlersForContinueAsNew(ctx)
	if err != nil {
		return nil, err
	}

	// this is for an optimization for StateId Search attribute, see updateStateIdSearchAttribute in stateExecutionCounter
	defer stateExecutionCounter.ClearStateIdSearchAttributeFinally()

	for len(statesToExecuteQueue) > 0 {
		// copy the whole slice(pointer)
		statesToExecute := statesToExecuteQueue
		err = stateExecutionCounter.MarkStateExecutionsPending(statesToExecuteQueue)
		if err != nil {
			return nil, err
		}
		//reset to empty slice since each iteration will process all current states in the queue
		statesToExecuteQueue = nil

		for _, stateToExecuteForLoopingOnly := range statesToExecute {
			// execute in another thread for parallelism
			// state must be passed via parameter https://stackoverflow.com/questions/67263092
			stateCtx := provider.ExtendContextWithValue(ctx, "state", stateToExecuteForLoopingOnly)
			provider.GoNamed(stateCtx, stateToExecuteForLoopingOnly.GetStateId(), func(ctx UnifiedContext) {
				state, ok := provider.GetContextValue(ctx, "state").(iwfidl.StateMovement)
				if !ok {
					errToFailWf = provider.NewApplicationError(
						string(iwfidl.SERVER_INTERNAL_ERROR_TYPE),
						"critical code bug when passing state via context",
					)
					return
				}

				stateExeId := stateExecutionCounter.CreateNextExecutionId(state.GetStateId())
				decision, stateExecStatus, err := executeState(
					ctx, provider, state, execution, stateExeId, persistenceManager,
					interStateChannel, signalReceiver, timerProcessor, continueAsNewer, continueAsNewCounter)
				if err != nil {
					errToFailWf = err
					// state execution fail should fail the workflow, no more processing
					return
				}

				if stateExecStatus == service.DecideApiCompletedStateExecutionStatus {
					// NOTE: decision is only available on this DecideApiCompletedStateExecutionStatus

					shouldClose, gracefulComplete, forceComplete, forceFail, output, err := checkClosingWorkflow(provider, decision, state.GetStateId(), stateExeId)
					if err != nil {
						errToFailWf = err
						// no return so that it can fall through to call MarkStateExecutionCompleted
					}
					if gracefulComplete || forceComplete || forceFail {
						outputsToReturnWf = append(outputsToReturnWf, *output)
					}
					if forceComplete {
						forceCompleteWf = true
					}
					if forceFail {
						errToFailWf = provider.NewApplicationError(
							string(iwfidl.STATE_DECISION_FAILING_WORKFLOW_ERROR_TYPE),
							outputsToReturnWf,
						)
						// no return so that it can fall through to call MarkStateExecutionCompleted
					}
					if !shouldClose && decision.HasNextStates() {
						statesToExecuteQueue = append(statesToExecuteQueue, decision.GetNextStates()...)
					}

					// finally, mark state completed and may also update system search attribute(IwfExecutingStateIds)
					// doing this at last because upsertSearchAttribute is blocking call for workflow [decision] task
					err = stateExecutionCounter.MarkStateExecutionCompleted(state)
					if err != nil {
						errToFailWf = err
					}
				} else {
					// if not completed (without error), then it's because of cancellation by continueAsNew
					continueAsNewer.ProcessUncompletedStateExecution(stateExecStatus, stateExeId, state)
				}
			}) // end of executing one state
		} // end loop of executing all states from the queue for one iteration

		// The conditions here are quite tricky:
		// For len(statesToExecuteQueue) > 0: We need some condition to wait here because all the state execution are running in different thread.
		//    Right after the queue are popped, the len(...) becomes zero. So when the len(...) >0, it means there are new states to execute pushed into the queue,
		//    and it's time to wake up the outer loop to go to next iteration. Alternatively, waiting for all current started in this iteration to complete will also work,
		//    but not as efficient as this one because it will take much longer time.
		// For errToFailWf != nil || forceCompleteWf: this means we need to close workflow immediately
		// For stateExecutionCounter.GetTotalPendingStateExecutions() == 0: this means all the state executions have reach "Dead Ends" so the workflow can complete gracefully without output
		// For continueAsNewCounter.IsThresholdMet(): this means workflow need to continueAsNew
		awaitError := provider.Await(ctx, func() bool {
			failByApi, errStr := signalReceiver.IsFailWorkflowRequested()
			if failByApi {
				errToFailWf = provider.NewApplicationError(
					string(iwfidl.CLIENT_API_FAILING_WORKFLOW_ERROR_TYPE),
					errStr,
				)

				return true
			}
			return len(statesToExecuteQueue) > 0 || errToFailWf != nil || forceCompleteWf || stateExecutionCounter.GetTotalPendingStateExecutions() == 0 || continueAsNewCounter.IsThresholdMet()
		})
		if continueAsNewCounter.IsThresholdMet() {
			// NOTE: drain signals+thread before checking errToFailWf/forceCompleteWf so that we can close the workflow if possible
			err := continueAsNewer.DrainAllSignalsAndThreads(ctx)
			if err != nil {
				awaitError = err
			}
		}

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
		if continueAsNewCounter.IsThresholdMet() {
			// at here, all signals + threads are drained, so it's safe to continueAsNew
			return nil, continueAsNewer.ContinueToNewRun(ctx, execution, input.Config)
		}
	} // end main loop -- loop until no more state can be executed (dead end)

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
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
	}
	if shouldClose && len(decision.NextStates) > 1 {
		// Illegal decision
		err = provider.NewApplicationError(
			string(iwfidl.INVALID_USER_WORKFLOW_CODE_ERROR_TYPE),
			"invalid state decisions. Closing workflow decision cannot be combined with other state decisions",
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
	persistenceManager *PersistenceManager,
	interStateChannel *InterStateChannel,
	signalReceiver *SignalReceiver,
	timerProcessor *TimerProcessor,
	continueAsNewer *ContinueAsNewer,
	continueAsNewCounter *ContinueAsNewCounter,
) (*iwfidl.StateDecision, service.StateExecutionStatus, error) {
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
	errStartApi := provider.ExecuteActivity(ctx, StateStart, provider.GetBackendType(), service.StateStartActivityInput{
		IwfWorkerUrl: execution.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateStartRequest{
			Context:          exeCtx,
			WorkflowType:     execution.WorkflowType,
			WorkflowStateId:  state.StateId,
			StateInput:       state.StateInput,
			SearchAttributes: persistenceManager.LoadSearchAttributes(state.StateOptions),
			DataObjects:      persistenceManager.LoadDataObjects(state.StateOptions),
		},
	}).Get(ctx, &startResponse)

	if errStartApi != nil && !shouldProceedOnStartApiError(state) {
		return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, errStartApi)
	}

	err := persistenceManager.ProcessUpsertSearchAttribute(ctx, startResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	err = persistenceManager.ProcessUpsertDataObject(startResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	interStateChannel.ProcessPublishing(startResponse.GetPublishToInterStateChannel())

	commandReq := startResponse.GetCommandRequest()
	commandReqDoneOrCanceled := false

	completedTimerCmds := map[int]bool{}
	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.StartTimers(stateExeId, commandReq.GetTimerCommands())
		for idx, cmd := range commandReq.GetTimerCommands() {
			cmdCtx := provider.ExtendContextWithValue(ctx, "idx", idx)
			provider.GoNamed(cmdCtx, getThreadName("timer", cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				// Note that commandReqDoneOrCanceled is needed for two cases:
				// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
				// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
				completed := timerProcessor.WaitForTimerFiredOrSkipped(ctx, stateExeId, idx, &commandReqDoneOrCanceled)
				if completed {
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
				received := false
				_ = provider.Await(ctx, func() bool {
					received = signalReceiver.HasSignal(cmd.SignalChannelName)
					// Note that commandReqDoneOrCanceled is needed for two cases:
					// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
					// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
					return received || commandReqDoneOrCanceled
				})
				if received {
					completedSignalCmds[idx] = signalReceiver.Retrieve(cmd.SignalChannelName)
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
					// Note that commandReqDoneOrCanceled is needed for two cases:
					// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
					// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
					return received || commandReqDoneOrCanceled
				})

				if received {
					completedInterStateChannelCmds[idx] = interStateChannel.Retrieve(cmd.ChannelName)
				}
			})
		}
	}

	continueAsNewer.AddPendingStateExecutionCommandStatus(
		stateExeId,
		completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds,
		commandReq.GetTimerCommands(), commandReq.GetSignalCommands(), commandReq.GetInterStateChannelCommands(),
	)
	_ = provider.Await(ctx, func() bool {
		return IsDeciderTriggerConditionMet(provider, ctx, commandReq, completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds) || continueAsNewCounter.IsThresholdMet()
	})
	commandReqDoneOrCanceled = true
	if continueAsNewCounter.IsThresholdMet() {
		return nil, service.WaitingCommandsStateExecutionStatus, nil
	}

	commandRes := &iwfidl.CommandResults{}
	commandRes.StateStartApiSucceeded = iwfidl.PtrBool(errStartApi == nil)

	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.RemovePendingTimersOfState(stateExeId)

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
			SearchAttributes: persistenceManager.LoadSearchAttributes(state.StateOptions),
			DataObjects:      persistenceManager.LoadDataObjects(state.StateOptions),
			StateInput:       state.StateInput,
		},
	}).Get(ctx, &decideResponse)
	if err != nil {
		return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, err)
	}

	decision := decideResponse.GetStateDecision()
	err = persistenceManager.ProcessUpsertSearchAttribute(ctx, decideResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	err = persistenceManager.ProcessUpsertDataObject(decideResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	interStateChannel.ProcessPublishing(decideResponse.GetPublishToInterStateChannel())

	continueAsNewer.ClearPendingStateExecutionCommandStatus(stateExeId)

	return &decision, service.DecideApiCompletedStateExecutionStatus, nil
}

func shouldProceedOnStartApiError(state iwfidl.StateMovement) bool {
	if state.StateOptions == nil {
		return false
	}

	if state.StateOptions.StartApiFailurePolicy == nil {
		return false
	}

	return state.StateOptions.GetStartApiFailurePolicy() == iwfidl.PROCEED_TO_DECIDE_ON_START_API_FAILURE
}

func convertStateApiActivityError(provider WorkflowProvider, err error) error {
	if provider.IsApplicationError(err) {
		return err
	}
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE), err.Error())
}

func getThreadName(prefix string, cmdId string, idx int) string {
	return fmt.Sprintf("%v-%v-%v", prefix, cmdId, idx)
}
