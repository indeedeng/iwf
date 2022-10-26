package interpreter

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"io/ioutil"
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
	if checkError(err, httpResp) {
		return nil, composeError(provider, "state start API failed", err, httpResp)
	}

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
	if checkError(err, httpResp) {
		return nil, composeError(provider, "state decide API failed", err, httpResp)
	}
	return resp, nil
}

func checkError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}

func composeError(provider ActivityProvider, errType string, err error, httpResp *http.Response) error {
	responseBody := "None"
	var statusCode int
	if httpResp != nil {
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			responseBody = "cannot read body from http response"
		} else {
			responseBody = string(body)
		}
		statusCode = httpResp.StatusCode
	}
	return provider.NewApplicationError(fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", statusCode, responseBody, err), errType)

}
