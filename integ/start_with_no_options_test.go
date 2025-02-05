package integ

import (
	"context"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/basic"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestStartWorkflowNoOptionsTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	doTestStartWorkflowWithoutStartOptions(t, service.BackendTypeTemporal)
}

func TestStartWorkflowNoOptionsCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	doTestStartWorkflowWithoutStartOptions(t, service.BackendTypeCadence)
}

func doTestStartWorkflowWithoutStartOptions(t *testing.T, backendType service.BackendType) {
	if !*cadenceIntegTest {
		t.Skip()
	}

	wfHandler := basic.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	client, closeFunc2 := startIwfServiceWithClient(backendType)
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := "TestStartWorkflowWithoutStartOptions" + strconv.Itoa(int(time.Now().UnixNano()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test data"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        basic.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(basic.State1),
		StateInput:             wfInput,
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	requestedSAs := []iwfidl.SearchAttributeKeyAndType{
		{
			Key:       ptr.Any(service.SearchAttributeIwfWorkflowType),
			ValueType: iwfidl.KEYWORD.Ptr(),
		},
	}
	response, err := client.DescribeWorkflowExecution(context.Background(), wfId, "", requestedSAs)
	assertions := assert.New(t)
	attribute := response.SearchAttributes[service.SearchAttributeIwfWorkflowType]
	assertions.Equal(basic.WorkflowType, attribute.GetStringValue())

	// Terminate the workflow once tests completed
	stopReq := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	_, err = stopReq.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
		StopType:   iwfidl.TERMINATE.Ptr(),
	}).Execute()
}
