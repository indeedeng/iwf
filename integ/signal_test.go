package integ

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestSignalWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func TestSignalWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
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
	wfId := signal.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           signal.State1,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// signal the workflow
	var signalVals []iwfidl.EncodedObject
	for i := 0; i < 4; i++ {
		signalVal := iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(fmt.Sprintf("test-data-%v", i)),
		}
		signalVals = append(signalVals, signalVal)

		req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		httpResp2, err := req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: signal.SignalName,
			SignalValue:       &signalVal,
		}).Execute()

		panicAtHttpError(err, httpResp2)
	}

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "signal test fail, %v", history)

	assertions.Equal(fmt.Sprintf("signal-cmd-id%v", 0), data[fmt.Sprintf("signalId%v", 0)])
	assertions.Equal(fmt.Sprintf("signal-cmd-id%v", 1), data[fmt.Sprintf("signalId%v", 1)])
	assertions.Equal("", data[fmt.Sprintf("signalId%v", 2)])
	assertions.Equal("", data[fmt.Sprintf("signalId%v", 3)])
	for i := 0; i < 4; i++ {
		assertions.Equal(signalVals[i], data[fmt.Sprintf("signalValue%v", i)])
	}
}
