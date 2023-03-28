package interpreter

import "github.com/indeedeng/iwf/service"

type ContinueAsNewCounter struct {
	executedStateExecution          int
	signalsReceived                 int
	executedStateExecutionThreshold int
	signalsReceivedThreshold        int
}

func NewContinueAsCounter(config service.WorkflowConfig) *ContinueAsNewCounter {
	return &ContinueAsNewCounter{
		executedStateExecutionThreshold: config.ContinueAsNewThresholdExecutedStateExecution,
		signalsReceivedThreshold:        config.ContinueAsNewThresholdSignalsReceived,
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
	return (c.executedStateExecutionThreshold > 0 && c.executedStateExecution > c.executedStateExecutionThreshold) ||
		(c.signalsReceivedThreshold > 0 && c.signalsReceived > c.signalsReceivedThreshold)
}
