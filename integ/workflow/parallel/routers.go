package parallel

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

/**
 * This test workflow has eight states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method does nothing
 * 		- Execute method delays 1s then moves to State11, State12, & State13
 * State11:
 *		- WaitUntil method does nothing
 * 		- Execute method delays 2s then moves to State111 & State112
 * State12:
 *		- WaitUntil method does nothing
 * 		- Execute method delays 2s then moves to State121 & State122
 * State13:
 *		- WaitUntil method does nothing
 *      - Execute method will delay 1s then gracefully complete workflow
 * State111:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 * State112:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 * State121:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 * State122:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "parallel"
	State1       = "S1"
	State11      = "S11"
	State12      = "S12"
	State13      = "S13"
	State111     = "S111"
	State112     = "S112"
	State121     = "S121"
	State122     = "S122"
)

type handler struct {
	invokeHistory sync.Map
}

func NewHandler() common.WorkflowHandler {
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
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
		}

		// Go straight to the decide methods without any commands
		c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
			CommandRequest: &iwfidl.CommandRequest{
				DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
			},
		})
		return
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

		var nextStates []iwfidl.StateMovement
		switch req.GetWorkflowStateId() {
		case State1:
			// Cause graceful complete to wait
			time.Sleep(time.Second * 1)

			// Move to 3 states (which will all move to this decide method without commands)
			nextStates = []iwfidl.StateMovement{
				{
					StateId: State11,
				},
				{
					StateId: State12,
				},
				{
					StateId: State13,
				},
			}
		case State11:
			// Cause graceful complete to wait
			time.Sleep(time.Second * 2)

			// Move to 2 states (which will all move to this decide method without commands)
			nextStates = []iwfidl.StateMovement{
				{
					StateId: State111,
				},
				{
					StateId: State112,
				},
			}
		case State12:
			// Cause graceful complete to wait
			time.Sleep(time.Second * 2)

			// Move to 2 states (which will all move to this decide method without commands)
			nextStates = []iwfidl.StateMovement{
				{
					StateId: State121,
				},
				{
					StateId: State122,
				},
			}
		case State13:
			// Cause graceful complete to wait
			time.Sleep(time.Second * 1)

			// Move to completion after updating the state input
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.GracefulCompletingWorkflowStateId,
					StateInput: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("from " + req.GetWorkflowStateId()),
					},
				},
			}
		case State111, State112, State121, State122:
			// Move to completion after updating the state input
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.GracefulCompletingWorkflowStateId,
					StateInput: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("from " + req.GetWorkflowStateId()),
					},
				},
			}
		default:
			// Fail workflow due to unknown or unexpected state
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.ForceFailingWorkflowStateId,
				},
			}
		}

		c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
			StateDecision: &iwfidl.StateDecision{
				NextStates: nextStates,
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	invokeHistory := make(map[string]int64)
	h.invokeHistory.Range(func(key, value interface{}) bool {
		invokeHistory[key.(string)] = value.(int64)
		return true
	})
	return invokeHistory, nil
}
