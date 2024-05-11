package integ

import (
	"github.com/indeedeng/iwf/service/common/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

func createTestConfig(failAtMemoCompatibility bool, optimizedVersioning *bool) config.Config {
	return config.Config{
		Api: config.ApiConfig{
			Port:                9715,
			MaxWaitSeconds:      10, // use 10 so that we can test it in the waiting test
			OptimizedVersioning: optimizedVersioning,
		},
		Interpreter: config.Interpreter{
			VerboseDebug:              false,
			FailAtMemoIncompatibility: failAtMemoCompatibility,
		},
	}
}
