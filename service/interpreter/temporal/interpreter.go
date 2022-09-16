package temporal

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"go.temporal.io/sdk/temporal"
	"net/http"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"

	// TODO(cretz): Remove when tagged
	_ "go.temporal.io/sdk/contrib/tools/workflowcheck/determinism"
)

func Interpreter(ctx workflow.Context, input service.InterpreterWorkflowInput) (*service.InterpreterWorkflowOutput, error) {
	currentStates := []iwfidl.StateMovement{
		{
			StateId:          &input.StartStateId,
			NextStateOptions: &input.StateOptions,
			NextStateInput:   &input.StateInput,
		},
	}

	var err error
	for len(currentStates) > 0 {
		statesToExecute := currentStates
		//reset to empty slice since each iteration will process all current states in the queue
		currentStates = nil

		for _, state := range statesToExecute {
			decision, err := executeState(state)
			if err != nil {
				return nil, err
			}
			// TODO process search attributes
			// TODO process query attributes

			isClosing, output, err := checkClosingWorkflow(decision)
			if isClosing {
				return output, err
			}
			if decision.HasNextStates() {
				currentStates = append(currentStates, decision.GetNextStates()...)
			}

		}
		err = workflow.Await(ctx, func() bool {
			return len(currentStates) > 0
		})
		if err != nil {
			break
		}
	}

	return nil, err
}

func checkClosingWorkflow(decision *iwfidl.StateDecision) (bool, *service.InterpreterWorkflowOutput, error) {
	hasClosingDecision := false
	var output *service.InterpreterWorkflowOutput
	for _, movement := range decision.GetNextStates() {
		stateId := movement.GetStateId()
		if stateId == service.CompletingWorkflowStateId || stateId == service.FailingWorkflowStateId {
			hasClosingDecision = true
			output = &service.InterpreterWorkflowOutput{
				CompletedStateId: stateId,
				StateOutput:      movement.GetNextStateInput(),
			}
		}
	}
	if hasClosingDecision && len(decision.NextStates) > 1 {
		// Illegal decision, should fail the workflow
		return true, nil, temporal.NewApplicationError(
			"closing workflow decision shouldn't have other state movements",
			"Illegal closing workflow decision",
		)
	}
	return hasClosingDecision, output, nil
}

func executeState(state iwfidl.StateMovement) (*iwfidl.StateDecision, error) {
	return nil, nil
}

func StateStartActivity(ctx context.Context, input service.StateStartActivityInput) (*iwfidl.WorkflowStateStartResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("StateStartActivity", "input", input)

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: input.IwfWorkerUrl,
			},
		},
	})
	req := apiClient.DefaultApi.ApiV1WorkflowStateStartPost(context.Background())
	resp, httpResp, err := req.WorkflowStateStartRequest(input.Request).Execute()
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, temporal.NewApplicationError("state start API failed", "api failed", httpResp)
	}
	return resp, nil
}

func StateDecideActivity(ctx context.Context, input service.StateDecideActivityInput) (*iwfidl.WorkflowStateDecideResponse, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("StateStartActivity", "input", input)

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: input.IwfWorkerUrl,
			},
		},
	})
	req := apiClient.DefaultApi.ApiV1WorkflowStateDecidePost(context.Background())
	resp, httpResp, err := req.WorkflowStateDecideRequest(input.Request).Execute()
	if err != nil {
		return nil, err
	}
	if httpResp.StatusCode != http.StatusOK {
		return nil, temporal.NewApplicationError("state decide API failed", "api failed", httpResp)
	}
	return resp, nil
}
