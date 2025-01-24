package anycommandclose

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"testing"
)

/**
 * This test workflow has 2 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil wait until a signal is received
 *      - Execute method will fire the signal and move the State2
 * State2:
 *		- Waits on nothing. Will execute momentarily
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "any_command_close"
	State1       = "S1"
	State2       = "S2"
	SignalName1  = "test-signal-name1"
	SignalName2  = "test-signal-name2"
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
			// Proceed after either signal is received
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         ptr.Any("signal-cmd-id1"),
							SignalChannelName: SignalName1,
						},
						{
							CommandId:         ptr.Any("signal-cmd-id2"),
							SignalChannelName: SignalName2,
						},
					},
					CommandWaitingType: ptr.Any(iwfidl.ANY_COMPLETED),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State2 {
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
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

		if req.GetWorkflowStateId() == State1 {
			signalResults := req.GetCommandResults()
			h.invokeData["signalCommandResultsLength"] = len(signalResults.SignalResults)

			// Trigger signals
			h.invokeData["signalChannelName0"] = signalResults.SignalResults[0].GetSignalChannelName()
			h.invokeData["signalCommandId0"] = signalResults.SignalResults[0].GetCommandId()
			h.invokeData["signalStatus0"] = signalResults.SignalResults[0].GetSignalRequestStatus()

			h.invokeData["signalChannelName1"] = signalResults.SignalResults[1].GetSignalChannelName()
			h.invokeData["signalCommandId1"] = signalResults.SignalResults[1].GetCommandId()
			h.invokeData["signalStatus1"] = signalResults.SignalResults[1].GetSignalRequestStatus()
			h.invokeData["signalValue1"] = signalResults.SignalResults[1].GetSignalValue()

			// Move to State 2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State2,
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
