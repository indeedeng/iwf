package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type stateExecutingManager struct {
	ctx      UnifiedContext
	provider WorkflowProvider

	stateIdCount        map[string]int
	totalExecutingCount int
}

func newStateExecutingManager(ctx UnifiedContext, provider WorkflowProvider) *stateExecutingManager {
	return &stateExecutingManager{
		ctx:                 ctx,
		provider:            provider,
		stateIdCount:        map[string]int{},
		totalExecutingCount: 0,
	}
}

func (e *stateExecutingManager) startStates(states []iwfidl.StateMovement) error {
	needsUpdate := false
	for _, s := range states {
		e.stateIdCount[s.StateId]++
		if e.stateIdCount[s.StateId] == 1 {
			// first time the stateId show up
			needsUpdate = true
		}
	}
	e.totalExecutingCount += len(states)
	if needsUpdate {
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *stateExecutingManager) completeStates(state iwfidl.StateMovement) error {
	e.stateIdCount[state.StateId]--
	e.totalExecutingCount -= 1
	if e.stateIdCount[state.StateId] == 0 {
		delete(e.stateIdCount, state.StateId)
		return e.updateSearchAttribute()
	}
	return nil
}

func (e *stateExecutingManager) getTotalExecutingStates() int {
	return e.totalExecutingCount
}

func (e *stateExecutingManager) updateSearchAttribute() error {
	var executingStateIds []string
	for sid := range e.stateIdCount {
		executingStateIds = append(executingStateIds, sid)
	}
	return e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
		service.SearchAttributeExecutingStateIds: executingStateIds,
	})
}
