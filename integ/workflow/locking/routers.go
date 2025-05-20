package locking

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
	"strconv"
	"sync"
	"testing"
	"time"
)

/**
 * This test workflow has three states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method does nothing
 * 		- Execute method will move to State Waiting, and 10 instances of State 2
 * State2:
 * 		- WaitUntil update SA
 * 		- Execute method will update data attributes and will gracefully complete workflow
 * StateWaiting:
 * 		- WaitUntil will proceed once the internal channel has been published to
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType                  = "locking"
	State1                        = "S1"
	State2                        = "S2"
	StateWaiting                  = "StateWaiting"
	TestDataAttributeKey1         = "test-data-attribute-1"
	TestDataAttributeKey2         = "test-data-attribute-2"
	RPCName                       = "increase-counter"
	InternalChannelName           = "test-channel"
	TestSearchAttributeKeywordKey = "CustomKeywordField"
	TestSearchAttributeIntKey     = "CustomIntField"

	ShouldUnblockStateWaiting = "shouldUnblockStateWaiting"

	InParallelS2 = 10

	NumUnusedSignals = 4

	UnusedSignalChannelName   = "test-unused-signal-channel"
	UnusedInternalChannelName = "test-unused-internal-channel"
)

var TestValue = &iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString("data"),
}

var UnblockValue = &iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("json"),
	Data:     iwfidl.PtrString(ShouldUnblockStateWaiting),
}

var state2Options = &iwfidl.WorkflowStateOptions{
	SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
		PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
		PartialLoadingKeys: []string{
			TestSearchAttributeIntKey,
			TestSearchAttributeKeywordKey,
		},
		LockingKeys: []string{
			TestSearchAttributeIntKey,
		},
	},
	DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
		PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
		PartialLoadingKeys: []string{
			TestDataAttributeKey1,
			TestDataAttributeKey2,
		},
		LockingKeys: []string{
			TestDataAttributeKey1,
		},
	},
}

var state2Movement = iwfidl.StateMovement{
	StateId:      State2,
	StateOptions: state2Options,
}

type handler struct {
	invokeHistory sync.Map
	rpcInvokes    int32
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: sync.Map{},
	}
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)

	if req.WorkflowType != WorkflowType || (req.RpcName != RPCName) {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("invalid rpc name: %s", req.RpcName), t)
	}

	input := req.Input
	if input.GetEncoding() != TestValue.GetEncoding() {
		helpers.FailTestWithErrorMessage("input is incorrect", t)
	}

	// Publish to internal channel
	if input.GetData() == ShouldUnblockStateWaiting {
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
				{
					ChannelName: InternalChannelName,
					Value:       TestValue,
				},
			},
		})
		return
	}

	signalChannelInfo := (*req.SignalChannelInfos)[UnusedSignalChannelName]
	if signalChannelInfo.GetSize() != NumUnusedSignals {
		// the 4 messages are sent from the beginning of "locking_test"
		helpers.FailTestWithErrorMessage("incorrect signal channel size", t)
	}

	if h.rpcInvokes > 0 {
		internalChannelInfo := (*req.InternalChannelInfos)[UnusedInternalChannelName]
		if h.rpcInvokes != internalChannelInfo.GetSize() {
			helpers.FailTestWithErrorMessage("incorrect internal channel size", t)
		}
	}
	h.rpcInvokes++

	time.Sleep(time.Millisecond)

	// This RPC will increase both SA and DA
	saInt := int64(0)
	for _, sa := range req.GetSearchAttributes() {
		if sa.GetKey() == TestSearchAttributeIntKey {
			saInt = sa.GetIntegerValue()
		}
	}
	saInt++

	context := req.GetContext()
	upsertSearchAttributes := []iwfidl.SearchAttribute{
		{
			Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
			StringValue: iwfidl.PtrString(context.GetStateExecutionId()),
			ValueType:   ptr.Any(iwfidl.KEYWORD),
		},
		{
			Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
			IntegerValue: iwfidl.PtrInt64(saInt),
			ValueType:    ptr.Any(iwfidl.INT),
		},
	}

	daInt := 0
	for _, da := range req.DataAttributes {
		if da.GetKey() == TestDataAttributeKey1 {
			value := da.GetValue()
			data := value.GetData()
			if data != "" {
				i, err := strconv.ParseInt(data, 10, 32)
				if err != nil {
					helpers.FailTestWithError(err, t)
				}
				daInt = int(i)
			}
		}
	}
	daInt++

	upsertDataAttributes := []iwfidl.KeyValue{
		{
			Key: iwfidl.PtrString(TestDataAttributeKey1),
			Value: &iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString(fmt.Sprintf("%v", daInt)),
			},
		},
		{
			Key: iwfidl.PtrString(TestDataAttributeKey2),
			Value: &iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString(context.GetStateExecutionId()),
			},
		},
	}

	response := iwfidl.WorkflowWorkerRpcResponse{
		Output: TestValue,
		StateDecision: &iwfidl.StateDecision{NextStates: []iwfidl.StateMovement{
			state2Movement,
		}},
		UpsertSearchAttributes: upsertSearchAttributes,
		UpsertDataAttributes:   upsertDataAttributes,
		RecordEvents: []iwfidl.KeyValue{
			{
				Key:   iwfidl.PtrString("test-key"),
				Value: TestValue,
			},
		},
		PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
			{
				ChannelName: UnusedInternalChannelName,
				Value:       TestValue,
			},
		},
	}
	c.JSON(http.StatusOK, response)

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
		if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
		} else {
			h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
		}

		// Go straight to the decide methods without any commands
		if req.GetWorkflowStateId() == State1 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		// Will proceed once the internal channel has been published to
		if req.GetWorkflowStateId() == StateWaiting {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							ChannelName: InternalChannelName,
						},
					},
				},
			})
			return

		}
		if req.GetWorkflowStateId() == State2 {
			// This state API is to increase SA
			time.Sleep(time.Second)
			saInt := int64(0)
			for _, sa := range req.GetSearchAttributes() {
				if sa.GetKey() == TestSearchAttributeIntKey {
					saInt = sa.GetIntegerValue()
				}
			}
			saInt++

			var sa []iwfidl.SearchAttribute
			context := req.GetContext()
			sa = []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
					StringValue: iwfidl.PtrString(context.GetStateExecutionId()),
					ValueType:   ptr.Any(iwfidl.KEYWORD),
				},
				{
					Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
					IntegerValue: iwfidl.PtrInt64(saInt),
					ValueType:    ptr.Any(iwfidl.INT),
				},
			}

			// Go straight to the decide methods after updating SA
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				UpsertSearchAttributes: sa,
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

			stms := []iwfidl.StateMovement{
				{
					StateId: StateWaiting,
				},
			}
			for i := 0; i < InParallelS2; i++ {
				stms = append(stms, state2Movement)
			}

			// Move to State Waiting, and 10 instances of State 2
			// State Waiting will not complete until the internal channel has been published to
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: stms,
				},
			})
			return
		}
		// Move to completion
		if req.GetWorkflowStateId() == StateWaiting {
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
		if req.GetWorkflowStateId() == State2 {
			// This API is to increase DA
			time.Sleep(time.Second)
			daInt := 0
			for _, da := range req.DataObjects {
				if da.GetKey() == TestDataAttributeKey1 {
					value := da.GetValue()
					data := value.GetData()
					if data != "" {
						i, err := strconv.ParseInt(data, 10, 32)
						if err != nil {
							helpers.FailTestWithError(err, t)
						}
						daInt = int(i)
					}
				}
			}
			daInt++
			context := req.GetContext()

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key: iwfidl.PtrString(TestDataAttributeKey1),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(fmt.Sprintf("%v", daInt)),
						},
					},
					{
						Key: iwfidl.PtrString(TestDataAttributeKey2),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(context.GetStateExecutionId()),
						},
					},
				},

				// Move to completion
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
	return invokeHistory, nil
}
