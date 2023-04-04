package cadence

import (
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/interpreter"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
	"log"
)

type InterpreterWorker struct {
	service   workflowserviceclient.Interface
	closeFunc func()
	domain    string
	worker    worker.Worker
	tasklist  string
}

func NewInterpreterWorker(service workflowserviceclient.Interface, domain, tasklist string, closeFunc func()) *InterpreterWorker {
	apiAddress := config.GetApiServiceAddress()
	if apiAddress == "" {
		panic("empty api address, must be initialized through config.SetApiServiceAddress()")
	}

	return &InterpreterWorker{
		service:   service,
		domain:    domain,
		tasklist:  tasklist,
		closeFunc: closeFunc,
	}
}

func (iw *InterpreterWorker) Close() {
	iw.closeFunc()
	iw.worker.Stop()
}

func (iw *InterpreterWorker) Start() {
	iw.worker = worker.New(iw.service, iw.domain, iw.tasklist, worker.Options{})

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterActivity(interpreter.StateStart)
	iw.worker.RegisterActivity(interpreter.StateDecide)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
