package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/wf_state_options_search_attributes_loading"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestWfStateOptionsSearchAttributesLoading_PARTIAL_WITHOUT_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestWfStateOptionsSearchAttributesLoading(t, backendType, iwfidl.PARTIAL_WITHOUT_LOCKING)
			smallWaitForFastTest()
			doTestWfStateOptionsSearchAttributesLoading(t, backendType, iwfidl.PARTIAL_WITHOUT_LOCKING)
			smallWaitForFastTest()
		}
	}
}

func TestWfStateOptionsSearchAttributesLoading_PARTIAL_WITH_LOCK(t *testing.T) {
	for _, backendType := range getBackendTypes() {
		for i := 0; i < *repeatIntegTest; i++ {
			doTestWfStateOptionsSearchAttributesLoading(t, backendType, iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK)
			smallWaitForFastTest()
			doTestWfStateOptionsSearchAttributesLoading(t, backendType, iwfidl.PARTIAL_WITHOUT_LOCKING)
			smallWaitForFastTest()
		}
	}
}

func doTestWfStateOptionsSearchAttributesLoading(
	t *testing.T, backendType service.BackendType, loadingType iwfidl.PersistenceLoadingType,
) {
	assertions := assert.New(t)

	wfHandler := wf_state_options_search_attributes_loading.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler, t)
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

	wfId := wf_state_options_search_attributes_loading.WorkflowType + "_" + string(loadingType) + "_" + strconv.Itoa(int(time.Now().UnixNano()))

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(string(loadingType)),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())

	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wf_state_options_search_attributes_loading.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wf_state_options_search_attributes_loading.State1),
		StateInput:             wfInput,
	}

	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()

	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
		"S3_start":  1,
		"S3_decide": 1,
		"S4_start":  1,
		"S4_decide": 1,
		"S5_start":  1,
		"S5_decide": 1,
	}, history, "state options search attributes loading, %v", history)
}
