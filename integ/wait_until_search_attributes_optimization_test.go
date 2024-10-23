package integ

import (
	"context"
	"github.com/indeedeng/iwf/integ/workflow/wait_until_search_attributes_optimization"
	"github.com/indeedeng/iwf/service/common/ptr"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	history "go.temporal.io/api/history/v1"
	"go.temporal.io/api/workflowservice/v1"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestWaitUntilSearchAttributesOptimizationWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilHistoryCompleted(t, service.BackendTypeTemporal, &iwfidl.WorkflowConfig{
			ExecutingStateIdMode: ptr.Any(iwfidl.DISABLED),
		})
		smallWaitForFastTest()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilHistoryCompleted(t, service.BackendTypeTemporal, &iwfidl.WorkflowConfig{
			ExecutingStateIdMode: ptr.Any(iwfidl.ENABLED_FOR_ALL),
		})
		smallWaitForFastTest()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilHistoryCompleted(t, service.BackendTypeTemporal, nil) // defaults to ExecutingStateIdMode: ENABLED_FOR_STATES_WITH_WAIT_UNTIL
		smallWaitForFastTest()
	}
}

func doTestWaitUntilHistoryCompleted(
	t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig,
) {
	assertions := assert.New(t)
	wfHandler := wait_until_search_attributes_optimization.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
	defer closeFunc1()

	uclient, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:         backendType,
		OptimizedVersioning: ptr.Any(true),
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
	wfId := wait_until_search_attributes_optimization.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	reqStart := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	wfReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wait_until_search_attributes_optimization.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wait_until_search_attributes_optimization.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	_, httpResp, err := reqStart.WorkflowStartRequest(wfReq).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(time.Second * 5)

	signalValue := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test"),
	}

	reqSignal := apiClient.DefaultApi.ApiV1WorkflowSignalPost(context.Background())
	httpResp, err = reqSignal.WorkflowSignalRequest(iwfidl.WorkflowSignalRequest{
		WorkflowId:        wfId,
		SignalChannelName: wait_until_search_attributes_optimization.SignalName,
		SignalValue:       &signalValue,
	}).Execute()

	panicAtHttpError(err, httpResp)

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for workflow to complete
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp)

	api := uclient.GetApiService().(workflowservice.WorkflowServiceClient)
	reqHistory := &workflowservice.GetWorkflowExecutionHistoryRequest{
		Namespace: "default",
		Execution: &common.WorkflowExecution{
			WorkflowId: wfId,
		},
	}
	eventHistory, err := api.GetWorkflowExecutionHistory(context.Background(), reqHistory)
	if err != nil {
		panic("couldn't load eventHistory")
	}

	var upsertSAEvents []*history.HistoryEvent

	for _, e := range eventHistory.History.Events {
		if e.EventType == enums.EVENT_TYPE_UPSERT_WORKFLOW_SEARCH_ATTRIBUTES {
			upsertSAEvents = append(upsertSAEvents, e)
		}
	}

	switch mode := config.GetExecutingStateIdMode(); mode {
	case iwfidl.ENABLED_FOR_ALL:
		assertions.Equal(6, len(upsertSAEvents))
		assertions.Equal("[\"S1\"]", historyEventSAs(upsertSAEvents[0]))
		assertions.Equal("[\"S2\"]", historyEventSAs(upsertSAEvents[1]))
		assertions.Equal("[\"S3\",\"S4\"]", historyEventSAs(upsertSAEvents[2]))
		assertions.Equal("[\"S4\"]", historyEventSAs(upsertSAEvents[3]))
		assertions.Equal("[\"S5\"]", historyEventSAs(upsertSAEvents[4]))
		assertions.Equal("null", historyEventSAs(upsertSAEvents[5]))
	case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
		assertions.Equal(5, len(upsertSAEvents))
		assertions.Equal("[\"S1\"]", historyEventSAs(upsertSAEvents[0]))
		assertions.Equal("[\"S2\"]", historyEventSAs(upsertSAEvents[1]))
		assertions.Equal("[\"S3\",\"S4\"]", historyEventSAs(upsertSAEvents[2]))
		assertions.Equal("[\"S4\"]", historyEventSAs(upsertSAEvents[3]))
		assertions.Equal("null", historyEventSAs(upsertSAEvents[4]))
	case iwfidl.DISABLED:
		assertions.Equal(0, len(upsertSAEvents))
	}
}

func historyEventSAs(e *history.HistoryEvent) string {
	attrs := e.GetAttributes().(*history.HistoryEvent_UpsertWorkflowSearchAttributesEventAttributes)
	return string(attrs.UpsertWorkflowSearchAttributesEventAttributes.GetSearchAttributes().GetIndexedFields()[service.SearchAttributeExecutingStateIds].GetData())
}
