package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	conditionalClose "github.com/indeedeng/iwf/integ/workflow/conditional_close"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestConditionalForceCompleteOnInternalChannelEmptyWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestConditionalForceCompleteOnInternalChannelEmptyWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestConditionalForceCompleteOnInternalChannelEmptyWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		// TODO not sure why using minimumContinueAsNewConfig(true) will fail
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func TestConditionalForceCompleteOnInternalChannelEmptyWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(
	t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig,
) {
	doTestConditionalForceCompleteOnChannelEmptyWorkflow(t, backendType, config, false)
	doTestConditionalForceCompleteOnChannelEmptyWorkflow(t, backendType, config, true)
}

func doTestConditionalForceCompleteOnChannelEmptyWorkflow(
	t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig, useSignalChannel bool,
) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := conditionalClose.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceWithClient(backendType)
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
	channelType := "_internal_channel_"
	if useSignalChannel {
		channelType = "_signal_channel_"
	}
	wfId := conditionalClose.WorkflowType + channelType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        conditionalClose.WorkflowType,
		WorkflowTimeoutSeconds: 20,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(conditionalClose.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	if useSignalChannel {
		startReq.StateInput = &iwfidl.EncodedObject{
			Data: iwfidl.PtrString("use-signal-channel"),
		} // this will tell the workflow to use signal
	}

	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	// Wait for a second so that query handler is ready for executing PRC
	time.Sleep(time.Second)
	// invoke RPC to send 1 messages to the internal channel to unblock the waitUntil
	// then send another two messages
	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	reqSignal := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	for i := 0; i < 3; i++ {
		if useSignalChannel {
			httpResp, err = reqSignal.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
				WorkflowId:        wfId,
				SignalChannelName: conditionalClose.TestChannelName,
			}).Execute()
		} else {
			_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
				WorkflowId: wfId,
				RpcName:    conditionalClose.RpcPublishInternalChannel,
			}).Execute()
		}

		failTestAtHttpError(err, httpResp, t)
		if i == 0 {
			// Wait for a second so that the workflow is in execute state
			time.Sleep(time.Second)
		}
	}

	// Wait for the workflow to complete
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()

	expectMap := map[string]int64{
		"S1_start":  4,
		"S1_decide": 4,
	}
	if useSignalChannel {
		expectMap = map[string]int64{
			"S1_start":  3,
			"S1_decide": 3,
		}
	}
	if !useSignalChannel {
		expectMap[conditionalClose.RpcPublishInternalChannel] = 3
	}
	assertions.Equalf(expectMap, history, "rpc test fail, %v", history)

	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())
	assertions.Equal(1, len(resp2.GetResults()))
	expectedOutput := iwfidl.StateCompletionOutput{
		CompletedStateId:          "S1",
		CompletedStateExecutionId: "S1-4",
		CompletedStateOutput:      &conditionalClose.TestInput,
	}
	if useSignalChannel {
		expectedOutput = iwfidl.StateCompletionOutput{
			CompletedStateId:          "S1",
			CompletedStateExecutionId: "S1-3",
			CompletedStateOutput:      &conditionalClose.TestInput,
		}
	}
	assertions.Equal(expectedOutput, resp2.GetResults()[0])
}
