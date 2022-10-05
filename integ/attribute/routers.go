package attribute

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	WorkflowType               = "attribute"
	State1                     = "S1"
	State2                     = "S2"
	TestQueryAttributeKey      = "test-query-attribute"
	TestStateLocalAttributeKey = "test-state-local-attribute"

	// Here use builtin search attribute for testing
	// NOTE: SA is not supported in temporalite: https://github.com/temporalio/temporalite/pull/28
	EnableTestingSearchAttribute     = false
	TestSearchAttributeKeywordKey    = "CustomKeywordField"
	TestSearchAttributeKeywordValue1 = "keyword-value1"
	TestSearchAttributeKeywordValue2 = "keyword-value2"
	TestSearchAttributeIntKey        = "CustomIntField"
	TestSearchAttributeIntValue1     = "1"
	TestSearchAttributeIntValue2     = "2"
)

var TestQueryVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-query-value1"),
}

var TestQueryVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-query-value2"),
}

var testStateLocalAttributeVal = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-state-local-value"),
}

func NewAttributeWorkflow() (*Handler, *gin.Engine) {
	router := gin.Default()

	handler := newHandler()

	router.POST(service.StateStartApi, handler.apiV1WorkflowStateStart)
	router.POST(service.StateDecideApi, handler.apiV1WorkflowStateDecide)

	return handler, router
}

type Handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func newHandler() *Handler {
	return &Handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *Handler) apiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			var sa []iwfidl.SearchAttribute
			if EnableTestingSearchAttribute {
				sa = []iwfidl.SearchAttribute{
					{
						Key:       iwfidl.PtrString(TestSearchAttributeKeywordKey),
						Value:     iwfidl.PtrString(TestSearchAttributeKeywordValue1),
						ValueType: iwfidl.PtrString(service.SearchAttributeValueTypeKeyword),
					},
					{
						Key:       iwfidl.PtrString(TestSearchAttributeIntKey),
						Value:     iwfidl.PtrString(TestSearchAttributeIntValue1),
						ValueType: iwfidl.PtrString(service.SearchAttributeValueTypeInt),
					},
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
				},
				UpsertQueryAttributes: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestQueryAttributeKey),
						Value: &TestQueryVal1,
					},
				},
				UpsertSearchAttributes: sa,
				UpsertStateLocalAttributes: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestStateLocalAttributeKey),
						Value: &testStateLocalAttributeVal,
					},
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State2 {
			sas := req.GetSearchAttributes()
			kwSaFounds := 0
			intSaFounds := 0
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetValue() == TestSearchAttributeKeywordValue2 && sa.GetValueType() == service.SearchAttributeValueTypeKeyword {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetValue() == TestSearchAttributeIntValue2 && sa.GetValueType() == service.SearchAttributeValueTypeInt {
					intSaFounds++
				}
			}
			h.invokeData["S2_start_kwSaFounds"] = kwSaFounds
			h.invokeData["S2_start_intSaFounds"] = intSaFounds

			queryAttFound := false
			queryAtt := req.GetQueryAttributes()[0]
			value := queryAtt.GetValue()
			if queryAtt.GetKey() == TestQueryAttributeKey && value.GetData() == TestQueryVal2.GetData() && value.GetEncoding() == TestQueryVal2.GetEncoding() {
				queryAttFound = true
			}
			h.invokeData["S2_start_queryAttFound"] = queryAttFound

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) apiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		if req.GetWorkflowStateId() == State1 {
			sas := req.GetSearchAttributes()
			kwSaFounds := 0
			intSaFounds := 0
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetValue() == TestSearchAttributeKeywordValue1 && sa.GetValueType() == service.SearchAttributeValueTypeKeyword {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetValue() == TestSearchAttributeIntValue1 && sa.GetValueType() == service.SearchAttributeValueTypeInt {
					intSaFounds++
				}
			}
			h.invokeData["S1_decide_kwSaFounds"] = kwSaFounds
			h.invokeData["S1_decide_intSaFounds"] = intSaFounds

			queryAttFound := false
			queryAtt := req.GetQueryAttributes()[0]
			value := queryAtt.GetValue()
			if queryAtt.GetKey() == TestQueryAttributeKey && value.GetData() == TestQueryVal1.GetData() && value.GetEncoding() == TestQueryVal1.GetEncoding() {
				queryAttFound = true
			}
			h.invokeData["S1_decide_queryAttFound"] = queryAttFound

			localAttFound := false
			localAtt := req.GetStateLocalAttributes()[0]
			value = localAtt.GetValue()
			if localAtt.GetKey() == TestStateLocalAttributeKey && value.GetData() == testStateLocalAttributeVal.GetData() && value.GetEncoding() == testStateLocalAttributeVal.GetEncoding() {
				localAttFound = true
			}
			h.invokeData["S1_decide_localAttFound"] = localAttFound

			var sa []iwfidl.SearchAttribute
			if EnableTestingSearchAttribute {
				sa = []iwfidl.SearchAttribute{
					{
						Key:       iwfidl.PtrString(TestSearchAttributeKeywordKey),
						Value:     iwfidl.PtrString(TestSearchAttributeKeywordValue2),
						ValueType: iwfidl.PtrString(service.SearchAttributeValueTypeKeyword),
					},
					{
						Key:       iwfidl.PtrString(TestSearchAttributeIntKey),
						Value:     iwfidl.PtrString(TestSearchAttributeIntValue2),
						ValueType: iwfidl.PtrString(service.SearchAttributeValueTypeInt),
					},
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State2,
						},
					},
					UpsertQueryAttributes: []iwfidl.KeyValue{
						{
							Key:   iwfidl.PtrString(TestQueryAttributeKey),
							Value: &TestQueryVal2,
						},
					},
					UpsertSearchAttributes: sa,
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			sas := req.GetSearchAttributes()
			kwSaFounds := 0
			intSaFounds := 0
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetValue() == TestSearchAttributeKeywordValue2 && sa.GetValueType() == service.SearchAttributeValueTypeKeyword {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetValue() == TestSearchAttributeIntValue2 && sa.GetValueType() == service.SearchAttributeValueTypeInt {
					intSaFounds++
				}
			}
			h.invokeData["S2_decide_kwSaFounds"] = kwSaFounds
			h.invokeData["S2_decide_intSaFounds"] = intSaFounds

			queryAttFound := false
			queryAtt := req.GetQueryAttributes()[0]
			value := queryAtt.GetValue()
			if queryAtt.GetKey() == TestQueryAttributeKey && value.GetData() == TestQueryVal2.GetData() && value.GetEncoding() == TestQueryVal2.GetEncoding() {
				queryAttFound = true
			}
			h.invokeData["S2_decide_queryAttFound"] = queryAttFound

			// go to complete
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.GracefulCompletingWorkflowStateId,
						},
					},
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
