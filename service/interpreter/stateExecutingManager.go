package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type stateExecutingManager struct {
	ctx      UnifiedContext
	provider WorkflowProvider

	stateIdCountMap     map[string]int // count the stateId being executed so that we can create stateExecutionId
	pendingStateIdCount map[string]int // keep counting the pending stateIds so that we know times to upsert search attributes
	totalPendingCount   int            // count the total pending states so that we know the workflow can complete when all threads reach "dead ends"
}

func newStateExecutionManager(ctx UnifiedContext, provider WorkflowProvider) *stateExecutingManager {
	return &stateExecutingManager{
		ctx:                 ctx,
		provider:            provider,
		pendingStateIdCount: make(map[string]int),
		stateIdCountMap:     make(map[string]int),
		totalPendingCount:   0,
	}
}

func rebuildStateExecutionManager(ctx UnifiedContext, provider WorkflowProvider,
	stateIdCountMap map[string]int, pendingStateIdCount map[string]int, totalPendingCount int,
) *stateExecutingManager {
	return &stateExecutingManager{
		ctx:                 ctx,
		provider:            provider,
		pendingStateIdCount: pendingStateIdCount,
		stateIdCountMap:     stateIdCountMap,
		totalPendingCount:   totalPendingCount,
	}
}

func (e *stateExecutingManager) getCarryOverData() (stateIdCountMap map[string]int, pendingStateIdCount map[string]int, totalPendingCount int) {
	return e.stateIdCountMap, e.pendingStateIdCount, e.totalPendingCount
}

func (e *stateExecutingManager) createNextExecutionId(stateId string) string {
	e.stateIdCountMap[stateId]++
	id := e.stateIdCountMap[stateId]
	return fmt.Sprintf("%v-%v", stateId, id)
}

func (e *stateExecutingManager) startStates(states []iwfidl.StateMovement) error {
	needsUpdate := false
	for _, s := range states {
		e.pendingStateIdCount[s.StateId]++
		if e.pendingStateIdCount[s.StateId] == 1 {
			// first time the stateId show up
			needsUpdate = true
		}
	}
	e.totalPendingCount += len(states)
	if needsUpdate {
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *stateExecutingManager) completeStates(state iwfidl.StateMovement) error {
	e.pendingStateIdCount[state.StateId]--
	e.totalPendingCount -= 1
	if e.pendingStateIdCount[state.StateId] == 0 {
		delete(e.pendingStateIdCount, state.StateId)
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *stateExecutingManager) getTotalPendingStates() int {
	return e.totalPendingCount
}

func (e *stateExecutingManager) updateSearchAttribute() error {
	var executingStateIds []string
	for sid := range e.pendingStateIdCount {
		executingStateIds = append(executingStateIds, sid)
	}
	return e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
		service.SearchAttributeExecutingStateIds: executingStateIds,
	})
}
