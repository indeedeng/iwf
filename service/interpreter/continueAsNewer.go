package interpreter

import (
	"encoding/json"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
)

type ContinueAsNewer struct {
	provider                                WorkflowProvider
	pendingStateExecution                   map[string]service.PendingStateExecution
	pendingStateExecutionsRequestCommands   map[string]service.PendingStateExecutionRequestCommands
	pendingStateExecutionsCompletedCommands map[string]service.PendingStateExecutionCompletedCommands
	interStateChannel                       *InterStateChannel
	stateExecutionCounter                   *StateExecutionCounter
	persistenceManager                      *PersistenceManager
	signalReceiver                          *SignalReceiver
}

func NewContinueAsNewer(
	provider WorkflowProvider,
	interStateChannel *InterStateChannel, signalReceiver *SignalReceiver, stateExecutionCounter *StateExecutionCounter, persistenceManager *PersistenceManager,
) *ContinueAsNewer {
	return &ContinueAsNewer{
		provider:                                provider,
		interStateChannel:                       interStateChannel,
		signalReceiver:                          signalReceiver,
		stateExecutionCounter:                   stateExecutionCounter,
		persistenceManager:                      persistenceManager,
		pendingStateExecutionsCompletedCommands: map[string]service.PendingStateExecutionCompletedCommands{},
		pendingStateExecutionsRequestCommands:   map[string]service.PendingStateExecutionRequestCommands{},
	}
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext) error {
	err := c.provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func() (*service.DumpAllInternalResponse, error) {
		return &service.DumpAllInternalResponse{
			InterStateChannelReceived:               c.interStateChannel.ReadReceived(nil),
			SignalChannelReceived:                   c.signalReceiver.ReadReceived(nil),
			StateExecutionCounterInfo:               c.stateExecutionCounter.Dump(),
			PendingStateExecutionsCompletedCommands: c.pendingStateExecutionsCompletedCommands,
			PendingStateExecutionsRequestCommands:   c.pendingStateExecutionsRequestCommands,
			DataObjects:                             c.persistenceManager.GetAllDataObjects(),
			SearchAttributes:                        c.persistenceManager.GetAllSearchAttributes(),
		}, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *ContinueAsNewer) AddPendingStateExecutionCommandStatus(
	stateExecutionId string,
	completedTimerCommands map[int]bool, completedSignalCommands, completedInterStateChannelCommands map[int]*iwfidl.EncodedObject,
	timerCommands []iwfidl.TimerCommand, signalCommands []iwfidl.SignalCommand, interStateChannelCommands []iwfidl.InterStateChannelCommand,
) {
	c.pendingStateExecutionsCompletedCommands[stateExecutionId] = service.PendingStateExecutionCompletedCommands{
		CompletedTimerCommands:             completedTimerCommands,
		CompletedSignalCommands:            completedSignalCommands,
		CompletedInterStateChannelCommands: completedInterStateChannelCommands,
	}
	c.pendingStateExecutionsRequestCommands[stateExecutionId] = service.PendingStateExecutionRequestCommands{
		TimerCommands:             timerCommands,
		SignalCommands:            signalCommands,
		InterStateChannelCommands: interStateChannelCommands,
	}
}

func (c *ContinueAsNewer) ClearPendingStateExecutionCommandStatus(stateExecutionId string) {
	delete(c.pendingStateExecutionsCompletedCommands, stateExecutionId)
	delete(c.pendingStateExecutionsRequestCommands, stateExecutionId)
}

func (c *ContinueAsNewer) DrainAllSignalsAndThreads(ctx UnifiedContext) error {
	// NOTE: consider using AwaitWithTimeout to get an alert when workflow stuck due to a bug in the draining logic for continueAsNew
	return c.provider.Await(ctx, func() bool {
		return c.canContinueAsNew(ctx)
	})
}

func (c *ContinueAsNewer) canContinueAsNew(ctx UnifiedContext) bool {
	// drain all signals + all threads
	return c.signalReceiver.HaveAllUserAndSystemSignalsToReceive(ctx) && c.provider.GetThreadCount() == 0
}

func (c *ContinueAsNewer) ProcessUncompletedStateExecution(stateExecStatus service.StateExecutionStatus, stateExeId string, state iwfidl.StateMovement) {
	c.pendingStateExecution[stateExeId] = service.PendingStateExecution{
		State:                state,
		StateExecutionStatus: stateExecStatus,
	}
}

func (c *ContinueAsNewer) ContinueToNewRun(ctx UnifiedContext, execution service.IwfWorkflowExecution, config iwfidl.WorkflowConfig) error {
	return c.provider.NewInterpreterContinueAsNewError(ctx, service.InterpreterWorkflowInput{
		ContinueAsNew: true,
		ContinueAsNewInput: service.ContinueAsNewInput{
			Config:               config,
			IwfWorkflowExecution: execution,
		},
	})
}

func ResumeFromPreviousRun(input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	// TODO this is for test only, will be implemented in later PRs
	data, err := json.Marshal(input)
	return &service.InterpreterWorkflowOutput{
		StateCompletionOutputs: []iwfidl.StateCompletionOutput{
			{
				CompletedStateOutput: &iwfidl.EncodedObject{
					Data: ptr.Any(string(data)),
				},
			},
		},
	}, err
}
