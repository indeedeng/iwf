package blobstore

import "context"

type BlobStore interface {
	// WriteObject will write to the current active store
	// returns the active storeId
	WriteObject(ctx context.Context, workflowId, data string) (storeId, path string, err error)
	// ReadObject will read from the store by storeId
	ReadObject(ctx context.Context, storeId, path string) (string, error)
	// DeleteObjectPath will delete all the objects of the path
	DeleteObjectPath(ctx context.Context, storeId, path string) error
	// ListObjectPaths will list the paths, delimiter by "/"
	ListObjectPaths(ctx context.Context, input ListObjectPathsInput) (*ListObjectPathsOutput, error)
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
