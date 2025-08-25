package integ

import (
	"fmt"
	"github.com/indeedeng/iwf/integ/helpers"
	cadenceapi "github.com/indeedeng/iwf/service/client/cadence"
	temporalapi "github.com/indeedeng/iwf/service/client/temporal"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/api"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/log/loggerimpl"
	"github.com/indeedeng/iwf/service/interpreter/cadence"
	"github.com/indeedeng/iwf/service/interpreter/temporal"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.uber.org/cadence/encoded"
)

const testNamespace = "default"

func createTemporalClient(dataConverter converter.DataConverter) client.Client {
	temporalClient, err := client.Dial(client.Options{
		HostPort:      *temporalHostPort,
		Namespace:     testNamespace,
		DataConverter: dataConverter,
	})
	if err != nil {
		log.Fatalf("unable to connect to Temporal %v", err)
	}
	return temporalClient
}

func startWorkflowWorkerWithRpc(handler common.WorkflowHandlerWithRpc, t *testing.T) (closeFunc func()) {
	router := gin.Default()
	router.POST(service.WorkflowWorkerRpcApi, func(c *gin.Context) {
		handler.ApiV1WorkflowWorkerRpc(c, t)
	})
	return doStartWorkflowWorker(handler, t, router)
}

func startWorkflowWorker(handler common.WorkflowHandler, t *testing.T) (closeFunc func()) {
	router := gin.Default()
	return doStartWorkflowWorker(handler, t, router)
}
func doStartWorkflowWorker(handler common.WorkflowHandler, t *testing.T, router *gin.Engine) (closeFunc func()) {
	router.POST(service.StateStartApi, func(c *gin.Context) {
		handler.ApiV1WorkflowStateStart(c, t)
	})
	router.POST(service.StateDecideApi, func(c *gin.Context) {
		handler.ApiV1WorkflowStateDecide(c, t)
	})

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
	BackendType                      service.BackendType
	MemoEncryption                   bool
	DisableFailAtMemoIncompatibility bool // default to false so that we will fail at test
	DefaultHeaders                   map[string]string
}

func startIwfService(backendType service.BackendType) (closeFunc func()) {
	_, cf := startIwfServiceWithClient(backendType)
	return cf
}

func startIwfServiceByConfig(config IwfServiceTestConfig) (uclient uclient.UnifiedClient, closeFunc func()) {
	return doStartIwfServiceWithClient(config)
}

func startIwfServiceWithClient(backendType service.BackendType) (uclient uclient.UnifiedClient, closeFunc func()) {
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

func doStartIwfServiceWithClient(config IwfServiceTestConfig) (uclient uclient.UnifiedClient, closeFunc func()) {
	if config.BackendType == service.BackendTypeTemporal {
		dataConverter := converter.GetDefaultDataConverter()
		if config.MemoEncryption {
			dataConverter = encryptionDataConverter
		}

		temporalClient := createTemporalClient(dataConverter)
		logger, err := loggerimpl.NewDevelopment()
		if err != nil {
			panic(err)
		}

		testCfg := createTestConfig(config)

		uclient = temporalapi.NewTemporalClient(temporalClient, testNamespace, dataConverter, config.MemoEncryption, &testCfg.Api.QueryWorkflowFailedRetryPolicy)
		iwfService := api.NewService(testCfg, uclient, logger, nil, "") // TODO pass s3 client for integ test
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
		interpreter := temporal.NewInterpreterWorker(testCfg, temporalClient, service.TaskQueue, config.MemoEncryption, dataConverter, uclient)
		if *disableStickyCache {
			interpreter.StartWithStickyCacheDisabledForTest()
		} else {
			interpreter.Start()
		}
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

		testCfg := createTestConfig(config)

		uclient = cadenceapi.NewCadenceClient(iwf.DefaultCadenceDomain, cadenceClient, serviceClient, encoded.GetDefaultDataConverter(), closeFunc, &testCfg.Api.QueryWorkflowFailedRetryPolicy)
		iwfService := api.NewService(testCfg, uclient, logger, nil, "") // pass in for integ tests
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
		interpreter := cadence.NewInterpreterWorker(testCfg, serviceClient, iwf.DefaultCadenceDomain, service.TaskQueue, closeFunc, uclient)
		if *disableStickyCache {
			interpreter.StartWithStickyCacheDisabledForTest()
		} else {
			interpreter.Start()
		}
		return uclient, func() {
			iwfServer.Close()
			interpreter.Close()
		}
	} else {
		panic("not supported backend type " + config.BackendType)
	}
}

func failTestAtError(err error, t *testing.T) {
	if err != nil {
		helpers.FailTestWithError(err, t)
	}
}

func failTestAtHttpError(err error, httpResp *http.Response, t *testing.T) {
	if err != nil {
		helpers.FailTestWithError(err, t)
	}
	if httpResp.StatusCode != http.StatusOK {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("HTTP status not success: %v", httpResp.Status), t)
	}
}

func failTestAtHttpErrorOrWorkflowUncompleted(err error, httpResp *http.Response, resp *iwfidl.WorkflowGetResponse, t *testing.T) {
	if err != nil {
		helpers.FailTestWithError(err, t)
	}
	if httpResp.StatusCode != http.StatusOK {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("HTTP status not success: %v", httpResp.Status), t)
	}
	if resp.WorkflowStatus != iwfidl.COMPLETED {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("Workflow uncompleted: %v", resp.WorkflowStatus), t)
	}
}

func smallWaitForFastTest() {
	du := time.Millisecond * time.Duration(*repeatInterval)
	if *repeatIntegTest == 0 {
		du = time.Millisecond
	}
	time.Sleep(du)
}

func minimumContinueAsNewConfig(optimizeActivity bool) *iwfidl.WorkflowConfig {
	return &iwfidl.WorkflowConfig{
		ContinueAsNewThreshold: iwfidl.PtrInt32(1),
		OptimizeActivity:       iwfidl.PtrBool(optimizeActivity),
	}
}

func minimumGreedyTimerConfig() *iwfidl.WorkflowConfig {
	return greedyTimerConfig(false)
}

func greedyTimerConfig(continueAsNew bool) *iwfidl.WorkflowConfig {
	if continueAsNew {
		return &iwfidl.WorkflowConfig{
			ContinueAsNewThreshold: iwfidl.PtrInt32(1),
			OptimizeTimer:          iwfidl.PtrBool(true),
		}
	}

	return &iwfidl.WorkflowConfig{
		OptimizeTimer: iwfidl.PtrBool(true),
	}
}

func minimumContinueAsNewConfigV0() *iwfidl.WorkflowConfig {
	return minimumContinueAsNewConfig(false)
}

func getBackendTypes() []service.BackendType {
	backendTypesToTest := []service.BackendType{}

	if *temporalIntegTest {
		backendTypesToTest = append(backendTypesToTest, service.BackendTypeTemporal)
	}

	if *cadenceIntegTest {
		backendTypesToTest = append(backendTypesToTest, service.BackendTypeCadence)
	}

	return backendTypesToTest
}
