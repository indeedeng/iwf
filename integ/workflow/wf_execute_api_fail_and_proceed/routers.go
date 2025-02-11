package wf_execute_api_fail_and_proceed

import (
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
)

/**
 * This test workflow has one state, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method is skipped
 *      - Execute method will intentionally fail
 * StateRecover:
 *		- Execute method will gracefully complete workflow
 */
const (
	WorkflowType      = "wf_execute_api_fail_and_proceed"
	State1            = "S1"
	StateRecover      = "Recover"
	InputData         = "test-data"
	InputDataEncoding = "test-encoding"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	helpers.FailTestWithErrorMessage("should not get here", t)
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
	}
	input := req.StateInput
	if req.WorkflowStateId == State1 {
		if input.GetData() == InputData && input.GetEncoding() == InputDataEncoding {
			c.JSON(http.StatusBadRequest, map[string]string{"error": "test-error"})
		} else {
			helpers.FailTestWithErrorMessage("input is not correct: "+input.GetData()+", "+input.GetEncoding(), t)
		}

		return
	}
	if req.WorkflowStateId == StateRecover {
		if input.GetData() == InputData && input.GetEncoding() == InputDataEncoding {
			// Move to completion
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: []iwfidl.StateMovement{
						{
							StateId: service.GracefulCompletingWorkflowStateId,
						},
					},
				},
			})
		} else {
			helpers.FailTestWithErrorMessage("input is not correct: "+input.GetData()+", "+input.GetEncoding(), t)
		}
		return
	}

	helpers.FailTestWithErrorMessage("should not get here", t)
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
