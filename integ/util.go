package integ

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	cadenceapi "github.com/indeedeng/iwf/service/api/cadence"
	temporalapi "github.com/indeedeng/iwf/service/api/temporal"
	"github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/log/loggerimpl"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.uber.org/cadence/encoded"
	"log"
	"net/http"
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

func startWorkflowWorker(handler common.WorkflowHandler) (closeFunc func()) {
	router := gin.Default()
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

func startIwfService(backendType service.BackendType) (closeFunc func()) {
	_, cf := doStartIwfServiceWithClient(backendType)
	return cf
}

func doStartIwfServiceWithClient(backendType service.BackendType) (uclient api.UnifiedClient, closeFunc func()) {
	if backendType == service.BackendTypeTemporal {
		if integTemporalUclientCached == nil {
			return doOnceStartIwfServiceWithClient(backendType)
		}
		return integTemporalUclientCached, func() {}
	}

	if integCadenceUclientCached == nil {
		return doOnceStartIwfServiceWithClient(backendType)
	}
	return integCadenceUclientCached, func() {}
}

var integCadenceUclientCached api.UnifiedClient
var integTemporalUclientCached api.UnifiedClient

func doOnceStartIwfServiceWithClient(backendType service.BackendType) (uclient api.UnifiedClient, closeFunc func()) {
	if backendType == service.BackendTypeTemporal {
		temporalClient := createTemporalClient()
		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}
		uclient = temporalapi.NewTemporalClient(temporalClient, testNamespace, converter.GetDefaultDataConverter())
		iwfService := api.NewService(&config.Config{}, uclient, logger)
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
		interpreter := temporal.NewInterpreterWorker(temporalClient, service.TaskQueue)
		interpreter.Start()
		return uclient, func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else if backendType == service.BackendTypeCadence {
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
		iwfService := api.NewService(&config.Config{}, uclient, logger)
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
		interpreter := cadence.NewInterpreterWorker(serviceClient, iwf.DefaultCadenceDomain, service.TaskQueue, closeFunc)
		interpreter.Start()
		return uclient, func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else {
		panic("not supported backend type " + backendType)
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
