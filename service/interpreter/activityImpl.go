package interpreter

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/rpc"
	"github.com/indeedeng/iwf/service/common/urlautofix"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"io/ioutil"
	"net/http"
	"os"
)

// StateStart is Deprecated, will be removed in next release
func StateStart(
	ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput,
) (*iwfidl.WorkflowStateStartResponse, error) {
	return StateApiWaitUntil(ctx, backendType, input)
}

func StateApiWaitUntil(
	ctx context.Context, backendType service.BackendType, input service.StateStartActivityInput,
) (*iwfidl.WorkflowStateStartResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateStartActivity", "input", input)
	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(input.IwfWorkerUrl)

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
	printDebugMsg(logger, err, iwfWorkerBaseUrl)
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, err, httpResp, string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE))
	}

	if err := checkCommandRequestFromWaitUntilResponse(resp); err != nil {
		return nil, composeStartApiRespError(provider, err, resp)
	}

	return resp, nil
}

// StateDecide is deprecated. Will be removed in next release
func StateDecide(
	ctx context.Context,
	backendType service.BackendType,
	input service.StateDecideActivityInput,
	shouldSendSignalOnCompletion bool,
	timeout int,
) (*iwfidl.WorkflowStateDecideResponse, error) {
	return StateApiExecute(ctx, backendType, input, shouldSendSignalOnCompletion, timeout)
}

func StateApiExecute(
	ctx context.Context,
	backendType service.BackendType,
	input service.StateDecideActivityInput,
	_ bool, // no used anymore, keep for compatibility
	_ int, // no used anymore, keep for compatibility
) (*iwfidl.WorkflowStateDecideResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateDecideActivity", "input", input)

	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(input.IwfWorkerUrl)
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
	printDebugMsg(logger, err, iwfWorkerBaseUrl)
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider, err, httpResp, string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE))
	}

	if err = checkStateDecisionFromResponse(resp); err != nil {
		return nil, composeExecuteApiRespError(provider, err, resp)
	}

	return resp, nil
}

func checkStateDecisionFromResponse(resp *iwfidl.WorkflowStateDecideResponse) error {
	if resp == nil || resp.StateDecision == nil || len(resp.StateDecision.NextStates) == 0 {
		return fmt.Errorf("empty state decision is no longer supported. If it's from old SDKs then upgrade the SDK to newer versions")
	}
	return nil
}

func printDebugMsg(logger UnifiedLogger, err error, url string) {
	debugMode := os.Getenv(service.EnvNameDebugMode)
	if debugMode != "" {
		logger.Info("check error at http request", err, url)
	}
}

func composeStartApiRespError(provider ActivityProvider, err error, resp *iwfidl.WorkflowStateStartResponse) error {
	respStr, _ := resp.MarshalJSON()
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE),
		fmt.Sprintf("err msg: %v, response: %v", err, string(respStr)))
}

func composeExecuteApiRespError(provider ActivityProvider, err error, resp *iwfidl.WorkflowStateDecideResponse) error {
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

func checkCommandRequestFromWaitUntilResponse(resp *iwfidl.WorkflowStateStartResponse) error {
	if resp == nil || resp.CommandRequest == nil {
		return nil
	}
	commandReq := resp.CommandRequest
	if len(commandReq.GetTimerCommands())+len(commandReq.GetSignalCommands())+len(commandReq.GetInterStateChannelCommands()) > 0 {
		dtt := compatibility.GetDeciderTriggerType(*commandReq)
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

func DumpWorkflowInternal(
	ctx context.Context, backendType service.BackendType, req iwfidl.WorkflowDumpRequest,
) (*iwfidl.WorkflowDumpResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("DumpWorkflowInternal", "input", req)

	apiAddress := config.GetApiServiceAddressWithDefault(env.GetSharedConfig())

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

func InvokeWorkerRpc(
	ctx context.Context, backendType service.BackendType, rpcPrep *service.PrepareRpcQueryResponse,
	req iwfidl.WorkflowRpcRequest,
) (*InvokeRpcActivityOutput, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("invoke worker RPC by activity", "input", req)

	apiMaxSeconds := env.GetSharedConfig().Api.MaxWaitSeconds

	resp, statusErr := rpc.InvokeWorkerRpc(ctx, rpcPrep, req, apiMaxSeconds)
	return &InvokeRpcActivityOutput{
		RpcOutput:   resp,
		StatusError: statusErr,
	}, nil
}
