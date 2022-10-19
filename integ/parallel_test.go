package integ

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/integ/workflow/parallel"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestParallelWorkflowTemporal(t *testing.T) {
	doTestParallelWorkflow(t, service.BackendTypeTemporal)
}

func TestParallelWorkflowCadence(t *testing.T) {
	doTestParallelWorkflow(t, service.BackendTypeCadence)
}

func doTestParallelWorkflow(t *testing.T, backendType service.BackendType) {
	// start test workflow server
	wfHandler := parallel.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

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
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
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

	// wait for the workflow
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke get with long wait api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	history, _ := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
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
