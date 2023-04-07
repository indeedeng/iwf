package interpreter

type ContinueAsNewCounter struct {
	executedStateExecution int32
	signalsReceived        int32

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

func (c *ContinueAsNewCounter) IncExecutedStateExecution() {
	c.executedStateExecution++
}
func (c *ContinueAsNewCounter) IncSignalsReceived() {
	c.signalsReceived++
}

func (c *ContinueAsNewCounter) IsThresholdMet() bool {
	// Note: when threshold == 0, it means unlimited

	config := c.configer.Get()
	if config.GetContinueAsNewThreshold() == 0 {
		return false
	}
	totalOperations := c.signalsReceived + c.executedStateExecution*2

	return totalOperations >= config.GetContinueAsNewThreshold()
}
