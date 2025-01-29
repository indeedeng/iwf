package wf_force_fail

import (
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
 * This test workflow has one state, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method does nothing
 *      - Execute method will intentionally force-fail
 */
const (
	WorkflowType = "wf_force_fail"
	State1       = "S1"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

var TestData = &iwfidl.EncodedObject{
	Encoding: iwfidl.PtrString("test-encoding"),
	Data:     iwfidl.PtrString("test-data"),
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			// Empty response
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{})
			return
		}
	}

	panic("should not get here")
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType && req.GetWorkflowStateId() == State1 {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		// Force fail
		c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
			StateDecision: &iwfidl.StateDecision{
				NextStates: []iwfidl.StateMovement{
					{
						StateId:    service.ForceFailingWorkflowStateId,
						StateInput: TestData,
					},
				},
			},
		})
		return
	}

	helpers.FailTestWithErrorMessage("should not get here", t)
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
