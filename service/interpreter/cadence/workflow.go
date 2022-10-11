package cadence

import (
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter"
	"go.uber.org/cadence/workflow"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	return interpreter.InterpreterImpl(interpreter.NewUnifiedContext(ctx), defaultWorkflowProvider, input)
}
