package replayTests

import (
	"testing"

	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"github.com/stretchr/testify/assert"

	"go.temporal.io/sdk/worker"
)

var jsonHistoryFiles = []string{
	"eval.json",
}

func TestTemporalReplay(t *testing.T) {
	worker.EnableVerboseLogging(true)

	replayer, err := worker.NewWorkflowReplayerWithOptions(
		worker.WorkflowReplayerOptions{
			EnableLoggingInReplay: true,
		})

	if err != nil {
		panic(err)
	}

	replayer.RegisterWorkflow(temporal.Interpreter)

	for _, f := range jsonHistoryFiles {
		err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history/"+f)
		assertions := assert.New(t)
		assertions.Nil(err, "fail at replay history for: "+f)
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
