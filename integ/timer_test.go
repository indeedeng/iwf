package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/timer"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"log"
	"strconv"
	"testing"
	"time"
)

func TestTimerWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func TestTimerWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func doTestTimerWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := timer.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
	defer closeFunc1()

	uclient, closeFunc2 := doStartIwfServiceWithClient(backendType)
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := timer.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	nowTimestamp := time.Now().Unix()
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        timer.WorkflowType,
		WorkflowTimeoutSeconds: 30,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           timer.State1,
		StateInput: &iwfidl.EncodedObject{
			Data: iwfidl.PtrString(strconv.Itoa(int(nowTimestamp))),
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			Config: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 3)
	timerInfos := service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	assertions := assert.New(t)
	timer2 := &service.TimerInfo{
		CommandId:                  "timer-cmd-id-2",
		FiringUnixTimestampSeconds: nowTimestamp + 86400,
		Status:                     service.TimerPending,
	}
	timer3 := &service.TimerInfo{
		CommandId:                  "timer-cmd-id-3",
		FiringUnixTimestampSeconds: nowTimestamp + 86400*365,
		Status:                     service.TimerPending,
	}
	expectedTimerInfos := service.GetCurrentTimerInfosQueryResponse{
		StateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{
			"S1-1": {
				{
					CommandId:                  "timer-cmd-id",
					FiringUnixTimestampSeconds: nowTimestamp + 10,
					Status:                     service.TimerPending,
				},
				timer2,
				timer3,
			},
		},
	}
	assertions.Equal(expectedTimerInfos, timerInfos)

	req3 := apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandId:           iwfidl.PtrString("timer-cmd-id-2"),
	}).Execute()
	panicAtHttpError(err, httpResp)

	timerInfos = service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer2.Status = service.TimerSkipped
	assertions.Equal(expectedTimerInfos, timerInfos)

	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandIndex:        iwfidl.PtrInt32(2),
	}).Execute()
	panicAtHttpError(err, httpResp)

	timerInfos = service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer3.Status = service.TimerSkipped
	assertions.Equal(expectedTimerInfos, timerInfos)

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

	// reset with all signals reserved (default behavior)
	// however, the skip timer won't be able to re-apply because the timers won't be ready at that moment
	req4 := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
	_, httpResp, err = req4.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
		WorkflowId: wfId,
		ResetType:  iwfidl.BEGINNING,
	}).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second)
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer2.Status = service.TimerPending
	timer3.Status = service.TimerPending
	assertions.Equal(expectedTimerInfos, timerInfos)

	req3 = apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandId:           iwfidl.PtrString("timer-cmd-id-2"),
	}).Execute()
	panicAtHttpError(err, httpResp)

	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandIndex:        iwfidl.PtrInt32(2),
	}).Execute()
	panicAtHttpError(err, httpResp)

	req2 = apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp)
}
