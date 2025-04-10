package persistence_loading_policy

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/helpers"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"sync"
	"testing"
)

/**
 * This test workflow has two states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil skipped
 * 		- Execute method verifies the loaded attributes then moves to a dead-end.
 * State2:
 * 		- WaitUntil method verifies the loaded attributes
 * 		- Execute method verifies the loaded attributes then gracefully completes the workflow
 */
const (
	WorkflowType = "persistence_loading_policy"
	State1       = "S1"
	State2       = "S2"
)

type handler struct {
	invokeHistory sync.Map
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: sync.Map{},
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("persistence_loading_policy: received state start request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_start"); ok {
		h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", value.(int64)+1)
	} else {
		h.invokeHistory.Store(req.GetWorkflowStateId()+"_start", int64(1))
	}

	if req.GetWorkflowStateId() == State2 {
		// Dynamically get the loadingType from input
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

		verifyLoadedAttributes(t, req.GetSearchAttributes(), req.GetDataObjects(), loadingType)
	}

	// Go straight to the decide methods without any commands
	c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
		CommandRequest: &iwfidl.CommandRequest{
			DeciderTriggerType: iwfidl.ANY_COMMAND_COMPLETED.Ptr(),
		},
	})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("persistence_loading_policy: received state decide request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if value, ok := h.invokeHistory.Load(req.GetWorkflowStateId() + "_decide"); ok {
		h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", value.(int64)+1)
	} else {
		h.invokeHistory.Store(req.GetWorkflowStateId()+"_decide", int64(1))
	}

	// Dynamically get the loadingType from input
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	if req.GetWorkflowStateId() == State2 {
		verifyLoadedAttributes(t, req.GetSearchAttributes(), req.GetDataObjects(), loadingType)
	}

	var upsertSearchAttributes []iwfidl.SearchAttribute
	var upsertDataAttributes []iwfidl.KeyValue

	// Set search attributes and data attributes in State1
	if req.GetWorkflowStateId() == State1 {
		upsertSearchAttributes = []iwfidl.SearchAttribute{
			{
				Key:         iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType:   ptr.Any(iwfidl.KEYWORD),
				StringValue: iwfidl.PtrString("test-search-attribute-1"),
			},
			{
				Key:         iwfidl.PtrString(persistence.TestSearchAttributeTextKey),
				ValueType:   ptr.Any(iwfidl.TEXT),
				StringValue: iwfidl.PtrString("test-search-attribute-2"),
			},
		}

		upsertDataAttributes = []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_1"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-value1"),
				},
			},
			{
				Key: iwfidl.PtrString("da_2"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-value2"),
				},
			},
		}
	}

	// Move to dead-end (state 1) or completion (state 2)
	var nextStateId string
	if req.GetWorkflowStateId() == State1 {
		nextStateId = service.DeadEndWorkflowStateId
	} else if req.GetWorkflowStateId() == State2 {
		nextStateId = service.GracefulCompletingWorkflowStateId
	}

	c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
		StateDecision:          getStateDecision(nextStateId, loadingTypeFromInput, loadingType),
		UpsertSearchAttributes: upsertSearchAttributes,
		UpsertDataObjects:      upsertDataAttributes,
	})
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context, t *testing.T) {
	var req iwfidl.WorkflowWorkerRpcRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("persistence_loading_policy: received rpc request, ", req)

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	if value, ok := h.invokeHistory.Load("rpc"); ok {
		h.invokeHistory.Store("rpc", value.(int64)+1)
	} else {
		h.invokeHistory.Store("rpc", int64(1))
	}

	// dynamically get the loadingType from input
	loadingTypeFromInput := req.GetInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	verifyLoadedAttributes(t, req.GetSearchAttributes(), req.GetDataAttributes(), loadingType)

	c.JSON(http.StatusOK, iwfidl.WorkflowWorkerRpcResponse{
		StateDecision: getStateDecision(State2, loadingTypeFromInput, loadingType),
	})

}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	invokeHistory := make(map[string]int64)
	h.invokeHistory.Range(func(key, value interface{}) bool {
		invokeHistory[key.(string)] = value.(int64)
		return true
	})
	return invokeHistory, nil
}
func verifyLoadedAttributes(
	t *testing.T,
	searchAttributes []iwfidl.SearchAttribute,
	dataAttributes []iwfidl.KeyValue,
	loadingType iwfidl.PersistenceLoadingType) {

	var expectedSearchAttributes []iwfidl.SearchAttribute
	var expectedDataAttributes []iwfidl.KeyValue

	if loadingType == iwfidl.ALL_WITHOUT_LOCKING || loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
		expectedSearchAttributes = []iwfidl.SearchAttribute{
			{
				Key:         iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType:   ptr.Any(iwfidl.KEYWORD),
				StringValue: iwfidl.PtrString("test-search-attribute-1"),
			},
			{
				Key:         iwfidl.PtrString(persistence.TestSearchAttributeTextKey),
				ValueType:   ptr.Any(iwfidl.TEXT),
				StringValue: iwfidl.PtrString("test-search-attribute-2"),
			},
		}
		expectedDataAttributes = []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_1"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-value1"),
				},
			},
			{
				Key: iwfidl.PtrString("da_2"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-value2"),
				},
			},
		}
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING {
		expectedSearchAttributes = []iwfidl.SearchAttribute{
			{
				Key:         iwfidl.PtrString(persistence.TestSearchAttributeKeywordKey),
				ValueType:   ptr.Any(iwfidl.KEYWORD),
				StringValue: iwfidl.PtrString("test-search-attribute-1"),
			},
		}
		expectedDataAttributes = []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_1"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-value1"),
				},
			},
		}
	} else if loadingType == iwfidl.NONE {
		expectedSearchAttributes = []iwfidl.SearchAttribute{}
		expectedDataAttributes = []iwfidl.KeyValue{}
	}

	// use ElementsMatch so that the order won't be a problem.
	// Internally the SAs are stored as a map and as a result, Golang return it without ordering guarantee
	if !assert.ElementsMatch(common.DummyT{}, expectedSearchAttributes, searchAttributes) {
		helpers.FailTestWithErrorMessage("Search attributes should be the same", t)
	}

	if !assert.ElementsMatch(common.DummyT{}, expectedDataAttributes, dataAttributes) {
		helpers.FailTestWithErrorMessage("Data attributes should be the same", t)
	}
}

func getStateDecision(nextStateId string, loadingTypeFromInput iwfidl.EncodedObject, loadingType iwfidl.PersistenceLoadingType) *iwfidl.StateDecision {
	return &iwfidl.StateDecision{
		NextStates: []iwfidl.StateMovement{
			{
				StateId: nextStateId,
				StateOptions: &iwfidl.WorkflowStateOptions{
					SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
						PersistenceLoadingType: &loadingType,
						PartialLoadingKeys: []string{
							persistence.TestSearchAttributeKeywordKey,
						},
						LockingKeys: []string{
							persistence.TestSearchAttributeTextKey,
						},
					},
					DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
						PersistenceLoadingType: &loadingType,
						PartialLoadingKeys: []string{
							"da_1",
						},
						LockingKeys: []string{
							"da_2",
						},
					},
				},
				StateInput: &loadingTypeFromInput,
			},
		},
	}
}
