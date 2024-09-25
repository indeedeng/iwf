package common

import "github.com/gin-gonic/gin"

type WorkflowHandler interface {
	ApiV1WorkflowStateStart(c *gin.Context)
	ApiV1WorkflowStateDecide(c *gin.Context)
	GetTestResult() (map[string]int64, map[string]interface{})
}

type WorkflowHandlerWithRpc interface {
	WorkflowHandler
	ApiV1WorkflowWorkerRpc(c *gin.Context)
}
