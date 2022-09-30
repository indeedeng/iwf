# WorkflowResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompletedStateId** | **string** |  | 
**CompletedStateExecutionId** | **string** |  | 
**CompletedStateOutput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewWorkflowResult

`func NewWorkflowResult(completedStateId string, completedStateExecutionId string, ) *WorkflowResult`

NewWorkflowResult instantiates a new WorkflowResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowResultWithDefaults

`func NewWorkflowResultWithDefaults() *WorkflowResult`

NewWorkflowResultWithDefaults instantiates a new WorkflowResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompletedStateId

`func (o *WorkflowResult) GetCompletedStateId() string`

GetCompletedStateId returns the CompletedStateId field if non-nil, zero value otherwise.

### GetCompletedStateIdOk

`func (o *WorkflowResult) GetCompletedStateIdOk() (*string, bool)`

GetCompletedStateIdOk returns a tuple with the CompletedStateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateId

`func (o *WorkflowResult) SetCompletedStateId(v string)`

SetCompletedStateId sets CompletedStateId field to given value.


### GetCompletedStateExecutionId

`func (o *WorkflowResult) GetCompletedStateExecutionId() string`

GetCompletedStateExecutionId returns the CompletedStateExecutionId field if non-nil, zero value otherwise.

### GetCompletedStateExecutionIdOk

`func (o *WorkflowResult) GetCompletedStateExecutionIdOk() (*string, bool)`

GetCompletedStateExecutionIdOk returns a tuple with the CompletedStateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateExecutionId

`func (o *WorkflowResult) SetCompletedStateExecutionId(v string)`

SetCompletedStateExecutionId sets CompletedStateExecutionId field to given value.


### GetCompletedStateOutput

`func (o *WorkflowResult) GetCompletedStateOutput() EncodedObject`

GetCompletedStateOutput returns the CompletedStateOutput field if non-nil, zero value otherwise.

### GetCompletedStateOutputOk

`func (o *WorkflowResult) GetCompletedStateOutputOk() (*EncodedObject, bool)`

GetCompletedStateOutputOk returns a tuple with the CompletedStateOutput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateOutput

`func (o *WorkflowResult) SetCompletedStateOutput(v EncodedObject)`

SetCompletedStateOutput sets CompletedStateOutput field to given value.

### HasCompletedStateOutput

`func (o *WorkflowResult) HasCompletedStateOutput() bool`

HasCompletedStateOutput returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


