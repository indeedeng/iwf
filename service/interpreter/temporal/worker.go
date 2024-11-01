package temporal

import (
	"github.com/indeedeng/iwf/config"
	"log"

	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/interpreter"
	"github.com/indeedeng/iwf/service/interpreter/env"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/worker"
)

type InterpreterWorker struct {
	temporalClient client.Client
	worker         worker.Worker
	taskQueue      string
}

type StartOptions struct {
	DisableStickyCache bool
}

func NewInterpreterWorker(
	config config.Config, temporalClient client.Client, taskQueue string, memoEncryption bool,
	memoEncryptionConverter converter.DataConverter, unifiedClient uclient.UnifiedClient,
) *InterpreterWorker {
	env.SetSharedEnv(config, memoEncryption, memoEncryptionConverter, unifiedClient, taskQueue)

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
	var options StartOptions

	// default options
	options.DisableStickyCache = false

	iw.StartWithOptions(options)
}

func (iw *InterpreterWorker) StartWithOptions(startOptions StartOptions) {
	config := env.GetSharedConfig()
	options := worker.Options{
		MaxConcurrentActivityTaskPollers: 10,
		// TODO: this cannot be too small otherwise the persistence_test for continueAsNew will fail, probably a bug in Temporal goSDK.
		// It seems work as "parallelism" of something... need to report a bug ticket...
		MaxConcurrentWorkflowTaskPollers: 10,
	}
	if config.Interpreter.Temporal != nil && config.Interpreter.Temporal.WorkerOptions != nil {
		options = *config.Interpreter.Temporal.WorkerOptions
	}
	iw.worker = worker.New(iw.temporalClient, iw.taskQueue, options)
	worker.EnableVerboseLogging(config.Interpreter.VerboseDebug)

	if startOptions.DisableStickyCache {
		worker.SetStickyWorkflowCacheSize(0)
	}

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterWorkflow(WaitforStateCompletionWorkflow)
	iw.worker.RegisterActivity(interpreter.StateStart)  // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateDecide) // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateApiWaitUntil)
	iw.worker.RegisterActivity(interpreter.StateApiExecute)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)
	iw.worker.RegisterActivity(interpreter.InvokeWorkerRpc)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
