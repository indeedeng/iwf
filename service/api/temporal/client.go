package temporal

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/api"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

type temporalClient struct {
	tClient client.Client
}

func (t *temporalClient) Close() {
	t.tClient.Close()
}

func (t *temporalClient) ExecuteWorkflow(ctx context.Context, options api.StartWorkflowOptions, workflow interface{}, args ...interface{}) (runId string, err error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:                 options.ID,
		TaskQueue:          options.TaskQueue,
		WorkflowRunTimeout: options.WorkflowRunTimeout,
	}

	run, err := t.tClient.ExecuteWorkflow(ctx, workflowOptions, workflow, args...)
	if err != nil {
		return "", err
	}
	return run.GetRunID(), nil
}

func (t *temporalClient) SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error {
	return t.tClient.SignalWorkflow(ctx, workflowID, runID, signalName, arg)
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
	return &api.DescribeWorkflowExecutionResponse{
		RunId:  resp.GetWorkflowExecutionInfo().GetExecution().GetRunId(),
		Status: status,
	}, nil
}

func mapToIwfWorkflowStatus(status enums.WorkflowExecutionStatus) (string, error) {
	switch status {
	case enums.WORKFLOW_EXECUTION_STATUS_CANCELED:
		return service.WorkflowStatusCanceled, nil
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

func NewTemporalClient(tClient client.Client) api.UnifiedClient {
	return &temporalClient{
		tClient: tClient,
	}
}
