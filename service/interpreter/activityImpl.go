package interpreter

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/config"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// StateStart is Deprecated, will be removed in next release
func StateStart(ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput) (*iwfidl.WorkflowStateStartResponse, error) {
	return StateApiWaitUntil(ctx, backendType, input)
}

func StateApiWaitUntil(ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput) (*iwfidl.WorkflowStateStartResponse, error) {
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

	attempt := provider.GetActivityInfo(ctx).Attempt
	scheduledTs := provider.GetActivityInfo(ctx).ScheduledTime.Unix()
	input.Request.Context.Attempt = &attempt
	input.Request.Context.FirstAttemptTimestamp = &scheduledTs

	req := apiClient.DefaultApi.ApiV1WorkflowStateStartPost(ctx)
	resp, httpResp, err := req.WorkflowStateStartRequest(input.Request).Execute()
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, err, httpResp, string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE))
	}

	if err := checkResp(resp); err != nil {
		return nil, composeStartApiRespError(provider, err, resp)
	}

	return resp, nil
}

// StateDecide is deprecated. Will be removed in next release
func StateDecide(ctx context.Context, backendType service.BackendType, input service.StateDecideActivityInput) (*iwfidl.WorkflowStateDecideResponse, error) {
	return StateApiExecute(ctx, backendType, input)
}

func StateApiExecute(ctx context.Context, backendType service.BackendType, input service.StateDecideActivityInput) (*iwfidl.WorkflowStateDecideResponse, error) {
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

	attempt := provider.GetActivityInfo(ctx).Attempt
	scheduledTs := provider.GetActivityInfo(ctx).ScheduledTime.Unix()
	input.Request.Context.Attempt = &attempt
	input.Request.Context.FirstAttemptTimestamp = &scheduledTs

	req := apiClient.DefaultApi.ApiV1WorkflowStateDecidePost(ctx)
	resp, httpResp, err := req.WorkflowStateDecideRequest(input.Request).Execute()
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, err, httpResp, string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE))
	}
	return resp, nil
}

func getIwfWorkerBaseUrlWithFix(url string) string {
	autofixUrl := os.Getenv("AUTO_FIX_WORKER_URL")
	if autofixUrl != "" {
		url = strings.Replace(url, "localhost", autofixUrl, 1)
		url = strings.Replace(url, "127.0.0.1", autofixUrl, 1)
	}
	autofixPortEnv := os.Getenv("AUTO_FIX_WORKER_PORT_FROM_ENV")
	if autofixPortEnv != "" {
		envVal := os.Getenv(autofixPortEnv)
		url = strings.Replace(url, "$"+autofixPortEnv+"$", envVal, 1)
	}

	return url
}

func composeStartApiRespError(provider ActivityProvider, err error, resp *iwfidl.WorkflowStateStartResponse) error {
	respStr, _ := resp.MarshalJSON()
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE),
		fmt.Sprintf("err msg: %v, response: %v", err, string(respStr)))
}

func checkHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}

func composeHttpError(provider ActivityProvider, err error, httpResp *http.Response, errType string) error {
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
	return provider.NewApplicationError(errType,
		fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", statusCode, responseBody, err))
}

func checkResp(resp *iwfidl.WorkflowStateStartResponse) error {
	if resp == nil || resp.CommandRequest == nil {
		return nil
	}
	commandReq := resp.CommandRequest
	if len(commandReq.GetTimerCommands())+len(commandReq.GetSignalCommands())+len(commandReq.GetInterStateChannelCommands()) > 0 {
		dtt := commandReq.GetDeciderTriggerType()
		if dtt != iwfidl.ANY_COMMAND_COMPLETED && dtt != iwfidl.ALL_COMMAND_COMPLETED && dtt != iwfidl.ANY_COMMAND_COMBINATION_COMPLETED {
			return fmt.Errorf("unsupported decider trigger type %s", dtt)
		}
		if dtt == iwfidl.ANY_COMMAND_COMBINATION_COMPLETED {
			// every command must have an id for this type
			err := fmt.Errorf("ANY_COMMAND_COMBINATION_COMPLETED can only be used when every command has an commandId, and the combination list cannot be empty")
			if len(commandReq.GetCommandCombinations()) == 0 {
				return err
			}
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

func DumpWorkflowInternal(ctx context.Context, backendType service.BackendType, req iwfidl.WorkflowDumpRequest) (*iwfidl.WorkflowDumpResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("DumpWorkflowInternal", "input", req)

	apiAddress := config.GetApiServiceAddressWithDefault(GetSharedConfig())

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: apiAddress,
			},
		},
	})

	request := apiClient.DefaultApi.ApiV1WorkflowInternalDumpPost(ctx)
	resp, httpResp, err := request.WorkflowDumpRequest(req).Execute()
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, err, httpResp, string(iwfidl.SERVER_INTERNAL_ERROR_TYPE))
	}
	return resp, nil
}
