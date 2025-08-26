package blobstore

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/service/common/log"
	"io"
	"strings"
)

type blobStoreImpl struct {
	s3Client       *s3.Client
	pathPrefix     string // the Temporal namespace or Cadence domain + "/"
	activeStorage  config.BlobStorageConfig
	supportedStore map[string]config.BlobStorageConfig // storeId as key
	logger         log.Logger
}

func NewBlobStore(
	s3Client *s3.Client,
	temporalOrCadenceNamespace string,
	storeConfig config.ExternalStorageConfig,
	logger log.Logger,
) BlobStore {
	if !storeConfig.Enabled {
		return nil
	}

	var activeStorage *config.BlobStorageConfig
	supportedStores := map[string]config.BlobStorageConfig{}
	for _, storage := range storeConfig.SupportedStorages {
		if storage.Status == config.StorageStatusActive {
			if activeStorage != nil {
				panic("cannot have more than one active storage configured")
			}
			activeStorage = &storage
		}
		supportedStores[storage.StorageId] = storage
		if storage.StorageType != config.StorageTypeS3 {
			panic("only S3 storage type is supported")
		}
	}
	if activeStorage == nil {
		panic("no active storage found")
	}

	return &blobStoreImpl{
		s3Client:       s3Client,
		pathPrefix:     temporalOrCadenceNamespace + "/",
		activeStorage:  *activeStorage,
		supportedStore: supportedStores,
		logger:         logger,
	}
}

func (b blobStoreImpl) WriteObject(ctx context.Context, path, data string) (string, error) {
	err := putObject(ctx, b.s3Client, b.activeStorage.S3Bucket, b.pathPrefix+path, data)
	if err != nil {
		return "", err
	}
	return b.activeStorage.StorageId, nil
}

func (b blobStoreImpl) ReadObject(ctx context.Context, storeId, path string) (string, error) {
	storeConfig, ok := b.supportedStore[storeId]
	if !ok {
		return "", errors.New("store not found for " + storeId)
	}
	return getObject(ctx, b.s3Client, storeConfig.S3Bucket, b.pathPrefix+path)
}

func putObject(ctx context.Context, client *s3.Client, bucketName string, key, content string) error {
	_, err := client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        strings.NewReader(content),
		ContentType: aws.String("application/json"),
	})
	return err
}

func getObject(ctx context.Context, client *s3.Client, bucketName, key string) (string, error) {
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
