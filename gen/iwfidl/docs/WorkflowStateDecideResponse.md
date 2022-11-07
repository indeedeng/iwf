# WorkflowStateDecideResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StateDecision** | Pointer to [**StateDecision**](StateDecision.md) |  | [optional] 
**UpsertSearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**UpsertQueryAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**RecordEvents** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**UpsertStateLocalAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**PublishToInterStateChannel** | Pointer to [**[]InterStateChannelPublishing**](InterStateChannelPublishing.md) |  | [optional] 

## Methods

### NewWorkflowStateDecideResponse

`func NewWorkflowStateDecideResponse() *WorkflowStateDecideResponse`

NewWorkflowStateDecideResponse instantiates a new WorkflowStateDecideResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStateDecideResponseWithDefaults

`func NewWorkflowStateDecideResponseWithDefaults() *WorkflowStateDecideResponse`

NewWorkflowStateDecideResponseWithDefaults instantiates a new WorkflowStateDecideResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStateDecision

`func (o *WorkflowStateDecideResponse) GetStateDecision() StateDecision`

GetStateDecision returns the StateDecision field if non-nil, zero value otherwise.

### GetStateDecisionOk

`func (o *WorkflowStateDecideResponse) GetStateDecisionOk() (*StateDecision, bool)`

GetStateDecisionOk returns a tuple with the StateDecision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateDecision

`func (o *WorkflowStateDecideResponse) SetStateDecision(v StateDecision)`

SetStateDecision sets StateDecision field to given value.

### HasStateDecision

`func (o *WorkflowStateDecideResponse) HasStateDecision() bool`

HasStateDecision returns a boolean if a field has been set.

### GetUpsertSearchAttributes

`func (o *WorkflowStateDecideResponse) GetUpsertSearchAttributes() []SearchAttribute`

GetUpsertSearchAttributes returns the UpsertSearchAttributes field if non-nil, zero value otherwise.

### GetUpsertSearchAttributesOk

`func (o *WorkflowStateDecideResponse) GetUpsertSearchAttributesOk() (*[]SearchAttribute, bool)`

GetUpsertSearchAttributesOk returns a tuple with the UpsertSearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertSearchAttributes

`func (o *WorkflowStateDecideResponse) SetUpsertSearchAttributes(v []SearchAttribute)`

SetUpsertSearchAttributes sets UpsertSearchAttributes field to given value.

### HasUpsertSearchAttributes

`func (o *WorkflowStateDecideResponse) HasUpsertSearchAttributes() bool`

HasUpsertSearchAttributes returns a boolean if a field has been set.

### GetUpsertQueryAttributes

`func (o *WorkflowStateDecideResponse) GetUpsertQueryAttributes() []KeyValue`

GetUpsertQueryAttributes returns the UpsertQueryAttributes field if non-nil, zero value otherwise.

### GetUpsertQueryAttributesOk

`func (o *WorkflowStateDecideResponse) GetUpsertQueryAttributesOk() (*[]KeyValue, bool)`

GetUpsertQueryAttributesOk returns a tuple with the UpsertQueryAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertQueryAttributes

`func (o *WorkflowStateDecideResponse) SetUpsertQueryAttributes(v []KeyValue)`

SetUpsertQueryAttributes sets UpsertQueryAttributes field to given value.

### HasUpsertQueryAttributes

`func (o *WorkflowStateDecideResponse) HasUpsertQueryAttributes() bool`

HasUpsertQueryAttributes returns a boolean if a field has been set.

### GetRecordEvents

`func (o *WorkflowStateDecideResponse) GetRecordEvents() []KeyValue`

GetRecordEvents returns the RecordEvents field if non-nil, zero value otherwise.

### GetRecordEventsOk

`func (o *WorkflowStateDecideResponse) GetRecordEventsOk() (*[]KeyValue, bool)`

GetRecordEventsOk returns a tuple with the RecordEvents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordEvents

`func (o *WorkflowStateDecideResponse) SetRecordEvents(v []KeyValue)`

SetRecordEvents sets RecordEvents field to given value.

### HasRecordEvents

`func (o *WorkflowStateDecideResponse) HasRecordEvents() bool`

HasRecordEvents returns a boolean if a field has been set.

### GetUpsertStateLocalAttributes

`func (o *WorkflowStateDecideResponse) GetUpsertStateLocalAttributes() []KeyValue`

GetUpsertStateLocalAttributes returns the UpsertStateLocalAttributes field if non-nil, zero value otherwise.

### GetUpsertStateLocalAttributesOk

`func (o *WorkflowStateDecideResponse) GetUpsertStateLocalAttributesOk() (*[]KeyValue, bool)`

GetUpsertStateLocalAttributesOk returns a tuple with the UpsertStateLocalAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertStateLocalAttributes

`func (o *WorkflowStateDecideResponse) SetUpsertStateLocalAttributes(v []KeyValue)`

SetUpsertStateLocalAttributes sets UpsertStateLocalAttributes field to given value.

### HasUpsertStateLocalAttributes

`func (o *WorkflowStateDecideResponse) HasUpsertStateLocalAttributes() bool`

HasUpsertStateLocalAttributes returns a boolean if a field has been set.

### GetPublishToInterStateChannel

`func (o *WorkflowStateDecideResponse) GetPublishToInterStateChannel() []InterStateChannelPublishing`

GetPublishToInterStateChannel returns the PublishToInterStateChannel field if non-nil, zero value otherwise.

### GetPublishToInterStateChannelOk

`func (o *WorkflowStateDecideResponse) GetPublishToInterStateChannelOk() (*[]InterStateChannelPublishing, bool)`

GetPublishToInterStateChannelOk returns a tuple with the PublishToInterStateChannel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishToInterStateChannel

`func (o *WorkflowStateDecideResponse) SetPublishToInterStateChannel(v []InterStateChannelPublishing)`

SetPublishToInterStateChannel sets PublishToInterStateChannel field to given value.

### HasPublishToInterStateChannel

`func (o *WorkflowStateDecideResponse) HasPublishToInterStateChannel() bool`

HasPublishToInterStateChannel returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


