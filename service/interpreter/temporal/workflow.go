package temporal

import (
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/interpreter"
	"go.temporal.io/sdk/workflow"
	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	return interpreter.InterpreterImpl(interpreter.NewUnifiedContext(ctx), DefaultWorkflowProvider, input)
}
