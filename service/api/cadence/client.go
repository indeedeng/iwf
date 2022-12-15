package cadence

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/cadence/workflow"
	"time"

	"github.com/google/uuid"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/.gen/go/shared"
	"go.uber.org/cadence/client"
)

type cadenceClient struct {
	domain        string
	cClient       client.Client
	closeFunc     func()
	serviceClient workflowserviceclient.Interface
}

func NewCadenceClient(domain string, cClient client.Client, serviceClient workflowserviceclient.Interface, closeFunc func()) api.UnifiedClient {
	return &cadenceClient{
		domain:        domain,
		cClient:       cClient,
		closeFunc:     closeFunc,
		serviceClient: serviceClient,
	}
}

func (t *cadenceClient) Close() {
	t.closeFunc()
}

func (t *cadenceClient) StartInterpreterWorkflow(ctx context.Context, options api.StartWorkflowOptions, args ...interface{}) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                           options.ID,
		TaskList:                     options.TaskQueue,
		ExecutionStartToCloseTimeout: options.WorkflowRunTimeout,
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
		workflowOptions.RetryPolicy = &workflow.RetryPolicy{
			InitialInterval:    time.Second * time.Duration(options.RetryPolicy.GetInitialIntervalSeconds()),
			MaximumInterval:    time.Second * time.Duration(options.RetryPolicy.GetMaximumIntervalSeconds()),
			MaximumAttempts:    options.RetryPolicy.GetMaximumAttempts(),
			BackoffCoefficient: float64(options.RetryPolicy.GetBackoffCoefficient()),
		}
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
		PageSize: &request.PageSize,
		Query:    &request.Query,
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
		Executions: executions,
	}, nil
}

func (t *cadenceClient) QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error {
	qres, err := t.cClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func (t *cadenceClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*api.DescribeWorkflowExecutionResponse, error) {
	resp, err := t.cClient.DescribeWorkflowExecution(ctx, workflowID, runID)
	if err != nil {
		return nil, err
	}
	status, err := mapToIwfWorkflowStatus(resp.GetWorkflowExecutionInfo().CloseStatus)
	if err != nil {
		return nil, err
	}
	searchAttributes, err := mapToIwfSearchAttributes(resp.GetWorkflowExecutionInfo().GetSearchAttributes())
	if err != nil {
		return nil, err
	}

	return &api.DescribeWorkflowExecutionResponse{
		RunId:            resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		Status:           status,
		SearchAttributes: searchAttributes,
	}, nil
}

func mapToIwfSearchAttributes(searchAttributes *shared.SearchAttributes) (map[string]iwfidl.SearchAttribute, error) {
	result := make(map[string]iwfidl.SearchAttribute)
	if searchAttributes == nil {
		return result, nil
	}

	for key, value := range searchAttributes.IndexedFields {
		var object interface{}
		err := client.NewValue(value).Get(&object)
		if err != nil {
			return make(map[string]iwfidl.SearchAttribute), nil
		}

		str, ok := object.(string)
		if ok {
			result[key] = iwfidl.SearchAttribute{
				Key:         iwfidl.PtrString(key),
				StringValue: iwfidl.PtrString(str),
				ValueType:   iwfidl.PtrString(service.SearchAttributeValueTypeKeyword),
			}
		} else {
			number, ok := object.(json.Number)
			if ok {
				integer, _ := number.Int64()
				result[key] = iwfidl.SearchAttribute{
					Key:          iwfidl.PtrString(key),
					IntegerValue: iwfidl.PtrInt64(integer),
					ValueType:    iwfidl.PtrString(service.SearchAttributeValueTypeInt),
				}
			}
		}
	}

	return result, nil
}

func mapToCadenceWorkflowIdReusePolicy(workflowIdReusePolicy string) (*client.WorkflowIDReusePolicy, error) {
	var res client.WorkflowIDReusePolicy
	switch workflowIdReusePolicy {
	case service.WorkflowIDReusePolicyAllowDuplicate:
		res = client.WorkflowIDReusePolicyAllowDuplicate
		return &res, nil
	case service.WorkflowIDReusePolicyAllowDuplicateFailedOnly:
		res = client.WorkflowIDReusePolicyAllowDuplicateFailedOnly
		return &res, nil
	case service.WorkflowIDReusePolicyRejectDuplicate:
		res = client.WorkflowIDReusePolicyRejectDuplicate
		return &res, nil
	case service.WorkflowIDReusePolicyTerminateIfRunning:
		res = client.WorkflowIDReusePolicyTerminateIfRunning
		return &res, nil
	default:
		return nil, fmt.Errorf("unsupported workflow id reuse policy %s", workflowIdReusePolicy)
	}
}

func mapToIwfWorkflowStatus(status *shared.WorkflowExecutionCloseStatus) (string, error) {
	if status == nil {
		return service.WorkflowStatusRunning, nil
	}

	switch *status {
	case shared.WorkflowExecutionCloseStatusCanceled:
		return service.WorkflowStatusCanceled, nil
	case shared.WorkflowExecutionCloseStatusContinuedAsNew:
		return service.WorkflowStatusContinueAsNew, nil
	case shared.WorkflowExecutionCloseStatusFailed:
		return service.WorkflowStatusFailed, nil
	case shared.WorkflowExecutionCloseStatusTimedOut:
		return service.WorkflowStatusTimeout, nil
	case shared.WorkflowExecutionCloseStatusTerminated:
		return service.WorkflowStatusTerminated, nil
	case shared.WorkflowExecutionCloseStatusCompleted:
		return service.WorkflowStatusCompleted, nil
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
		resp, err := t.DescribeWorkflowExecution(ctx, request.GetWorkflowId(), "")
		if err != nil {
			return "", err
		}
		reqRunId = resp.RunId
	}

	resetType := service.ResetType(request.GetResetType())
	resetBaseRunID, decisionFinishID, err := getResetIDsByType(ctx, resetType, t.domain, request.GetWorkflowId(),
		reqRunId, t.serviceClient, request.GetHistoryEventId(), request.GetHistoryEventTime())

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
