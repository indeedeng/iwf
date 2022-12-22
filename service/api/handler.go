package api

import (
	"github.com/indeedeng/iwf/service"
	"log"
	"net/http"

	"github.com/indeedeng/iwf/gen/iwfidl"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc ApiService
}

func newHandler(client UnifiedClient) *handler {
	svc, err := NewApiService(client, service.TaskQueue)
	if err != nil {
		panic(err)
	}
	return &handler{
		svc: svc,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowStartPost(c.Request.Context(), req)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	errResp := h.svc.ApiV1WorkflowStopPost(c.Request.Context(), req)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowSearchPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowGetDataObjects(c *gin.Context) {
	var req iwfidl.WorkflowGetDataObjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowGetQueryAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowGetSearchAttributes(c *gin.Context) {
	var req iwfidl.WorkflowGetSearchAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowGetSearchAttributesPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	var resp *iwfidl.WorkflowGetResponse
	var errResp *ErrorAndStatus
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowResetPost(c.Request.Context(), req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) processError(c *gin.Context, resp *ErrorAndStatus) {
	c.JSON(resp.StatusCode, resp.Error)
}
