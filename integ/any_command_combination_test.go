package integ

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	anycommandconbination "github.com/indeedeng/iwf/integ/workflow/any_command_combination"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestAnyCommandCombinationWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCombinationWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func TestAnyCommandCombinationWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCloseWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func doTestAnyCommandCombinationWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := anycommandconbination.NewHandler()
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
	wfId := anycommandconbination.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        anycommandconbination.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           anycommandconbination.State1,
	}).Execute()
	panicAtHttpError(err, httpResp)

	signalValue := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-data-1"),
	}

	// send the signal
	req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId1,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// skip the timer
	time.Sleep(time.Second) // wait for a second so that timer is ready to be skipped
	req3 := apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandId:           iwfidl.PtrString(anycommandconbination.TimerId1),
	}).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second)
	// now it should be running at S2
	// TODO maybe we check it is already done S2??

	// send two signals for s2
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId1,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId2,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait and check the workflow, it should be still running
	time.Sleep(time.Second)
	reqDesc := apiClient.DefaultApi.ApiV1WorkflowGetPost(context.Background())
	_, httpResp, err = reqDesc.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// skip the other timer
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S2-1",
		TimerCommandId:           iwfidl.PtrString(anycommandconbination.TimerId1),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// workflow should be completed now
	descResp, httpResp, err := reqDesc.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions := assert.New(t)

	assertions.Equal(iwfidl.COMPLETED, descResp.GetWorkflowStatus())

	history, data := wfHandler.GetTestResult()

	assertions.Equalf(map[string]int64{
		"S1_start":  2,
		"S1_decide": 1,
		"S2_start":  2,
		"S2_decide": 1,
	}, history, "anycommandconbination test fail, %v", history)

	dataStr, _ := json.Marshal(data)
	fmt.Println(string(dataStr))
}
