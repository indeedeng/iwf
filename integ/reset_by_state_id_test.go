package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/reset"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestResetByStateIdWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestResetByStatIdWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()

		//TODO: uncomment below when IWF-403 implementation is done.
		//TODO cont.: Reset with state id & state execution id is broken for local activities.
		//doTestResetByStatIdWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(true))
		//smallWaitForFastTest()
	}
}

func TestResetByStateIdWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestResetByStatIdWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()

		//TODO: uncomment below when IWF-403 implementation is done.
		//TODO cont.: Reset with state id & state execution id is broken for local activities.
		//doTestResetByStatIdWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		//smallWaitForFastTest()
	}
}

func doTestResetByStatIdWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := reset.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType: backendType,
	})
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := reset.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("1"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        reset.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(reset.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
			WorkflowIDReusePolicy:  ptr.Any(iwfidl.REJECT_DUPLICATE),
		},
		StateOptions: &iwfidl.WorkflowStateOptions{
			//Skipping wait until for state1
			SkipWaitUntil: iwfidl.PtrBool(true),
		},
	}
	startResp, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	assertions := assert.New(t)

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()
	//expect no starts in history as WaitUntil api is skipped.
	assertions.Equalf(map[string]int64{
		"S1_decide": 1,
		"S2_decide": 5,
	}, history, "reset test fail, %v", history)

	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())
	assertions.Equal(1, len(resp2.GetResults()))
	assertions.Equal("S2", resp2.GetResults()[0].CompletedStateId)
	assertions.Equal("S2-5", resp2.GetResults()[0].CompletedStateExecutionId)
	assertions.Equal("5", resp2.GetResults()[0].CompletedStateOutput.GetData())

	//reset workflow by state id
	resetReq := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
	_, httpResp, err = resetReq.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
		WorkflowRunId: iwfidl.PtrString(startResp.GetWorkflowRunId()),
		WorkflowId:    wfId,
		ResetType:     iwfidl.STATE_ID,
		StateId:       iwfidl.PtrString(reset.State2),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	req3 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp3, httpResp, err := req3.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	resetHistory, _ := wfHandler.GetTestResult()
	//expect no starts in history as WaitUntil api is skipped.
	assertions.Equalf(map[string]int64{
		"S1_decide": 1,
		"S2_decide": 10,
	}, resetHistory, "reset test fail, %v", resetHistory)

	assertions.Equal(iwfidl.COMPLETED, resp3.GetWorkflowStatus())
	assertions.Equal(1, len(resp3.GetResults()))
	assertions.Equal("S2", resp3.GetResults()[0].CompletedStateId)
	assertions.Equal("S2-5", resp3.GetResults()[0].CompletedStateExecutionId)
	assertions.Equal("5", resp3.GetResults()[0].CompletedStateOutput.GetData())

	//reset workflow by state execution id
	reset2Req := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
	_, httpResp, err = reset2Req.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
		WorkflowRunId:    iwfidl.PtrString(startResp.GetWorkflowRunId()),
		WorkflowId:       wfId,
		ResetType:        iwfidl.STATE_EXECUTION_ID,
		StateExecutionId: iwfidl.PtrString(reset.State2 + "-4"),
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	req4 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp4, httpResp, err := req4.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	reset2History, _ := wfHandler.GetTestResult()
	//expect no starts in history as WaitUntil api is skipped.
	assertions.Equalf(map[string]int64{
		"S1_decide": 1,
		"S2_decide": 12,
	}, reset2History, "reset test fail, %v", reset2History)

	assertions.Equal(iwfidl.COMPLETED, resp4.GetWorkflowStatus())
	assertions.Equal(1, len(resp4.GetResults()))
	assertions.Equal("S2", resp4.GetResults()[0].CompletedStateId)
	assertions.Equal("S2-5", resp4.GetResults()[0].CompletedStateExecutionId)
	assertions.Equal("5", resp4.GetResults()[0].CompletedStateOutput.GetData())
}
