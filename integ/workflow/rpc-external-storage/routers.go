package rpcStorage

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
)

/**
 * Test workflow for RPC external storage functionality.
 * Tests updating data attributes with both small and large data via RPC methods.
 *
 * State1:
 *   - Sets up initial data attributes (small and large)
 *   - Waits for RPC to update the data attributes
 * State2:
 *   - Completes the workflow
 */

const (
	WorkflowType            = "rpc-external-storage"
	State1                  = "S1"
	State2                  = "S2"
	UpdateDataAttributesRPC = "update-data-attributes"

	SmallDataKey = "small-data"
	LargeDataKey = "large-data"

	// Small data stays in Temporal (under threshold)
	SmallDataContent = "small-data-content"

	// Initial data for testing
	InitialSmallDataContent = "initial-small-data"
)

var (
	// Large data goes to external storage (over threshold) - 1KB+
	LargeDataContent = "large-data-content-" + strings.Repeat("x", 1000)

	// Initial large data for testing - 1KB+
	InitialLargeDataContent = "initial-large-data-" + strings.Repeat("y", 1000)
)

var (
	SmallDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(SmallDataContent),
	}

	LargeDataValue = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(LargeDataContent),
	}

	InitialSmallData = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(InitialSmallDataContent),
	}

	InitialLargeData = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString(InitialLargeDataContent),
	}

	TestInput = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-input-value"),
	}

	TestOutput = iwfidl.EncodedObject{
		Encoding: iwfidl.PtrString("json"),
		Data:     iwfidl.PtrString("test-output-value"),
	}
)

type handler struct {
	testData sync.Map
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		testData: sync.Map{},
	}
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received workflow worker rpc request, ", req)

	wfCtx := req.Context
	if wfCtx.WorkflowId == "" || wfCtx.WorkflowRunId == "" {
		helpers.FailTestWithErrorMessage("invalid context in the request", t)
	}
	if req.WorkflowType != WorkflowType {
		helpers.FailTestWithErrorMessage(fmt.Sprintf("invalid workflow type: %s", req.WorkflowType), t)
	}

	// Store received data for verification
	h.testData.Store(req.RpcName+"-input", req.Input)
	h.testData.Store(req.RpcName+"-received-data", req.DataAttributes)

	if req.RpcName == UpdateDataAttributesRPC {
		// Verify we received the current data attributes (loaded from external storage)
		if req.DataAttributes != nil {
			for _, attr := range req.DataAttributes {
				log.Printf("Received data attribute: key=%s, hasData=%t, hasExtStoreId=%t",
					*attr.Key, attr.Value.Data != nil, attr.Value.ExtStoreId != nil)

				// Verify we received actual data content, not just external storage references
				if attr.Value.Data == nil {
					helpers.FailTestWithErrorMessage(fmt.Sprintf("RPC should receive actual data content for key %s, not external storage references", *attr.Key), t)
				}
			}
		}

		// Update data attributes with new values and send signal to close workflow
		c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
			Output: &TestOutput,
			UpsertDataAttributes: []iwfidl.KeyValue{
				{
					Key:   iwfidl.PtrString(SmallDataKey),
					Value: &SmallDataValue,
				},
				{
					Key:   iwfidl.PtrString(LargeDataKey),
					Value: &LargeDataValue,
				},
			},
			PublishToInterStateChannel: []iwfidl.InterStateChannelPublishing{
				{
					ChannelName: "close-workflow",
					Value: &iwfidl.EncodedObject{
						Encoding: iwfidl.PtrString("json"),
						Data:     iwfidl.PtrString("\"close\""),
					},
				},
			},
		})
		return
	}

	helpers.FailTestWithErrorMessage(fmt.Sprintf("unknown RPC name: %s", req.RpcName), t)
}

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if req.GetWorkflowStateId() == State1 {
		// Set up initial data attributes and wait for internal signal
		c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
			CommandRequest: &iwfidl.CommandRequest{
				DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
					{
						ChannelName: "close-workflow",
					},
				},
			},
			UpsertDataObjects: []iwfidl.KeyValue{
				{
					Key:   iwfidl.PtrString(SmallDataKey),
					Value: &InitialSmallData,
				},
				{
					Key:   iwfidl.PtrString(LargeDataKey),
					Value: &InitialLargeData,
				},
			},
		})
		return
	}

	if req.GetWorkflowStateId() == State2 {
		// Final state - no commands needed
		c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
			CommandRequest: &iwfidl.CommandRequest{
				DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if req.GetWorkflowStateId() == State1 {
		// Only complete workflow when we receive the close-workflow signal
		if req.CommandResults != nil &&
			req.CommandResults.InterStateChannelResults != nil &&
			len(req.CommandResults.InterStateChannelResults) > 0 {
			// We received the internal signal to close workflow
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
			// Should not happen - wait until signal is received
			c.JSON(http.StatusBadRequest, gin.H{"error": "Expected internal signal"})
		}
		return
	}

	if req.GetWorkflowStateId() == State2 {
		// Complete the workflow
		c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
			StateDecision: &iwfidl.StateDecision{
				NextStates: []iwfidl.StateMovement{
					{
						StateId: service.GracefulCompletingWorkflowStateId,
					},
				},
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	// Return empty history (not tracking state invocations for this test)
	history := make(map[string]int64)

	// Return test data collected from RPC calls
	testData := make(map[string]interface{})
	h.testData.Range(func(key, value interface{}) bool {
		testData[key.(string)] = value
		return true
	})
	return history, testData
}
