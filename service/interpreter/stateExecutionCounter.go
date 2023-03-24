package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateExecutionCounter struct {
	ctx             UnifiedContext
	provider        WorkflowProvider
	config          service.WorkflowConfig
	globalVersioner *globalVersioner

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
		globalVersioner:           NewGlobalVersionProvider(provider, ctx),
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
		return e.updateStateIdSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) MarkStateExecutionCompleted(state iwfidl.StateMovement) error {
	e.pendingStateIdCount[state.StateId]--
	e.totalPendingStateExeCount--
	if e.pendingStateIdCount[state.StateId] == 0 {
		delete(e.pendingStateIdCount, state.StateId)
		return e.updateStateIdSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) GetTotalPendingStateExecutions() int {
	return e.totalPendingStateExeCount
}

func (e *StateExecutionCounter) updateStateIdSearchAttribute() error {
	var executingStateIds []string
	for sid := range e.pendingStateIdCount {
		executingStateIds = append(executingStateIds, sid)
	}
	if e.config.DisableSystemSearchAttributes {
		return nil
	}
	if e.globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() && len(executingStateIds) == 0 {
		// we don't clear search attributes because there are only two possible cases:
		// 1. there will be another stateId being upsert right after this. So this will avoid calling the upsertSA twice
		// 2. there will not be another stateId being upsert. Then this will be cleared before the workflow is closed.
		// see workflowImpl.go to call ClearStateIdSearchAttributeFinally at the end
		return nil
	}
	return e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
		service.SearchAttributeExecutingStateIds: executingStateIds,
	})
}

// ClearStateIdSearchAttributeFinally should only be called at the end of workflow
func (e *StateExecutionCounter) ClearStateIdSearchAttributeFinally() {
	if e.globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() && e.totalPendingStateExeCount == 0 {
		err := e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
			service.SearchAttributeExecutingStateIds: []string{},
		})
		if err != nil {
			e.provider.GetLogger(e.ctx).Error("error for upseart SearchAttributeExecutingStateIds", err)
		}
	}
}
