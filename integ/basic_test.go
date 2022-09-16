package integ

import (
	"bytes"
	"encoding/json"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/basic"
	"github.com/cadence-oss/iwf-server/service/api"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"log"
	"net/http"
	"testing"
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
	httpClient := &http.Client{}

	reqStr, err := json.Marshal(iwfidl.WorkflowStartRequest{
		IwfWorkflowType: iwfidl.PtrString(basic.WorkflowType),
	})
	if err != nil {
		log.Fatalf("Failed to marshal request %v", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:"+testIwfServerPort+""+api.WorkflowStartApiPath, bytes.NewBuffer(reqStr))
	if err != nil {
		log.Fatalf("Failed to create request %v", err)

	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatalf("failed to start workflow %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("failed to start workflow, status code: %v", resp.StatusCode)
	}
}
