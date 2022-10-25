package api

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"time"
)

type ApiService interface {
	ApiV1WorkflowStartPost(request iwfidl.WorkflowStartRequest) (*iwfidl.WorkflowStartResponse, *ErrorAndStatus)
	ApiV1WorkflowSignalPost(request iwfidl.WorkflowSignalRequest) *ErrorAndStatus
	ApiV1WorkflowGetQueryAttributesPost(request iwfidl.WorkflowGetQueryAttributesRequest) (*iwfidl.WorkflowGetQueryAttributesResponse, *ErrorAndStatus)
	ApiV1WorkflowGetPost(request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus)
	ApiV1WorkflowGetWithWaitPost(request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *ErrorAndStatus)
	ApiV1WorkflowSearchPost(request iwfidl.WorkflowSearchRequest) (*iwfidl.WorkflowSearchResponse, *ErrorAndStatus)
	ApiV1WorkflowResetPost(request iwfidl.WorkflowResetRequest) (*iwfidl.WorkflowResetResponse, *ErrorAndStatus)
	Close()
}

type ErrorAndStatus struct {
	StatusCode int
	Error      iwfidl.ErrorResponse
}

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
