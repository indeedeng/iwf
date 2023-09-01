package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/integ/workflow/persistence_loading_policy"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestPersistenceLoadingPolicy_ALL(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.ALL_WITHOUT_LOCKING, true)
			smallWaitForFastTest()
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.ALL_WITHOUT_LOCKING, false)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicy_PARTIAL_WITHOUT_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.PARTIAL_WITHOUT_LOCKING, true)
			smallWaitForFastTest()
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.PARTIAL_WITHOUT_LOCKING, false)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicy_PARTIAL_WITH_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK, true)
			smallWaitForFastTest()
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK, false)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicy_NONE(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.NONE, true)
			smallWaitForFastTest()
			doTestPersistenceLoadingPolicy(t, backendType, iwfidl.NONE, false)
			smallWaitForFastTest()
		}
	}
}

func doTestPersistenceLoadingPolicy(t *testing.T, backendType service.BackendType, loadingType iwfidl.PersistenceLoadingType, rpcUseMemo bool) {
	assertions := assert.New(t)

	wfHandler := persistence_loading_policy.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
	defer closeFunc1()
	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	wfId := persistence_loading_policy.WorkflowType + "_" + string(loadingType) + "_" + strconv.Itoa(int(time.Now().UnixNano()))

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(string(loadingType)),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())

	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        persistence_loading_policy.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(persistence_loading_policy.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			UseMemoForDataAttributes: ptr.Any(rpcUseMemo),
		},
	}

	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	if rpcUseMemo && backendType == service.BackendTypeCadence {
		if err == nil {
			panic("err should not be nil when Memo is not supported with Cadence")
		}
		return
	}
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 2)

	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())

	rpcReq := iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    persistence_loading_policy.WorkflowType + "_rpc",
		Input:      wfInput,
		SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: &loadingType,
			PartialLoadingKeys: []string{
				persistence.TestSearchAttributeKeywordKey,
			},
			LockingKeys: []string{
				persistence.TestSearchAttributeTextKey,
			},
		},
		DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: &loadingType,
			PartialLoadingKeys: []string{
				"da_1",
			},
			LockingKeys: []string{
				"da_2",
			},
		},
		TimeoutSeconds:           iwfidl.PtrInt32(3),
		UseMemoForDataAttributes: ptr.Any(rpcUseMemo),
		SearchAttributes:         getSearchAttributesToGetFromMemo(loadingType),
	}

	_, httpResp, err = reqRpc.WorkflowRpcRequest(rpcReq).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 2)

	history, _ := wfHandler.GetTestResult()

	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"rpc":       1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "persistence loading policy test fail, %v", history)
}

func getSearchAttributesToGetFromMemo(loadingType iwfidl.PersistenceLoadingType) []iwfidl.SearchAttributeKeyAndType {
	if loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING {
		return []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType: iwfidl.KEYWORD.Ptr(),
			},
		}
	}

	return []iwfidl.SearchAttributeKeyAndType{
		{
			Key:       iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
			ValueType: iwfidl.KEYWORD.Ptr(),
		},
		{
			Key:       iwfidl.PtrString(persistence.TestSearchAttributeTextKey),
			ValueType: iwfidl.TEXT.Ptr(),
		},
	}
}
