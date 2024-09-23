package integ

import (
	"github.com/indeedeng/iwf/service/common/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

func createTestConfig(testCfg IwfServiceTestConfig) config.Config {
	return config.Config{
		Api: config.ApiConfig{
			Port:                9715,
			MaxWaitSeconds:      10, // use 10 so that we can test it in the waiting test
			OptimizedVersioning: testCfg.OptimizedVersioning,
		},
		Interpreter: config.Interpreter{
			VerboseDebug:              false,
			FailAtMemoIncompatibility: !testCfg.DisableFailAtMemoIncompatibility,
			InterpreterActivityConfig: config.InterpreterActivityConfig{
				DefaultHeader: testCfg.DefaultHeaders,
			},
		},
	}
}
