package integ

import (
	"context"
	"encoding/json"
	"fmt"
	config2 "github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSignalWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestSignalWorkflowTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func TestSignalWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func TestSignalWorkflowCadenceContinueAsNew(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestSignalWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestSignalWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)

	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorkerWithRpc(wfHandler)
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
	wfId := signal.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 20,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	// test update config
	var debugDump service.DebugDumpResponse
	err = uclient.QueryWorkflow(context.Background(), &debugDump, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		panic(err)
	}
	expectedConfig := *config2.DefaultWorkflowConfig
	if config != nil {
		expectedConfig = *config
	}
	assertions.Equal(service.DebugDumpResponse{
		Config: expectedConfig,
	}, debugDump)

	// update the disable system SA
	reqUpdateConfig := apiClient.DefaultApi.ApiV1WorkflowConfigUpdatePost(context.Background())
	httpResp, err = reqUpdateConfig.WorkflowConfigUpdateRequest(iwfidl.WorkflowConfigUpdateRequest{
		WorkflowId: wfId,
		WorkflowConfig: iwfidl.WorkflowConfig{
			DisableSystemSearchAttribute: iwfidl.PtrBool(true),
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	err = uclient.QueryWorkflow(context.Background(), &debugDump, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		panic(err)
	}
	expectedConfig.DisableSystemSearchAttribute = iwfidl.PtrBool(true)
	assertions.Equal(service.DebugDumpResponse{
		Config: expectedConfig,
	}, debugDump)

	// update the pagination size
	reqUpdateConfig = apiClient.DefaultApi.ApiV1WorkflowConfigUpdatePost(context.Background())
	httpResp, err = reqUpdateConfig.WorkflowConfigUpdateRequest(iwfidl.WorkflowConfigUpdateRequest{
		WorkflowId: wfId,
		WorkflowConfig: iwfidl.WorkflowConfig{
			ContinueAsNewPageSizeInBytes: iwfidl.PtrInt32(300),
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	err = uclient.QueryWorkflow(context.Background(), &debugDump, wfId, "", service.DebugDumpQueryType)
	if err != nil {
		panic(err)
	}
	expectedConfig.ContinueAsNewPageSizeInBytes = iwfidl.PtrInt32(300)
	assertions.Equal(service.DebugDumpResponse{
		Config: expectedConfig,
	}, debugDump)

	// signal for testing unhandled signals
	var unhandledSignalVals []*iwfidl.EncodedObject
	for i := 0; i < 10; i++ {
		sigVal := &iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(fmt.Sprintf("test-data-%v", i)),
		}
		req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		httpResp2, _ := req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: signal.UnhandledSignalName,
			SignalValue:       sigVal,
		}).Execute()
		if httpResp2.StatusCode == http.StatusOK {
			// see why in https://github.com/temporalio/temporal/issues/4801
			unhandledSignalVals = append(unhandledSignalVals, sigVal)
		}
		reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
		rpcResp, httpResp2, err2 := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
			WorkflowId: wfId,
			RpcName:    signal.RPCNameGetSignalChannelInfo,
		}).Execute()
		panicAtHttpError(err2, httpResp2)
		var infos map[string]iwfidl.ChannelInfo
		err = json.Unmarshal([]byte(rpcResp.Output.GetData()), &infos)
		panicAError(err)
		assertions.Equal(
			map[string]iwfidl.ChannelInfo{signal.UnhandledSignalName: {Size: ptr.Any(int32(i + 1))}}, infos)
	}

	reqRpc := apiClient.DefaultApi.ApiV1WorkflowRpcPost(context.Background())
	rpcResp, httpResp2, err2 := reqRpc.WorkflowRpcRequest(iwfidl.WorkflowRpcRequest{
		WorkflowId: wfId,
		RpcName:    signal.RPCNameGetInternalChannelInfo,
	}).Execute()
	panicAtHttpError(err2, httpResp2)
	var infos map[string]iwfidl.ChannelInfo
	err = json.Unmarshal([]byte(rpcResp.Output.GetData()), &infos)
	panicAError(err)
	assertions.Equal(
		map[string]iwfidl.ChannelInfo{signal.InternalChannelName: {Size: ptr.Any(int32(10))}}, infos)

	// signal the workflow
	var signalVals []iwfidl.EncodedObject
	for i := 0; i < 4; i++ {
		signalVal := iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(fmt.Sprintf("test-data-%v", i)),
		}
		signalVals = append(signalVals, signalVal)

		req2 := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
		httpResp2, err := req2.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
			WorkflowId:        wfId,
			SignalChannelName: signal.SignalName,
			SignalValue:       &signalVal,
		}).Execute()

		panicAtHttpError(err, httpResp2)
	}

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	history, data := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "signal test fail, %v", history)

	assertions.Equal(fmt.Sprintf("signal-cmd-id%v", 0), data[fmt.Sprintf("signalId%v", 0)])
	assertions.Equal(fmt.Sprintf("signal-cmd-id%v", 1), data[fmt.Sprintf("signalId%v", 1)])
	assertions.Equal("", data[fmt.Sprintf("signalId%v", 2)])
	assertions.Equal("", data[fmt.Sprintf("signalId%v", 3)])
	for i := 0; i < 4; i++ {
		assertions.Equal(signalVals[i], data[fmt.Sprintf("signalValue%v", i)])
	}

	var dump service.ContinueAsNewDumpResponse
	err = uclient.QueryWorkflow(context.Background(), &dump, wfId, "", service.ContinueAsNewDumpQueryType)
	if err != nil {
		panic(err)
	}
	assertions.Equal(unhandledSignalVals, dump.SignalsReceived[signal.UnhandledSignalName])
	assertions.True(len(unhandledSignalVals) > 0)

	if config == nil {
		// TODO add assertion for continueAsNew case

		// reset with all signals reserved (default behavior)
		// reset to beginning
		req4 := apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
		_, httpResp, err = req4.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
			WorkflowId: wfId,
			ResetType:  iwfidl.BEGINNING,
		}).Execute()
		panicAtHttpError(err, httpResp)

		reqWait = apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
		resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp)

		// reset to STATE_EXECUTION_ID
		req4 = apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
		_, httpResp, err = req4.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
			WorkflowId:       wfId,
			ResetType:        iwfidl.STATE_EXECUTION_ID,
			StateExecutionId: iwfidl.PtrString("S2-1"),
		}).Execute()
		panicAtHttpError(err, httpResp)

		reqWait = apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
		resp, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp)

		// reset to STATE_ID
		req4 = apiClient.DefaultApi.ApiV1WorkflowResetPost(context.Background())
		_, httpResp, err = req4.WorkflowResetRequest(iwfidl.WorkflowResetRequest{
			WorkflowId: wfId,
			ResetType:  iwfidl.STATE_ID,
			StateId:    iwfidl.PtrString("S2"),
		}).Execute()
		panicAtHttpError(err, httpResp)

		reqWait = apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
		resp, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
			WorkflowId: wfId,
		}).Execute()
		panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp)
	}

}
