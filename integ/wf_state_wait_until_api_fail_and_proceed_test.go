package integ

import (
	"context"
	"github.com/indeedeng/iwf/integ/workflow/wf_state_api_fail_and_proceed"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestStateApiFailAndProceedTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateApiFailAndProceed(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestStateApiFailAndProceed(t, service.BackendTypeTemporal, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestStateApiFailAndProceedCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateApiFailAndProceed(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestStateApiFailAndProceed(t, service.BackendTypeCadence, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func doTestStateApiFailAndProceed(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := wf_state_api_fail_and_proceed.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
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
	wfId := wf_state_api_fail_and_proceed.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	stateOptions := &iwfidl.WorkflowStateOptions{
		StartApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttempts: iwfidl.PtrInt32(1),
		},
		StartApiFailurePolicy: iwfidl.PROCEED_TO_DECIDE_ON_START_API_FAILURE.Ptr(),
	}
	if config != nil { // use this hack to test the compatibility
		stateOptions = &iwfidl.WorkflowStateOptions{
			StartApiRetryPolicy: &iwfidl.RetryPolicy{
				MaximumAttempts: iwfidl.PtrInt32(1),
			},
			WaitUntilApiFailurePolicy: iwfidl.PROCEED_ON_FAILURE.Ptr(),
		}
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wf_state_api_fail_and_proceed.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wf_state_api_fail_and_proceed.State1),
		StateOptions:           stateOptions,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, _ := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
	}, history, "wf state api fail and proceed test fail, %v", history)

	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.COMPLETED,
		Results: []iwfidl.StateCompletionOutput{
			{
				CompletedStateId:          wf_state_api_fail_and_proceed.State1,
				CompletedStateExecutionId: wf_state_api_fail_and_proceed.State1 + "-1",
			},
		},
	}, resp, "response not expected")
}
