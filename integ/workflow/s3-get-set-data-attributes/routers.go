package s3GetSetDataAttributes

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
 * Test workflow for S3 external storage with get/set data attributes APIs.
 * Tests both small data (stays in Temporal) and large data (goes to S3).
 *
 * State1:
 *   - Simple workflow that waits and completes
 *
 * The main testing is done via direct API calls to get/set data attributes,
 * not through workflow state transitions.
 */

const (
	WorkflowType = "s3-get-set-data-attributes"
	State1       = "S1"

	SmallDataKey        = "small-data"
	LargeDataKey        = "large-data"
	AnotherLargeDataKey = "another-large-data"

	// Small data content (stays in Temporal - under 50 byte threshold)
	SmallDataContent = "small"

	// Large data content (goes to S3 - over 50 byte threshold)
	LargeDataContent        = "large-data-content-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx" // Over 50 bytes
	AnotherLargeDataContent = "another-large-data-content-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy"         // Over 50 bytes

	// Updated values for testing updates
	UpdatedSmallDataContent = "updated-small"
	UpdatedLargeDataContent = "updated-large-data-content-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz" // Over 50 bytes
)

var (
	SmallDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"" + SmallDataContent + "\""),
	}

	LargeDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"" + LargeDataContent + "\""),
	}

	AnotherLargeDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"" + AnotherLargeDataContent + "\""),
	}

	UpdatedSmallDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"" + UpdatedSmallDataContent + "\""),
	}

	UpdatedLargeDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("\"" + UpdatedLargeDataContent + "\""),
	}
)

type handler struct {
	invokeHistory sync.Map
}

func NewHandler() *handler {
	return &handler{}
}

// GetTestResult returns the test result
func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	outInvokehistory := make(map[string]interface{})
	h.invokeHistory.Range(func(key, value interface{}) bool {
		outInvokehistory[key.(string)] = value
		return true
	})
	return nil, outInvokehistory
}

// ApiV1WorkflowStartPost - Define workflow states
func (h *handler) ApiV1WorkflowStartPost(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GetIwfWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	c.JSON(http.StatusOK, iwfidl.WorkflowStartResponse{
		WorkflowRunId: iwfidl.PtrString("test-run-id"),
	})
}

// ApiV1WorkflowStateStart - Handle state start (waitUntil)
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if req.GetWorkflowStateId() == State1 {
		h.invokeHistory.Store("S1_start", int64(1))

		// Simple waitUntil - no commands, just proceed
		c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
			CommandRequest: &iwfidl.CommandRequest{
				DeciderTriggerType: iwfidl.ANY_COMMAND_COMPLETED.Ptr(),
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

// ApiV1WorkflowStateDecide - Handle state execution (execute)
func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if req.GetWorkflowStateId() == State1 {
		h.invokeHistory.Store("S1_decide", int64(1))

		// Complete the workflow
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

	c.JSON(http.StatusBadRequest, struct{}{})
}
