package s3_state_input_optimization

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
 * This test workflow has 3 states, testing S3 state input optimization functionality.
 * All states use the same large input data to test deduplication.
 *
 * State1:
 *		- WaitUntil method stores input data for verification
 *      - Execute method transitions to State2 with same input
 *
 * State2:
 *		- WaitUntil method stores input data for verification
 *      - Execute method transitions to State3 with same input
 *
 * State3:
 *		- WaitUntil method stores input data for verification
 *      - Execute method completes workflow
 */
const (
	WorkflowType = "s3-state-input-optimization"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
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

		// Store input data for verification
		if req.StateInput != nil && req.StateInput.Data != nil {
			h.invokeData.Store(stateId+"_input_data", *req.StateInput.Data)
			log.Printf("%s WaitUntil: Received input data (length: %d): %s", stateId, len(*req.StateInput.Data), *req.StateInput.Data)
		}

		if stateId == State1 || stateId == State2 || stateId == State3 {
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
			// Transition to State2 with same input (should reuse S3 object)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State2,
							StateInput: req.StateInput, // Same input - should trigger optimization
						},
					},
				},
			})
			return
		}

		if stateId == State2 {
			// Transition to State3 with same input (should reuse S3 object again)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State3,
							StateInput: req.StateInput, // Same input - should trigger optimization
						},
					},
				},
			})
			return
		}

		if stateId == State3 {
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
