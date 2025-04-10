package api

import (
	"github.com/indeedeng/iwf/config"
	"net/http"

	"github.com/indeedeng/iwf/service"
	uclient "github.com/indeedeng/iwf/service/client"
	"github.com/indeedeng/iwf/service/common/errors"
	"github.com/indeedeng/iwf/service/common/log"
	"github.com/indeedeng/iwf/service/common/log/tag"

	"github.com/indeedeng/iwf/gen/iwfidl"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc    ApiService
	logger log.Logger
}

func newHandler(config config.Config, client uclient.UnifiedClient, logger log.Logger) *handler {
	svc, err := NewApiService(config, client, service.TaskQueue, logger)
	if err != nil {
		panic(err)
	}
	return &handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *handler) close() {
	h.svc.Close()
}

// Index is the index handler.
func (h *handler) index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from iWF server!")
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) apiV1WorkflowStart(c *gin.Context) {
	var req iwfidl.WorkflowStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowStartPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowWaitForStateCompletion(c *gin.Context) {
	var req iwfidl.WorkflowWaitForStateCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowWaitForStateCompletion(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowSignal(c *gin.Context) {
	var req iwfidl.WorkflowSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowSignalPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowStop(c *gin.Context) {
	var req iwfidl.WorkflowStopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowStopPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowInternalDump(c *gin.Context) {
	var req iwfidl.WorkflowDumpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowDumpPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowConfigUpdate(c *gin.Context) {
	var req iwfidl.WorkflowConfigUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowConfigUpdate(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowTriggerContinueAsNew(c *gin.Context) {
	var req iwfidl.TriggerContinueAsNewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowTriggerContinueAsNew(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowSearch(c *gin.Context) {
	var req iwfidl.WorkflowSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowSearchPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowRpc(c *gin.Context) {
	var req iwfidl.WorkflowRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowRpcPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) infoHealthCheck(c *gin.Context) {
	h.logger.Debug("received Health check request")

	resp := h.svc.ApiInfoHealth(c.Request.Context())
	c.JSON(http.StatusOK, resp)

	return
}

func (h *handler) apiV1WorkflowGetDataAttributes(c *gin.Context) {
	var req iwfidl.WorkflowGetDataObjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowGetQueryAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowSetDataAttributes(c *gin.Context) {
	var req iwfidl.WorkflowSetDataObjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowSetQueryAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowGetSearchAttributes(c *gin.Context) {
	var req iwfidl.WorkflowGetSearchAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowGetSearchAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowSetSearchAttributes(c *gin.Context) {
	var req iwfidl.WorkflowSetSearchAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	errResp := h.svc.ApiV1WorkflowSetSearchAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func (h *handler) apiV1WorkflowGet(c *gin.Context) {
	h.doApiV1WorkflowGetPost(c, false)
}

func (h *handler) apiV1WorkflowGetWithWait(c *gin.Context) {
	h.doApiV1WorkflowGetPost(c, true)
}

func (h *handler) doApiV1WorkflowGetPost(c *gin.Context, waitIfStillRunning bool) {
	var req iwfidl.WorkflowGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	var resp *iwfidl.WorkflowGetResponse
	var errResp *errors.ErrorAndStatus
	if waitIfStillRunning {
		resp, errResp = h.svc.ApiV1WorkflowGetWithWaitPost(c.Request.Context(), req)
	} else {
		resp, errResp = h.svc.ApiV1WorkflowGetPost(c.Request.Context(), req)
	}

	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowReset(c *gin.Context) {
	var req iwfidl.WorkflowResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	h.logger.Debug("received API request", tag.Value(log.ToJsonAndTruncateForLogging(req)))

	resp, errResp := h.svc.ApiV1WorkflowResetPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowSkipTimer(c *gin.Context) {
	var req iwfidl.WorkflowSkipTimerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		invalidRequestSchema(c)
		return
	}
	errResp := h.svc.ApiV1WorkflowSkipTimerPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
	return
}

func invalidRequestSchema(c *gin.Context) {
	c.JSON(http.StatusBadRequest, iwfidl.ErrorResponse{
		Detail: iwfidl.PtrString("invalid request schema"),
	})
}

func (h *handler) processError(c *gin.Context, resp *errors.ErrorAndStatus) {
	c.JSON(resp.StatusCode, resp.Error)
}
