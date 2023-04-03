package interpreter

import (
	"encoding/json"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"strings"
	"time"
)

type ContinueAsNewer struct {
	provider WorkflowProvider

	StateExecutionToResumeMap map[string]service.StateExecutionResumeInfo // stateExeId to StateExecutionResumeInfo

	stateRequestQueue     *StateRequestQueue
	interStateChannel     *InterStateChannel
	stateExecutionCounter *StateExecutionCounter
	persistenceManager    *PersistenceManager
	signalReceiver        *SignalReceiver
	outputCollector       *OutputCollector
}

func NewContinueAsNewer(
	provider WorkflowProvider,
	interStateChannel *InterStateChannel, signalReceiver *SignalReceiver, stateExecutionCounter *StateExecutionCounter,
	persistenceManager *PersistenceManager, stateRequestQueue *StateRequestQueue, collector *OutputCollector,
) *ContinueAsNewer {
	return &ContinueAsNewer{
		provider: provider,

		StateExecutionToResumeMap: map[string]service.StateExecutionResumeInfo{},

		stateRequestQueue:     stateRequestQueue,
		interStateChannel:     interStateChannel,
		signalReceiver:        signalReceiver,
		stateExecutionCounter: stateExecutionCounter,
		persistenceManager:    persistenceManager,
		outputCollector:       collector,
	}
}

func LoadInternalsFromPreviousRun(ctx UnifiedContext, provider WorkflowProvider, input service.InterpreterWorkflowInput) (*service.DumpAllInternalResponse, error) {
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = provider.WithActivityOptions(ctx, activityOptions)
	workflowId := provider.GetWorkflowInfo(ctx).WorkflowExecution.ID
	runId := input.ContinueAsNewInput.PreviousInternalRunId
	pageSize := input.Config.GetContinueAsNewPageSizeInBytes()
	if pageSize == 0 {
		pageSize = service.DefaultContinueAsNewPageSizeInBytes
	}
	var sb strings.Builder
	lastChecksum := ""
	pageNum := int32(0)
	for {
		var resp iwfidl.WorkflowDumpResponse
		err := provider.ExecuteActivity(ctx, DumpWorkflowInternal, provider.GetBackendType(),
			iwfidl.WorkflowDumpRequest{
				WorkflowId:      workflowId,
				WorkflowRunId:   runId,
				PageNum:         pageNum,
				PageSizeInBytes: pageSize,
			}).Get(ctx, &resp)
		if err != nil {
			return nil, err
		}
		if lastChecksum != "" && lastChecksum != resp.Checksum {
			// reset to start from beginning
			pageNum = 0
			lastChecksum = ""
			sb.Reset()
			provider.GetLogger(ctx).Error("checksum has changed during the loading", lastChecksum, resp.Checksum)
			continue
		} else {
			lastChecksum = resp.Checksum
			sb.WriteString(resp.JsonData)
			pageNum++
			if pageNum >= resp.TotalPages {
				break
			}
		}
	}

	var resp service.DumpAllInternalResponse
	err := json.Unmarshal([]byte(sb.String()), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *ContinueAsNewer) createDumpAllInternalResponse() *service.DumpAllInternalResponse {
	return &service.DumpAllInternalResponse{
		InterStateChannelReceived:  c.interStateChannel.ReadReceived(nil),
		SignalsReceived:            c.signalReceiver.DumpReceived(nil),
		StateExecutionCounterInfo:  c.stateExecutionCounter.Dump(),
		DataObjects:                c.persistenceManager.GetAllDataObjects(),
		SearchAttributes:           c.persistenceManager.GetAllSearchAttributes(),
		StatesToStartFromBeginning: c.stateRequestQueue.GetAllStateStartRequests(),
		StateExecutionsToResume:    c.StateExecutionToResumeMap,
		StateOutputs:               c.outputCollector.GetAll(),
	}
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext) error {
	return c.provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func() (*service.DumpAllInternalResponse, error) {
		return c.createDumpAllInternalResponse(), nil
	})
}

func (c *ContinueAsNewer) AddPotentialStateExecutionToResume(
	stateExecutionId string, state iwfidl.StateMovement, stateExecLocals []iwfidl.KeyValue, commandRequest iwfidl.CommandRequest,
	completedTimerCommands map[int]bool, completedSignalCommands, completedInterStateChannelCommands map[int]*iwfidl.EncodedObject,
) {
	c.StateExecutionToResumeMap[stateExecutionId] = service.StateExecutionResumeInfo{
		StateExecutionId:     stateExecutionId,
		State:                state,
		StateExecutionLocals: stateExecLocals,
		CommandRequest:       commandRequest,
		StateExecutionCompletedCommands: service.StateExecutionCompletedCommands{
			CompletedTimerCommands:             completedTimerCommands,
			CompletedSignalCommands:            completedSignalCommands,
			CompletedInterStateChannelCommands: completedInterStateChannelCommands,
		},
	}
}

func (c *ContinueAsNewer) RemoveStateExecutionToResume(stateExecutionId string) {
	delete(c.StateExecutionToResumeMap, stateExecutionId)
}

func (c *ContinueAsNewer) DrainAllSignalsAndThreads(ctx UnifiedContext) error {
	// TODO: add metric for before and after Await to monitor stuck
	// NOTE: consider using AwaitWithTimeout to get an alert when workflow stuck due to a bug in the draining logic for continueAsNew
	return c.provider.Await(ctx, func() bool {
		return c.canContinueAsNew(ctx)
	})
}

func (c *ContinueAsNewer) canContinueAsNew(ctx UnifiedContext) bool {
	// drain all signals + all threads
	return c.signalReceiver.HaveAllUserAndSystemSignalsToReceive(ctx) && c.provider.GetThreadCount() == 0
}
