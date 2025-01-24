package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/wf_state_api_timeout"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestStateApiTimeoutTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateApiTimeout(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestStateApiTimeout(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func TestStateApiTimeoutCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateApiTimeout(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestStateApiTimeout(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestStateApiTimeout(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := wf_state_api_timeout.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := wf_state_api_timeout.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wf_state_api_timeout.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wf_state_api_timeout.State1),
		StateOptions: &iwfidl.WorkflowStateOptions{
			WaitUntilApiTimeoutSeconds: iwfidl.PtrInt32(1),
			WaitUntilApiRetryPolicy: &iwfidl.RetryPolicy{
				MaximumAttempts: iwfidl.PtrInt32(1),
			},
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp, t)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start": 1,
	}, history, "wf state api fail test fail, %v", history)

	assertions.True(strings.Contains(resp.GetErrorMessage(), "StartToClose") || strings.Contains(resp.GetErrorMessage(), "START_TO_CLOSE"), resp.GetErrorMessage())
	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.FAILED,
		ErrorType:      ptr.Any(iwfidl.STATE_API_FAIL_MAX_OUT_RETRY_ERROR_TYPE),
		ErrorMessage:   resp.ErrorMessage,
	}, resp, "response not expected")
}
