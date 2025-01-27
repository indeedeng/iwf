package integ

import (
	"context"
	"encoding/json"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/signal"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
	"time"
)

func TestWorkflowGetWithWaitTimeoutTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithWaitTimeout(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
	}
}

func TestWorkflowWithWaitTimeoutCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWorkflowWithWaitTimeout(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
	}
}

func doTestWorkflowWithWaitTimeout(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	assertions := assert.New(t)

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
	wfId := "wf-wait-timeout-test" + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        signal.WorkflowType,
		WorkflowTimeoutSeconds: 15,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(signal.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	// wait for the workflow
	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	startTimeUnix := time.Now().Unix()
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	elapsedSeconds := time.Now().Unix() - startTimeUnix

	assertions.NotNil(err)
	assertions.Equalf(service.HttpStatusCodeSpecial4xxError1, httpResp.StatusCode, "http code")
	var errResp iwfidl.ErrorResponse
	body, err := ioutil.ReadAll(httpResp.Body)
	assertions.Nil(err)
	err = json.Unmarshal(body, &errResp)
	assertions.Equalf(iwfidl.ErrorResponse{
		Detail:    ptr.Any("workflow is still running, waiting has exceeded timeout limit, please retry"),
		SubStatus: iwfidl.LONG_POLL_TIME_OUT_SUB_STATUS.Ptr(),
	}, errResp, "body")

	assertions.True(elapsedSeconds >= 5 && elapsedSeconds <= 12, "expect to poll for ~8 seconds, actual value is ", elapsedSeconds)
}
