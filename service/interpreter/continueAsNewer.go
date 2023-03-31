package interpreter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"math"
	"strings"
	"time"
)

type ContinueAsNewer struct {
	rootCtx  UnifiedContext
	provider WorkflowProvider

	pendingStateExecution                   []service.PendingStateExecution
	pendingStateExecutionsRequestCommands   map[string]service.PendingStateExecutionRequestCommands
	pendingStateExecutionsCompletedCommands map[string]service.PendingStateExecutionCompletedCommands
	stateRequestQueue                       *StateRequestQueue

	interStateChannel     *InterStateChannel
	stateExecutionCounter *StateExecutionCounter
	persistenceManager    *PersistenceManager
	signalReceiver        *SignalReceiver
}

func NewContinueAsNewer(
	provider WorkflowProvider,
	interStateChannel *InterStateChannel, signalReceiver *SignalReceiver, stateExecutionCounter *StateExecutionCounter,
	persistenceManager *PersistenceManager, stateRequestQueue *StateRequestQueue,
) *ContinueAsNewer {
	return &ContinueAsNewer{
		provider: provider,

		pendingStateExecution:                   nil,
		pendingStateExecutionsCompletedCommands: map[string]service.PendingStateExecutionCompletedCommands{},
		pendingStateExecutionsRequestCommands:   map[string]service.PendingStateExecutionRequestCommands{},
		stateRequestQueue:                       stateRequestQueue,

		interStateChannel:     interStateChannel,
		signalReceiver:        signalReceiver,
		stateExecutionCounter: stateExecutionCounter,
		persistenceManager:    persistenceManager,
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
	pageNum := 0
	for {
		var resp service.DumpAllInternalWithPaginationResponse
		err := provider.ExecuteActivity(ctx, DumpWorkflowInternal, workflowId, runId, service.DumpAllInternalWithPaginationRequest{
			PageSizeInBytes: int(pageSize),
			PageNum:         pageNum,
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
		InterStateChannelReceived:               c.interStateChannel.ReadReceived(nil),
		SignalsReceived:                         c.signalReceiver.DumpReceived(nil),
		StateExecutionCounterInfo:               c.stateExecutionCounter.Dump(),
		PendingStateExecutionsCompletedCommands: c.pendingStateExecutionsCompletedCommands,
		PendingStateExecutionsRequestCommands:   c.pendingStateExecutionsRequestCommands,
		DataObjects:                             c.persistenceManager.GetAllDataObjects(),
		SearchAttributes:                        c.persistenceManager.GetAllSearchAttributes(),
		NonStartedStates:                        c.stateRequestQueue.GetAllNonPendingRequest(),
		PendingStateExecution:                   c.pendingStateExecution,
	}
}

func (c *ContinueAsNewer) SetQueryHandlersForContinueAsNew(ctx UnifiedContext) error {
	err := c.provider.SetQueryHandler(ctx, service.DumpAllInternalQueryType, func() (*service.DumpAllInternalResponse, error) {
		return c.createDumpAllInternalResponse(), nil
	})
	if err != nil {
		return err
	}
	return c.provider.SetQueryHandler(ctx, service.DumpAllInternalWithPaginationQueryType, func(req service.DumpAllInternalWithPaginationRequest) (*service.DumpAllInternalWithPaginationResponse, error) {
		resp := c.createDumpAllInternalResponse()
		data, err := json.Marshal(resp)
		if err != nil {
			return nil, err
		}
		checksum := md5.Sum(data)
		pageSize := service.DefaultContinueAsNewPageSizeInBytes
		if req.PageSizeInBytes > 0 {
			pageSize = req.PageSizeInBytes
		}
		lenInDouble := float64(len(data))
		totalPages := int(math.Ceil(lenInDouble / float64(pageSize)))
		if req.PageNum >= totalPages {
			return nil, fmt.Errorf("wrong pageNum, max is %v", totalPages-1)
		}
		start := pageSize * req.PageNum
		end := start + pageSize
		if end > len(data) {
			end = len(data)
		}
		return &service.DumpAllInternalWithPaginationResponse{
			Checksum:   string(checksum[:]),
			TotalPages: totalPages,
			JsonData:   string(data[start:end]),
		}, nil
	})
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

func (c *ContinueAsNewer) ProcessUncompletedStateExecution(stateExecStatus service.StateExecutionStatus, stateExeId string, state iwfidl.StateMovement) {
	c.pendingStateExecution = append(c.pendingStateExecution, service.PendingStateExecution{
		StateExecutionId:     stateExeId,
		State:                state,
		StateExecutionStatus: stateExecStatus,
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
