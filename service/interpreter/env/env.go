package env

import (
	"github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/config"
	"go.temporal.io/sdk/converter"
)

// this env package is for creating some handy global objects so that workflow worker can make use of,
// to avoid passing throw worker context(esp hard because we have to do it with both Cadence & Temporal)
// Note it's not a good practice to use global so should keep them minimum

var sharedConfig config.Config

var temporalDataConverter converter.DataConverter

var temporalMemoEncryption bool

var unifiedClient client.UnifiedClient

var taskQueue string

func SetSharedEnv(
	config config.Config,
	memoEncryption bool,
	temporalMemoEncryptionDataConverter converter.DataConverter,
	unifiedClient client.UnifiedClient,
	taskQueue string,
) {
	sharedConfig = config
	temporalDataConverter = temporalMemoEncryptionDataConverter
	temporalMemoEncryption = memoEncryption
	unifiedClient = unifiedClient
	taskQueue = taskQueue
}

func GetUnifiedClient() client.UnifiedClient {
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
