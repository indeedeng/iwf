package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateExecutionCounter struct {
	ctx                  UnifiedContext
	provider             WorkflowProvider
	config               iwfidl.WorkflowConfig
	globalVersioner      *GlobalVersioner
	continueAsNewCounter *ContinueAsNewCounter

	executedStateIdCount      map[string]int // For creating stateExecutionId: count the stateId for how many times that have been executed
	pendingStateIdCount       map[string]int // For system search attributes service.SearchAttributeExecutingStateIds: keep counting the pending stateIds
	totalPendingStateExeCount int            // For "dead ends": count the total pending states
}

func NewStateExecutionCounter(ctx UnifiedContext, provider WorkflowProvider, config iwfidl.WorkflowConfig, continueAsNewCounter *ContinueAsNewCounter) *StateExecutionCounter {
	return &StateExecutionCounter{
		ctx:                       ctx,
		provider:                  provider,
		pendingStateIdCount:       make(map[string]int),
		executedStateIdCount:      make(map[string]int),
		totalPendingStateExeCount: 0,
		config:                    config,
		globalVersioner:           NewGlobalVersioner(provider, ctx),
		continueAsNewCounter:      continueAsNewCounter,
	}
}

func RebuildStateExecutionCounter(ctx UnifiedContext, provider WorkflowProvider,
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

func (e *StateExecutionCounter) MarkStateExecutionsPendingIfHavenot(stateReqs []StateRequest) error {
	needsUpdate := false
	numOfNew := 0
	for _, sr := range stateReqs {
		if sr.IsPendingFromContinueAsNew() {
			continue
		}
		s := sr.GetNewRequest()
		numOfNew++
		e.pendingStateIdCount[s.StateId]++
		if e.pendingStateIdCount[s.StateId] == 1 {
			// first time the stateId show up
			needsUpdate = true
		}
	}
	e.totalPendingStateExeCount += numOfNew
	if needsUpdate {
		return e.updateStateIdSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) MarkStateExecutionCompleted(state iwfidl.StateMovement) error {
	e.pendingStateIdCount[state.StateId]--
	e.totalPendingStateExeCount--
	e.continueAsNewCounter.IncExecutedStateExecution()
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
	if e.config.GetDisableSystemSearchAttribute() {
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
