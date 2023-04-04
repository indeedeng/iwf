package parallel

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"time"
)

const (
	WorkflowType = "parallel"
	State1       = "S1"
	State11      = "S11"
	State12      = "S12"
	State13      = "S13"
	State111     = "S111"
	State112     = "S112"
	State121     = "S121"
	State122     = "S122"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
	}
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
		c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
			CommandRequest: &iwfidl.CommandRequest{
				DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED,
			},
		})
		return
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
		var nextStates []iwfidl.StateMovement
		switch req.GetWorkflowStateId() {
		case State1:
			// cause graceful complete to wait
			time.Sleep(time.Second * 1)

			nextStates = []iwfidl.StateMovement{
				{
					StateId: State11,
				},
				{
					StateId: State12,
				},
				{
					StateId: State13,
				},
			}
		case State11:
			// cause graceful complete to wait
			time.Sleep(time.Second * 2)

			nextStates = []iwfidl.StateMovement{
				{
					StateId: State111,
				},
				{
					StateId: State112,
				},
			}
		case State12:
			// cause graceful complete to wait
			time.Sleep(time.Second * 2)
			nextStates = []iwfidl.StateMovement{
				{
					StateId: State121,
				},
				{
					StateId: State122,
				},
			}
		case State13:
			// cause graceful complete to wait
			time.Sleep(time.Second * 1)
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.GracefulCompletingWorkflowStateId,
					StateInput: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("from " + req.GetWorkflowStateId()),
					},
				},
			}
		case State112, State121, State122, State111:
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.GracefulCompletingWorkflowStateId,
					StateInput: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("from " + req.GetWorkflowStateId()),
					},
				},
			}
		default:
			nextStates = []iwfidl.StateMovement{
				{
					StateId: service.ForceFailingWorkflowStateId,
				},
			}
		}

		c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
			StateDecision: &iwfidl.StateDecision{
				NextStates: nextStates,
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, nil
}
