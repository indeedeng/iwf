package s3_init_data_attributes

import (
	"log"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

/**
 * This test workflow has 2 states, testing S3 data attribute loading functionality.
 *
 * State1:
 *		- WaitUntil method loads and validates data attributes from S3
 *      - Execute method transitions to State2
 *
 * State2:
 *		- WaitUntil method does nothing
 *      - Execute method loads and validates data attributes from S3, then completes workflow
 */
const (
	WorkflowType      = "s3-init-data-attributes"
	State1            = "S1"
	State2            = "S2"
	TestDataAttrKey1  = "test-da-key1"
	TestDataAttrKey2  = "test-da-key2"
	TestDataAttrKey3  = "test-da-key3"
	LargeDataContent1 = "this_is_a_large_data_content_that_should_be_stored_in_s3_for_testing_purposes_with_more_than_10_characters"
	LargeDataContent2 = "another_large_data_content_for_second_attribute_that_exceeds_the_s3_threshold_for_external_storage_testing"
	SmallDataContent3 = "small"
)

var TestDataAttributeVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("\"" + LargeDataContent1 + "\""),
}

var TestDataAttributeVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("\"" + LargeDataContent2 + "\""),
}

var TestDataAttributeVal3 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("\"" + SmallDataContent3 + "\""),
}

type handler struct {
	invokeHistory sync.Map
	invokeData    sync.Map
}

func NewHandler() *handler {
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

	if req.GetWorkflowType() == WorkflowType {
		if req.GetWorkflowStateId() == State1 {
			// Increment invoke count
			if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
			} else {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
			}

			// Store the state input for verification
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start_input", req.GetStateInput())

			// Validate that data attributes received match exactly the initial values provided at workflow start
			queryAtts := req.GetDataObjects()
			log.Printf("S1 WaitUntil: Received %d data attributes, validating they match initial values", len(queryAtts))

			foundAttr1 := false
			foundAttr2 := false
			foundAttr3 := false
			validationErrors := []string{}

			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataAttrKey1 {
					expectedData := *TestDataAttributeVal1.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr1 = true
						h.invokeData.Store("S1_start_attr1_data", receivedData)
						log.Printf("S1 WaitUntil: ✅ %s value matches initial data (length: %d)", TestDataAttrKey1, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr1 mismatch")
						log.Printf("S1 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey1, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataAttrKey2 {
					expectedData := *TestDataAttributeVal2.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr2 = true
						h.invokeData.Store("S1_start_attr2_data", receivedData)
						log.Printf("S1 WaitUntil: ✅ %s value matches initial data (length: %d)", TestDataAttrKey2, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr2 mismatch")
						log.Printf("S1 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey2, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataAttrKey3 {
					expectedData := *TestDataAttributeVal3.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr3 = true
						h.invokeData.Store("S1_start_attr3_data", receivedData)
						log.Printf("S1 WaitUntil: ✅ %s value matches initial data (length: %d)", TestDataAttrKey3, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr3 mismatch")
						log.Printf("S1 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey3, expectedData, receivedData)
					}
				}
			}

			allValidationsPass := foundAttr1 && foundAttr2 && foundAttr3 && len(validationErrors) == 0
			log.Printf("S1 WaitUntil: Data attribute validation complete - all match initial values: %t", allValidationsPass)

			h.invokeData.Store("S1_start_attr1_found", foundAttr1)
			h.invokeData.Store("S1_start_attr2_found", foundAttr2)
			h.invokeData.Store("S1_start_attr3_found", foundAttr3)
			h.invokeData.Store("S1_start_total_attrs", len(queryAtts))
			h.invokeData.Store("S1_start_validation_pass", allValidationsPass)

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State2 {
			// Increment invoke count
			if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
			} else {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
			}

			// State2 waitUntil doesn't need to check data attributes
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
		if req.GetWorkflowStateId() == State1 {
			// Increment invoke count
			if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_decide"); ok {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", value.(int64)+1)
			} else {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", int64(1))
			}

			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide_input", req.GetStateInput())

			// Transition to State2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State2,
							StateInput: req.StateInput,
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State2 {
			// Increment invoke count
			if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_decide"); ok {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", value.(int64)+1)
			} else {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", int64(1))
			}

			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide_input", req.GetStateInput())

			// Validate that data attributes received match exactly the initial values provided at workflow start
			queryAtts := req.GetDataObjects()
			log.Printf("S2 Execute: Received %d data attributes, validating they match initial values", len(queryAtts))

			foundAttr1 := false
			foundAttr2 := false
			foundAttr3 := false
			validationErrors := []string{}

			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataAttrKey1 {
					expectedData := *TestDataAttributeVal1.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr1 = true
						h.invokeData.Store("S2_decide_attr1_data", receivedData)
						log.Printf("S2 Execute: ✅ %s value matches initial data (length: %d)", TestDataAttrKey1, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr1 mismatch")
						log.Printf("S2 Execute: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey1, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataAttrKey2 {
					expectedData := *TestDataAttributeVal2.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr2 = true
						h.invokeData.Store("S2_decide_attr2_data", receivedData)
						log.Printf("S2 Execute: ✅ %s value matches initial data (length: %d)", TestDataAttrKey2, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr2 mismatch")
						log.Printf("S2 Execute: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey2, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataAttrKey3 {
					expectedData := *TestDataAttributeVal3.Data
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundAttr3 = true
						h.invokeData.Store("S2_decide_attr3_data", receivedData)
						log.Printf("S2 Execute: ✅ %s value matches initial data (length: %d)", TestDataAttrKey3, len(receivedData))
					} else {
						validationErrors = append(validationErrors, "attr3 mismatch")
						log.Printf("S2 Execute: ❌ %s value mismatch - expected: %s, received: %s", TestDataAttrKey3, expectedData, receivedData)
					}
				}
			}

			allValidationsPass := foundAttr1 && foundAttr2 && foundAttr3 && len(validationErrors) == 0
			log.Printf("S2 Execute: Data attribute validation complete - all match initial values: %t", allValidationsPass)

			h.invokeData.Store("S2_decide_attr1_found", foundAttr1)
			h.invokeData.Store("S2_decide_attr2_found", foundAttr2)
			h.invokeData.Store("S2_decide_attr3_found", foundAttr3)
			h.invokeData.Store("S2_decide_total_attrs", len(queryAtts))
			h.invokeData.Store("S2_decide_validation_pass", allValidationsPass)

			// Complete workflow
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    service.GracefulCompletingWorkflowStateId,
							StateInput: req.StateInput,
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
	outInvokehistory := make(map[string]interface{})
	h.invokeHistory.Range(func(key, value interface{}) bool {
		outInvokehistory[key.(string)] = value
		return true
	})

	outInvokeData := make(map[string]interface{})
	h.invokeData.Range(func(key, value interface{}) bool {
		outInvokeData[key.(string)] = value
		return true
	})

	// Merge both maps
	for k, v := range outInvokeData {
		outInvokehistory[k] = v
	}

	return nil, outInvokehistory
}
