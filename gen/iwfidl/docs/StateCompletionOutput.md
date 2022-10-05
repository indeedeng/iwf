# StateCompletionOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CompletedStateId** | **string** |  | 
**CompletedStateExecutionId** | **string** |  | 
**CompletedStateOutput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewStateCompletionOutput

`func NewStateCompletionOutput(completedStateId string, completedStateExecutionId string, ) *StateCompletionOutput`

NewStateCompletionOutput instantiates a new StateCompletionOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStateCompletionOutputWithDefaults

`func NewStateCompletionOutputWithDefaults() *StateCompletionOutput`

NewStateCompletionOutputWithDefaults instantiates a new StateCompletionOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCompletedStateId

`func (o *StateCompletionOutput) GetCompletedStateId() string`

GetCompletedStateId returns the CompletedStateId field if non-nil, zero value otherwise.

### GetCompletedStateIdOk

`func (o *StateCompletionOutput) GetCompletedStateIdOk() (*string, bool)`

GetCompletedStateIdOk returns a tuple with the CompletedStateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateId

`func (o *StateCompletionOutput) SetCompletedStateId(v string)`

SetCompletedStateId sets CompletedStateId field to given value.


### GetCompletedStateExecutionId

`func (o *StateCompletionOutput) GetCompletedStateExecutionId() string`

GetCompletedStateExecutionId returns the CompletedStateExecutionId field if non-nil, zero value otherwise.

### GetCompletedStateExecutionIdOk

`func (o *StateCompletionOutput) GetCompletedStateExecutionIdOk() (*string, bool)`

GetCompletedStateExecutionIdOk returns a tuple with the CompletedStateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateExecutionId

`func (o *StateCompletionOutput) SetCompletedStateExecutionId(v string)`

SetCompletedStateExecutionId sets CompletedStateExecutionId field to given value.


### GetCompletedStateOutput

`func (o *StateCompletionOutput) GetCompletedStateOutput() EncodedObject`

GetCompletedStateOutput returns the CompletedStateOutput field if non-nil, zero value otherwise.

### GetCompletedStateOutputOk

`func (o *StateCompletionOutput) GetCompletedStateOutputOk() (*EncodedObject, bool)`

GetCompletedStateOutputOk returns a tuple with the CompletedStateOutput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCompletedStateOutput

`func (o *StateCompletionOutput) SetCompletedStateOutput(v EncodedObject)`

SetCompletedStateOutput sets CompletedStateOutput field to given value.

### HasCompletedStateOutput

`func (o *StateCompletionOutput) HasCompletedStateOutput() bool`

HasCompletedStateOutput returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


