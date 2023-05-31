package env

import (
	"github.com/indeedeng/iwf/service/common/config"
	"go.temporal.io/sdk/converter"
)

// this env package is for creating some handy global objects so that workflow worker can make use of,
// to avoid passing throw worker context(esp hard because we have to do it with both Cadence & Temporal)
// Note it's not a good practice to use global so should keep them minimum

var sharedConfig config.Config

var temporalDataConverter converter.DataConverter

func SetSharedEnv(config config.Config, temporalMemoEncryptionDataConverter converter.DataConverter) {
	sharedConfig = config
	temporalDataConverter = temporalMemoEncryptionDataConverter
}

func GetSharedConfig() config.Config {
	return sharedConfig
}

func CheckAndGetTemporalMemoEncryptionDataConverter() (converter.DataConverter, bool) {
	exists := temporalDataConverter != nil
	return temporalDataConverter, exists
}
