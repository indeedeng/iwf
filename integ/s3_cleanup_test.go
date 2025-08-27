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

func TestS3CleanupTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3Cleanup(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3CleanupCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithS3Cleanup(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

func doTestWorkflowWithS3Cleanup(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3_start_input.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

		
	uclient, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 10,
	})
	defer closeFunc2()

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	workflowIds := make([]string, 0)
	for i := 0; i < 12; i++ {
		wfId := s3_start_input.WorkflowType + strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(i)
		workflowIds = append(workflowIds, wfId)
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
	}

	for i := 0; i < 12; i++ {
		wfId := s3_start_input.WorkflowType + strconv.Itoa(int(time.Now().UnixNano())) + strconv.Itoa(i)
		workflowIds = append(workflowIds, wfId) // the last 12 workflows are not started, so the workflowIds are not existing workflows
	}

	// 1. use globalBlobStore to insert a lot of workflow objects
	//    24*1000 = 24000 workflows, and each workflow has 100, 200, ... 2400 objects
	//    Use the workflowIds above to insert the objects.
	// 2. use globalBlobStore to verify the number of objects, and the path of the objects. For the first 12 workflows, also include the object from the start input)
	// 3. use uclient to start StartBlobStoreCleanupWorkflow to cleanup the objects
	// 4. use uclient to wait for the workflow to complete
	// 5. use globalBlobStore to verify the number of objects, and the path of the objects

	// objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	// assertions.Nil(err)
	// assertions.Equal(int64(1), objectCount)
}
