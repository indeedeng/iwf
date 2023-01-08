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

const stateApiFailErrorType = "stateStartApiFailed"
const stateApiInvalidResponseErrorType = "stateStartApiInvalidResponse"

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
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, stateApiFailErrorType, err, httpResp)
	}

	if err := checkResp(resp); err != nil {
		return nil, composeRespError(provider, stateApiInvalidResponseErrorType, err, resp)
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
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, "state decide API failed", err, httpResp)
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

func checkHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}

func composeHttpError(provider ActivityProvider, errType string, err error, httpResp *http.Response) error {
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

func checkResp(resp *iwfidl.WorkflowStateStartResponse) error {
	if resp == nil || resp.CommandRequest == nil {
		return fmt.Errorf("empty response or command request")
	}
	commandReq := resp.CommandRequest
	if len(commandReq.GetTimerCommands())+len(commandReq.GetSignalCommands())+len(commandReq.GetInterStateChannelCommands()) > 0 {
		dtt := commandReq.GetDeciderTriggerType()
		if dtt != iwfidl.ANY_COMMAND_COMPLETED && dtt != iwfidl.ALL_COMMAND_COMPLETED && dtt != iwfidl.ANY_COMMAND_COMBINATION_COMPLETED {
			return fmt.Errorf("unsupported decider trigger type %s", dtt)
		}
		if dtt == iwfidl.ANY_COMMAND_COMBINATION_COMPLETED {
			// every command must have an id for this type
			err := fmt.Errorf("ANY_COMMAND_COMBINATION_COMPLETED can only be used when every command has an commandId")
			for _, cmd := range commandReq.GetTimerCommands() {
				if cmd.GetCommandId() == "" {
					return err
				}
			}
			for _, cmd := range commandReq.GetSignalCommands() {
				if cmd.GetCommandId() == "" {
					return err
				}
			}
			for _, cmd := range commandReq.GetInterStateChannelCommands() {
				if cmd.GetCommandId() == "" {
					return err
				}
			}
		}
	}
	// NOTE: we don't require decider trigger type when there is no commands
	return nil
}

func composeRespError(provider ActivityProvider, errType string, err error, resp *iwfidl.WorkflowStateStartResponse) error {
	respStr, _ := resp.MarshalJSON()
	return provider.NewApplicationError(fmt.Sprintf("err msg: %v, response: %v", err, string(respStr)), errType)
}
