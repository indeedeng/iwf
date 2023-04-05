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

func NewInterpreterWorker(config config.Config, service workflowserviceclient.Interface, domain, tasklist string, closeFunc func()) *InterpreterWorker {
	interpreter.SetSharedConfig(config)
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
	config := interpreter.GetSharedConfig()
	options := worker.Options{
		MaxConcurrentActivityTaskPollers: 10,
		MaxConcurrentDecisionTaskPollers: 10,
	}
	if config.Interpreter.Cadence != nil && config.Interpreter.Cadence.WorkerOptions != nil {
		options = *config.Interpreter.Cadence.WorkerOptions
	}
	iw.worker = worker.New(iw.service, iw.domain, iw.tasklist, options)
	worker.EnableVerboseLogging(config.Interpreter.VerboseDebug)

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterActivity(interpreter.StateStart)
	iw.worker.RegisterActivity(interpreter.StateDecide)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
