package persistence_loading_policy

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/integ/workflow/persistence"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
)

const (
	WorkflowType = "persistence_loading_policy"
	State1       = "S1"
	State2       = "S2"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
	}
}

// ApiV1WorkflowStartPost - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
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

	h.invokeHistory[req.GetWorkflowStateId()+"_start"]++

	if req.GetWorkflowStateId() == State2 {
		// dynamically get the loadingType from input
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

		verifyLoadedAttributes(req.GetSearchAttributes(), req.GetDataObjects(), loadingType)
	}

	// go straight to decide methods without any commands
	c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
		CommandRequest: &iwfidl.CommandRequest{
			DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
		},
	})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
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

	h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

	// dynamically get the loadingType from input
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	if req.GetWorkflowStateId() == State2 {
		verifyLoadedAttributes(req.GetSearchAttributes(), req.GetDataObjects(), loadingType)
	}

	var upsertSearchAttributes []iwfidl.SearchAttribute
	var upsertDataObjects []iwfidl.KeyValue

	// set search attributes and data attributes in State1
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

		upsertDataObjects = []iwfidl.KeyValue{
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

	var nextStateId string
	if req.GetWorkflowStateId() == State1 {
		nextStateId = State2
	} else if req.GetWorkflowStateId() == State2 {
		nextStateId = service.GracefulCompletingWorkflowStateId
	}

	c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
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
		},
		UpsertSearchAttributes: upsertSearchAttributes,
		UpsertDataObjects:      upsertDataObjects,
	})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, nil
}

func verifyLoadedAttributes(
	searchAttributes []iwfidl.SearchAttribute,
	dataAttributes []iwfidl.KeyValue,
	loadingType iwfidl.PersistenceLoadingType) {

	var expectedSearchAttributes []iwfidl.SearchAttribute
	var expectedDataAttributes []iwfidl.KeyValue

	if loadingType == iwfidl.ALL_WITHOUT_LOCKING {
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
	} else if loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING || loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK {
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
		panic("Search attributes should be the same")
	}

	if !assert.ElementsMatch(common.DummyT{}, expectedDataAttributes, dataAttributes) {
		panic("Data attributes should be the same")
	}
}
