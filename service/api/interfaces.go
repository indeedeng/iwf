package api

import (
	"context"
	"time"

	"github.com/indeedeng/iwf/service/common/errors"

	"github.com/indeedeng/iwf/gen/iwfidl"
)

type ApiService interface {
	ApiV1WorkflowStartPost(ctx context.Context, request iwfidl.WorkflowStartRequest) (*iwfidl.WorkflowStartResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSignalPost(ctx context.Context, request iwfidl.WorkflowSignalRequest) *errors.ErrorAndStatus
	ApiV1WorkflowStopPost(ctx context.Context, request iwfidl.WorkflowStopRequest) *errors.ErrorAndStatus
	ApiV1WorkflowConfigUpdate(ctx context.Context, request iwfidl.WorkflowConfigUpdateRequest) *errors.ErrorAndStatus
	ApiV1WorkflowGetQueryAttributesPost(ctx context.Context, request iwfidl.WorkflowGetDataObjectsRequest) (*iwfidl.WorkflowGetDataObjectsResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetSearchAttributesPost(ctx context.Context, request iwfidl.WorkflowGetSearchAttributesRequest) (*iwfidl.WorkflowGetSearchAttributesResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetPost(ctx context.Context, request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetWithWaitPost(ctx context.Context, request iwfidl.WorkflowGetRequest) (*iwfidl.WorkflowGetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSearchPost(ctx context.Context, request iwfidl.WorkflowSearchRequest) (*iwfidl.WorkflowSearchResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowRpcPost(ctx context.Context, request iwfidl.WorkflowRpcRequest) (*iwfidl.WorkflowRpcResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowResetPost(ctx context.Context, request iwfidl.WorkflowResetRequest) (*iwfidl.WorkflowResetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSkipTimerPost(ctx context.Context, request iwfidl.WorkflowSkipTimerRequest) *errors.ErrorAndStatus
	ApiV1WorkflowDumpPost(ctx context.Context, request iwfidl.WorkflowDumpRequest) (*iwfidl.WorkflowDumpResponse, *errors.ErrorAndStatus)
	Close()
}

type UnifiedClient interface {
	Close()
	errorHandler
	StartInterpreterWorkflow(ctx context.Context, options StartWorkflowOptions, args ...interface{}) (runId string, err error)
	SignalWorkflow(ctx context.Context, workflowID string, runID string, signalName string, arg interface{}) error
	CancelWorkflow(ctx context.Context, workflowID string, runID string) error
	TerminateWorkflow(ctx context.Context, workflowID string, runID string, reason string) error
	ListWorkflow(ctx context.Context, request *ListWorkflowExecutionsRequest) (*ListWorkflowExecutionsResponse, error)
	QueryWorkflow(ctx context.Context, valuePtr interface{}, workflowID string, runID string, queryType string, args ...interface{}) error // TODO it doesn't return error correctly... the error is nil when query handler is not implemented
	DescribeWorkflowExecution(ctx context.Context, workflowID, runID string, requestedSearchAttributes []iwfidl.SearchAttributeKeyAndType) (*DescribeWorkflowExecutionResponse, error)
	GetWorkflowResult(ctx context.Context, valuePtr interface{}, workflowID string, runID string) error
	ResetWorkflow(ctx context.Context, request iwfidl.WorkflowResetRequest) (runId string, err error)
}

type errorHandler interface {
	GetApplicationErrorTypeIfIsApplicationError(err error) string
	GetApplicationErrorDetails(err error, detailsPtr interface{}) error
	IsWorkflowAlreadyStartedError(error) bool
	IsNotFoundError(error) bool
}

type StartWorkflowOptions struct {
	ID                       string
	TaskQueue                string
	WorkflowExecutionTimeout time.Duration
	WorkflowIDReusePolicy    *iwfidl.WorkflowIDReusePolicy
	CronSchedule             *string
	RetryPolicy              *iwfidl.WorkflowRetryPolicy
	SearchAttributes         map[string]interface{}
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
	Status           iwfidl.WorkflowStatus
	RunId            string
	SearchAttributes map[string]iwfidl.SearchAttribute
}
