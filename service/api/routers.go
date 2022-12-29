package api

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/service/common/log"
)

const WorkflowStartApiPath = "/api/v1/workflow/start"
const WorkflowSignalApiPath = "/api/v1/workflow/signal"
const WorkflowGetDataObjectsApiPath = "/api/v1/workflow/dataobjects/get"
const WorkflowGetSearchAttributesApiPath = "/api/v1/workflow/searchattributes/get"
const WorkflowGetApiPath = "/api/v1/workflow/get"
const WorkflowGetWithWaitApiPath = "/api/v1/workflow/getWithWait"
const WorkflowSearchApiPath = "/api/v1/workflow/search"
const WorkflowResetApiPath = "/api/v1/workflow/reset"
const WorkflowStopApiPath = "/api/v1/workflow/stop"

// NewService returns a new router.
func NewService(client UnifiedClient, logger log.Logger) *gin.Engine {
	router := gin.Default()

	handler := newHandler(client, logger)

	router.GET("/", handler.index)
	router.POST(WorkflowStartApiPath, handler.apiV1WorkflowStart)
	router.POST(WorkflowSignalApiPath, handler.apiV1WorkflowSignal)
	router.POST(WorkflowGetDataObjectsApiPath, handler.apiV1WorkflowGetDataObjects)
	router.POST(WorkflowGetSearchAttributesApiPath, handler.apiV1WorkflowGetSearchAttributes)
	router.POST(WorkflowGetApiPath, handler.apiV1WorkflowGet)
	router.POST(WorkflowGetWithWaitApiPath, handler.apiV1WorkflowGetWithWait)
	router.POST(WorkflowSearchApiPath, handler.apiV1WorkflowSearch)
	router.POST(WorkflowResetApiPath, handler.apiV1WorkflowReset)
	router.POST(WorkflowStopApiPath, handler.apiV1WorkflowStop)

	return router
}
