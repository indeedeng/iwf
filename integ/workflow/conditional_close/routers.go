package conditional_close

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"testing"
	"time"
)

/**
 * This test workflow has 1 state, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil will proceed when the channel or signal is published to
 *      - Execute method will continuously retry State1 until the 3rd attempt which will send a message to the channel or
 *        signal, making the state empty and force-complete.
 */
const (
	WorkflowType              = "conditional_close"
	RpcPublishInternalChannel = "publish_internal_channel"

	TestChannelName = "test-channel-name"

	State1 = "S1"
)

var TestInput = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data"),
}

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)
	h.invokeHistory[req.RpcName]++

	// Return channel name
	c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
		PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
			{
				ChannelName: TestChannelName,
			},
		},
	})
}

// ApiV1WorkflowStateStart - for a workflow
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
			// Proceed when channel is published to
			cmdReq := &iwfidl.CommandRequest{
				InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
					{
						ChannelName: TestChannelName,
					},
				},
				CommandWaitingType: ptr.Any(iwfidl.ANY_COMPLETED),
			}
			input := req.GetStateInput()

			if input.GetData() == "use-signal-channel" {
				// Proceed when signal is published to
				cmdReq = &iwfidl.CommandRequest{
					SignalCommands: []iwfidl.SignalCommand{
						{
							SignalChannelName: TestChannelName,
						},
					},
					CommandWaitingType: ptr.Any(iwfidl.ANY_COMPLETED),
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: cmdReq,
			})
			return
		}
	}

	helpers.FailTestWithErrorMessage("error request", t)
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

			var internalChanPub []iwfidl.InterStateChannelPublishing
			context := req.GetContext()
			if context.GetStateExecutionId() == "S1-1" {
				// Wait for 3 seconds so that the channel can have a new message
				time.Sleep(time.Second * 3)
			} else if context.GetStateExecutionId() == "S1-3" {
				// Send internal channel message within the state execution
				// and expecting the messages are processed by the conditional check
				internalChanPub = []iwfidl.InterStateChannelPublishing{
					{
						ChannelName: TestChannelName,
						Value:       &TestInput,
					}}
			}

			conditionalClose := &iwfidl.WorkflowConditionalClose{
				ConditionalCloseType: iwfidl.FORCE_COMPLETE_ON_INTERNAL_CHANNEL_EMPTY.Ptr(),
				ChannelName:          iwfidl.PtrString(TestChannelName),
				CloseInput:           &TestInput,
			}
			input := req.GetStateInput()

			// Use signal instead
			if input.GetData() == "use-signal-channel" {
				conditionalClose = &iwfidl.WorkflowConditionalClose{
					ConditionalCloseType: iwfidl.FORCE_COMPLETE_ON_SIGNAL_CHANNEL_EMPTY.Ptr(),
					ChannelName:          iwfidl.PtrString(TestChannelName),
					CloseInput:           &TestInput,
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				PublishToInterStateChannel: internalChanPub,
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State1,
							StateInput: req.StateInput,
						},
					},
					ConditionalClose: conditionalClose,
				},
			})
			return
		}
	}

	helpers.FailTestWithErrorMessage("error request", t)
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
