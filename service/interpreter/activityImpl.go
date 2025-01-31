package interpreter

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/event"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/common/rpc"
	"github.com/indeedeng/iwf/service/common/urlautofix"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"io"
	"net/http"
	"os"
	"slices"
	"time"
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
	stateApiWaitUntilStartTime := time.Now().UnixMilli()
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateWaitUntilActivity", "input", log.ToJsonAndTruncateForLogging(input))
	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(input.IwfWorkerUrl)

	svcCfg := env.GetSharedConfig()
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		DefaultHeader: svcCfg.Interpreter.InterpreterActivityConfig.DefaultHeaders,
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})

	activityInfo := provider.GetActivityInfo(ctx)
	attempt := activityInfo.Attempt
	scheduledTs := activityInfo.ScheduledTime.Unix()
	input.Request.Context.Attempt = &attempt
	input.Request.Context.FirstAttemptTimestamp = &scheduledTs

	req := apiClient.DefaultApi.ApiV1WorkflowStateStartPost(ctx)
	resp, httpResp, err := req.WorkflowStateStartRequest(input.Request).Execute()
	printDebugMsg(logger, err, iwfWorkerBaseUrl)
	if checkHttpError(err, httpResp) {
		event.Handle(iwfidl.IwfEvent{
			EventType:        iwfidl.STATE_WAIT_UNTIL_ATTEMPT_FAIL_EVENT,
			WorkflowType:     input.Request.WorkflowType,
			WorkflowId:       activityInfo.WorkflowExecution.ID,
			WorkflowRunId:    activityInfo.WorkflowExecution.RunID,
			StateId:          ptr.Any(input.Request.WorkflowStateId),
			StateExecutionId: ptr.Any(input.Request.Context.GetStateExecutionId()),
		})
		return nil, composeHttpError(
			activityInfo.IsLocalActivity,
			provider, err, httpResp, string(iwfidl.STATE_API_FAIL_ERROR_TYPE))
	}

	if err := checkCommandRequestFromWaitUntilResponse(resp); err != nil {
		event.Handle(iwfidl.IwfEvent{
			EventType:        iwfidl.STATE_WAIT_UNTIL_ATTEMPT_FAIL_EVENT,
			WorkflowType:     input.Request.WorkflowType,
			WorkflowId:       activityInfo.WorkflowExecution.ID,
			WorkflowRunId:    activityInfo.WorkflowExecution.RunID,
			StateId:          ptr.Any(input.Request.WorkflowStateId),
			StateExecutionId: ptr.Any(input.Request.Context.GetStateExecutionId()),
		})
		return nil, composeStartApiRespError(provider, err, resp)
	}

	// Before returning successful results, check if it's local activity then compose some info for debug purpose
	// This is because local activity doesn't record input into the history.
	// But there are some small info that are important to record
	if activityInfo.IsLocalActivity {
		resp.LocalActivityInput = composeInputForDebug(input.Request.Context.GetStateExecutionId())
	}

	event.Handle(iwfidl.IwfEvent{
		EventType:          iwfidl.STATE_WAIT_UNTIL_ATTEMPT_SUCC_EVENT,
		WorkflowType:       input.Request.WorkflowType,
		WorkflowId:         activityInfo.WorkflowExecution.ID,
		WorkflowRunId:      activityInfo.WorkflowExecution.RunID,
		StateId:            ptr.Any(input.Request.WorkflowStateId),
		StateExecutionId:   ptr.Any(input.Request.Context.GetStateExecutionId()),
		StartTimestampInMs: ptr.Any(stateApiWaitUntilStartTime),
		EndTimestampInMs:   ptr.Any(time.Now().UnixMilli()),
	})
	return resp, nil
}

// StateDecide is deprecated. Will be removed in next release
func StateDecide(
	ctx context.Context,
	backendType service.BackendType,
	input service.StateDecideActivityInput,
) (*iwfidl.WorkflowStateDecideResponse, error) {
	return StateApiExecute(ctx, backendType, input)
}

func StateApiExecute(
	ctx context.Context,
	backendType service.BackendType,
	input service.StateDecideActivityInput,
) (*iwfidl.WorkflowStateDecideResponse, error) {
	stateApiExecuteStartTime := time.Now().UnixMilli()
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("StateExecuteActivity", "input", log.ToJsonAndTruncateForLogging(input))

	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(input.IwfWorkerUrl)
	svcCfg := env.GetSharedConfig()
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		DefaultHeader: svcCfg.Interpreter.InterpreterActivityConfig.DefaultHeaders,
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})

	activityInfo := provider.GetActivityInfo(ctx)
	attempt := activityInfo.Attempt
	scheduledTs := activityInfo.ScheduledTime.Unix()
	input.Request.Context.Attempt = &attempt
	input.Request.Context.FirstAttemptTimestamp = &scheduledTs

	req := apiClient.DefaultApi.ApiV1WorkflowStateDecidePost(ctx)
	resp, httpResp, err := req.WorkflowStateDecideRequest(input.Request).Execute()
	printDebugMsg(logger, err, iwfWorkerBaseUrl)
	if checkHttpError(err, httpResp) {
		event.Handle(iwfidl.IwfEvent{
			EventType:        iwfidl.STATE_EXECUTE_ATTEMPT_FAIL_EVENT,
			WorkflowType:     input.Request.WorkflowType,
			WorkflowId:       activityInfo.WorkflowExecution.ID,
			WorkflowRunId:    activityInfo.WorkflowExecution.RunID,
			StateId:          ptr.Any(input.Request.WorkflowStateId),
			StateExecutionId: input.Request.Context.StateExecutionId,
		})
		return nil, composeHttpError(
			activityInfo.IsLocalActivity,
			provider, err, httpResp, string(iwfidl.STATE_API_FAIL_ERROR_TYPE))
	}

	if err = checkStateDecisionFromResponse(resp); err != nil {
		event.Handle(iwfidl.IwfEvent{
			EventType:        iwfidl.STATE_EXECUTE_ATTEMPT_FAIL_EVENT,
			WorkflowType:     input.Request.WorkflowType,
			WorkflowId:       activityInfo.WorkflowExecution.ID,
			WorkflowRunId:    activityInfo.WorkflowExecution.RunID,
			StateId:          ptr.Any(input.Request.WorkflowStateId),
			StateExecutionId: input.Request.Context.StateExecutionId,
		})
		return nil, composeExecuteApiRespError(provider, err, resp)
	}

	// Before returning successful results, check if it's local activity then compose some info for debug purpose
	// This is because local activity doesn't record input into the history.
	// But there are some small info that are important to record
	if activityInfo.IsLocalActivity {
		resp.LocalActivityInput = composeInputForDebug(input.Request.Context.GetStateExecutionId())
	}

	event.Handle(iwfidl.IwfEvent{
		EventType:          iwfidl.STATE_EXECUTE_ATTEMPT_SUCC_EVENT,
		WorkflowType:       input.Request.WorkflowType,
		WorkflowId:         activityInfo.WorkflowExecution.ID,
		WorkflowRunId:      activityInfo.WorkflowExecution.RunID,
		StateId:            ptr.Any(input.Request.WorkflowStateId),
		StateExecutionId:   input.Request.Context.StateExecutionId,
		StartTimestampInMs: ptr.Any(stateApiExecuteStartTime),
		EndTimestampInMs:   ptr.Any(time.Now().UnixMilli()),
	})
	return resp, nil
}

func composeInputForDebug(stateExeId string) *string {
	// NOTE: only use the stateExecutionId for now, but we can add more later if needed
	return iwfidl.PtrString(fmt.Sprintf("stateExeId: %s", stateExeId))
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
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_ERROR_TYPE),
		fmt.Sprintf("err msg: %v, response: %v", err, string(respStr)))
}

func composeExecuteApiRespError(provider ActivityProvider, err error, resp *iwfidl.WorkflowStateDecideResponse) error {
	respStr, _ := resp.MarshalJSON()
	return provider.NewApplicationError(string(iwfidl.STATE_API_FAIL_ERROR_TYPE),
		fmt.Sprintf("err msg: %v, response: %v", err, string(respStr)))
}

func checkHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}

func composeHttpError(
	isLocalActivity bool, provider ActivityProvider, err error, httpResp *http.Response, errType string,
) error {
	responseBody := "None"
	var statusCode int
	if httpResp != nil {
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			responseBody = "cannot read body from http response"
		} else {
			responseBody = string(body)
		}
		statusCode = httpResp.StatusCode
	}
	errMsg := err.Error()
	var trimmedResponseBody, trimmedErrMsg string
	if isLocalActivity {
		trimmedErrMsg = trimText(errMsg, 5)
		trimmedResponseBody = trimText(responseBody, 50)
		errType = "1st-attempt-failure"
	} else {
		trimmedErrMsg = trimText(errMsg, 50)
		trimmedResponseBody = trimText(responseBody, 500)
	}

	return provider.NewApplicationError(errType,
		fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", statusCode, trimmedResponseBody, trimmedErrMsg))
}

func trimText(msg string, maxLength int) string {
	if len(msg) > maxLength {
		return msg[:maxLength] + "..."
	}
	return msg
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
			// Check if each command in the combinations has a matching command in one of the lists
			if !areAllCommandCombinationsIdsValid(commandReq) {
				return fmt.Errorf("ANY_COMMAND_COMBINATION_COMPLETED can only be used when every command has an commandId that is found in TimerCommands, SignalCommands or InternalChannelCommand")
			}
		}
	}
	// NOTE: we don't require decider trigger type when there is no commands
	return nil
}

func areAllCommandCombinationsIdsValid(commandReq *iwfidl.CommandRequest) bool {
	timerSignalInternalChannelCmdIds := listTimerSignalInternalChannelCommandIds(commandReq)
	for _, commandCombo := range commandReq.GetCommandCombinations() {
		for _, cmdId := range commandCombo.GetCommandIds() {
			if !slices.Contains(timerSignalInternalChannelCmdIds, cmdId) {
				return false
			}
		}
	}
	return true
}

func listTimerSignalInternalChannelCommandIds(commandReq *iwfidl.CommandRequest) []string {
	var ids []string
	for _, timerCmd := range commandReq.GetTimerCommands() {
		ids = append(ids, timerCmd.GetCommandId())
	}
	for _, signalCmd := range commandReq.GetSignalCommands() {
		ids = append(ids, signalCmd.GetCommandId())
	}
	for _, internalChannelCmd := range commandReq.GetInterStateChannelCommands() {
		ids = append(ids, internalChannelCmd.GetCommandId())
	}
	return ids
}

func DumpWorkflowInternal(
	ctx context.Context, backendType service.BackendType, req iwfidl.WorkflowDumpRequest,
) (*iwfidl.WorkflowDumpResponse, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("DumpWorkflowInternalActivity", "input", log.ToJsonAndTruncateForLogging(req))

	svcCfg := env.GetSharedConfig()
	apiAddress := svcCfg.GetApiServiceAddressWithDefault()

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		DefaultHeader: svcCfg.Interpreter.InterpreterActivityConfig.DefaultHeaders,
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: apiAddress,
			},
		},
	})

	request := apiClient.DefaultApi.ApiV1WorkflowInternalDumpPost(ctx)
	resp, httpResp, err := request.WorkflowDumpRequest(req).Execute()
	if checkHttpError(err, httpResp) {
		return nil, composeHttpError(provider.GetActivityInfo(ctx).IsLocalActivity,
			provider, err, httpResp, string(iwfidl.SERVER_INTERNAL_ERROR_TYPE))
	}
	return resp, nil
}

func InvokeWorkerRpc(
	ctx context.Context, backendType service.BackendType, rpcPrep *service.PrepareRpcQueryResponse,
	req iwfidl.WorkflowRpcRequest,
) (*InvokeRpcActivityOutput, error) {
	provider := getActivityProviderByType(backendType)
	logger := provider.GetLogger(ctx)
	logger.Info("InvokeWorkerRpcActivity", "input", log.ToJsonAndTruncateForLogging(req))

	apiMaxSeconds := env.GetSharedConfig().Api.MaxWaitSeconds

	resp, statusErr := rpc.InvokeWorkerRpc(ctx, rpcPrep, req, apiMaxSeconds)
	return &InvokeRpcActivityOutput{
		RpcOutput:   resp,
		StatusError: statusErr,
	}, nil
}
