package interpreter

type ContinueAsNewCounter struct {
	executedStateApis  int32
	signalsReceived    int32
	syncUpdateReceived int32

	configer *WorkflowConfiger
	rootCtx  UnifiedContext
	provider WorkflowProvider
}

func NewContinueAsCounter(configer *WorkflowConfiger, rootCtx UnifiedContext, provider WorkflowProvider) *ContinueAsNewCounter {
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
	// Note: when threshold == 0, it means unlimited

	config := c.configer.Get()
	if config.GetContinueAsNewThreshold() == 0 {
		return false
	}
	totalOperations := c.signalsReceived + c.executedStateApis + c.syncUpdateReceived

	return totalOperations >= config.GetContinueAsNewThreshold()
}
