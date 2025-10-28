package command_thread_completion

import (
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
)

/**
 * This test workflow validates that all command threads complete before continue-as-new snapshots state.
 * It tests the fix for the bug where internal channel signals were lost during continue-as-new.
 *
 * Workflow structure:
 * State1:
 *   - WaitUntil: Set up timer, signal, and internal channel commands with ANY_COMMAND_COMPLETED
 *   - Execute: Publish to internal channel, move to State2
 * State2:
 *   - WaitUntil: Wait for the internal channel from State1
 *   - Execute: Complete workflow
 *
 * The test triggers continue-as-new after State1 execute but before State2 starts.
 * This ensures the internal channel signal published by State1 is captured before continue-as-new.
 */
const (
	WorkflowType = "command_thread_completion"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
	StateAnyCmd  = "StateAnyCmd" // Tests ANY_COMMAND_COMPLETED with CAN

	testChannel    = "test-channel"
	testSignal     = "test-signal"
	testTimerCmd   = "test-timer"
	testChannelCmd = "test-channel-cmd"
	testSignalCmd  = "test-signal-cmd"
)

var testChannelValue = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("channel-data"),
}

type handler struct {
	invokeHistory sync.Map
	invokeData    sync.Map
}

func NewHandler() *handler {
	return &handler{
		invokeHistory: sync.Map{},
		invokeData:    sync.Map{},
	}
}

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.recordInvoke(req.GetWorkflowStateId() + "_start")

		if req.GetWorkflowStateId() == State1 {
			// State1: Set up all three command types with ALL_COMMAND_COMPLETED
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					TimerCommands: []iwfidl.TimerCommand{
						{
							CommandId:                  ptr.Any(testTimerCmd),
							FiringUnixTimestampSeconds: ptr.Any(time.Now().Add(2 * time.Second).Unix()),
						},
					},
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         ptr.Any(testSignalCmd),
							SignalChannelName: testSignal,
						},
					},
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any(testChannelCmd),
							ChannelName: testChannel,
						},
					},
				},
				// Immediately publish to internal channel so it's available for the thread to retrieve
				PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
					{
						ChannelName: testChannel,
						Value:       &testChannelValue,
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State2 {
			// State2: Wait for the channel published by State1's Execute
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("s2-channel-cmd"),
							ChannelName: testChannel + "2",
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State3 {
			// State3: Only wait for a timer command (tests timer thread in isolation)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					TimerCommands: []iwfidl.TimerCommand{
						{
							CommandId:                  ptr.Any("s3-timer-cmd"),
							FiringUnixTimestampSeconds: ptr.Any(time.Now().Add(2 * time.Second).Unix()),
						},
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == StateAnyCmd {
			// StateAnyCmd: Tests ANY_COMMAND_COMPLETED with long timer + quick signal
			// This validates that we don't wait for the timer when signal completes
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ANY_COMMAND_COMPLETED.Ptr(), // ANY, not ALL!
					TimerCommands: []iwfidl.TimerCommand{
						{
							CommandId:                  ptr.Any("any-cmd-timer"),
							FiringUnixTimestampSeconds: ptr.Any(time.Now().Add(20 * time.Second).Unix()), // Long timer
						},
					},
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         ptr.Any("any-cmd-signal-cmd"),
							SignalChannelName: "any-cmd-signal",
						},
					},
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
		h.recordInvoke(req.GetWorkflowStateId() + "_decide")

		if req.GetWorkflowStateId() == State1 {
			// With ALL_COMMAND_COMPLETED, all three command types must complete
			cmdResults := req.GetCommandResults()

			// Check timer - should be FIRED
			timerFired := false
			if cmdResults.HasTimerResults() {
				for _, tr := range cmdResults.GetTimerResults() {
					if tr.GetTimerStatus() == iwfidl.FIRED {
						timerFired = true
						h.recordData("s1_timer_fired", true)
					}
				}
			}
			if !timerFired {
				log.Println("ERROR: Timer should have fired in State1")
			}

			// Check signal - should be RECEIVED
			signalReceived := false
			if cmdResults.HasSignalResults() {
				for _, sr := range cmdResults.GetSignalResults() {
					if sr.GetSignalChannelName() == testSignal && sr.GetSignalRequestStatus() == iwfidl.RECEIVED {
						signalReceived = true
						h.recordData("s1_signal_received", true)
					}
				}
			}
			if !signalReceived {
				log.Println("ERROR: Signal should have been received in State1")
			}

			// Check internal channel - should be RECEIVED
			channelReceived := false
			if cmdResults.HasInterStateChannelResults() {
				for _, cr := range cmdResults.GetInterStateChannelResults() {
					if cr.GetChannelName() == testChannel && cr.GetRequestStatus() == iwfidl.RECEIVED {
						channelReceived = true
						h.recordData("s1_channel_received", true)
					}
				}
			}
			if !channelReceived {
				log.Println("ERROR: Internal channel should have been received in State1")
			}

			// Move to both State2 and State3 - publish channel for State2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State2,
						},
						{
							StateId: State3,
						},
					},
				},
				PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
					{
						ChannelName: testChannel + "2",
						Value:       &testChannelValue,
					},
				},
			})
			return
		}

		if req.GetWorkflowStateId() == State2 {
			// Verify the channel was received (this tests continue-as-new preservation)
			cmdResults := req.GetCommandResults()

			channelReceived := false
			if cmdResults.HasInterStateChannelResults() {
				for _, cr := range cmdResults.GetInterStateChannelResults() {
					if cr.GetChannelName() == testChannel+"2" && cr.GetRequestStatus() == iwfidl.RECEIVED {
						channelReceived = true
						h.recordData("s2_channel_received", true)
						h.recordData("s2_channel_value", cr.GetValue())
					}
				}
			}

			if !channelReceived {
				log.Println("ERROR: State2 channel was NOT received! This indicates the bug exists.")
				h.recordData("s2_channel_received", false)
			}

			// Dead end - don't complete workflow yet, let State3 complete
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
		}

		if req.GetWorkflowStateId() == State3 {
			// Verify the timer fired (tests timer thread in isolation)
			cmdResults := req.GetCommandResults()

			timerFired := false
			if cmdResults.HasTimerResults() {
				for _, tr := range cmdResults.GetTimerResults() {
					if tr.GetTimerStatus() == iwfidl.FIRED {
						timerFired = true
						h.recordData("s3_timer_fired", true)
					}
				}
			}

			if !timerFired {
				log.Println("ERROR: Timer should have fired in State3")
				h.recordData("s3_timer_fired", false)
			}

			// Complete workflow
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

		if req.GetWorkflowStateId() == StateAnyCmd {
			// Verify that with ANY_COMMAND_COMPLETED, we proceeded when signal was received
			// without waiting for the long timer
			cmdResults := req.GetCommandResults()

			signalReceived := false
			timerFired := false

			if cmdResults.HasSignalResults() {
				for _, sr := range cmdResults.GetSignalResults() {
					if sr.GetSignalChannelName() == "any-cmd-signal" && sr.GetSignalRequestStatus() == iwfidl.RECEIVED {
						signalReceived = true
						h.recordData("any_cmd_signal_received", true)
					}
				}
			}

			if cmdResults.HasTimerResults() {
				for _, tr := range cmdResults.GetTimerResults() {
					if tr.GetCommandId() == "any-cmd-timer" && tr.GetTimerStatus() == iwfidl.FIRED {
						timerFired = true
					}
				}
			}

			if !signalReceived {
				log.Println("ERROR: Signal should have been received in StateAnyCmd (ANY_COMMAND_COMPLETED)")
				h.recordData("any_cmd_signal_received", false)
			}

			if timerFired {
				log.Println("WARNING: Timer fired in StateAnyCmd - this suggests we waited for it instead of proceeding with signal")
			}

			// Complete workflow
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State3,
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
	history := make(map[string]int64)
	h.invokeHistory.Range(func(key, value interface{}) bool {
		history[key.(string)] = value.(int64)
		return true
	})

	data := make(map[string]interface{})
	h.invokeData.Range(func(key, value interface{}) bool {
		data[key.(string)] = value
		return true
	})

	return history, data
}

func (h *handler) recordInvoke(key string) {
	if value, ok := h.invokeHistory.Load(key); ok {
		h.invokeHistory.Store(key, value.(int64)+1)
	} else {
		h.invokeHistory.Store(key, int64(1))
	}
}

func (h *handler) recordData(key string, value interface{}) {
	h.invokeData.Store(key, value)
}
