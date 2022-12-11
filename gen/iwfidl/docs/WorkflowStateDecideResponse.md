# WorkflowStateDecideResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StateDecision** | Pointer to [**StateDecision**](StateDecision.md) |  | [optional] 
**UpsertSearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**UpsertDataObjects** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**RecordEvents** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**UpsertStateLocals** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
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

### GetUpsertDataObjects

`func (o *WorkflowStateDecideResponse) GetUpsertDataObjects() []KeyValue`

GetUpsertDataObjects returns the UpsertDataObjects field if non-nil, zero value otherwise.

### GetUpsertDataObjectsOk

`func (o *WorkflowStateDecideResponse) GetUpsertDataObjectsOk() (*[]KeyValue, bool)`

GetUpsertDataObjectsOk returns a tuple with the UpsertDataObjects field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertDataObjects

`func (o *WorkflowStateDecideResponse) SetUpsertDataObjects(v []KeyValue)`

SetUpsertDataObjects sets UpsertDataObjects field to given value.

### HasUpsertDataObjects

`func (o *WorkflowStateDecideResponse) HasUpsertDataObjects() bool`

HasUpsertDataObjects returns a boolean if a field has been set.

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

### GetUpsertStateLocals

`func (o *WorkflowStateDecideResponse) GetUpsertStateLocals() []KeyValue`

GetUpsertStateLocals returns the UpsertStateLocals field if non-nil, zero value otherwise.

### GetUpsertStateLocalsOk

`func (o *WorkflowStateDecideResponse) GetUpsertStateLocalsOk() (*[]KeyValue, bool)`

GetUpsertStateLocalsOk returns a tuple with the UpsertStateLocals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertStateLocals

`func (o *WorkflowStateDecideResponse) SetUpsertStateLocals(v []KeyValue)`

SetUpsertStateLocals sets UpsertStateLocals field to given value.

### HasUpsertStateLocals

`func (o *WorkflowStateDecideResponse) HasUpsertStateLocals() bool`

HasUpsertStateLocals returns a boolean if a field has been set.

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


