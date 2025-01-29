package common

import (
	"github.com/gin-gonic/gin"
	"testing"
)

type WorkflowHandler interface {
	ApiV1WorkflowStateStart(c *gin.Context, t *testing.T)
	ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T)
	GetTestResult() (map[string]int64, map[string]interface{})
}

type WorkflowHandlerWithRpc interface {
	WorkflowHandler
	ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T)
}
