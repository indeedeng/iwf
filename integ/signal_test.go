package integ

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/workflow/signal"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSignalWorkflowTemporal(t *testing.T) {
	doTestSignalWorkflow(t, service.BackendTypeTemporal)
}

func TestSignalWorkflowCadence(t *testing.T) {
	doTestSignalWorkflow(t, service.BackendTypeCadence)
}

func doTestSignalWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := signal.NewHandler()
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
	wfId := signal.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           signal.State1,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}

	// signal the workflow
	signalVal := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-data"),
	}
	//err = temporalClient.SignalWorkflow(context.Background(), wfId, "", signal.SignalName, signalVal)
	req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp2, err := req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:  wfId,
		SignalName:  signal.SignalName,
		SignalValue: &signalVal,
	}).Execute()

	if err != nil || httpResp2.StatusCode != 200 {
		log.Fatalf("Fail to signal the workflow %v %v", err, httpResp2)
	}

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithLongWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId:   wfId,
		NeedsResults: iwfidl.PtrBool(true),
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
	}, history, "signal test fail, %v", history)
	assertions.Equal("signal-cmd-id", data["signalId"])
	assertions.Equal(signalVal, data["signalValue"])
}
