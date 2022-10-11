package integ

import (
	"github.com/cadence-oss/iwf-server/cmd/server/iwf"
	"github.com/cadence-oss/iwf-server/integ/workflow/common"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/api"
	cadenceapi "github.com/cadence-oss/iwf-server/service/api/cadence"
	temporalapi "github.com/cadence-oss/iwf-server/service/api/temporal"
	"github.com/cadence-oss/iwf-server/service/interpreter/cadence"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
)

func createTemporalClient() client.Client {
	temporalClient, err := client.Dial(client.Options{})
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
	if backendType == service.BackendTypeTemporal {
		temporalClient := createTemporalClient()
		iwfService := api.NewService(temporalapi.NewTemporalClient(temporalClient))
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
		interpreter := temporal.NewInterpreterWorker(temporalClient)
		interpreter.Start()
		return func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else if backendType == service.BackendTypeCadence {
		serviceClient, closeFunc, err := iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
		if err != nil {
			log.Fatalf("cannot connnect to Cadence %v", err)
		}

		cadenceClient, err := iwf.BuildCadenceClient(serviceClient, iwf.DefaultCadenceDomain)

		iwfService := api.NewService(cadenceapi.NewCadenceClient(cadenceClient, closeFunc))
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
		interpreter := cadence.NewInterpreterWorker(serviceClient, iwf.DefaultCadenceDomain, closeFunc)
		interpreter.Start()
		return func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else {
		panic("not supported backend type " + backendType)
	}
}
