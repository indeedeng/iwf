package cadence

import (
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"go.uber.org/cadence/workflow"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	return interpreter.InterpreterImpl(interfaces.NewUnifiedContext(ctx), newCadenceWorkflowProvider(), input)
}

func WaitforStateCompletionWorkflow(ctx workflow.Context) (*service.WaitForStateCompletionWorkflowOutput, error) {
	return interpreter.WaitForStateCompletionWorkflowImpl(interfaces.NewUnifiedContext(ctx), newCadenceWorkflowProvider())
}

func BlobStoreCleanup(ctx workflow.Context, storeId string) error {
	return interpreter.BlobStoreCleanup(interfaces.NewUnifiedContext(ctx), newCadenceWorkflowProvider(), storeId)
}
