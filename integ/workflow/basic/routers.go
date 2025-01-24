package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"testing"
)

/**
 * This test workflow has 2 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- Waits on nothing. Will execute momentarily
 *      - Execute method will move to State2
 * State2:
 *		- Waits on nothing. Will execute momentarily
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "basic"
	State1       = "S1"
	State2       = "S2"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() *handler {
	return &handler{
		invokeHistory: make(map[string]int64),
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

	context := req.GetContext()
	if context.GetAttempt() <= 0 || context.GetFirstAttemptTimestamp() <= 0 {
		t.Fatal("attempt and firstAttemptTimestamp should be greater than zero")
	}

	if req.GetWorkflowType() == WorkflowType {
		// Basic workflow go straight to decide methods without any commands
		if req.GetWorkflowStateId() == State1 || req.GetWorkflowStateId() == State2 {
			h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
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
	context := req.GetContext()
	if context.GetAttempt() <= 0 || context.GetFirstAttemptTimestamp() <= 0 {
		t.Fatal("attempt and firstAttemptTimestamp should be greater than zero")
	}

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

		if req.GetWorkflowStateId() == State1 {
			// Move to next state
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State2,
							StateInput: req.StateInput,
							StateOptions: &iwfidl.WorkflowStateOptions{
								StartApiTimeoutSeconds:   iwfidl.PtrInt32(14),
								ExecuteApiTimeoutSeconds: iwfidl.PtrInt32(15),
								StartApiRetryPolicy: &iwfidl.RetryPolicy{
									InitialIntervalSeconds: iwfidl.PtrInt32(14),
									BackoffCoefficient:     iwfidl.PtrFloat32(14),
									MaximumAttempts:        iwfidl.PtrInt32(14),
									MaximumIntervalSeconds: iwfidl.PtrInt32(14),
								},
								ExecuteApiRetryPolicy: &iwfidl.RetryPolicy{
									InitialIntervalSeconds: iwfidl.PtrInt32(15),
									BackoffCoefficient:     iwfidl.PtrFloat32(15),
									MaximumAttempts:        iwfidl.PtrInt32(15),
									MaximumIntervalSeconds: iwfidl.PtrInt32(15),
								},
							},
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
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
	return h.invokeHistory, nil
}
