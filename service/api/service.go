package api

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/service/common/event"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"github.com/indeedeng/iwf/service/interpreter/interfaces"
	"net/http"
	"os"
	"strings"
	"time"

	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/log/tag"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/indeedeng/iwf/service/common/rpc"
	"github.com/indeedeng/iwf/service/common/utils"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type serviceImpl struct {
	client    uclient.UnifiedClient
	taskQueue string
	logger    log.Logger
	config    config.Config
}

func (s *serviceImpl) Close() {
	s.client.Close()
}

func NewApiService(
	config config.Config, client uclient.UnifiedClient, taskQueue string, logger log.Logger,
) (ApiService, error) {
	return &serviceImpl{
		client:    client,
		taskQueue: taskQueue,
		logger:    logger,
		config:    config,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowStartPost(
	ctx context.Context, req iwfidl.WorkflowStartRequest,
) (wresp *iwfidl.WorkflowStartResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	sysSAs := map[string]interface{}{
		service.SearchAttributeIwfWorkflowType: req.IwfWorkflowType,
	}

	workflowOptions := uclient.StartWorkflowOptions{
		ID:                       req.GetWorkflowId(),
		TaskQueue:                s.taskQueue,
		WorkflowExecutionTimeout: time.Duration(req.WorkflowTimeoutSeconds) * time.Second,
		SearchAttributes:         sysSAs,
	}

	var workflowConfig iwfidl.WorkflowConfig
	if s.config.Interpreter.DefaultWorkflowConfig == nil {
		workflowConfig = *config.DefaultWorkflowConfig
	} else {
		workflowConfig = *s.config.Interpreter.DefaultWorkflowConfig
	}
	var initCustomSAs []iwfidl.SearchAttribute
	var initCustomDAs []iwfidl.KeyValue
	// workerUrl is always needed, for optimizing None as persistence loading type
	workflowOptions.Memo = map[string]interface{}{
		service.WorkerUrlMemoKey: iwfidl.EncodedObject{
			Data: iwfidl.PtrString(req.IwfWorkerUrl),
		},
	}

	ignoreAlreadyStartedError := false
	var requestId *string

	useMemoForDAs := false
	if req.WorkflowStartOptions != nil {
		startOptions := req.WorkflowStartOptions
		workflowOptions.WorkflowIDReusePolicy = compatibility.GetWorkflowIdReusePolicy(*startOptions)
		workflowOptions.CronSchedule = startOptions.CronSchedule
		workflowOptions.RetryPolicy = startOptions.RetryPolicy
		var err error
		initialCustomSAInternal, err := mapper.MapToInternalSearchAttributes(startOptions.SearchAttributes)
		if err != nil {
			return nil, s.handleError(err, WorkflowStartApiPath, req.GetWorkflowId())
		}
		workflowOptions.SearchAttributes = utils.MergeMap(initialCustomSAInternal, workflowOptions.SearchAttributes)

		initCustomSAs = startOptions.SearchAttributes
		initCustomDAs = startOptions.DataAttributes
		if startOptions.HasWorkflowConfigOverride() {
			configOverride := startOptions.GetWorkflowConfigOverride()
			overrideWorkflowConfig(configOverride, &workflowConfig)
		}

		workflowAlreadyStartedOptions := startOptions.WorkflowAlreadyStartedOptions

		if workflowAlreadyStartedOptions != nil {
			ignoreAlreadyStartedError = req.WorkflowStartOptions.WorkflowAlreadyStartedOptions.IgnoreAlreadyStartedError
			if workflowAlreadyStartedOptions.RequestId != nil {
				requestId = workflowAlreadyStartedOptions.RequestId
			}
		}

		if startOptions.GetUseMemoForDataAttributes() {
			useMemoForDAs = true
			workflowOptions.Memo[service.UseMemoForDataAttributesKey] = iwfidl.EncodedObject{
				// Note: the value is actually not too important, we will check the presence of the key only as today
				Data: iwfidl.PtrString("true"),
			}
			for _, da := range initCustomDAs {
				workflowOptions.Memo[da.GetKey()] = da.GetValue()
			}
		}
		if requestId != nil {
			workflowOptions.Memo[service.WorkflowRequestId] = iwfidl.EncodedObject{
				Data: requestId,
			}
		}
		if startOptions.WorkflowStartDelaySeconds != nil {
			workflowOptions.WorkflowStartDelay =
				ptr.Any(time.Duration(*startOptions.WorkflowStartDelaySeconds) * time.Second)
		}
	}

	input := service.InterpreterWorkflowInput{
		IwfWorkflowType:                    req.GetIwfWorkflowType(),
		IwfWorkerUrl:                       req.GetIwfWorkerUrl(),
		StartStateId:                       req.StartStateId,
		StateInput:                         req.StateInput,
		StateOptions:                       req.StateOptions,
		InitSearchAttributes:               initCustomSAs,
		InitDataAttributes:                 initCustomDAs,
		Config:                             workflowConfig,
		UseMemoForDataAttributes:           useMemoForDAs,
		WaitForCompletionStateExecutionIds: req.GetWaitForCompletionStateExecutionIds(),
		WaitForCompletionStateIds:          req.GetWaitForCompletionStateIds(),
	}

	runId, err := s.client.StartInterpreterWorkflow(ctx, workflowOptions, input)
	if err != nil {
		shouldReturnError := true

		if s.client.IsWorkflowAlreadyStartedError(err) && ignoreAlreadyStartedError {
			alreadyRunningRunId, _ := s.client.GetRunIdFromWorkflowAlreadyStartedError(err)
			runId = alreadyRunningRunId

			if requestId == nil {
				shouldReturnError = false
			} else {
				response, descErr := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), runId, nil)
				if descErr != nil {
					return nil, s.handleError(err, WorkflowStartApiPath, req.WorkflowId)
				}

				if *response.Memos[service.WorkflowRequestId].Data == *requestId {
					shouldReturnError = false
				}
			}
		}

		if shouldReturnError {
			return nil, s.handleError(err, WorkflowStartApiPath, req.GetWorkflowId())
		}
	} else {
		s.logger.Info("Started workflow", tag.WorkflowID(req.WorkflowId), tag.WorkflowRunID(runId))
	}

	return &iwfidl.WorkflowStartResponse{
		WorkflowRunId: iwfidl.PtrString(runId),
	}, nil
}

func overrideWorkflowConfig(configOverride iwfidl.WorkflowConfig, workflowConfig *iwfidl.WorkflowConfig) {
	if configOverride.ExecutingStateIdMode != nil {
		workflowConfig.ExecutingStateIdMode = configOverride.ExecutingStateIdMode
	}
	if configOverride.ContinueAsNewThreshold != nil {
		workflowConfig.ContinueAsNewThreshold = configOverride.ContinueAsNewThreshold
	}
	if configOverride.ContinueAsNewPageSizeInBytes != nil {
		workflowConfig.ContinueAsNewPageSizeInBytes = configOverride.ContinueAsNewPageSizeInBytes
	}
	if configOverride.DisableSystemSearchAttribute != nil {
		workflowConfig.DisableSystemSearchAttribute = configOverride.DisableSystemSearchAttribute
	}
	if configOverride.OptimizeActivity != nil {
		workflowConfig.OptimizeActivity = configOverride.OptimizeActivity
	}
	if configOverride.OptimizeTimer != nil {
		workflowConfig.OptimizeTimer = configOverride.OptimizeTimer
	}
}

func (s *serviceImpl) ApiV1WorkflowWaitForStateCompletion(
	ctx context.Context, req iwfidl.WorkflowWaitForStateCompletionRequest,
) (wresp *iwfidl.WorkflowWaitForStateCompletionResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	var workflowId string

	sharedConfig := env.GetSharedConfig()
	waitForOn := sharedConfig.GetWaitForOnWithDefault()

	if waitForOn == "old" {
		workflowId = utils.GetWorkflowIdForWaitForStateExecution(req.WorkflowId, req.StateExecutionId, req.WaitForKey, req.StateId)
	} else { // waitForOn == "new"
		var parentId string
		if s.client.GetBackendType() == service.BackendTypeTemporal { // Temporal
			response, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), "", nil)
			if err != nil {
				return nil, s.handleError(err, WorkflowWaitForStateCompletionApiPath, req.WorkflowId)
			}
			parentId = response.FirstRunId
		} else { // Cadence
			parentId = req.WorkflowId
		}

		workflowId = utils.GetWorkflowIdForWaitForStateExecution(parentId, req.StateExecutionId, req.WaitForKey, req.StateId)
	}

	options := uclient.StartWorkflowOptions{
		ID:        workflowId,
		TaskQueue: s.taskQueue,
		// TODO: https://github.com/indeedeng/iwf-java-sdk/issues/218
		// it doesn't seem to have a way for SDK to know the timeout at this API
		// So hardcoded to 1 hour for now. If it timeouts, the IDReusePolicy will restart a new one
		WorkflowExecutionTimeout: 60 * time.Minute,
	}

	runId, err := s.client.StartWaitForStateCompletionWorkflow(ctx, options)

	if err != nil {
		return nil, s.handleError(err, WorkflowWaitForStateCompletionApiPath, req.WorkflowId)
	}

	subCtx, cancFunc := utils.TrimContextByTimeoutWithCappedDDL(ctx, req.WaitTimeSeconds, s.config.Api.MaxWaitSeconds)
	defer cancFunc()
	var output service.WaitForStateCompletionWorkflowOutput
	getErr := s.client.GetWorkflowResult(subCtx, &output, workflowId, runId)

	if s.client.IsRequestTimeoutError(getErr) || s.client.IsWorkflowTimeoutError(getErr) {
		// the workflow is still running, but the wait has exceeded limit
		return nil, errors.NewErrorAndStatus(
			service.HttpStatusCodeSpecial4xxError1,
			iwfidl.LONG_POLL_TIME_OUT_SUB_STATUS,
			"waiting has exceeded timeout limit, please retry")
	}

	if getErr != nil {
		return nil, s.handleError(getErr, WorkflowWaitForStateCompletionApiPath, req.WorkflowId)
	}

	return &iwfidl.WorkflowWaitForStateCompletionResponse{
		StateCompletionOutput: &output.StateCompletionOutput,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSignalPost(
	ctx context.Context, req iwfidl.WorkflowSignalRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.SignalWorkflow(ctx,
		req.GetWorkflowId(), req.GetWorkflowRunId(), req.GetSignalChannelName(), req.GetSignalValue())
	if err != nil {
		return s.handleError(err, WorkflowSignalApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowPublishToInternalChannelPost(
	ctx context.Context, req iwfidl.PublishToInternalChannelRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	sigVal := service.ExecuteRpcSignalRequest{
		InterStateChannelPublishing: req.GetMessages(),
	}
	err := s.client.SignalWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteRpcSignalChannelName, sigVal)
	if err != nil {
		return s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowConfigUpdate(
	ctx context.Context, req iwfidl.WorkflowConfigUpdateRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.SignalWorkflow(ctx,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.UpdateConfigSignalChannelName, req)
	if err != nil {
		return s.handleError(err, WorkflowConfigUpdateApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowTriggerContinueAsNew(
	ctx context.Context, req iwfidl.TriggerContinueAsNewRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.SignalWorkflow(ctx,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.TriggerContinueAsNewSignalChannelName, nil)
	if err != nil {
		return s.handleError(err, WorkflowTriggerContinueAsNewApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowStopPost(
	ctx context.Context, req iwfidl.WorkflowStopRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	wfId := req.GetWorkflowId()
	runId := req.GetWorkflowRunId()
	stopType := iwfidl.CANCEL
	if req.StopType != nil {
		stopType = req.GetStopType()
	}

	switch stopType {
	case iwfidl.CANCEL:
		err := s.client.CancelWorkflow(ctx, wfId, runId)
		if err != nil {
			return s.handleError(err, WorkflowStopApiPath, req.GetWorkflowId())
		}
	case iwfidl.TERMINATE:
		err := s.client.TerminateWorkflow(ctx, wfId, runId, req.GetReason())
		if err != nil {
			return s.handleError(err, WorkflowStopApiPath, req.GetWorkflowId())
		}
	case iwfidl.FAIL:
		err := s.client.SignalWorkflow(ctx, wfId, runId, service.FailWorkflowSignalChannelName, service.FailWorkflowSignalRequest{Reason: req.GetReason()})
		if err != nil {
			return s.handleError(err, WorkflowStopApiPath, req.GetWorkflowId())
		}
	default:
		return s.handleError(
			fmt.Errorf("unsupported stop type: %v", stopType),
			WorkflowStopApiPath, req.GetWorkflowId())
	}

	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetQueryAttributesPost(
	ctx context.Context, req iwfidl.WorkflowGetDataObjectsRequest,
) (wresp *iwfidl.WorkflowGetDataObjectsResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	var queryResp service.GetDataAttributesQueryResponse
	queryToGetDataAttributes := true
	if req.GetUseMemoForDataAttributes() {
		requestedKeys := map[string]bool{}
		for _, k := range req.Keys {
			requestedKeys[k] = true
		}
		// Note that when the requested keys is empty, it means all

		var dataAttributes []iwfidl.KeyValue

		response, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), nil)
		if err != nil {
			return nil, s.handleError(err, WorkflowGetDataAttributesApiPath, req.GetWorkflowId())
		}

		for k, v := range response.Memos {
			if strings.HasPrefix(k, service.IwfSystemConstPrefix) {
				continue
			}
			if len(requestedKeys) > 0 && !requestedKeys[k] {
				continue
			}
			dataAttributes = append(dataAttributes, iwfidl.KeyValue{
				Key:   iwfidl.PtrString(k),
				Value: ptr.Any(v),
			})
		}

		_, ok := response.Memos[service.UseMemoForDataAttributesKey]
		if ok {
			// using memo is enough
			queryToGetDataAttributes = false
		} else {
			// this means that we cannot use memo to continue, need to fall back to use query
			s.logger.Warn("workflow attempt to use memo but probably isn't started with it", tag.WorkflowID(req.WorkflowId))
			if s.config.Interpreter.FailAtMemoIncompatibility {
				return nil, s.handleError(
					fmt.Errorf("memo is not set correctly to use"),
					WorkflowGetDataAttributesApiPath, req.GetWorkflowId())
			}
		}

		queryResp = service.GetDataAttributesQueryResponse{
			DataAttributes: dataAttributes,
		}
	}

	if queryToGetDataAttributes {
		err := s.client.QueryWorkflow(ctx, &queryResp,
			req.GetWorkflowId(), req.GetWorkflowRunId(), service.GetDataAttributesWorkflowQueryType,
			service.GetDataAttributesQueryRequest{
				Keys: req.Keys,
			})

		if err != nil {
			return nil, s.handleError(err, WorkflowGetDataAttributesApiPath, req.GetWorkflowId())
		}
	}

	return &iwfidl.WorkflowGetDataObjectsResponse{
		Objects: queryResp.DataAttributes,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSetQueryAttributesPost(
	ctx context.Context, req iwfidl.WorkflowSetDataObjectsRequest,
) (retError *errors.ErrorAndStatus) {
	sigVal := service.ExecuteRpcSignalRequest{
		UpsertDataObjects: req.Objects,
	}

	err := s.client.SignalWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteRpcSignalChannelName, sigVal)
	if err != nil {
		return s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetSearchAttributesPost(
	ctx context.Context, req iwfidl.WorkflowGetSearchAttributesRequest,
) (wresp *iwfidl.WorkflowGetSearchAttributesResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	response, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), req.Keys)
	if err != nil {
		return nil, s.handleError(err, WorkflowGetSearchAttributesApiPath, req.GetWorkflowId())
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

func (s *serviceImpl) ApiV1WorkflowSetSearchAttributesPost(
	ctx context.Context, req iwfidl.WorkflowSetSearchAttributesRequest,
) (retError *errors.ErrorAndStatus) {
	sigVal := service.ExecuteRpcSignalRequest{
		UpsertSearchAttributes: req.SearchAttributes,
	}

	err := s.client.SignalWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteRpcSignalChannelName, sigVal)
	if err != nil {
		return s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetPost(
	ctx context.Context, req iwfidl.WorkflowGetRequest,
) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	return s.doApiV1WorkflowGetPost(ctx, req, false)
}

func (s *serviceImpl) ApiV1WorkflowGetWithWaitPost(
	ctx context.Context, req iwfidl.WorkflowGetRequest,
) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	return s.doApiV1WorkflowGetPost(ctx, req, true)
}

// withWait:
//
//	 because s.client.GetWorkflowResult will wait for the completion if workflow is running --
//		when withWait is false, if workflow is not running and needResults is true, it will then call s.client.GetWorkflowResult to get results
//		when withWait is true, it will do everything
func (s *serviceImpl) doApiV1WorkflowGetPost(
	ctx context.Context, req iwfidl.WorkflowGetRequest, withWait bool,
) (wresp *iwfidl.WorkflowGetResponse, retError *errors.ErrorAndStatus) {
	descResp, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), nil)
	if err != nil {
		return nil, s.handleError(err, WorkflowGetApiPath, req.GetWorkflowId())
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
		subCtx, cancFunc := utils.TrimContextByTimeoutWithCappedDDL(ctx, req.WaitTimeSeconds, s.config.Api.MaxWaitSeconds)
		defer cancFunc()
		getErr = s.client.GetWorkflowResult(subCtx, &output, req.GetWorkflowId(), req.GetWorkflowRunId())
		if getErr == nil {
			status = iwfidl.COMPLETED
		}
	}

	if getErr == nil {
		return &iwfidl.WorkflowGetResponse{
			WorkflowRunId:  descResp.RunId,
			WorkflowStatus: status,
			Results:        output.StateCompletionOutputs,
		}, nil
	}

	if s.client.IsRequestTimeoutError(getErr) {
		// the workflow is still running, but the wait has exceeded limit
		return nil, errors.NewErrorAndStatus(
			service.HttpStatusCodeSpecial4xxError1,
			iwfidl.LONG_POLL_TIME_OUT_SUB_STATUS,
			"workflow is still running, waiting has exceeded timeout limit, please retry")
	}

	var outputsToReturnWf []iwfidl.StateCompletionOutput
	var errMsg string
	errType := s.client.GetApplicationErrorTypeIfIsApplicationError(getErr)
	if errType != "" {
		// workflow failed by interpreter decision, or by user workflow state decision
		errTypeEnum := iwfidl.WorkflowErrorType(errType)
		if errTypeEnum == iwfidl.STATE_DECISION_FAILING_WORKFLOW_ERROR_TYPE {
			err = s.client.GetApplicationErrorDetails(getErr, &outputsToReturnWf)
			if err != nil {
				return nil, s.handleError(err, WorkflowGetApiPath, req.GetWorkflowId())
			}
		} else {
			err = s.client.GetApplicationErrorDetails(getErr, &errMsg)
			if err != nil {
				return nil, s.handleError(err, WorkflowGetApiPath, req.GetWorkflowId())
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
			return nil, s.handleError(err, WorkflowGetApiPath, req.GetWorkflowId())
		}
		errMsg = ""
		if descResp.Status == iwfidl.RUNNING || descResp.Status == iwfidl.CONTINUED_AS_NEW || descResp.Status == iwfidl.COMPLETED {
			errMsg = "impossible/very rare status, maybe caused by racing conditions"
			s.logger.Error(errMsg, tag.WorkflowID(req.GetWorkflowId()), tag.WorkflowRunID(descResp.RunId))
			// we cannot return these status, which will be a wrong results
			// TODO: maybe return 4xx
			return nil, s.handleError(fmt.Errorf(errMsg), WorkflowGetApiPath, req.GetWorkflowId())
		}

		if descResp.Status == iwfidl.FAILED {
			errMsg = "unknown workflow failure from interpreter implementation"
			s.logger.Error(errMsg, tag.WorkflowID(req.GetWorkflowId()), tag.WorkflowRunID(descResp.RunId))
		}

		var errMsgPtr *string
		if errMsg != "" {
			errMsgPtr = ptr.Any(errMsg)
		}

		return &iwfidl.WorkflowGetResponse{
			WorkflowRunId:  descResp.RunId,
			WorkflowStatus: descResp.Status,
			ErrorMessage:   errMsgPtr,
		}, nil
	}

}

func (s *serviceImpl) ApiV1WorkflowSearchPost(
	ctx context.Context, req iwfidl.WorkflowSearchRequest,
) (wresp *iwfidl.WorkflowSearchResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	pageSize := int32(1000)
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	resp, err := s.client.ListWorkflow(ctx, &uclient.ListWorkflowExecutionsRequest{
		PageSize:      pageSize,
		Query:         req.GetQuery(),
		NextPageToken: []byte(req.GetNextPageToken()),
	})
	if err != nil {
		return nil, s.handleError(err, WorkflowSearchApiPath, "N/A")
	}
	return &iwfidl.WorkflowSearchResponse{
		WorkflowExecutions: resp.Executions,
		NextPageToken:      ptr.Any(string(resp.NextPageToken)),
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowRpcPost(
	ctx context.Context, req iwfidl.WorkflowRpcRequest,
) (wresp *iwfidl.WorkflowRpcResponse, retError *errors.ErrorAndStatus) {
	defer func() {
		log.CapturePanic(recover(), s.logger, &retError)
	}()

	if needLocking(req) {
		return s.handleRpcBySynchronousUpdate(ctx, req)
	}

	var rpcPrep *service.PrepareRpcQueryResponse

	err := s.client.QueryWorkflow(ctx, &rpcPrep, req.GetWorkflowId(), req.GetWorkflowRunId(), service.PrepareRpcQueryType, service.PrepareRpcQueryRequest{
		DataObjectsLoadingPolicy:      req.DataAttributesLoadingPolicy,
		SearchAttributesLoadingPolicy: req.SearchAttributesLoadingPolicy,
	})
	if err != nil {
		return nil, s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
	}

	stateApiExecuteStartTime := time.Now().UnixMilli()

	defer func() {
		event.Handle(iwfidl.IwfEvent{
			EventType:          iwfidl.RPC_EXECUTION_EVENT,
			RpcName:            &req.RpcName,
			WorkflowType:       rpcPrep.IwfWorkflowType,
			WorkflowId:         req.GetWorkflowId(),
			StartTimestampInMs: ptr.Any(stateApiExecuteStartTime),
			// search attributes are not available at this time
		})
	}()

	resp, retError := rpc.InvokeWorkerRpc(ctx, rpcPrep, req, s.config.Api.MaxWaitSeconds)
	if retError != nil {
		return nil, retError
	}

	decision := resp.GetStateDecision()
	if len(resp.UpsertDataAttributes)+len(resp.UpsertSearchAttributes)+len(resp.PublishToInterStateChannel)+len(resp.RecordEvents)+len(decision.GetNextStates()) > 0 {
		// if there is no mutation on the workflow, this RPC is "readonly", then don't send the signal

		// send the signal
		sigVal := service.ExecuteRpcSignalRequest{
			RpcInput:                    req.Input,
			RpcOutput:                   resp.Output,
			UpsertDataObjects:           resp.UpsertDataAttributes,
			UpsertSearchAttributes:      resp.UpsertSearchAttributes,
			StateDecision:               resp.StateDecision,
			RecordEvents:                resp.RecordEvents,
			InterStateChannelPublishing: resp.PublishToInterStateChannel,
		}
		if s.config.Api.OmitRpcInputOutputInHistory != nil && *s.config.Api.OmitRpcInputOutputInHistory {
			// the input/output is only for debugging purpose but could be too expensive to store
			sigVal.RpcInput = nil
			sigVal.RpcOutput = nil
		}
		err := s.client.SignalWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteRpcSignalChannelName, sigVal)
		if err != nil {
			return nil, s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
		}
	}

	return &iwfidl.WorkflowRpcResponse{Output: resp.Output}, nil
}

func needLocking(req iwfidl.WorkflowRpcRequest) bool {
	if req.SearchAttributesLoadingPolicy != nil {
		if doNeedLocking(req.SearchAttributesLoadingPolicy) {
			return true
		}
	}
	if req.DataAttributesLoadingPolicy != nil {
		if doNeedLocking(req.DataAttributesLoadingPolicy) {
			return true
		}
	}
	return false
}

func (s *serviceImpl) handleRpcBySynchronousUpdate(
	ctx context.Context, req iwfidl.WorkflowRpcRequest,
) (resp *iwfidl.WorkflowRpcResponse, retError *errors.ErrorAndStatus) {
	req.TimeoutSeconds = ptr.Any(utils.TrimRpcTimeoutSeconds(ctx, req))
	var output interfaces.HandlerOutput
	err := s.client.SynchronousUpdateWorkflow(ctx, &output, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteOptimisticLockingRpcUpdateType, req)
	if err != nil {
		errType := s.client.GetApplicationErrorTypeIfIsApplicationError(err)
		if errType != "" {
			errTypeEnum := iwfidl.WorkflowErrorType(errType)
			if errTypeEnum == iwfidl.RPC_ACQUIRE_LOCK_FAILURE {
				var details string
				err2 := s.client.GetApplicationErrorDetails(err, &details)
				if err2 != nil {
					details = err2.Error()
				}
				return nil, errors.NewErrorAndStatus(service.HttpStatusCodeSpecial4xxError2, iwfidl.WORKER_API_ERROR, details)
			}
		}
		return nil, s.handleError(err, WorkflowRpcApiPath, req.GetWorkflowId())
	}
	return output.RpcOutput, output.StatusError
}

func doNeedLocking(policy *iwfidl.PersistenceLoadingPolicy) bool {
	if policy.GetPersistenceLoadingType() == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK ||
		policy.GetPersistenceLoadingType() == iwfidl.ALL_WITH_PARTIAL_LOCK {
		return true
	}
	return false
}

func (s *serviceImpl) ApiV1WorkflowResetPost(
	ctx context.Context, req iwfidl.WorkflowResetRequest,
) (wresp *iwfidl.WorkflowResetResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	runId, err := s.client.ResetWorkflow(ctx, req)
	if err != nil {
		return nil, s.handleError(err, WorkflowResetApiPath, req.GetWorkflowId())
	}
	return &iwfidl.WorkflowResetResponse{
		WorkflowRunId: runId,
	}, nil
}

func (s *serviceImpl) ApiV1WorkflowSkipTimerPost(
	ctx context.Context, request iwfidl.WorkflowSkipTimerRequest,
) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	if request.GetTimerCommandId() == "" && request.TimerCommandIndex == nil {
		return makeInvalidRequestError("must provide either a timerCommandId or index")
	}

	timerInfos := service.GetCurrentTimerInfosQueryResponse{}
	err := s.client.QueryWorkflow(ctx, &timerInfos, request.GetWorkflowId(), request.GetWorkflowRunId(), service.GetCurrentTimerInfosQueryType)
	if err != nil {
		return s.handleError(err, WorkflowSkipTimerApiPath, request.GetWorkflowId())
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
		return s.handleError(err, WorkflowSkipTimerApiPath, request.GetWorkflowId())
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowDumpPost(
	ctx context.Context, request iwfidl.WorkflowDumpRequest,
) (*iwfidl.WorkflowDumpResponse, *errors.ErrorAndStatus) {
	var pageOfSnapshot iwfidl.WorkflowDumpResponse

	err := s.client.QueryWorkflow(ctx, &pageOfSnapshot, request.GetWorkflowId(), request.GetWorkflowRunId(), service.ContinueAsNewDumpByPageQueryType, request)
	if err != nil {
		return nil, s.handleError(err, WorkflowInternalDumpApiPath, request.GetWorkflowId())
	}
	return &pageOfSnapshot, nil
}

func (s *serviceImpl) ApiInfoHealth(ctx context.Context) *iwfidl.HealthInfo {
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "Hostname Not Available"
	}
	return &iwfidl.HealthInfo{
		Condition: iwfidl.PtrString("OK"),
		Hostname:  iwfidl.PtrString(hostName),
		Duration:  iwfidl.PtrInt32(0),
	}
}

func makeInvalidRequestError(msg string) *errors.ErrorAndStatus {
	return errors.NewErrorAndStatus(http.StatusBadRequest,
		iwfidl.UNCATEGORIZED_SUB_STATUS,
		"invalid request - "+msg)
}

func (s *serviceImpl) handleError(err error, apiPath string, workflowId string) *errors.ErrorAndStatus {
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
	if workflowId == "" && status == http.StatusInternalServerError {
		status = http.StatusBadRequest
		subStatus = iwfidl.WORKFLOW_NOT_EXISTS_SUB_STATUS
	}
	if status >= 500 {
		s.logger.Error("encounter server error for API",
			tag.StatusCode(status), tag.Error(err),
			tag.Name(apiPath), tag.WorkflowID(workflowId),
			tag.SubStatus(string(subStatus)))
	} else {
		s.logger.Warn("encounter client error for API",
			tag.StatusCode(status), tag.Error(err),
			tag.Name(apiPath), tag.WorkflowID(workflowId),
			tag.SubStatus(string(subStatus)))
	}

	return errors.NewErrorAndStatus(
		status,
		subStatus,
		err.Error(),
	)
}
