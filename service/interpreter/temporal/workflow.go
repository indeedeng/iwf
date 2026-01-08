package temporal

import (
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/interpreter"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	return interpreter.InterpreterImpl(interfaces.NewUnifiedContext(ctx), newTemporalWorkflowProvider(), input)
}

func WaitforStateCompletionWorkflow(ctx workflow.Context) (*service.WaitForStateCompletionWorkflowOutput, error) {
	return interpreter.WaitForStateCompletionWorkflowImpl(interfaces.NewUnifiedContext(ctx), newTemporalWorkflowProvider())
}

func BlobStoreCleanup(ctx workflow.Context, storeId string) (int, error) {
	return interpreter.BlobStoreCleanup(interfaces.NewUnifiedContext(ctx), newTemporalWorkflowProvider(), storeId)
}
