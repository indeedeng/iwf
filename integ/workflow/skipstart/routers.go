package skipstart

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
)

/**
 * This test workflow has 2 states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- Wait until is skipped.
 *      - Execute method will go to State2
 * State2:
 *		- Wait until is skipped.
 *      - Execute method will gracefully complete workflow
 */
const (
	WorkflowType = "skipstart"
	State1       = "S1"
	State2       = "S2"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() *handler {
	return &handler{
		invokeHistory: make(map[string]int64),
	}
}

// ApiV1WorkflowStateStart - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	panic("start API should be skipped")
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
			// Move to State 2 with the provided input & options
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    State2,
							StateInput: req.StateInput,
							StateOptions: &iwfidl.WorkflowStateOptions{
								SkipStartApi: iwfidl.PtrBool(true),
							},
						},
					},
				},
			})
			return
		} else if req.GetWorkflowStateId() == State2 {
			// Move to completion with the provided input
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId:    service.GracefulCompletingWorkflowStateId,
							StateInput: req.StateInput,
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
	return h.invokeHistory, nil
}
