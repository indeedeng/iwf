# WorkflowWaitForStateCompletionRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**StateExecutionId** | **string** |  | 

## Methods

### NewWorkflowWaitForStateCompletionRequest

`func NewWorkflowWaitForStateCompletionRequest(workflowId string, stateExecutionId string, ) *WorkflowWaitForStateCompletionRequest`

NewWorkflowWaitForStateCompletionRequest instantiates a new WorkflowWaitForStateCompletionRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowWaitForStateCompletionRequestWithDefaults

`func NewWorkflowWaitForStateCompletionRequestWithDefaults() *WorkflowWaitForStateCompletionRequest`

NewWorkflowWaitForStateCompletionRequestWithDefaults instantiates a new WorkflowWaitForStateCompletionRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowWaitForStateCompletionRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowWaitForStateCompletionRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowWaitForStateCompletionRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetStateExecutionId

`func (o *WorkflowWaitForStateCompletionRequest) GetStateExecutionId() string`

GetStateExecutionId returns the StateExecutionId field if non-nil, zero value otherwise.

### GetStateExecutionIdOk

`func (o *WorkflowWaitForStateCompletionRequest) GetStateExecutionIdOk() (*string, bool)`

GetStateExecutionIdOk returns a tuple with the StateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateExecutionId

`func (o *WorkflowWaitForStateCompletionRequest) SetStateExecutionId(v string)`

SetStateExecutionId sets StateExecutionId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


