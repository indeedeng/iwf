package integ

import (
	"context"
	"github.com/indeedeng/iwf/helpers"
	"github.com/indeedeng/iwf/integ/workflow/wait_until_search_attributes_optimization"
	"github.com/indeedeng/iwf/service/common/ptr"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	history "go.temporal.io/api/history/v1"
	"go.temporal.io/api/workflowservice/v1"
	"slices"
	"strconv"
	"strings"
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
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	uclient, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
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
	wfId := wait_until_search_attributes_optimization.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	reqStart := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	wfReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wait_until_search_attributes_optimization.WorkflowType,
		WorkflowTimeoutSeconds: 15,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wait_until_search_attributes_optimization.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	_, httpResp, err := reqStart.WorkflowStartRequest(wfReq).Execute()
	failTestAtHttpError(err, httpResp, t)

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

	failTestAtHttpError(err, httpResp, t)

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// wait for workflow to complete
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp, t)

	api := uclient.GetApiService().(workflowservice.WorkflowServiceClient)
	reqHistory := &workflowservice.GetWorkflowExecutionHistoryRequest{
		Namespace: "default",
		Execution: &common.WorkflowExecution{
			WorkflowId: wfId,
		},
	}
	eventHistory, err := api.GetWorkflowExecutionHistory(context.Background(), reqHistory)
	if err != nil {
		helpers.FailTestWithErrorMessage("couldn't load eventHistory", t)
	}

	var upsertSAEvents []*history.HistoryEvent

	for _, e := range eventHistory.History.Events {
		if e.EventType == enums.EVENT_TYPE_UPSERT_WORKFLOW_SEARCH_ATTRIBUTES {
			upsertSAEvents = append(upsertSAEvents, e)
		}
	}

	switch mode := config.GetExecutingStateIdMode(); mode {
	case iwfidl.ENABLED_FOR_ALL:
		assertions.Equal(10, len(upsertSAEvents))
		assertions.Equal([]string{"S1"}, historyEventSAs(upsertSAEvents[1]))
		assertions.Equal([]string{"S2"}, historyEventSAs(upsertSAEvents[2]))
		assertions.Equal([]string{"S2", "S3"}, historyEventSAs(upsertSAEvents[3]))
		assertions.Equal([]string{"S3", "S4"}, historyEventSAs(upsertSAEvents[4]))
		assertions.Equal([]string{"S3", "S5"}, historyEventSAs(upsertSAEvents[5]))
		assertions.Equal([]string{"S3", "S6", "S7"}, historyEventSAs(upsertSAEvents[6]))
		assertions.Equal([]string{"S3", "S6"}, historyEventSAs(upsertSAEvents[7]))
		assertions.Equal([]string{"S3"}, historyEventSAs(upsertSAEvents[8]))
		assertions.Equal([]string{"null"}, historyEventSAs(upsertSAEvents[9]))
	case iwfidl.DISABLED:
		assertions.Equal(1, len(upsertSAEvents))
	case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
	default:
		assertions.Equal(9, len(upsertSAEvents))
		assertions.Equal([]string{"S1"}, historyEventSAs(upsertSAEvents[1]))
		assertions.Equal([]string{"null"}, historyEventSAs(upsertSAEvents[2]))
		assertions.Equal([]string{"S3"}, historyEventSAs(upsertSAEvents[3]))
		assertions.Equal([]string{"S3", "S4"}, historyEventSAs(upsertSAEvents[4]))
		assertions.Equal([]string{"S3"}, historyEventSAs(upsertSAEvents[5]))
		assertions.Equal([]string{"S3", "S6"}, historyEventSAs(upsertSAEvents[6]))
		assertions.Equal([]string{"S3"}, historyEventSAs(upsertSAEvents[7]))
		assertions.Equal([]string{"null"}, historyEventSAs(upsertSAEvents[8]))
	}
}

func historyEventSAs(e *history.HistoryEvent) []string {
	attrs := e.GetAttributes().(*history.HistoryEvent_UpsertWorkflowSearchAttributesEventAttributes)
	data := string(attrs.UpsertWorkflowSearchAttributesEventAttributes.GetSearchAttributes().GetIndexedFields()[service.SearchAttributeExecutingStateIds].GetData())
	data = strings.ReplaceAll(data, "[", "")
	data = strings.ReplaceAll(data, "]", "")
	data = strings.ReplaceAll(data, "\"", "")
	dataSlice := strings.Split(data, ",")
	slices.Sort(dataSlice)
	return dataSlice
}
