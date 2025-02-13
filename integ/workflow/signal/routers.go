package signal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"sync"
	"testing"
)

/**
 * This test workflow has two states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil waits until 4 signals are received
 * 		- Execute method publishes the 4 signals & moves to State2
 * State2:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType                  = "signal"
	State1                        = "S1"
	State2                        = "S2"
	SignalName                    = "test-signal-name"
	InternalChannelName           = "test-internal-channel-name"
	UnhandledSignalName           = "test-unhandled-signal-name"
	RPCNameGetSignalChannelInfo   = "RPCNameGetSignalChannelInfo"
	RPCNameGetInternalChannelInfo = "RPCNameGetInternalChannelInfo"
)

var StateOptionsForLargeDataAttributes = iwfidl.WorkflowStateOptions{
	DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
		PersistenceLoadingType: ptr.Any(iwfidl.NONE),
	},
}

type handler struct {
	invokeHistory sync.Map
	invokeData    sync.Map
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: sync.Map{},
		invokeData:    sync.Map{},
	}
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.RpcName == RPCNameGetSignalChannelInfo {
		signalInfos := req.SignalChannelInfos
		data, err := json.Marshal(signalInfos)
		if err != nil {
			helpers.FailTestWithError(err, t)
		}
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
				{
					ChannelName: InternalChannelName,
				},
			},
			Output: &iwfidl.EncodedObject{
				Data: ptr.Any(string(data)),
			},
		})
		return
	}
	if req.RpcName == RPCNameGetInternalChannelInfo {
		icInfos := req.InternalChannelInfos
		data, err := json.Marshal(icInfos)
		if err != nil {
			helpers.FailTestWithError(err, t)
		}
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			Output: &iwfidl.EncodedObject{
				Data: ptr.Any(string(data)),
			},
		})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{})
	return
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
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
		}

		if req.GetWorkflowStateId() == State1 {
			// Proceed when 4 signals are received
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					SignalCommands: []iwfidl.SignalCommand{
						{
							CommandId:         ptr.Any("signal-cmd-id0"),
							SignalChannelName: SignalName,
						},
						{
							CommandId:         ptr.Any("signal-cmd-id1"),
							SignalChannelName: SignalName,
						},
						{
							SignalChannelName: SignalName,
						},
						{
							SignalChannelName: SignalName,
						},
					},
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
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
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_decide"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", int64(1))
		}

		if req.GetWorkflowStateId() == State1 {
			signalResults := req.GetCommandResults()

			// Publish 4 signals
			for i := 0; i < 4; i++ {
				signalId := signalResults.SignalResults[i].GetCommandId()
				signalValue := signalResults.SignalResults[i].GetSignalValue()

				h.invokeData.Store(fmt.Sprintf("signalId%v", i), signalId)
				h.invokeData.Store(fmt.Sprintf("signalValue%v", i), signalValue)
			}

			// Move to State 2
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:      State2,
							StateOptions: &StateOptionsForLargeDataAttributes,
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
