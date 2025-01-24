package integ

import (
	"context"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestWorkflowCanceledTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowCanceled(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestWorkflowCanceled(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()

		doTestWorkflowTerminated(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestWorkflowTerminated(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()

		doTestWorkflowFail(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestWorkflowFail(t, service.BackendTypeTemporal, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func TestWorkflowCanceledCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowCanceled(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestWorkflowCanceled(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()

		doTestWorkflowTerminated(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestWorkflowTerminated(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()

		doTestWorkflowFail(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestWorkflowFail(t, service.BackendTypeCadence, minimumContinueAsNewConfigV0())
		smallWaitForFastTest()
	}
}

func doTestWorkflowCanceled(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
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
	wfId := "wf-cancel-test" + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqCancel := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	httpResp, err = reqCancel.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
		StopType:   iwfidl.CANCEL.Ptr(),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions := assert.New(t)

	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.CANCELED,
		ErrorType:      nil,
		ErrorMessage:   nil,
	}, resp, "response not expected")
}

func doTestWorkflowTerminated(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
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
	wfId := "wf-cancel-test" + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	request := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	startResp, httpResp, err := req.WorkflowStartRequest(request).Execute()
	panicAtHttpError(err, httpResp)

	assertions := assert.New(t)
	// start it again should return already started error
	_, _, err = req.WorkflowStartRequest(request).Execute()
	assertions.NotNil(err)

	// using terminate if running should go through
	request.WorkflowStartOptions.WorkflowIDReusePolicy = iwfidl.TERMINATE_IF_RUNNING.Ptr()
	_, httpResp, err = req.WorkflowStartRequest(request).Execute()
	panicAtHttpError(err, httpResp)

	// using the compatibility
	request.WorkflowStartOptions.WorkflowIDReusePolicy = nil
	request.WorkflowStartOptions.IdReusePolicy = iwfidl.ALLOW_TERMINATE_IF_RUNNING.Ptr()
	startResp, httpResp, err = req.WorkflowStartRequest(request).Execute()
	panicAtHttpError(err, httpResp)

	reqCancel := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	httpResp, err = reqCancel.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
		StopType:   iwfidl.TERMINATE.Ptr(),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.TERMINATED,
		ErrorType:      nil,
		ErrorMessage:   nil,
	}, resp, "response not expected")
}

func doTestWorkflowFail(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := signal.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
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
	wfId := "wf-cancel-test" + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startResp, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	panicAtHttpError(err, httpResp)

	reqCancel := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	httpResp, err = reqCancel.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
		StopType:   iwfidl.FAIL.Ptr(),
		Reason:     iwfidl.PtrString("fail reason"),
	}).Execute()
	panicAtHttpError(err, httpResp)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	assertions := assert.New(t)

	assertions.Equalf(&iwfidl.WorkflowGetResponse{
		WorkflowRunId:  startResp.GetWorkflowRunId(),
		WorkflowStatus: iwfidl.FAILED,
		ErrorType:      iwfidl.CLIENT_API_FAILING_WORKFLOW_ERROR_TYPE.Ptr(),
		ErrorMessage:   iwfidl.PtrString("fail reason"),
	}, resp, "response not expected")
}
