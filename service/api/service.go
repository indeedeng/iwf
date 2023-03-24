package api

import (
	"context"
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/log/tag"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/ptr"
	"net/http"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type serviceImpl struct {
	client    UnifiedClient
	taskQueue string
	logger    log.Logger
	config    *config.Config
}

func (s *serviceImpl) Close() {
	s.client.Close()
}

func NewApiService(config *config.Config, client UnifiedClient, taskQueue string, logger log.Logger) (ApiService, error) {
	return &serviceImpl{
		client:    client,
		taskQueue: taskQueue,
		logger:    logger,
		config:    config,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowStartPost(ctx context.Context, req iwfidl.WorkflowStartRequest) (wresp *iwfidl.WorkflowStartResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	workflowOptions := StartWorkflowOptions{
		ID:                       req.GetWorkflowId(),
		TaskQueue:                s.taskQueue,
		WorkflowExecutionTimeout: time.Duration(req.WorkflowTimeoutSeconds) * time.Second,
	}

	var initSAs []iwfidl.SearchAttribute
	if req.WorkflowStartOptions != nil {
		workflowOptions.WorkflowIDReusePolicy = req.WorkflowStartOptions.WorkflowIDReusePolicy
		workflowOptions.CronSchedule = req.WorkflowStartOptions.CronSchedule
		workflowOptions.RetryPolicy = req.WorkflowStartOptions.RetryPolicy
		var err error
		workflowOptions.SearchAttributes, err = mapper.MapToInternalSearchAttributes(req.WorkflowStartOptions.SearchAttributes)
		if err != nil {
			return nil, s.handleError(err)
		}
		workflowOptions.SearchAttributes[service.SearchAttributeIwfWorkflowType] = req.IwfWorkflowType
		initSAs = req.WorkflowStartOptions.SearchAttributes
	}

	disableSystemSearchAttributes := false
	if s.config.Backend.Temporal == nil && s.config.Backend.Cadence != nil {
		disableSystemSearchAttributes = s.config.Backend.Cadence.DisableSystemSearchAttributes
	}

	input := service.InterpreterWorkflowInput{
		IwfWorkflowType:      req.GetIwfWorkflowType(),
		IwfWorkerUrl:         req.GetIwfWorkerUrl(),
		StartStateId:         req.GetStartStateId(),
		StateInput:           req.GetStateInput(),
		StateOptions:         req.GetStateOptions(),
		InitSearchAttributes: initSAs,
		Config: service.WorkflowConfig{
			DisableSystemSearchAttributes: disableSystemSearchAttributes,
		},
	}
	runId, err := s.client.StartInterpreterWorkflow(ctx, workflowOptions, input)
	if err != nil {
		return nil, s.handleError(err)
	}

	s.logger.Info("Started workflow", tag.WorkflowID(req.WorkflowId), tag.WorkflowRunID(runId))

	return &iwfidl.WorkflowStartResponse{
		WorkflowRunId: iwfidl.PtrString(runId),
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSignalPost(ctx context.Context, req iwfidl.WorkflowSignalRequest) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.SignalWorkflow(ctx,
		req.GetWorkflowId(), req.GetWorkflowRunId(), req.GetSignalChannelName(), req.GetSignalValue())
	if err != nil {
		return s.handleError(err)
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowStopPost(ctx context.Context, req iwfidl.WorkflowStopRequest) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.CancelWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId())
	if err != nil {
		return s.handleError(err)
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetQueryAttributesPost(ctx context.Context, req iwfidl.WorkflowGetDataObjectsRequest) (wresp *iwfidl.WorkflowGetDataObjectsResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	var queryResult1 service.GetDataObjectsQueryResponse
	err := s.client.QueryWorkflow(ctx, &queryResult1,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.GetDataObjectsWorkflowQueryType,
		service.GetDataObjectsQueryRequest{
			Keys: req.Keys,
		})

	if err != nil {
		return nil, s.handleError(err)
	}

	return &iwfidl.WorkflowGetDataObjectsResponse{
		Objects: queryResult1.DataObjects,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowGetSearchAttributesPost(ctx context.Context, req iwfidl.WorkflowGetSearchAttributesRequest) (wresp *iwfidl.WorkflowGetSearchAttributesResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	response, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), req.Keys)
	if err != nil {
		return nil, s.handleError(err)
	}

	var searchAttributes []iwfidl.SearchAttribute
	for _, v := range req.Keys {
		searchAttribute, exist := response.SearchAttributes[*v.Key]
		if exist {
			searchAttributes = append(searchAttributes, searchAttribute)
		}
	}

	return &iwfidl.WorkflowGetSearchAttributesResponse{
		SearchAttributes: searchAttributes,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowGetPost(ctx context.Context, req iwfidl.WorkflowGetRequest) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	return s.doApiV1WorkflowGetPost(ctx, req, false)
}

func (s *serviceImpl) ApiV1WorkflowGetWithWaitPost(ctx context.Context, req iwfidl.WorkflowGetRequest) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	return s.doApiV1WorkflowGetPost(ctx, req, true)
}

// withWait:
//
//	 because s.client.GetWorkflowResult will wait for the completion if workflow is running --
//		when withWait is false, if workflow is not running and needResults is true, it will then call s.client.GetWorkflowResult to get results
//		when withWait is true, it will do everything
func (s *serviceImpl) doApiV1WorkflowGetPost(ctx context.Context, req iwfidl.WorkflowGetRequest, withWait bool) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	descResp, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), nil)
	if err != nil {
		return nil, s.handleError(err)
	}

	status := descResp.Status
	var output service.InterpreterWorkflowOutput
	var getErr error
	if !withWait {
		if descResp.Status != iwfidl.RUNNING && req.GetNeedsResults() {
			getErr = s.client.GetWorkflowResult(ctx, &output, req.GetWorkflowId(), req.GetWorkflowRunId())
			if getErr == nil {
				status = iwfidl.COMPLETED
			}
		}
	} else {
		getErr = s.client.GetWorkflowResult(ctx, &output, req.GetWorkflowId(), req.GetWorkflowRunId())
		if getErr == nil {
			status = iwfidl.COMPLETED
		}
	}

	if getErr != nil { // workflow closed at an abnormal state(failed/timeout/terminated/canceled)
		var outputsToReturnWf []iwfidl.StateCompletionOutput
		var errMsg string
		errType := s.client.GetApplicationErrorTypeIfIsApplicationError(getErr)
		if errType != "" {
			errTypeEnum := iwfidl.WorkflowErrorType(errType)
			if errTypeEnum == iwfidl.STATE_DECISION_FAILING_WORKFLOW_ERROR_TYPE {
				err = s.client.GetApplicationErrorDetails(getErr, &outputsToReturnWf)
				if err != nil {
					return nil, s.handleError(err)
				}
			} else {
				err = s.client.GetApplicationErrorDetails(getErr, &errMsg)
				if err != nil {
					return nil, s.handleError(err)
				}
			}

			var errMsgPtr *string
			if errMsg != "" {
				errMsgPtr = iwfidl.PtrString(errMsg)
			}
			return &iwfidl.WorkflowGetResponse{
				WorkflowRunId:  descResp.RunId,
				WorkflowStatus: iwfidl.FAILED,
				ErrorType:      ptr.Any(errTypeEnum),
				ErrorMessage:   errMsgPtr,
				Results:        outputsToReturnWf,
			}, nil
		} else {
			// it could be timeout/terminated/canceled/etc. We need to describe again to get the final status
			descResp, err = s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), nil)
			if err != nil {
				return nil, s.handleError(err)
			}
			errMsg = ""
			if descResp.Status == iwfidl.FAILED {
				errMsg = "unknown workflow failure from interpreter implementation"
				s.logger.Error(errMsg, tag.WorkflowID(req.GetWorkflowId()), tag.WorkflowRunID(descResp.RunId))
			}
			var errMsgPtr *string
			if errMsg != "" {
				errMsgPtr = iwfidl.PtrString(errMsg)
			}
			return &iwfidl.WorkflowGetResponse{
				WorkflowRunId:  descResp.RunId,
				WorkflowStatus: descResp.Status,
				ErrorMessage:   errMsgPtr,
			}, nil
		}
	}

	return &iwfidl.WorkflowGetResponse{
		WorkflowRunId:  descResp.RunId,
		WorkflowStatus: status,
		Results:        output.StateCompletionOutputs,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSearchPost(ctx context.Context, req iwfidl.WorkflowSearchRequest) (wresp *iwfidl.WorkflowSearchResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	pageSize := int32(1000)
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	resp, err := s.client.ListWorkflow(ctx, &ListWorkflowExecutionsRequest{
		PageSize:      pageSize,
		Query:         req.GetQuery(),
		NextPageToken: []byte(req.GetNextPageToken()),
	})
	if err != nil {
		return nil, s.handleError(err)
	}
	return &iwfidl.WorkflowSearchResponse{
		WorkflowExecutions: resp.Executions,
		NextPageToken:      ptr.Any(string(resp.NextPageToken)),
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowResetPost(ctx context.Context, req iwfidl.WorkflowResetRequest) (wresp *iwfidl.WorkflowResetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	runId, err := s.client.ResetWorkflow(ctx, req)
	if err != nil {
		return nil, s.handleError(err)
	}
	return &iwfidl.WorkflowResetResponse{
		WorkflowRunId: runId,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSkipTimerPost(ctx context.Context, request iwfidl.WorkflowSkipTimerRequest) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	if request.GetTimerCommandId() == "" && request.TimerCommandIndex == nil {
		return makeInvalidRequestError("must provide either a timerCommandId or index")
	}

	timerInfos := service.GetCurrentTimerInfosQueryResponse{}
	err := s.client.QueryWorkflow(ctx, &timerInfos, request.GetWorkflowId(), request.GetWorkflowRunId(), service.GetCurrentTimerInfosQueryType)
	if err != nil {
		return s.handleError(err)
	}
	_, valid := service.ValidateTimerSkipRequest(timerInfos.StateExecutionCurrentTimerInfos, request.GetWorkflowStateExecutionId(), request.GetTimerCommandId(), int(request.GetTimerCommandIndex()))
	if !valid {
		return makeInvalidRequestError("requested timer command doesn't exist")
	}
	signal := service.SkipTimerSignalRequest{
		StateExecutionId: request.GetWorkflowStateExecutionId(),
		CommandId:        request.GetTimerCommandId(),
		CommandIndex:     int(request.GetTimerCommandIndex()),
	}
	err = s.client.SignalWorkflow(ctx, request.GetWorkflowId(), request.GetWorkflowRunId(), service.SkipTimerSignalChannelName, signal)
	if err != nil {
		return s.handleError(err)
	}
	return nil
}

func makeInvalidRequestError(msg string) *errors.ErrorAndStatus {
	return errors.NewErrorAndStatus(http.StatusBadRequest,
		iwfidl.UNCATEGORIZED_SUB_STATUS,
		"invalid request - "+msg)
}

func (s *serviceImpl) handleError(err error) *errors.ErrorAndStatus {
	s.logger.Error("encounter error for API", tag.Error(err))
	status := http.StatusInternalServerError
	subStatus := iwfidl.UNCATEGORIZED_SUB_STATUS
	if s.client.IsNotFoundError(err) {
		status = http.StatusBadRequest
		subStatus = iwfidl.WORKFLOW_NOT_EXISTS_SUB_STATUS
	}
	if s.client.IsWorkflowAlreadyStartedError(err) {
		status = http.StatusBadRequest
		subStatus = iwfidl.WORKFLOW_ALREADY_STARTED_SUB_STATUS
	}
	return errors.NewErrorAndStatus(
		status,
		subStatus,
		err.Error(),
	)
}
