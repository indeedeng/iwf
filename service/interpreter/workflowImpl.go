package interpreter

import (
	"context"
	"fmt"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"time"

	"github.com/indeedeng/iwf/service/common/compatibility"
	"golang.org/x/exp/slices"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

func InterpreterImpl(
	ctx UnifiedContext, provider WorkflowProvider, input service.InterpreterWorkflowInput,
) (*service.InterpreterWorkflowOutput, error) {
	var err error
	globalVersioner := NewGlobalVersioner(provider, ctx)
	if globalVersioner.IsAfterVersionOfUsingGlobalVersioning() {
		err = globalVersioner.UpsertGlobalVersionSearchAttribute()
		if err != nil {
			return nil, err
		}
	}

	if !input.Config.GetDisableSystemSearchAttribute() {
		if !globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() {
			// we have stopped upsert here in new versions, because it's done in start workflow request
			err = provider.UpsertSearchAttributes(ctx, map[string]interface{}{
				service.SearchAttributeIwfWorkflowType: input.IwfWorkflowType,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	workflowConfiger := NewWorkflowConfiger(input.Config)
	basicInfo := service.BasicInfo{
		IwfWorkflowType: input.IwfWorkflowType,
		IwfWorkerUrl:    input.IwfWorkerUrl,
	}

	var interStateChannel *InterStateChannel
	var stateRequestQueue *StateRequestQueue
	var persistenceManager *PersistenceManager
	var timerProcessor *TimerProcessor
	var continueAsNewCounter *ContinueAsNewCounter
	var signalReceiver *SignalReceiver
	var stateExecutionCounter *StateExecutionCounter
	var outputCollector *OutputCollector
	var continueAsNewer *ContinueAsNewer
	if input.IsResumeFromContinueAsNew {
		canInput := input.ContinueAsNewInput
		config := workflowConfiger.Get()
		previous, err := LoadInternalsFromPreviousRun(ctx, provider, canInput.PreviousInternalRunId, config.GetContinueAsNewPageSizeInBytes())
		if err != nil {
			return nil, err
		}

		// The below initialization order should be the same as for non-continueAsNew

		interStateChannel = RebuildInterStateChannel(previous.InterStateChannelReceived)
		stateRequestQueue = NewStateRequestQueueWithResumeRequests(previous.StatesToStartFromBeginning, previous.StateExecutionsToResume)
		persistenceManager = RebuildPersistenceManager(provider, previous.DataObjects, previous.SearchAttributes, input.UseMemoForDataAttributes)
		timerProcessor = NewTimerProcessor(ctx, provider, previous.StaleSkipTimerSignals)
		continueAsNewCounter = NewContinueAsCounter(workflowConfiger, ctx, provider)
		signalReceiver = NewSignalReceiver(ctx, provider, interStateChannel, stateRequestQueue, persistenceManager, timerProcessor, continueAsNewCounter, workflowConfiger, previous.SignalsReceived)
		counterInfo := previous.StateExecutionCounterInfo
		stateExecutionCounter = RebuildStateExecutionCounter(ctx, provider,
			counterInfo.StateIdStartedCount, counterInfo.StateIdCurrentlyExecutingCount, counterInfo.TotalCurrentlyExecutingCount,
			workflowConfiger, continueAsNewCounter)
		outputCollector = NewOutputCollector(previous.StateOutputs)
		continueAsNewer = NewContinueAsNewer(provider, interStateChannel, signalReceiver, stateExecutionCounter, persistenceManager, stateRequestQueue, outputCollector, timerProcessor)
	} else {
		interStateChannel = NewInterStateChannel()
		stateRequestQueue = NewStateRequestQueue()
		persistenceManager = NewPersistenceManager(provider, input.InitSearchAttributes, input.UseMemoForDataAttributes)
		timerProcessor = NewTimerProcessor(ctx, provider, nil)
		continueAsNewCounter = NewContinueAsCounter(workflowConfiger, ctx, provider)
		signalReceiver = NewSignalReceiver(ctx, provider, interStateChannel, stateRequestQueue, persistenceManager, timerProcessor, continueAsNewCounter, workflowConfiger, nil)
		stateExecutionCounter = NewStateExecutionCounter(ctx, provider, workflowConfiger, continueAsNewCounter)
		outputCollector = NewOutputCollector(nil)
		continueAsNewer = NewContinueAsNewer(provider, interStateChannel, signalReceiver, stateExecutionCounter, persistenceManager, stateRequestQueue, outputCollector, timerProcessor)
	}

	_, err = NewWorkflowUpdater(ctx, provider, persistenceManager, stateRequestQueue, continueAsNewer, continueAsNewCounter, interStateChannel, basicInfo)
	if err != nil {
		return nil, err
	}
	err = SetQueryHandlers(ctx, provider, persistenceManager, continueAsNewer, workflowConfiger, basicInfo)
	if err != nil {
		return nil, err
	}

	var errToFailWf error // Note that today different errors could overwrite each other, we only support last one wins. we may use multiError to improve.
	var forceCompleteWf bool
	var shouldGracefulComplete bool

	// this is for an optimization for StateId Search attribute, see updateStateIdSearchAttribute in stateExecutionCounter
	// Because it will check totalCurrentlyExecutingCount == 0, so it will also work for continueAsNew case
	defer stateExecutionCounter.ClearExecutingStateIdsSearchAttributeFinally()

	if !input.IsResumeFromContinueAsNew {
		// it's possible that a workflow is started without any starting state
		// it will wait for a new state coming in (by RPC results)
		if input.StartStateId != nil {
			startingState := iwfidl.StateMovement{
				StateId:      *input.StartStateId,
				StateOptions: input.StateOptions,
				StateInput:   input.StateInput,
			}
			stateRequestQueue.AddStateStartRequests([]iwfidl.StateMovement{startingState})
		}
	}

	for {
		err = provider.Await(ctx, func() bool {
			failWorkflowByClient, _ := signalReceiver.IsFailWorkflowRequested()
			if globalVersioner.IsAfterVersionOfContinueAsNewOnNoStates() {
				return !stateRequestQueue.IsEmpty() || failWorkflowByClient || shouldGracefulComplete || continueAsNewCounter.IsThresholdMet()
			}
			// below was a bug in the older version that workflow didn't continue as new
			// but have to keep workflow deterministic
			return !stateRequestQueue.IsEmpty() || failWorkflowByClient || shouldGracefulComplete
		})
		if err != nil {
			return nil, err
		}
		failWorkflowByClient, failErr := signalReceiver.IsFailWorkflowRequested()
		if failWorkflowByClient {
			return nil, failErr
		}
		if shouldGracefulComplete && stateRequestQueue.IsEmpty() {
			break
		}

		for !stateRequestQueue.IsEmpty() {

			var statesToExecute []StateRequest
			if !continueAsNewCounter.IsThresholdMet() {
				statesToExecute = stateRequestQueue.TakeAll()
				err = stateExecutionCounter.MarkStateIdExecutingIfNotYet(statesToExecute)
				if err != nil {
					return nil, err
				}
			}

			for _, stateReqForLoopingOnly := range statesToExecute {
				// execute in another thread for parallelism
				// state must be passed via parameter https://stackoverflow.com/questions/67263092
				stateCtx := provider.ExtendContextWithValue(ctx, "stateReq", stateReqForLoopingOnly)
				provider.GoNamed(stateCtx, "state-execution-thread:"+stateReqForLoopingOnly.GetStateId(), func(ctx UnifiedContext) {
					stateReq, ok := provider.GetContextValue(ctx, "stateReq").(StateRequest)
					if !ok {
						errToFailWf = provider.NewApplicationError(
							string(iwfidl.SERVER_INTERNAL_ERROR_TYPE),
							"critical code bug when passing state request via context",
						)
						return
					}

					var state iwfidl.StateMovement
					var stateExeId string
					if stateReq.IsResumeRequest() {
						resumeReq := stateReq.GetStateResumeRequest()
						state = resumeReq.State
						stateExeId = resumeReq.StateExecutionId
					} else {
						state = stateReq.GetStateStartRequest()
						stateExeId = stateExecutionCounter.CreateNextExecutionId(state.GetStateId())
					}

					shouldSendSignalOnCompletion := slices.Contains(input.WaitForCompletionStateExecutionIds, stateExeId)

					decision, stateExecStatus, err := executeState(
						ctx, provider, basicInfo, stateReq, stateExeId, persistenceManager, interStateChannel,
						signalReceiver, timerProcessor, continueAsNewer, continueAsNewCounter, shouldSendSignalOnCompletion)
					if err != nil {
						// this is the case where stateExecStatus == FailureStateExecutionStatus
						errToFailWf = err
						// state execution fail should fail the workflow, no more processing
						return
					}

					if stateExecStatus == service.CompletedStateExecutionStatus {
						// NOTE: decision is only available on this CompletedStateExecutionStatus

						canGoNext, gracefulComplete, forceComplete, forceFail, output, err :=
							checkClosingWorkflow(ctx, provider, decision, state.GetStateId(), stateExeId, interStateChannel, signalReceiver)
						if err != nil {
							errToFailWf = err
							// no return so that it can fall through to call MarkStateExecutionCompleted
						}
						if gracefulComplete {
							shouldGracefulComplete = true
						}
						if (gracefulComplete || forceComplete || forceFail) && output != nil {
							outputCollector.Add(*output)
						}
						if forceComplete {
							forceCompleteWf = true
						}
						if forceFail {
							errToFailWf = provider.NewApplicationError(
								string(iwfidl.STATE_DECISION_FAILING_WORKFLOW_ERROR_TYPE),
								outputCollector.GetAll(),
							)
							// no return so that it can fall through to call MarkStateExecutionCompleted
						}
						if canGoNext && decision.HasNextStates() {
							stateRequestQueue.AddStateStartRequests(decision.GetNextStates())
						}

						// finally, mark state completed and may also update system search attribute(IwfExecutingStateIds)
						err = stateExecutionCounter.MarkStateExecutionCompleted(state)
						if err != nil {
							errToFailWf = err
						}
					} else if stateExecStatus == service.ExecuteApiFailedAndProceed {
						options := state.GetStateOptions()
						stateRequestQueue.AddSingleStateStartRequest(options.GetExecuteApiFailureProceedStateId(), state.StateInput, options.ExecuteApiFailureProceedStateOptions)
						// finally, mark state completed and may also update system search attribute(IwfExecutingStateIds)
						err = stateExecutionCounter.MarkStateExecutionCompleted(state)
						if err != nil {
							errToFailWf = err
						}
					}
					// noop for WaitingCommandsStateExecutionStatus, because it means continueAsNew
				}) // end of executing one state
			} // end loop of executing all states from the queue for one iteration

			// The conditions here are quite tricky:
			// For !stateRequestQueue.IsEmpty(): We need some condition to wait here because all the state execution are running in different thread.
			//    Right after the queue are popped it becomes empty. When it's not empty, it means there are new states to execute pushed into the queue,
			//    and it's time to wake up the outer loop to go to next iteration. Alternatively, waiting for all current started in this iteration to complete will also work,
			//    but not as efficient as this one because it will take much longer time.
			// For errToFailWf != nil || forceCompleteWf: this means we need to close workflow immediately
			// For stateExecutionCounter.GetTotalCurrentlyExecutingCount() == 0: this means all the state executions have reach "Dead Ends" so the workflow can complete gracefully without output
			// For continueAsNewCounter.IsThresholdMet(): this means workflow need to continueAsNew
			awaitError := provider.Await(ctx, func() bool {
				failByApi, failErr := signalReceiver.IsFailWorkflowRequested()
				if failByApi {
					errToFailWf = failErr
					return true
				}
				return !stateRequestQueue.IsEmpty() || errToFailWf != nil || forceCompleteWf || stateExecutionCounter.GetTotalCurrentlyExecutingCount() == 0 || continueAsNewCounter.IsThresholdMet()
			})
			if continueAsNewCounter.IsThresholdMet() {
				// NOTE: drain thread before checking errToFailWf/forceCompleteWf so that we can close the workflow if possible
				err := continueAsNewer.DrainThreads(ctx)
				if err != nil {
					awaitError = err
				}
			}

			if errToFailWf != nil || forceCompleteWf {
				return &service.InterpreterWorkflowOutput{
					StateCompletionOutputs: outputCollector.GetAll(),
				}, errToFailWf
			}

			if awaitError != nil {
				// this could happen for cancellation
				errToFailWf = awaitError
				break
			}
			if continueAsNewCounter.IsThresholdMet() {
				// the outer logic will do the actual continue as new
				break
			}
		} // end loop until no more state can be executed (dead end)

		if continueAsNewCounter.IsThresholdMet() {
			// we have to drain this again because this can be from non-state cases
			err := continueAsNewer.DrainThreads(ctx)
			if err != nil {
				errToFailWf = err
				break
			}

			// NOTE: This must be the last thing before continueAsNew!!!
			// Otherwise, there could be signals unhandled
			signalReceiver.DrainAllUnreceivedSignals(ctx)

			// after draining signals, there could be some changes
			// last fail workflow signal, return the workflow so that we don't carry over the fail request
			failByApi, failErr := signalReceiver.IsFailWorkflowRequested()
			if failByApi {
				return &service.InterpreterWorkflowOutput{
					StateCompletionOutputs: outputCollector.GetAll(),
				}, failErr
			}
			if stateRequestQueue.IsEmpty() && !continueAsNewer.HasAnyStateExecutionToResume() && shouldGracefulComplete {
				// if it is empty and no stateExecutionsToResume and request a graceful complete just complete the loop
				// so that we don't carry over shouldGracefulComplete
				break
			}
			// last update config, do it here because we use input to carry over config, not continueAsNewer query
			input.Config = workflowConfiger.Get() // update config to the lastest before continueAsNew to carry over
			input.IsResumeFromContinueAsNew = true
			input.ContinueAsNewInput = &service.ContinueAsNewInput{
				PreviousInternalRunId: provider.GetWorkflowInfo(ctx).WorkflowExecution.RunID,
			}
			// nix the unused data
			input.StateInput = nil
			input.StateOptions = nil
			input.StartStateId = nil
			return nil, provider.NewInterpreterContinueAsNewError(ctx, input)
		}
	} // end main loop

	// gracefully complete workflow when all states are executed to dead ends
	return &service.InterpreterWorkflowOutput{
		StateCompletionOutputs: outputCollector.GetAll(),
	}, errToFailWf
}

func checkClosingWorkflow(
	ctx UnifiedContext, provider WorkflowProvider, decision *iwfidl.StateDecision,
	currentStateId, currentStateExeId string,
	internalChannel *InterStateChannel, signalReceiver *SignalReceiver,
) (canGoNext, gracefulComplete, forceComplete, forceFail bool, completeOutput *iwfidl.StateCompletionOutput, err error) {
	if decision.HasConditionalClose() {
		conditionClose := decision.ConditionalClose
		if conditionClose.GetConditionalCloseType() == iwfidl.FORCE_COMPLETE_ON_INTERNAL_CHANNEL_EMPTY ||
			conditionClose.GetConditionalCloseType() == iwfidl.FORCE_COMPLETE_ON_SIGNAL_CHANNEL_EMPTY {
			// trigger a signal draining so that all the signal/internal channel messages are processed
			// TODO https://github.com/indeedeng/iwf/issues/289
			// https://github.com/indeedeng/iwf/issues/290
			// if a messages from internal channel is published via State execution,
			// we don't do any draining here yet, so the conditional completion could still lose the messages
			// to workaround, user code will have to use persistence locking
			signalReceiver.DrainAllUnreceivedSignals(ctx)

			conditionMet := false
			if conditionClose.GetConditionalCloseType() == iwfidl.FORCE_COMPLETE_ON_INTERNAL_CHANNEL_EMPTY &&
				!internalChannel.HasData(conditionClose.GetChannelName()) {
				conditionMet = true
			}
			if conditionClose.GetConditionalCloseType() == iwfidl.FORCE_COMPLETE_ON_SIGNAL_CHANNEL_EMPTY &&
				!signalReceiver.HasSignal(conditionClose.GetChannelName()) {
				conditionMet = true
			}

			if conditionMet {
				// condition is met, force complete the workflow
				forceComplete = true
				completeOutput = &iwfidl.StateCompletionOutput{
					CompletedStateId:          currentStateId,
					CompletedStateExecutionId: currentStateExeId,
					CompletedStateOutput:      conditionClose.CloseInput,
				}
				return
			} else {
				for _, st := range decision.GetNextStates() {
					if service.ValidClosingWorkflowStateId[st.GetStateId()] {
						err = createUserWorkflowError(provider, "invalid ConditionUnmetDecision with stateId: "+st.GetStateId())
						return
					}
				}

				canGoNext = true
				return
			}
		} else {
			msg := "invalid state decisions. Unsupported ConditionalCloseType " + string(conditionClose.GetConditionalCloseType())
			err = createUserWorkflowError(provider, msg)
			return
		}
	}

	canGoNext = true
	systemStateId := false
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.GracefulCompletingWorkflowStateId {
			canGoNext = false
			gracefulComplete = true
			systemStateId = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceCompletingWorkflowStateId {
			canGoNext = false
			forceComplete = true
			systemStateId = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.ForceFailingWorkflowStateId {
			canGoNext = false
			forceFail = true
			systemStateId = true
			completeOutput = &iwfidl.StateCompletionOutput{
				CompletedStateId:          currentStateId,
				CompletedStateExecutionId: currentStateExeId,
				CompletedStateOutput:      movement.StateInput,
			}
		}
		if stateId == service.DeadEndWorkflowStateId {
			canGoNext = false
			systemStateId = true
		}
	}
	if len(decision.GetNextStates()) == 0 {
		// legacy to keep compatibility for old code that use empty decision as graceful complete
		gracefulComplete = true
		canGoNext = false
	}
	if systemStateId && len(decision.NextStates) > 1 {
		// Illegal decision
		err = createUserWorkflowError(provider, "invalid state decisions. Closing workflow decision cannot be combined with other state decisions")
		return
	}
	return
}

func executeState(
	ctx UnifiedContext,
	provider WorkflowProvider,
	basicInfo service.BasicInfo,
	stateReq StateRequest,
	stateExeId string,
	persistenceManager *PersistenceManager,
	interStateChannel *InterStateChannel,
	signalReceiver *SignalReceiver,
	timerProcessor *TimerProcessor,
	continueAsNewer *ContinueAsNewer,
	continueAsNewCounter *ContinueAsNewCounter,
	shouldSendSignalOnCompletion bool,
) (*iwfidl.StateDecision, service.StateExecutionStatus, error) {
	globalVersioner := NewGlobalVersioner(provider, ctx)
	waitUntilApi := StateStart
	executeApi := StateDecide
	if globalVersioner.IsAfterVersionOfRenamedStateApi() {
		waitUntilApi = StateApiWaitUntil
		executeApi = StateApiExecute
	}

	info := provider.GetWorkflowInfo(ctx) // TODO use firstRunId instead
	executionContext := iwfidl.Context{
		WorkflowId:               info.WorkflowExecution.ID,
		WorkflowRunId:            info.WorkflowExecution.RunID,
		WorkflowStartedTimestamp: info.WorkflowStartTime.Unix(),
		StateExecutionId:         &stateExeId,
	}
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}

	var errStartApi error
	var startResponse *iwfidl.WorkflowStateStartResponse
	var stateExecutionLocal []iwfidl.KeyValue
	var commandReq iwfidl.CommandRequest
	commandReqDoneOrCanceled := false
	completedTimerCmds := map[int]service.InternalTimerStatus{}
	completedSignalCmds := map[int]*iwfidl.EncodedObject{}
	completedInterStateChannelCmds := map[int]*iwfidl.EncodedObject{}

	state := stateReq.GetStateMovement()
	isResumeFromContinueAsNew := stateReq.IsResumeRequest()

	options := state.GetStateOptions()
	skipStart := compatibility.GetSkipStartApi(&options)
	if skipStart {
		return executeStateDecide(ctx, provider, basicInfo, state, stateExeId, persistenceManager, interStateChannel, executionContext,
			nil, continueAsNewer, executeApi, stateExecutionLocal, shouldSendSignalOnCompletion)
	}

	if isResumeFromContinueAsNew {
		resumeStateRequest := stateReq.GetStateResumeRequest()
		stateExecutionLocal = resumeStateRequest.StateExecutionLocals
		commandReq = resumeStateRequest.CommandRequest
		completedCmds := resumeStateRequest.StateExecutionCompletedCommands
		completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds = completedCmds.CompletedTimerCommands, completedCmds.CompletedSignalCommands, completedCmds.CompletedInterStateChannelCommands
	} else {
		if state.StateOptions != nil {
			startApiTimeout := compatibility.GetStartApiTimeoutSeconds(state.StateOptions)
			if startApiTimeout > 0 {
				activityOptions.StartToCloseTimeout = time.Duration(startApiTimeout) * time.Second
			}
			activityOptions.RetryPolicy = compatibility.GetStartApiRetryPolicy(state.StateOptions)
		}

		ctx = provider.WithActivityOptions(ctx, activityOptions)

		saLoadingPolicy := state.GetStateOptions().SearchAttributesLoadingPolicy
		doLoadingPolicy := compatibility.GetDataObjectsLoadingPolicy(state.StateOptions)

		errStartApi = provider.ExecuteActivity(ctx, waitUntilApi, provider.GetBackendType(), service.StateStartActivityInput{
			IwfWorkerUrl: basicInfo.IwfWorkerUrl,
			Request: iwfidl.WorkflowStateStartRequest{
				Context:          executionContext,
				WorkflowType:     basicInfo.IwfWorkflowType,
				WorkflowStateId:  state.StateId,
				StateInput:       state.StateInput,
				SearchAttributes: persistenceManager.LoadSearchAttributes(ctx, saLoadingPolicy),
				DataObjects:      persistenceManager.LoadDataObjects(ctx, doLoadingPolicy),
			},
		}).Get(ctx, &startResponse)
		persistenceManager.UnlockPersistence(saLoadingPolicy, doLoadingPolicy)
		if errStartApi != nil && !shouldProceedOnStartApiError(state) {
			return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, errStartApi)
		}

		err := persistenceManager.ProcessUpsertSearchAttribute(ctx, startResponse.GetUpsertSearchAttributes())
		if err != nil {
			return nil, service.FailureStateExecutionStatus, err
		}
		err = persistenceManager.ProcessUpsertDataObject(ctx, startResponse.GetUpsertDataObjects())
		if err != nil {
			return nil, service.FailureStateExecutionStatus, err
		}
		interStateChannel.ProcessPublishing(startResponse.GetPublishToInterStateChannel())

		commandReq = startResponse.GetCommandRequest()
		stateExecutionLocal = startResponse.GetUpsertStateLocals()
	}

	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.AddTimers(stateExeId, commandReq.GetTimerCommands(), completedTimerCmds)
		for idx, cmd := range commandReq.GetTimerCommands() {
			if _, ok := completedTimerCmds[idx]; ok {
				// skip the completed timers(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("timer", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
				idx, ok := provider.GetContextValue(ctx, "idx").(int)
				if !ok {
					panic("critical code bug")
				}

				// Note that commandReqDoneOrCanceled is needed for two cases:
				// 1. will be true when trigger type of the commandReq is completed(e.g. AnyCommandCompleted) so we don't need to wait for all commands. Returning the thread to avoid thread leakage.
				// 2. will be true to cancel the wait for unblocking continueAsNew(continueAsNew will wait for all threads to complete)
				status := timerProcessor.WaitForTimerFiredOrSkipped(ctx, stateExeId, idx, &commandReqDoneOrCanceled)
				if status == service.TimerSkipped || status == service.TimerFired {
					completedTimerCmds[idx] = status
				}
			})
		}
	}

	if len(commandReq.GetSignalCommands()) > 0 {
		for idx, cmd := range commandReq.GetSignalCommands() {
			if _, ok := completedSignalCmds[idx]; ok {
				// skip completed signal(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("signal", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
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

	if len(commandReq.GetInterStateChannelCommands()) > 0 {
		for idx, cmd := range commandReq.GetInterStateChannelCommands() {
			if _, ok := completedInterStateChannelCmds[idx]; ok {
				// skip completed interStateChannelCommand(from continueAsNew)
				continue
			}
			cmdCtx := provider.ExtendContextWithValue(ctx, "cmd", cmd)
			cmdCtx = provider.ExtendContextWithValue(cmdCtx, "idx", idx)
			provider.GoNamed(cmdCtx, getCommandThreadName("interstate", stateExeId, cmd.GetCommandId(), idx), func(ctx UnifiedContext) {
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

	continueAsNewer.AddPotentialStateExecutionToResume(
		stateExeId, state, stateExecutionLocal, commandReq,
		completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds,
	)
	_ = provider.Await(ctx, func() bool {
		return IsDeciderTriggerConditionMet(commandReq, completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds) || continueAsNewCounter.IsThresholdMet()
	})
	commandReqDoneOrCanceled = true
	if !IsDeciderTriggerConditionMet(commandReq, completedTimerCmds, completedSignalCmds, completedInterStateChannelCmds) {
		// this means continueAsNewCounter.IsThresholdMet == true
		// not using continueAsNewCounter.IsThresholdMet because deciderTrigger is higher prioritized
		// it won't continueAsNew in those cases 1. start Api fail with proceed policy, 2. empty commands, 3. both commands and continueAsNew are met
		return nil, service.WaitingCommandsStateExecutionStatus, nil
	}

	commandRes := &iwfidl.CommandResults{}
	commandRes.StateStartApiSucceeded = iwfidl.PtrBool(errStartApi == nil)

	if len(commandReq.GetTimerCommands()) > 0 {
		timerProcessor.RemovePendingTimersOfState(stateExeId)

		var timerResults []iwfidl.TimerResult
		for idx, cmd := range commandReq.GetTimerCommands() {
			status := iwfidl.SCHEDULED
			if _, ok := completedTimerCmds[idx]; ok {
				// TODO expose skipped status to external
				status = iwfidl.FIRED
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

	return executeStateDecide(ctx, provider, basicInfo, state, stateExeId, persistenceManager, interStateChannel, executionContext,
		commandRes, continueAsNewer, executeApi, stateExecutionLocal, shouldSendSignalOnCompletion)
}
func executeStateDecide(
	ctx UnifiedContext,
	provider WorkflowProvider,
	basicInfo service.BasicInfo,
	state iwfidl.StateMovement,
	stateExeId string,
	persistenceManager *PersistenceManager,
	interStateChannel *InterStateChannel,
	executionContext iwfidl.Context,
	commandRes *iwfidl.CommandResults,
	continueAsNewer *ContinueAsNewer,
	executeApi interface{},
	stateExecutionLocal []iwfidl.KeyValue,
	shouldSendSignalOnCompletion bool,
) (*iwfidl.StateDecision, service.StateExecutionStatus, error) {
	var err error
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 30 * time.Second,
	}
	if state.StateOptions != nil {
		decideApiTimeout := compatibility.GetDecideApiTimeoutSeconds(state.StateOptions)
		if decideApiTimeout > 0 {
			activityOptions.StartToCloseTimeout = time.Duration(decideApiTimeout) * time.Second
		}
		activityOptions.RetryPolicy = compatibility.GetDecideApiRetryPolicy(state.StateOptions)
	}

	saLoadingPolicy := state.GetStateOptions().SearchAttributesLoadingPolicy
	doLoadingPolicy := compatibility.GetDataObjectsLoadingPolicy(state.StateOptions)

	ctx = provider.WithActivityOptions(ctx, activityOptions)
	var decideResponse *iwfidl.WorkflowStateDecideResponse
	err = provider.ExecuteActivity(ctx, executeApi, provider.GetBackendType(), service.StateDecideActivityInput{
		IwfWorkerUrl: basicInfo.IwfWorkerUrl,
		Request: iwfidl.WorkflowStateDecideRequest{
			Context:          executionContext,
			WorkflowType:     basicInfo.IwfWorkflowType,
			WorkflowStateId:  state.StateId,
			CommandResults:   commandRes,
			StateLocals:      stateExecutionLocal,
			SearchAttributes: persistenceManager.LoadSearchAttributes(ctx, saLoadingPolicy),
			DataObjects:      persistenceManager.LoadDataObjects(ctx, doLoadingPolicy),
			StateInput:       state.StateInput,
		},
	}, false, 0).Get(ctx, &decideResponse)
	persistenceManager.UnlockPersistence(saLoadingPolicy, doLoadingPolicy)
	if err == nil && shouldSendSignalOnCompletion && !provider.IsReplaying(ctx) {
		// NOTE: here uses NOT IsReplaying to signalWithStart, to save an activity for this operation
		// this is not a problem because the signalWithStart will be very fast and highly available
		unifiedClient := env.GetUnifiedClient()
		err := unifiedClient.SignalWithStartWaitForStateCompletionWorkflow(
			context.Background(),
			uclient.StartWorkflowOptions{
				ID:                       service.IwfSystemConstPrefix + executionContext.WorkflowId + "_" + *executionContext.StateExecutionId,
				TaskQueue:                env.GetTaskQueue(),
				WorkflowExecutionTimeout: 60 * time.Second, // timeout doesn't matter here as it will complete immediate with the signal
			},
			iwfidl.StateCompletionOutput{
				CompletedStateExecutionId: *executionContext.StateExecutionId,
			})
		if err != nil && !unifiedClient.IsWorkflowAlreadyStartedError(err) {
			// WorkflowAlreadyStartedError is returned when the started workflow is closed and the signal is not sent
			// panic will let the workflow task will retry until the signal is sent
			panic(fmt.Errorf("failed to signal on completion %w", err))
		}
	}
	if err != nil {
		if shouldProceedOnExecuteApiError(state) {
			return nil, service.ExecuteApiFailedAndProceed, nil
		}
		return nil, service.FailureStateExecutionStatus, convertStateApiActivityError(provider, err)
	}

	err = persistenceManager.ProcessUpsertSearchAttribute(ctx, decideResponse.GetUpsertSearchAttributes())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	err = persistenceManager.ProcessUpsertDataObject(ctx, decideResponse.GetUpsertDataObjects())
	if err != nil {
		return nil, service.FailureStateExecutionStatus, err
	}
	interStateChannel.ProcessPublishing(decideResponse.GetPublishToInterStateChannel())

	continueAsNewer.RemoveStateExecutionToResume(stateExeId)

	decision := decideResponse.GetStateDecision()
	return &decision, service.CompletedStateExecutionStatus, nil
}

func shouldProceedOnStartApiError(state iwfidl.StateMovement) bool {
	if state.StateOptions == nil {
		return false
	}

	policy := compatibility.GetStartApiFailurePolicy(state.StateOptions)
	if policy == nil {
		return false
	}

	return *policy == iwfidl.PROCEED_TO_DECIDE_ON_START_API_FAILURE
}

func shouldProceedOnExecuteApiError(state iwfidl.StateMovement) bool {
	if state.StateOptions == nil {
		return false
	}

	options := state.GetStateOptions()
	return options.GetExecuteApiFailureProceedStateId() != "" &&
		options.GetExecuteApiFailurePolicy() == iwfidl.PROCEED_TO_CONFIGURED_STATE
}

func convertStateApiActivityError(provider WorkflowProvider, err error) error {
	if provider.IsApplicationError(err) {
		return err
	}
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE), err.Error())
}

func getCommandThreadName(prefix string, stateExecId, cmdId string, idx int) string {
	return fmt.Sprintf("%v-%v-%v-%v", prefix, stateExecId, cmdId, idx)
}

func createUserWorkflowError(provider WorkflowProvider, message string) error {
	return provider.NewApplicationError(
		string(iwfidl.INVALID_USER_WORKFLOW_CODE_ERROR_TYPE),
		message,
	)
}

func WaitForStateCompletionWorkflowImpl(
	ctx UnifiedContext, provider WorkflowProvider,
) (*service.WaitForStateCompletionWorkflowOutput, error) {
	signalReceiveChannel := provider.GetSignalChannel(ctx, service.StateCompletionSignalChannelName)
	var signalValue iwfidl.StateCompletionOutput
	signalReceiveChannel.ReceiveBlocking(ctx, &signalValue)

	return &service.WaitForStateCompletionWorkflowOutput{
		StateCompletionOutput: signalValue,
	}, nil
}
