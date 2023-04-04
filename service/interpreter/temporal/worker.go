package temporal

import (
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

type InterpreterWorker struct {
	temporalClient client.Client
	worker         worker.Worker
	taskQueue      string
}

func NewInterpreterWorker(temporalClient client.Client, taskQueue string) *InterpreterWorker {
	apiAddress := config.GetApiServiceAddress()
	if apiAddress == "" {
		panic("empty api address, must be initialized through config.SetApiServiceAddress()")
	}

	return &InterpreterWorker{
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

func (iw *InterpreterWorker) Close() {
	iw.temporalClient.Close()
	iw.worker.Stop()
}

func (iw *InterpreterWorker) Start() {
	iw.worker = worker.New(iw.temporalClient, iw.taskQueue, worker.Options{})

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterActivity(interpreter.StateStart)
	iw.worker.RegisterActivity(interpreter.StateDecide)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
