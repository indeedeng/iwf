package greedy_timer

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
	"strconv"
	"sync"
	"testing"
)

/*
*
This workflow will accept an array of integers representing durations and execute a state that waits on a timer corresponding to each duration provided
*/
const (
	WorkflowType       = "greedy_timer"
	ScheduleTimerState = "schedule"
	SubmitDurationsRPC = "submitDurationsRPC"
)

type handler struct {
	invokeHistory sync.Map
	invokeData    sync.Map
}

type Input struct {
	Durations []int64 `json:"durations"`
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
	log.Println("received workflow worker rpc request, ", req)

	wfCtx := req.Context
	if wfCtx.WorkflowId == "" || wfCtx.WorkflowRunId == "" {
		helpers.FailTestWithErrorMessage("invalid context in the request", t)
	}
	if req.WorkflowType != WorkflowType {
		panic("invalid WorkflowType:" + req.WorkflowType)
	}

	if req.RpcName == SubmitDurationsRPC {

		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			StateDecision: &iwfidl.StateDecision{NextStates: []iwfidl.StateMovement{
				{
					StateId:    ScheduleTimerState,
					StateInput: req.Input,
				},
			}},
		})
		return
	}

	helpers.FailTestWithErrorMessage(fmt.Sprintf("invalid rpc name: %s", req.RpcName), t)
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

		if req.GetWorkflowStateId() == ScheduleTimerState {

			var input Input
			err := json.Unmarshal([]byte(req.StateInput.GetData()), &input)
			if err != nil {
				panic(err)
			}

			timers := make([]iwfidl.TimerCommand, len(input.Durations))
			for i, duration := range input.Durations {
				timers[i] = iwfidl.TimerCommand{
					CommandId:       ptr.Any("duration-" + strconv.FormatInt(duration, 10)),
					DurationSeconds: iwfidl.PtrInt64(duration),
				}
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					TimerCommands:      timers,
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

		if req.GetWorkflowStateId() == ScheduleTimerState {
			h.invokeData.Store("completed_state_id", req.GetContext().StateExecutionId)
			results := req.GetCommandResults()
			timerResults := results.GetTimerResults()
			h.invokeData.Store("completed_timer_id", timerResults[0].CommandId)

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.ForceCompletingWorkflowStateId,
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
