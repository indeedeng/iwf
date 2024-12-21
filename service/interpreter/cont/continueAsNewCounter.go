package cont

import (
	"github.com/indeedeng/iwf/service/interpreter/config"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
)

type ContinueAsNewCounter struct {
	executedStateApis  int32
	signalsReceived    int32
	syncUpdateReceived int32
	triggeredByAPI     bool

	configer *config.WorkflowConfiger
	rootCtx  interfaces.UnifiedContext
	provider interfaces.WorkflowProvider
}

func NewContinueAsCounter(
	configer *config.WorkflowConfiger, rootCtx interfaces.UnifiedContext, provider interfaces.WorkflowProvider,
) *ContinueAsNewCounter {
	return &ContinueAsNewCounter{
		configer: configer,

		rootCtx:  rootCtx,
		provider: provider,
	}
}

func (c *ContinueAsNewCounter) IncExecutedStateExecution(skipStart bool) {
	if skipStart {
		c.executedStateApis++
	} else {
		c.executedStateApis += 2
	}
}
func (c *ContinueAsNewCounter) IncSignalsReceived() {
	c.signalsReceived++
}

func (c *ContinueAsNewCounter) IncSyncUpdateReceived() {
	c.syncUpdateReceived++
}

func (c *ContinueAsNewCounter) IsThresholdMet() bool {
	if c.triggeredByAPI {
		return true
	}

	// Note: when threshold == 0, it means unlimited

	config := c.configer.Get()
	if config.GetContinueAsNewThreshold() == 0 {
		return false
	}
	totalOperations := c.signalsReceived + c.executedStateApis + c.syncUpdateReceived

	return totalOperations >= config.GetContinueAsNewThreshold()
}

func (c *ContinueAsNewCounter) TriggerByAPI() {
	c.triggeredByAPI = true
}
