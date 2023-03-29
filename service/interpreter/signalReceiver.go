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
}

func NewSignalReceiver(ctx UnifiedContext, provider WorkflowProvider, tp *TimerProcessor, continueAsNewCounter *ContinueAsNewCounter) *SignalReceiver {
	sr := &SignalReceiver{
		provider:             provider,
		receivedSignals:      map[string][]*iwfidl.EncodedObject{},
		failWorkflowByClient: false,
	}

	provider.GoNamed(ctx, "fail-workflow-system-signal-handler", func(ctx UnifiedContext) {
		ch := provider.GetSignalChannel(ctx, service.FailWorkflowSignalChanncelName)

		val := service.FailWorkflowSignalRequest{}
		err := provider.Await(ctx, func() bool {
			sr.failWorkflowByClient = ch.ReceiveAsync(&val)
			// NOTE: continueAsNew will wait for all threads to complete, so we must stop this thread for continueAsNew when no more signals to process
			return sr.failWorkflowByClient || continueAsNewCounter.IsThresholdMet()
		})
		if err != nil {
			return
		}
		if sr.failWorkflowByClient {
			sr.reasonFailWorkflowByClient = &val.Reason
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
				var sigVal iwfidl.EncodedObject
				ch := provider.GetSignalChannel(ctx, sigName)
				ch.Receive(ctx, &sigVal)
				receivedChan := sr.receivedSignals[sigName]
				receivedChan = append(receivedChan, &sigVal)
				sr.receivedSignals[sigName] = receivedChan
			}
			toProcess = nil
		}
	})
	return sr
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

func (sr *SignalReceiver) ReadReceived(channelNames []string) map[string][]*iwfidl.EncodedObject {
	if len(channelNames) == 0 {
		return sr.receivedSignals
	}
	data := make(map[string][]*iwfidl.EncodedObject)
	for _, n := range channelNames {
		data[n] = sr.receivedSignals[n]
	}
	return data
}

// HaveAllUserAndSystemSignalsToReceive will check for if signals are received for a safe continueAsNew
// this includes both regular user signals and system signals
// Note that being received doesn't mean being processed completed. ContinueAsNew should also wait for processing the received signals properly
func (sr *SignalReceiver) HaveAllUserAndSystemSignalsToReceive(ctx UnifiedContext) bool {
	unhandledSigs := sr.provider.GetUnhandledSignalNames(ctx)
	if len(unhandledSigs) == 0 {
		return true
	}

	for _, sigName := range unhandledSigs {
		if strings.HasPrefix(sigName, service.IwfSystemSignalPrefix) {
			if service.ValidIwfSystemSignalNames[sigName] {
				// found a valid system signal, return false so that continueAsNew can wait for it
				return false
			}
			// ignore invalid system signals because we can't process it
			sr.provider.GetLogger(ctx).Error("ignore the invalid system signal", sigName)
			continue
		} else {
			// found a regular signal, return false
			return false
		}
	}
	// no unhandled system or user signals
	return true
}

func (sr *SignalReceiver) GetFailWorklowAndReasonByClient() (bool, string) {
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
