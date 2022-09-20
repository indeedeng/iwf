package interpreter

import "fmt"

type StateExecutionIdManager struct {
	stateIdMap map[string]int
}

func NewStateExecutionIdManager() *StateExecutionIdManager {
	return &StateExecutionIdManager{
		stateIdMap: make(map[string]int),
	}
}

func (sm *StateExecutionIdManager) IncAndGetNextExecutionId(stateId string) string {
	sm.stateIdMap[stateId]++
	id := sm.stateIdMap[stateId]
	return fmt.Sprintf("%v-%v", stateId, id)
}
