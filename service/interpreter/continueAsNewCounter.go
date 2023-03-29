package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

type ContinueAsNewCounter struct {
	executedStateExecution          int32
	signalsReceived                 int32
	executedStateExecutionThreshold int32
	signalsReceivedThreshold        int32
}

func NewContinueAsCounter(config iwfidl.WorkflowConfig) *ContinueAsNewCounter {
	return &ContinueAsNewCounter{
		executedStateExecutionThreshold: config.GetContinueAsNewThresholdExecutedStateExecution(),
		signalsReceivedThreshold:        config.GetContinueAsNewThresholdSignalsReceived(),
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
