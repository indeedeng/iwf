package locking

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/integ/workflow/common"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"log"
	"net/http"
	"strconv"
)

const (
	WorkflowType       = "locking"
	State1             = "S1"
	State2             = "S2"
	TestDataObjectKey1 = "test-data-object-1"
	TestDataObjectKey2 = "test-data-object-2"

	TestSearchAttributeKeywordKey = "CustomKeywordField"
	TestSearchAttributeIntKey     = "CustomIntField"
)

type handler struct {
	invokeHistory map[string]int64
	invokeData    map[string]interface{}
}

func NewHandler() common.WorkflowHandler {
	return &handler{
		invokeHistory: make(map[string]int64),
		invokeData:    make(map[string]interface{}),
	}
}

// ApiV1WorkflowStateStart - for a workflow
func (h *handler) ApiV1WorkflowStateStart(c *gin.Context) {
	var req iwfidl.WorkflowStateStartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state start request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_start"]++
		if req.GetWorkflowStateId() == State1 {
			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State2 {
			saInt := int64(0)
			for _, sa := range req.GetSearchAttributes() {
				if sa.GetKey() == TestSearchAttributeIntKey {
					saInt = sa.GetIntegerValue()
				}
			}
			saInt++

			var sa []iwfidl.SearchAttribute
			context := req.GetContext()
			sa = []iwfidl.SearchAttribute{
				{
					Key:         iwfidl.PtrString(TestSearchAttributeKeywordKey),
					StringValue: iwfidl.PtrString(context.GetStateExecutionId()),
					ValueType:   ptr.Any(iwfidl.KEYWORD),
				},
				{
					Key:          iwfidl.PtrString(TestSearchAttributeIntKey),
					IntegerValue: iwfidl.PtrInt64(saInt),
					ValueType:    ptr.Any(iwfidl.INT),
				},
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateStartResponse{
				CommandRequest: &iwfidl.CommandRequest{
					DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
				},
				UpsertSearchAttributes: sa,
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) ApiV1WorkflowStateDecide(c *gin.Context) {
	var req iwfidl.WorkflowStateDecideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println("received state decide request, ", req)

	if req.GetWorkflowType() == WorkflowType {
		h.invokeHistory[req.GetWorkflowStateId()+"_decide"]++
		if req.GetWorkflowStateId() == State1 {

			var stms []iwfidl.StateMovement
			for i := 0; i < 10; i++ {
				stms = append(stms, iwfidl.StateMovement{
					StateId: State2,
					StateOptions: &iwfidl.WorkflowStateOptions{
						SearchAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
							PartialLoadingKeys: []string{
								TestSearchAttributeIntKey,
								TestSearchAttributeKeywordKey,
							},
							LockingKeys: []string{
								TestSearchAttributeIntKey,
							},
						},
						DataAttributesLoadingPolicy: &iwfidl.PersistenceLoadingPolicy{
							PersistenceLoadingType: iwfidl.PARTIAL_WITH_EXCLUSIVE_LOCK.Ptr(),
							PartialLoadingKeys: []string{
								TestDataObjectKey1,
								TestDataObjectKey2,
							},
							LockingKeys: []string{
								TestDataObjectKey1,
							},
						},
					},
				})
			}

			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				StateDecision: &iwfidl.StateDecision{
					NextStates: stms,
				},
			})
			return
		}
		if req.GetWorkflowStateId() == State2 {
			daInt := 0
			for _, da := range req.DataObjects {
				if da.GetKey() == TestDataObjectKey1 {
					value := da.GetValue()
					data := value.GetData()
					if data != "" {
						i, err := strconv.ParseInt(data, 10, 32)
						if err != nil {
							panic(err)
						}
						daInt = int(i)
					}
				}
			}
			daInt++
			context := req.GetContext()
			c.JSON(http.StatusOK, iwfidl.WorkflowStateDecideResponse{
				UpsertDataObjects: []iwfidl.KeyValue{
					{
						Key: iwfidl.PtrString(TestDataObjectKey1),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(fmt.Sprintf("%v", daInt)),
						},
					},
					{
						Key: iwfidl.PtrString(TestDataObjectKey2),
						Value: &iwfidl.EncodedObject{
							Encoding: iwfidl.PtrString("json"),
							Data:     iwfidl.PtrString(context.GetStateExecutionId()),
						},
					},
				},

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
	}

	c.JSON(http.StatusBadRequest, struct{}{})
}

func (h *handler) GetTestResult() (map[string]int64, map[string]interface{}) {
	return h.invokeHistory, h.invokeData
}
