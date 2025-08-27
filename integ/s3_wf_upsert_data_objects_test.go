package integ

import (
	"context"
	"strconv"
	"testing"
	"time"

	s3_upsert_data_objects "github.com/indeedeng/iwf/integ/workflow/s3-upsert-data-objects"

	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestS3WorkflowUpsertDataObjectsTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3UpsertDataObjects(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3WorkflowUpsertDataObjectsCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3UpsertDataObjects(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

func doTestWorkflowWithS3UpsertDataObjects(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3_upsert_data_objects.NewHandler()
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
	wfId := s3_upsert_data_objects.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	// Create small input
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"test\""),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3_upsert_data_objects.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3_upsert_data_objects.State1),
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

	// Verify all states were executed
	assertions.Equal(history["S1_start"], int64(1), "S1_start should be called once")
	assertions.Equal(history["S1_decide"], int64(1), "S1_decide should be called once")
	assertions.Equal(history["S2_start"], int64(1), "S2_start should be called once")
	assertions.Equal(history["S2_decide"], int64(1), "S2_decide should be called once")

	// Verify State2 received the large data objects that were upserted by State1
	assertions.Equal(history["S2_received_large_obj1"], true, "S2 should receive large_obj1 from State1's upsert")
	assertions.Equal(history["S2_received_large_obj2"], true, "S2 should receive large_obj2 from State1's upsert")
	assertions.Equal(history["S2_received_small_obj3"], true, "S2 should receive small_obj3 from State1's upsert")

	// Verify the data content matches what State1 upserted
	assertions.Equal(history["S2_large_obj1_data"], s3_upsert_data_objects.LargeDataContent1, "S2 large_obj1 data should match")
	assertions.Equal(history["S2_large_obj2_data"], s3_upsert_data_objects.LargeDataContent2, "S2 large_obj2 data should match")
	assertions.Equal(history["S2_small_obj3_data"], s3_upsert_data_objects.SmallDataContent3, "S2 small_obj3 data should match")

	// Verify external storage usage: 2 large objects should be in S3, small one should stay in memory
	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	assertions.Equal(int64(2), objectCount, "Should have 2 objects in S3 (large_obj1 and large_obj2)")
}
