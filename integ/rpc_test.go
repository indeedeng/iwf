package integ

import (
	"context"
	"encoding/json"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/rpc"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestRpcWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, false, false, nil)
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, false, false, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowTemporalWithMemo(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, true, false, nil)
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowTemporalContinueAsNewWithMemo(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, true, false, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeCadence, false, false, nil)
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeCadence, false, false, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func doTestRpcWorkflow(t *testing.T, backendType service.BackendType, useMemo, memoEncryption bool, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := rpc.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:    backendType,
		MemoEncryption: memoEncryption,
	})
	defer closeFunc2()

	// create client
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	// start a workflow
	wfId := rpc.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        rpc.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(rpc.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride:   config,
			UseMemoForDataAttributes: ptr.Any(useMemo),
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	allSearchAttributes := []iwfidl.SearchAttributeKeyAndType{
		{
			Key:       iwfidl.PtrString(rpc.TestSearchAttributeKeywordKey),
			ValueType: iwfidl.KEYWORD.Ptr(),
		},
		{
			Key:       iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
			ValueType: iwfidl.INT.Ptr(),
		},
		{
			Key:       iwfidl.PtrString(rpc.TestSearchAttributeBoolKey),
			ValueType: iwfidl.BOOL.Ptr(),
		},
	}
	time.Sleep(time.Second * 1)
	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	rpcRespReadOnly, httpResp, err := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    rpc.RPCNameReadOnly,
		Input:      &rpc.TestInput,
		SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: iwfidl.PARTIAL_WITHOUT_LOCKING.Ptr(),
			PartialLoadingKeys: []string{
				rpc.TestSearchAttributeIntKey,
			},
		},
		TimeoutSeconds:           iwfidl.PtrInt32(2),
		UseMemoForDataAttributes: ptr.Any(useMemo),
		SearchAttributes:         allSearchAttributes,
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqRpc = apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    rpc.RPCNameError,
		Input:      &rpc.TestInput,
		SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: iwfidl.PARTIAL_WITHOUT_LOCKING.Ptr(),
			PartialLoadingKeys: []string{
				rpc.TestSearchAttributeIntKey,
			},
		},
		TimeoutSeconds:           iwfidl.PtrInt32(2),
		UseMemoForDataAttributes: ptr.Any(useMemo),
		SearchAttributes:         allSearchAttributes,
	}).Execute()
	assertions.NotNil(err)
	assertions.Equalf(service.HttpStatusCodeWorkerApiError, httpResp.StatusCode, "http code")
	var errResp iwfidl.ErrorResponse
	body, err := ioutil.ReadAll(httpResp.Body)
	assertions.Nil(err)
	err = json.Unmarshal(body, &errResp)
	assertions.Equalf(iwfidl.ErrorResponse{
		Detail:                    ptr.Any("worker API error, status:502, errorType:test-type"),
		SubStatus:                 iwfidl.WORKER_API_ERROR.Ptr(),
		OriginalWorkerErrorStatus: iwfidl.PtrInt32(502),
		OriginalWorkerErrorType:   iwfidl.PtrString("test-type"),
		OriginalWorkerErrorDetail: iwfidl.PtrString("test-details"),
	}, errResp, "body")

	reqRpc = apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	rpcResp, httpResp, err := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    rpc.RPCName,
		Input:      &rpc.TestInput,
		SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: iwfidl.PARTIAL_WITHOUT_LOCKING.Ptr(),
			PartialLoadingKeys: []string{
				rpc.TestSearchAttributeIntKey,
			},
		},
		TimeoutSeconds:           iwfidl.PtrInt32(2),
		UseMemoForDataAttributes: ptr.Any(useMemo),
		SearchAttributes:         allSearchAttributes,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	respWait, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, respWait)

	history, data := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  2,
		"S2_decide": 2,
	}, history, "rpc test fail, %v", history)

	assertions.Equalf(&iwfidl.WorkflowRpcResponse{
		Output: &rpc.TestOutput,
	}, rpcResp, "rpc test fail, %v", rpcResp)

	assertions.Equalf(&iwfidl.WorkflowRpcResponse{
		Output: &rpc.TestOutput,
	}, rpcResp, "rpc test fail, %v", rpcRespReadOnly)

	assertions.Equalf(map[string]interface{}{
		rpc.RPCName + "-data-attributes": []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString(rpc.TestDataObjectKey),
				Value: &rpc.TestDataObjectVal1,
			},
		},
		rpc.RPCName + "-search-attributes": []iwfidl.SearchAttribute{
			{
				Key:          iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
				IntegerValue: iwfidl.PtrInt64(rpc.TestSearchAttributeIntValue1),
				ValueType:    ptr.Any(iwfidl.INT),
			},
		},
		rpc.RPCName + "-input":        &rpc.TestInput,
		rpc.TestInterStateChannelName: &rpc.TestInterstateChannelValue,

		rpc.RPCNameReadOnly + "-data-attributes": []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString(rpc.TestDataObjectKey),
				Value: &rpc.TestDataObjectVal1,
			},
		},
		rpc.RPCNameReadOnly + "-search-attributes": []iwfidl.SearchAttribute{
			{
				Key:          iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
				IntegerValue: iwfidl.PtrInt64(rpc.TestSearchAttributeIntValue1),
				ValueType:    ptr.Any(iwfidl.INT),
			},
		},
		rpc.RPCNameReadOnly + "-input": &rpc.TestInput,

		rpc.RPCNameError + "-data-attributes": []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString(rpc.TestDataObjectKey),
				Value: &rpc.TestDataObjectVal1,
			},
		},
		rpc.RPCNameError + "-search-attributes": []iwfidl.SearchAttribute{
			{
				Key:          iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
				IntegerValue: iwfidl.PtrInt64(rpc.TestSearchAttributeIntValue1),
				ValueType:    ptr.Any(iwfidl.INT),
			},
		},
		rpc.RPCNameError + "-input": &rpc.TestInput,
	}, data, "rpc test fail, %v", data)

	reqQry := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	allDos, httpResp, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqSearch := apiClient.DefaultApi.ApiV1WorkflowSearchattributesGetPost(context.Background())
	allSAs, httpResp, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId: wfId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(rpc.TestSearchAttributeKeywordKey),
				ValueType: ptr.Any(iwfidl.KEYWORD),
			},
			{
				Key:       iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
				ValueType: ptr.Any(iwfidl.INT),
			},
			{
				Key:       iwfidl.PtrString(rpc.TestSearchAttributeBoolKey),
				ValueType: ptr.Any(iwfidl.BOOL),
			},
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions.Equalf([]iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(rpc.TestDataObjectKey),
			Value: &rpc.TestDataObjectVal2,
		},
	}, allDos.Objects, "rpc test fail")

	assertions.ElementsMatchf([]iwfidl.SearchAttribute{
		{
			Key:         iwfidl.PtrString(rpc.TestSearchAttributeKeywordKey),
			StringValue: iwfidl.PtrString(rpc.TestSearchAttributeKeywordValue2),
			ValueType:   ptr.Any(iwfidl.KEYWORD),
		},
		{
			Key:          iwfidl.PtrString(rpc.TestSearchAttributeIntKey),
			IntegerValue: iwfidl.PtrInt64(rpc.TestSearchAttributeIntValue2),
			ValueType:    ptr.Any(iwfidl.INT),
		},
		{
			Key:       iwfidl.PtrString(rpc.TestSearchAttributeBoolKey),
			ValueType: ptr.Any(iwfidl.BOOL),
			BoolValue: iwfidl.PtrBool(false),
		},
	}, allSAs.SearchAttributes, "rpc test fail")
}
