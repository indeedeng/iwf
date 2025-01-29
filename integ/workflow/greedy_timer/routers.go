package greedy_timer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
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
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

type Input struct {
	Durations []int64 `json:"durations"`
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

	wfCtx := req.Context
	if wfCtx.WorkflowId == "" || wfCtx.WorkflowRunId == "" {
		panic("invalid context in the request")
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

	panic("invalid rpc name:" + req.RpcName)
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {

		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
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

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {

		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		if req.GetWorkflowStateId() == ScheduleTimerState {
			h.invokeData["completed_state_id"] = req.GetContext().StateExecutionId
			results := req.GetCommandResults()
			timerResults := results.GetTimerResults()
			h.invokeData["completed_timer_id"] = timerResults[0].CommandId

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
	return h.invokeHistory, h.invokeData
}
