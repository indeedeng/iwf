package interpreter

import "github.com/indeedeng/iwf/gen/iwfidl"

type OutputCollector struct {
	outputs []iwfidl.StateCompletionOutput
}

func NewOutputCollector(initOutputs []iwfidl.StateCompletionOutput) *OutputCollector {
	if initOutputs == nil {
		initOutputs = []iwfidl.StateCompletionOutput{}
	}
	return &OutputCollector{
		outputs: initOutputs,
	}
}

func (o *OutputCollector) Add(output iwfidl.StateCompletionOutput) {
	if output.CompletedStateOutput != nil {
		o.outputs = append(o.outputs, output)
	}
}

func (o *OutputCollector) GetAll() []iwfidl.StateCompletionOutput {
	return o.outputs
}
