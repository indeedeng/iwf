package replayTests

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
	"v2-any-timer-signal.json",
	"v3-any-timer-signal-continue-as-new.json",
	"v3-basic.json",
	"v3-skip-start.json",
	"v3-bug-no-state-stuck.json",
	"v4-continue-as-new-on-no-state.json",
	"v4-continued-as-new-before-versioning-optimization.json",
	"v4-local-activity-optimization.json",
	"v5-basic.json",
	"v6-search-attributes-optimization-enabled-for-all.json",
	"v6-search-attributes-optimization-default.json",
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
