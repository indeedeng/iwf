package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/timer"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestTimerWorkflowTemporal(t *testing.T) {
	doTestTimerWorkflow(t, service.BackendTypeTemporal)
}

func TestTimerWorkflowCadence(t *testing.T) {
	doTestTimerWorkflow(t, service.BackendTypeCadence)
}

func doTestTimerWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := timer.NewHandler()
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
	wfId := timer.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        timer.WorkflowType,
		WorkflowTimeoutSeconds: 30,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           timer.State1,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "timer test fail, %v", history)
	duration := (data["fired_at"]).(int64) - (data["scheduled_at"]).(int64)
	assertions.Equal("timer-cmd-id", data["timer_id"])
	assertions.True(duration >= 9 && duration <= 11)
}
