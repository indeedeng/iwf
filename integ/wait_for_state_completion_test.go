package integ

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/wait_for_state_completion"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestWaitForStateCompletionTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitForStateCompletion(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestWaitForStateCompletionCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitForStateCompletion(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func doTestWaitForStateCompletion(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := wait_for_state_completion.NewHandler()
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
	wfId := wait_for_state_completion.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	nowTimestamp := time.Now().Unix()
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wait_for_state_completion.WorkflowType,
		WorkflowTimeoutSeconds: 30,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wait_for_state_completion.State1),
		StateInput: &iwfidl.EncodedObject{
			Data: iwfidl.PtrString(strconv.Itoa(int(nowTimestamp))),
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
		WaitForCompletionStateExecutionIds: []string{"S1-1"},
	}).Execute()
	panicAtHttpError(err, httpResp)

	req1 := apiClient.DefaultApi.ApiV1WorkflowWaitForStateCompletionPost(context.Background())
	_, httpResp, err = req1.WorkflowWaitForStateCompletionRequest(
		iwfidl.WorkflowWaitForStateCompletionRequest{
			WorkflowId:       wfId,
			StateExecutionId: "S1-1",
			WaitTimeSeconds:  iwfidl.PtrInt32(30),
		}).Execute()
	panicAtHttpError(err, httpResp)

	assertions := assert.New(t)
	assertions.Equal(200, httpResp.StatusCode)
	// read httpResp body
	var output iwfidl.WorkflowWaitForStateCompletionResponse
	defer httpResp.Body.Close()
	err = json.NewDecoder(httpResp.Body).Decode(&output)
	if err != nil {
		log.Fatalf("Failed to decode the response: %v", err)
	}
	assertions.Equal("S1-1", output.StateCompletionOutput.CompletedStateExecutionId)

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, data := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "timer test fail, %v", history)
	duration := (data["fired_at"]).(int64) - (data["scheduled_at"]).(int64)
	assertions.Equal("timer-cmd-id", data["timer_id"])
	assertions.True(duration >= 9 && duration <= 11, duration)
}
