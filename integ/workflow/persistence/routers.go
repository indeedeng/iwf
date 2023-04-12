package persistence

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
)

const (
	WorkflowType       = "persistence"
	State1             = "S1"
	State2             = "S2"
	State3             = "S3"
	TestDataObjectKey  = "test-data-object"
	TestDataObjectKey2 = "test-data-object-2"
	TestStateLocalKey  = "test-state-local"

	TestSearchAttributeKeywordKey    = "CustomKeywordField"
	TestSearchAttributeKeywordValue1 = "keyword-value1"
	TestSearchAttributeKeywordValue2 = "keyword-value2"

	TestSearchAttributeIntKey      = "CustomIntField"
	TestSearchAttributeBoolKey     = "CustomBoolField"
	TestSearchAttributeDoubleKey   = "CustomDoubleField"
	TestSearchAttributeDatetimeKey = "CustomDatetimeField"
	TestSearchAttributeTextKey     = "CustomStringField"
	TestSearchAttributeIntValue1   = 1
	TestSearchAttributeIntValue2   = 2
)

var TestDataObjectVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-object-value1"),
}

var TestDataObjectVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-object-value2"),
}

var testStateLocalVal = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-state-local-value"),
}

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	initSas := req.GetSearchAttributes()
	if len(initSas) < 1 {
		panic("should have at least one init search attribute")
	}
	for _, sa := range initSas {
		if sa.GetKey() == "CustomDatetimeField" {
			if sa.GetValueType() != iwfidl.DATETIME {
				panic("key and value type not match")
			}
		}
	}

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			var sa []iwfidl.SearchAttribute
			sa = []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
					StringValue: iwfidl.PtrString(TestSearchAttributeKeywordValue1),
					ValueType:   ptr.Any(iwfidl.KEYWORD),
				},
				{
					Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
					IntegerValue: iwfidl.PtrInt64(TestSearchAttributeIntValue1),
					ValueType:    ptr.Any(iwfidl.INT),
				},
				{
					Key:       iwfidl.PtrString(TestSearchAttributeBoolKey),
					ValueType: ptr.Any(iwfidl.BOOL),
					BoolValue: iwfidl.PtrBool(false),
				},
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestDataObjectKey),
						Value: &TestDataObjectVal1,
					},
					{
						Key:   iwfidl.PtrString(TestDataObjectKey2),
						Value: &TestDataObjectVal1,
					},
				},
				UpsertSearchAttributes: sa,
				UpsertStateLocals: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestStateLocalKey),
						Value: &testStateLocalVal,
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
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetStringValue() == TestSearchAttributeKeywordValue2 &&
					sa.GetValueType() == iwfidl.KEYWORD {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					intSaFounds++
				}
			}
			h.invokeData["S2_start_kwSaFounds"] = kwSaFounds
			h.invokeData["S2_start_intSaFounds"] = intSaFounds

			queryAttFound := false
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataObjectKey && value.GetData() == TestDataObjectVal2.GetData() && value.GetEncoding() == TestDataObjectVal2.GetEncoding() {
					queryAttFound = true
				}
				if queryAtt.GetKey() == TestDataObjectKey2 {
					panic("should not load key that is not included in partial loading")
				}
			}
			h.invokeData["S2_start_queryAttFound"] = queryAttFound

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State3 {
			sas := req.GetSearchAttributes()
			found := false
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey {
					panic("should not load key that is not included in partial loading")
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					found = true
				}
			}
			if !found {
				panic("should see the requested partial loading key")
			}

			queryAttFound := 0
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataObjectKey {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataObjectKey2 {
					queryAttFound++
				}
			}
			if queryAttFound != 2 {
				panic("missing query attribute requested by partial loading keys")
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
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
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetStringValue() == TestSearchAttributeKeywordValue1 &&
					sa.GetValueType() == iwfidl.KEYWORD {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue1 &&
					sa.GetValueType() == iwfidl.INT {
					intSaFounds++
				}
			}
			h.invokeData["S1_decide_kwSaFounds"] = kwSaFounds
			h.invokeData["S1_decide_intSaFounds"] = intSaFounds

			queryAttFound := 0
			queryAtts := req.GetDataObjects()

			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataObjectKey && value.GetData() == TestDataObjectVal1.GetData() && value.GetEncoding() == TestDataObjectVal1.GetEncoding() {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataObjectKey2 && value.GetData() == TestDataObjectVal1.GetData() && value.GetEncoding() == TestDataObjectVal1.GetEncoding() {
					queryAttFound++
				}
			}
			h.invokeData["S1_decide_queryAttFound"] = queryAttFound

			localAttFound := false
			localAtt := req.GetStateLocals()[0]
			value := localAtt.GetValue()
			if localAtt.GetKey() == TestStateLocalKey && value.GetData() == testStateLocalVal.GetData() && value.GetEncoding() == testStateLocalVal.GetEncoding() {
				localAttFound = true
			}
			h.invokeData["S1_decide_localAttFound"] = localAttFound

			var sa []iwfidl.SearchAttribute
			sa = []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
					StringValue: iwfidl.PtrString(TestSearchAttributeKeywordValue2),
					ValueType:   ptr.Any(iwfidl.KEYWORD),
				},
				{
					Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
					IntegerValue: iwfidl.PtrInt64(TestSearchAttributeIntValue2),
					ValueType:    ptr.Any(iwfidl.INT),
				},
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State2,
							StateOptions: &iwfidl.WorkflowStateOptions{
								SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
									PersistenceLoadingType: ptr.Any(iwfidl.PARTIAL_WITHOUT_LOCKING),
									PartialLoadingKeys: []string{
										TestSearchAttributeIntKey,
										TestSearchAttributeKeywordKey,
									},
								},
								DataObjectsLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
									PersistenceLoadingType: ptr.Any(iwfidl.PARTIAL_WITHOUT_LOCKING),
									PartialLoadingKeys: []string{
										TestDataObjectKey,
									},
								},
							},
						},
					},
				},
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestDataObjectKey),
						Value: &TestDataObjectVal2,
					},
				},
				UpsertSearchAttributes: sa,
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			sas := req.GetSearchAttributes()
			kwSaFounds := 0
			intSaFounds := 0
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey && sa.GetStringValue() == TestSearchAttributeKeywordValue2 &&
					sa.GetValueType() == iwfidl.KEYWORD {
					kwSaFounds++
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					intSaFounds++
				}
			}
			h.invokeData["S2_decide_kwSaFounds"] = kwSaFounds
			h.invokeData["S2_decide_intSaFounds"] = intSaFounds

			queryAttFound := false
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataObjectKey && value.GetData() == TestDataObjectVal2.GetData() && value.GetEncoding() == TestDataObjectVal2.GetEncoding() {
					queryAttFound = true
				}
				if queryAtt.GetKey() == TestDataObjectKey2 {
					panic("should not load key that is not included in partial loading")
				}
			}

			h.invokeData["S2_decide_queryAttFound"] = queryAttFound

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State3,
							StateOptions: &iwfidl.WorkflowStateOptions{
								SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
									PersistenceLoadingType: ptr.Any(iwfidl.PARTIAL_WITHOUT_LOCKING),
									PartialLoadingKeys: []string{
										TestSearchAttributeIntKey,
									},
								},
								DataObjectsLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
									PersistenceLoadingType: ptr.Any(iwfidl.PARTIAL_WITHOUT_LOCKING),
									PartialLoadingKeys: []string{
										TestDataObjectKey,
										TestDataObjectKey2,
									},
								},
							},
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State3 {
			sas := req.GetSearchAttributes()
			found := false
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey {
					panic("should not load key that is not included in partial loading")
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					found = true
				}
			}
			if !found {
				panic("should see the requested partial loading key")
			}

			queryAttFound := 0
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataObjectKey {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataObjectKey2 {
					queryAttFound++
				}
			}
			if queryAttFound != 2 {
				panic("missing query attribute requested by partial loading keys")
			}

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

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
