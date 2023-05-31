package api

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/indeedeng/iwf/service/common/compatibility"
	"github.com/indeedeng/iwf/service/common/urlautofix"
	"io/ioutil"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/log/tag"
	"github.com/indeedeng/iwf/service/common/mapper"
	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type serviceImpl struct {
	client    UnifiedClient
	taskQueue string
	logger    log.Logger
	config    config.Config
}

func (s *serviceImpl) Close() {
	s.client.Close()
}

func NewApiService(config config.Config, client UnifiedClient, taskQueue string, logger log.Logger) (ApiService, error) {
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
	workflowConfig := s.config.Interpreter.DefaultWorkflowConfig

	useMemo := false
	if req.WorkflowStartOptions != nil {
		startOptions := req.WorkflowStartOptions
		workflowOptions.WorkflowIDReusePolicy = compatibility.GetWorkflowIdReusePolicy(*startOptions)
		workflowOptions.CronSchedule = startOptions.CronSchedule
		workflowOptions.RetryPolicy = startOptions.RetryPolicy
		var err error
		workflowOptions.SearchAttributes, err = mapper.MapToInternalSearchAttributes(startOptions.SearchAttributes)
		if err != nil {
			return nil, s.handleError(err)
		}
		workflowOptions.SearchAttributes[service.SearchAttributeIwfWorkflowType] = req.IwfWorkflowType
		initSAs = startOptions.SearchAttributes
		if startOptions.HasWorkflowConfigOverride() {
			workflowConfig = startOptions.GetWorkflowConfigOverride()
		}
		if startOptions.GetUseMemoForDataAttributes() {
			workflowOptions.Memo = map[string]interface{}{
				service.WorkerUrlMemoKey: iwfidl.EncodedObject{
					Data: iwfidl.PtrString(req.IwfWorkerUrl), // this is hack to ensure all memos are with the same type
				},
			}
			useMemo = true
		}
	}

	input := service.InterpreterWorkflowInput{
		IwfWorkflowType:          req.GetIwfWorkflowType(),
		IwfWorkerUrl:             req.GetIwfWorkerUrl(),
		StartStateId:             req.StartStateId,
		StateInput:               req.GetStateInput(),
		StateOptions:             req.GetStateOptions(),
		InitSearchAttributes:     initSAs,
		Config:                   workflowConfig,
		UseMemoForDataAttributes: useMemo,
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

func (s *serviceImpl) ApiV1WorkflowConfigUpdate(ctx context.Context, req iwfidl.WorkflowConfigUpdateRequest) (retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	err := s.client.SignalWorkflow(ctx,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.UpdateConfigSignalChannelName, req)
	if err != nil {
		return s.handleError(err)
	}
	return nil
}

func (s *serviceImpl) ApiV1WorkflowStopPost(ctx context.Context, req iwfidl.WorkflowStopRequest) (retError *errors.ErrorAndStatus) {
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
			return s.handleError(err)
		}
	case iwfidl.TERMINATE:
		err := s.client.TerminateWorkflow(ctx, wfId, runId, req.GetReason())
		if err != nil {
			return s.handleError(err)
		}
	case iwfidl.FAIL:
		err := s.client.SignalWorkflow(ctx, wfId, runId, service.FailWorkflowSignalChannelName, service.FailWorkflowSignalRequest{Reason: req.GetReason()})
		if err != nil {
			return s.handleError(err)
		}
	default:
		return s.handleError(fmt.Errorf("unsupported stop type: %v", stopType))
	}

	return nil
}

func (s *serviceImpl) ApiV1WorkflowGetQueryAttributesPost(ctx context.Context, req iwfidl.WorkflowGetDataObjectsRequest) (wresp *iwfidl.WorkflowGetDataObjectsResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	var queryResp service.GetDataObjectsQueryResponse
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
			return nil, s.handleError(err)
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

		_, ok := response.Memos[service.WorkerUrlMemoKey]
		if ok {
			// using memo is enough
			queryToGetDataAttributes = false
		} else {
			// this means that we cannot use memo to continue, need to fall back to use query
			s.logger.Warn("workflow attempt to use memo but probably isn't started with it", tag.WorkflowID(req.WorkflowId))
			if s.config.Interpreter.FailAtMemoIncompatibility {
				return nil, s.handleError(fmt.Errorf("memo is not set correctly to use"))
			}
		}

		queryResp = service.GetDataObjectsQueryResponse{
			DataObjects: dataAttributes,
		}
	}

	if queryToGetDataAttributes {
		err := s.client.QueryWorkflow(ctx, &queryResp,
			req.GetWorkflowId(), req.GetWorkflowRunId(), service.GetDataObjectsWorkflowQueryType,
			service.GetDataObjectsQueryRequest{
				Keys: req.Keys,
			})

		if err != nil {
			return nil, s.handleError(err)
		}
	}

	return &iwfidl.WorkflowGetDataObjectsResponse{
		Objects: queryResp.DataObjects,
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

func (s *serviceImpl) ApiV1WorkflowRpcPost(ctx context.Context, req iwfidl.WorkflowRpcRequest) (wresp *iwfidl.WorkflowRpcResponse, retError *errors.ErrorAndStatus) {
	defer func() { log.CapturePanic(recover(), s.logger, &retError) }()

	if err := checkPersistenceLoadingPolicy(req); err != nil {
		return nil, s.handleError(err)
	}

	var queryResp service.PrepareRpcQueryResponse
	queryToPrepare := true
	if req.GetUseMemoForDataAttributes() {
		var searchAttributes []iwfidl.SearchAttribute
		var dataAttributes []iwfidl.KeyValue

		requestedSAs := req.SearchAttributes
		saPolicy := req.GetSearchAttributesLoadingPolicy()
		if saPolicy.GetPersistenceLoadingType() != iwfidl.ALL_WITHOUT_LOCKING {
			requestedSAKeys := map[string]bool{}
			for _, saKey := range saPolicy.PartialLoadingKeys {
				requestedSAKeys[saKey] = true
			}
			requestedSAs = []iwfidl.SearchAttributeKeyAndType{}
			for _, sa := range req.SearchAttributes {
				if requestedSAKeys[sa.GetKey()] {
					requestedSAs = append(requestedSAs, sa)
				}
			}
		}

		requestedSAs = append(requestedSAs, iwfidl.SearchAttributeKeyAndType{
			Key:       ptr.Any(service.SearchAttributeIwfWorkflowType),
			ValueType: iwfidl.KEYWORD.Ptr(),
		})
		response, err := s.client.DescribeWorkflowExecution(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), requestedSAs)
		if err != nil {
			return nil, s.handleError(err)
		}

		for _, sa := range requestedSAs {
			if sa.GetKey() == service.SearchAttributeIwfWorkflowType {
				continue
			}
			searchAttribute, exist := response.SearchAttributes[sa.GetKey()]
			if exist {
				searchAttributes = append(searchAttributes, searchAttribute)
			}
		}

		for k, v := range response.Memos {
			if strings.HasPrefix(k, service.IwfSystemConstPrefix) {
				continue
			}
			dataAttributes = append(dataAttributes, iwfidl.KeyValue{
				Key:   iwfidl.PtrString(k),
				Value: ptr.Any(v), //NOTE: using &v is WRONG: must avoid using & for the iteration item
			})
		}

		attribute := response.SearchAttributes[service.SearchAttributeIwfWorkflowType]
		workflowType := attribute.GetStringValue()
		workerUrlMemoObj, ok := response.Memos[service.WorkerUrlMemoKey]
		if ok {
			// using memo is enough
			queryToPrepare = false
		} else {
			// this means that we cannot use memo to continue, need to fall back to use query
			s.logger.Warn("workflow attempt to use memo but probably isn't started with it", tag.WorkflowID(req.WorkflowId))
			if s.config.Interpreter.FailAtMemoIncompatibility {
				return nil, s.handleError(fmt.Errorf("memo is not set correctly to use"))
			}
		}
		workerUrl := workerUrlMemoObj.GetData()

		queryResp = service.PrepareRpcQueryResponse{
			DataObjects:              dataAttributes,
			SearchAttributes:         searchAttributes,
			WorkflowStartedTimestamp: response.WorkflowStartedTimestamp,
			WorkflowRunId:            response.RunId,
			IwfWorkflowType:          workflowType,
			IwfWorkerUrl:             workerUrl,
		}
	}

	if queryToPrepare {
		// use query to load, this is expensive. So it tries to avoid if possible
		err := s.client.QueryWorkflow(ctx, &queryResp, req.GetWorkflowId(), req.GetWorkflowRunId(), service.PrepareRpcQueryType, service.PrepareRpcQueryRequest{
			DataObjectsLoadingPolicy:      req.DataAttributesLoadingPolicy,
			SearchAttributesLoadingPolicy: req.SearchAttributesLoadingPolicy,
		})
		if err != nil {
			return nil, s.handleError(err)
		}
	}

	iwfWorkerBaseUrl := urlautofix.GetIwfWorkerBaseUrlWithFix(queryResp.IwfWorkerUrl)
	// invoke worker rpc
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})
	rpcCtx := ctx
	var cancel context.CancelFunc
	if req.GetTimeoutSeconds() > 0 {
		rpcCtx, cancel = context.WithTimeout(rpcCtx, time.Duration(req.GetTimeoutSeconds())*time.Second)
		defer cancel()
	}
	workerReq := apiClient.DefaultApi.ApiV1WorkflowWorkerRpcPost(rpcCtx)
	workerRequest := iwfidl.WorkflowWorkerRpcRequest{
		Context: iwfidl.Context{
			WorkflowId:               req.WorkflowId,
			WorkflowRunId:            queryResp.WorkflowRunId,
			WorkflowStartedTimestamp: queryResp.WorkflowStartedTimestamp,
		},
		WorkflowType:     queryResp.IwfWorkflowType,
		RpcName:          req.RpcName,
		Input:            req.Input,
		SearchAttributes: queryResp.SearchAttributes,
		DataAttributes:   queryResp.DataObjects,
	}
	resp, httpResp, err := workerReq.WorkflowWorkerRpcRequest(workerRequest).Execute()
	if checkHttpError(err, httpResp) {
		return nil, s.handleWorkerRpcApiError(err, httpResp)
	}
	decision := resp.GetStateDecision()
	for _, st := range decision.GetNextStates() {
		if service.ValidClosingWorkflowStateId[st.GetStateId()] {
			// TODO this need more work in workflow to support
			return nil, s.handleError(fmt.Errorf("closing workflow in RPC is not supported yet"))
		}
	}

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
		err = s.client.SignalWorkflow(ctx, req.GetWorkflowId(), req.GetWorkflowRunId(), service.ExecuteRpcSignalChannelName, sigVal)
		if err != nil {
			return nil, s.handleError(err)
		}
	}

	return &iwfidl.WorkflowRpcResponse{Output: resp.Output}, nil
}

func checkPersistenceLoadingPolicy(req iwfidl.WorkflowRpcRequest) error {
	if req.SearchAttributesLoadingPolicy != nil {
		if err := doCheckPersistenceLoadingPolicy(req.SearchAttributesLoadingPolicy); err != nil {
			return err
		}
	}
	if req.DataAttributesLoadingPolicy != nil {
		if err := doCheckPersistenceLoadingPolicy(req.DataAttributesLoadingPolicy); err != nil {
			return err
		}
	}
	return nil
}

func doCheckPersistenceLoadingPolicy(policy *iwfidl.PersistenceLoadingPolicy) error {
	if policy.GetPersistenceLoadingType() == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		return fmt.Errorf("PARTIAL_WITH_EXCLUSIVE_LOCK is not supported in RPC yet")
	}
	return nil
}

func checkHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
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

func (s *serviceImpl) ApiV1WorkflowDumpPost(ctx context.Context, request iwfidl.WorkflowDumpRequest) (*iwfidl.WorkflowDumpResponse, *errors.ErrorAndStatus) {
	var internals service.ContinueAsNewDumpResponse

	err := s.client.QueryWorkflow(ctx, &internals, request.GetWorkflowId(), request.GetWorkflowRunId(), service.ContinueAsNewDumpQueryType)
	if err != nil {
		return nil, s.handleError(err)
	}

	data, err := json.Marshal(internals)
	if err != nil {
		return nil, s.handleError(err)
	}
	checksum := md5.Sum(data)
	pageSize := int32(service.DefaultContinueAsNewPageSizeInBytes)
	if request.PageSizeInBytes > 0 {
		pageSize = request.PageSizeInBytes
	}
	lenInDouble := float64(len(data))
	totalPages := int32(math.Ceil(lenInDouble / float64(pageSize)))
	if request.PageNum >= totalPages {
		return nil, s.handleError(fmt.Errorf("wrong pageNum, max is %v", totalPages-1))
	}
	start := pageSize * request.PageNum
	end := start + pageSize
	if end > int32(len(data)) {
		end = int32(len(data))
	}
	return &iwfidl.WorkflowDumpResponse{
		Checksum:   string(checksum[:]),
		TotalPages: totalPages,
		JsonData:   string(data[start:end]),
	}, nil
}

func makeInvalidRequestError(msg string) *errors.ErrorAndStatus {
	return errors.NewErrorAndStatus(http.StatusBadRequest,
		iwfidl.UNCATEGORIZED_SUB_STATUS,
		"invalid request - "+msg)
}

func (s *serviceImpl) handleWorkerRpcApiError(err error, httpResp *http.Response) *errors.ErrorAndStatus {
	detailedMessage := err.Error()
	if err != nil {
		detailedMessage = err.Error()
	}

	var originalStatusCode int
	var workerError iwfidl.WorkerErrorResponse
	if httpResp != nil {
		originalStatusCode = httpResp.StatusCode
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			detailedMessage = "cannot read body from http response"
		} else {
			err := json.Unmarshal(body, &workerError)
			if err != nil {
				detailedMessage = "unable to decode worker response body to WorkerErrorResponse: body" + string(body)
			} else {
				detailedMessage = fmt.Sprintf("worker API error, status:%v, errorType:%v", originalStatusCode, workerError.GetErrorType())
			}
		}

	}

	return errors.NewErrorAndStatusWithWorkerError(
		service.HttpStatusCodeWorkerApiError,
		iwfidl.WORKER_API_ERROR,
		detailedMessage,
		workerError.GetDetail(),
		workerError.GetErrorType(),
		int32(originalStatusCode),
	)
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
