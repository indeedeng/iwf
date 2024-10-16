package wf_state_options_data_attributes_loading

import (
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
)

/**
 * This test has four states
 *
 * State1: Creates all Data Attributes that will be used in this test
 * State2:
 * 		- Starts State with State Options containing WaitUntilApiDataAttributesLoadingPolicy
 * 		- Checks that (ApiV1WorkflowStateStart) WaitUntil method will load with expected Data Attributes
 * 		- Checks that (ApiV1WorkflowStateDecide) Execute method will not load any DataAttributes
 * State3:
 * 		- Starts State with State Options containing ExecuteApiDataAttributesLoadingPolicy
 * 		- Checks that (ApiV1WorkflowStateStart) WaitUntil method will not load any DataAttributes
 * 		- Checks that (ApiV1WorkflowStateDecide) Execute method will load with expected Data Attributes
 * State4:
 * 		- Starts State with State Options containing DataAttributesLoadingPolicy
 * 		- Checks that (ApiV1WorkflowStateStart) WaitUntil method will load with expected Data Attributes
 * 		- Checks that (ApiV1WorkflowStateDecide) Execute method will load with expected Data Attributes
 */
const (
	WorkflowType = "state_options_data_attributes_loading"
	State1       = "S1"
	State2       = "S2"
	State3       = "S3"
	State4       = "S4"
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

	log.Println("state_options_data_attributes_loading: received state decide request, ", req)

	h.invokeHistory[req.GetWorkflowStateId()+"_start"]++

	if req.GetWorkflowStateId() == State2 {
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
		verifyLoadedDataAttributes(req.GetWorkflowStateId(), req.GetDataObjects(), loadingType)
	}

	if req.GetWorkflowStateId() == State3 {
		verifyEmptyDataAttributes(req.GetDataObjects())
	}

	if req.GetWorkflowStateId() == State4 {
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
		verifyLoadedDataAttributes(req.GetWorkflowStateId(), req.GetDataObjects(), loadingType)
	}

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

	log.Println("state_options_data_attributes_loading: received state decide request, ", req)

	h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++

	var response iwfidl.WorkflowStateDecideResponse

	switch req.GetWorkflowStateId() {
	case State1:
		response = getState1DecideResponse(req)
	case State2:
		verifyEmptyDataAttributes(req.GetDataObjects())
		response = getState2DecideResponse(req)
	case State3:
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
		// Assert: State3 Execute method loads with expected Data Attributes
		verifyLoadedDataAttributes(req.GetWorkflowStateId(), req.GetDataObjects(), loadingType)
		response = getState3DecideResponse(req)
	case State4:
		loadingTypeFromInput := req.GetStateInput()
		loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
		// Assert: State4 Execute method loads with expected Data Attributes
		verifyLoadedDataAttributes(req.GetWorkflowStateId(), req.GetDataObjects(), loadingType)
		response = getState4DecideResponse()
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

	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State2,
					StateOptions: &iwfidl.WorkflowStateOptions{
						WaitUntilApiDataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"da_wait_until1"},
						},
						ExecuteApiDataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &noneLoadingType,
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
		UpsertDataObjects: getUpsertDataObjects(),
	}
}

func getState2DecideResponse(req iwfidl.WorkflowStateDecideRequest) iwfidl.WorkflowStateDecideResponse {
	loadingTypeFromInput := req.GetStateInput()
	loadingType := iwfidl.PersistenceLoadingType(loadingTypeFromInput.GetData())
	noneLoadingType := iwfidl.NONE

	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State3,
					StateOptions: &iwfidl.WorkflowStateOptions{
						WaitUntilApiDataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &noneLoadingType,
						},
						ExecuteApiDataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"da_execute1"},
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

	return iwfidl.WorkflowStateDecideResponse{
		StateDecision: &iwfidl.StateDecision{
			NextStates: []iwfidl.StateMovement{
				{
					StateId: State4,
					StateOptions: &iwfidl.WorkflowStateOptions{
						DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: &loadingType,
							PartialLoadingKeys:     []string{"da_other_key"},
						},
					},
					StateInput: &loadingTypeFromInput,
				},
			},
		},
	}
}

func getState4DecideResponse() iwfidl.WorkflowStateDecideResponse {
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

func verifyEmptyDataAttributes(dataAttributes []iwfidl.KeyValue) {
	var expectedDataAttributes []iwfidl.KeyValue
	if !assert.ElementsMatch(common.DummyT{}, expectedDataAttributes, dataAttributes) {
		panic("Data attributes should be empty")
	}
}

func verifyLoadedDataAttributes(stateId string, dataAttributes []iwfidl.KeyValue, loadingType iwfidl.PersistenceLoadingType) {
	expectedDataAttributes := getExpectedDataAttributes(stateId, loadingType)
	if !assert.ElementsMatch(common.DummyT{}, expectedDataAttributes, dataAttributes) {
		panic("Data attributes should be the same")
	}
}

func getUpsertDataObjects() []iwfidl.KeyValue {
	return []iwfidl.KeyValue{
		{
			Key:   iwfidl.PtrString("da_wait_until1"),
			Value: &iwfidl.EncodedObject{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("test-data-object-wait-until")},
		},
		{
			Key:   iwfidl.PtrString("da_execute1"),
			Value: &iwfidl.EncodedObject{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("test-data-object-execute")},
		},
		{
			Key:   iwfidl.PtrString("da_other_key"),
			Value: &iwfidl.EncodedObject{Encoding: iwfidl.PtrString("json"), Data: iwfidl.PtrString("random-value")},
		},
	}
}

func getExpectedDataAttributes(stateId string, loadingType iwfidl.PersistenceLoadingType) []iwfidl.KeyValue {
	if stateId == State2 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_wait_until1"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-wait-until"),
				},
			},
		}
	}
	if stateId == State3 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_execute1"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("test-data-object-execute"),
				},
			},
		}
	}

	if stateId == State4 && (loadingType == iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK || loadingType == iwfidl.PARTIAL_WITHOUT_LOCKING) {
		return []iwfidl.KeyValue{
			{
				Key: iwfidl.PtrString("da_other_key"),
				Value: &iwfidl.EncodedObject{
					Encoding: iwfidl.PtrString("json"),
					Data:     iwfidl.PtrString("random-value"),
				},
			},
		}
	}

	return getUpsertDataObjects()
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, nil
}
