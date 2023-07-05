package conditional_close

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"time"
)

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

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)
	h.invokeHistory[req.RpcName]++

	c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
		PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
			{
				ChannelName: TestChannelName,
			},
		},
	})
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
					InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
						{
							ChannelName: TestChannelName,
						},
					},
					CommandWaitingType: ptr.Any(iwfidl.ANY_COMPLETED),
				},
			})
			return
		}
	}

	panic("error request")
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
			if context.GetStateExecutionId() == "S1-1" {
				// wait for 3 seconds so that the channel can have a new message
				time.Sleep(time.Second * 3)
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: State1,
						},
					},
					ConditionalClose: &iwfidl.WorkflowConditionalClose{
						ConditionalCloseType: iwfidl.FORCE_COMPLETE_ON_INTERNAL_CHANNEL_EMPTY.Ptr(),
						ChannelName:          iwfidl.PtrString(TestChannelName),
						CloseInput:           &TestInput,
					},
				},
			})
			return
		}
	}

	panic("error request")
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
