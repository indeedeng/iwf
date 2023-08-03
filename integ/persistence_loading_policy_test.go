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
			doTestPersistenceLoadingPolicy(t, backendType, persistence_loading_policy.LoadingTypeAll)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicy_PARTIAL_WITHOUT_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, persistence_loading_policy.LoadingTypePartialWithoutLock)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicy_PARTIAL_WITH_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, persistence_loading_policy.LoadingTypePartialWithLock)
			smallWaitForFastTest()
		}
	}
}

func TestPersistenceLoadingPolicyTemporal_NONE(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestPersistenceLoadingPolicy(t, backendType, persistence_loading_policy.LoadingTypeNone)
			smallWaitForFastTest()
		}
	}
}

func doTestPersistenceLoadingPolicy(t *testing.T, backendType service.BackendType, loadingType string) {
	wfHandler := persistence_loading_policy.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
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

	wfId := persistence_loading_policy.WorkflowType + "_" + loadingType + "_" + strconv.Itoa(int(time.Now().UnixNano()))

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(loadingType),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())

	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        persistence_loading_policy.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(persistence_loading_policy.State1),
		StateInput:             wfInput,
	}

	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 1)

	history, data := wfHandler.GetTestResult()

	assertions := assert.New(t)

	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "persistence loading policy test fail, %v", history)

	if loadingType == persistence_loading_policy.LoadingTypeAll {
		assertions.Equalf(map[string]interface{}{
			"sa_state_start_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_start_S2_" + persistence.TestSearchAttributeTextKey:    "test-search-attribute-2",
			"sa_state_start_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},
			"sa_state_start_S2_da_2": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value2"),
			},
			"sa_state_decide_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_decide_S2_" + persistence.TestSearchAttributeTextKey:    "test-search-attribute-2",
			"sa_state_decide_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},
			"sa_state_decide_S2_da_2": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value2"),
			},
		}, data, "persistence loading policy test fail, %v", data)
	} else if loadingType == persistence_loading_policy.LoadingTypePartialWithoutLock {
		assertions.Equalf(map[string]interface{}{
			"sa_state_start_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_start_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},

			"sa_state_decide_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_decide_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},
		}, data, "persistence loading policy test fail, %v", data)
	} else if loadingType == persistence_loading_policy.LoadingTypePartialWithLock {
		assertions.Equalf(map[string]interface{}{
			"sa_state_start_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_start_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},

			"sa_state_decide_S2_" + persistence.TestSearchAttributeKeywordKey: "test-search-attribute-1",
			"sa_state_decide_S2_da_1": iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString("test-data-object-value1"),
			},
		}, data, "persistence loading policy test fail, %v", data)
	} else if loadingType == persistence_loading_policy.LoadingTypeNone {
		assertions.Equalf(map[string]interface{}{}, data, "persistence loading policy test fail, %v", data)
	}
}
