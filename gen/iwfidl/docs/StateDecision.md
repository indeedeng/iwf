# StateDecision

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NextStates** | Pointer to [**[]StateMovement**](StateMovement.md) |  | [optional] 

## Methods

### NewStateDecision

`func NewStateDecision() *StateDecision`

NewStateDecision instantiates a new StateDecision object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStateDecisionWithDefaults

`func NewStateDecisionWithDefaults() *StateDecision`

NewStateDecisionWithDefaults instantiates a new StateDecision object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNextStates

`func (o *StateDecision) GetNextStates() []StateMovement`

GetNextStates returns the NextStates field if non-nil, zero value otherwise.

### GetNextStatesOk

`func (o *StateDecision) GetNextStatesOk() (*[]StateMovement, bool)`

GetNextStatesOk returns a tuple with the NextStates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextStates

`func (o *StateDecision) SetNextStates(v []StateMovement)`

SetNextStates sets NextStates field to given value.

### HasNextStates

`func (o *StateDecision) HasNextStates() bool`

HasNextStates returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


