package persistence

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"sync"
	"testing"
)

/**
 * This test workflow has three states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method will update DA, SA, & SL
 * 		- Execute method will move to State2 with partially loaded data
 * State2:
 * 		- WaitUntil method will store attribute data
 * 		- Execute method will move to State3 with partially loaded data
 * State3:
 * 		- WaitUntil method performs some attribute checks
 * 		- Execute method performs checks on the attribute data and then gracefully completes the workflow
 */
const (
	WorkflowType          = "persistence"
	State1                = "S1"
	State2                = "S2"
	State3                = "S3"
	TestDataAttributeKey  = "test-data-object"
	TestDataAttributeKey2 = "test-data-object-2"
	TestStateLocalKey     = "test-state-local"

	TestSearchAttributeKeywordKey    = "CustomKeywordField"
	TestSearchAttributeKeywordValue1 = "keyword-value1"
	TestSearchAttributeKeywordValue2 = "keyword-value2"

	TestSearchAttributeKeywordArrayKey = "CustomKeywordArrayField"
	TestSearchAttributeIntKey          = "CustomIntField"
	TestSearchAttributeBoolKey         = "CustomBoolField"
	TestSearchAttributeDoubleKey       = "CustomDoubleField"
	TestSearchAttributeDatetimeKey     = "CustomDatetimeField"
	TestSearchAttributeTextKey         = "CustomStringField"
	TestSearchAttributeIntValue1       = 1
	TestSearchAttributeIntValue2       = 2
)

var TestDataAttributeVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-object-value1"),
}

var TestDataAttributeVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-object-value2"),
}

var testStateLocalVal = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-state-local-value"),
}

type handler struct {
	invokeHistory sync.Map
	invokeData    sync.Map
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: sync.Map{},
		invokeData:    sync.Map{},
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	initSas := req.GetSearchAttributes()
	if len(initSas) < 1 {
		helpers.FailTestWithErrorMessage("should have at least one init search attribute", t)
	}
	for _, sa := range initSas {
		if sa.GetKey() == "CustomDatetimeField" {
			if sa.GetValueType() != iwfidl.DATETIME {
				helpers.FailTestWithErrorMessage("key and value type not match", t)
			}
		}
	}

	if req.GetWorkflowType() == WorkflowType {
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
		}

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

			// Go to the decide methods after updating DA, SA, & SL
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestDataAttributeKey),
						Value: &TestDataAttributeVal1,
					},
					{
						Key:   iwfidl.PtrString(TestDataAttributeKey2),
						Value: &TestDataAttributeVal1,
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

			// Determine how many keywords and ints are found in the search attributes
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
			h.invokeData.Store("S2_start_kwSaFounds", kwSaFounds)
			h.invokeData.Store("S2_start_intSaFounds", intSaFounds)

			// Determine if the attribute is found in the request
			queryAttFound := false
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataAttributeKey && value.GetData() == TestDataAttributeVal2.GetData() && value.GetEncoding() == TestDataAttributeVal2.GetEncoding() {
					queryAttFound = true
				}
				if queryAtt.GetKey() == TestDataAttributeKey2 {
					helpers.FailTestWithErrorMessage("should not load key that is not included in partial loading", t)
				}
			}
			h.invokeData.Store("S2_start_queryAttFound", queryAttFound)

			// Go straight to the decide methods without any commands
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State3 {
			sas := req.GetSearchAttributes()

			// Determine if the INT attribute is found in the request
			found := false
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey {
					helpers.FailTestWithErrorMessage("should not load key that is not included in partial loading", t)
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					found = true
				}
			}
			if !found {
				helpers.FailTestWithErrorMessage("should see the requested partial loading key", t)
			}

			queryAttFound := 0
			queryAtts := req.GetDataObjects()
			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataAttributeKey {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataAttributeKey2 {
					queryAttFound++
				}
			}
			if queryAttFound != 2 {
				helpers.FailTestWithErrorMessage("missing query attribute requested by partial loading keys", t)
			}

			// Go straight to the decide methods without any commands
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

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_decide"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", int64(1))
		}

		if req.GetWorkflowStateId() == State1 {
			sas := req.GetSearchAttributes()

			// Determine how many keywords and ints are found in the search attributes
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
			h.invokeData.Store("S1_decide_kwSaFounds", kwSaFounds)
			h.invokeData.Store("S1_decide_intSaFounds", intSaFounds)

			queryAttFound := 0
			queryAtts := req.GetDataObjects()

			// Determine how many query attributes are found
			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataAttributeKey && value.GetData() == TestDataAttributeVal1.GetData() && value.GetEncoding() == TestDataAttributeVal1.GetEncoding() {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataAttributeKey2 && value.GetData() == TestDataAttributeVal1.GetData() && value.GetEncoding() == TestDataAttributeVal1.GetEncoding() {
					queryAttFound++
				}
			}
			h.invokeData.Store("S1_decide_queryAttFound", queryAttFound)

			// Determine if local attribute is found
			localAttFound := false
			localAtt := req.GetStateLocals()[0]
			value := localAtt.GetValue()
			if localAtt.GetKey() == TestStateLocalKey && value.GetData() == testStateLocalVal.GetData() && value.GetEncoding() == testStateLocalVal.GetEncoding() {
				localAttFound = true
			}
			h.invokeData.Store("S1_decide_localAttFound", localAttFound)

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

			// Move to state 2 with set options after updating values
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
								DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
									PersistenceLoadingType: ptr.Any(iwfidl.PARTIAL_WITHOUT_LOCKING),
									PartialLoadingKeys: []string{
										TestDataAttributeKey,
									},
								},
							},
						},
					},
				},
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestDataAttributeKey),
						Value: &TestDataAttributeVal2,
					},
				},
				UpsertSearchAttributes: sa,
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			sas := req.GetSearchAttributes()

			// Determine how many keywords and ints are found in the search attributes
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
			h.invokeData.Store("S2_decide_kwSaFounds", kwSaFounds)
			h.invokeData.Store("S2_decide_intSaFounds", intSaFounds)

			queryAttFound := false
			queryAtts := req.GetDataObjects()

			// Determine how many query attributes are found
			for _, queryAtt := range queryAtts {
				value := queryAtt.GetValue()
				if queryAtt.GetKey() == TestDataAttributeKey && value.GetData() == TestDataAttributeVal2.GetData() && value.GetEncoding() == TestDataAttributeVal2.GetEncoding() {
					queryAttFound = true
				}
				if queryAtt.GetKey() == TestDataAttributeKey2 {
					helpers.FailTestWithErrorMessage("should not load key that is not included in partial loading", t)
				}
			}

			h.invokeData.Store("S2_decide_queryAttFound", queryAttFound)

			// Move to state 3 after with set options
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
										TestDataAttributeKey,
										TestDataAttributeKey2,
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

			// Determine if the INT attribute is found in the request
			found := false
			for _, sa := range sas {
				if sa.GetKey() == TestSearchAttributeKeywordKey {
					helpers.FailTestWithErrorMessage("should not load key that is not included in partial loading", t)
				}
				if sa.GetKey() == TestSearchAttributeIntKey && sa.GetIntegerValue() == TestSearchAttributeIntValue2 &&
					sa.GetValueType() == iwfidl.INT {
					found = true
				}
			}
			if !found {
				helpers.FailTestWithErrorMessage("should see the requested partial loading key", t)
			}

			queryAttFound := 0
			queryAtts := req.GetDataObjects()

			// Determine how many query attributes are found
			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataAttributeKey {
					queryAttFound++
				}
				if queryAtt.GetKey() == TestDataAttributeKey2 {
					queryAttFound++
				}
			}
			if queryAttFound != 2 {
				helpers.FailTestWithErrorMessage("missing query attribute requested by partial loading keys", t)
			}

			// Move to completion
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
	invokeHistory := make(map[string]int64)
	h.invokeHistory.Range(func(key, value interface{}) bool {
		invokeHistory[key.(string)] = value.(int64)
		return true
	})
	invokeData := make(map[string]interface{})
	h.invokeData.Range(func(key, value interface{}) bool {
		invokeData[key.(string)] = value
		return true
	})
	return invokeHistory, invokeData
}
