package temporal

import (
	"github.com/cadence-oss/iwf-server/service"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type InterpreterWorker struct {
	temporalClient client.Client
	worker         worker.Worker
}

func NewInterpreterWorker(temporalClient client.Client) *InterpreterWorker {
	return &InterpreterWorker{
		temporalClient: temporalClient,
	}
}

func (iw *InterpreterWorker) Close() {
	iw.temporalClient.Close()
	iw.worker.Stop()
}

func (iw *InterpreterWorker) Start() {
	iw.worker = worker.New(iw.temporalClient, service.TaskQueue, worker.Options{})

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterActivity(StateStartActivity)
	iw.worker.RegisterActivity(StateDecideActivity)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
