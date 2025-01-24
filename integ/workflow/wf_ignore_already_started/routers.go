package wf_ignore_already_started

import (
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

/**
 * This test workflow has 1 state, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- Waits on a timer
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "wf_ignore_already_started"
	State1       = "S1"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
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

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			nowInt, err := strconv.Atoi(req.StateInput.GetData())
			if err != nil {
				t.Fatal(err)
			}
			now := int64(nowInt)
			h.invokeData["scheduled_at"] = now
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					TimerCommands: []iwfidl.TimerCommand{
						{
							CommandId:       ptr.Any("timer-cmd-id"),
							DurationSeconds: iwfidl.PtrInt64(10), // fire after 10s
						},
					},
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
