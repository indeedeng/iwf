package wf_execute_api_fail_and_proceed

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
)

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

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	panic("should not get here")
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
	}
	if req.WorkflowStateId == State1 {
		c.JSON(http.StatusBadRequest, map[string]string{"error": "test-error"})
		return
	}
	if req.WorkflowStateId == StateRecover {
		input := req.StateInput
		if input.GetData() == InputData && input.GetEncoding() == InputDataEncoding {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{})
		} else {
			panic("input is not correct")
		}
		return
	}

	panic("should not get here")
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
