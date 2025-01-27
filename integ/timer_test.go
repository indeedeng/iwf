package integ

import (
	"context"
	"log"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/timer"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestTimerWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestTimerWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestTimerWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func TestTimerWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestTimerWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestTimerWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := timer.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	uclient, closeFunc2 := startIwfServiceWithClient(backendType)
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
		StartStateId:           ptr.Any(timer.State1),
		StateInput: &iwfidl.EncodedObject{
			Data: iwfidl.PtrString(strconv.Itoa(int(nowTimestamp))),
		},
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	time.Sleep(time.Second * 1)
	timerInfos := service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	assertions := assert.New(t)
	timer2 := &service.TimerInfo{
		CommandId:                  ptr.Any("timer-cmd-id-2"),
		FiringUnixTimestampSeconds: nowTimestamp + 86400,
		Status:                     service.TimerPending,
	}
	timer3 := &service.TimerInfo{
		CommandId:                  ptr.Any("timer-cmd-id-3"),
		FiringUnixTimestampSeconds: nowTimestamp + 86400*365,
		Status:                     service.TimerPending,
	}
	expectedTimerInfos := service.GetCurrentTimerInfosQueryResponse{
		StateExecutionCurrentTimerInfos: map[string][]*service.TimerInfo{
			"S1-1": {
				{
					CommandId:                  ptr.Any("timer-cmd-id"),
					FiringUnixTimestampSeconds: nowTimestamp + 10,
					Status:                     service.TimerPending,
				},
				timer2,
				timer3,
			},
		},
	}
	assertTimerQueryResponseEqual(assertions, expectedTimerInfos, timerInfos)

	req3 := apiClient.DefaultApi.ApiV1WorkflowTimerSkipPost(context.Background())
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandId:           iwfidl.PtrString("timer-cmd-id-2"),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	time.Sleep(time.Second * 1)
	timerInfos = service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer2.Status = service.TimerSkipped
	assertTimerQueryResponseEqual(assertions, expectedTimerInfos, timerInfos)

	time.Sleep(time.Second * 1)
	httpResp, err = req3.WorkflowSkipTimerRequest(iwfidl.WorkflowSkipTimerRequest{
		WorkflowId:               wfId,
		WorkflowStateExecutionId: "S1-1",
		TimerCommandIndex:        iwfidl.PtrInt32(2),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	time.Sleep(time.Second * 1)
	timerInfos = service.GetCurrentTimerInfosQueryResponse{}
	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer3.Status = service.TimerSkipped
	assertTimerQueryResponseEqual(assertions, expectedTimerInfos, timerInfos)

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

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
	// Therefore, the skip timer would be reapplied
	req4 := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
	_, httpResp, err = req4.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
		WorkflowId: wfId,
		ResetType:  iwfidl.BEGINNING,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	err = uclient.QueryWorkflow(context.Background(), &timerInfos, wfId, "", service.GetCurrentTimerInfosQueryType)
	if err != nil {
		log.Fatalf("Fail to invoke query %v", err)
	}
	timer2.Status = service.TimerSkipped
	timer3.Status = service.TimerSkipped
	assertTimerQueryResponseEqual(assertions, expectedTimerInfos, timerInfos)

	req2 = apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp, t)
}

func assertTimerQueryResponseEqual(
	assertions *assert.Assertions, resp1 service.GetCurrentTimerInfosQueryResponse,
	resp2 service.GetCurrentTimerInfosQueryResponse,
) {
	for k, infos1 := range resp1.StateExecutionCurrentTimerInfos {
		infos2 := resp2.StateExecutionCurrentTimerInfos[k]
		assertions.Equal(len(infos1), len(infos2))
		for idx, info1 := range infos1 {
			info2 := infos2[idx]
			abs := math.Abs(float64(info1.FiringUnixTimestampSeconds - info2.FiringUnixTimestampSeconds))
			assertions.True(abs <= 1)
			info1.FiringUnixTimestampSeconds = info2.FiringUnixTimestampSeconds
			assertions.Equal(info1, info2)
		}
	}
}
