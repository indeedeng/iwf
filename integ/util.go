package integ

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	cadenceapi "github.com/indeedeng/iwf/service/api/cadence"
	temporalapi "github.com/indeedeng/iwf/service/api/temporal"
	"github.com/indeedeng/iwf/service/common/log/loggerimpl"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.uber.org/cadence/encoded"
	"log"
	"net/http"
	"time"
)

const testNamespace = "default"

func createTemporalClient() client.Client {
	temporalClient, err := client.Dial(client.Options{
		HostPort:  *temporalHostPort,
		Namespace: testNamespace,
	})
	if err != nil {
		log.Fatalf("unable to connect to Temporal %v", err)
	}
	return temporalClient
}

func startWorkflowWorkerWithRpc(handler common.WorkflowHandlerWithRpc) (closeFunc func()) {
	router := gin.Default()
	router.POST(service.WorkflowWorkerRpcApi, handler.ApiV1WorkflowWorkerRpc)
	return doStartWorkflowWorker(handler, router)
}

func startWorkflowWorker(handler common.WorkflowHandler) (closeFunc func()) {
	router := gin.Default()
	return doStartWorkflowWorker(handler, router)
}
func doStartWorkflowWorker(handler common.WorkflowHandler, router *gin.Engine) (closeFunc func()) {
	router.POST(service.StateStartApi, handler.ApiV1WorkflowStateStart)
	router.POST(service.StateDecideApi, handler.ApiV1WorkflowStateDecide)

	wfServer := &http.Server{
		Addr:    ":" + testWorkflowServerPort,
		Handler: router,
	}
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return func() { wfServer.Close() }
}

type IwfServiceTestConfig struct {
	BackendType    service.BackendType
	MemoEncryption bool
}

func startIwfService(backendType service.BackendType) (closeFunc func()) {
	_, cf := startIwfServiceWithClient(backendType)
	return cf
}

func startIwfServiceByConfig(config IwfServiceTestConfig) (uclient api.UnifiedClient, closeFunc func()) {
	return doStartIwfServiceWithClient(config)
}

func startIwfServiceWithClient(backendType service.BackendType) (uclient api.UnifiedClient, closeFunc func()) {
	return doStartIwfServiceWithClient(IwfServiceTestConfig{BackendType: backendType})

	//if backendType == service.BackendTypeTemporal {
	//if integTemporalUclientCached == nil {
	//	return doStartIwfServiceWithClient(backendType)
	//}
	//return integTemporalUclientCached, func() {}
	//}
	//if integCadenceUclientCached == nil {
	//	return doStartIwfServiceWithClient(backendType)
	//}
	//return integCadenceUclientCached, func() {}
}

// disable caching for now as it makes it difficult to test memo
//var integCadenceUclientCached api.UnifiedClient
//var integTemporalUclientCached api.UnifiedClient

func doStartIwfServiceWithClient(config IwfServiceTestConfig) (uclient api.UnifiedClient, closeFunc func()) {
	if config.BackendType == service.BackendTypeTemporal {
		temporalClient := createTemporalClient()
		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}
		uclient = temporalapi.NewTemporalClient(temporalClient, testNamespace, converter.GetDefaultDataConverter(), config.MemoEncryption) // TODO set correct data converter for memo encryption
		iwfService := api.NewService(testConfig, uclient, logger)
		iwfServer := &http.Server{
			Addr:    ":" + testIwfServerPort,
			Handler: iwfService,
		}
		go func() {
			if err := iwfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// start iwf interpreter worker
		interpreter := temporal.NewInterpreterWorker(testConfig, temporalClient, service.TaskQueue, nil) // TODO set correct data converter for memo encryption
		interpreter.Start()
		return uclient, func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else if config.BackendType == service.BackendTypeCadence {
		serviceClient, closeFunc, err := iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
		if err != nil {
			log.Fatalf("cannot connnect to Cadence %v", err)
		}

		cadenceClient, err := iwf.BuildCadenceClient(serviceClient, iwf.DefaultCadenceDomain)

		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}
		uclient = cadenceapi.NewCadenceClient(iwf.DefaultCadenceDomain, cadenceClient, serviceClient, encoded.GetDefaultDataConverter(), closeFunc)
		iwfService := api.NewService(testConfig, uclient, logger)
		iwfServer := &http.Server{
			Addr:    ":" + testIwfServerPort,
			Handler: iwfService,
		}
		go func() {
			if err := iwfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen: %s\n", err)
			}
		}()

		// start iwf interpreter worker
		interpreter := cadence.NewInterpreterWorker(testConfig, serviceClient, iwf.DefaultCadenceDomain, service.TaskQueue, closeFunc)
		interpreter.Start()
		return uclient, func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else {
		panic("not supported backend type " + config.BackendType)
	}
}

func panicAtHttpError(err error, httpResp *http.Response) {
	if err != nil {
		panic(err)
	}
	if httpResp.StatusCode != http.StatusOK {
		panic("Status not success" + httpResp.Status)
	}
}

func panicAtHttpErrorOrWorkflowUncompleted(err error, httpResp *http.Response, resp *iwfidl.WorkflowGetResponse) {
	if err != nil {
		panic(err)
	}
	if httpResp.StatusCode != http.StatusOK {
		panic("Status not success" + httpResp.Status)
	}
	if resp.WorkflowStatus != iwfidl.COMPLETED {
		panic("Workflow uncompleted:" + resp.WorkflowStatus)
	}
}

func smallWaitForFastTest() {
	du := time.Millisecond * time.Duration(*repeatInterval)
	if *repeatIntegTest == 0 {
		du = time.Millisecond
	}
	time.Sleep(du)
}

func minimumContinueAsNewConfig() *iwfidl.WorkflowConfig {
	return &iwfidl.WorkflowConfig{
		ContinueAsNewThreshold: iwfidl.PtrInt32(1),
	}
}
