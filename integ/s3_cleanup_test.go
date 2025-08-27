package integ

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	s3_start_input "github.com/indeedeng/iwf/integ/workflow/s3-start-input"

	"github.com/indeedeng/iwf/service/common/blobstore"
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
		wfId := fmt.Sprintf("test-cleanup-wf-%d-%d", i, time.Now().UnixNano())
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
		wfId := fmt.Sprintf("test-cleanup-wf-%d-%d", 12+i, time.Now().UnixNano())
		workflowIds = append(workflowIds, wfId) // the last 12 workflows are not started, so the workflowIds are not existing workflows
	}

	ctx := context.Background()
	const storeId = "s3-store-id"

	// 1. Use globalBlobStore to insert a lot of workflow objects
	// For each of the 24 workflows, create varying numbers of objects (100, 200, 300, ... up to 2400)
	t.Logf("Creating objects for %d workflows", len(workflowIds))
	for i, wfId := range workflowIds {
		objectCount := (i + 1) * 100 // 100, 200, 300, ... 2400
		for j := 0; j < objectCount; j++ {
			testData := fmt.Sprintf("test-data-workflow-%d-object-%d", i, j)
			_, _, err := globalBlobStore.WriteObject(ctx, wfId, testData)
			assert.NoError(t, err)
		}
		t.Logf("Created %d objects for workflow %s", objectCount, wfId)
	}

	// 2. Verify the number of objects and paths
	t.Logf("Verifying initial object counts")
	for i, wfId := range workflowIds {
		expectedCount := int64((i + 1) * 100)
		// For the first 12 workflows (the ones that were started), include the object from start input
		if i < 12 {
			expectedCount += 1 // +1 for the start input object
		}

		actualCount, err := globalBlobStore.CountWorkflowObjectsForTesting(ctx, wfId)
		assert.NoError(t, err)
		if expectedCount > 1000 {
			expectedCount = 1000
		}
		assert.Equal(t, expectedCount, actualCount, "Workflow %s should have %d at least 1000 objects", wfId, expectedCount)
	}

	// Verify workflow paths exist (handle pagination)
	t.Logf("Verifying workflow paths with pagination")
	allWorkflowPaths := make([]string, 0)
	var continuationToken *string

	for {
		listInput := blobstore.ListObjectPathsInput{
			StoreId:           storeId,
			ContinuationToken: continuationToken,
		}
		output, err := globalBlobStore.ListWorkflowPaths(ctx, listInput)
		assert.NoError(t, err)
		assert.NotNil(t, output)

		allWorkflowPaths = append(allWorkflowPaths, output.WorkflowPaths...)
		t.Logf("Retrieved %d workflow paths in this page", len(output.WorkflowPaths))

		// Check if there are more pages
		if output.ContinuationToken == nil {
			break
		}
		continuationToken = output.ContinuationToken
	}

	t.Logf("Total workflow paths retrieved: %d", len(allWorkflowPaths))

	// Should have paths for all workflows
	todayPrefix := time.Now().Format("20060102")
	expectedPaths := make(map[string]bool)
	for _, wfId := range workflowIds {
		expectedPath := fmt.Sprintf("%s$%s", todayPrefix, wfId)
		expectedPaths[expectedPath] = false
	}

	for _, path := range allWorkflowPaths {
		if _, exists := expectedPaths[path]; exists {
			expectedPaths[path] = true
		}
	}

	foundCount := 0
	for path, found := range expectedPaths {
		if found {
			foundCount++
		} else {
			t.Logf("Missing expected path: %s", path)
		}
	}
	assert.True(t, foundCount >= len(workflowIds), "Should find paths for all workflows")

	// 3. Start the BlobStoreCleanupWorkflow
	t.Logf("Starting BlobStoreCleanupWorkflow")
	cleanupWorkflowId := "test-cleanup-" + strconv.Itoa(int(time.Now().UnixNano()))
	err := uclient.StartBlobStoreCleanupWorkflow(
		ctx,
		service.TaskQueue,
		cleanupWorkflowId,
		"",
		storeId,
	)
	assert.Nil(t, err)

	// 4. Wait for the cleanup workflow to complete
	t.Logf("Waiting for cleanup workflow to complete")
	_ = uclient.GetWorkflowResult(ctx, nil, cleanupWorkflowId, "")

	// 5. verify all the objects are deleted for the last 12 workflows, but not for the first 12 workflows
	t.Logf("Verifying cleanup results - objects should be deleted for non-existing workflows only")

	for i, wfId := range workflowIds {
		count, err := globalBlobStore.CountWorkflowObjectsForTesting(ctx, wfId)
		assert.NoError(t, err)

		if i < 12 {

			// First 12 workflows were never started, so their objects should remain
			expectedCount := int64((i + 1) * 100)
			if expectedCount > 1000 {
				expectedCount = 1000
			}
			assert.Equal(t, expectedCount, count,
				"Non-started workflow %s (index %d) should still have %d objects", wfId, i, expectedCount)
			t.Logf("✓ Non-started workflow %s (index %d): %d objects (expected %d)", wfId, i, count, expectedCount)
		} else {
			// Last 12 workflows were started, so their objects should NOT be cleaned up
			assert.Equal(t, int64(0), count,
				"Started workflow %s (index %d) should have 0 objects after cleanup", wfId, i)
			t.Logf("✓ Started workflow %s (index %d): %d objects (expected 0)", wfId, i, count)
		}
	}

	t.Logf("Cleanup test completed successfully!")
}
