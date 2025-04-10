package interpreter

import (
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/interpreter/config"
	"github.com/indeedeng/iwf/service/interpreter/cont"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"strings"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type SignalReceiver struct {
	// key is channel name
	receivedSignals            map[string][]*iwfidl.EncodedObject
	failWorkflowByClient       bool
	reasonFailWorkflowByClient *string
	provider                   interfaces.WorkflowProvider
	timerProcessor             interfaces.TimerProcessor
	workflowConfiger           *config.WorkflowConfiger
	interStateChannel          *InternalChannel
	stateRequestQueue          *StateRequestQueue
	persistenceManager         *PersistenceManager
}

func NewSignalReceiver(
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, interStateChannel *InternalChannel,
	stateRequestQueue *StateRequestQueue,
	persistenceManager *PersistenceManager, tp interfaces.TimerProcessor, continueAsNewCounter *cont.ContinueAsNewCounter,
	workflowConfiger *config.WorkflowConfiger,
	initReceivedSignals map[string][]*iwfidl.EncodedObject,
) *SignalReceiver {
	if initReceivedSignals == nil {
		initReceivedSignals = map[string][]*iwfidl.EncodedObject{}
	}
	sr := &SignalReceiver{
		provider:             provider,
		receivedSignals:      initReceivedSignals,
		failWorkflowByClient: false,
		timerProcessor:       tp,
		workflowConfiger:     workflowConfiger,
		interStateChannel:    interStateChannel,
		stateRequestQueue:    stateRequestQueue,
		persistenceManager:   persistenceManager,
	}

	//The thread waits until a FailWorkflowSignalChannelName signal has been
	//received or a continueAsNew run is triggered. When a signal has been received it sets
	//SignalReceiver.failWorkflowByClient to true and sets SignalReceiver.reasonFailWorkflowByClient to the reason
	//given in the signal's value. If continueIsNew is triggered, the thread completes after all signals have been processed.
	provider.GoNamed(ctx, "fail-workflow-system-signal-handler", func(ctx interfaces.UnifiedContext) {
		for {
			ch := provider.GetSignalChannel(ctx, service.FailWorkflowSignalChannelName)

			val := service.FailWorkflowSignalRequest{}
			received := false
			err := provider.Await(ctx, func() bool {
				received = ch.ReceiveAsync(&val)
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				return received || continueAsNewCounter.IsThresholdMet()
			})
			if err != nil {
				break
			}
			if received {
				continueAsNewCounter.IncSignalsReceived()
				sr.failWorkflowByClient = true
				sr.reasonFailWorkflowByClient = &val.Reason
			} else {
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				break
			}
		}
	})

	//The thread waits until a SkipTimerSignalChannelName signal has been
	//received or a continueAsNew run is triggered. When a signal has been received it skips the specific timer
	//described in the signal's value. If continueIsNew is triggered, the thread completes after all signals have been processed.
	provider.GoNamed(ctx, "skip-timer-system-signal-handler", func(ctx interfaces.UnifiedContext) {
		for {
			ch := provider.GetSignalChannel(ctx, service.SkipTimerSignalChannelName)
			val := service.SkipTimerSignalRequest{}

			received := false
			err := provider.Await(ctx, func() bool {
				received = ch.ReceiveAsync(&val)
				return received || continueAsNewCounter.IsThresholdMet()
			})
			if err != nil {
				// break the loop to prevent goroutine leakage
				break
			}
			if received {
				continueAsNewCounter.IncSignalsReceived()
				tp.SkipTimer(val.StateExecutionId, val.CommandId, val.CommandIndex)
			} else {
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				return
			}
		}
	})

	//The thread waits until a UpdateConfigSignalChannelName signal has been
	//received or a continueAsNew run is triggered. When a signal has been received it updates the workflow config
	//defined in the signal's value. If continueIsNew is triggered, the thread completes after all signals have been processed.
	provider.GoNamed(ctx, "update-config-system-signal-handler", func(ctx interfaces.UnifiedContext) {
		for {
			ch := provider.GetSignalChannel(ctx, service.UpdateConfigSignalChannelName)
			val := iwfidl.WorkflowConfigUpdateRequest{}

			received := false
			err := provider.Await(ctx, func() bool {
				received = ch.ReceiveAsync(&val)
				return received || continueAsNewCounter.IsThresholdMet()
			})
			if err != nil {
				// break the loop to prevent goroutine leakage
				break
			}
			if received {
				continueAsNewCounter.IncSignalsReceived()
				workflowConfiger.UpdateByAPI(val.WorkflowConfig)
			} else {
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				return
			}
		}
	})

	//The thread waits until a TriggerContinueAsNewSignalChannelName signal has
	//been received or a continueAsNew run is triggered. When a signal has been received it triggers a continueAsNew run.
	//Since this thread is triggering a continueAsNew run it doesn't need to wait for signals to drain from the channel.
	provider.GoNamed(ctx, "trigger-continue-as-new-handler", func(ctx interfaces.UnifiedContext) {
		// NOTE: unlike other signal channels, this one doesn't need to drain during continueAsNew
		// because if there is a continueAsNew, this signal is not needed anymore
		ch := provider.GetSignalChannel(ctx, service.TriggerContinueAsNewSignalChannelName)

		received := false
		err := provider.Await(ctx, func() bool {
			received = ch.ReceiveAsync(nil)
			return received || continueAsNewCounter.IsThresholdMet()
		})
		if err != nil {
			return
		}
		if received {
			continueAsNewCounter.TriggerByAPI()
			return
		}
		return
	})

	//The thread waits until a ExecuteRpcSignalChannelName signal has been
	//received or a continueAsNew run is triggered. When a signal has been received it upserts data objects
	//(if they exist in the signal value), upserts search attributes (if they exist in the signal value),
	//and/or publishes a message to an internal channel (if InterStateChannelPublishing is set in the signal value).
	//If continueIsNew is triggered, the thread completes after all signals have been processed.
	provider.GoNamed(ctx, "execute-rpc-signal-handler", func(ctx interfaces.UnifiedContext) {
		for {
			ch := provider.GetSignalChannel(ctx, service.ExecuteRpcSignalChannelName)
			var val service.ExecuteRpcSignalRequest

			received := false
			err := provider.Await(ctx, func() bool {
				received = ch.ReceiveAsync(&val)
				return received || continueAsNewCounter.IsThresholdMet()
			})
			if err != nil {
				// break the loop to prevent goroutine leakage
				break
			}
			if received {
				continueAsNewCounter.IncSignalsReceived()
				_ = sr.persistenceManager.ProcessUpsertDataAttribute(ctx, val.UpsertDataObjects)
				_ = sr.persistenceManager.ProcessUpsertSearchAttribute(ctx, val.UpsertSearchAttributes)
				sr.interStateChannel.ProcessPublishing(val.InterStateChannelPublishing)
				if val.StateDecision != nil {
					sr.stateRequestQueue.AddStateStartRequests(val.StateDecision.NextStates)
				}
			} else {
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				return
			}
		}
	})

	//The thread waits until a signal has been received that is not an IWF
	//system signal name or a continueAsNew run is triggered. When a signal has been received it processes the
	//external signal. If continueIsNew is triggered, the thread completes after all signals have been processed.
	provider.GoNamed(ctx, "user-signal-receiver-handler", func(ctx interfaces.UnifiedContext) {
		for {
			var toProcess []string
			err := provider.Await(ctx, func() bool {
				unhandledSigs := provider.GetUnhandledSignalNames(ctx)

				for _, sigName := range unhandledSigs {
					if strings.HasPrefix(sigName, service.IwfSystemConstPrefix) {
						// skip this because it will be processed in a different thread
						if !service.ValidIwfSystemSignalNames[sigName] {
							provider.GetLogger(ctx).Error("found an invalid system signal", sigName)
						}
						continue
					}
					toProcess = append(toProcess, sigName)
				}
				return len(toProcess) > 0 || continueAsNewCounter.IsThresholdMet()
			})
			if err != nil {
				break
			}
			// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
			if len(toProcess) == 0 && continueAsNewCounter.IsThresholdMet() {
				return
			}

			for _, sigName := range toProcess {
				continueAsNewCounter.IncSignalsReceived()
				sr.receiveSignal(ctx, sigName)
			}
			toProcess = nil
		}
	})
	return sr
}

func (sr *SignalReceiver) receiveSignal(ctx interfaces.UnifiedContext, sigName string) {
	ch := sr.provider.GetSignalChannel(ctx, sigName)
	for {
		var sigVal iwfidl.EncodedObject
		ok := ch.ReceiveAsync(&sigVal)
		if ok {
			sr.receivedSignals[sigName] = append(sr.receivedSignals[sigName], &sigVal)
		} else {
			break
		}
	}
}

func (sr *SignalReceiver) HasSignal(channelName string) bool {
	l := sr.receivedSignals[channelName]
	return len(l) > 0
}

func (sr *SignalReceiver) Retrieve(channelName string) *iwfidl.EncodedObject {
	l := sr.receivedSignals[channelName]
	if len(l) <= 0 {
		panic("critical bug, this shouldn't happen")
	}
	sigVal := l[0]
	l = l[1:]
	if len(l) == 0 {
		delete(sr.receivedSignals, channelName)
	} else {
		sr.receivedSignals[channelName] = l
	}

	return sigVal
}

func (sr *SignalReceiver) GetAllReceived() map[string][]*iwfidl.EncodedObject {
	return sr.receivedSignals
}

func (sr *SignalReceiver) GetInfos() map[string]iwfidl.ChannelInfo {
	infos := make(map[string]iwfidl.ChannelInfo, len(sr.receivedSignals))
	for name, l := range sr.receivedSignals {
		infos[name] = iwfidl.ChannelInfo{
			Size: ptr.Any(int32(len(l))),
		}
	}
	return infos
}

// DrainAllReceivedButUnprocessedSignals will process all the signals that are received but not processed in the current
// workflow task.
// There are two cases this is needed:
// 1. ContinueAsNew:
// retrieve signals that after signal handler threads are stopped,
// so that the signals can be carried over to next run by continueAsNew.
// This includes both regular user signals and system signals
// 2. Conditional close/complete workflow on signal/internal channel:
// retrieve all signal/internal channel messages before checking the signal/internal channels
func (sr *SignalReceiver) DrainAllReceivedButUnprocessedSignals(ctx interfaces.UnifiedContext) {
	unhandledSigs := sr.provider.GetUnhandledSignalNames(ctx)
	if len(unhandledSigs) == 0 {
		return
	}

	for _, sigName := range unhandledSigs {
		if strings.HasPrefix(sigName, service.IwfSystemConstPrefix) {
			if service.ValidIwfSystemSignalNames[sigName] {

				sr.provider.GetLogger(ctx).Info("found a valid system signal before continueAsNew to carry over", sigName)
				if sigName == service.SkipTimerSignalChannelName {
					ch := sr.provider.GetSignalChannel(ctx, service.SkipTimerSignalChannelName)
					for {
						val := service.SkipTimerSignalRequest{}
						ok := ch.ReceiveAsync(&val)
						if ok {
							sr.timerProcessor.SkipTimer(val.StateExecutionId, val.CommandId, val.CommandIndex)
						} else {
							break
						}
					}
				} else if sigName == service.UpdateConfigSignalChannelName {
					ch := sr.provider.GetSignalChannel(ctx, service.UpdateConfigSignalChannelName)
					for {
						val := iwfidl.WorkflowConfigUpdateRequest{}
						ok := ch.ReceiveAsync(&val)
						if ok {
							sr.workflowConfiger.UpdateByAPI(val.WorkflowConfig)
						} else {
							break
						}
					}
				} else if sigName == service.FailWorkflowSignalChannelName {
					ch := sr.provider.GetSignalChannel(ctx, service.FailWorkflowSignalChannelName)
					for {
						val := service.FailWorkflowSignalRequest{}
						ok := ch.ReceiveAsync(&val)
						if ok {
							sr.failWorkflowByClient = true
							sr.reasonFailWorkflowByClient = &val.Reason
						} else {
							break
						}
					}
				} else if sigName == service.ExecuteRpcSignalChannelName {
					ch := sr.provider.GetSignalChannel(ctx, service.ExecuteRpcSignalChannelName)
					for {
						val := service.ExecuteRpcSignalRequest{}
						ok := ch.ReceiveAsync(&val)
						if ok {
							_ = sr.persistenceManager.ProcessUpsertDataAttribute(ctx, val.UpsertDataObjects)
							_ = sr.persistenceManager.ProcessUpsertSearchAttribute(ctx, val.UpsertSearchAttributes)
							sr.interStateChannel.ProcessPublishing(val.InterStateChannelPublishing)
							if val.StateDecision != nil {
								sr.stateRequestQueue.AddStateStartRequests(val.StateDecision.NextStates)
							}
						} else {
							break
						}
					}
				}
				continue
			}
			// ignore invalid system signals because we can't process it
			sr.provider.GetLogger(ctx).Error("ignore the invalid system signal", sigName)
			continue
		} else {
			sr.provider.GetLogger(ctx).Info("found a valid user signal before continueAsNew to carry over", sigName)
			sr.receiveSignal(ctx, sigName)
			continue
		}
	}
}

func (sr *SignalReceiver) IsFailWorkflowRequested() (bool, error) {
	reason := "fail by client"
	if sr.reasonFailWorkflowByClient != nil {
		reason = *sr.reasonFailWorkflowByClient
	}
	if sr.failWorkflowByClient {
		return true, sr.provider.NewApplicationError(
			string(iwfidl.CLIENT_API_FAILING_WORKFLOW_ERROR_TYPE),
			reason,
		)
	} else {
		return false, nil
	}
}
