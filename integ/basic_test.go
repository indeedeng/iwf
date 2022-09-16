package integ

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/basic"
	"github.com/cadence-oss/iwf-server/service/api"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestBasicWorkflow(t *testing.T) {
	// start test workflow server
	basicWorkflow := basic.NewBasicWorkflow()
	testWorkflowServerPort := "9714"
	wfServer := &http.Server{
		Addr:    ":" + testWorkflowServerPort,
		Handler: basicWorkflow,
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
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	resp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             iwfidl.PtrString("test-wf-id"),
		IwfWorkflowType:        iwfidl.PtrString("test-iwf-wf-type"),
		WorkflowTimeoutSeconds: iwfidl.PtrInt32(10),
		IwfWorkerUrl:           iwfidl.PtrString("http://localhost:" + testWorkflowServerPort),
		StartStateId:           iwfidl.PtrString(basic.State1),
	}).Execute()
	if err != nil {
		fmt.Println("fail to invoke start api", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		fmt.Println("status not success" + httpResp.Status)
	}
	fmt.Println(*resp)

	// give some time for workflow to execute
	time.Sleep(time.Second * 10)
}
