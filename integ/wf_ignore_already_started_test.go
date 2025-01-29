package integ

import (
	"context"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/wf_ignore_already_started"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestIgnoreAlreadyStartedWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		// Case 1: secondReq does not ignore AlreadyStartedError; second start request should return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeTemporal, nil, nil, true)
		smallWaitForFastTest()

		// Case 2: secondReq does ignore AlreadyStartedError; second start request should not return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeTemporal, nil, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
		}, false)
		smallWaitForFastTest()

		// Case 3: secondReq does ignore AlreadyStartedError only if requestId match; they do, so second start request should not return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeTemporal, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: false,
			RequestId:                 iwfidl.PtrString("test"),
		}, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
			RequestId:                 iwfidl.PtrString("test"),
		}, false)
		smallWaitForFastTest()

		// Case 4: secondReq does ignore AlreadyStartedError only if requestId match; they do not, so second start request should return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeTemporal, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: false,
			RequestId:                 iwfidl.PtrString("test1"),
		}, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
			RequestId:                 iwfidl.PtrString("test2"),
		}, true)
		smallWaitForFastTest()
	}
}

func TestIgnoreAlreadyStartedWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		// Case 1: secondReq does not ignore AlreadyStartedError; second start request should return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeCadence, nil, nil, true)
		smallWaitForFastTest()

		// Case 2: secondReq does ignore AlreadyStartedError; second start request should not return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeCadence, nil, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
		}, false)
		smallWaitForFastTest()

		// Case 3: secondReq does ignore AlreadyStartedError only if requestId match; they do, so second start request should not return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeCadence, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: false,
			RequestId:                 iwfidl.PtrString("test"),
		}, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
			RequestId:                 iwfidl.PtrString("test"),
		}, false)
		smallWaitForFastTest()

		// Case 4: secondReq does ignore AlreadyStartedError only if requestId match; they do not, so second start request should return error
		doIgnoreAlreadyStartedWorkflow(t, service.BackendTypeCadence, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: false,
			RequestId:                 iwfidl.PtrString("test1"),
		}, &iwfidl.WorkflowAlreadyStartedOptions{
			IgnoreAlreadyStartedError: true,
			RequestId:                 iwfidl.PtrString("test2"),
		}, true)
		smallWaitForFastTest()
	}
}

func doIgnoreAlreadyStartedWorkflow(t *testing.T, backendType service.BackendType, firstReqConfig *iwfidl.WorkflowAlreadyStartedOptions, secondReqConfig *iwfidl.WorkflowAlreadyStartedOptions, errorExpected bool) {
	// start test workflow server
	wfHandler := wf_ignore_already_started.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	_, closeFunc2 := startIwfServiceByConfig(IwfServiceTestConfig{
		BackendType: backendType,
	})
	defer closeFunc2()

	// start a workflow
	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})
	wfId := wf_ignore_already_started.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())

	firstReq := createReq(wfId, firstReqConfig)

	firstRes, firstHttpResp, err := req.WorkflowStartRequest(firstReq).Execute()
	failTestAtHttpError(err, firstHttpResp, t)

	secondReq := createReq(wfId, secondReqConfig)
	secondRes, secondHttpResp, err := req.WorkflowStartRequest(secondReq).Execute()

	assertions := assert.New(t)
	if errorExpected {
		apiErr, ok := err.(*iwfidl.GenericOpenAPIError)
		if !ok {
			log.Fatalf("Should fail to invoke start api %v", err)
		}
		errResp, ok := apiErr.Model().(iwfidl.ErrorResponse)
		if !ok {
			log.Fatalf("should be error response")
		}
		assertions.Equal(iwfidl.WORKFLOW_ALREADY_STARTED_SUB_STATUS, errResp.GetSubStatus())
		assertions.Equal(400, secondHttpResp.StatusCode)
	} else {
		assertions.Equal(nil, err)
		assertions.Equal(firstRes.GetWorkflowRunId(), secondRes.GetWorkflowRunId())
		assertions.Equal(200, secondHttpResp.StatusCode)
	}

	// Terminate the workflow once tests completed
	stopReq := apiClient.DefaultApi.ApiV1WorkflowStopPost(context.Background())
	_, err = stopReq.WorkflowStopRequest(iwfidl.WorkflowStopRequest{
		WorkflowId: wfId,
		StopType:   iwfidl.TERMINATE.Ptr(),
	}).Execute()
}

func createReq(wfId string, options *iwfidl.WorkflowAlreadyStartedOptions) iwfidl.WorkflowStartRequest {
	return iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wf_ignore_already_started.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wf_ignore_already_started.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowAlreadyStartedOptions: options,
		},
	}
}
