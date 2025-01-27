package deadend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"testing"
)

/**
 * This test workflow has 3 states, using REST controller to implement the workflow directly.
 *
 * RPCWriteData:
 *		- WaitUntil will upsert data attributes
 * RPCTriggerState:
 *		- WaitUntil will move to State1
 * State1:
 *		- WaitUntil is skipped
 *      - Execute method will put the state into a dead-end.
 */
const (
	WorkflowType    = "deadend"
	RPCTriggerState = "test-RPCTriggerState"
	RPCWriteData    = "RPCWriteData"

	State1 = "S1"
)

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

	wfCtx := req.Context
	if wfCtx.WorkflowId == "" || wfCtx.WorkflowRunId == "" {
		helpers.FailTestWithErrorMessage("invalid context in the request", t)
	}
	if req.WorkflowType != WorkflowType {
		helpers.FailTestWithErrorMessage("invalid workflow type", t)
	}

	if req.RpcName == RPCTriggerState {
		// Move to State 1
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			StateDecision: &iwfidl.StateDecision{NextStates: []iwfidl.StateMovement{
				{
					StateId: State1,
					StateOptions: &iwfidl.WorkflowStateOptions{
						SkipStartApi: iwfidl.PtrBool(true),
					},
				},
			}},
		})
	} else if req.RpcName == RPCWriteData {
		// Upsert data attributes
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			UpsertDataAttributes: []iwfidl.KeyValue{
				{
					Key: iwfidl.PtrString("any key"),
					Value: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("encoding"),
						Data:     iwfidl.PtrString("data"),
					},
				},
			},
		})
	} else {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("invalid rpc name: %s", req.RpcName), t)
	}
}

// ApiV1WorkflowStateStart - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	helpers.FailTestWithErrorMessage("should not be called", t)
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

		// Move to the dead-end state
		if req.GetWorkflowStateId() == State1 {

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
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
