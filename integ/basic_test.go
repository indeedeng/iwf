package integ

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/basic"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/api"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestBasicWorkflow(t *testing.T) {
	// start test workflow server
	wfHandler, basicWorkflow := basic.NewBasicWorkflow()
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
	temporalClient := createTemporalClient()
	iwfService := api.NewService(temporalClient)
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
	interpreter := temporal.NewInterpreterWorker(temporalClient)
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
	wfId := basic.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test data"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	resp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        basic.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           basic.State1,
		StateInput:             wfInput,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}
	fmt.Println(*resp)
	defer temporalClient.TerminateWorkflow(context.Background(), wfId, "", "terminate incase not completed")

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId:         wfId,
		NeedsResults:       iwfidl.PtrBool(true),
		WaitIfStillRunning: iwfidl.PtrBool(true),
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	history := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
	}, history, "basic test fail, %v", history)

	assertions.Equal(service.WorkflowStatusCompleted, resp2.GetWorkflowStatus())
	assertions.Equal(1, len(resp2.GetResults()))
	assertions.Equal(iwfidl.StateCompletionOutput{
		CompletedStateId:          "S2",
		CompletedStateExecutionId: "S2-1",
		CompletedStateOutput:      wfInput,
	}, resp2.GetResults()[0])
}
