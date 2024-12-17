package integ

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/indeedeng/iwf/integ/workflow/locking"
	"github.com/indeedeng/iwf/service/common/ptr"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestLockingWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestLockingWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestLockingWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		// TODO not sure why using minimumContinueAsNewConfig(true) will fail
		doTestLockingWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func TestLockingWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestLockingWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestLockingWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestLockingWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func doTestLockingWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := locking.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:                      backendType,
		DisableFailAtMemoIncompatibility: true,
	})
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := locking.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        locking.WorkflowType,
		WorkflowTimeoutSeconds: 300,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(locking.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	panicAtHttpError(err, httpResp)

	for i := 0; i < locking.NumUnusedSignals; i++ {
		// send 4 unused signals at the beginning to validate the ChannelInfo feature
		reqSignal := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		httpResp, err = reqSignal.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: locking.UnusedSignalChannelName,
			SignalValue:       nil,
		}).Execute()
		panicAtHttpError(err, httpResp)
	}

	assertions := assert.New(t)

	if config != nil && backendType == service.BackendTypeTemporal {
		// special waiting time for continue as new
		// the first run will have to take more time to finish all the in parallel waitUntil APIs before continueAsNew
		time.Sleep(locking.InParallelS2 * time.Second)
	}
	rpcIncrease := 0
	rpcLockingFailure := 0
	if backendType == service.BackendTypeTemporal {
		// only test rpc locking with Temporal
		for i := 0; i < 25; i++ {
			allSearchAttributes := []iwfidl.SearchAttributeKeyAndType{
				{
					Key:       iwfidl.PtrString(locking.TestSearchAttributeKeywordKey),
					ValueType: iwfidl.KEYWORD.Ptr(),
				},
				{
					Key:       iwfidl.PtrString(locking.TestSearchAttributeIntKey),
					ValueType: iwfidl.INT.Ptr(),
				},
			}
			time.Sleep(time.Second * 2)
			reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
			rpcResp, httpResp, err := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
				WorkflowId: wfId,
				RpcName:    locking.RPCName,
				Input:      locking.TestValue,
				SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
					PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
					PartialLoadingKeys: []string{
						locking.TestSearchAttributeKeywordKey,
						locking.TestSearchAttributeIntKey,
					},
					LockingKeys: []string{
						locking.TestSearchAttributeIntKey,
					},
				},
				DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
					PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
					PartialLoadingKeys: []string{
						locking.TestDataObjectKey2,
						locking.TestDataObjectKey1,
					},
					LockingKeys: []string{
						locking.TestDataObjectKey1,
					},
				},
				TimeoutSeconds:   iwfidl.PtrInt32(2),
				SearchAttributes: allSearchAttributes,
			}).Execute()
			if err != nil || httpResp.StatusCode != 200 {
				if httpResp.StatusCode == service.HttpStatusCodeSpecial4xxError2 {
					var errResp iwfidl.ErrorResponse
					body, err := ioutil.ReadAll(httpResp.Body)
					assertions.Nil(err)
					err = json.Unmarshal(body, &errResp)
					lockingErrorMsg := "requested data or search attributes are being locked by other operations"
					assertions.Equal(lockingErrorMsg, errResp.GetDetail())
					assertions.Equal(iwfidl.WORKER_API_ERROR, errResp.GetSubStatus())
					fmt.Println(lockingErrorMsg)
					rpcLockingFailure++
					continue
				} else {
					panicAtHttpError(err, httpResp)
				}
			}
			fmt.Println("rpc execution succeeded")
			rpcIncrease++
			assertions.Equal(rpcResp.Output, locking.TestValue)
		}
		assertions.True(rpcIncrease > 0)
		assertions.True(rpcLockingFailure > 0)
		fmt.Println("rpc results, success, failure:", rpcIncrease, rpcLockingFailure)
	}

	time.Sleep(time.Second * 1)
	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    locking.RPCName,
		Input:      locking.UnblockValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 20)
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	s2StartsDecides := locking.InParallelS2 + rpcIncrease // locking.InParallelS2 original state executions, and a new trigger from rpc
	finalCounterValue := int64(locking.InParallelS2 + 2*rpcIncrease)
	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":            1,
		"S1_decide":           1,
		"StateWaiting_start":  1,
		"StateWaiting_decide": 1,
		"S2_start":            int64(s2StartsDecides),
		"S2_decide":           int64(s2StartsDecides),
	}, history, "locking.test fail, %v", history)

	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())
	// State completions with empty output are ignored
	assertions.Equal(0, len(resp2.GetResults()))

	reqSearch := apiClient.DefaultApi.ApiV1WorkflowSearchattributesGetPost(context.Background())
	searchResult2, httpResp, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId: wfId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(locking.TestSearchAttributeIntKey),
				ValueType: ptr.Any(iwfidl.INT),
			},
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	expectedSearchAttributeInt := iwfidl.SearchAttribute{
		Key:          iwfidl.PtrString(locking.TestSearchAttributeIntKey),
		ValueType:    ptr.Any(iwfidl.INT),
		IntegerValue: iwfidl.PtrInt64(finalCounterValue),
	}
	assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeInt}, searchResult2.GetSearchAttributes())

	reqQry := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	// force to test compatibility of memo
	useMemo := false
	if backendType == service.BackendTypeTemporal {
		useMemo = true
	}
	queryResult1, httpResp, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			locking.TestDataObjectKey1,
		},
		UseMemoForDataAttributes: ptr.Any(useMemo),
	}).Execute()
	panicAtHttpError(err, httpResp)

	expected1 := []iwfidl.KeyValue{
		{
			Key: iwfidl.PtrString(locking.TestDataObjectKey1),
			Value: &iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString(fmt.Sprintf("%v", finalCounterValue)),
			},
		},
	}
	assertions.ElementsMatch(expected1, queryResult1.GetObjects())

	//reset here with reapply and compare counter
	resetReq := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
	_, httpResp, err = resetReq.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
		WorkflowId: wfId,
		ResetType:  iwfidl.BEGINNING,
		//SkipSignalReapply: ptr.Any(true),
	}).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 20)
	req2Reset := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2Reset, httpResp, err := req2Reset.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions.Equal(iwfidl.COMPLETED, resp2Reset.GetWorkflowStatus())

	//TODO: There is a bug in the Temporal go SDK where only the first update method is actually executed. When that is fixed the following code can be uncommented to test resetting update methods.
	//time.Sleep(time.Second * 10)
	//reqRpcReset := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	//_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
	//	WorkflowId: wfId,
	//	RpcName:    locking.RPCName,
	//	Input:      locking.UnblockValue,
	//}).Execute()
	//panicAtHttpError(err, httpResp)

	//time.Sleep(time.Second * 20)
	//req2Reset := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	//resp2Reset, httpResp, err := req2Reset.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
	//	WorkflowId: wfId,
	//}).Execute()
	//panicAtHttpError(err, httpResp)

	//s2StartsDecides := locking.InParallelS2 + rpcIncrease // locking.InParallelS2 original state executions, and a new trigger from rpc
	//finalCounterValue := int64(locking.InParallelS2 + 2*rpcIncrease)
	//stateCompletionCount := locking.InParallelS2 + rpcIncrease + 1
	//resetHistory, _ := wfHandler.GetTestResult()
	//assertions.Equalf(map[string]int64{
	//	"S1_start":            1,
	//	"S1_decide":           1,
	//	"StateWaiting_start":  1,
	//	"StateWaiting_decide": 1,
	//	"S2_start":            int64(s2StartsDecides),
	//	"S2_decide":           int64(s2StartsDecides),
	//}, resetHistory, "locking.test fail, %v", history)
	//
	//assertions.Equal(iwfidl.COMPLETED, resp2Reset.GetWorkflowStatus())
	//assertions.Equal(stateCompletionCount, len(resp2Reset.GetResults()))
	//
	//reqSearchReset := apiClient.DefaultApi.ApiV1WorkflowSearchattributesGetPost(context.Background())
	//searchResultReset, httpResp, err := reqSearchReset.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
	//	WorkflowId: wfId,
	//	Keys: []iwfidl.SearchAttributeKeyAndType{
	//		{
	//			Key:       iwfidl.PtrString(locking.TestSearchAttributeIntKey),
	//			ValueType: ptr.Any(iwfidl.INT),
	//		},
	//	},
	//}).Execute()
	//panicAtHttpError(err, httpResp)
	//
	//expectedSearchAttributeIntReset := iwfidl.SearchAttribute{
	//	Key:          iwfidl.PtrString(locking.TestSearchAttributeIntKey),
	//	ValueType:    ptr.Any(iwfidl.INT),
	//	IntegerValue: iwfidl.PtrInt64(finalCounterValue),
	//}
	//assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeIntReset}, searchResultReset.GetSearchAttributes())

	//reset here without update reapply and counter should be less
}
