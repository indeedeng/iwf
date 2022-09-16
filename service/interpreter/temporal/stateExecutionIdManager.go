package temporal

import "fmt"

type stateExecutionIdManager struct {
	stateIdMap map[string]int
}

func newStateExecutionIdManager() *stateExecutionIdManager {
	return &stateExecutionIdManager{
		stateIdMap: make(map[string]int),
	}
}

func (sm *stateExecutionIdManager) incAndGetNextExecutionId(stateId string) string {
	sm.stateIdMap[stateId]++
	id := sm.stateIdMap[stateId]
	return fmt.Sprintf("%v-%v", stateId, id)
}
