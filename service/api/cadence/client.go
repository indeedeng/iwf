package cadence

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/api"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/retry"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
	"go.uber.org/cadence/encoded"
	"time"
)

type cadenceClient struct {
	domain        string
	cClient       client.Client
	closeFunc     func()
	serviceClient workflowserviceclient.Interface
	converter     encoded.DataConverter
}

func NewCadenceClient(domain string, cClient client.Client, serviceClient workflowserviceclient.Interface, converter encoded.DataConverter, closeFunc func()) api.UnifiedClient {
	return &cadenceClient{
		domain:        domain,
		cClient:       cClient,
		closeFunc:     closeFunc,
		serviceClient: serviceClient,
		converter:     converter,
	}
}

func (t *cadenceClient) Close() {
	t.closeFunc()
}

func (t *cadenceClient) StartInterpreterWorkflow(ctx context.Context, options api.StartWorkflowOptions, args ...interface{}) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                           options.ID,
		TaskList:                     options.TaskQueue,
		ExecutionStartToCloseTimeout: options.WorkflowExecutionTimeout,
		SearchAttributes:             options.SearchAttributes,
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
		workflowOptions.RetryPolicy = retry.ConvertCadenceRetryPolicy(options.RetryPolicy)
	}

	run, err := t.cClient.ExecuteWorkflow(ctx, workflowOptions, cadence.Interpreter, args...)
	if err != nil {
		return "", err
	}
	return run.GetRunID(), nil
}

func (t *cadenceClient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error {
	return t.cClient.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
}

func (t *cadenceClient) CancelWorkflow(ctx context.Context, workflowID string, runID string) error {
	return t.cClient.CancelWorkflow(ctx, workflowID, runID)
}

func (t *cadenceClient) ListWorkflow(ctx context.Context, request *api.ListWorkflowExecutionsRequest) (*api.ListWorkflowExecutionsResponse, error) {
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
	return &api.ListWorkflowExecutionsResponse{
		Executions:    executions,
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (t *cadenceClient) QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error {
	qres, err := t.cClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func (t *cadenceClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (*api.DescribeWorkflowExecutionResponse, error) {
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

	return &api.DescribeWorkflowExecutionResponse{
		RunId:            resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		Status:           status,
		SearchAttributes: searchAttributes,
	}, nil
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

func (t *cadenceClient) GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error {
	run := t.cClient.GetWorkflow(ctx, workflowID, runID)
	return run.Get(ctx, valuePtr)
}

func (t *cadenceClient) ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (newRunId string, err error) {

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
