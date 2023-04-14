package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/rpc"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestRpcWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestRpcWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestRpcWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func doTestRpcWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := rpc.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
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
		StartStateId:           rpc.State1,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

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
		TimeoutSeconds: iwfidl.PtrInt32(2),
	}).Execute()
	panicAtHttpError(err, httpResp)

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
		TimeoutSeconds: iwfidl.PtrInt32(2),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	respWait, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, respWait)

	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
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
