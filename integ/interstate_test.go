package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/interstate"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestInterStateWorkflowTemporal(t *testing.T) {
	doTestInterStateWorkflow(t, service.BackendTypeTemporal)
}

func TestInterStateWorkflowCadence(t *testing.T) {
	doTestInterStateWorkflow(t, service.BackendTypeCadence)
}

func doTestInterStateWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := interstate.NewHandler()
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
	wfId := interstate.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        interstate.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           interstate.State1,
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
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":   1,
		"S1_decide":  1,
		"S21_start":  1,
		"S21_decide": 1,
		"S22_start":  1,
		"S22_decide": 1,
		"S31_start":  1,
		"S31_decide": 1,
	}, history, "interstate test fail, %v", history)

	assertions.Equal(service.WorkflowStatusCompleted, resp2.GetWorkflowStatus())
	assertions.Equal(0, len(resp2.GetResults()))
	assertions.Equal(map[string]interface{}{
		interstate.State21 + "received": interstate.TestVal1,
		interstate.State31 + "received": interstate.TestVal2,
	}, data)
}
