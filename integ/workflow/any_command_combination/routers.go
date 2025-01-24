package anycommandcombination

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"testing"
	"time"
)

/**
 * This test workflow has 2 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil will fail its first attempt and then retry which will proceed when a combination is completed
 *      - Execute method will invoke the combination and move the State2
 * State2:
 *		- WaitUntil will fail its first attempt and then retry which will proceed when a combination is completed
 *      - Execute method will invoke the combination and gracefully complete workflow
 */
const (
	WorkflowType     = "any_command_combination"
	State1           = "S1"
	State2           = "S2"
	TimerId1         = "test-timer-1"
	SignalNameAndId1 = "test-signal-name1"
	SignalNameAndId2 = "test-signal-name2"
	SignalNameAndId3 = "test-signal-name3"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
	//we want to confirm that the interpreter workflow activity will fail when the commandId is empty with ANY_COMMAND_COMBINATION_COMPLETED
	hasS1RetriedForInvalidCommandId bool
	hasS2RetriedForInvalidCommandId bool
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory:                   make(map[string]int64),
		invokeData:                      make(map[string]interface{}),
		hasS1RetriedForInvalidCommandId: false,
		hasS2RetriedForInvalidCommandId: false,
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

	invalidTimerCommands := []iwfidl.TimerCommand{
		{
			FiringUnixTimestampSeconds: iwfidl.PtrInt64(time.Now().Unix() + 86400*365), // one year later
		},
	}
	validTimerCommands := []iwfidl.TimerCommand{
		{
			CommandId:                  ptr.Any(TimerId1),
			FiringUnixTimestampSeconds: iwfidl.PtrInt64(time.Now().Unix() + 86400*365), // one year later
		},
	}
	invalidSignalCommands := []iwfidl.SignalCommand{
		{
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         ptr.Any(SignalNameAndId2),
			SignalChannelName: SignalNameAndId2,
		},
	}
	validSignalCommands := []iwfidl.SignalCommand{
		{
			CommandId:         ptr.Any(SignalNameAndId1),
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         ptr.Any(SignalNameAndId1),
			SignalChannelName: SignalNameAndId1,
		},
		{
			CommandId:         ptr.Any(SignalNameAndId2),
			SignalChannelName: SignalNameAndId2,
		},
		{
			CommandId:         ptr.Any(SignalNameAndId3),
			SignalChannelName: SignalNameAndId3,
		},
	}

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++

		if req.GetWorkflowStateId() == State1 {
			// If the state has already retried an invalid command, proceed on combination completed
			if h.hasS1RetriedForInvalidCommandId {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
						CommandCombinations: []iwfidl.CommandCombination{
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId1, // wait for two SignalNameAndId1
								},
							},
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId2,
								},
							},
						},
					},
				}

				c.JSON(http.StatusOK, startResp)
			} else {
				// If the state has not already retried an invalid command, return invalid trigger signals, which will fail
				// and cause a retry
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      invalidTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
					},
				}
				h.hasS1RetriedForInvalidCommandId = true
				c.JSON(http.StatusOK, startResp)
			}
			return
		}

		if req.GetWorkflowStateId() == State2 {
			// If the state has already retried an invalid command, return signals and completion metrics
			if h.hasS2RetriedForInvalidCommandId {
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     validSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
						CommandCombinations: []iwfidl.CommandCombination{
							{
								CommandIds: []string{
									SignalNameAndId2, SignalNameAndId1,
								},
							},
							{
								CommandIds: []string{
									TimerId1, SignalNameAndId1, SignalNameAndId2,
								},
							},
						},
					},
				}

				c.JSON(http.StatusOK, startResp)
			} else {
				// If the state has not already retried an invalid command, return invalid trigger signals, which will fail
				// and cause a retry
				startResp := iwfidl.WorkflowStateStartResponse{
					CommandRequest: &iwfidl.CommandRequest{
						SignalCommands:     invalidSignalCommands,
						TimerCommands:      validTimerCommands,
						DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
					},
				}
				h.hasS2RetriedForInvalidCommandId = true
				c.JSON(http.StatusOK, startResp)
			}
			return
		}
	}

	t.Fatal("invalid workflow type")
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

		// Trigger signals and move to State 2
		if req.GetWorkflowStateId() == State1 {
			h.invokeData["s1_commandResults"] = req.GetCommandResults()

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
			// Trigger data and move to completion
			h.invokeData["s2_commandResults"] = req.GetCommandResults()
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

	t.Fatal("invalid workflow type")
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
