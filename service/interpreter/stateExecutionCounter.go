package interpreter

import (
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/interpreter/config"
	"github.com/indeedeng/iwf/service/interpreter/cont"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"reflect"
	"slices"
)

type StateExecutionCounter struct {
	ctx                  interfaces.UnifiedContext
	provider             interfaces.WorkflowProvider
	configer             *config.WorkflowConfiger
	globalVersioner      *GlobalVersioner
	continueAsNewCounter *cont.ContinueAsNewCounter

	stateIdCompletedCounts          map[string]int
	stateIdStartedCounts            map[string]int // For creating stateExecutionId: count the stateId for how many times that have been executed
	stateIdCurrentlyExecutingCounts map[string]int // For system search attribute IwfExecutingStateId: keep counting the stateIds that are executing based on the ExecutingStateIdMode
	totalCurrentlyExecutingCount    int            // For "dead ends": count the total pending states
}

func NewStateExecutionCounter(
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, globalVersioner *GlobalVersioner,
	configer *config.WorkflowConfiger, continueAsNewCounter *cont.ContinueAsNewCounter,
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
	ctx interfaces.UnifiedContext, provider interfaces.WorkflowProvider, globalVersioner *GlobalVersioner,
	stateIdStartedCounts map[string]int, stateIdCurrentlyExecutingCounts map[string]int,
	totalCurrentlyExecutingCount int,
	configer *config.WorkflowConfiger, continueAsNewCounter *cont.ContinueAsNewCounter,
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

		if e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() {
			switch mode := config.GetExecutingStateIdMode(); mode {
			case iwfidl.DISABLED:
				// do nothing
			case iwfidl.ENABLED_FOR_ALL:
				e.increaseStateIdCurrentlyExecutingCounts(s)
				needsUpdateSA = true
			case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
				fallthrough
			default:
				options := s.GetStateOptions()
				if !compatibility.GetSkipWaitUntilApi(&options) {
					e.increaseStateIdCurrentlyExecutingCounts(s)
					needsUpdateSA = true
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
		return e.refreshIwfExecutingStateIdSearchAttribute()
	}
	return nil
}

func (e *StateExecutionCounter) increaseStateIdCurrentlyExecutingCounts(s iwfidl.StateMovement) bool {
	e.stateIdCurrentlyExecutingCounts[s.StateId]++
	// first time the stateId show up
	return e.stateIdCurrentlyExecutingCounts[s.StateId] == 1
}

func (e *StateExecutionCounter) MarkStateExecutionCompleted(currentState iwfidl.StateMovement, nextStates []iwfidl.StateMovement) error {
	e.totalCurrentlyExecutingCount--

	options := currentState.GetStateOptions()
	skipStart := compatibility.GetSkipWaitUntilApi(&options)
	e.continueAsNewCounter.IncExecutedStateExecution(skipStart)

	config := e.configer.Get()

	if e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() {
		switch mode := config.GetExecutingStateIdMode(); mode {
		case iwfidl.DISABLED:
			return nil
		case iwfidl.ENABLED_FOR_ALL:
			e.decreaseStateIdCurrentlyExecutingCounts(currentState)
			shouldSkipUpsert := determineIfShouldSkipRefreshOnCompleted(nextStates, true)
			if shouldSkipUpsert {
				return nil
			}
		case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
			fallthrough
		default:
			if compatibility.GetSkipWaitUntilApi(&options) {
				return nil
			} else {
				e.decreaseStateIdCurrentlyExecutingCounts(currentState)
				shouldSkipRefresh := determineIfShouldSkipRefreshOnCompleted(nextStates, false)
				if shouldSkipRefresh {
					return nil
				}
			}
		}
	} else {
		if config.GetDisableSystemSearchAttribute() {
			return nil
		} else {
			e.decreaseStateIdCurrentlyExecutingCounts(currentState)
		}
	}

	return e.refreshIwfExecutingStateIdSearchAttribute()
}

func determineIfShouldSkipRefreshOnCompleted(nextStates []iwfidl.StateMovement, enabledForAll bool) bool {
	var nonClosingNextStates []iwfidl.StateMovement
	for _, s := range nextStates {
		if _, ok := service.ValidClosingWorkflowStateId[s.GetStateId()]; !ok {
			// s is not a ValidClosingWorkflowStateId
			nonClosingNextStates = append(nonClosingNextStates, s)
		}
	}
	if enabledForAll {
		if len(nonClosingNextStates) > 0 {
			return true
		}
	} else {
		for _, s := range nonClosingNextStates {
			options := s.GetStateOptions()
			if !compatibility.GetSkipWaitUntilApi(&options) {
				return true
			}
		}
	}

	return false
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

func (e *StateExecutionCounter) refreshIwfExecutingStateIdSearchAttribute() error {
	// Optimization: don't upsert SAs if currentSAsValues == stateIdCurrentlyExecutingCounts keys
	if e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() {
		sas, err := e.provider.GetSearchAttributes(e.ctx, []iwfidl.SearchAttributeKeyAndType{
			{Key: ptr.Any(service.SearchAttributeExecutingStateIds),
				ValueType: ptr.Any(iwfidl.KEYWORD_ARRAY)},
		})
		if err != nil {
			e.provider.GetLogger(e.ctx).Error("error for GetSearchAttributes", err)
			return err
		}

		var currentSAsValues []string

		currentSAs, ok := sas[service.SearchAttributeExecutingStateIds]
		if ok {
			currentSAsValues = currentSAs.StringArrayValue
		}

		var executingStateIds []string
		executingStateIds = append(executingStateIds, DeterministicKeys(e.stateIdCurrentlyExecutingCounts)...)

		slices.Sort(currentSAsValues)
		slices.Sort(executingStateIds)
		if reflect.DeepEqual(currentSAsValues, executingStateIds) {
			return nil
		}
	}

	var executingStateIds []string
	executingStateIds = append(executingStateIds, DeterministicKeys(e.stateIdCurrentlyExecutingCounts)...)

	if e.globalVersioner.IsAfterVersionOfOptimizedUpsertSearchAttribute() && !e.globalVersioner.IsAfterVersionOfExecutingStateIdMode() && len(executingStateIds) == 0 {
		// we don't clear search attributes because there are only two possible cases:
		// 1. there will be another stateId being upsert right after this. So this will avoid calling the upsertSA twice
		// 2. there will not be another stateId being upsert. Then this will be cleared before the workflow is closed.
		// see workflowImpl.go to call ClearExecutingStateIdsSearchAttributeFinally at the end
		return nil
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
