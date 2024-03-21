package integ

import (
	"context"
	"github.com/indeedeng/iwf/cmd/server/iwf"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/basic"
	"github.com/indeedeng/iwf/service"
	config2 "github.com/indeedeng/iwf/service/common/config"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/client"
	rawLog "log"
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
		BackendType:         backendType,
		OptimizationVersion: ptr.Any(config2.OptimizationVersionNone),
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
			// TODO: need more work to write integ test for cron
			// manual testing for now by uncomment the following line
			// CronSchedule:           iwfidl.PtrString("* * * * *"),
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

	resp, httpResp, err := req.WorkflowStartRequest(startReq).Execute()
	panicAtHttpError(err, httpResp)

	time.Sleep(10 * time.Second)
	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp)

	var delay time.Duration
	if backendType == service.BackendTypeTemporal {
		temporalClient, err := client.Dial(client.Options{})
		if err != nil {
			t.Fatal(err)
		}
		defer temporalClient.Close()

		describeResp, err := temporalClient.DescribeWorkflowExecution(context.Background(), wfId, *resp.WorkflowRunId)
		if err != nil {
			t.Fatal(err)
		}
		if describeResp.WorkflowExecutionInfo.GetStartTime() == nil {
			t.Fatal("start time is nil")
		}

		delay = describeResp.WorkflowExecutionInfo.GetExecutionTime().Sub(timeSentReq)
	} else {
		serviceClient, closeFunc, err := iwf.BuildCadenceServiceClient("localhost:7833")
		if err != nil {
			rawLog.Fatalf("Unable to connect to Cadence because of error %v", err)
		}
		defer closeFunc()

		cadenceClient, err := iwf.BuildCadenceClient(serviceClient, "default")
		if err != nil {
			rawLog.Fatalf("Unable to connect to Cadence because of error %v", err)
		}

		describeResp, err := cadenceClient.DescribeWorkflowExecution(context.Background(), wfId, *resp.WorkflowRunId)
		if err != nil {
			t.Fatal(err)
		}
		if describeResp.GetWorkflowExecutionInfo().ExecutionTime == nil {
			t.Fatal("start time is nil")
		}

		delay = time.Unix(0, describeResp.GetWorkflowExecutionInfo().GetExecutionTime()).Sub(timeSentReq)
	}

	assertions := assert.New(t)
	assertions.True(delay.Seconds() > 8, "delay is %v", delay)
	assertions.True(delay.Seconds() < 12, "delay is %v", delay)
}
