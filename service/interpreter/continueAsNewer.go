package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type ContinueAsNewer struct {
	stateExecutionsCompletedCommands map[string]service.StateExecutionCompletedCommands
	interStateChannel                *InterStateChannel
	stateExecutionCounter            *StateExecutionCounter
}

func NewContinueAsNewer(interStateChannel *InterStateChannel, stateExecutionCounter *StateExecutionCounter) *ContinueAsNewer {
	return &ContinueAsNewer{
		interStateChannel:                interStateChannel,
		stateExecutionCounter:            stateExecutionCounter,
		stateExecutionsCompletedCommands: map[string]service.StateExecutionCompletedCommands{},
	}
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext, provider WorkflowProvider) error {
	err := provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func(request service.DumpAllInternalRequest) (*service.DumpAllInternalResponse, error) {
		// TODO use request for pagination
		return &service.DumpAllInternalResponse{
			InterStateChannelReceived:        c.interStateChannel.ReadReceived(nil),
			StateExecutionCounterInfo:        c.stateExecutionCounter.Dump(),
			StateExecutionsCompletedCommands: c.stateExecutionsCompletedCommands,
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ContinueAsNewer) AddStateExecutionCompletedCommands(
	stateExecutionId string, completedTimerCommands map[int]bool, completedSignalCommands, completedInterStateChannelCommands map[int]*iwfidl.EncodedObject) {
	c.stateExecutionsCompletedCommands[stateExecutionId] = service.StateExecutionCompletedCommands{
		CompletedTimerCommands:             completedTimerCommands,
		CompletedSignalCommands:            completedSignalCommands,
		CompletedInterStateChannelCommands: completedInterStateChannelCommands,
	}
}

func (c *ContinueAsNewer) DeleteStateExecutionCompletedCommands(stateExecutionId string) {
	delete(c.stateExecutionsCompletedCommands, stateExecutionId)
}
