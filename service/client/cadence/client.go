package cadence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/indeedeng/iwf/config"
	"time"

	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/google/uuid"
	"github.com/indeedeng/iwf/gen/iwfidl"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	realcadence "go.uber.org/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/encoded"
)

type cadenceClient struct {
	domain                         string
	cClient                        client.Client
	closeFunc                      func()
	serviceClient                  workflowserviceclient.Interface
	converter                      encoded.DataConverter
	queryWorkflowFailedRetryPolicy config.QueryWorkflowFailedRetryPolicy
}

func (t *cadenceClient) IsWorkflowAlreadyStartedError(err error) bool {
	var workflowExecutionAlreadyStartedError *shared.WorkflowExecutionAlreadyStartedError
	ok := errors.As(err, &workflowExecutionAlreadyStartedError)
	return ok
}

func (t *cadenceClient) GetRunIdFromWorkflowAlreadyStartedError(err error) (string, bool) {
	var res *shared.WorkflowExecutionAlreadyStartedError
	ok := errors.As(err, &res)
	runId := ""
	if ok {
		runId = *res.RunId
	}
	return runId, ok
}

func (t *cadenceClient) IsNotFoundError(err error) bool {
	var entityNotExistsError *shared.EntityNotExistsError
	ok := errors.As(err, &entityNotExistsError)
	return ok
}

func (t *cadenceClient) isQueryFailedError(err error) bool {
	var serviceError *shared.QueryFailedError
	ok := errors.As(err, &serviceError)
	return ok
}

func (t *cadenceClient) IsWorkflowTimeoutError(err error) bool {
	return realcadence.IsTimeoutError(err)
}

func (t *cadenceClient) IsRequestTimeoutError(err error) bool {
	return errors.Is(err, context.DeadlineExceeded)
}

func (t *cadenceClient) GetApplicationErrorTypeIfIsApplicationError(err error) string {
	var cErr *realcadence.CustomError
	ok := errors.As(err, &cErr)
	if ok {
		return cErr.Reason()
	}
	return ""
}

func (t *cadenceClient) GetApplicationErrorDetails(err error, detailsPtr interface{}) error {
	var cErr *realcadence.CustomError
	ok := errors.As(err, &cErr)
	if ok {
		if cErr.HasDetails() {
			return cErr.Details(detailsPtr)
		}
		return fmt.Errorf("application error doesn't have details. Critical code bug")
	}
	return fmt.Errorf("not an application error. Critical code bug")
}

func (t *cadenceClient) GetApplicationErrorTypeAndDetails(err error) (string, string) {
	errType := t.GetApplicationErrorTypeIfIsApplicationError(err)

	var errDetailsPtr interface{}
	var errDetails string

	// Get error details into a generic interface{} pointer that can hold any type
	err2 := t.GetApplicationErrorDetails(err, &errDetailsPtr)
	if err2 != nil {
		errDetails = err2.Error()
	} else {
		// Check if the error details is a string
		errDetailsString, ok := errDetailsPtr.(string)
		// If it is a string, use it as the error details
		if ok {
			errDetails = errDetailsString
		} else {
			// For all other types, try to Marshal the object to JSON
			var err error
			jsonBytes, err := json.Marshal(errDetailsPtr)
			if err == nil {
				errDetails = string(jsonBytes)
			} else {
				// If Marshal fails, error message will say "couldn't parse the error details"
				errDetails = "couldn't parse error details to JSON. Critical code bug"
			}
		}
	}

	return errType, errDetails
}

func NewCadenceClient(
	domain string, cClient client.Client, serviceClient workflowserviceclient.Interface,
	converter encoded.DataConverter, closeFunc func(), retryPolicy *config.QueryWorkflowFailedRetryPolicy,
) uclient.UnifiedClient {
	return &cadenceClient{
		domain:                         domain,
		cClient:                        cClient,
		closeFunc:                      closeFunc,
		serviceClient:                  serviceClient,
		converter:                      converter,
		queryWorkflowFailedRetryPolicy: config.QueryWorkflowFailedRetryPolicyWithDefaults(retryPolicy),
	}
}

func (t *cadenceClient) Close() {
	t.closeFunc()
}

func (t *cadenceClient) StartInterpreterWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions, args ...interface{},
) (runId string, err error) {
	_, ok := options.Memo[service.UseMemoForDataAttributesKey]
	if ok {
		return "", fmt.Errorf("using Memo is not supported with Cadence, see https://github.com/uber/cadence/issues/3729")
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:                           options.ID,
		TaskList:                     options.TaskQueue,
		ExecutionStartToCloseTimeout: options.WorkflowExecutionTimeout,
		SearchAttributes:             options.SearchAttributes,
		Memo:                         options.Memo,
	}

	if options.WorkflowIDReusePolicy != nil {
		workflowIdReusePolicy, err := mapToCadenceWorkflowIdReusePolicy(*options.WorkflowIDReusePolicy)
		if err != nil {
			return "", nil
		}

		workflowOptions.WorkflowIDReusePolicy = *workflowIdReusePolicy
	}

	if options.CronSchedule != nil {
		workflowOptions.CronSchedule = *options.CronSchedule
	}

	if options.RetryPolicy != nil {
		workflowOptions.RetryPolicy = retry.ConvertCadenceWorkflowRetryPolicy(options.RetryPolicy)
	}

	if options.WorkflowStartDelay != nil {
		workflowOptions.DelayStart = *options.WorkflowStartDelay
	}

	run, err := t.cClient.StartWorkflow(ctx, workflowOptions, cadence.Interpreter, args...)
	if err != nil {
		return "", err
	}
	return run.RunID, nil
}

func (t *cadenceClient) StartWaitForStateCompletionWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions,
) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                           options.ID,
		TaskList:                     options.TaskQueue,
		WorkflowIDReusePolicy:        client.WorkflowIDReusePolicyAllowDuplicateFailedOnly, // the workflow could be timeout, so we allow duplicate
		ExecutionStartToCloseTimeout: options.WorkflowExecutionTimeout,
	}
	run, err := t.cClient.StartWorkflow(ctx, workflowOptions, cadence.WaitforStateCompletionWorkflow)
	if err != nil {
		if t.IsWorkflowAlreadyStartedError(err) {
			// if the workflow is already started, we return the runId
			return *err.(*shared.WorkflowExecutionAlreadyStartedError).RunId, nil
		}
		return "", err
	}
	return run.RunID, nil
}

func (t *cadenceClient) SignalWithStartWaitForStateCompletionWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions, stateCompletionOutput iwfidl.StateCompletionOutput,
) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:                           options.ID,
		TaskList:                     options.TaskQueue,
		WorkflowIDReusePolicy:        client.WorkflowIDReusePolicyAllowDuplicateFailedOnly, // the workflow could be timeout, so we allow duplicate
		ExecutionStartToCloseTimeout: options.WorkflowExecutionTimeout,
	}

	_, err := t.cClient.SignalWithStartWorkflow(ctx, options.ID, service.StateCompletionSignalChannelName, stateCompletionOutput, workflowOptions, cadence.WaitforStateCompletionWorkflow)
	if err != nil {
		return err
	}
	return nil
}

func (t *cadenceClient) SignalWorkflow(
	ctx context.Context, workflowID string, runID string, signalName string, arg interface{},
) error {
	return t.cClient.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

func (t *cadenceClient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	return t.cClient.CancelWorkflow(ctx, workflowID, runID)
}

func (t *cadenceClient) TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error {
	var reasonStr string
	if reason == "" {
		reasonStr = "Force termiantion from user"
	} else {
		reasonStr = reason
	}

	return t.cClient.TerminateWorkflow(ctx, workflowID, runID, reasonStr, nil)
}

func (t *cadenceClient) ListWorkflow(
	ctx context.Context, request *uclient.ListWorkflowExecutionsRequest,
) (*uclient.ListWorkflowExecutionsResponse, error) {
	listReq := &shared.ListWorkflowExecutionsRequest{
		PageSize:      &request.PageSize,
		Query:         &request.Query,
		NextPageToken: request.NextPageToken,
	}
	resp, err := t.cClient.ListWorkflow(ctx, listReq)
	if err != nil {
		return nil, err
	}
	var executions []iwfidl.WorkflowSearchResponseEntry
	for _, exe := range resp.GetExecutions() {
		executions = append(executions, iwfidl.WorkflowSearchResponseEntry{
			WorkflowId:    *exe.Execution.WorkflowId,
			WorkflowRunId: *exe.Execution.RunId,
		})
	}
	return &uclient.ListWorkflowExecutionsResponse{
		Executions:    executions,
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (t *cadenceClient) QueryWorkflow(
	ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{},
) error {
	var qres encoded.Value
	var err error

	attempt := 1
	// Only QueryFailed error causes retry; all other errors make the loop to finish immediately
	for attempt <= t.queryWorkflowFailedRetryPolicy.MaximumAttempts {
		qres, err = t.cClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
		if err == nil {
			break
		} else {
			if t.isQueryFailedError(err) {
				time.Sleep(time.Duration(t.queryWorkflowFailedRetryPolicy.InitialIntervalSeconds) * time.Second)
				attempt++
				continue
			}
			return err
		}
	}
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func queryWorkflowWithStrongConsistency(
	t *cadenceClient, ctx context.Context, workflowID string, runID string, queryType string, args []interface{},
) (encoded.Value, error) {
	queryWorkflowWithOptionsRequest := &client.QueryWorkflowWithOptionsRequest{
		WorkflowID:            workflowID,
		RunID:                 runID,
		QueryType:             queryType,
		Args:                  args,
		QueryConsistencyLevel: ptr.Any(shared.QueryConsistencyLevelStrong),
	}
	result, err := t.cClient.QueryWorkflowWithOptions(ctx, queryWorkflowWithOptionsRequest)
	if err != nil {
		return nil, err
	}
	return result.QueryResult, nil
}

func (t *cadenceClient) DescribeWorkflowExecution(
	ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType,
) (*uclient.DescribeWorkflowExecutionResponse, error) {
	resp, err := t.cClient.DescribeWorkflowExecution(ctx, workflowID, runID)
	if err != nil {
		return nil, err
	}
	status, err := mapToIwfWorkflowStatus(resp.GetWorkflowExecutionInfo().CloseStatus)
	if err != nil {
		return nil, err
	}
	searchAttributes, err := mapper.MapCadenceToIwfSearchAttributes(resp.GetWorkflowExecutionInfo().GetSearchAttributes(), requestedSearchAttributes)
	if err != nil {
		return nil, err
	}

	memo, err := t.decodeMemo(resp.GetWorkflowExecutionInfo().GetMemo())

	return &uclient.DescribeWorkflowExecutionResponse{
		RunId:            resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		FirstRunId:       "", // Cadence does not provide FirstRunId
		Status:           status,
		SearchAttributes: searchAttributes,
		Memos:            memo,
	}, nil
}

func (t *cadenceClient) decodeMemo(memo *shared.Memo) (map[string]iwfidl.EncodedObject, error) {
	if memo == nil || len(memo.GetFields()) == 0 {
		return nil, nil
	}

	out := map[string]iwfidl.EncodedObject{}
	for k, payload := range memo.GetFields() {
		var value iwfidl.EncodedObject
		err := encoded.GetDefaultDataConverter().FromData(payload, &value)
		if err != nil {
			return nil, err
		}
		out[k] = value
	}
	return out, nil
}

func mapToCadenceWorkflowIdReusePolicy(workflowIdReusePolicy iwfidl.WorkflowIDReusePolicy) (*client.WorkflowIDReusePolicy, error) {
	var res client.WorkflowIDReusePolicy
	switch workflowIdReusePolicy {
	case iwfidl.ALLOW_DUPLICATE:
		res = client.WorkflowIDReusePolicyAllowDuplicate
		return &res, nil
	case iwfidl.ALLOW_DUPLICATE_FAILED_ONLY:
		res = client.WorkflowIDReusePolicyAllowDuplicateFailedOnly
		return &res, nil
	case iwfidl.REJECT_DUPLICATE:
		res = client.WorkflowIDReusePolicyRejectDuplicate
		return &res, nil
	case iwfidl.TERMINATE_IF_RUNNING:
		res = client.WorkflowIDReusePolicyTerminateIfRunning
		return &res, nil
	default:
		return nil, fmt.Errorf("unsupported workflow id reuse policy %s", workflowIdReusePolicy)
	}
}

func mapToIwfWorkflowStatus(status *shared.WorkflowExecutionCloseStatus) (iwfidl.WorkflowStatus, error) {
	if status == nil {
		return iwfidl.RUNNING, nil
	}

	switch *status {
	case shared.WorkflowExecutionCloseStatusCanceled:
		return iwfidl.CANCELED, nil
	case shared.WorkflowExecutionCloseStatusContinuedAsNew:
		return iwfidl.CONTINUED_AS_NEW, nil
	case shared.WorkflowExecutionCloseStatusFailed:
		return iwfidl.FAILED, nil
	case shared.WorkflowExecutionCloseStatusTimedOut:
		return iwfidl.TIMEOUT, nil
	case shared.WorkflowExecutionCloseStatusTerminated:
		return iwfidl.TERMINATED, nil
	case shared.WorkflowExecutionCloseStatusCompleted:
		return iwfidl.COMPLETED, nil
	default:
		return "", fmt.Errorf("not supported status %s", status)
	}
}

func (t *cadenceClient) GetWorkflowResult(
	ctx context.Context, valuePtr interface{}, workflowID string, runID string,
) error {
	run := t.cClient.GetWorkflow(ctx, workflowID, runID)
	return run.Get(ctx, valuePtr)
}

func (t *cadenceClient) SynchronousUpdateWorkflow(
	ctx context.Context, valuePtr interface{}, workflowID, runID, updateType string, input interface{},
) error {
	return fmt.Errorf("not supported in Cadence")
}

func (t *cadenceClient) ResetWorkflow(
	ctx context.Context, request iwfidl.WorkflowResetRequest,
) (newRunId string, err error) {

	reqRunId := request.GetWorkflowRunId()
	if reqRunId == "" {
		// set default runId to current
		resp, err := t.cClient.DescribeWorkflowExecution(ctx, request.GetWorkflowId(), "")
		if err != nil {
			return "", err
		}
		reqRunId = resp.GetWorkflowExecutionInfo().GetExecution().GetRunId()
	}

	// TODO not sure why Cadence reset API requires this for GetWorkflowExecutionHistory API....
	ctx, cancelFn := context.WithTimeout(ctx, time.Second*120)
	defer cancelFn()

	resetType := request.GetResetType()
	resetBaseRunID, decisionFinishID, err := getResetIDsByType(ctx, resetType, t.domain, request.GetWorkflowId(),
		reqRunId, t.serviceClient, t.converter, request.GetHistoryEventId(), request.GetHistoryEventTime(), request.GetStateId(), request.GetStateExecutionId())

	if err != nil {
		return "", err
	}

	requestId := uuid.New().String()
	resetReq := &shared.ResetWorkflowExecutionRequest{
		Domain: &t.domain,
		WorkflowExecution: &shared.WorkflowExecution{
			WorkflowId: &request.WorkflowId,
			RunId:      &resetBaseRunID,
		},
		Reason:                request.Reason,
		DecisionFinishEventId: iwfidl.PtrInt64(decisionFinishID),
		RequestId:             &requestId,
		SkipSignalReapply:     iwfidl.PtrBool(request.GetSkipSignalReapply()),
	}

	resp, err := t.serviceClient.ResetWorkflowExecution(ctx, resetReq)
	if err != nil {
		return "", err
	}
	return resp.GetRunId(), nil
}

func (t *cadenceClient) GetBackendType() (backendType service.BackendType) {
	return service.BackendTypeCadence
}

func (t *cadenceClient) GetApiService() interface{} {
	return t.cClient
}
