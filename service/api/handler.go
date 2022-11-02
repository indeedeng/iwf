package api

import (
	"log"
	"net/http"

	"github.com/indeedeng/iwf/gen/iwfidl"

	"github.com/gin-gonic/gin"
)

type handler struct {
	svc ApiService
}

func newHandler(client UnifiedClient) *handler {
	svc, err := NewApiService(client)
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

	resp, errResp := h.svc.ApiV1WorkflowStartPost(req)
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

	errResp := h.svc.ApiV1WorkflowSignalPost(req)
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

	resp, errResp := h.svc.ApiV1WorkflowSearchPost(req)
	if errResp != nil {
		h.processError(c, errResp)
		return
	}
	c.JSON(http.StatusOK, resp)
	return
}

func (h *handler) apiV1WorkflowGetQueryAttributes(c *gin.Context) {
	var req iwfidl.WorkflowGetQueryAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, errResp := h.svc.ApiV1WorkflowGetQueryAttributesPost(req)
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

	resp, errResp := h.svc.ApiV1WorkflowGetSearchAttributesPost(req)
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
		resp, errResp = h.svc.ApiV1WorkflowGetWithWaitPost(req)
	} else {
		resp, errResp = h.svc.ApiV1WorkflowGetPost(req)
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

	resp, errResp := h.svc.ApiV1WorkflowResetPost(req)
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
