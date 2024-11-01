package cadence

import (
	"github.com/indeedeng/iwf/config"
	"log"

	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/interpreter"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"
)

type InterpreterWorker struct {
	service   workflowserviceclient.Interface
	closeFunc func()
	domain    string
	worker    worker.Worker
	tasklist  string
}

type StartOptions struct {
	DisableStickyCache bool
}

func NewInterpreterWorker(
	config config.Config, service workflowserviceclient.Interface, domain, tasklist string, closeFunc func(),
	unifiedClient uclient.UnifiedClient,
) *InterpreterWorker {
	env.SetSharedEnv(config, false, nil, unifiedClient, tasklist)
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
	var options StartOptions

	// default options
	options.DisableStickyCache = false

	iw.StartWithOptions(options)
}

func (iw *InterpreterWorker) StartWithOptions(startOptions StartOptions) {
	config := env.GetSharedConfig()
	options := worker.Options{
		MaxConcurrentActivityTaskPollers: 10,
		MaxConcurrentDecisionTaskPollers: 10,
	}

	if startOptions.DisableStickyCache {
		options.DisableStickyExecution = true
	}

	if config.Interpreter.Cadence != nil && config.Interpreter.Cadence.WorkerOptions != nil {
		options = *config.Interpreter.Cadence.WorkerOptions
	}
	iw.worker = worker.New(iw.service, iw.domain, iw.tasklist, options)
	worker.EnableVerboseLogging(config.Interpreter.VerboseDebug)

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterWorkflow(WaitforStateCompletionWorkflow)
	iw.worker.RegisterActivity(interpreter.StateStart)  // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateDecide) // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateApiWaitUntil)
	iw.worker.RegisterActivity(interpreter.StateApiExecute)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
