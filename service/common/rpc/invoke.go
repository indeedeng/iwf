package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/blobstore"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/urlautofix"
	"github.com/indeedeng/iwf/service/common/utils"
)

func InvokeWorkerRpc(
	ctx context.Context, rpcPrep *service.PrepareRpcQueryResponse, req iwfidl.WorkflowRpcRequest, apiMaxSeconds int64, blobStore blobstore.BlobStore, externalStorageConfig config.ExternalStorageConfig,
) (*iwfidl.WorkflowWorkerRpcResponse, *errors.ErrorAndStatus) {
	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(rpcPrep.IwfWorkerUrl)
	// invoke worker rpc
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})

	err := blobstore.LoadDataObjectsFromExternalStorage(ctx, rpcPrep.DataObjects, blobStore)

	rpcCtx, cancel := utils.TrimContextByTimeoutWithCappedDDL(ctx, req.TimeoutSeconds, apiMaxSeconds)
	defer cancel()
	workerReq := apiClient.DefaultApi.ApiV1WorkflowWorkerRpcPost(rpcCtx)

	// creating empty maps for signalChannelInfos & internalChannelInfos instead of passing in nils
	// using nil causes problems when converting to map model defined with OpenAPI
	var signalChannelInfos map[string]iwfidl.ChannelInfo
	if rpcPrep.SignalChannelInfo == nil {
		signalChannelInfos = make(map[string]iwfidl.ChannelInfo)
	} else {
		signalChannelInfos = rpcPrep.SignalChannelInfo
	}

	var internalChannelInfos map[string]iwfidl.ChannelInfo
	if rpcPrep.InternalChannelInfo == nil {
		internalChannelInfos = make(map[string]iwfidl.ChannelInfo)
	} else {
		internalChannelInfos = rpcPrep.InternalChannelInfo
	}

	workerRequest := iwfidl.WorkflowWorkerRpcRequest{
		Context: iwfidl.Context{
			WorkflowId:               req.WorkflowId,
			WorkflowRunId:            rpcPrep.WorkflowRunId,
			WorkflowStartedTimestamp: rpcPrep.WorkflowStartedTimestamp,
		},
		WorkflowType:         rpcPrep.IwfWorkflowType,
		RpcName:              req.RpcName,
		Input:                req.Input,
		SearchAttributes:     rpcPrep.SearchAttributes,
		DataAttributes:       rpcPrep.DataObjects,
		SignalChannelInfos:   &signalChannelInfos,
		InternalChannelInfos: &internalChannelInfos,
	}
	resp, httpResp, err := workerReq.WorkflowWorkerRpcRequest(workerRequest).Execute()
	if utils.CheckHttpError(err, httpResp) {
		return nil, handleWorkerRpcResponseError(err, httpResp)
	}
	decision := resp.GetStateDecision()
	if decision.HasConditionalClose() {
		return nil, handleWorkerRpcResponseError(fmt.Errorf("closing workflow in RPC is not supported yet"), nil)
	}

	if resp.UpsertDataAttributes != nil {
		err = blobstore.WriteDataObjectsToExternalStorage(ctx, resp.UpsertDataAttributes, req.WorkflowId, externalStorageConfig.ThresholdInBytes, blobStore, externalStorageConfig.Enabled)
		if err != nil {
			return nil, handleWorkerRpcResponseError(err, nil)
		}
	}

	for _, st := range decision.GetNextStates() {
		if service.ValidClosingWorkflowStateId[st.GetStateId()] {
			// TODO this need more work in workflow to support
			return nil, handleWorkerRpcResponseError(fmt.Errorf("closing workflow in RPC is not supported yet"), nil)
		}
	}
	return resp, nil
}

func handleWorkerRpcResponseError(err error, httpResp *http.Response) *errors.ErrorAndStatus {
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
		service.HttpStatusCodeSpecial4xxError1,
		iwfidl.WORKER_API_ERROR,
		detailedMessage,
		workerError.GetDetail(),
		workerError.GetErrorType(),
		int32(originalStatusCode),
	)
}
