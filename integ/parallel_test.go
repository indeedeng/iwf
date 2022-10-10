package integ

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/parallel"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/cadence-oss/iwf-server/service/api"
	temporalapi "github.com/cadence-oss/iwf-server/service/api/temporal"
	"github.com/cadence-oss/iwf-server/service/interpreter/temporal"
	"github.com/stretchr/testify/assert"
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
	temporalClient := createTemporalClient()
	iwfService := api.NewService(temporalapi.NewTemporalClient(temporalClient))
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
	wfId := parallel.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	resp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        parallel.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           parallel.State1,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}
	fmt.Println(*resp)
	defer temporalClient.TerminateWorkflow(context.Background(), wfId, "", "terminate incase not completed")

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithLongWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId:   wfId,
		NeedsResults: iwfidl.PtrBool(true),
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

	assertions.Equal(service.WorkflowStatusCompleted, resp2.GetWorkflowStatus())
	assertions.Equal(4, len(resp2.GetResults()))
	//assertions.Equal([]iwfidl.StateCompletionOutput{
	//	{
	//		CompletedStateId:          parallel.State111,
	//		CompletedStateExecutionId: parallel.State111 + "-1",
	//		CompletedStateOutput: &iwfidl.EncodedObject{
	//			Encoding: iwfidl.PtrString("json"),
	//			Data:     iwfidl.PtrString("from " + parallel.State111),
	//		},
	//	},
	//	{
	//		CompletedStateId:          parallel.State112,
	//		CompletedStateExecutionId: parallel.State112 + "-1",
	//		CompletedStateOutput: &iwfidl.EncodedObject{
	//			Encoding: iwfidl.PtrString("json"),
	//			Data:     iwfidl.PtrString("from " + parallel.State112),
	//		},
	//	},
	//	{
	//		CompletedStateId:          parallel.State121,
	//		CompletedStateExecutionId: parallel.State121 + "-1",
	//		CompletedStateOutput: &iwfidl.EncodedObject{
	//			Encoding: iwfidl.PtrString("json"),
	//			Data:     iwfidl.PtrString("from " + parallel.State121),
	//		},
	//	},
	//	{
	//		CompletedStateId:          parallel.State122,
	//		CompletedStateExecutionId: parallel.State122 + "-1",
	//		CompletedStateOutput: &iwfidl.EncodedObject{
	//			Encoding: iwfidl.PtrString("json"),
	//			Data:     iwfidl.PtrString("from " + parallel.State122),
	//		},
	//	},
	//}, resp2.GetResults())
}
