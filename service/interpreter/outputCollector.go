package interpreter

import "github.com/indeedeng/iwf/gen/iwfidl"

type OutputCollector struct {
	outputs []iwfidl.StateCompletionOutput
}

func NewOutputCollector(initOutputs []iwfidl.StateCompletionOutput) *OutputCollector {
	filteredOutputs := []iwfidl.StateCompletionOutput{}

	if initOutputs == nil {
		return &OutputCollector{
			outputs: filteredOutputs,
		}
	} else {
		for _, output := range initOutputs {
			if output.CompletedStateOutput != nil {
				filteredOutputs = append(filteredOutputs, output)
			}
		}

		return &OutputCollector{
			outputs: filteredOutputs,
		}
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
