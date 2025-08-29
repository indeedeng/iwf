package integ

import (
	"context"
	"strconv"
	"testing"
	"time"

	s3GetSetDataAttributes "github.com/indeedeng/iwf/integ/workflow/s3-get-set-data-attributes"

	"github.com/indeedeng/iwf/service/common/ptr"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestS3GetSetDataAttributesTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestS3GetSetDataAttributes(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3GetSetDataAttributesCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestS3GetSetDataAttributes(t, service.BackendTypeCadence)
		smallWaitForFastTest()
	}
}

// useMemo tests - only supported on Temporal (Cadence doesn't support memo feature)
func TestS3GetSetDataAttributesWithUseMemoTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestS3GetSetDataAttributesWithUseMemo(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func TestS3GetSetDataAttributesWithUseMemoAndInitialDataTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestS3GetSetDataAttributesWithUseMemoAndInitialData(t, service.BackendTypeTemporal)
		smallWaitForFastTest()
	}
}

func doTestS3GetSetDataAttributes(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3GetSetDataAttributes.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 50, // Set low threshold so large test data gets stored in S3
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
	wfId := s3GetSetDataAttributes.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-input"),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3GetSetDataAttributes.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3GetSetDataAttributes.State1),
		StateInput:             wfInput,
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	assertions := assert.New(t)

	// Test 1: Set data attributes with mix of small and large data
	testDataAttributes := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(s3GetSetDataAttributes.SmallDataKey),
			Value: &s3GetSetDataAttributes.SmallDataValue, // Small - stays in Temporal
		},
		{
			Key:   iwfidl.PtrString(s3GetSetDataAttributes.LargeDataKey),
			Value: &s3GetSetDataAttributes.LargeDataValue, // Large - goes to S3
		},
		{
			Key:   iwfidl.PtrString(s3GetSetDataAttributes.AnotherLargeDataKey),
			Value: &s3GetSetDataAttributes.AnotherLargeDataValue, // Large - goes to S3
		},
	}

	setReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsSetPost(context.Background())
	httpResp2, err := setReq.WorkflowSetDataObjectsRequest(iwfidl.WorkflowSetDataObjectsRequest{
		WorkflowId: wfId,
		Objects:    testDataAttributes,
	}).Execute()
	failTestAtHttpError(err, httpResp2, t)

	// Wait for data to be processed
	time.Sleep(time.Millisecond * 500)

	// Test 2: Get all data attributes and verify they were loaded from external storage correctly
	getReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	getResult, httpRespGet, err := getReq.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			s3GetSetDataAttributes.SmallDataKey,
			s3GetSetDataAttributes.LargeDataKey,
			s3GetSetDataAttributes.AnotherLargeDataKey,
		},
	}).Execute()
	failTestAtHttpError(err, httpRespGet, t)

	// Verify we got all 3 data attributes back
	assertions.Equal(3, len(getResult.Objects), "Should return exactly 3 data attributes")

	// Create a map for easier lookup
	retrievedAttrs := make(map[string]iwfidl.EncodedObject)
	for _, attr := range getResult.Objects {
		retrievedAttrs[*attr.Key] = *attr.Value
	}

	// Test 3: Verify small data attribute (should have actual data, no external storage references)
	smallAttr, exists := retrievedAttrs[s3GetSetDataAttributes.SmallDataKey]
	assertions.True(exists, "Small data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Data, *smallAttr.Data, "Small data content should match")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Encoding, *smallAttr.Encoding, "Small data encoding should match")

	// Test 4: Verify large data attributes (should have actual data loaded from S3)
	largeAttr, exists := retrievedAttrs[s3GetSetDataAttributes.LargeDataKey]
	assertions.True(exists, "Large data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Data, *largeAttr.Data, "Large data content should match")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Encoding, *largeAttr.Encoding, "Large data encoding should match")
	// Large data should also preserve external storage references for optimization
	assertions.NotNil(largeAttr.ExtStoreId, "Large data should preserve ExtStoreId reference")
	assertions.NotNil(largeAttr.ExtPath, "Large data should preserve ExtPath reference")

	anotherLargeAttr, exists := retrievedAttrs[s3GetSetDataAttributes.AnotherLargeDataKey]
	assertions.True(exists, "Another large data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.AnotherLargeDataValue.Data, *anotherLargeAttr.Data, "Another large data content should match")
	assertions.Equal(*s3GetSetDataAttributes.AnotherLargeDataValue.Encoding, *anotherLargeAttr.Encoding, "Another large data encoding should match")
	assertions.NotNil(anotherLargeAttr.ExtStoreId, "Another large data should preserve ExtStoreId reference")
	assertions.NotNil(anotherLargeAttr.ExtPath, "Another large data should preserve ExtPath reference")

	// Test 5: Get specific keys only
	getSpecificReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	getSpecificResult, httpRespGetSpecific, err := getSpecificReq.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			s3GetSetDataAttributes.LargeDataKey,
		},
	}).Execute()
	failTestAtHttpError(err, httpRespGetSpecific, t)

	assertions.Equal(1, len(getSpecificResult.Objects), "Should return exactly 1 data attribute when requesting specific key")
	assertions.Equal(s3GetSetDataAttributes.LargeDataKey, *getSpecificResult.Objects[0].Key, "Should return the requested key")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Data, *getSpecificResult.Objects[0].Value.Data, "Specific large data content should match")

	// Note: Skip updating existing data attributes test since the workflow has completed.
	// This is expected behavior as our test workflow is designed to complete after the first state.

	// Test 6: Verify S3 objects were created for large data attributes
	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	// Should have exactly 2 objects for the 2 large data attributes
	assertions.Equal(int64(2), objectCount, "Should have exactly 2 S3 objects for large data attributes")

	// Complete the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp2, err2 := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err2, httpResp2, t)
}

func doTestS3GetSetDataAttributesWithUseMemo(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3GetSetDataAttributes.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 50, // Set low threshold so large test data gets stored in S3
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
	wfId := s3GetSetDataAttributes.WorkflowType + strconv.Itoa(int(time.Now().UnixNano())) + "-memo"

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-input"),
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3GetSetDataAttributes.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3GetSetDataAttributes.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			UseMemoForDataAttributes: ptr.Any(true), // Enable memo-based data attributes
		},
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	assertions := assert.New(t)

	// Test 1: Set data attributes with mix of small and large data (using memo)
	testDataAttributes := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(s3GetSetDataAttributes.SmallDataKey),
			Value: &s3GetSetDataAttributes.SmallDataValue, // Small - stays in memo
		},
		{
			Key:   iwfidl.PtrString(s3GetSetDataAttributes.LargeDataKey),
			Value: &s3GetSetDataAttributes.LargeDataValue, // Large - goes to S3
		},
	}

	setReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsSetPost(context.Background())
	httpResp2, err := setReq.WorkflowSetDataObjectsRequest(iwfidl.WorkflowSetDataObjectsRequest{
		WorkflowId: wfId,
		Objects:    testDataAttributes,
	}).Execute()
	failTestAtHttpError(err, httpResp2, t)

	// Wait for data to be processed
	time.Sleep(time.Millisecond * 500)

	// Test 2: Get data attributes using memo
	getReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	getResult, httpRespGet, err := getReq.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId:               wfId,
		UseMemoForDataAttributes: ptr.Any(true), // Use memo for retrieval
		Keys: []string{
			s3GetSetDataAttributes.SmallDataKey,
			s3GetSetDataAttributes.LargeDataKey,
		},
	}).Execute()
	failTestAtHttpError(err, httpRespGet, t)

	// Verify we got both data attributes back
	assertions.Equal(2, len(getResult.Objects), "Should return exactly 2 data attributes")

	// Create a map for easier lookup
	retrievedAttrs := make(map[string]iwfidl.EncodedObject)
	for _, attr := range getResult.Objects {
		retrievedAttrs[*attr.Key] = *attr.Value
	}

	// Test 3: Verify small data attribute (memo-based)
	smallAttr, exists := retrievedAttrs[s3GetSetDataAttributes.SmallDataKey]
	assertions.True(exists, "Small data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Data, *smallAttr.Data, "Small data content should match")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Encoding, *smallAttr.Encoding, "Small data encoding should match")

	// Test 4: Verify large data attribute (should have actual data loaded from S3)
	largeAttr, exists := retrievedAttrs[s3GetSetDataAttributes.LargeDataKey]
	assertions.True(exists, "Large data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Data, *largeAttr.Data, "Large data content should match")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Encoding, *largeAttr.Encoding, "Large data encoding should match")
	// Large data should also preserve external storage references for optimization
	assertions.NotNil(largeAttr.ExtStoreId, "Large data should preserve ExtStoreId reference")
	assertions.NotNil(largeAttr.ExtPath, "Large data should preserve ExtPath reference")

	// Test 5: Verify S3 objects were created only for large data
	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	assertions.Equal(int64(1), objectCount, "Should have exactly 1 S3 object for large data attribute")

	// Complete the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp2, err2 := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err2, httpResp2, t)
}

func doTestS3GetSetDataAttributesWithUseMemoAndInitialData(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := s3GetSetDataAttributes.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType:     backendType,
		S3TestThreshold: 50, // Set low threshold so large test data gets stored in S3
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
	wfId := s3GetSetDataAttributes.WorkflowType + strconv.Itoa(int(time.Now().UnixNano())) + "-memo-initial"

	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-input"),
	}

	// Initial data attributes to be set at workflow start
	initialDataAttributes := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString("initial-small"),
			Value: &s3GetSetDataAttributes.SmallDataValue, // Small - stays in memo
		},
		{
			Key:   iwfidl.PtrString("initial-large"),
			Value: &s3GetSetDataAttributes.LargeDataValue, // Large - goes to S3
		},
	}

	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        s3GetSetDataAttributes.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(s3GetSetDataAttributes.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			UseMemoForDataAttributes: ptr.Any(true),         // Enable memo-based data attributes
			DataAttributes:           initialDataAttributes, // Set initial data attributes
		},
	}
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	failTestAtHttpError(err, httpResp, t)

	assertions := assert.New(t)

	// Wait for workflow to start and process initial data
	time.Sleep(time.Millisecond * 500)

	// Test 1: Get initial data attributes using memo
	getReq := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	getResult, httpRespGet, err := getReq.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId:               wfId,
		UseMemoForDataAttributes: ptr.Any(true), // Use memo for retrieval
		Keys: []string{
			"initial-small",
			"initial-large",
		},
	}).Execute()
	failTestAtHttpError(err, httpRespGet, t)

	// Verify we got both initial data attributes back
	assertions.Equal(2, len(getResult.Objects), "Should return exactly 2 initial data attributes")

	// Create a map for easier lookup
	retrievedAttrs := make(map[string]iwfidl.EncodedObject)
	for _, attr := range getResult.Objects {
		retrievedAttrs[*attr.Key] = *attr.Value
	}

	// Test 2: Verify initial small data attribute (memo-based)
	smallAttr, exists := retrievedAttrs["initial-small"]
	assertions.True(exists, "Initial small data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Data, *smallAttr.Data, "Initial small data content should match")
	assertions.Equal(*s3GetSetDataAttributes.SmallDataValue.Encoding, *smallAttr.Encoding, "Initial small data encoding should match")

	// Test 3: Verify initial large data attribute (should have actual data loaded from S3)
	largeAttr, exists := retrievedAttrs["initial-large"]
	assertions.True(exists, "Initial large data attribute should exist")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Data, *largeAttr.Data, "Initial large data content should match")
	assertions.Equal(*s3GetSetDataAttributes.LargeDataValue.Encoding, *largeAttr.Encoding, "Initial large data encoding should match")
	// Large data should also preserve external storage references for optimization
	assertions.NotNil(largeAttr.ExtStoreId, "Initial large data should preserve ExtStoreId reference")
	assertions.NotNil(largeAttr.ExtPath, "Initial large data should preserve ExtPath reference")

	// Note: Skip setting additional data attributes test since the workflow has completed.
	// This is expected behavior as our test workflow is designed to complete after the first state.

	// Test 4: Verify S3 objects were created for initial large data attribute
	objectCount, err := globalBlobStore.CountWorkflowObjectsForTesting(context.Background(), wfId)
	assertions.Nil(err)
	assertions.Equal(int64(1), objectCount, "Should have exactly 1 S3 object for initial large data attribute")

	// Complete the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp2, err2 := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err2, httpResp2, t)
}
