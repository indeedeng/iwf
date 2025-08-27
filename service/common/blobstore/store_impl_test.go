package blobstore

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/service/common/log/loggerimpl"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
)

const (
	testBucket    = "iwf-test-bucket"
	testRegion    = "us-east-1"
	testEndpoint  = "http://localhost:9000"
	testAccessKey = "minioadmin"
	testSecretKey = "minioadmin"
	testNamespace = "default"
	testStorageId = "test-storage-id"
)

func createTestBlobStore(t *testing.T) BlobStore {
	// Create S3 client for MinIO
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(testRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(testAccessKey, testSecretKey, "")),
		awsconfig.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: testEndpoint,
				}, nil
			})),
	)
	assert.NoError(t, err)

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// Create bucket if it doesn't exist
	_, err = s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(testBucket),
	})
	if err != nil {
		// Bucket doesn't exist, create it
		_, err = s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String(testBucket),
		})
		assert.NoError(t, err)
	}

	// Create test configuration
	storeConfig := config.ExternalStorageConfig{
		Enabled:          true,
		ThresholdInBytes: 100,
		SupportedStorages: []config.BlobStorageConfig{
			{
				Status:      config.StorageStatusActive,
				StorageId:   testStorageId,
				StorageType: config.StorageTypeS3,
				S3Endpoint:  testEndpoint,
				S3Bucket:    testBucket,
				S3Region:    testRegion,
				S3AccessKey: testAccessKey,
				S3SecretKey: testSecretKey,
			},
		},
	}

	logger, err := loggerimpl.NewDevelopment()
	assert.NoError(t, err)
	blobStore := NewBlobStore(s3Client, testNamespace, storeConfig, logger, client.MetricsNopHandler)
	assert.NotNil(t, blobStore)

	return blobStore
}

func TestBlobStoreIntegration(t *testing.T) {
	// Skip if not running integration tests or MinIO is not available
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	blobStore := createTestBlobStore(t)
	ctx := context.Background()

	// Test data
	workflowId1 := "test-workflow-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	workflowId2 := "test-workflow-" + strconv.FormatInt(time.Now().UnixNano()+1, 10)
	testData := "test data content"

	t.Run("WriteAndCountObjects", func(t *testing.T) {
		// Initial count should be 0
		count, err := blobStore.CountWorkflowObjectsForTesting(ctx, workflowId1)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)

		// Write first object
		storeId1, path1, err := blobStore.WriteObject(ctx, workflowId1, testData)
		assert.NoError(t, err)
		assert.Equal(t, testStorageId, storeId1)
		assert.NotEmpty(t, path1)

		// Count should be 1
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, workflowId1)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Write second object for same workflow
		storeId2, path2, err := blobStore.WriteObject(ctx, workflowId1, testData+"2")
		assert.NoError(t, err)
		assert.Equal(t, testStorageId, storeId2)
		assert.NotEmpty(t, path2)
		assert.NotEqual(t, path1, path2)

		// Count should be 2
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, workflowId1)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// Write object for different workflow
		storeId3, path3, err := blobStore.WriteObject(ctx, workflowId2, testData+"3")
		assert.NoError(t, err)
		assert.Equal(t, testStorageId, storeId3)
		assert.NotEmpty(t, path3)

		// Count for first workflow should still be 2
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, workflowId1)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), count)

		// Count for second workflow should be 1
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, workflowId2)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("ReadObjects", func(t *testing.T) {
		// Write an object
		storeId, path, err := blobStore.WriteObject(ctx, workflowId1, testData)
		assert.NoError(t, err)

		// Read it back
		retrievedData, err := blobStore.ReadObject(ctx, storeId, path)
		assert.NoError(t, err)
		assert.Equal(t, testData, retrievedData)
	})

	t.Run("ListWorkflowPaths", func(t *testing.T) {
		// Write objects for multiple workflows
		_, _, err := blobStore.WriteObject(ctx, workflowId1, testData)
		assert.NoError(t, err)
		_, _, err = blobStore.WriteObject(ctx, workflowId2, testData)
		assert.NoError(t, err)

		// List workflow paths
		input := ListObjectPathsInput{
			StoreId: testStorageId,
		}
		output, err := blobStore.ListWorkflowPaths(ctx, input)
		assert.NoError(t, err)
		assert.NotNil(t, output)

		// Should contain both workflow paths
		assert.True(t, len(output.WorkflowPaths) >= 2)

		// Verify workflow paths contain expected patterns
		todayPrefix := time.Now().Format("20060102")
		expectedPath1 := fmt.Sprintf("%s$%s", todayPrefix, workflowId1)
		expectedPath2 := fmt.Sprintf("%s$%s", todayPrefix, workflowId2)

		foundPath1 := false
		foundPath2 := false
		for _, path := range output.WorkflowPaths {
			if path == expectedPath1 {
				foundPath1 = true
			}
			if path == expectedPath2 {
				foundPath2 = true
			}
		}
		assert.True(t, foundPath1, "Expected path for workflowId1 not found")
		assert.True(t, foundPath2, "Expected path for workflowId2 not found")
	})

	t.Run("DeleteWorkflowObjects", func(t *testing.T) {
		// Create a new workflow ID for this test
		deleteTestWorkflowId := "delete-test-workflow-" + strconv.FormatInt(time.Now().UnixNano(), 10)

		// Write multiple objects for the workflow
		_, _, err := blobStore.WriteObject(ctx, deleteTestWorkflowId, "test data 1")
		assert.NoError(t, err)
		_, _, err = blobStore.WriteObject(ctx, deleteTestWorkflowId, "test data 2")
		assert.NoError(t, err)
		_, _, err = blobStore.WriteObject(ctx, deleteTestWorkflowId, "test data 3")
		assert.NoError(t, err)

		// Verify objects exist
		count, err := blobStore.CountWorkflowObjectsForTesting(ctx, deleteTestWorkflowId)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)

		// Delete all objects for the workflow
		todayPrefix := time.Now().Format("20060102")
		workflowPath := fmt.Sprintf("%s$%s", todayPrefix, deleteTestWorkflowId)
		err = blobStore.DeleteWorkflowObjects(ctx, testStorageId, workflowPath)
		assert.NoError(t, err)

		// Verify objects are deleted
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, deleteTestWorkflowId)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("DeleteWorkflowObjectsMultiple", func(t *testing.T) {
		// Create a workflow with more objects to test pagination handling
		multiDeleteWorkflowId := "multi-delete-workflow-" + strconv.FormatInt(time.Now().UnixNano(), 10)

		// Create 10 objects for this workflow
		numObjects := 10
		for i := 0; i < numObjects; i++ {
			data := fmt.Sprintf("test data %d", i)
			_, _, err := blobStore.WriteObject(ctx, multiDeleteWorkflowId, data)
			assert.NoError(t, err)
		}

		// Verify all objects exist
		count, err := blobStore.CountWorkflowObjectsForTesting(ctx, multiDeleteWorkflowId)
		assert.NoError(t, err)
		assert.Equal(t, int64(numObjects), count)

		// Delete all objects for the workflow
		todayPrefix := time.Now().Format("20060102")
		workflowPath := fmt.Sprintf("%s$%s", todayPrefix, multiDeleteWorkflowId)
		err = blobStore.DeleteWorkflowObjects(ctx, testStorageId, workflowPath)
		assert.NoError(t, err)

		// Verify all objects are deleted
		count, err = blobStore.CountWorkflowObjectsForTesting(ctx, multiDeleteWorkflowId)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("DeleteWorkflowObjectsNonExistent", func(t *testing.T) {
		// Try to delete objects for a workflow that doesn't exist
		todayPrefix := time.Now().Format("20060102")
		workflowPath := fmt.Sprintf("%s$%s", todayPrefix, "non-existent-workflow")
		err := blobStore.DeleteWorkflowObjects(ctx, testStorageId, workflowPath)
		assert.NoError(t, err) // Should succeed even if no objects to delete
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test with invalid store ID
		input := ListObjectPathsInput{
			StoreId: "invalid-store-id",
		}
		_, err := blobStore.ListWorkflowPaths(ctx, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store not found")

		// Test reading with invalid store ID
		_, err = blobStore.ReadObject(ctx, "invalid-store-id", "some-path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store not found")

		// Test delete with invalid store ID
		err = blobStore.DeleteWorkflowObjects(ctx, "invalid-store-id", "some-workflow-path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store not found")
	})
}
