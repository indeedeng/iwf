package integ

import (
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"github.com/stretchr/testify/assert"
	"testing"

	"go.temporal.io/sdk/worker"
)

var jsonHistoryFiles = []string{
	"v1-persistence.json",
	"v2-persistence.json",
	"v2-basic.json",
	"v2-basic-disable-system-searchattributes.json",
	"v2-any-timer-signal-continue-as-new.json",
	"v2-any-timer-signal.json",
}

func TestTemporalReplay(t *testing.T) {
	replayer := worker.NewWorkflowReplayer()

	replayer.RegisterWorkflow(temporal.Interpreter)

	for _, f := range jsonHistoryFiles {
		err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history/"+f)
		assertions := assert.New(t)
		assertions.Nil(err)
	}

}

// NOTE: set TEMPORAL_DEBUG=true
//func TestDebugTemporalReplay(t *testing.T) {
//	replayer := worker.NewWorkflowReplayer()
//
//	replayer.RegisterWorkflow(temporal.Interpreter)
//
//	err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history/debug.json")
//	assertions := assert.New(t)
//	assertions.Nil(err)
//}
