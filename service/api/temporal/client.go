package temporal

import (
	"context"
	"fmt"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

type temporalClient struct {
	tClient client.Client
}

func NewTemporalClient(tClient client.Client) api.UnifiedClient {
	return &temporalClient{
		tClient: tClient,
	}
}

func (t *temporalClient) Close() {
	t.tClient.Close()
}

func (t *temporalClient) StartInterpreterWorkflow(ctx context.Context, options api.StartWorkflowOptions, args ...interface{}) (runId string, err error) {
	workflowIdReusePolicy, err := mapToTemporalWorkflowIdReusePolicy(options.WorkflowIDReusePolicy)
	if err != nil {
		return "", nil
	}
	workflowOptions := client.StartWorkflowOptions{
		ID:                    options.ID,
		TaskQueue:             options.TaskQueue,
		WorkflowRunTimeout:    options.WorkflowRunTimeout,
		WorkflowIDReusePolicy: workflowIdReusePolicy,
		CronSchedule:          options.CronSchedule,
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

func (t *temporalClient) ListWorkflow(ctx context.Context, request *api.ListWorkflowExecutionsRequest) (*api.ListWorkflowExecutionsResponse, error) {
	listReq := &workflowservice.ListWorkflowExecutionsRequest{
		PageSize: request.PageSize,
		Query:    request.Query,
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
		Executions: executions,
	}, nil
}

func (t *temporalClient) QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error {
	qres, err := t.tClient.QueryWorkflow(ctx, workflowID, runID, queryType, args...)
	if err != nil {
		return err
	}
	return qres.Get(valuePtr)
}

func (t *temporalClient) DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*api.DescribeWorkflowExecutionResponse, error) {
	resp, err := t.tClient.DescribeWorkflowExecution(ctx, workflowID, runID)
	if err != nil {
		return nil, err
	}
	status, err := mapToIwfWorkflowStatus(resp.GetWorkflowExecutionInfo().GetStatus())
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

func mapToIwfSearchAttributes(searchAttributes *common.SearchAttributes) (map[string]iwfidl.SearchAttribute, error) {
	result := make(map[string]iwfidl.SearchAttribute)
	if searchAttributes == nil {
		return result, nil
	}

	for key, value := range searchAttributes.IndexedFields {
		var object interface{}
		err := converter.GetDefaultDataConverter().FromPayload(value, &object)
		if err != nil {
			return make(map[string]iwfidl.SearchAttribute), nil
		}

		str, isString := object.(string)
		if isString {
			result[key] = iwfidl.SearchAttribute{
				Key:         iwfidl.PtrString(key),
				StringValue: iwfidl.PtrString(str),
				ValueType:   iwfidl.PtrString(service.SearchAttributeValueTypeKeyword),
			}
		}
		number, isInt := object.(float64)
		if isInt {
			result[key] = iwfidl.SearchAttribute{
				Key:          iwfidl.PtrString(key),
				IntegerValue: iwfidl.PtrInt64(int64(number)),
				ValueType:    iwfidl.PtrString(service.SearchAttributeValueTypeInt),
			}
		}
	}

	return result, nil
}

func mapToTemporalWorkflowIdReusePolicy(workflowIdReusePolicy string) (enums.WorkflowIdReusePolicy, error) {
	switch workflowIdReusePolicy {
	case service.WorkflowIDReusePolicyAllowDuplicate:
		return enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE, nil
	case service.WorkflowIDReusePolicyAllowDuplicateFailedOnly:
		return enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE_FAILED_ONLY, nil
	case service.WorkflowIDReusePolicyRejectDuplicate:
		return enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE, nil
	case service.WorkflowIDReusePolicyTerminateIfRunning:
		return enums.WORKFLOW_ID_REUSE_POLICY_TERMINATE_IF_RUNNING, nil
	default:
		return enums.WORKFLOW_ID_REUSE_POLICY_UNSPECIFIED, fmt.Errorf("unsupported workflow id reuse policy %s", workflowIdReusePolicy)
	}
}

func mapToIwfWorkflowStatus(status enums.WorkflowExecutionStatus) (string, error) {
	switch status {
	case enums.WORKFLOW_EXECUTION_STATUS_CANCELED:
		return service.WorkflowStatusCanceled, nil
	case enums.WORKFLOW_EXECUTION_STATUS_COMPLETED:
		return service.WorkflowStatusCompleted, nil
	case enums.WORKFLOW_EXECUTION_STATUS_CONTINUED_AS_NEW:
		return service.WorkflowStatusContinueAsNew, nil
	case enums.WORKFLOW_EXECUTION_STATUS_FAILED:
		return service.WorkflowStatusFailed, nil
	case enums.WORKFLOW_EXECUTION_STATUS_RUNNING:
		return service.WorkflowStatusRunning, nil
	case enums.WORKFLOW_EXECUTION_STATUS_TIMED_OUT:
		return service.WorkflowStatusTimeout, nil
	case enums.WORKFLOW_EXECUTION_STATUS_TERMINATED:
		return service.WorkflowStatusTerminated, nil
	default:
		return "", fmt.Errorf("not supported status %s", status)
	}
}

func (t *temporalClient) GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error {
	run := t.tClient.GetWorkflow(ctx, workflowID, runID)
	return run.Get(ctx, valuePtr)
}

func (t *temporalClient) ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error) {
	return "", fmt.Errorf("not supported")
}
