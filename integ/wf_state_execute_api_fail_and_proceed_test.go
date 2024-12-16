package integ

import (
	"context"
	"github.com/indeedeng/iwf/integ/workflow/wf_execute_api_fail_and_proceed"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

// TODO: Fix
func _TestStateExecuteApiFailAndProceedTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateExecuteApiFailAndProceed(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestStateExecuteApiFailAndProceed(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func TestStateExecuteApiFailAndProceedCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestStateExecuteApiFailAndProceed(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestStateExecuteApiFailAndProceed(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestStateExecuteApiFailAndProceed(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := wf_execute_api_fail_and_proceed.NewHandler()
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
	wfId := wf_execute_api_fail_and_proceed.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	stateOptions := &iwfidl.WorkflowStateOptions{
		ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
			MaximumAttempts: iwfidl.PtrInt32(1),
		},
		SkipWaitUntil:                   ptr.Any(true),
		ExecuteApiFailurePolicy:         iwfidl.PROCEED_TO_CONFIGURED_STATE.Ptr(),
		ExecuteApiFailureProceedStateId: ptr.Any(wf_execute_api_fail_and_proceed.StateRecover),
		ExecuteApiFailureProceedStateOptions: &iwfidl.WorkflowStateOptions{
			SkipWaitUntil: ptr.Any(true),
		},
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wf_execute_api_fail_and_proceed.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wf_execute_api_fail_and_proceed.State1),
		StateOptions:           stateOptions,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
		StateInput: &iwfidl.EncodedObject{
			Data:     ptr.Any(wf_execute_api_fail_and_proceed.InputData),
			Encoding: ptr.Any(wf_execute_api_fail_and_proceed.InputDataEncoding),
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
		"S1_decide":      1,
		"Recover_decide": 1,
	}, history, "wf state api fail and proceed test fail, %v", history)

	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.COMPLETED,
		Results: []iwfidl.StateCompletionOutput{
			{
				CompletedStateId:          wf_execute_api_fail_and_proceed.StateRecover,
				CompletedStateExecutionId: wf_execute_api_fail_and_proceed.StateRecover + "-1",
			},
		},
	}, resp, "response not expected")
}
