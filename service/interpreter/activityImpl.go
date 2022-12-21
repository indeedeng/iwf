package interpreter

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func StateStart(ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput) (*iwfidl.WorkflowStateStartResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateStartActivity", "input", input)
	iwfWorkerBaseUrl := getIwfWorkerBaseUrlWithFix(input.IwfWorkerUrl)

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})
	req := apiClient.DefaultApi.ApiV1WorkflowStateStartPost(ctx)
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

	iwfWorkerBaseUrl := getIwfWorkerBaseUrlWithFix(input.IwfWorkerUrl)
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})
	req := apiClient.DefaultApi.ApiV1WorkflowStateDecidePost(ctx)
	resp, httpResp, err := req.WorkflowStateDecideRequest(input.Request).Execute()
	if checkError(err, httpResp) {
		return nil, composeError(provider, "state decide API failed", err, httpResp)
	}
	return resp, nil
}

func getIwfWorkerBaseUrlWithFix(url string) string {
	autofix := os.Getenv("AUTO_FIX_WORKER_URL")
	if autofix != "" {
		url = strings.Replace(url, "localhost", autofix, 1)
		url = strings.Replace(url, "127.0.0.1", autofix, 1)
	}
	return url
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
