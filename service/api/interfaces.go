package api

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"time"
)

type UnifiedClient interface {
	Close()
	StartInterpreterWorkflow(ctx context.Context, options StartWorkflowOptions, args ...interface{}) (runId string, err error)
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	ListWorkflow(ctx context.Context, request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error)
	QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error
	DescribeWorkflowExecution(ctx context.Context, workflowID, runID string) (*DescribeWorkflowExecutionResponse, error)
	GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error
	ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error)
}

type StartWorkflowOptions struct {
	ID                 string
	TaskQueue          string
	WorkflowRunTimeout time.Duration
}

type ListWorkflowExecutionsRequest struct {
	PageSize int32
	Query    string
}

type ListWorkflowExecutionsResponse struct {
	Executions []iwfidl.WorkflowSearchResponseEntry
}

type DescribeWorkflowExecutionResponse struct {
	Status string
	RunId  string
}
