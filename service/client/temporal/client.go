package temporal

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/common/utils"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/api/common/v1"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	realtemporal "go.temporal.io/sdk/temporal"
)

type temporalClient struct {
	tClient        client.Client
	namespace      string
	dataConverter  converter.DataConverter
	memoEncryption bool // this is a workaround for https://github.com/temporalio/sdk-go/issues/1045
}

func NewTemporalClient(
	tClient client.Client, namespace string, dataConverter converter.DataConverter, memoEncryption bool,
) uclient.UnifiedClient {
	return &temporalClient{
		tClient:        tClient,
		namespace:      namespace,
		dataConverter:  dataConverter,
		memoEncryption: memoEncryption,
	}
}

func (t *temporalClient) Close() {
	t.tClient.Close()
}

func (t *temporalClient) IsWorkflowAlreadyStartedError(err error) bool {
	if err.Error() == "schedule with this ID is already registered" {
		// there is no type to check, just a string
		// https://github.com/temporalio/sdk-go/blob/d10e87118a07b44fd09bf88d39a628f0e6e70c34/internal/error.go#L336
		return true
	}
	return realtemporal.IsWorkflowExecutionAlreadyStartedError(err)
}

func (t *temporalClient) GetRunIdFromWorkflowAlreadyStartedError(err error) (string, bool) {
	var workflowExecutionAlreadyStarted *serviceerror.WorkflowExecutionAlreadyStarted
	ok := errors.As(err, &workflowExecutionAlreadyStarted)

	runId := ""
	if ok {
		runId = workflowExecutionAlreadyStarted.RunId
	}

	return runId, ok
}

func (t *temporalClient) IsNotFoundError(err error) bool {
	var notFound *serviceerror.NotFound
	ok := errors.As(err, &notFound)
	return ok
}

func (t *temporalClient) IsRequestTimeoutError(err error) bool {
	var deadlineExceeded *serviceerror.DeadlineExceeded
	ok := errors.As(err, &deadlineExceeded)
	if ok {
		return ok
	}
	var canceled *serviceerror.Canceled
	ok = errors.As(err, &canceled)
	return ok
}

func (t *temporalClient) IsWorkflowTimeoutError(err error) bool {
	return realtemporal.IsTimeoutError(err)
}

func (t *temporalClient) GetApplicationErrorTypeIfIsApplicationError(err error) string {
	var applicationError *realtemporal.ApplicationError
	isAppErr := errors.As(err, &applicationError)
	if isAppErr {
		return applicationError.Type()
	}
	return ""
}

func (t *temporalClient) GetApplicationErrorDetails(err error, detailsPtr interface{}) error {
	var applicationError *realtemporal.ApplicationError
	isAppErr := errors.As(err, &applicationError)
	if !isAppErr {
		return fmt.Errorf("not an application error. Critical code bug")
	}
	if applicationError.HasDetails() {
		return applicationError.Details(detailsPtr)
	}
	return fmt.Errorf("application error doesn't have details. Critical code bug")
}

func (t *temporalClient) StartInterpreterWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions, args ...interface{},
) (runId string, err error) {
	memo, err := t.encryptMemoIfNeeded(options.Memo)
	if err != nil {
		return "", err
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                                       options.ID,
		TaskQueue:                                options.TaskQueue,
		WorkflowExecutionTimeout:                 options.WorkflowExecutionTimeout,
		SearchAttributes:                         options.SearchAttributes,
		Memo:                                     memo,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	if options.WorkflowIDReusePolicy != nil {
		workflowIdReusePolicy, err := mapToTemporalWorkflowIdReusePolicy(*options.WorkflowIDReusePolicy)
		if err != nil {
			return "", nil
		}

		workflowOptions.WorkflowIDReusePolicy = *workflowIdReusePolicy
	}

	if options.RetryPolicy != nil {
		workflowOptions.RetryPolicy = retry.ConvertTemporalWorkflowRetryPolicy(options.RetryPolicy)
	}

	if options.CronSchedule != nil && *options.CronSchedule != "" {
		// use temporal schedule instead of cron
		// https://temporal.io/blog/how-do-i-convert-my-cron-into-a-schedule
		// workflowOptions.CronSchedule = *options.CronSchedule

		_, err := t.tClient.ScheduleClient().Create(ctx, client.ScheduleOptions{
			ID: "schedule for workflow: " + options.ID,
			Spec: client.ScheduleSpec{
				CronExpressions: []string{*options.CronSchedule},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:                       workflowOptions.ID,
				TaskQueue:                workflowOptions.TaskQueue,
				Workflow:                 temporal.Interpreter,
				Args:                     args,
				WorkflowExecutionTimeout: workflowOptions.WorkflowExecutionTimeout,
				RetryPolicy:              workflowOptions.RetryPolicy,
				Memo:                     workflowOptions.Memo,
				TypedSearchAttributes:    workflowOptions.TypedSearchAttributes,
			},
		})

		return "", err
	}

	if options.WorkflowStartDelay != nil {
		workflowOptions.StartDelay = *options.WorkflowStartDelay
	}

	run, err := t.tClient.ExecuteWorkflow(ctx, workflowOptions, temporal.Interpreter, args...)
	if err != nil {
		return "", err
	}
	return run.GetRunID(), nil
}

func (t *temporalClient) StartWaitForStateCompletionWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions,
) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                       options.ID,
		WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY, // the workflow could be timeout, so we allow duplicate
		TaskQueue:                options.TaskQueue,
		WorkflowExecutionTimeout: options.WorkflowExecutionTimeout,
	}

	run, err := t.tClient.ExecuteWorkflow(ctx, workflowOptions, temporal.WaitforStateCompletionWorkflow)
	if err != nil {
		// because of WorkflowExecutionErrorWhenAlreadyStarted: false, we won't get WorkflowAlreadyStartedError as we do in Cadence
		return "", err
	}
	return run.GetRunID(), nil
}

func (t *temporalClient) SignalWithStartWaitForStateCompletionWorkflow(
	ctx context.Context, options uclient.StartWorkflowOptions, stateCompletionOutput iwfidl.StateCompletionOutput,
) error {
	workflowOptions := client.StartWorkflowOptions{
		ID:                       options.ID,
		WorkflowIDReusePolicy:    enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY, // the workflow could be timeout, so we allow duplicate
		TaskQueue:                options.TaskQueue,
		WorkflowExecutionTimeout: options.WorkflowExecutionTimeout,
	}

	_, err := t.tClient.SignalWithStartWorkflow(ctx, options.ID, service.StateCompletionSignalChannelName, stateCompletionOutput, workflowOptions, temporal.WaitforStateCompletionWorkflow)
	if err != nil {
		return err
	}
	return nil
}

func (t *temporalClient) SignalWorkflow(
	ctx context.Context, workflowID string, runID string, signalName string, arg interface{},
) error {
	return t.tClient.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

func (t *temporalClient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	return t.tClient.CancelWorkflow(ctx, workflowID, runID)
}

func (t *temporalClient) TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error {
	var reasonStr string
	if reason == "" {
		reasonStr = "Force termiantion from user"
	} else {
		reasonStr = reason
	}

	return t.tClient.TerminateWorkflow(ctx, workflowID, runID, reasonStr)
}

func (t *temporalClient) ListWorkflow(
	ctx context.Context, request *uclient.ListWorkflowExecutionsRequest,
) (*uclient.ListWorkflowExecutionsResponse, error) {
	listReq := &workflowservice.ListWorkflowExecutionsRequest{
		PageSize:      request.PageSize,
		Query:         request.Query,
		NextPageToken: request.NextPageToken,
	}
	resp, err := t.tClient.ListWorkflow(ctx, listReq)
	if err != nil {
		return nil, err
	}
	var executions []iwfidl.WorkflowSearchResponseEntry
	for _, exe := range resp.GetExecutions() {
		executions = append(executions, iwfidl.WorkflowSearchResponseEntry{
			WorkflowId:    exe.Execution.WorkflowId,
			WorkflowRunId: exe.Execution.RunId,
		})
	}
	return &uclient.ListWorkflowExecutionsResponse{
		Executions:    executions,
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (t *temporalClient) QueryWorkflow(
	ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{},
) error {
	qres, err := t.tClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func (t *temporalClient) DescribeWorkflowExecution(
	ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType,
) (*uclient.DescribeWorkflowExecutionResponse, error) {
	resp, err := t.tClient.DescribeWorkflowExecution(ctx, workflowID, runID)
	if err != nil {
		return nil, err
	}
	status, err := mapToIwfWorkflowStatus(resp.GetWorkflowExecutionInfo().GetStatus())
	if err != nil {
		return nil, err
	}

	searchAttributes, err := mapper.MapTemporalToIwfSearchAttributes(resp.GetWorkflowExecutionInfo().GetSearchAttributes(), requestedSearchAttributes)
	if err != nil {
		return nil, err
	}

	memo, err := t.getMemoAndDecryptIfNeeded(resp.GetWorkflowExecutionInfo().GetMemo())

	return &uclient.DescribeWorkflowExecutionResponse{
		RunId:                    resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		FirstRunId:               resp.GetWorkflowExecutionInfo().GetFirstRunId(),
		Status:                   status,
		SearchAttributes:         searchAttributes,
		Memos:                    memo,
		WorkflowStartedTimestamp: utils.ToNanoSeconds(resp.GetWorkflowExecutionInfo().GetStartTime()),
	}, err
}

func (t *temporalClient) encryptMemoIfNeeded(rawMemo map[string]interface{}) (map[string]interface{}, error) {
	if !t.memoEncryption || rawMemo == nil {
		return rawMemo, nil
	}

	out := map[string]interface{}{}
	for k, v := range rawMemo {

		pl, err := t.dataConverter.ToPayload(v)
		if err != nil {
			return nil, err
		}
		out[k] = pl
	}
	return out, nil
}

func (t *temporalClient) getMemoAndDecryptIfNeeded(memo *common.Memo) (map[string]iwfidl.EncodedObject, error) {
	if memo == nil || len(memo.GetFields()) == 0 {
		return nil, nil
	}

	out := map[string]iwfidl.EncodedObject{}
	for k, payload := range memo.GetFields() {

		if t.memoEncryption {
			var encryptedPayload commonpb.Payload
			err := converter.GetDefaultDataConverter().FromPayload(payload, &encryptedPayload)
			if err != nil {
				return nil, err
			}

			var value iwfidl.EncodedObject
			err = t.dataConverter.FromPayload(&encryptedPayload, &value)
			if err != nil {
				return nil, err
			}
			out[k] = value
		} else {
			var value iwfidl.EncodedObject
			err := converter.GetDefaultDataConverter().FromPayload(payload, &value)
			if err != nil {
				return nil, err
			}
			out[k] = value
		}
	}
	return out, nil
}

func mapToTemporalWorkflowIdReusePolicy(workflowIdReusePolicy iwfidl.WorkflowIDReusePolicy) (*enums.WorkflowIdReusePolicy, error) {
	var res enums.WorkflowIdReusePolicy
	switch workflowIdReusePolicy {
	case iwfidl.ALLOW_DUPLICATE:
		res = enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE
		return &res, nil
	case iwfidl.ALLOW_DUPLICATE_FAILED_ONLY:
		res = enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY
		return &res, nil
	case iwfidl.REJECT_DUPLICATE:
		res = enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE
		return &res, nil
	case iwfidl.TERMINATE_IF_RUNNING:
		res = enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING
		return &res, nil
	default:
		return nil, fmt.Errorf("unsupported workflow id reuse policy %s", workflowIdReusePolicy)
	}
}

func mapToIwfWorkflowStatus(status enums.WorkflowExecutionStatus) (iwfidl.WorkflowStatus, error) {
	switch status {
	case enums.WORKFLOW_EXECUTION_STATUS_CANCELED:
		return iwfidl.CANCELED, nil
	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		return iwfidl.COMPLETED, nil
	case enums.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW:
		return iwfidl.CONTINUED_AS_NEW, nil
	case enums.WORKFLOW_EXECUTION_STATUS_FAILED:
		return iwfidl.FAILED, nil
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		return iwfidl.RUNNING, nil
	case enums.WORKFLOW_EXECUTION_STATUS_TIMED_OUT:
		return iwfidl.TIMEOUT, nil
	case enums.WORKFLOW_EXECUTION_STATUS_TERMINATED:
		return iwfidl.TERMINATED, nil
	default:
		return "", fmt.Errorf("not supported status %s", status)
	}
}

func (t *temporalClient) GetWorkflowResult(
	ctx context.Context, valuePtr interface{}, workflowID string, runID string,
) error {
	run := t.tClient.GetWorkflow(ctx, workflowID, runID)
	return run.Get(ctx, valuePtr)
}

func (t *temporalClient) SynchronousUpdateWorkflow(
	ctx context.Context, valuePtr interface{}, workflowID, runID, updateType string, input interface{},
) error {
	args := []interface{}{input}
	options := client.UpdateWorkflowOptions{
		WorkflowID: workflowID,
		RunID:      runID,
		UpdateName: updateType,
		Args:       args,
		// TODO: Leaving this as Accepted that was a default value before WaitForStage became required argument, but Completed might be a better choice
		WaitForStage: client.WorkflowUpdateStageAccepted,
	}
	handle, err := t.tClient.UpdateWorkflow(ctx, options)
	if err != nil {
		return err
	}
	return handle.Get(context.Background(), valuePtr)
}

func (t *temporalClient) ResetWorkflow(
	ctx context.Context, request iwfidl.WorkflowResetRequest,
) (runId string, err error) {
	reqRunId := request.GetWorkflowRunId()
	if reqRunId == "" {
		// set default runId to current
		resp, err := t.tClient.DescribeWorkflowExecution(ctx, request.GetWorkflowId(), "")
		if err != nil {
			return "", err
		}
		reqRunId = resp.GetWorkflowExecutionInfo().GetExecution().GetRunId()
	}

	resetType := request.GetResetType()
	resetBaseRunID, resetEventId, err := getResetEventIDByType(ctx, resetType,
		t.namespace, request.GetWorkflowId(), reqRunId,
		t.tClient.WorkflowService(), t.dataConverter,
		request.GetHistoryEventId(), request.GetHistoryEventTime(), request.GetStateId(), request.GetStateExecutionId())

	if err != nil {
		return "", err
	}

	requestId := uuid.New().String()
	resetReapplyType := enums.RESET_REAPPLY_TYPE_SIGNAL
	if request.GetSkipSignalReapply() {
		resetReapplyType = enums.RESET_REAPPLY_TYPE_NONE
	}

	resp, err := t.tClient.ResetWorkflowExecution(ctx, &workflowservice.ResetWorkflowExecutionRequest{
		Namespace: t.namespace,
		WorkflowExecution: &common.WorkflowExecution{
			WorkflowId: request.WorkflowId,
			RunId:      resetBaseRunID,
		},
		Reason:                    request.GetReason(),
		WorkflowTaskFinishEventId: resetEventId,
		RequestId:                 requestId,
		ResetReapplyType:          resetReapplyType,
	})

	if err != nil {
		return "", err
	}
	return resp.GetRunId(), nil
}

func (t *temporalClient) GetBackendType() (backendType service.BackendType) {
	return service.BackendTypeTemporal
}

func (t *temporalClient) GetApiService() interface{} {
	return t.tClient.WorkflowService()
}
