package temporal

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/api"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	realtemporal "go.temporal.io/sdk/temporal"
)

type temporalClient struct {
	tClient       client.Client
	namespace     string
	dataConverter converter.DataConverter
}

func (t *temporalClient) IsWorkflowAlreadyStartedError(err error) bool {
	return realtemporal.IsWorkflowExecutionAlreadyStartedError(err)
}

func (t *temporalClient) IsNotFoundError(err error) bool {
	_, ok := err.(*serviceerror.NotFound)
	return ok
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

func NewTemporalClient(tClient client.Client, namespace string, dataConverter converter.DataConverter) api.UnifiedClient {
	return &temporalClient{
		tClient:       tClient,
		namespace:     namespace,
		dataConverter: dataConverter,
	}
}

func (t *temporalClient) Close() {
	t.tClient.Close()
}

func (t *temporalClient) StartInterpreterWorkflow(ctx context.Context, options api.StartWorkflowOptions, args ...interface{}) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                                       options.ID,
		TaskQueue:                                options.TaskQueue,
		WorkflowExecutionTimeout:                 options.WorkflowExecutionTimeout,
		SearchAttributes:                         options.SearchAttributes,
		WorkflowExecutionErrorWhenAlreadyStarted: true,
	}

	if options.WorkflowIDReusePolicy != nil {
		workflowIdReusePolicy, err := mapToTemporalWorkflowIdReusePolicy(*options.WorkflowIDReusePolicy)
		if err != nil {
			return "", nil
		}

		workflowOptions.WorkflowIDReusePolicy = *workflowIdReusePolicy
	}

	if options.CronSchedule != nil {
		workflowOptions.CronSchedule = *options.CronSchedule
	}

	if options.RetryPolicy != nil {
		workflowOptions.RetryPolicy = retry.ConvertTemporalWorkflowRetryPolicy(options.RetryPolicy)
	}

	run, err := t.tClient.ExecuteWorkflow(ctx, workflowOptions, temporal.Interpreter, args...)
	if err != nil {
		return "", err
	}
	return run.GetRunID(), nil
}

func (t *temporalClient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error {
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

func (t *temporalClient) ListWorkflow(ctx context.Context, request *api.ListWorkflowExecutionsRequest) (*api.ListWorkflowExecutionsResponse, error) {
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
	return &api.ListWorkflowExecutionsResponse{
		Executions:    executions,
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (t *temporalClient) QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error {
	qres, err := t.tClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func (t *temporalClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (*api.DescribeWorkflowExecutionResponse, error) {
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

	return &api.DescribeWorkflowExecutionResponse{
		RunId:            resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		Status:           status,
		SearchAttributes: searchAttributes,
	}, nil
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

func (t *temporalClient) GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error {
	run := t.tClient.GetWorkflow(ctx, workflowID, runID)
	return run.Get(ctx, valuePtr)
}

func (t *temporalClient) ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error) {
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
