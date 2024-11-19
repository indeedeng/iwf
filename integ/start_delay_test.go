package integ

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/basic"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func TestStartDelayTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	doTestStartDelay(t, service.BackendTypeTemporal, nil)
}

func TestStartDelayCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	doTestStartDelay(t, service.BackendTypeCadence, nil)
}

func doTestStartDelay(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	// start test workflow server
	wfHandler := basic.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler)
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
	wfId := basic.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	wfInput := &iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test data"),
	}
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	startReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        basic.WorkflowType,
		WorkflowTimeoutSeconds: 100,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(basic.State1),
		StateInput:             wfInput,
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowStartDelaySeconds: iwfidl.PtrInt32(10),
			WorkflowConfigOverride:    config,
			WorkflowIDReusePolicy:     ptr.Any(iwfidl.REJECT_DUPLICATE),
			RetryPolicy: &iwfidl.WorkflowRetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(11),
				BackoffCoefficient:     iwfidl.PtrFloat32(11),
				MaximumAttempts:        iwfidl.PtrInt32(11),
				MaximumIntervalSeconds: iwfidl.PtrInt32(11),
			},
		},
		StateOptions: &iwfidl.WorkflowStateOptions{
			StartApiTimeoutSeconds:  iwfidl.PtrInt32(12),
			DecideApiTimeoutSeconds: iwfidl.PtrInt32(13),
			StartApiRetryPolicy: &iwfidl.RetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(12),
				BackoffCoefficient:     iwfidl.PtrFloat32(12),
				MaximumAttempts:        iwfidl.PtrInt32(12),
				MaximumIntervalSeconds: iwfidl.PtrInt32(12),
			},
			DecideApiRetryPolicy: &iwfidl.RetryPolicy{
				InitialIntervalSeconds: iwfidl.PtrInt32(13),
				BackoffCoefficient:     iwfidl.PtrFloat32(13),
				MaximumAttempts:        iwfidl.PtrInt32(13),
				MaximumIntervalSeconds: iwfidl.PtrInt32(13),
			},
		},
	}

	timeSentReq := time.Now()
	_, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(5 * time.Second)
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	// here the delay is startDelay + execution time period, and the execution time period is negligible
	delay := time.Since(timeSentReq)

	assertions := assert.New(t)
	assertions.True(delay.Seconds() > 8, "delay is %v", delay)
	assertions.True(delay.Seconds() < 12, "delay is %v", delay)
}
