package cadence

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/config"
	"github.com/indeedeng/iwf/service/common/blobstore"
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

func NewInterpreterWorker(
	config config.Config, service workflowserviceclient.Interface, domain, tasklist string, closeFunc func(),
	unifiedClient uclient.UnifiedClient,
	store blobstore.BlobStore,
) *InterpreterWorker {
	env.SetSharedEnv(config, false, nil, unifiedClient, tasklist, store)
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

func (iw *InterpreterWorker) StartWithStickyCacheDisabledForTest() {
	iw.start(true)
}

func (iw *InterpreterWorker) Start() {
	iw.start(false)
}

func (iw *InterpreterWorker) start(disableStickyCache bool) {
	cfg := env.GetSharedConfig()
	var options worker.Options

	if cfg.Interpreter.Cadence != nil && cfg.Interpreter.Cadence.WorkerOptions != nil {
		options = *cfg.Interpreter.Cadence.WorkerOptions
	}

	// override default
	if options.MaxConcurrentActivityTaskPollers == 0 {
		options.MaxConcurrentActivityTaskPollers = 10
	}

	// override default
	if options.MaxConcurrentDecisionTaskPollers == 0 {
		options.MaxConcurrentDecisionTaskPollers = 10
	}

	// When DisableStickyCache is true it can harm performance; should not be used in production environment
	if disableStickyCache {
		options.DisableStickyExecution = true
		fmt.Println("Cadence worker: Sticky cache disabled")
	}

	iw.worker = worker.New(iw.service, iw.domain, iw.tasklist, options)
	worker.EnableVerboseLogging(cfg.Interpreter.VerboseDebug)

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterWorkflow(WaitforStateCompletionWorkflow)
	iw.worker.RegisterActivity(interpreter.StateStart)  // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateDecide) // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateApiWaitUntil)
	iw.worker.RegisterActivity(interpreter.StateApiExecute)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)
	iw.worker.RegisterActivity(interpreter.CleanupBlobStore)

	err := iw.worker.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

	if cfg.ExternalStorage.Enabled {
		schedule := cfg.ExternalStorage.CleanupCronSchedule
		if schedule == "" {
			schedule = "0 * * * * *"
		}
		for _, storeCfg := range cfg.ExternalStorage.SupportedStorages {
			err = env.GetUnifiedClient().StartBlobStoreCleanupWorkflow(
				context.Background(), iw.tasklist,
				"blobstore-cleanup-"+storeCfg.StorageId,
				schedule,
				storeCfg.StorageId)
			if err != nil {
				log.Fatalln("Unable to start blobstore cleanup workflow", err)
			}
		}
	}
}
