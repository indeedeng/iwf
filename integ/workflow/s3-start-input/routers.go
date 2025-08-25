package s3_start_input

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
 * This test workflow has 1 state, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "s3-start-input"
	State1       = "S1"
)

type handler struct {
	invokeHistory sync.Map
}

func NewHandler() *handler {
	return &handler{
		invokeHistory: sync.Map{},
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
		// go straight to the decide(execute) without any commands
		if req.GetWorkflowStateId() == State1 {
			if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
			} else {
				h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
			}

			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start_input", req.GetStateInput())

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

		h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide_input", req.GetStateInput())

		if req.GetWorkflowStateId() == State1 {
			// Move to completion
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
	return nil, outInvokehistory
}
