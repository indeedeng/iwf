/*
 * Workflow APIs
 *
 * This APIs for iwf SDKs to operate workflows
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"fmt"
	"github.com/cadence-oss/iwf-server/gen/server/workflow"
	temporalimpl "github.com/cadence-oss/iwf-server/service/interpreter/temporalImpl"

	"github.com/cadence-oss/iwf-server/gen/client/workflow/state"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handler struct {
	temporalClient client.Client
}

func newHandler() *handler {
	// The client is a heavyweight object that should be created once per process.
	// TODO use config for connection options and merge with api handler
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	return &handler{
		temporalClient: temporalClient,
	}
}

func (h *handler) close() {
	h.temporalClient.Close()
}

// Index is the index handler.
func (h *handler) index(c *gin.Context) {
	// for test only, will be removed
	runTestRestApi()

	c.String(http.StatusOK, "Hello World!")
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) apiV1WorkflowStartPost(c *gin.Context) {
	var req workflow.WorkflowStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received request", req)

	workflowOptions := client.StartWorkflowOptions{
		ID:        req.WorkflowId,
		TaskQueue: temporalimpl.TaskQueue,
	}

	we, err := h.temporalClient.ExecuteWorkflow(context.Background(), workflowOptions, temporalimpl.Interpreter, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	c.JSON(http.StatusOK, workflow.WorkflowStartResponse{
		WorkflowRunId: we.GetRunID(),
	})
}

func runTestRestApi() {
	apiClient := state.NewAPIClient(&state.Configuration{})
	req := apiClient.DefaultApi.ApiV1WorkflowStateStartPost(context.Background())
	wfType := "123"
	resp, httpResp, err := req.WorkflowStateStartRequest(state.WorkflowStateStartRequest{
		WorkflowType: &wfType,
	}).Execute()
	fmt.Println("test REST API", resp.GetCommandRequest(), httpResp, err)
}
