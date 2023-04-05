package interpreter

import "github.com/indeedeng/iwf/service/common/config"

var sharedConfig config.Config

func SetSharedConfig(config config.Config) {
	sharedConfig = config
}

func GetSharedConfig() config.Config {
	return sharedConfig
}
