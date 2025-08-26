package blobstore

import "context"

type BlobStore interface {
	// WriteObject will write to the current active store
	// returns the active storeId
	WriteObject(ctx context.Context, path, data string) (storeId string, err error)
	// ReadObject will read from the store by storeId
	ReadObject(ctx context.Context, storeId, path string) (string, error)
}
