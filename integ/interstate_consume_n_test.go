package integ

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/interstate_consume_n"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestInterStateConsumeNWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestInterStateConsumeNWorkflow(t, service.BackendTypeTemporal, nil)
		smallWaitForFastTest()
		doTestInterStateConsumeNWorkflow(t, service.BackendTypeTemporal, minimumContinueAsNewConfig(true))
		smallWaitForFastTest()
	}
}

func TestInterStateConsumeNWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestInterStateConsumeNWorkflow(t, service.BackendTypeCadence, nil)
		smallWaitForFastTest()
		doTestInterStateConsumeNWorkflow(t, service.BackendTypeCadence, minimumContinueAsNewConfig(false))
		smallWaitForFastTest()
	}
}

func doTestInterStateConsumeNWorkflow(t *testing.T, backendType service.BackendType, config *iwfidl.WorkflowConfig) {
	wfHandler := interstate_consume_n.NewHandler()
	closeFunc1 := startWorkflowWorker(wfHandler, t)
	defer closeFunc1()

	closeFunc2 := startIwfService(backendType)
	defer closeFunc2()

	apiClient := iwfidl.NewAPIClient(&iwfidl.Configuration{
		Servers: []iwfidl.ServerConfiguration{
			{
				URL: "http://localhost:" + testIwfServerPort,
			},
		},
	})

	wfId := interstate_consume_n.WorkflowType + strconv.Itoa(int(time.Now().UnixNano()))
	req := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := req.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        interstate_consume_n.WorkflowType,
		WorkflowTimeoutSeconds: 20,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           ptr.Any(interstate_consume_n.State1),
		WorkflowStartOptions: &iwfidl.WorkflowStartOptions{
			WorkflowConfigOverride: config,
		},
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	req2 := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	resp2, httpResp, err := req2.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	failTestAtHttpError(err, httpResp, t)

	assertions := assert.New(t)
	assertions.Equal(iwfidl.COMPLETED, resp2.GetWorkflowStatus())

	history, data := wfHandler.GetTestResult()
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
		"S3_start":  1,
		"S3_decide": 1,
		"S4_start":  1,
		"S4_decide": 1,
		"S5_start":  1,
		"S5_decide": 1,
		"S6_start":  1,
		"S6_decide": 1,
	}, history, "consume N test fail, %v", history)

	// ExactN (AtLeast=3, AtMost=3): should consume exactly 3 of the 5 published messages
	exactNValues := data["exactN_values"].([]iwfidl.EncodedObject)
	assertions.Equal(3, len(exactNValues), "ExactN should consume exactly 3 messages")
	assertions.Equal(*interstate_consume_n.TestValues[0].Data, *exactNValues[0].Data)
	assertions.Equal(*interstate_consume_n.TestValues[1].Data, *exactNValues[1].Data)
	assertions.Equal(*interstate_consume_n.TestValues[2].Data, *exactNValues[2].Data)

	// OneToAll (AtLeast=1, no AtMost): should consume all remaining 2 messages
	oneToAllValues := data["oneToAll_values"].([]iwfidl.EncodedObject)
	assertions.Equal(2, len(oneToAllValues), "OneToAll should consume all remaining messages")
	assertions.Equal(*interstate_consume_n.TestValues[3].Data, *oneToAllValues[0].Data)
	assertions.Equal(*interstate_consume_n.TestValues[4].Data, *oneToAllValues[1].Data)

	// ZeroToAll (AtLeast=0, no AtMost): channel empty, should consume 0 messages
	zeroToAllValues := data["zeroToAll_values"]
	if zeroToAllValues == nil {
		assertions.Nil(zeroToAllValues)
	} else {
		vals := zeroToAllValues.([]iwfidl.EncodedObject)
		assertions.Equal(0, len(vals), "ZeroToAll should consume 0 messages from empty channel")
	}

	// AtMostOnly (AtMost=2, no AtLeast): waits for late messages from S6, should consume 2 of 3
	atMostOnlyValues := data["atMostOnly_values"].([]iwfidl.EncodedObject)
	assertions.Equal(2, len(atMostOnlyValues), "AtMostOnly should consume exactly 2 messages")
	assertions.Equal(*interstate_consume_n.TestValuesCh2[0].Data, *atMostOnlyValues[0].Data)
	assertions.Equal(*interstate_consume_n.TestValuesCh2[1].Data, *atMostOnlyValues[1].Data)
}
