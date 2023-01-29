package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"strings"
)

type SignalReceiver struct {
	// key is channel name
	receivedSignals map[string][]*iwfidl.EncodedObject
	provider        WorkflowProvider
	logger          UnifiedLogger
}

func NewSignalReceiver(ctx UnifiedContext, provider WorkflowProvider) *SignalReceiver {
	sr := &SignalReceiver{
		provider:        provider,
		receivedSignals: map[string][]*iwfidl.EncodedObject{},
	}
	provider.GoNamed(ctx, "signal-receiver-handler", func(ctx UnifiedContext) {
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
				return len(toProcess) > 0
			})
			if err != nil {
				break
			}

			for _, sigName := range toProcess {
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

// DrainedAllSignals will wait for all signals are processed before a safe continueAsNew
func (sr *SignalReceiver) DrainedAllSignals(ctx UnifiedContext) error {
	return sr.provider.Await(ctx, func() bool {
		unhandledSigs := sr.provider.GetUnhandledSignalNames(ctx)

		for _, sigName := range unhandledSigs {
			if strings.HasPrefix(sigName, service.IwfSystemSignalPrefix) {
				if service.ValidIwfSystemSignalNames[sigName] {
					return false
				}
				// ignore invalid system signals because we can't process it
				sr.provider.GetLogger(ctx).Error("ignore the invalid system signal", sigName)
				continue
			}
			return false
		}
		return true
	})
}
