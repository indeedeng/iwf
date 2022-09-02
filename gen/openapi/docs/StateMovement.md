# StateMovement

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StateId** | Pointer to **string** |  | [optional] 
**NextStateInput** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewStateMovement

`func NewStateMovement() *StateMovement`

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

### HasStateId

`func (o *StateMovement) HasStateId() bool`

HasStateId returns a boolean if a field has been set.

### GetNextStateInput

`func (o *StateMovement) GetNextStateInput() EncodedObject`

GetNextStateInput returns the NextStateInput field if non-nil, zero value otherwise.

### GetNextStateInputOk

`func (o *StateMovement) GetNextStateInputOk() (*EncodedObject, bool)`

GetNextStateInputOk returns a tuple with the NextStateInput field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextStateInput

`func (o *StateMovement) SetNextStateInput(v EncodedObject)`

SetNextStateInput sets NextStateInput field to given value.

### HasNextStateInput

`func (o *StateMovement) HasNextStateInput() bool`

HasNextStateInput returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


