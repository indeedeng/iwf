package temporal

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/client/workflow/state"
	"github.com/cadence-oss/iwf-server/service"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

const TaskQueue = "Interpreter"

type stateExecution struct {
	stateId      string
	stateInput   state.EncodedObject
	stateOptions state.WorkflowStateOptions
}

// Interpreter is a interpreter workflow definition.
func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Interpreter workflow started", "input", input)

	var result string
	err := workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return nil, err
	}

	logger.Info("Interpreter workflow completed.", "result", result)

	return nil, nil
}

func Activity(ctx context.Context, input service.InterpreterWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", input)
	return "Hello " + input.StartStateId + "!", nil
}
