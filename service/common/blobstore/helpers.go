package blobstore

import (
	"context"

	"github.com/indeedeng/iwf/gen/iwfidl"
)

func WriteDataObjectsToExternalStorage(ctx context.Context, dataObjects []iwfidl.KeyValue, workflowId string, threashold int, blobStore BlobStore, isExternalStorageEnabled bool) error {
	if !isExternalStorageEnabled {
		return nil
	}

	for i := range dataObjects {
		if dataObjects[i].Value != nil && dataObjects[i].Value.Data != nil &&
			len(*dataObjects[i].Value.Data) > threashold {
			// Save data to external storage
			storeId, path, writeErr := blobStore.WriteObject(ctx, workflowId, *dataObjects[i].Value.Data)
			if writeErr != nil {
				return writeErr
			}
			dataObjects[i].Value.ExtStoreId = &storeId
			dataObjects[i].Value.ExtPath = &path
			dataObjects[i].Value.Data = nil // Clear data since it's now in external storage
		}
	}
	return nil
}

func LoadDataObjectsFromExternalStorage(ctx context.Context, dataObjects []iwfidl.KeyValue, blobStore BlobStore) error {
	for i := range dataObjects {
		if dataObjects[i].Value != nil && dataObjects[i].Value.ExtStoreId != nil && dataObjects[i].Value.ExtPath != nil {
			data, err := blobStore.ReadObject(ctx, *dataObjects[i].Value.ExtStoreId, *dataObjects[i].Value.ExtPath)
			if err != nil {
				return err
			}

			dataObjects[i].Value.Data = &data
			dataObjects[i].Value.ExtPath = nil
			dataObjects[i].Value.ExtStoreId = nil
		}
	}
	return nil
}
