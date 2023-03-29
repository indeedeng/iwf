package integ

import (
	"context"
	"github.com/indeedeng/iwf/integ/workflow"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/basic"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestBasicWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		//doTestBasicWorkflow(t, service.BackendTypeTemporal, iwfidl.WorkflowConfig{})
		smallWaitForFastTest()

		workflow.SetMockApiLatencyMs(3000)
		doTestBasicWorkflow(t, service.BackendTypeTemporal, iwfidl.WorkflowConfig{
			ContinueAsNewThresholdExecutedStateExecution: iwfidl.PtrInt32(1),
		})
		workflow.ClearMockApiLatencyMs()
		smallWaitForFastTest()
	}
}

func TestBasicWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestBasicWorkflow(t, service.BackendTypeCadence, iwfidl.WorkflowConfig{})
		smallWaitForFastTest()
	}
}

func doTestBasicWorkflow(t *testing.T, backendType service.BackendType, config iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := basic.NewHandler()
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
	wfId := basic.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test data"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        basic.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           basic.State1,
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			Config:                &config,
			WorkflowIDReusePolicy: ptr.Any(iwfidl.REJECT_DUPLICATE),
			// CronSchedule:          iwfidl.PtrString("* * * * *"),
			RetryPolicy: &iwfidl.WorkflowRetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(11),
				BackoffCoefficient:     iwfidl.PtrFloat32(11),
				MaximumAttempts:        iwfidl.PtrInt32(11),
				MaximumIntervalSeconds: iwfidl.PtrInt32(11),
			},
		},
		StateOptions: &iwfidl.WorkflowStateOptions{
			StartApiTimeoutSeconds:  iwfidl.PtrInt32(12),
			DecideApiTimeoutSeconds: iwfidl.PtrInt32(13),
			StartApiRetryPolicy: &iwfidl.RetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(12),
				BackoffCoefficient:     iwfidl.PtrFloat32(12),
				MaximumAttempts:        iwfidl.PtrInt32(12),
				MaximumIntervalSeconds: iwfidl.PtrInt32(12),
			},
			DecideApiRetryPolicy: &iwfidl.RetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(13),
				BackoffCoefficient:     iwfidl.PtrFloat32(13),
				MaximumAttempts:        iwfidl.PtrInt32(13),
				MaximumIntervalSeconds: iwfidl.PtrInt32(13),
			},
		},
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	panicAtHttpError(err, httpResp)

	// start it again should return already started error
	_, _, err = req.WorkflowStartRequest(startReq).Execute()
	apiErr, ok := err.(*iwfidl.GenericOpenAPIError)
	if !ok {
		log.Fatalf("Should fail to invoke start api %v", err)
	}
	errResp, ok := apiErr.Model().(iwfidl.ErrorResponse)
	if !ok {
		log.Fatalf("should be error response")
	}
	assertions := assert.New(t)
	assertions.Equal(errResp.GetSubStatus(), iwfidl.WORKFLOW_ALREADY_STARTED_SUB_STATUS)

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// use a wrong workflowId to test the error case
	_, _, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: "a wrong workflowId",
	}).Execute()
	apiErr, ok = err.(*iwfidl.GenericOpenAPIError)
	if !ok {
		log.Fatalf("Should fail to invoke get api %v", err)
	}
	errResp, ok = apiErr.Model().(iwfidl.ErrorResponse)
	if !ok {
		log.Fatalf("should be error response")
	}
	assertions.Equal(errResp.GetSubStatus(), iwfidl.WORKFLOW_NOT_EXISTS_SUB_STATUS)

	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "basic test fail, %v", history)

	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())
	assertions.Equal(1, len(resp2.GetResults()))
	assertions.Equal(iwfidl.StateCompletionOutput{
		CompletedStateId:          "S2",
		CompletedStateExecutionId: "S2-1",
		CompletedStateOutput:      wfInput,
	}, resp2.GetResults()[0])
}
