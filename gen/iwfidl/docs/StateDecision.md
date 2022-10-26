# StateDecision

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NextStates** | Pointer to [**[]StateMovement**](StateMovement.md) |  | [optional] 
**UpsertSearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**UpsertQueryAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**PublishToInterStateChannel** | Pointer to [**[]InterStateChannelPublishing**](InterStateChannelPublishing.md) |  | [optional] 

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

### GetUpsertSearchAttributes

`func (o *StateDecision) GetUpsertSearchAttributes() []SearchAttribute`

GetUpsertSearchAttributes returns the UpsertSearchAttributes field if non-nil, zero value otherwise.

### GetUpsertSearchAttributesOk

`func (o *StateDecision) GetUpsertSearchAttributesOk() (*[]SearchAttribute, bool)`

GetUpsertSearchAttributesOk returns a tuple with the UpsertSearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertSearchAttributes

`func (o *StateDecision) SetUpsertSearchAttributes(v []SearchAttribute)`

SetUpsertSearchAttributes sets UpsertSearchAttributes field to given value.

### HasUpsertSearchAttributes

`func (o *StateDecision) HasUpsertSearchAttributes() bool`

HasUpsertSearchAttributes returns a boolean if a field has been set.

### GetUpsertQueryAttributes

`func (o *StateDecision) GetUpsertQueryAttributes() []KeyValue`

GetUpsertQueryAttributes returns the UpsertQueryAttributes field if non-nil, zero value otherwise.

### GetUpsertQueryAttributesOk

`func (o *StateDecision) GetUpsertQueryAttributesOk() (*[]KeyValue, bool)`

GetUpsertQueryAttributesOk returns a tuple with the UpsertQueryAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertQueryAttributes

`func (o *StateDecision) SetUpsertQueryAttributes(v []KeyValue)`

SetUpsertQueryAttributes sets UpsertQueryAttributes field to given value.

### HasUpsertQueryAttributes

`func (o *StateDecision) HasUpsertQueryAttributes() bool`

HasUpsertQueryAttributes returns a boolean if a field has been set.

### GetPublishToInterStateChannel

`func (o *StateDecision) GetPublishToInterStateChannel() []InterStateChannelPublishing`

GetPublishToInterStateChannel returns the PublishToInterStateChannel field if non-nil, zero value otherwise.

### GetPublishToInterStateChannelOk

`func (o *StateDecision) GetPublishToInterStateChannelOk() (*[]InterStateChannelPublishing, bool)`

GetPublishToInterStateChannelOk returns a tuple with the PublishToInterStateChannel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishToInterStateChannel

`func (o *StateDecision) SetPublishToInterStateChannel(v []InterStateChannelPublishing)`

SetPublishToInterStateChannel sets PublishToInterStateChannel field to given value.

### HasPublishToInterStateChannel

`func (o *StateDecision) HasPublishToInterStateChannel() bool`

HasPublishToInterStateChannel returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


