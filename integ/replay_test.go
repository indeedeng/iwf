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
