package temporal

import (
	"context"
	"fmt"

	"github.com/cadence-oss/iwf-server/gen/client/workflow/state"
	iwf "github.com/cadence-oss/iwf-server/gen/server/workflow"
	"github.com/cadence-oss/iwf-server/service"

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
	currentStates := []stateExecution{
		{
			stateId:      input.StartStateId,
			stateOptions: convertStateOption(input.StateOptions),
			stateInput:   convertStateInput(input.StateInput),
		},
	}

	for len(currentStates) > 0 {
		statesToExecute := currentStates
		currentStates = nil //reset to empty slice
		for _, state := range statesToExecute {

		}
		workflow.Await(func() bool {
			len(currentStates) > 0
		})
	}

	return nil, fmt.Errorf("we should never run into this line")
	// ao := workflow.ActivityOptions{
	// 	StartToCloseTimeout: 10 * time.Second,
	// }
	// ctx = workflow.WithActivityOptions(ctx, ao)

	// logger := workflow.GetLogger(ctx)
	// logger.Info("Interpreter workflow started", "input", input)

	// var result string
	// err := workflow.ExecuteActivity(ctx, Activity, input).Get(ctx, &result)
	// if err != nil {
	// 	logger.Error("Activity failed.", "Error", err)
	// 	return nil, err
	// }

	// logger.Info("Interpreter workflow completed.", "result", result)

	//return nil, nil
}

func convertStateInput(input iwf.EncodedObject) state.EncodedObject {
	return state.EncodedObject{
		Data:     &input.Data,
		Encoding: &input.Encoding,
	}
}

func convertStateOption(options iwf.WorkflowStateOptions) state.WorkflowStateOptions {
	return state.WorkflowStateOptions{
		SearchAttributesLoadingPolicy: convertAttributesLoadingPolicy(options.SearchAttributesLoadingPolicy),
		QueryAttributesLoadingPolicy:  convertAttributesLoadingPolicy(options.QueryAttributesLoadingPolicy),
		CommandCarryOverPolicy: &state.CommandCarryOverPolicy{
			CommandCarryOverType: &options.CommandCarryOverPolicy.CommandCarryOverType,
		},
	}
}

func convertAttributesLoadingPolicy(policy iwf.AttributesLoadingPolicy) *state.AttributesLoadingPolicy {
	return &state.AttributesLoadingPolicy{
		AttributeLoadingType: &policy.AttributeLoadingType,
		AttributeKeys:        policy.AttributeKeys,
	}
}

func Activity(ctx context.Context, input service.InterpreterWorkflowInput) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "input", input)
	return "Hello " + input.StartStateId + "!", nil
}