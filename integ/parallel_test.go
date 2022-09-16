package integ

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/parallel"
	"github.com/cadence-oss/iwf-server/service/api"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestPrallelWorkflow(t *testing.T) {
	// start test workflow server
	wfHandler, wfSvc := parallel.NewParallelWorkflow()
	testWorkflowServerPort := "9714"
	wfServer := &http.Server{
		Addr:    ":" + testWorkflowServerPort,
		Handler: wfSvc,
	}
	defer wfServer.Close()
	go func() {
		if err := wfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// start iwf api server
	iwfService := api.NewService()
	testIwfServerPort := "9715"
	iwfServer := &http.Server{
		Addr:    ":" + testIwfServerPort,
		Handler: iwfService,
	}
	defer iwfServer.Close()
	go func() {
		if err := iwfServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// start iwf interpreter worker
	interpreter := temporal.NewInterpreterWorker()
	interpreter.Start()
	defer interpreter.Close()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := parallel.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	resp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             iwfidl.PtrString(wfId),
		IwfWorkflowType:        iwfidl.PtrString(parallel.WorkflowType),
		WorkflowTimeoutSeconds: iwfidl.PtrInt32(10),
		IwfWorkerUrl:           iwfidl.PtrString("http://localhost:" + testWorkflowServerPort),
		StartStateId:           iwfidl.PtrString(parallel.State1),
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}
	fmt.Println(*resp)

	// wait for the workflow
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	run := temporalClient.GetWorkflow(context.Background(), wfId, "")
	_ = run.Get(context.Background(), nil)

	history := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int{
		"S1_start":  1,
		"S1_decide": 1,

		"S11_start":  1,
		"S11_decide": 1,
		"S12_start":  1,
		"S12_decide": 1,

		"S111_start":  1,
		"S111_decide": 1,

		"S112_start":  1,
		"S112_decide": 1,

		"S121_start":  1,
		"S121_decide": 1,

		"S122_start":  1,
		"S122_decide": 1,
	}, history, "parallel test fail, %v", history)
}
