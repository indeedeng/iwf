package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/deadend"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestDeadEndWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestDeadEndWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestDeadEndWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestDeadEndWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestDeadEndWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestDeadEndWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(true))
		smallWaitForFastTest()
	}
}

func TestDeadEndWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestDeadEndWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func doTestDeadEndWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := deadend.NewHandler()
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
	wfId := deadend.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        deadend.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())

	// invoke RPC to trigger write to verify continue as new is happening with no states
	for i := 0; i < 3; i++ {
		// Delay between rpc requests
		time.Sleep(time.Second * 2)
		_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
			WorkflowId: wfId,
			RpcName:    deadend.RPCWriteData,
		}).Execute()
		failTestAtHttpError(err, httpResp, t)
	}

	if config != nil {
		reqDesc := apiClient.DefaultApi.ApiV1WorkflowGetPost(context.Background())
		descResp, httpResp, err := reqDesc.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		failTestAtHttpError(err, httpResp, t)
		assertions.True(startResp.GetWorkflowRunId() != descResp.GetWorkflowRunId())
	}

	// invoke an RPC to trigger the state execution
	for i := 0; i < 3; i++ {
		// Delay between rpc requests
		time.Sleep(time.Second * 2)
		_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
			WorkflowId: wfId,
			RpcName:    deadend.RPCTriggerState,
		}).Execute()
		failTestAtHttpError(err, httpResp, t)
	}

	// Short wait for workflow
	time.Sleep(time.Second * 2)

	reqCancel := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	httpResp, err = reqCancel.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_decide": 3,
	}, history, "rpc test fail, %v", history)
}
