package temporal

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type InterpreterWorker struct {
	temporalClient client.Client
	worker         worker.Worker
}

func NewInterpreterWorker() *InterpreterWorker {
	// TODO use config for connection options and merge with api handler
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return &InterpreterWorker{
		temporalClient: temporalClient,
	}
}

func (iw *InterpreterWorker) Close() {
	iw.temporalClient.Close()
	iw.worker.Stop()
}

func (iw *InterpreterWorker) Start() {
	iw.worker = worker.New(iw.temporalClient, TaskQueue, worker.Options{})

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterActivity(Activity)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
