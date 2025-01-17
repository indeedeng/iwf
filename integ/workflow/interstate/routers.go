package interstate

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"time"
)

const (
	WorkflowType = "interstate"
	State1       = "S1"
	State21      = "S21"
	State22      = "S22"
	State31      = "S31"

	channel1 = "channel1"
	channel2 = "channel2"
)

var TestVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-value1"),
}

var TestVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-value2"),
}

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func NewHandler() *handler {
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

		// Go straight to the decide methods without any commands
		if req.GetWorkflowStateId() == State1 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		// Will proceed once channel 1 has been published to
		if req.GetWorkflowStateId() == State21 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-1"),
							ChannelName: channel1,
						},
					},
				},
			})
			return
		}
		// Will proceed once channel 2 has been published to
		if req.GetWorkflowStateId() == State31 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-2"),
							ChannelName: channel2,
						},
					},
				},
			})
			return
		}

		// Wait 2 seconds then publish on the first channel
		if req.GetWorkflowStateId() == State22 {
			time.Sleep(time.Second * 2)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},

				PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
					{
						ChannelName: channel1,
						Value:       &TestVal1,
					},
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
			// First state requires no pre-reqs
			// Move to state 21 & 22:
			// 21 - Will wait for channel 1
			// 22 - Will wait 3 seconds then publish to channel 1
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State21,
						},
						{
							StateId: State22,
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State21 {
			results := req.GetCommandResults()
			h.invokeData[State21+"received"] = results.GetInterStateChannelResults()[0].GetValue()

			// Move to state 31, which will wait for channel 2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State31,
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State31 {
			results := req.GetCommandResults()
			h.invokeData[State31+"received"] = results.GetInterStateChannelResults()[0].GetValue()

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

		if req.GetWorkflowStateId() == State22 {
			time.Sleep(time.Second * 2)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				// Move to the dead-end state and publish on channel 2 (to unlock State 31)
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.DeadEndWorkflowStateId,
						},
					},
				},
				PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
					{
						ChannelName: channel2,
						Value:       &TestVal2,
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
