package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	anycommandclose "github.com/indeedeng/iwf/integ/workflow/any_command_close"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestAnyCommandCloseWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCloseWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Second * time.Duration(*repeatInterval))
	}
}

func TestAnyCommandCloseWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCloseWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Second * time.Duration(*repeatInterval))
	}
}

func doTestAnyCommandCloseWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := anycommandclose.NewHandler()
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
	wfId := anycommandclose.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        anycommandclose.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           anycommandclose.State1,
	}).Execute()
	panicAtHttpError(err, httpResp)

	signalValue := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-data-1"),
	}

	req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandclose.SignalName2,
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
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "anycommandclose test fail, %v", history)

	assertions.Equal(anycommandclose.SignalName2, data["signalChannelName1"])
	assertions.Equal("signal-cmd-id2", data["signalCommandId1"])
	assertions.Equal(signalValue, data["signalValue1"])
	assertions.Equal(iwfidl.RECEIVED, data["signalStatus1"])

	assertions.Equal(anycommandclose.SignalName1, data["signalChannelName0"])
	assertions.Equal("signal-cmd-id1", data["signalCommandId0"])
	assertions.Equal(iwfidl.WAITING, data["signalStatus0"])
}
