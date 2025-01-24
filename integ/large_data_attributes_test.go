package integ

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestLargeDataAttributesTemporalContinueAsNew(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestLargeQueryAttributes(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestLargeQueryAttributes(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	if !*temporalIntegTest {
		t.Skip()
	}
	assertions := assert.New(t)

	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceWithClient(backendType)
	defer closeFunc2()

	wfId := signal.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 86400,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		// this is necessary for large DAs
		// otherwise the workflow task will fail when trying to execute a stateAPI with data attributes >4MB
		StateOptions: &signal.StateOptionsForLargeDataAttributes,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions.Equal(httpResp.StatusCode, http.StatusOK)

	// Define the size of the string in bytes (1 MB = 1024 * 1024 bytes)
	const size = 1024 * 1024

	OneMbDataObject := iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(strings.Repeat("a", size)),
	}

	// setting a large data object to test, especially continueAsNew
	// because there is a 4MB limit for GRPC in temporal
	setReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsSetPost(context.Background())
	for i := 0; i < 5; i++ {

		httpResp2, err := setReq.WorkflowSetDataObjectsRequest(iwfidl.WorkflowSetDataObjectsRequest{
			WorkflowId: wfId,
			Objects: []iwfidl.KeyValue{
				{
					Key:   iwfidl.PtrString("large-data-object-" + strconv.Itoa(i)),
					Value: &OneMbDataObject,
				},
			},
		}).Execute()

		panicAtHttpError(err, httpResp2)
	}

	// signal the workflow to complete
	for i := 0; i < 4; i++ {
		signalVal := iwfidl.EncodedObject{
			Encoding: iwfidl.PtrString("json"),
			Data:     iwfidl.PtrString(fmt.Sprintf("test-data-%v", i)),
		}

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
}
