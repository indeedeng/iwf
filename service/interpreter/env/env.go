package env

import (
	"github.com/indeedeng/iwf/config"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/blobstore"
	"go.temporal.io/sdk/converter"
)

// this env package is for creating some handy global objects so that workflow worker can make use of,
// to avoid passing throw worker context(esp hard because we have to do it with both Cadence & Temporal)
// Note it's not a good practice to use global so should keep them minimum

var sharedConfig config.Config

var temporalDataConverter converter.DataConverter

var temporalMemoEncryption bool

var unifiedClient uclient.UnifiedClient

var taskQueue string

var blobStore blobstore.BlobStore

func SetSharedEnv(
	config config.Config,
	memoEncryption bool,
	temporalMemoEncryptionDataConverter converter.DataConverter,
	client uclient.UnifiedClient,
	queue string,
	store blobstore.BlobStore,
) {
	sharedConfig = config
	temporalDataConverter = temporalMemoEncryptionDataConverter
	temporalMemoEncryption = memoEncryption
	unifiedClient = client
	taskQueue = queue
	blobStore = store
}

func GetUnifiedClient() uclient.UnifiedClient {
	return unifiedClient
}

func GetTaskQueue() string {
	return taskQueue
}

func GetSharedConfig() config.Config {
	return sharedConfig
}

func CheckAndGetTemporalMemoEncryptionDataConverter() (converter.DataConverter, bool) {
	return temporalDataConverter, temporalMemoEncryption
}

func GetBlobStore() blobstore.BlobStore {
	return blobStore
}
