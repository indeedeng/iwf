package integ

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/cmd/server/iwf"
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
	if err == nil {
		return temporalClient
	}

	for i := 0; i < *dependencyWaitSeconds; i++ {
		fmt.Println("wait for Temporal to be up...last err: ", err)
		time.Sleep(time.Second)
		temporalClient, err = client.Dial(client.Options{
			HostPort:  *temporalHostPort,
			Namespace: testNamespace,
		})
		if err == nil {
			return temporalClient
		}
	}
	log.Fatalf("unable to connect to Temporal %v", err)
	return nil
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
		temporalClient := createTemporalClient()
		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}
		uclient = temporalapi.NewTemporalClient(temporalClient, testNamespace, converter.GetDefaultDataConverter())
		iwfService := api.NewService(uclient, logger)
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
			for i := 0; i < *dependencyWaitSeconds; i++ {
				fmt.Println("wait for Cadence to be up...last err: ", err)
				time.Sleep(time.Second)

				serviceClient, closeFunc, err = iwf.BuildCadenceServiceClient(iwf.DefaultCadenceHostPort)
				if err == nil {
					break
				}
			}
			if err != nil {
				log.Fatalf("cannot connnect to Cadence %v", err)
			}
		}

		for i := 0; i < *dependencyWaitSeconds; i++ {
			fmt.Println("wait for Cadence domain/Search attributes to be ready...")
			time.Sleep(time.Second)
			resp, err := serviceClient.GetSearchAttributes(context.Background())
			ready := false
			if err == nil {
				for key, _ := range resp.GetKeys() {
					// NOTE: this is the last one we registered in init-ci-cadence.sh
					if key == service.SearchAttributeIwfWorkflowType {
						ready = true
						break
					}
				}
			}
			if ready {
				break
			}
		}

		cadenceClient, err := iwf.BuildCadenceClient(serviceClient, iwf.DefaultCadenceDomain)

		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}
		uclient = cadenceapi.NewCadenceClient(iwf.DefaultCadenceDomain, cadenceClient, serviceClient, encoded.GetDefaultDataConverter(), closeFunc)
		iwfService := api.NewService(uclient, logger)
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
