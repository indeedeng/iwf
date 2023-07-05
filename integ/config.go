package integ

import (
	"github.com/indeedeng/iwf/service/common/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

func createTestConfig(failAtMemoCompatibility bool) config.Config {
	return config.Config{
		Api: config.ApiConfig{
			Port:           9715,
			MaxWaitSeconds: 10,
		},
		Interpreter: config.Interpreter{
			VerboseDebug:              false,
			FailAtMemoIncompatibility: failAtMemoCompatibility,
		},
	}
}
