package temporal

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/service/common/blobstore"
	"log"

	"github.com/indeedeng/iwf/config"
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

func NewInterpreterWorker(
	config config.Config, temporalClient client.Client, taskQueue string, memoEncryption bool,
	memoEncryptionConverter converter.DataConverter, unifiedClient uclient.UnifiedClient,
	store blobstore.BlobStore,
) *InterpreterWorker {
	env.SetSharedEnv(config, memoEncryption, memoEncryptionConverter, unifiedClient, taskQueue, store)

	return &InterpreterWorker{
		temporalClient: temporalClient,
		taskQueue:      taskQueue,
	}
}

func (iw *InterpreterWorker) Close() {
	iw.temporalClient.Close()
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

	if cfg.Interpreter.Temporal != nil && cfg.Interpreter.Temporal.WorkerOptions != nil {
		options = *cfg.Interpreter.Temporal.WorkerOptions
	}

	// override default
	if options.MaxConcurrentActivityTaskPollers == 0 {
		options.MaxConcurrentActivityTaskPollers = 10
	}

	// override default
	if options.MaxConcurrentWorkflowTaskPollers == 0 {
		// TODO: this cannot be too small otherwise the persistence_test for continueAsNew will fail, probably a bug in Temporal goSDK.
		// It seems work as "parallelism" of something... need to report a bug ticket...
		options.MaxConcurrentWorkflowTaskPollers = 10
	}

	// When DisableStickyCache is true it can harm performance; should not be used in production environment
	if disableStickyCache {
		worker.SetStickyWorkflowCacheSize(0)
		fmt.Println("Temporal worker: Sticky cache disabled")
	}

	iw.worker = worker.New(iw.temporalClient, iw.taskQueue, options)
	worker.EnableVerboseLogging(cfg.Interpreter.VerboseDebug)

	iw.worker.RegisterWorkflow(Interpreter)
	iw.worker.RegisterWorkflow(WaitforStateCompletionWorkflow)
	iw.worker.RegisterActivity(interpreter.StateStart)  // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateDecide) // TODO: remove in next release
	iw.worker.RegisterActivity(interpreter.StateApiWaitUntil)
	iw.worker.RegisterActivity(interpreter.StateApiExecute)
	iw.worker.RegisterActivity(interpreter.DumpWorkflowInternal)
	iw.worker.RegisterActivity(interpreter.InvokeWorkerRpc)
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
				context.Background(), iw.taskQueue,
				"blobstore-cleanup-"+storeCfg.StorageId,
				schedule,
				storeCfg.StorageId)
			if err != nil {
				log.Fatalln("Unable to start blobstore cleanup workflow", err)
			}
		}
	}
}
