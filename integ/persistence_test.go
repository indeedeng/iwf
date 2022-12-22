package integ

import (
	"context"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
)

func TestPersistenceWorkflowTemporal(t *testing.T) {
	if !*temporalIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestPersistenceWorkflow(t, service.BackendTypeTemporal)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func TestPersistenceWorkflowCadence(t *testing.T) {
	if !*cadenceIntegTest {
		t.Skip()
	}
	for i := 0; i < *repeatIntegTest; i++ {
		doTestPersistenceWorkflow(t, service.BackendTypeCadence)
		time.Sleep(time.Millisecond * time.Duration(*repeatInterval))
	}
}

func doTestPersistenceWorkflow(t *testing.T, backendType service.BackendType) {
	wfHandler := persistence.NewHandler()
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
	wfId := persistence.WorkflowType + strconv.Itoa(int(time.Now().Unix()))
	reqStart := apiClient.DefaultApi.ApiV1WorkflowStartPost(context.Background())
	_, httpResp, err := reqStart.WorkflowStartRequest(iwfidl.WorkflowStartRequest{
		WorkflowId:             wfId,
		IwfWorkflowType:        persistence.WorkflowType,
		WorkflowTimeoutSeconds: 10,
		IwfWorkerUrl:           "http://localhost:" + testWorkflowServerPort,
		StartStateId:           persistence.State1,
		StateOptions: &iwfidl.WorkflowStateOptions{
			SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
				PersistenceLoadingType: ptr.Any(iwfidl.ALL_WITHOUT_LOCKING),
			},
			DataObjectsLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
				PersistenceLoadingType: ptr.Any(iwfidl.ALL_WITHOUT_LOCKING),
			},
		},
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Status not success" + httpResp.Status)
	}

	reqWait := apiClient.DefaultApi.ApiV1WorkflowGetWithWaitPost(context.Background())
	wfResponse, httpResp, err := reqWait.WorkflowGetRequest(iwfidl.WorkflowGetRequest{
		WorkflowId: wfId,
	}).Execute()
	if err != nil {
		log.Fatalf("Fail to invoke start api %v", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Fatalf("Fail to get workflow" + httpResp.Status)
	}

	reqQry := apiClient.DefaultApi.ApiV1WorkflowDataobjectsGetPost(context.Background())
	queryResult1, httpResp2, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
		Keys: []string{
			persistence.TestDataObjectKey,
		},
	}).Execute()

	if err != nil || httpResp2.StatusCode != 200 {
		log.Fatalf("Fail to invoke query workflow for sigle attr %v %v", err, httpResp2)
	}

	queryResult2, httpResp2, err := reqQry.WorkflowGetDataObjectsRequest(iwfidl.WorkflowGetDataObjectsRequest{
		WorkflowId: wfId,
	}).Execute()

	if err != nil || httpResp2.StatusCode != 200 {
		log.Fatalf("Fail to invoke query workflow for sigle attr %v %v", err, httpResp2)
	}

	reqSearch := apiClient.DefaultApi.ApiV1WorkflowSearchattributesGetPost(context.Background())
	searchResult1, httpSearchResponse1, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId:    wfId,
		WorkflowRunId: &wfResponse.WorkflowRunId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType: ptr.Any(iwfidl.KEYWORD),
			},
		},
	}).Execute()

	if err != nil || httpSearchResponse1.StatusCode != 200 {
		log.Fatalf("Fail to invoke query workflow for sigle attr %v %v", err, httpSearchResponse1)
	}

	searchResult2, httpSearchResponse2, err := reqSearch.WorkflowGetSearchAttributesRequest(iwfidl.WorkflowGetSearchAttributesRequest{
		WorkflowId:    wfId,
		WorkflowRunId: &wfResponse.WorkflowRunId,
		Keys: []iwfidl.SearchAttributeKeyAndType{
			{
				Key:       iwfidl.PtrString(persistence.TestSearchAttributeIntKey),
				ValueType: ptr.Any(iwfidl.INT),
			},
		},
	}).Execute()

	if err != nil || httpSearchResponse2.StatusCode != 200 {
		log.Fatalf("Fail to invoke query workflow for sigle attr %v %v", err, httpSearchResponse2)
	}

	// assertion
	history, data := wfHandler.GetTestResult()
	assertions := assert.New(t)
	assertions.Equalf(map[string]int64{
		"S1_start":  1,
		"S1_decide": 1,
		"S2_start":  1,
		"S2_decide": 1,
		"S3_start":  1,
		"S3_decide": 1,
	}, history, "attribute test fail, %v", history)

	if persistence.EnableTestingSearchAttribute {
		assertions.Equal(map[string]interface{}{
			"S1_decide_intSaFounds": 1,
			"S1_decide_kwSaFounds":  1,
			"S2_decide_intSaFounds": 1,
			"S2_decide_kwSaFounds":  1,
			"S2_start_intSaFounds":  1,
			"S2_start_kwSaFounds":   1,

			"S1_decide_localAttFound": true,
			"S1_decide_queryAttFound": 2,
			"S2_decide_queryAttFound": true,
			"S2_start_queryAttFound":  true,
		}, data)
	} else {
		// map[S1_decide_intSaFounds:0 S1_decide_kwSaFounds:0 S1_decide_localAttFound:false
		//S1_decide_queryAttFound:true S2_decide_intSaFounds:0 S2_decide_kwSaFounds:0
		//S2_decide_queryAttFound:false S2_start_intSaFounds:0 S2_start_kwSaFounds:0 S2_start_queryAttFound:false]
		assertions.Equal(map[string]interface{}{
			"S1_decide_intSaFounds": 0,
			"S1_decide_kwSaFounds":  0,
			"S2_decide_intSaFounds": 0,
			"S2_decide_kwSaFounds":  0,
			"S2_start_intSaFounds":  0,
			"S2_start_kwSaFounds":   0,

			"S1_decide_localAttFound": true,
			"S1_decide_queryAttFound": 2,
			"S2_decide_queryAttFound": true,
			"S2_start_queryAttFound":  true,
		}, data)
	}

	expected1 := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey),
			Value: &persistence.TestDataObjectVal2,
		},
	}
	expected2 := []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey),
			Value: &persistence.TestDataObjectVal2,
		},
		{
			Key:   iwfidl.PtrString(persistence.TestDataObjectKey2),
			Value: &persistence.TestDataObjectVal1,
		},
	}
	assertions.ElementsMatch(expected1, queryResult1.GetObjects())
	assertions.ElementsMatch(expected2, queryResult2.GetObjects())

	expectedSearchAttributeInt := iwfidl.SearchAttribute{
		Key:          iwfidl.PtrString(persistence.TestSearchAttributeIntKey),
		ValueType:    ptr.Any(iwfidl.INT),
		IntegerValue: iwfidl.PtrInt64(persistence.TestSearchAttributeIntValue2),
	}

	expectedSearchAttributeKeyword := iwfidl.SearchAttribute{
		Key:         iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
		ValueType:   ptr.Any(iwfidl.KEYWORD),
		StringValue: iwfidl.PtrString(persistence.TestSearchAttributeKeywordValue2),
	}

	assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeKeyword}, searchResult1.GetSearchAttributes())
	assertions.Equal([]iwfidl.SearchAttribute{expectedSearchAttributeInt}, searchResult2.GetSearchAttributes())
}
