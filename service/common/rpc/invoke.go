package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/urlautofix"
	"github.com/indeedeng/iwf/service/common/utils"
	"io/ioutil"
	"net/http"
)

func InvokeWorkerRpc(ctx context.Context, rpcPrep *service.PrepareRpcQueryResponse, req iwfidl.WorkflowRpcRequest, apiMaxSeconds int64) (*iwfidl.WorkflowWorkerRpcResponse, *errors.ErrorAndStatus) {
	iwfWorkerBaseUrl := urlautofix.FixWorkerUrl(rpcPrep.IwfWorkerUrl)
	// invoke worker rpc
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: iwfWorkerBaseUrl,
			},
		},
	})

	rpcCtx, cancel := utils.TrimContextByTimeoutWithCappedDDL(ctx, req.TimeoutSeconds, apiMaxSeconds)
	defer cancel()
	workerReq := apiClient.DefaultApi.ApiV1WorkflowWorkerRpcPost(rpcCtx)
	workerRequest := iwfidl.WorkflowWorkerRpcRequest{
		Context: iwfidl.Context{
			WorkflowId:               req.WorkflowId,
			WorkflowRunId:            rpcPrep.WorkflowRunId,
			WorkflowStartedTimestamp: rpcPrep.WorkflowStartedTimestamp,
		},
		WorkflowType:     rpcPrep.IwfWorkflowType,
		RpcName:          req.RpcName,
		Input:            req.Input,
		SearchAttributes: rpcPrep.SearchAttributes,
		DataAttributes:   rpcPrep.DataObjects,
	}
	resp, httpResp, err := workerReq.WorkflowWorkerRpcRequest(workerRequest).Execute()
	if utils.CheckHttpError(err, httpResp) {
		return nil, handleWorkerRpcResponseError(err, httpResp)
	}
	decision := resp.GetStateDecision()
	if decision.HasConditionalClose() {
		return nil, handleWorkerRpcResponseError(fmt.Errorf("closing workflow in RPC is not supported yet"), nil)
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
