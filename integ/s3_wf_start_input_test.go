package integ

import (
	"context"
	"strconv"
	"testing"
	"time"

	s3_start_input "github.com/indeedeng/iwf/integ/workflow/s3-start-input"

	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestS3WorkflowStartInputTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3StartInput(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3WorkflowStartInputCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3StartInput(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

func doTestWorkflowWithS3StartInput(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3_start_input.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 10,
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
	wfId := s3_start_input.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"12345678901\""), //11 + 2bytes
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3_start_input.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3_start_input.State1),
		StateInput:             wfInput,
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp2, err2 := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err2, httpResp2, t)

	assertions := assert.New(t)

	_, history := wfHandler.GetTestResult()

	// The handler should receive objects with both the loaded data AND preserved external storage references
	s1StartInput := history["S1_start_input"].(iwfidl.EncodedObject)
	s1DecideInput := history["S1_decide_input"].(iwfidl.EncodedObject)

	// Verify the data content is correct
	assertions.Equal(*s1StartInput.Data, "\"12345678901\"", "S1_start_input data should match")
	assertions.Equal(*s1StartInput.Encoding, "json", "S1_start_input encoding should match")
	assertions.NotNil(s1StartInput.ExtStoreId, "S1_start_input should have ExtStoreId preserved")
	assertions.NotNil(s1StartInput.ExtPath, "S1_start_input should have ExtPath preserved")

	assertions.Equal(*s1DecideInput.Data, "\"12345678901\"", "S1_decide_input data should match")
	assertions.Equal(*s1DecideInput.Encoding, "json", "S1_decide_input encoding should match")
	assertions.NotNil(s1DecideInput.ExtStoreId, "S1_decide_input should have ExtStoreId preserved")
	assertions.NotNil(s1DecideInput.ExtPath, "S1_decide_input should have ExtPath preserved")

	// Verify that both start and decide inputs reference the same external storage location
	assertions.Equal(*s1StartInput.ExtStoreId, *s1DecideInput.ExtStoreId, "Both inputs should have the same ExtStoreId")
	assertions.Equal(*s1StartInput.ExtPath, *s1DecideInput.ExtPath, "Both inputs should have the same ExtPath")

	assertions.Equal(history["S1_start"], int64(1), "S1_start is not equal")
	assertions.Equal(history["S1_decide"], int64(1), "S1_decide is not equal")

	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	assertions.Equal(int64(1), objectCount)
}
