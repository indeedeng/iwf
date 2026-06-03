package interstate_consume_n

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
 * This test workflow verifies consuming N messages from an inter-state channel using AtLeast/AtMost.
 *
 * S1 (Setup):
 *   - Start: no commands
 *   - Decide: publishes 5 messages to "ch", moves to S2
 *
 * S2 (ExactN — AtLeast=3, AtMost=3):
 *   - Start: channel command on "ch" with AtLeast=3, AtMost=3
 *   - Decide: verifies 3 values consumed, moves to S3
 *
 * S3 (OneToAll — AtLeast=1, no AtMost):
 *   - Start: channel command on "ch" with AtLeast=1
 *   - Decide: verifies remaining 2 values consumed, moves to S4
 *
 * S4 (ZeroToAll — AtLeast=0, no AtMost on empty channel):
 *   - Start: channel command on "ch" with AtLeast=0
 *   - Decide: verifies 0 values consumed, moves to S5 + S6
 *
 * S5 (AtMostOnly — AtMost=2, no AtLeast; also tests late message arrival):
 *   - Start: channel command on "ch2" with AtMost=2 only
 *   - Decide: verifies 2 values consumed (out of 3 published by S6), completes workflow
 *
 * S6 (Delayed publisher):
 *   - Start: delays 2s, publishes 3 messages to "ch2", no commands
 *   - Decide: dead-end
 */
const (
	WorkflowType = "interstate_consume_n"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
	State4       = "S4"
	State5       = "S5"
	State6       = "S6"

	channel1 = "ch"
	channel2 = "ch2"
)

var TestValues = []iwfidl.EncodedObject{
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("val-0")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("val-1")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("val-2")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("val-3")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("val-4")},
}

var TestValuesCh2 = []iwfidl.EncodedObject{
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("ch2-val-0")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("ch2-val-1")},
	{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("ch2-val-2")},
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

		switch req.GetWorkflowStateId() {
		case State1:
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return

		case State2:
			// ExactN: wait for exactly 3, consume exactly 3
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-1"),
							ChannelName: channel1,
							AtLeast:     iwfidl.PtrInt32(3),
							AtMost:      iwfidl.PtrInt32(3),
						},
					},
				},
			})
			return

		case State3:
			// OneToAll: wait for at least 1, consume all available
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-2"),
							ChannelName: channel1,
							AtLeast:     iwfidl.PtrInt32(1),
						},
					},
				},
			})
			return

		case State4:
			// ZeroToAll: don't wait, consume all available (channel is empty)
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-3"),
							ChannelName: channel1,
							AtLeast:     iwfidl.PtrInt32(0),
						},
					},
				},
			})
			return

		case State5:
			// AtMostOnly: only AtMost set (no AtLeast), waits for late messages from S6
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							CommandId:   ptr.Any("cmd-4"),
							ChannelName: channel2,
							AtMost:      iwfidl.PtrInt32(2),
						},
					},
				},
			})
			return

		case State6:
			// Delayed publisher: wait 2s then publish 3 messages to ch2
			time.Sleep(time.Second * 2)
			publishes := make([]iwfidl.InterStateChannelPublishing, len(TestValuesCh2))
			for i := range TestValuesCh2 {
				v := TestValuesCh2[i]
				publishes[i] = iwfidl.InterStateChannelPublishing{
					ChannelName: channel2,
					Value:       &v,
				}
			}
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				PublishToInterStateChannel: publishes,
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

		switch req.GetWorkflowStateId() {
		case State1:
			// Publish 5 messages to ch and move to S2
			publishes := make([]iwfidl.InterStateChannelPublishing, len(TestValues))
			for i := range TestValues {
				v := TestValues[i]
				publishes[i] = iwfidl.InterStateChannelPublishing{
					ChannelName: channel1,
					Value:       &v,
				}
			}
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: State2},
					},
				},
				PublishToInterStateChannel: publishes,
			})
			return

		case State2:
			results := req.GetCommandResults()
			channelResult := results.GetInterStateChannelResults()[0]
			h.invokeData.Store("exactN_values", channelResult.MultiValues)

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: State3},
					},
				},
			})
			return

		case State3:
			results := req.GetCommandResults()
			channelResult := results.GetInterStateChannelResults()[0]
			h.invokeData.Store("oneToAll_values", channelResult.MultiValues)

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: State4},
					},
				},
			})
			return

		case State4:
			results := req.GetCommandResults()
			channelResult := results.GetInterStateChannelResults()[0]
			h.invokeData.Store("zeroToAll_values", channelResult.MultiValues)

			// Move to S5 (waiter) and S6 (delayed publisher) concurrently
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: State5},
						{StateId: State6},
					},
				},
			})
			return

		case State5:
			results := req.GetCommandResults()
			channelResult := results.GetInterStateChannelResults()[0]
			h.invokeData.Store("atMostOnly_values", channelResult.MultiValues)

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: service.GracefulCompletingWorkflowStateId},
					},
				},
			})
			return

		case State6:
			// Dead-end after publishing
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{StateId: service.DeadEndWorkflowStateId},
					},
				},
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	invokeHistory := make(map[string]int64)
	h.invokeHistory.Range(func(key, value interface{}) bool {
		invokeHistory[key.(string)] = value.(int64)
		return true
	})
	invokeData := make(map[string]interface{})
	h.invokeData.Range(func(key, value interface{}) bool {
		invokeData[key.(string)] = value
		return true
	})
	return invokeHistory, invokeData
}

func (h *handler) recordInvoke(key string) {
	if value, ok := h.invokeHistory.Load(key); ok {
		h.invokeHistory.Store(key, value.(int64)+1)
	} else {
		h.invokeHistory.Store(key, int64(1))
	}
}
