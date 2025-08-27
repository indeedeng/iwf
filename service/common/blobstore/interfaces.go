package blobstore

import (
	"context"
	"fmt"
	"strings"
	"time"
)

var reservedCharacters = []string{"/", "$"}

func ValidateWorkflowId(workflowId string) error {
	for _, reservedCharacter := range reservedCharacters {
		if strings.Contains(workflowId, reservedCharacter) {
			return fmt.Errorf("workflowId contains reserved character: %s", reservedCharacter)
		}
	}
	return nil
}

func MustExtractWorkflowId(workflowPath string) string {
	workflowId, err := ExtractWorkflowId(workflowPath)
	if err != nil {
		panic(err)
	}
	return workflowId
}

func ExtractWorkflowId(workflowPath string) (string, error) {
	parts := strings.Split(workflowPath, "$")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid workflow path: %s", workflowPath)
	}
	return parts[1], nil
}

func MustExtractYyyymmddToUnixSeconds(workflowPath string) int64 {
	// yyyymmdd$workflowId
	yyyymmdd, err := ExtractYyyymmdd(workflowPath)
	if err != nil {
		panic(err)
	}
	parsedTime, err := time.Parse("20060102", yyyymmdd)
	if err != nil {
		panic(err)
	}
	return parsedTime.Unix()
}

func ExtractYyyymmdd(workflowPath string) (string, error) {
	parts := strings.Split(workflowPath, "$")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid workflow path: %s", workflowPath)
	}
	return parts[0], nil
}

type BlobStore interface {
	// WriteObject will write to the current active store
	// returns the active storeId
	// The final path pattern is pathPrefix + yyyymmdd$workflowId/uuid
	// But the returned path doesn't include pathPrefix, only yymmdd$workflowId/uuid
	WriteObject(ctx context.Context, workflowId, data string) (storeId, path string, err error)
	// ReadObject will read from the store by storeId and path
	// The path should be the one returned from WriteObject, in format of yyyymmdd$workflowId/uuid
	ReadObject(ctx context.Context, storeId, path string) (string, error)
	// DeleteWorkflowObjects will delete all the objects of the workflowId
	// workflowPath is yyyymmdd$workflowId, where yymmdd is needed to compose the path
	DeleteWorkflowObjects(ctx context.Context, storeId, workflowPath string) error
	// ListWorkflowPaths will list the workflowPaths ( yyyymmdd$workflowId ) as CommonPrefixes from S3
	// It uses of delimiter "/" before the uuid to get all the CommonPrefixes
	// StartAfterYyyymmdd is the yyyymmdd to exclude the date when listing
	ListWorkflowPaths(ctx context.Context, input ListObjectPathsInput) (*ListObjectPathsOutput, error)
	// CountWorkflowObjectsForTesting is for testing ONLY.
	// count the number of S3 objects for this workflowId
	// Limitation:
	//  1. It doesn't count across two days(so expect test to fail if you happen to run the test across day boundary :)
	//  2. Only count less than 1000 objects(because it only make one API call to S3 which return at most 1000 objects)
	CountWorkflowObjectsForTesting(ctx context.Context, workflowId string) (int64, error)
}

type ListObjectPathsInput struct {
	StoreId           string
	ContinuationToken *string
}

type ListObjectPathsOutput struct {
	ContinuationToken *string
	WorkflowPaths     []string
}
