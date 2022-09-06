package temporalimpl

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type InterpreterWorker struct {
	temporalClient client.Client
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
}

func (iw *InterpreterWorker) Start() {
	w := worker.New(iw.temporalClient, TaskQueue, worker.Options{})

	w.RegisterWorkflow(Interpreter)
	w.RegisterActivity(Activity)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
