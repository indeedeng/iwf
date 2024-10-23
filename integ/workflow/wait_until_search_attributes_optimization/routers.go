package wait_until_search_attributes_optimization

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

/**
 * This test workflow has 5 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- Waits on nothing. Will execute momentarily
 *      - Executes method will loop back to State1 five times then execute method will go to State2
 * State2:
 *		- Waits 5 seconds
 *      - Execute method will go to State3 and State4
 * State3:
 *		- Does nothing
 * State4:
 *		- Waits on command trigger
 *      - Execute method will go to State5
 * State5:
 *		- Skips waitUntil and executes momentarily
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "wait_until_search_optimization"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
	State4       = "S4"
	State5       = "S5"

	SignalName = "test-signal"

	TestSearchAttributeExecutingStateIdsKey = "IwfExecutingStateIds"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
	}
}

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
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
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
		if req.GetWorkflowStateId() == State3 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State4 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         ptr.Any("test"),
							SignalChannelName: SignalName,
						},
					},
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State5 {
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

			context := req.GetContext()
			if context.GetStateExecutionId() == "S1-5" {
				c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
					StateDecision: &iwfidl.StateDecision{
						NextStates: []iwfidl.StateMovement{
							{
								StateId: State2,
							},
						},
					},
				})
			} else {
				time.Sleep(time.Second * 1)
				c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
					StateDecision: &iwfidl.StateDecision{
						NextStates: []iwfidl.StateMovement{
							{
								StateId: State1,
							},
						},
					},
				})
			}
			return
		} else if req.GetWorkflowStateId() == State2 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State3,
						},
						{
							StateId: State4,
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State3 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.DeadEndWorkflowStateId,
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State4 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State5,
							StateOptions: &iwfidl.WorkflowStateOptions{
								SkipWaitUntil: iwfidl.PtrBool(true),
							},
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State5 {
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
	return h.invokeHistory, nil
}
