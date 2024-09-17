# WorkflowWaitForStateCompletionRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**StateExecutionId** | Pointer to **string** |  | [optional] 
**StateId** | Pointer to **string** |  | [optional] 
**WaitForKey** | Pointer to **string** |  | [optional] 
**WaitTimeSeconds** | Pointer to **int32** |  | [optional] 

## Methods

### NewWorkflowWaitForStateCompletionRequest

`func NewWorkflowWaitForStateCompletionRequest(workflowId string, ) *WorkflowWaitForStateCompletionRequest`

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

### HasStateExecutionId

`func (o *WorkflowWaitForStateCompletionRequest) HasStateExecutionId() bool`

HasStateExecutionId returns a boolean if a field has been set.

### GetStateId

`func (o *WorkflowWaitForStateCompletionRequest) GetStateId() string`

GetStateId returns the StateId field if non-nil, zero value otherwise.

### GetStateIdOk

`func (o *WorkflowWaitForStateCompletionRequest) GetStateIdOk() (*string, bool)`

GetStateIdOk returns a tuple with the StateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateId

`func (o *WorkflowWaitForStateCompletionRequest) SetStateId(v string)`

SetStateId sets StateId field to given value.

### HasStateId

`func (o *WorkflowWaitForStateCompletionRequest) HasStateId() bool`

HasStateId returns a boolean if a field has been set.

### GetWaitForKey

`func (o *WorkflowWaitForStateCompletionRequest) GetWaitForKey() string`

GetWaitForKey returns the WaitForKey field if non-nil, zero value otherwise.

### GetWaitForKeyOk

`func (o *WorkflowWaitForStateCompletionRequest) GetWaitForKeyOk() (*string, bool)`

GetWaitForKeyOk returns a tuple with the WaitForKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitForKey

`func (o *WorkflowWaitForStateCompletionRequest) SetWaitForKey(v string)`

SetWaitForKey sets WaitForKey field to given value.

### HasWaitForKey

`func (o *WorkflowWaitForStateCompletionRequest) HasWaitForKey() bool`

HasWaitForKey returns a boolean if a field has been set.

### GetWaitTimeSeconds

`func (o *WorkflowWaitForStateCompletionRequest) GetWaitTimeSeconds() int32`

GetWaitTimeSeconds returns the WaitTimeSeconds field if non-nil, zero value otherwise.

### GetWaitTimeSecondsOk

`func (o *WorkflowWaitForStateCompletionRequest) GetWaitTimeSecondsOk() (*int32, bool)`

GetWaitTimeSecondsOk returns a tuple with the WaitTimeSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitTimeSeconds

`func (o *WorkflowWaitForStateCompletionRequest) SetWaitTimeSeconds(v int32)`

SetWaitTimeSeconds sets WaitTimeSeconds field to given value.

### HasWaitTimeSeconds

`func (o *WorkflowWaitForStateCompletionRequest) HasWaitTimeSeconds() bool`

HasWaitTimeSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


