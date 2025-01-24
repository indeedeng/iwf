package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestSetDataAttributesTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	assertions := assert.New(t)

	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceWithClient(service.BackendTypeTemporal)
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
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions.Equal(httpResp.StatusCode, http.StatusOK)

	smallDataObjects := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey),
			Value: &persistence.TestDataObjectVal1,
		},
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey2),
			Value: &persistence.TestDataObjectVal2,
		},
	}

	setReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsSetPost(context.Background())
	httpResp2, err := setReq.WorkflowSetDataObjectsRequest(iwfidl.WorkflowSetDataObjectsRequest{
		WorkflowId: wfId,
		Objects:    smallDataObjects,
	}).Execute()

	panicAtHttpError(err, httpResp2)

	time.Sleep(time.Second)

	getReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	getResult, httpRespGet, err := getReq.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			persistence.TestDataObjectKey, persistence.TestDataObjectKey2,
		}}).Execute()
	panicAtHttpError(err, httpRespGet)

	assertions.ElementsMatch(smallDataObjects, getResult.Objects)
}
