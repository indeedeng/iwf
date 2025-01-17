package wf_state_options_search_attributes_loading

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
)

/**
 * This test workflow has five states, using REST controller to implement the workflow directly.
 *
 * State1:
 *		- WaitUntil method does nothing
 * 		- Execute method sets values for all Search Attributes used in this test
 * 			- CustomKeywordField
 * 			- CustomStringField
 *			- CustomBoolField
 * State2:
 * 		- State Options contains WaitUntilApiSearchAttributesLoadingPolicy
 * 		- WaitUntil method asserts that expected SearchAttributes are loaded
 * 		- Execute method asserts that no SearchAttributes are loaded
 * State3:
 * 		- State Options contains WaitUntilApiSearchAttributesLoadingPolicy
 * 		- WaitUntil method asserts that no SearchAttributes are loaded
 * 		- Execute method asserts that expected SearchAttributes are loaded
 * State4:
 * 		- State Options contains SearchAttributesLoadingPolicy
 * 		- WaitUntil method asserts that expected SearchAttributes are loaded
 * 		- Execute method asserts that expected SearchAttributes are loaded
 * State5:
 * 		- State Options contains AdSearchAttributesLoadingPolicy and WaitUntilApiSearchAttributesLoadingPolicy
 * 		- WaitUntil method asserts that WaitUntilApiSearchAttributesLoadingPolicy are loaded
 * 		- Execute method asserts that SearchAttributesLoadingPolicy are loaded
 */
const (
	WorkflowType = "state_options_search_attributes_loading"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
	State4       = "S4"
	State5       = "S5"
)

type handler struct {
	invokeHistory map[string]int64
}

func NewHandler() common.WorkflowHandlerWithRpc {
	return &handler{
		invokeHistory: make(map[string]int64),
	}
}

func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	log.Println("state_options_search_attributes_loading: received state decide request, ", req)

	h.invokeHistory[req.GetWorkflowStateId()+"_start"]++

	currentMethod := "WaitUntil"
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	if req.GetWorkflowStateId() == State2 || req.GetWorkflowStateId() == State4 || req.GetWorkflowStateId() == State5 {
		verifyLoadedSearchAttributes(req.GetWorkflowStateId(), currentMethod, req.GetSearchAttributes(), loadingType)
	}

	if req.GetWorkflowStateId() == State3 {
		verifyEmptySearchAttributes(req.GetSearchAttributes())
	}

	// Go straight to the decide methods without any commands
	c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
		CommandRequest: &iwfidl.CommandRequest{
			DeciderTriggerType: iwfidl.ANY_COMMAND_COMPLETED.Ptr(),
		},
	})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.GetWorkflowType() != WorkflowType {
		c.JSON(http.StatusBadRequest, struct{}{})
		return
	}

	log.Println("state_options_search_attributes_loading: received state decide request, ", req)

	h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

	currentMethod := "Execute"
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	var response iwfidl.WorkflowStateDecideResponse
	switch req.GetWorkflowStateId() {
	case State1:
		response = getState1DecideResponse(req)
	case State2:
		verifyEmptySearchAttributes(req.GetSearchAttributes())
		response = getState2DecideResponse(req)
	case State3:
		verifyLoadedSearchAttributes(req.GetWorkflowStateId(), currentMethod, req.GetSearchAttributes(), loadingType)
		response = getState3DecideResponse(req)
	case State4:
		verifyLoadedSearchAttributes(req.GetWorkflowStateId(), currentMethod, req.GetSearchAttributes(), loadingType)
		response = getState4DecideResponse(req)
	case State5:
		verifyLoadedSearchAttributes(req.GetWorkflowStateId(), currentMethod, req.GetSearchAttributes(), loadingType)
		response = getState5DecideResponse()
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) ApiV1WorkflowWorkerRpc(c *gin.Context) {
	c.JSON(http.StatusBadRequest, struct{}{})
}

func getState1DecideResponse(req iwfidl.WorkflowStateDecideRequest) iwfidl.WorkflowStateDecideResponse {
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
	noneLoadingType := iwfidl.NONE

	// Move to State 2 with the provided options & input after updating data objects
	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State2,
					StateOptions: &iwfidl.WorkflowStateOptions{
						WaitUntilApiSearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"CustomKeywordField"},
						},
						ExecuteApiSearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &noneLoadingType,
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
		UpsertSearchAttributes: getUpsertSearchAttributes(),
	}
}

func getState2DecideResponse(req iwfidl.WorkflowStateDecideRequest) iwfidl.WorkflowStateDecideResponse {
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
	noneLoadingType := iwfidl.NONE

	// Move to State 3 with the provided options & input
	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State3,
					StateOptions: &iwfidl.WorkflowStateOptions{
						WaitUntilApiSearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &noneLoadingType,
						},
						ExecuteApiSearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"CustomStringField"},
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
	}
}

func getState3DecideResponse(req iwfidl.WorkflowStateDecideRequest) iwfidl.WorkflowStateDecideResponse {
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	// Move to State 4 with the provided options & input
	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State4,
					StateOptions: &iwfidl.WorkflowStateOptions{
						SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"CustomBoolField"},
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
	}
}

func getState4DecideResponse(req iwfidl.WorkflowStateDecideRequest) iwfidl.WorkflowStateDecideResponse {
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())

	// Move to State 5 with the provided options & input
	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State5,
					StateOptions: &iwfidl.WorkflowStateOptions{
						WaitUntilApiSearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"CustomKeywordField"},
						},
						SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"CustomBoolField"},
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
	}
}

func getState5DecideResponse() iwfidl.WorkflowStateDecideResponse {
	// Move to completion
	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: service.GracefulCompletingWorkflowStateId,
				},
			},
		},
	}
}

func verifyEmptySearchAttributes(searchAttributes []iwfidl.SearchAttribute) {
	var expectedSearchAttributes []iwfidl.SearchAttribute
	if !assert.ElementsMatch(common.DummyT{}, expectedSearchAttributes, searchAttributes) {
		panic("Search attributes should be empty")
	}
}

func verifyLoadedSearchAttributes(stateId string, method string, searchAttributes []iwfidl.SearchAttribute, loadingType iwfidl.PersistenceLoadingType) {
	expectedSearchAttributes := getExpectedSearchAttributes(stateId, method, loadingType)
	if !assert.ElementsMatch(common.DummyT{}, expectedSearchAttributes, searchAttributes) {
		panic("Search attributes should be the same")
	}
}

func getUpsertSearchAttributes() []iwfidl.SearchAttribute {
	return []iwfidl.SearchAttribute{
		{
			Key:              iwfidl.PtrString("CustomKeywordField"),
			ValueType:        ptr.Any(iwfidl.KEYWORD_ARRAY),
			StringArrayValue: []string{"keyword1", "keyword2"},
		},
		{
			Key:         iwfidl.PtrString("CustomStringField"),
			ValueType:   ptr.Any(iwfidl.TEXT),
			StringValue: iwfidl.PtrString("I am a string"),
		},
		{
			Key:       iwfidl.PtrString("CustomBoolField"),
			ValueType: ptr.Any(iwfidl.BOOL),
			BoolValue: ptr.Any(true),
		},
	}
}

func getExpectedSearchAttributes(stateId string, method string, loadingType iwfidl.PersistenceLoadingType) []iwfidl.SearchAttribute {
	if stateId == State2 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.SearchAttribute{
			{
				Key:              iwfidl.PtrString("CustomKeywordField"),
				ValueType:        ptr.Any(iwfidl.KEYWORD_ARRAY),
				StringArrayValue: []string{"keyword1", "keyword2"},
			},
		}
	}
	if stateId == State3 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.SearchAttribute{
			{
				Key:         iwfidl.PtrString("CustomStringField"),
				ValueType:   ptr.Any(iwfidl.TEXT),
				StringValue: iwfidl.PtrString("I am a string"),
			},
		}
	}

	if stateId == State4 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.SearchAttribute{
			{
				Key:       iwfidl.PtrString("CustomBoolField"),
				ValueType: ptr.Any(iwfidl.BOOL),
				BoolValue: ptr.Any(true),
			},
		}
	}

	if stateId == State5 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		switch method {
		case "WaitUntil":
			return []iwfidl.SearchAttribute{
				{
					Key:              iwfidl.PtrString("CustomKeywordField"),
					ValueType:        ptr.Any(iwfidl.KEYWORD_ARRAY),
					StringArrayValue: []string{"keyword1", "keyword2"},
				},
			}
		case "Execute":
			return []iwfidl.SearchAttribute{
				{
					Key:       iwfidl.PtrString("CustomBoolField"),
					ValueType: ptr.Any(iwfidl.BOOL),
					BoolValue: ptr.Any(true),
				},
			}
		}
	}

	return getUpsertSearchAttributes()
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, nil
}
