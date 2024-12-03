package locking

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	WorkflowType                  = "locking"
	State1                        = "S1"
	State2                        = "S2"
	StateWaiting                  = "StateWaiting"
	TestDataObjectKey1            = "test-data-object-1"
	TestDataObjectKey2            = "test-data-object-2"
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
			TestDataObjectKey1,
			TestDataObjectKey2,
		},
		LockingKeys: []string{
			TestDataObjectKey1,
		},
	},
}

var state2Movement = iwfidl.StateMovement{
	StateId:      State2,
	StateOptions: state2Options,
}

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
	rpcInvokes    int32
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)

	if req.WorkflowType != WorkflowType || (req.RpcName != RPCName) {
		panic("invalid rpc name:" + req.RpcName)
	}

	input := req.Input
	if input.GetEncoding() != TestValue.GetEncoding() {
		panic("input is incorrect")
	}
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
		panic("incorrect signal channel size")
	}
	if h.rpcInvokes > 0 {
		internalChannelInfo := (*req.InternalChannelInfos)[UnusedInternalChannelName]
		if h.rpcInvokes != internalChannelInfo.GetSize() {
			panic("incorrect internal channel size")
		}
	}
	h.rpcInvokes++

	// this RPC will increase both SA and DA
	time.Sleep(time.Millisecond)

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
		if da.GetKey() == TestDataObjectKey1 {
			value := da.GetValue()
			data := value.GetData()
			if data != "" {
				i, err := strconv.ParseInt(data, 10, 32)
				if err != nil {
					panic(err)
				}
				daInt = int(i)
			}
		}
	}
	daInt++

	upsertDataAttributes := []iwfidl.KeyValue{
		{
			Key: iwfidl.PtrString(TestDataObjectKey1),
			Value: &iwfidl.EncodedObject{
				Encoding: iwfidl.PtrString("json"),
				Data:     iwfidl.PtrString(fmt.Sprintf("%v", daInt)),
			},
		},
		{
			Key: iwfidl.PtrString(TestDataObjectKey2),
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
			// this state API is to increase SA
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

			stms := []iwfidl.StateMovement{
				{
					StateId: StateWaiting,
				},
			}
			for i := 0; i < InParallelS2; i++ {
				stms = append(stms, state2Movement)
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: stms,
				},
			})
			return
		}
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
			// this API is to increase DA
			time.Sleep(time.Second)
			daInt := 0
			for _, da := range req.DataObjects {
				if da.GetKey() == TestDataObjectKey1 {
					value := da.GetValue()
					data := value.GetData()
					if data != "" {
						i, err := strconv.ParseInt(data, 10, 32)
						if err != nil {
							panic(err)
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
						Key: iwfidl.PtrString(TestDataObjectKey1),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(fmt.Sprintf("%v", daInt)),
						},
					},
					{
						Key: iwfidl.PtrString(TestDataObjectKey2),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(context.GetStateExecutionId()),
						},
					},
				},

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
	return h.invokeHistory, h.invokeData
}
