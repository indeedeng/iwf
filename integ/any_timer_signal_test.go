package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	anytimersignal "github.com/indeedeng/iwf/integ/workflow/any_timer_signal"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestAnyTimerSignalWorkflowTemporal(t *testing.T) {
	doTestAnyTimerSignalWorkflow(t, service.BackendTypeTemporal)
}

// TODO this bug in Cadence SDK may cause the test to fail https://github.com/uber-go/cadence-client/issues/1198
func TestAnyTimerSignalWorkflowCadence(t *testing.T) {
	doTestAnyTimerSignalWorkflow(t, service.BackendTypeCadence)
}

func doTestAnyTimerSignalWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := anytimersignal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	// create client
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	// start a workflow
	wfId := anytimersignal.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        anytimersignal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           anytimersignal.State1,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for 3 secs and send the signal
	time.Sleep(time.Second * 3)
	signalValue := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-data-1"),
	}
	req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anytimersignal.SignalName,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":  2,
		"S1_decide": 2,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "anytimersignal test fail, %v", history)

	assertions.Equal(anytimersignal.SignalName, data["signalChannelName1"])
	assertions.Equal("signal-cmd-id", data["signalCommandId1"])
	assertions.Equal(service.SignalStatusWaiting, data["signalStatus1"])

	assertions.Equal(anytimersignal.SignalName, data["signalChannelName2"])
	assertions.Equal("signal-cmd-id", data["signalCommandId2"])
	assertions.Equal(service.SignalStatusReceived, data["signalStatus2"])
	assertions.Equal(signalValue, data["signalValue2"])
}
