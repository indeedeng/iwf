package integ

import (
	"context"
	"strconv"
	"testing"
	"time"

	s3_init_data_attributes "github.com/indeedeng/iwf/integ/workflow/s3-init-data-attributes"

	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestS3WorkflowInitDataAttributesTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3InitDataAttributes(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3WorkflowInitDataAttributesCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3InitDataAttributes(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

func doTestWorkflowWithS3InitDataAttributes(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3_init_data_attributes.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 10, // Set low threshold so our test data gets stored in S3
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
	wfId := s3_init_data_attributes.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	// Create initial data attributes - mix of large (stored in S3) and small (kept in memory)
	initialDataAttributes := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(s3_init_data_attributes.TestDataAttrKey1),
			Value: &s3_init_data_attributes.TestDataAttributeVal1, // Large - will go to S3
		},
		{
			Key:   iwfidl.PtrString(s3_init_data_attributes.TestDataAttrKey2),
			Value: &s3_init_data_attributes.TestDataAttributeVal2, // Large - will go to S3
		},
		{
			Key:   iwfidl.PtrString(s3_init_data_attributes.TestDataAttrKey3),
			Value: &s3_init_data_attributes.TestDataAttributeVal3, // Small - will stay in memory
		},
	}

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test"),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3_init_data_attributes.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3_init_data_attributes.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			DataAttributes: initialDataAttributes,
		},
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

	// Verify State1 start (waitUntil) received data attributes from S3 matching initial values
	assertions.Equal(history["S1_start"], int64(1), "S1_start should be called once")
	assertions.Equal(history["S1_start_attr1_found"], true, "S1_start should find data attribute 1")
	assertions.Equal(history["S1_start_attr2_found"], true, "S1_start should find data attribute 2")
	assertions.Equal(history["S1_start_attr3_found"], true, "S1_start should find data attribute 3")
	assertions.Equal(history["S1_start_total_attrs"], 3, "S1_start should receive exactly 3 data attributes (no duplicates)")
	assertions.Equal(history["S1_start_attr1_data"], *s3_init_data_attributes.TestDataAttributeVal1.Data, "S1_start attr1 data should match initial value")
	assertions.Equal(history["S1_start_attr2_data"], *s3_init_data_attributes.TestDataAttributeVal2.Data, "S1_start attr2 data should match initial value")
	assertions.Equal(history["S1_start_attr3_data"], *s3_init_data_attributes.TestDataAttributeVal3.Data, "S1_start attr3 data should match initial value")
	assertions.Equal(history["S1_start_validation_pass"], true, "S1_start validation should pass - all data attributes match initial values exactly")

	// Verify State1 decide (execute) was called
	assertions.Equal(history["S1_decide"], int64(1), "S1_decide should be called once")

	// Verify State2 start (waitUntil) was called
	assertions.Equal(history["S2_start"], int64(1), "S2_start should be called once")

	// Verify State2 decide (execute) received data attributes from S3 matching initial values
	assertions.Equal(history["S2_decide"], int64(1), "S2_decide should be called once")
	assertions.Equal(history["S2_decide_attr1_found"], true, "S2_decide should find data attribute 1")
	assertions.Equal(history["S2_decide_attr2_found"], true, "S2_decide should find data attribute 2")
	assertions.Equal(history["S2_decide_attr3_found"], true, "S2_decide should find data attribute 3")
	assertions.Equal(history["S2_decide_total_attrs"], 3, "S2_decide should receive exactly 3 data attributes (no duplicates)")
	assertions.Equal(history["S2_decide_attr1_data"], *s3_init_data_attributes.TestDataAttributeVal1.Data, "S2_decide attr1 data should match initial value")
	assertions.Equal(history["S2_decide_attr2_data"], *s3_init_data_attributes.TestDataAttributeVal2.Data, "S2_decide attr2 data should match initial value")
	assertions.Equal(history["S2_decide_attr3_data"], *s3_init_data_attributes.TestDataAttributeVal3.Data, "S2_decide attr3 data should match initial value")
	assertions.Equal(history["S2_decide_validation_pass"], true, "S2_decide validation should pass - all data attributes match initial values exactly")

	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	assertions.Equal(int64(2), objectCount)
}
