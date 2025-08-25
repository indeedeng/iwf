package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// TODO remove this file when implementation of S3 is done

const (
	bucketName = "iwf-test-bucket"
	region     = "us-east-1"
	endpoint   = "http://localhost:9000"
	accessKey  = "minioadmin"
	secretKey  = "minioadmin"
)

func main() {
	ctx := context.Background()

	// Create S3 client configured for MinIO
	client, err := createS3Client(ctx)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Create bucket if it doesn't exist
	if err := createBucketIfNotExists(ctx, client); err != nil {
		log.Fatalf("Failed to create bucket: %v", err)
	}

	// Demonstrate S3 operations
	if err := demonstrateS3Operations(ctx, client); err != nil {
		log.Fatalf("Failed to demonstrate S3 operations: %v", err)
	}

	fmt.Println("S3 operations completed successfully!")
}

func createS3Client(ctx context.Context) (*s3.Client, error) {
	// Create custom resolver for MinIO endpoint
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID {
			return aws.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
				Source:            aws.EndpointSourceCustom,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})

	// Load AWS config with custom credentials and endpoint
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with path-style addressing (required for MinIO)
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}

func createBucketIfNotExists(ctx context.Context, client *s3.Client) error {
	// Check if bucket exists
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		// Bucket doesn't exist, create it
		fmt.Printf("Creating bucket: %s\n", bucketName)
		_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		fmt.Printf("Bucket created successfully: %s\n", bucketName)
	} else {
		fmt.Printf("Bucket already exists: %s\n", bucketName)
	}

	return nil
}

func demonstrateS3Operations(ctx context.Context, client *s3.Client) error {
	// Generate workflow objects with path pattern: /createdDayStamp/namespace/workflowType$workflowId/uuid
	workflowObjects := generateWorkflowObjects()

	// 1. Write workflow objects to S3
	fmt.Println("\n=== Writing workflow objects to S3 ===")
	for key, content := range workflowObjects {
		if err := putObject(ctx, client, key, content); err != nil {
			return fmt.Errorf("failed to put object %s: %w", key, err)
		}
		fmt.Printf("✓ Uploaded: %s\n", key)
	}

	// 2. Verify objects using GET API
	fmt.Println("\n=== Verifying objects using GET API ===")
	for key := range workflowObjects {
		content, err := getObject(ctx, client, key)
		if err != nil {
			return fmt.Errorf("failed to get object %s: %w", key, err)
		}
		fmt.Printf("✓ Verified: %s (size: %d bytes)\n", key, len(content))

		// Show content preview
		preview := content
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		fmt.Printf("  Content: %s\n", strings.ReplaceAll(preview, "\n", "\\n"))
	}

	// 3. List all objects to show the full structure
	fmt.Println("\n=== Listing all workflow objects ===")
	objects, err := listObjects(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to list objects: %w", err)
	}

	fmt.Printf("Found %d workflow objects in bucket '%s':\n", len(objects), bucketName)
	for _, obj := range objects {
		fmt.Printf("  - %s (size: %d bytes, modified: %s)\n",
			aws.ToString(obj.Key),
			obj.Size,
			obj.LastModified.Format(time.RFC3339))
	}

	// 4. Demonstrate hierarchical listing with delimiter
	fmt.Println("\n=== Listing day stamps (top level) ===")
	dayStamps, err := listWithDelimiter(ctx, client, "", "/")
	if err != nil {
		return fmt.Errorf("failed to list day stamps: %w", err)
	}

	fmt.Printf("Found %d day stamp directories:\n", len(dayStamps))
	for _, prefix := range dayStamps {
		fmt.Printf("  - %s\n", prefix)
	}

	// 5. For each day stamp, list namespaces
	for _, dayStamp := range dayStamps {
		fmt.Printf("\n=== Listing namespaces under %s ===\n", dayStamp)
		namespaces, err := listWithDelimiter(ctx, client, dayStamp, "/")
		if err != nil {
			return fmt.Errorf("failed to list namespaces under %s: %w", dayStamp, err)
		}

		fmt.Printf("Found %d namespace directories under %s:\n", len(namespaces), dayStamp)
		for _, namespace := range namespaces {
			fmt.Printf("  - %s\n", namespace)
		}

		// 6. For each namespace, demonstrate paginated listing of workflows
		for _, namespace := range namespaces {
			fmt.Printf("\n=== Paginated listing of workflows under %s ===\n", namespace)

			// Demonstrate pagination with page size of 2
			pageSize := 2
			err := listWorkflowsPaginated(ctx, client, namespace, pageSize)
			if err != nil {
				return fmt.Errorf("failed to paginate workflows under %s: %w", namespace, err)
			}

			// Demonstrate filtered pagination
			fmt.Printf("\n=== Filtered pagination example ===\n")
			err = listWorkflowsPaginatedWithFilter(ctx, client, namespace, "OrderProcessing", 3)
			if err != nil {
				return fmt.Errorf("failed to paginate filtered workflows under %s: %w", namespace, err)
			}
		}
	}

	return nil
}

func generateWorkflowObjects() map[string]string {
	// Generate different day stamps (days since Unix epoch)
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	workflowObjects := make(map[string]string)

	// Create objects across several different days
	days := []int{0, 1, 2, 5, 10} // Different day offsets
	workflowTypes := []string{"OrderProcessing", "PaymentFlow", "InventoryUpdate", "UserRegistration"}

	for _, dayOffset := range days {
		// Calculate day stamp (days since Unix epoch)
		currentDay := baseTime.AddDate(0, 0, dayOffset)
		dayStamp := int(currentDay.Unix() / 86400) // 86400 seconds in a day

		// Create multiple workflows for each day
		for i, workflowType := range workflowTypes {
			workflowId := fmt.Sprintf("wf-%d-%d", dayOffset, i+1)

			// Create multiple objects per workflow (different UUIDs)
			numObjects := rand.Intn(3) + 1 // 1-3 objects per workflow
			for j := 0; j < numObjects; j++ {
				objectUuid := uuid.New().String()

				// Path pattern: /createdDayStamp/namespace/workflowType$workflowId/uuid
				key := fmt.Sprintf("%d/testNamespace/%s$%s/%s",
					dayStamp,
					workflowType,
					workflowId,
					objectUuid)

				// Create meaningful content for each object
				content := fmt.Sprintf(`{
  "createdDayStamp": %d,
  "namespace": "testNamespace",
  "workflowType": "%s",
  "workflowId": "%s",
  "objectId": "%s",
  "createdAt": "%s",
  "data": {
    "step": %d,
    "status": "completed",
    "processingTime": "%dms"
  }
}`, dayStamp, workflowType, workflowId, objectUuid,
					currentDay.Format(time.RFC3339), j+1, rand.Intn(1000)+100)

				workflowObjects[key] = content
			}
		}
	}

	return workflowObjects
}

func listWithDelimiter(ctx context.Context, client *s3.Client, prefix, delimiter string) ([]string, error) {
	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String(delimiter),
	})
	if err != nil {
		return nil, err
	}

	var prefixes []string
	for _, commonPrefix := range result.CommonPrefixes {
		prefixes = append(prefixes, aws.ToString(commonPrefix.Prefix))
	}

	return prefixes, nil
}

func listWorkflowsPaginated(ctx context.Context, client *s3.Client, namespacePrefix string, pageSize int) error {
	var continuationToken *string
	pageNumber := 1
	totalWorkflows := 0

	fmt.Printf("📄 Paginating workflows with page size: %d\n", pageSize)

	for {
		fmt.Printf("\n--- Page %d ---\n", pageNumber)

		// Create input for ListObjectsV2 with pagination parameters
		input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucketName),
			Prefix:    aws.String(namespacePrefix),
			Delimiter: aws.String("/"),
			MaxKeys:   aws.Int32(int32(pageSize)), // Control page size
		}

		// Add continuation token for subsequent pages
		if continuationToken != nil {
			input.ContinuationToken = continuationToken
		}

		// Execute the request
		result, err := client.ListObjectsV2(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to list workflows page %d: %w", pageNumber, err)
		}

		// Display workflows in this page
		workflowsInPage := len(result.CommonPrefixes)
		fmt.Printf("Workflows in page %d (%d items):\n", pageNumber, workflowsInPage)

		for i, commonPrefix := range result.CommonPrefixes {
			workflowPath := aws.ToString(commonPrefix.Prefix)
			fmt.Printf("  %d. %s\n", i+1, workflowPath)

			// Extract just the workflow name from the path for cleaner display
			pathParts := strings.Split(strings.TrimSuffix(workflowPath, "/"), "/")
			if len(pathParts) >= 3 {
				workflowName := pathParts[2] // workflowType$workflowId
				fmt.Printf("     → Workflow: %s\n", workflowName)
			}
		}

		totalWorkflows += workflowsInPage

		// Check if there are more pages
		if !aws.ToBool(result.IsTruncated) {
			fmt.Printf("\n✅ End of results reached. Total workflows found: %d across %d pages\n", totalWorkflows, pageNumber)
			break
		}

		// Prepare for next page
		continuationToken = result.NextContinuationToken
		pageNumber++

		fmt.Printf("🔄 More results available. Continuation token: %s\n", aws.ToString(continuationToken)[:20]+"...")
	}

	return nil
}

// Advanced pagination example with custom filtering
func listWorkflowsPaginatedWithFilter(ctx context.Context, client *s3.Client, namespacePrefix, workflowTypeFilter string, pageSize int) error {
	var continuationToken *string
	pageNumber := 1
	totalMatched := 0

	fmt.Printf("📄 Paginating workflows filtered by type '%s' with page size: %d\n", workflowTypeFilter, pageSize)

	for {
		fmt.Printf("\n--- Page %d (Filtered) ---\n", pageNumber)

		input := &s3.ListObjectsV2Input{
			Bucket:    aws.String(bucketName),
			Prefix:    aws.String(namespacePrefix),
			Delimiter: aws.String("/"),
			MaxKeys:   aws.Int32(int32(pageSize)),
		}

		if continuationToken != nil {
			input.ContinuationToken = continuationToken
		}

		result, err := client.ListObjectsV2(ctx, input)
		if err != nil {
			return fmt.Errorf("failed to list filtered workflows page %d: %w", pageNumber, err)
		}

		// Filter workflows by type
		var matchedWorkflows []string
		for _, commonPrefix := range result.CommonPrefixes {
			workflowPath := aws.ToString(commonPrefix.Prefix)
			if strings.Contains(workflowPath, workflowTypeFilter+"$") {
				matchedWorkflows = append(matchedWorkflows, workflowPath)
			}
		}

		fmt.Printf("Matching workflows in page %d (%d/%d items):\n", pageNumber, len(matchedWorkflows), len(result.CommonPrefixes))
		for i, workflow := range matchedWorkflows {
			fmt.Printf("  %d. %s\n", i+1, workflow)
		}

		totalMatched += len(matchedWorkflows)

		if !aws.ToBool(result.IsTruncated) {
			fmt.Printf("\n✅ End of results. Total matching workflows: %d across %d pages\n", totalMatched, pageNumber)
			break
		}

		continuationToken = result.NextContinuationToken
		pageNumber++
	}

	return nil
}

func putObject(ctx context.Context, client *s3.Client, key, content string) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   strings.NewReader(content),
		//ContentType: aws.String("json"),
	})
	return err
}

func getObject(ctx context.Context, client *s3.Client, key string) (string, error) {
	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return "", err
	}
	defer func() { _ = result.Body.Close() }()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, result.Body)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func listObjects(ctx context.Context, client *s3.Client) ([]types.Object, error) {
	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	return result.Contents, nil
}

//func listObjectsWithPrefix(ctx context.Context, client *s3.Client, prefix string) ([]types.Object, error) {
//	result, err := client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
//		Bucket: aws.String(bucketName),
//		Prefix: aws.String(prefix),
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	return result.Contents, nil
//}
