package parallel

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	WorkflowType = "parallel"
	State1       = "S1"
	State11      = "S11"
	State12      = "S12"
	State111     = "S111"
	State112     = "S112"
	State121     = "S121"
	State122     = "S122"
)

func NewParallelWorkflow() (*Handler, *gin.Engine) {
	router := gin.Default()

	handler := newHandler()

	router.POST(service.StateStartApi, handler.apiV1WorkflowStateStart)
	router.POST(service.StateDecideApi, handler.apiV1WorkflowStateDecide)

	return handler, router
}

type Handler struct {
	invokeHistory map[string]int
}

func newHandler() *Handler {
	return &Handler{
		invokeHistory: make(map[string]int),
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *Handler) apiV1WorkflowStateStart(c *gin.Context) {
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
				DeciderTriggerType: service.DeciderTypeAllCommandCompleted,
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *Handler) apiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		var nexts []string
		switch req.GetWorkflowStateId() {
		case State1:
			// cause graceful complete to wait
			time.Sleep(time.Second * 1)
			nexts = append(nexts, State11, State12)
		case State11:
			// cause graceful complete to wait
			time.Sleep(time.Second * 2)
			nexts = append(nexts, State111, State112)
		case State12:
			// cause graceful complete to wait
			time.Sleep(time.Second * 2)
			nexts = append(nexts, State121, State122)
		case State111:
			nexts = append(nexts, service.GracefulCompletingWorkflowStateId)
		case State112, State121, State122:
			//empty
		default:
			nexts = append(nexts, service.ForceFailingWorkflowStateId)
		}

		var nextStates []iwfidl.StateMovement
		for _, nextId := range nexts {
			nextStates = append(nextStates, iwfidl.StateMovement{
				StateId: nextId,
			})
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

func (h *Handler) GetTestResult() map[string]int {
	return h.invokeHistory
}
