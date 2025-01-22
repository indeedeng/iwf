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
 * This test workflow has 7 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- Waits one second before executing
 *      - Execute method will loop back to State1 five times; then execute method will go to State2
 * State2:
 *		- First execution: loops back to State2 + goes to State3
 *      - Second execution (after 1 second): goes to State3 and State4
 * State3:
 *		- Waits 8 seconds
 *      - Execute method will gracefully complete workflow
 * State4:
 *		- Waits on command trigger
 *      - Execute method will go to State5
 * State5:
 *		- Skips waitUntil and executes momentarily
 *      - Execute method will go to State6 and State7
 * State6:
 *		- Waits 4 seconds
 *      - Execute method will gracefully complete workflow
 * State7:
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
	State6       = "S6"
	State7       = "S7"

	SignalName = "test-signal"
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
		if req.GetWorkflowStateId() == State1 || req.GetWorkflowStateId() == State2 || req.GetWorkflowStateId() == State3 ||
			req.GetWorkflowStateId() == State5 || req.GetWorkflowStateId() == State6 || req.GetWorkflowStateId() == State7 {
			// Go straight to the decide methods without any commands
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State4 {
			// Proceed after signal is received
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
				// Move to State 2
				c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
					StateDecision: &iwfidl.StateDecision{
						NextStates: []iwfidl.StateMovement{
							{
								StateId: State2,
								StateOptions: &iwfidl.WorkflowStateOptions{
									SkipWaitUntil: iwfidl.PtrBool(true),
								},
							},
						},
					},
				})
			} else {
				// Repeat State 1 (5 times)
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
			context := req.GetContext()
			if context.GetStateExecutionId() == "S2-2" {
				// Move to State 3 & 4
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
			} else {
				// Repeat State 2 and Move to State 3
				c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
					StateDecision: &iwfidl.StateDecision{
						NextStates: []iwfidl.StateMovement{
							{
								StateId: State2,
								StateOptions: &iwfidl.WorkflowStateOptions{
									SkipWaitUntil: iwfidl.PtrBool(true),
								},
							},
							{
								StateId: State3,
							},
						},
					},
				})
			}
			return
		} else if req.GetWorkflowStateId() == State3 {
			time.Sleep(time.Second * 8)
			// Move to Completion
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
		} else if req.GetWorkflowStateId() == State4 {
			// Move to State 5, skipping wait until
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
			// Move to State 6 and State 7 skipping wait until for 7
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State6,
						},
						{
							StateId: State7,
							StateOptions: &iwfidl.WorkflowStateOptions{
								SkipWaitUntil: iwfidl.PtrBool(true),
							},
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State6 {
			time.Sleep(time.Second * 4)
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
		} else if req.GetWorkflowStateId() == State7 {
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
	return h.invokeHistory, nil
}
