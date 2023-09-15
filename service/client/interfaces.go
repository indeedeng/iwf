package client

import (
	"context"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
)

type UnifiedClient interface {
	Close()
	errorHandler
	StartInterpreterWorkflow(ctx context.Context, options StartWorkflowOptions, args ...interface{}) (runId string, err error)
	StartWaitForStateCompletionWorkflow(ctx context.Context, options StartWorkflowOptions) (runId string, err error)
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	SignalWithStartWaitForStateCompletionWorkflow(ctx context.Context, options StartWorkflowOptions, stateCompletionOutput iwfidl.StateCompletionOutput) error
	CancelWorkflow(ctx context.Context, workflowID string, runID string) error
	TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error
	ListWorkflow(ctx context.Context, request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error)
	QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error // TODO it doesn't return error correctly... the error is nil when query handler is not implemented
	DescribeWorkflowExecution(ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (*DescribeWorkflowExecutionResponse, error)
	GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error
	SynchronousUpdateWorkflow(ctx context.Context, valuePtr interface{}, workflowID, runID, updateType string, input interface{}) error
	ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error)
}

type errorHandler interface {
	GetApplicationErrorTypeIfIsApplicationError(err error) string
	GetApplicationErrorDetails(err error, detailsPtr interface{}) error
	IsWorkflowAlreadyStartedError(error) bool
	IsNotFoundError(error) bool
	IsRequestTimeoutError(error) bool
}

type StartWorkflowOptions struct {
	ID                       string
	TaskQueue                string
	WorkflowExecutionTimeout time.Duration
	WorkflowIDReusePolicy    *iwfidl.WorkflowIDReusePolicy
	CronSchedule             *string
	RetryPolicy              *iwfidl.WorkflowRetryPolicy
	SearchAttributes         map[string]interface{}
	Memo                     map[string]interface{}
}

type ListWorkflowExecutionsRequest struct {
	PageSize      int32
	Query         string
	NextPageToken []byte
}

type ListWorkflowExecutionsResponse struct {
	Executions    []iwfidl.WorkflowSearchResponseEntry
	NextPageToken []byte
}

type DescribeWorkflowExecutionResponse struct {
	Status                   iwfidl.WorkflowStatus
	RunId                    string
	SearchAttributes         map[string]iwfidl.SearchAttribute
	Memos                    map[string]iwfidl.EncodedObject
	WorkflowStartedTimestamp int64
}
