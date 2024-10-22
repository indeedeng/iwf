package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/compatibility"
)

type StateExecutionCounter struct {
	ctx                  UnifiedContext
	provider             WorkflowProvider
	configer             *WorkflowConfiger
	globalVersioner      *GlobalVersioner
	continueAsNewCounter *ContinueAsNewCounter

	stateIdCompletedCounts          map[string]int
	stateIdStartedCounts            map[string]int // For creating stateExecutionId: count the stateId for how many times that have been executed
	stateIdCurrentlyExecutingCounts map[string]int // For system search attribute IwfExecutingStateId: keep counting the stateIds that are executing based on the ExecutingStateIdMode
	totalCurrentlyExecutingCount    int            // For "dead ends": count the total pending states
}

type StateTransition struct {
	current iwfidl.StateMovement
	next    []iwfidl.StateMovement
}

func NewStateExecutionCounter(
	ctx UnifiedContext, provider WorkflowProvider, globalVersioner *GlobalVersioner,
	configer *WorkflowConfiger, continueAsNewCounter *ContinueAsNewCounter,
) *StateExecutionCounter {
	return &StateExecutionCounter{
		ctx:                             ctx,
		provider:                        provider,
		stateIdStartedCounts:            make(map[string]int),
		stateIdCurrentlyExecutingCounts: make(map[string]int),
		totalCurrentlyExecutingCount:    0,
		configer:                        configer,
		globalVersioner:                 globalVersioner,
		continueAsNewCounter:            continueAsNewCounter,
	}
}

func RebuildStateExecutionCounter(
	ctx UnifiedContext, provider WorkflowProvider, globalVersioner *GlobalVersioner,
	stateIdStartedCounts map[string]int, stateIdCurrentlyExecutingCounts map[string]int,
	totalCurrentlyExecutingCount int,
	configer *WorkflowConfiger, continueAsNewCounter *ContinueAsNewCounter,
) *StateExecutionCounter {
	return &StateExecutionCounter{
		ctx:                             ctx,
		provider:                        provider,
		stateIdStartedCounts:            stateIdStartedCounts,
		stateIdCurrentlyExecutingCounts: stateIdCurrentlyExecutingCounts,
		totalCurrentlyExecutingCount:    totalCurrentlyExecutingCount,
		configer:                        configer,
		globalVersioner:                 globalVersioner,
		continueAsNewCounter:            continueAsNewCounter,
	}
}

func (e *StateExecutionCounter) Dump() service.StateExecutionCounterInfo {
	return service.StateExecutionCounterInfo{
		StateIdStartedCount:            e.stateIdStartedCounts,
		StateIdCurrentlyExecutingCount: e.stateIdCurrentlyExecutingCounts,
		TotalCurrentlyExecutingCount:   e.totalCurrentlyExecutingCount,
	}
}

func (e *StateExecutionCounter) CreateNextExecutionId(stateId string) string {
	e.stateIdStartedCounts[stateId]++
	id := e.stateIdStartedCounts[stateId]
	return fmt.Sprintf("%v-%v", stateId, id)
}

func (e *StateExecutionCounter) MarkStateIdExecutingIfNotYet(stateReqs []StateRequest) error {
	config := e.configer.Get()

	needsUpdateSA := false
	numOfNew := 0
	for _, sr := range stateReqs {
		if sr.IsResumeRequest() {
			continue
		}
		s := sr.GetStateStartRequest()

		// e.provider.GetSearchAttributes()
		// TODO check if SA in context and decide if need to update

		if e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() {
			switch mode := config.GetExecutingStateIdMode(); mode {
			case iwfidl.DISABLED:
				// do nothing
			case iwfidl.ENABLED_FOR_ALL:
				if e.increaseStateIdCurrentlyExecutingCounts(s) {
					needsUpdateSA = true
				}
			case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
				fallthrough
			default:
				options := s.GetStateOptions()
				if !compatibility.GetSkipWaitUntilApi(&options) {
					if e.increaseStateIdCurrentlyExecutingCounts(s) {
						needsUpdateSA = true
					}
				}
			}
		} else {
			if !config.GetDisableSystemSearchAttribute() {
				if e.increaseStateIdCurrentlyExecutingCounts(s) {
					needsUpdateSA = true
				}
			}
		}

		numOfNew++
	}
	e.totalCurrentlyExecutingCount += numOfNew

	if needsUpdateSA {
		return e.updateStateIdSearchAttribute(nil)
	}
	return nil
}

func (e *StateExecutionCounter) increaseStateIdCurrentlyExecutingCounts(s iwfidl.StateMovement) bool {
	e.stateIdCurrentlyExecutingCounts[s.StateId]++
	// first time the stateId show up
	return e.stateIdCurrentlyExecutingCounts[s.StateId] == 1
}

func (e *StateExecutionCounter) MarkStateExecutionCompleted(state iwfidl.StateMovement, nextStates []iwfidl.StateMovement) error {
	e.totalCurrentlyExecutingCount--

	options := state.GetStateOptions()
	skipStart := compatibility.GetSkipWaitUntilApi(&options)
	e.continueAsNewCounter.IncExecutedStateExecution(skipStart)

	config := e.configer.Get()

	var stateTransition *StateTransition

	if e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() {
		switch mode := config.GetExecutingStateIdMode(); mode {
		case iwfidl.DISABLED:
			return nil
		case iwfidl.ENABLED_FOR_ALL:
			e.decreaseStateIdCurrentlyExecutingCounts(state)
			stateTransition = &StateTransition{
				current: state,
				next:    nextStates,
			}
		case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
			fallthrough
		default:
			if compatibility.GetSkipWaitUntilApi(&options) {
				return nil
			} else {
				e.decreaseStateIdCurrentlyExecutingCounts(state)
				stateTransition = &StateTransition{
					current: state,
					next:    nextStates,
				}
			}
		}
	} else {
		if config.GetDisableSystemSearchAttribute() {
			return nil
		} else {
			e.decreaseStateIdCurrentlyExecutingCounts(state)
		}
	}

	return e.updateStateIdSearchAttribute(stateTransition)
}

func (e *StateExecutionCounter) decreaseStateIdCurrentlyExecutingCounts(state iwfidl.StateMovement) {
	e.stateIdCurrentlyExecutingCounts[state.StateId]--
	if e.stateIdCurrentlyExecutingCounts[state.StateId] == 0 {
		delete(e.stateIdCurrentlyExecutingCounts, state.StateId)
	}
}

func (e *StateExecutionCounter) GetTotalCurrentlyExecutingCount() int {
	return e.totalCurrentlyExecutingCount
}

// Next states of stateTransition argument will be evaluated. If the next states skip waitUntil, upsertSearchAttributes will not be invoked
// stateExecuted should be only provided when updateStateIdSearchAttribute is called after decreaseStateIdCurrentlyExecutingCounts
func (e *StateExecutionCounter) updateStateIdSearchAttribute(stateTransition *StateTransition) error {
	var executingStateIds []string
	for sid := range e.stateIdCurrentlyExecutingCounts {
		executingStateIds = append(executingStateIds, sid)
	}

	if e.globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() && !e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() && len(executingStateIds) == 0 {
		// we don't clear search attributes because there are only two possible cases:
		// 1. there will be another stateId being upsert right after this. So this will avoid calling the upsertSA twice
		// 2. there will not be another stateId being upsert. Then this will be cleared before the workflow is closed.
		// see workflowImpl.go to call ClearExecutingStateIdsSearchAttributeFinally at the end
		return nil
	}

	if stateTransition != nil {
		// UpsertSearchAttributes should be only invoked if when the next states do not skip waitUntil
		// Or if current and next states are the same

		// State transition loops back to the current state
		if len(stateTransition.next) == 1 && stateTransition.current.StateId == stateTransition.next[0].StateId {
			return nil
		}

		shouldSkipUpsertingSAs := true
		for _, s := range stateTransition.next {
			if !s.StateOptions.GetSkipWaitUntil() {
				shouldSkipUpsertingSAs = false
				break
			}
		}

		if shouldSkipUpsertingSAs {
			return nil
		}

	}

	return e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
		service.SearchAttributeExecutingStateIds: executingStateIds,
	})
}

// ClearExecutingStateIdsSearchAttributeFinally should only be called at the end of workflow
func (e *StateExecutionCounter) ClearExecutingStateIdsSearchAttributeFinally() {
	config := e.configer.Get()

	if e.globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() && !e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() && e.totalCurrentlyExecutingCount == 0 {
		if config.GetDisableSystemSearchAttribute() {
			return
		}

		err := e.provider.UpsertSearchAttributes(e.ctx, map[string]interface{}{
			service.SearchAttributeExecutingStateIds: []string{},
		})
		if err != nil {
			e.provider.GetLogger(e.ctx).Error("error for upsert SearchAttributeExecutingStateIds", err)
		}
	}
}
