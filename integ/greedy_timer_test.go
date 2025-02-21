package integ

import (
	"context"
	"encoding/json"
	"github.com/indeedeng/iwf/integ/workflow/greedy_timer"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
)

func TestGreedyTimerWorkflowBaseTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestGreedyTimerWorkflow(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestGreedyTimerWorkflowBaseCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestGreedyTimerWorkflow(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

func TestGreedyTimerWorkflowBaseTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestGreedyTimerWorkflowCustomConfig(t, service.BackendTypeTemporal, greedyTimerConfig(true))
		smallWaitForFastTest()
	}
}

func TestGreedyTimerWorkflowBaseCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestGreedyTimerWorkflowCustomConfig(t, service.BackendTypeCadence, greedyTimerConfig(true))
		smallWaitForFastTest()
	}
}

func doTestGreedyTimerWorkflow(t *testing.T, backendType service.BackendType) {
	doTestGreedyTimerWorkflowCustomConfig(t, backendType, minimumGreedyTimerConfig())
}

func doTestGreedyTimerWorkflowCustomConfig(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := greedy_timer.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler, t)
	defer closeFunc1()

	uClient, closeFunc2 := startIwfServiceWithClient(backendType)
	defer closeFunc2()

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	// start a workflow
	durations := []int64{15, 30}
	input := greedy_timer.Input{Durations: durations}

	wfId := greedy_timer.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	inputData, _ := json.Marshal(input)

	//schedule-1
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        greedy_timer.WorkflowType,
		WorkflowTimeoutSeconds: 30,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(greedy_timer.ScheduleTimerState),
		StateInput: &iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(string(inputData)),
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Short wait for workflow to initialize
	time.Sleep(time.Second * 1)

	// assertions
	debug := service.DebugDumpResponse{}
	err = uClient.QueryWorkflow(context.Background(), &debug, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	assertions.Equal(1, len(debug.FiringTimersUnixTimestamps))
	singleTimerScheduled := debug.FiringTimersUnixTimestamps[0]

	scheduleTimerAndAssertExpectedScheduled(t, apiClient, uClient, wfId, 20, 1)

	// skip next timer for state: schedule-1
	skipReq := apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = skipReq.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: greedy_timer.ScheduleTimerState + "-1",
		TimerCommandId:           iwfidl.PtrString("duration-15"),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Short wait for workflow to initialize
	time.Sleep(time.Second * 1)

	err = uClient.QueryWorkflow(context.Background(), &debug, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}

	// no second timer started
	assertions.Equal(1, len(debug.FiringTimersUnixTimestamps))
	// LessOrEqual due to continue as new workflow scheduling the next, not skipped timer
	assertions.LessOrEqual(singleTimerScheduled, debug.FiringTimersUnixTimestamps[0])
	scheduleTimerAndAssertExpectedScheduled(t, apiClient, uClient, wfId, 5, 2)

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"schedule_start":  3,
		"schedule_decide": 1,
	}, history, "history does not match expected")
}

func scheduleTimerAndAssertExpectedScheduled(
	t *testing.T,
	apiClient *iwfidl.APIClient,
	uClient uclient.UnifiedClient,
	wfId string,
	duration int64,
	noMoreThan int) {

	assertions := assert.New(t)
	input := greedy_timer.Input{Durations: []int64{duration}}
	inputData, _ := json.Marshal(input)

	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	_, httpResp, err := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    greedy_timer.SubmitDurationsRPC,
		Input: &iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(string(inputData)),
		},
		TimeoutSeconds: iwfidl.PtrInt32(2),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Short wait for workflow to initialize
	time.Sleep(time.Second * 1)

	debug := service.DebugDumpResponse{}
	err = uClient.QueryWorkflow(context.Background(), &debug, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}

	assertions.LessOrEqual(len(debug.FiringTimersUnixTimestamps), noMoreThan)
}
