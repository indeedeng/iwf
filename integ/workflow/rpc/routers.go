package rpc

import (
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
 *		- WaitUntil updates attribute data and data objects and then waits until the channel has been published to
 * 		- Execute method moves to State2
 * State2:
 *		- WaitUntil method does nothing
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType              = "rpc"
	State1                    = "S1"
	State2                    = "S2"
	TestInterStateChannelName = "test-TestInterStateChannelName"
	RPCName                   = "test-RPCName"
	RPCNameReadOnly           = "test-RPC-readonly"
	RPCNameError              = "test-RPC-error"

	TestDataAttributeKey = "test-data-attribute"

	TestSearchAttributeKeywordKey    = "CustomKeywordField"
	TestSearchAttributeKeywordValue1 = "keyword-value1"
	TestSearchAttributeKeywordValue2 = "keyword-value2"

	TestSearchAttributeIntKey    = "CustomIntField"
	TestSearchAttributeBoolKey   = "CustomBoolField"
	TestSearchAttributeIntValue1 = 1
	TestSearchAttributeIntValue2 = 2

	WorkerApiErrorDetails = "test-details"
	WorkerApiErrorType    = "test-type"
)

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

var TestDataAttributeVal1 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-attribute-value1"),
}

var TestDataAttributeVal2 = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-data-attribute-value2"),
}

var TestInput = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-input-value"),
}

var TestOutput = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-output-value"),
}

var TestRecordEvent = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-record-event-value"),
}

var TestInterstateChannelValue = iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("test-interstatechannel-value"),
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)

	wfCtx := req.Context
	if wfCtx.WorkflowId == "" || wfCtx.WorkflowRunId == "" {
		helpers.FailTestWithErrorMessage("invalid context in the request", t)
	}
	if req.WorkflowType != WorkflowType ||
		(req.RpcName != RPCName && req.RpcName != RPCNameReadOnly && req.RpcName != RPCNameError) {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("invalid rpc name: %s", req.RpcName), t)
	}

	h.invokeData.Store(req.RpcName+"-input", req.Input)
	h.invokeData.Store(req.RpcName+"-search-attributes", req.SearchAttributes)
	h.invokeData.Store(req.RpcName+"-data-attributes", req.DataAttributes)

	if req.RpcName == RPCNameReadOnly {
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			Output: &TestOutput,
		})
		return
	}
	if req.RpcName == RPCNameError {
		c.JSON(http.StatusBadGateway, iwfidl.WorkerErrorResponse{
			Detail:    iwfidl.PtrString(WorkerApiErrorDetails),
			ErrorType: iwfidl.PtrString(WorkerApiErrorType),
		})
		return
	}

	upsertSAs := []iwfidl.SearchAttribute{
		{
			Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
			StringValue: iwfidl.PtrString(TestSearchAttributeKeywordValue2),
			ValueType:   ptr.Any(iwfidl.KEYWORD),
		},
		{
			Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
			IntegerValue: iwfidl.PtrInt64(TestSearchAttributeIntValue2),
			ValueType:    ptr.Any(iwfidl.INT),
		},
	}

	// Proceed with State 2 after setting the attributes
	c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
		Output: &TestOutput,
		StateDecision: &iwfidl.StateDecision{NextStates: []iwfidl.StateMovement{
			{
				StateId: State2,
			},
		}},
		UpsertSearchAttributes: upsertSAs,
		UpsertDataAttributes: []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString(TestDataAttributeKey),
				Value: &TestDataAttributeVal2,
			},
		},
		RecordEvents: []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString("test-key"),
				Value: &TestRecordEvent,
			},
		},
		PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
			{
				ChannelName: TestInterStateChannelName,
				Value:       &TestInterstateChannelValue,
			},
		},
	})
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
			upsertSAs := []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
					StringValue: iwfidl.PtrString(TestSearchAttributeKeywordValue1),
					ValueType:   ptr.Any(iwfidl.KEYWORD),
				},
				{
					Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
					IntegerValue: iwfidl.PtrInt64(TestSearchAttributeIntValue1),
					ValueType:    ptr.Any(iwfidl.INT),
				},
				{
					Key:       iwfidl.PtrString(TestSearchAttributeBoolKey),
					ValueType: ptr.Any(iwfidl.BOOL),
					BoolValue: iwfidl.PtrBool(false),
				},
			}

			// Proceed after attributes and data objects have been updated and channel has been published to
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							ChannelName: TestInterStateChannelName,
						},
					},
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				UpsertSearchAttributes: upsertSAs,
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key:   iwfidl.PtrString(TestDataAttributeKey),
						Value: &TestDataAttributeVal1,
					},
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
			commandRes := req.GetCommandResults()
			res := commandRes.GetInterStateChannelResults()[0]
			if res.GetRequestStatus() != iwfidl.RECEIVED || res.GetChannelName() != TestInterStateChannelName {
				helpers.FailTestWithErrorMessage("the signal should be received", t)
			}
			h.invokeData.Store(TestInterStateChannelName, res.Value)

			// Move to state 2
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
