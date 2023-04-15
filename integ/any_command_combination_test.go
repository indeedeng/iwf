package integ

import (
	"context"
	"encoding/json"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	anycommandconbination "github.com/indeedeng/iwf/integ/workflow/any_command_combination"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestAnyCommandCombinationWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCombinationWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestAnyCommandCombinationWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestAnyCommandCombinationWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestAnyCommandCloseWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestAnyCommandCloseWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func doTestAnyCommandCombinationWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
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
		StartStateId:           ptr.Any(anycommandconbination.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	signalValue := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-data-1"),
	}

	// send the signals to S1
	req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId1,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId1,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// skip the timer for S1
	time.Sleep(time.Second * 2) // wait for a second so that timer is ready to be skipped
	req3 := apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandId:           iwfidl.PtrString(anycommandconbination.TimerId1),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// now it should be running at S2
	// Future: we can check it is already done S1

	// send first signal for s2
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId1,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait and check the workflow, it should be still running
	time.Sleep(time.Second)
	reqDesc := apiClient.DefaultApi.ApiV1WorkflowGetPost(context.Background())
	descResp, httpResp, err := reqDesc.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)
	assertions.Equal(iwfidl.RUNNING, descResp.GetWorkflowStatus())

	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId3,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// send 2nd signal for s2
	httpResp, err = req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: anycommandconbination.SignalNameAndId2,
		SignalValue:       &signalValue,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// workflow should be completed now
	if config == nil {
		time.Sleep(time.Second)
		descResp, httpResp, err = reqDesc.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		panicAtHttpError(err, httpResp)
		assertions.Equal(iwfidl.COMPLETED, descResp.GetWorkflowStatus())
	} else {
		reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
		respWait, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, respWait)
	}

	history, data := wfHandler.GetTestResult()

	assertions.Equalf(map[string]int64{
		"S1_start":  2,
		"S1_decide": 1,
		"S2_start":  2,
		"S2_decide": 1,
	}, history, "anycommandconbination test fail, %v", history)

	var s1CommandResults iwfidl.CommandResults
	var s2CommandResults iwfidl.CommandResults
	s1ResultJsonStr := "{\"signalResults\":[" +
		"{\"commandId\":\"test-signal-name1\",\"signalChannelName\":\"test-signal-name1\",\"signalRequestStatus\":\"RECEIVED\",\"signalValue\":{\"data\":\"test-data-1\",\"encoding\":\"json\"}}, " +
		"{\"commandId\":\"test-signal-name1\",\"signalChannelName\":\"test-signal-name1\",\"signalRequestStatus\":\"RECEIVED\",\"signalValue\":{\"data\":\"test-data-1\",\"encoding\":\"json\"}}, " +
		"{\"commandId\":\"test-signal-name2\",\"signalChannelName\":\"test-signal-name2\",\"signalRequestStatus\":\"WAITING\"}," +
		"{\"commandId\":\"test-signal-name3\",\"signalChannelName\":\"test-signal-name3\",\"signalRequestStatus\":\"WAITING\"}" +
		"],\"timerResults\":[" +
		"{\"commandId\":\"test-timer-1\",\"timerStatus\":\"FIRED\"}]," +
		"\"stateStartApiSucceeded\":true}"
	err = json.Unmarshal([]byte(s1ResultJsonStr), &s1CommandResults)
	if err != nil {
		panic(err)
	}
	s2ResultsJsonStr := "{\"signalResults\":[" +
		"{\"commandId\":\"test-signal-name1\",\"signalChannelName\":\"test-signal-name1\",\"signalRequestStatus\":\"RECEIVED\",\"signalValue\":{\"data\":\"test-data-1\",\"encoding\":\"json\"}}, " +
		"{\"commandId\":\"test-signal-name1\",\"signalChannelName\":\"test-signal-name1\",\"signalRequestStatus\":\"WAITING\"}," +
		"{\"commandId\":\"test-signal-name2\",\"signalChannelName\":\"test-signal-name2\",\"signalRequestStatus\":\"RECEIVED\",\"signalValue\":{\"data\":\"test-data-1\",\"encoding\":\"json\"}}," +
		"{\"commandId\":\"test-signal-name3\",\"signalChannelName\":\"test-signal-name3\",\"signalRequestStatus\":\"RECEIVED\",\"signalValue\":{\"data\":\"test-data-1\",\"encoding\":\"json\"}}" +
		"],\"timerResults\":[" +
		"{\"commandId\":\"test-timer-1\",\"timerStatus\":\"SCHEDULED\"}]," +
		"\"stateStartApiSucceeded\":true}"
	err = json.Unmarshal([]byte(s2ResultsJsonStr), &s2CommandResults)
	if err != nil {
		panic(err)
	}
	expectedData := map[string]interface{}{
		"s1_commandResults": s1CommandResults,
		"s2_commandResults": s2CommandResults,
	}
	assertions.Equal(expectedData, data)
}
