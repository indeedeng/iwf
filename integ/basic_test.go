package integ

import (
	"context"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
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
		doTestBasicWorkflow(t, service.BackendTypeTemporal)
		// NOTE: basic wf is too fast so we have to make sure to have enough interval
		du := time.Millisecond * time.Duration(*repeatInterval)
		if *repeatIntegTest > 1 && du < time.Second {
			du = time.Second
		}
		time.Sleep(du)
	}
}

func TestBasicWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestBasicWorkflow(t, service.BackendTypeCadence)
		// NOTE: basic wf is too fast so we have to make sure to have enough interval
		du := time.Millisecond * time.Duration(*repeatInterval)
		if *repeatIntegTest > 1 && du < time.Second {
			du = time.Second
		}
		time.Sleep(du)
	}
}

func doTestBasicWorkflow(t *testing.T, backendType service.BackendType) {
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
	wfId := basic.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test data"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        basic.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           basic.State1,
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowIDReusePolicy: ptr.Any(iwfidl.REJECT_DUPLICATE),
			// CronSchedule:          iwfidl.PtrString("* * * * *"),
			RetryPolicy: &iwfidl.RetryPolicy{
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
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke get api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	history, _ := wfHandler.GetTestResult()
	assertions := assert.New(t)
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
