package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

type ContinueAsNewCounter struct {
	executedStateExecution          int32
	signalsReceived                 int32
	executedStateExecutionThreshold int32
	signalsReceivedThreshold        int32

	rootCtx  UnifiedContext
	provider WorkflowProvider
}

func NewContinueAsCounter(config iwfidl.WorkflowConfig, rootCtx UnifiedContext, provider WorkflowProvider) *ContinueAsNewCounter {
	return &ContinueAsNewCounter{
		executedStateExecutionThreshold: config.GetContinueAsNewThresholdExecutedStateExecution(),
		signalsReceivedThreshold:        config.GetContinueAsNewThresholdSignalsReceived(),

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

	isMet := (c.executedStateExecutionThreshold > 0 && c.executedStateExecution >= c.executedStateExecutionThreshold) ||
		(c.signalsReceivedThreshold > 0 && c.signalsReceived >= c.signalsReceivedThreshold)
	if isMet {
		c.provider.GetLogger(c.rootCtx).Info("continueAsNew condition is met", c.executedStateExecution, c.signalsReceived, "called at:"+LastCaller())
	}

	return isMet
}
