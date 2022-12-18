# StateMovement

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StateId** | **string** |  | 
**StateInput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**StateOptions** | Pointer to [**WorkflowStateOptions**](WorkflowStateOptions.md) |  | [optional] 

## Methods

### NewStateMovement

`func NewStateMovement(stateId string, ) *StateMovement`

NewStateMovement instantiates a new StateMovement object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStateMovementWithDefaults

`func NewStateMovementWithDefaults() *StateMovement`

NewStateMovementWithDefaults instantiates a new StateMovement object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStateId

`func (o *StateMovement) GetStateId() string`

GetStateId returns the StateId field if non-nil, zero value otherwise.

### GetStateIdOk

`func (o *StateMovement) GetStateIdOk() (*string, bool)`

GetStateIdOk returns a tuple with the StateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateId

`func (o *StateMovement) SetStateId(v string)`

SetStateId sets StateId field to given value.


### GetStateInput

`func (o *StateMovement) GetStateInput() EncodedObject`

GetStateInput returns the StateInput field if non-nil, zero value otherwise.

### GetStateInputOk

`func (o *StateMovement) GetStateInputOk() (*EncodedObject, bool)`

GetStateInputOk returns a tuple with the StateInput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateInput

`func (o *StateMovement) SetStateInput(v EncodedObject)`

SetStateInput sets StateInput field to given value.

### HasStateInput

`func (o *StateMovement) HasStateInput() bool`

HasStateInput returns a boolean if a field has been set.

### GetStateOptions

`func (o *StateMovement) GetStateOptions() WorkflowStateOptions`

GetStateOptions returns the StateOptions field if non-nil, zero value otherwise.

### GetStateOptionsOk

`func (o *StateMovement) GetStateOptionsOk() (*WorkflowStateOptions, bool)`

GetStateOptionsOk returns a tuple with the StateOptions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateOptions

`func (o *StateMovement) SetStateOptions(v WorkflowStateOptions)`

SetStateOptions sets StateOptions field to given value.

### HasStateOptions

`func (o *StateMovement) HasStateOptions() bool`

HasStateOptions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


