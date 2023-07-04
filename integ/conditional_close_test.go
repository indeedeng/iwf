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
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func TestConditionalForceCompleteOnInternalChannelEmptyWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig())
		smallWaitForFastTest()
	}
}

func doTestConditionalForceCompleteOnInternalChannelEmptyWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := conditionalClose.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
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
	wfId := conditionalClose.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        conditionalClose.WorkflowType,
		WorkflowTimeoutSeconds: 20,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(conditionalClose.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for a second so that query handler is ready for executing PRC
	time.Sleep(time.Second)
	// invoke RPC to send 1 messages to the internal channel to unblock the waitUntil
	// then send another two messages
	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	for i := 0; i < 3; i++ {
		_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
			WorkflowId: wfId,
			RpcName:    conditionalClose.RpcPublishInternalChannel,
		}).Execute()
		panicAtHttpError(err, httpResp)
		if i == 0 {
			// wait for a second so that the workflow is in execute state
			time.Sleep(time.Second)
		}
	}

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  3,
		"S1_decide": 3,
		conditionalClose.RpcPublishInternalChannel: 3,
	}, history, "rpc test fail, %v", history)

	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())
	assertions.Equal(1, len(resp2.GetResults()))
	assertions.Equal(iwfidl.StateCompletionOutput{
		CompletedStateId:          "S1",
		CompletedStateExecutionId: "S1-3",
		CompletedStateOutput:      &conditionalClose.TestInput,
	}, resp2.GetResults()[0])
}
