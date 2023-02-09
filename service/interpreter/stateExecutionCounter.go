package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateExecutionCounter struct {
	ctx      UnifiedContext
	provider WorkflowProvider
	config   service.WorkflowConfig

	executedStateIdCount      map[string]int // count the stateId for how many times that have een executed so that we can create stateExecutionId
	pendingStateIdCount       map[string]int // keep counting the pending stateIds so that we know times to upsert system search attributes service.SearchAttributeExecutingStateIds
	totalPendingStateExeCount int            // count the total pending states so that we know the workflow can complete when all threads reach "dead ends"
}

func NewStateExecutionCounter(ctx UnifiedContext, provider WorkflowProvider, config service.WorkflowConfig) *StateExecutionCounter {
	return &StateExecutionCounter{
		ctx:                       ctx,
		provider:                  provider,
		pendingStateIdCount:       make(map[string]int),
		executedStateIdCount:      make(map[string]int),
		totalPendingStateExeCount: 0,
		config:                    config,
	}
}

func RebuildStateExecutionManager(ctx UnifiedContext, provider WorkflowProvider,
	executedStateIdCount map[string]int, pendingStateIdCount map[string]int, totalPendingStateExeCount int,
) *StateExecutionCounter {
	return &StateExecutionCounter{
		ctx:                       ctx,
		provider:                  provider,
		pendingStateIdCount:       pendingStateIdCount,
		executedStateIdCount:      executedStateIdCount,
		totalPendingStateExeCount: totalPendingStateExeCount,
	}
}

func (e *StateExecutionCounter) Dump() service.StateExecutionCounterInfo {
	return service.StateExecutionCounterInfo{
		ExecutedStateIdCount:      e.executedStateIdCount,
		PendingStateIdCount:       e.pendingStateIdCount,
		TotalPendingStateExeCount: e.totalPendingStateExeCount,
	}
}

func (e *StateExecutionCounter) CreateNextExecutionId(stateId string) string {
	e.executedStateIdCount[stateId]++
	id := e.executedStateIdCount[stateId]
	return fmt.Sprintf("%v-%v", stateId, id)
}

func (e *StateExecutionCounter) MarkStateExecutionsPending(states []iwfidl.StateMovement) error {
	needsUpdate := false
	for _, s := range states {
		e.pendingStateIdCount[s.StateId]++
		if e.pendingStateIdCount[s.StateId] == 1 {
			// first time the stateId show up
			needsUpdate = true
		}
	}
	e.totalPendingStateExeCount += len(states)
	if needsUpdate {
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) MarkStateExecutionCompleted(state iwfidl.StateMovement) error {
	e.pendingStateIdCount[state.StateId]--
	e.totalPendingStateExeCount--
	if e.pendingStateIdCount[state.StateId] == 0 {
		delete(e.pendingStateIdCount, state.StateId)
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) GetTotalPendingStateExecutions() int {
	return e.totalPendingStateExeCount
}

func (e *StateExecutionCounter) updateSearchAttribute() error {
	var executingStateIds []string
	for sid := range e.pendingStateIdCount {
		executingStateIds = append(executingStateIds, sid)
	}
	if e.config.DisableSystemSearchAttributes {
		return nil
	}
	return e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
		service.SearchAttributeExecutingStateIds: executingStateIds,
	})
}
