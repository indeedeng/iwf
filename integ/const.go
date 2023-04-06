package integ

import (
	"github.com/indeedeng/iwf/service/common/config"
)

const testWorkflowServerPort = "9714"
const testIwfServerPort = "9715"

var testConfig = config.Config{
	Api: config.ApiConfig{
		Port: 9715,
	},
	Interpreter: config.Interpreter{
		VerboseDebug: false,
	},
}
