package interpreter

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"net/http"
)

func StateStart(ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput) (*iwfidl.WorkflowStateStartResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
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
		return nil, provider.NewApplicationError("state start API failed", "api failed", httpResp)
	}
	// TODO validate commandId here

	return resp, nil
}

func StateDecide(ctx context.Context, backendType service.BackendType, input service.StateDecideActivityInput) (*iwfidl.WorkflowStateDecideResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateDecideActivity", "input", input)

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
		return nil, provider.NewApplicationError("state decide API failed", "api failed", httpResp)
	}
	return resp, nil
}
