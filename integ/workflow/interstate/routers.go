package interstate

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
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

		if req.GetWorkflowStateId() == State1 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State21 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   "cmd-1",
							ChannelName: channel1,
						},
					},
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State31 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   "cmd-2",
							ChannelName: channel2,
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State22 {
			time.Sleep(time.Second * 2)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
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

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				// dead end
				StateDecision: &iwfidl.StateDecision{},
			})
			return
		}

		if req.GetWorkflowStateId() == State22 {
			time.Sleep(time.Second * 2)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				// dead end
				StateDecision: &iwfidl.StateDecision{},
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
