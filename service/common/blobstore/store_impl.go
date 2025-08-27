package blobstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/service/common/log"
	"go.temporal.io/sdk/client"
)

type blobStoreImpl struct {
	s3Client                    *s3.Client
	pathPrefix                  string // the Temporal namespace or Cadence domain + "/"
	activeStorage               config.BlobStorageConfig
	supportedStore              map[string]config.BlobStorageConfig // storeId as key
	logger                      log.Logger
	writeObjectErrorCounter     client.MetricsCounter
	readObjectErrorCounter      client.MetricsCounter
	writeObjectSuccessHistogram client.MetricsTimer
	readObjectSuccessHistogram  client.MetricsTimer
}

func NewBlobStore(
	s3Client *s3.Client,
	temporalOrCadenceNamespace string,
	storeConfig config.ExternalStorageConfig,
	logger log.Logger,
	metrics client.MetricsHandler,
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

	metricsHandler := metrics.WithTags(map[string]string{"prefix": temporalOrCadenceNamespace})
	writeObjectErrorCounter := metricsHandler.Counter("write_object_error")
	readObjectErrorCounter := metricsHandler.Counter("read_object_error")
	writeObjectSuccessHistogram := metricsHandler.Timer("write_object_success")
	readObjectSuccessHistogram := metricsHandler.Timer("read_object_success")

	return &blobStoreImpl{
		s3Client:                    s3Client,
		pathPrefix:                  temporalOrCadenceNamespace + "/",
		activeStorage:               *activeStorage,
		supportedStore:              supportedStores,
		logger:                      logger,
		writeObjectErrorCounter:     writeObjectErrorCounter,
		readObjectErrorCounter:      readObjectErrorCounter,
		writeObjectSuccessHistogram: writeObjectSuccessHistogram,
		readObjectSuccessHistogram:  readObjectSuccessHistogram,
	}
}

func (b *blobStoreImpl) WriteObject(ctx context.Context, workflowId, data string) (storeId, path string, err error) {
	storeId = b.activeStorage.StorageId
	randomUuid := uuid.New().String()
	yyyymmdd := time.Now().Format("20060102")
	// yyyymmdd$workflowId/uuid
	// Note: using $ here so that the listing can be much easier to implement for pagination
	path = fmt.Sprintf("%s$%s/%s", yyyymmdd, workflowId, randomUuid)

	err = putObject(ctx, b.s3Client, b.activeStorage.S3Bucket, b.pathPrefix+path, data)
	if err != nil {
		b.writeObjectErrorCounter.Inc(1)
		return
	}
	b.writeObjectSuccessHistogram.Record(time.Duration(len(data)))
	return
}

func (b *blobStoreImpl) ReadObject(ctx context.Context, storeId, path string) (string, error) {
	storeConfig, ok := b.supportedStore[storeId]
	if !ok {
		b.readObjectErrorCounter.Inc(1)
		return "", errors.New("store not found for " + storeId)
	}
	data, err := getObject(ctx, b.s3Client, storeConfig.S3Bucket, b.pathPrefix+path)
	if err != nil {
		b.readObjectErrorCounter.Inc(1)
		return "", err
	}
	b.readObjectSuccessHistogram.Record(time.Duration(len(data)))
	return data, nil
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

func (b *blobStoreImpl) CountWorkflowObjectsForTesting(ctx context.Context, workflowId string) (int64, error) {
	// Create the prefix to match objects for this workflowId for today
	yyyymmdd := time.Now().Format("20060102")
	prefix := fmt.Sprintf("%s%s$%s/", b.pathPrefix, yyyymmdd, workflowId)

	// List objects with the prefix (limited to 1000 objects as documented)
	result, err := b.s3Client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(b.activeStorage.S3Bucket),
		Prefix: aws.String(prefix),
	})
	if err != nil {
		return 0, err
	}

	return int64(len(result.Contents)), nil
}

func (b *blobStoreImpl) DeleteWorkflowObjects(ctx context.Context, storeId, workflowPath string) error {
	//TODO implement me
	panic("implement me")
}

func (b *blobStoreImpl) ListWorkflowPaths(ctx context.Context, input ListObjectPathsInput) (*ListObjectPathsOutput, error) {
	storeConfig, ok := b.supportedStore[input.StoreId]
	if !ok {
		return nil, errors.New("store not found for " + input.StoreId)
	}

	listInput := &s3.ListObjectsV2Input{
		Bucket:    aws.String(storeConfig.S3Bucket),
		Prefix:    aws.String(b.pathPrefix),
		Delimiter: aws.String("/"),
	}

	// Set StartAfter if provided
	if input.StartAfterYyyymmdd != "" {
		startAfter := fmt.Sprintf("%s%s$", b.pathPrefix, input.StartAfterYyyymmdd)
		listInput.StartAfter = aws.String(startAfter)
	}

	// Set continuation token if provided
	if input.ContinuationToken != nil {
		listInput.ContinuationToken = input.ContinuationToken
	}

	result, err := b.s3Client.ListObjectsV2(ctx, listInput)
	if err != nil {
		return nil, err
	}

	// Extract workflow paths from common prefixes
	workflowPaths := make([]string, 0, len(result.CommonPrefixes))
	for _, commonPrefix := range result.CommonPrefixes {
		if commonPrefix.Prefix != nil {
			// Remove the pathPrefix to get the workflow path (yyyymmdd$workflowId)
			prefixStr := *commonPrefix.Prefix
			if strings.HasPrefix(prefixStr, b.pathPrefix) {
				workflowPath := strings.TrimPrefix(prefixStr, b.pathPrefix)
				// Remove trailing "/" if present
				workflowPath = strings.TrimSuffix(workflowPath, "/")
				workflowPaths = append(workflowPaths, workflowPath)
			}
		}
	}

	return &ListObjectPathsOutput{
		ContinuationToken: result.NextContinuationToken,
		WorkflowPaths:     workflowPaths,
	}, nil
}
