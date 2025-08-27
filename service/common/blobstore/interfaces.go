package blobstore

import (
	"context"
	"fmt"
	"strings"
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

type BlobStore interface {
	// WriteObject will write to the current active store
	// returns the active storeId
	WriteObject(ctx context.Context, workflowId, data string) (storeId, path string, err error)
	// ReadObject will read from the store by storeId
	ReadObject(ctx context.Context, storeId, path string) (string, error)
	// DeleteObjectPath will delete all the objects of the path
	DeleteObjectPath(ctx context.Context, storeId, path string) error
	// ListObjectPaths will list the paths of yyyymmdd$workflowId
	ListObjectPaths(ctx context.Context, input ListObjectPathsInput) (*ListObjectPathsOutput, error)
	// CountWorkflowObjectsForTesting is for testing ONLY.
	// count the number of S3 objects for this workflowId
	// Limitation:
	//  1. It doesn't count across two days(so expect test to fail if you happen to run the test across day boundary :)
	//  2. Only count less than 1000 objects(because it only make one API call to S3 which return at most 1000 objects)
	CountWorkflowObjectsForTesting(ctx context.Context, workflowId string) (int64, error)
}

type ListObjectPathsInput struct {
	StoreId           string
	StartAfter        string
	ContinuationToken *string
}

type ListObjectPathsOutput struct {
	ContinuationToken *string
	Paths             []string
}
