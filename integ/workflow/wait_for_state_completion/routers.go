package wait_for_state_completion

import (
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
)

/**
 * This test workflow has 2 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- 10-second delay is added before executing state
 *      - Execute method will go to State2
 * State2:
 *		- Waits on nothing. Will execute momentarily
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "wait_for_state_completion"
	State1       = "S1"
	State2       = "S2"
)

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
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			nowInt, err := strconv.Atoi(req.StateInput.GetData())
			if err != nil {
				t.Fatal(err)
			}
			now := int64(nowInt)
			h.invokeData["scheduled_at"] = now
			// Proceed after 10s
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					TimerCommands: []iwfidl.TimerCommand{
						{
							CommandId:                  ptr.Any("timer-cmd-id"),
							FiringUnixTimestampSeconds: iwfidl.PtrInt64(now + 10), // fire after 10s
						},
					},
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}

		// Go straight to the decide methods without any commands
		if req.GetWorkflowStateId() == State2 {
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
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		if req.GetWorkflowStateId() == State1 {
			now := time.Now().Unix()
			h.invokeData["fired_at"] = now
			timerResults := req.GetCommandResults()
			timerId := timerResults.GetTimerResults()[0].GetCommandId()
			h.invokeData["timer_id"] = timerId
			// Move to State 2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State2,
							WaitForKey: ptr.Any("testKey"),
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
