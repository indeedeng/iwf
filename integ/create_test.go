package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/rpc"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestCreateWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCreateWithoutStartingState(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestCreateWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCreateWithoutStartingState(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestCreateWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCreateWithoutStartingState(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(true))
		smallWaitForFastTest()
	}
}

func TestCreateWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestCreateWithoutStartingState(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func doTestCreateWithoutStartingState(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)
	// start test workflow server
	wfHandler := rpc.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
	defer closeFunc1()

	uclient, closeFunc2 := startIwfServiceWithClient(backendType)
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
	wfId := rpc.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        rpc.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	// workflow shouldn't executed any state
	var dump service.DebugDumpResponse
	err = uclient.QueryWorkflow(context.Background(), &dump, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		panic(err)
	}
	assertions.Equal(service.StateExecutionCounterInfo{
		StateIdStartedCount:            make(map[string]int),
		StateIdCurrentlyExecutingCount: make(map[string]int),
		TotalCurrentlyExecutingCount:   0,
	}, dump.Snapshot.StateExecutionCounterInfo)

	// invoke an RPC to trigger the state execution
	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	_, httpResp, err = reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    rpc.RPCName,
		Input:      &rpc.TestInput,
		SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
			PersistenceLoadingType: iwfidl.PARTIAL_WITHOUT_LOCKING.Ptr(),
			PartialLoadingKeys: []string{
				rpc.TestSearchAttributeIntKey,
			},
		},
		TimeoutSeconds: iwfidl.PtrInt32(2),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	respWait, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, respWait)

	history, _ := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "rpc test fail, %v", history)
}
