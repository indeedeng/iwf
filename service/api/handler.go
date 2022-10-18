package api

import (
	"context"
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
	"github.com/cadence-oss/iwf-server/service"
	"time"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	client UnifiedClient
}

func newHandler(client UnifiedClient) *handler {
	return &handler{
		client: client,
	}
}

func (h *handler) close() {
	h.client.Close()
}

// Index is the index handler.
func (h *handler) index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World from iWF server!")
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) apiV1WorkflowStartPost(c *gin.Context) {
	var req iwfidl.WorkflowStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	workflowOptions := StartWorkflowOptions{
		ID:                 req.GetWorkflowId(),
		TaskQueue:          service.TaskQueue,
		WorkflowRunTimeout: time.Duration(req.WorkflowTimeoutSeconds) * time.Second,
	}

	input := service.InterpreterWorkflowInput{
		IwfWorkflowType: req.GetIwfWorkflowType(),
		IwfWorkerUrl:    req.GetIwfWorkerUrl(),
		StartStateId:    req.GetStartStateId(),
		StateInput:      req.GetStateInput(),
		StateOptions:    req.GetStateOptions(),
	}
	runId, err := h.client.StartInterpreterWorkflow(context.Background(), workflowOptions, input)
	if err != nil {
		handleError(c, err)
		return
	}

	log.Println("Started workflow", "WorkflowID", req.WorkflowId, "RunID", runId)

	c.JSON(http.StatusOK, iwfidl.WorkflowStartResponse{
		WorkflowRunId: iwfidl.PtrString(runId),
	})
}

func handleError(c *gin.Context, err error) {
	// TODO differentiate different error for different codes
	log.Println("encounter error for API", err)
	c.JSON(http.StatusInternalServerError, iwfidl.ErrorResponse{
		Detail: iwfidl.PtrString(err.Error()),
	})
}

func (h *handler) apiV1WorkflowSignalPost(c *gin.Context) {
	var req iwfidl.WorkflowSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	err := h.client.SignalWorkflow(context.Background(),
		req.GetWorkflowId(), req.GetWorkflowRunId(), req.GetSignalChannelName(), req.GetSignalValue())
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, struct{}{})
}

func (h *handler) apiV1WorkflowSearchPost(c *gin.Context) {
	var req iwfidl.WorkflowSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	pageSize := int32(1000)
	if req.GetPageSize() > 0 {
		pageSize = req.GetPageSize()
	}
	resp, err := h.client.ListWorkflow(context.Background(), &ListWorkflowExecutionsRequest{
		PageSize: pageSize,
		Query:    req.GetQuery(),
	})
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, iwfidl.WorkflowSearchResponse{
		WorkflowExecutions: resp.Executions,
	})
}

func (h *handler) apiV1WorkflowQueryPost(c *gin.Context) {
	var req iwfidl.WorkflowGetQueryAttributesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	var queryResult1 service.QueryAttributeResponse
	err := h.client.QueryWorkflow(context.Background(), &queryResult1,
		req.GetWorkflowId(), req.GetWorkflowRunId(), service.AttributeQueryType,
		service.QueryAttributeRequest{
			Keys: req.AttributeKeys,
		})

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, iwfidl.WorkflowGetQueryAttributesResponse{
		QueryAttributes: queryResult1.AttributeValues,
	})
}

func (h *handler) apiV1WorkflowGetPost(c *gin.Context) {
	h.doApiV1WorkflowGetPost(c, false)
}

func (h *handler) apiV1WorkflowGetWithWaitPost(c *gin.Context) {
	h.doApiV1WorkflowGetPost(c, true)
}

func (h *handler) doApiV1WorkflowGetPost(c *gin.Context, waitIfStillRunning bool) {
	var req iwfidl.WorkflowGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	resp, err := h.client.DescribeWorkflowExecution(context.Background(), req.GetWorkflowId(), req.GetWorkflowRunId())
	if err != nil {
		handleError(c, err)
		return
	}

	var output service.InterpreterWorkflowOutput
	if req.GetNeedsResults() || waitIfStillRunning {
		if resp.Status == service.WorkflowStatusCompleted || waitIfStillRunning {
			err := h.client.GetWorkflowResult(context.Background(), &output, req.GetWorkflowId(), req.GetWorkflowRunId())
			if err != nil {
				handleError(c, err)
				return
			}
		}
	}

	status := resp.Status
	if waitIfStillRunning {
		// override because when GetWorkflowResult, the workflow is completed
		status = service.WorkflowStatusCompleted
	}

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, iwfidl.WorkflowGetResponse{
		WorkflowRunId:  resp.RunId,
		WorkflowStatus: status,
		Results:        output.StateCompletionOutputs,
	})
}

func (h *handler) apiV1WorkflowResetPost(c *gin.Context) {
	var req iwfidl.WorkflowResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received API request", req)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()
	runId, err := h.client.ResetWorkflow(ctx, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, iwfidl.WorkflowResetResponse{
		WorkflowRunId: runId,
	})
}
