package s3_upsert_data_objects

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
 * This test workflow has 2 states, testing S3 upsert data objects functionality.
 *
 * State1:
 *		- WaitUntil method does nothing
 *      - Execute method upserts large data objects that should go to S3, then transitions to State2
 *
 * State2:
 *		- WaitUntil method validates it receives the upserted data objects from S3
 *      - Execute method completes workflow
 */
const (
	WorkflowType      = "s3-upsert-data-objects"
	State1            = "S1"
	State2            = "S2"
	TestDataObjKey1   = "large_obj1"
	TestDataObjKey2   = "large_obj2"
	TestDataObjKey3   = "small_obj3"
	LargeDataContent1 = "this_is_a_large_data_content_that_should_be_stored_in_s3_for_upsert_testing_purposes_with_more_than_10_characters"
	LargeDataContent2 = "another_large_data_content_for_second_upserted_object_that_exceeds_the_s3_threshold_for_external_storage_testing"
	SmallDataContent3 = "small"
)

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
		stateId := req.GetWorkflowStateId()

		// Increment invoke count
		if value, ok := h.invokeHistory.Load(stateId + "_start"); ok {
			h.invokeHistory.Store(stateId+"_start", value.(int64)+1)
		} else {
			h.invokeHistory.Store(stateId+"_start", int64(1))
		}

		if stateId == State1 {
			// State1 waitUntil doesn't need to do anything
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}

		if stateId == State2 {
			// State2 waitUntil should validate data objects received from State1's upsert
			queryAtts := req.GetDataObjects()
			log.Printf("S2 WaitUntil: Received %d data objects, validating they match upserted values", len(queryAtts))

			foundLargeObj1 := false
			foundLargeObj2 := false
			foundSmallObj3 := false

			for _, queryAtt := range queryAtts {
				if queryAtt.GetKey() == TestDataObjKey1 {
					expectedData := "\"" + LargeDataContent1 + "\""
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundLargeObj1 = true
						h.invokeData.Store("S2_large_obj1_data", LargeDataContent1)
						log.Printf("S2 WaitUntil: ✅ %s value matches upserted data (length: %d)", TestDataObjKey1, len(receivedData))
					} else {
						log.Printf("S2 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataObjKey1, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataObjKey2 {
					expectedData := "\"" + LargeDataContent2 + "\""
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundLargeObj2 = true
						h.invokeData.Store("S2_large_obj2_data", LargeDataContent2)
						log.Printf("S2 WaitUntil: ✅ %s value matches upserted data (length: %d)", TestDataObjKey2, len(receivedData))
					} else {
						log.Printf("S2 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataObjKey2, expectedData, receivedData)
					}
				}
				if queryAtt.GetKey() == TestDataObjKey3 {
					expectedData := "\"" + SmallDataContent3 + "\""
					receivedData := *queryAtt.GetValue().Data
					if receivedData == expectedData {
						foundSmallObj3 = true
						h.invokeData.Store("S2_small_obj3_data", SmallDataContent3)
						log.Printf("S2 WaitUntil: ✅ %s value matches upserted data (length: %d)", TestDataObjKey3, len(receivedData))
					} else {
						log.Printf("S2 WaitUntil: ❌ %s value mismatch - expected: %s, received: %s", TestDataObjKey3, expectedData, receivedData)
					}
				}
			}

			h.invokeData.Store("S2_received_large_obj1", foundLargeObj1)
			h.invokeData.Store("S2_received_large_obj2", foundLargeObj2)
			h.invokeData.Store("S2_received_small_obj3", foundSmallObj3)

			log.Printf("S2 WaitUntil: Data object validation complete - found large_obj1: %t, large_obj2: %t, small_obj3: %t",
				foundLargeObj1, foundLargeObj2, foundSmallObj3)

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
		stateId := req.GetWorkflowStateId()

		// Increment invoke count
		if value, ok := h.invokeHistory.Load(stateId + "_decide"); ok {
			h.invokeHistory.Store(stateId+"_decide", value.(int64)+1)
		} else {
			h.invokeHistory.Store(stateId+"_decide", int64(1))
		}

		if stateId == State1 {
			// State1 Execute: Upsert large data objects that should go to S3
			log.Printf("S1 Execute: Upserting data objects - 2 large (should go to S3), 1 small (should stay in memory)")

			upsertDataObjects := []iwfidl.KeyValue{
				{
					Key: iwfidl.PtrString(TestDataObjKey1),
					Value: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("\"" + LargeDataContent1 + "\""), // Large - should go to S3
					},
				},
				{
					Key: iwfidl.PtrString(TestDataObjKey2),
					Value: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("\"" + LargeDataContent2 + "\""), // Large - should go to S3
					},
				},
				{
					Key: iwfidl.PtrString(TestDataObjKey3),
					Value: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("\"" + SmallDataContent3 + "\""), // Small - should stay in memory
					},
				},
			}

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
				UpsertDataObjects: upsertDataObjects,
			})
			return
		}

		if stateId == State2 {
			// State2 Execute: Complete workflow
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
