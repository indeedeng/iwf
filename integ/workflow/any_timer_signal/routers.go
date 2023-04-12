package anytimersignal

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"time"
)

const (
	WorkflowType = "any_timer_signal"
	State1       = "S1"
	State2       = "S2"
	SignalName   = "test-signal-name"
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
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			var timerCommands []iwfidl.TimerCommand
			context := req.GetContext()
			if context.GetStateExecutionId() == State1+"-"+"1" {
				now := time.Now().Unix()
				timerCommands = []iwfidl.TimerCommand{
					{
						FiringUnixTimestampSeconds: now + 1, // fire after 1s
					},
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         "signal-cmd-id",
							SignalChannelName: SignalName,
						},
					},
					TimerCommands:      timerCommands,
					CommandWaitingType: ptr.Any(iwfidl.ANY_COMPLETED),
				},
			})
			return
		}
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

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
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
			var movements []iwfidl.StateMovement

			context := req.GetContext()
			if context.GetStateExecutionId() == State1+"-"+"1" {
				h.invokeData["signalChannelName1"] = signalResults.SignalResults[0].GetSignalChannelName()
				h.invokeData["signalCommandId1"] = signalResults.SignalResults[0].GetCommandId()
				h.invokeData["signalStatus1"] = signalResults.SignalResults[0].GetSignalRequestStatus()
				movements = []iwfidl.StateMovement{{StateId: State1}}
			} else {
				h.invokeData["signalChannelName2"] = signalResults.SignalResults[0].GetSignalChannelName()
				h.invokeData["signalCommandId2"] = signalResults.SignalResults[0].GetCommandId()
				h.invokeData["signalStatus2"] = signalResults.SignalResults[0].GetSignalRequestStatus()
				h.invokeData["signalValue2"] = signalResults.SignalResults[0].GetSignalValue()
				movements = []iwfidl.StateMovement{{StateId: State2}}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: movements,
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			// go to complete
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
