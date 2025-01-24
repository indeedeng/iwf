package integ

import (
	"context"
	"fmt"
	"github.com/indeedeng/iwf/service/common/ptr"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/wait_until_search_attributes"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestWaitUntilSearchAttributesWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilSearchAttributes(t, service.BackendTypeTemporal, &iwfidl.WorkflowConfig{
			ExecutingStateIdMode: ptr.Any(iwfidl.DISABLED),
		})
		smallWaitForFastTest()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilSearchAttributes(t, service.BackendTypeTemporal, &iwfidl.WorkflowConfig{
			ExecutingStateIdMode: ptr.Any(iwfidl.ENABLED_FOR_ALL),
		})
		smallWaitForFastTest()
	}

	for i := 0; i < *repeatIntegTest; i++ {
		doTestWaitUntilSearchAttributes(t, service.BackendTypeTemporal, nil) // defaults to ExecutingStateIdMode: ENABLED_FOR_STATES_WITH_WAIT_UNTIL
		smallWaitForFastTest()
	}
}

func doTestWaitUntilSearchAttributes(
	t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig,
) {
	assertions := assert.New(t)
	wfHandler := wait_until_search_attributes.NewHandler()
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
	wfId := wait_until_search_attributes.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))

	reqStart := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	wfReq := iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        wait_until_search_attributes.WorkflowType,
		WorkflowTimeoutSeconds: 20,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(wait_until_search_attributes.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}
	_, httpResp, err := reqStart.WorkflowStartRequest(wfReq).Execute()
	panicAtHttpError(err, httpResp, t)

	// wait for the search attribute index to be ready in ElasticSearch
	time.Sleep(time.Duration(*searchWaitTimeIntegTest) * time.Millisecond)

	switch mode := config.GetExecutingStateIdMode(); mode {
	case iwfidl.ENABLED_FOR_ALL:
		assertSearch(fmt.Sprintf("WorkflowId='%v'", wfId), 1, apiClient, assertions)
		assertSearch(fmt.Sprintf("WorkflowId='%v' AND %v='%v'", wfId, wait_until_search_attributes.TestSearchAttributeExecutingStateIdsKey, wait_until_search_attributes.State2), 1, apiClient, assertions)
	case iwfidl.ENABLED_FOR_STATES_WITH_WAIT_UNTIL:
		assertSearch(fmt.Sprintf("WorkflowId='%v'", wfId), 1, apiClient, assertions)
		assertSearch(fmt.Sprintf("WorkflowId='%v' AND %v='%v'", wfId, wait_until_search_attributes.TestSearchAttributeExecutingStateIdsKey, wait_until_search_attributes.State2), 0, apiClient, assertions)
	case iwfidl.DISABLED:
		assertSearch(fmt.Sprintf("WorkflowId='%v'", wfId), 1, apiClient, assertions)
		assertSearch(fmt.Sprintf("WorkflowId='%v' AND %v='%v'", wfId, wait_until_search_attributes.TestSearchAttributeExecutingStateIdsKey, wait_until_search_attributes.State2), 0, apiClient, assertions)
	}

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	_, httpResp, err = reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpError(err, httpResp, t)

	// wait for workflow to complete
	resp, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	panicAtHttpErrorOrWorkflowUncompleted(err, httpResp, resp, t)
}
