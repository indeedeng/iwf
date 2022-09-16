package temporal

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"net/http"
)

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
