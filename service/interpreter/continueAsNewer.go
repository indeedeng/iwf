package interpreter

import (
	"encoding/json"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter/env"
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
	timerProcessor        *TimerProcessor
}

func NewContinueAsNewer(
	provider WorkflowProvider,
	interStateChannel *InterStateChannel, signalReceiver *SignalReceiver, stateExecutionCounter *StateExecutionCounter,
	persistenceManager *PersistenceManager, stateRequestQueue *StateRequestQueue, collector *OutputCollector, timerProcessor *TimerProcessor,
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
		timerProcessor:        timerProcessor,
	}
}

func LoadInternalsFromPreviousRun(ctx UnifiedContext, provider WorkflowProvider, previousRunId string, continueAsNewPageSizeInBytes int32) (*service.ContinueAsNewDumpResponse, error) {
	activityOptions := ActivityOptions{
		StartToCloseTimeout: 5 * time.Second,
		RetryPolicy: &iwfidl.RetryPolicy{
			MaximumIntervalSeconds: iwfidl.PtrInt32(5),
		},
	}
	config := env.GetSharedConfig()
	if config.Interpreter.InterpreterActivityConfig.DumpWorkflowInternalActivityConfig != nil {
		activityConfig := config.Interpreter.InterpreterActivityConfig.DumpWorkflowInternalActivityConfig
		activityOptions.StartToCloseTimeout = activityConfig.StartToCloseTimeout
		if activityConfig.RetryPolicy != nil {
			activityOptions.RetryPolicy = activityConfig.RetryPolicy
		}
	}

	ctx = provider.WithActivityOptions(ctx, activityOptions)
	workflowId := provider.GetWorkflowInfo(ctx).WorkflowExecution.ID
	pageSize := continueAsNewPageSizeInBytes
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
				WorkflowRunId:   previousRunId,
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

	var resp service.ContinueAsNewDumpResponse
	err := json.Unmarshal([]byte(sb.String()), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext) error {
	return c.provider.SetQueryHandler(ctx, service.ContinueAsNewDumpQueryType, func() (*service.ContinueAsNewDumpResponse, error) {
		return &service.ContinueAsNewDumpResponse{
			InterStateChannelReceived:  c.interStateChannel.ReadReceived(nil),
			SignalsReceived:            c.signalReceiver.DumpReceived(nil),
			StateExecutionCounterInfo:  c.stateExecutionCounter.Dump(),
			DataObjects:                c.persistenceManager.GetAllDataObjects(),
			SearchAttributes:           c.persistenceManager.GetAllSearchAttributes(),
			StatesToStartFromBeginning: c.stateRequestQueue.GetAllStateStartRequests(),
			StateExecutionsToResume:    c.StateExecutionToResumeMap,
			StateOutputs:               c.outputCollector.GetAll(),
			StaleSkipTimerSignals:      c.timerProcessor.Dump(),
		}, nil
	})
}

func (c *ContinueAsNewer) AddPotentialStateExecutionToResume(
	stateExecutionId string, state iwfidl.StateMovement, stateExecLocals []iwfidl.KeyValue, commandRequest iwfidl.CommandRequest,
	completedTimerCommands map[int]service.InternalTimerStatus, completedSignalCommands, completedInterStateChannelCommands map[int]*iwfidl.EncodedObject,
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

func (c *ContinueAsNewer) HasAnyStateExecutionToResume() bool {
	return len(c.StateExecutionToResumeMap) > 0
}
func (c *ContinueAsNewer) RemoveStateExecutionToResume(stateExecutionId string) {
	delete(c.StateExecutionToResumeMap, stateExecutionId)
}

func (c *ContinueAsNewer) DrainThreads(ctx UnifiedContext) error {
	// TODO: add metric for before and after Await to monitor stuck
	// NOTE: consider using AwaitWithTimeout to get an alert when workflow stuck due to a bug in the draining logic for continueAsNew

	errWait := c.provider.Await(ctx, func() bool {
		return c.allTHreadsDrained(ctx)
	})
	c.provider.GetLogger(ctx).Info("done draining threads for continueAsNew", errWait)

	return errWait
}

// if the DrainAllSignalsAndThreads await is being called more than a few times and cannot get through,
// there is very likely something wrong in the continueAsNew logic
// the key is runId, the value is how many times it has been called in this worker
// Using this in memory counter sot hat we don't have to use AwaitWithTimeout which will consume a timer
// TODO add TTL support because we don't have to keep the value in memory forever(likely a few hours or a day is enough)
var inMemoryContinueAsNewMonitor = make(map[string]time.Time)

const warnThreshold = time.Second * 5
const errThreshold = time.Second * 15

func (c *ContinueAsNewer) allTHreadsDrained(ctx UnifiedContext) bool {
	remainingThreadCount := c.provider.GetThreadCount()
	if remainingThreadCount == 0 {
		return true
	}

	// TODO using a flag to control this debugging info
	runId := c.provider.GetWorkflowInfo(ctx).WorkflowExecution.RunID

	c.provider.GetLogger(ctx).Debug("continueAsNew is in draining remainingThreadCount, attempt, threadNames", remainingThreadCount, inMemoryContinueAsNewMonitor[runId], c.provider.GetPendingThreadNames())

	initTime, ok := inMemoryContinueAsNewMonitor[runId]
	if !ok {
		inMemoryContinueAsNewMonitor[runId] = time.Now()
		return false
	}

	elapsed := time.Since(initTime)

	if elapsed >= errThreshold {
		c.provider.GetLogger(ctx).Warn("continueAsNew is VERY LIKELY stuck in draining remainingThreadCount, attempt, threadNames", remainingThreadCount, inMemoryContinueAsNewMonitor[runId], c.provider.GetPendingThreadNames())
		return false
	}
	if elapsed >= warnThreshold {
		c.provider.GetLogger(ctx).Warn("continueAsNew may be stuck in draining remainingThreadCount, attempt, threadNames", remainingThreadCount, inMemoryContinueAsNewMonitor[runId], c.provider.GetPendingThreadNames())
	}
	return false
}
