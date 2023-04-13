package interpreter

import (
	"strings"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type SignalReceiver struct {
	// key is channel name
	receivedSignals            map[string][]*iwfidl.EncodedObject
	failWorkflowByClient       bool
	reasonFailWorkflowByClient *string
	provider                   WorkflowProvider
	timerProcessor             *TimerProcessor
	workflowConfiger           *WorkflowConfiger
	interStateChannel          *InterStateChannel
	stateRequestQueue          *StateRequestQueue
	persistenceManager         *PersistenceManager
}

func NewSignalReceiver(ctx UnifiedContext, provider WorkflowProvider, interStateChannel *InterStateChannel, stateRequestQueue *StateRequestQueue,
	persistenceManager *PersistenceManager, tp *TimerProcessor, continueAsNewCounter *ContinueAsNewCounter, workflowConfiger *WorkflowConfiger,
	initReceivedSignals map[string][]*iwfidl.EncodedObject) *SignalReceiver {
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

	provider.GoNamed(ctx, "fail-workflow-system-signal-handler", func(ctx UnifiedContext) {
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

	provider.GoNamed(ctx, "skip-timer-system-signal-handler", func(ctx UnifiedContext) {
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

	provider.GoNamed(ctx, "update-config-system-signal-handler", func(ctx UnifiedContext) {
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
				workflowConfiger.SetIfPresent(val.WorkflowConfig)
			} else {
				// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
				return
			}
		}
	})

	provider.GoNamed(ctx, "execute-rpc-signal-handler", func(ctx UnifiedContext) {
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
				_ = sr.persistenceManager.ProcessUpsertDataObject(val.UpsertDataObjects)
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

	provider.GoNamed(ctx, "user-signal-receiver-handler", func(ctx UnifiedContext) {
		for {
			var toProcess []string
			err := provider.Await(ctx, func() bool {
				unhandledSigs := provider.GetUnhandledSignalNames(ctx)

				for _, sigName := range unhandledSigs {
					if strings.HasPrefix(sigName, service.IwfSystemSignalPrefix) {
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

func (sr *SignalReceiver) receiveSignal(ctx UnifiedContext, sigName string) {
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
	sr.receivedSignals[channelName] = l
	return sigVal
}

func (sr *SignalReceiver) DumpReceived(channelNames []string) map[string][]*iwfidl.EncodedObject {
	if len(channelNames) == 0 {
		return sr.receivedSignals
	}
	data := make(map[string][]*iwfidl.EncodedObject)
	for _, n := range channelNames {
		data[n] = sr.receivedSignals[n]
	}
	return data
}

// DrainAllUnreceivedSignals will retrieve signals that after signal handler threads are stopped,
// so that the signals can be carried over to next run by continueAsNew.
// This includes both regular user signals and system signals
func (sr *SignalReceiver) DrainAllUnreceivedSignals(ctx UnifiedContext) {
	unhandledSigs := sr.provider.GetUnhandledSignalNames(ctx)
	if len(unhandledSigs) == 0 {
		return
	}

	for _, sigName := range unhandledSigs {
		if strings.HasPrefix(sigName, service.IwfSystemSignalPrefix) {
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
							sr.workflowConfiger.SetIfPresent(val.WorkflowConfig)
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
							_ = sr.persistenceManager.ProcessUpsertDataObject(val.UpsertDataObjects)
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

func (sr *SignalReceiver) IsFailWorkflowRequested() (bool, string) {
	reason := "fail by client"
	if sr.reasonFailWorkflowByClient != nil {
		reason = *sr.reasonFailWorkflowByClient
	}
	if sr.failWorkflowByClient {
		return true, reason
	} else {
		return false, ""
	}
}
