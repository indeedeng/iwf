package integ

import (
	"github.com/indeedeng/iwf/service/common/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

const TestHeaderKey = "x-envoy-upstream-rq-timeout-ms"
const TestHeaderValue = "86400"

func createTestConfig(testCfg IwfServiceTestConfig) config.Config {
	failAtMemoIncompatibility := !testCfg.DisableFailAtMemoIncompatibility
	cfg := config.Config{
		Api: config.ApiConfig{
			Port:           9715,
			MaxWaitSeconds: 10, // use 10 so that we can test it in the waiting test
		},
		Interpreter: config.Interpreter{
			VerboseDebug:              false,
			FailAtMemoIncompatibility: failAtMemoIncompatibility,
			InterpreterActivityConfig: config.InterpreterActivityConfig{},
		},
	}
	if testCfg.SetTestHeader {
		cfg.Interpreter.InterpreterActivityConfig.DefaultHeader = map[string]string{
			TestHeaderKey: TestHeaderValue,
		}
	}
	return cfg
}
