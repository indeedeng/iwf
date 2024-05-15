package api

import (
	"context"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/errors"
)

type ApiService interface {
	ApiV1WorkflowStartPost(
		ctx context.Context, request iwfidl.WorkflowStartRequest,
	) (*iwfidl.WorkflowStartResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowWaitForStateCompletion(
		ctx context.Context, request iwfidl.WorkflowWaitForStateCompletionRequest,
	) (*iwfidl.WorkflowWaitForStateCompletionResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSignalPost(ctx context.Context, request iwfidl.WorkflowSignalRequest) *errors.ErrorAndStatus
	ApiV1WorkflowStopPost(ctx context.Context, request iwfidl.WorkflowStopRequest) *errors.ErrorAndStatus
	ApiV1WorkflowConfigUpdate(ctx context.Context, request iwfidl.WorkflowConfigUpdateRequest) *errors.ErrorAndStatus
	ApiV1WorkflowTriggerContinueAsNew(
		ctx context.Context, req iwfidl.TriggerContinueAsNewRequest,
	) (retError *errors.ErrorAndStatus)
	ApiV1WorkflowGetQueryAttributesPost(
		ctx context.Context, request iwfidl.WorkflowGetDataObjectsRequest,
	) (*iwfidl.WorkflowGetDataObjectsResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetSearchAttributesPost(
		ctx context.Context, request iwfidl.WorkflowGetSearchAttributesRequest,
	) (*iwfidl.WorkflowGetSearchAttributesResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetPost(
		ctx context.Context, request iwfidl.WorkflowGetRequest,
	) (*iwfidl.WorkflowGetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowGetWithWaitPost(
		ctx context.Context, request iwfidl.WorkflowGetRequest,
	) (*iwfidl.WorkflowGetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSearchPost(
		ctx context.Context, request iwfidl.WorkflowSearchRequest,
	) (*iwfidl.WorkflowSearchResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowRpcPost(
		ctx context.Context, request iwfidl.WorkflowRpcRequest,
	) (*iwfidl.WorkflowRpcResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowResetPost(
		ctx context.Context, request iwfidl.WorkflowResetRequest,
	) (*iwfidl.WorkflowResetResponse, *errors.ErrorAndStatus)
	ApiV1WorkflowSkipTimerPost(ctx context.Context, request iwfidl.WorkflowSkipTimerRequest) *errors.ErrorAndStatus
	ApiV1WorkflowDumpPost(
		ctx context.Context, request iwfidl.WorkflowDumpRequest,
	) (*iwfidl.WorkflowDumpResponse, *errors.ErrorAndStatus)
	ApiInfoHealth(ctx context.Context) *iwfidl.HealthInfo
	Close()
}
