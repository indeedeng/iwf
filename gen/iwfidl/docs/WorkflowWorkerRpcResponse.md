# WorkflowWorkerRpcResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Output** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 
**StateDecision** | Pointer to [**StateDecision**](StateDecision.md) |  | [optional] 
**UpsertSearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 
**UpsertDataAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**RecordEvents** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**UpsertStateLocals** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**PublishToInterStateChannel** | Pointer to [**[]InterStateChannelPublishing**](InterStateChannelPublishing.md) |  | [optional] 

## Methods

### NewWorkflowWorkerRpcResponse

`func NewWorkflowWorkerRpcResponse() *WorkflowWorkerRpcResponse`

NewWorkflowWorkerRpcResponse instantiates a new WorkflowWorkerRpcResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowWorkerRpcResponseWithDefaults

`func NewWorkflowWorkerRpcResponseWithDefaults() *WorkflowWorkerRpcResponse`

NewWorkflowWorkerRpcResponseWithDefaults instantiates a new WorkflowWorkerRpcResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetOutput

`func (o *WorkflowWorkerRpcResponse) GetOutput() EncodedObject`

GetOutput returns the Output field if non-nil, zero value otherwise.

### GetOutputOk

`func (o *WorkflowWorkerRpcResponse) GetOutputOk() (*EncodedObject, bool)`

GetOutputOk returns a tuple with the Output field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutput

`func (o *WorkflowWorkerRpcResponse) SetOutput(v EncodedObject)`

SetOutput sets Output field to given value.

### HasOutput

`func (o *WorkflowWorkerRpcResponse) HasOutput() bool`

HasOutput returns a boolean if a field has been set.

### GetStateDecision

`func (o *WorkflowWorkerRpcResponse) GetStateDecision() StateDecision`

GetStateDecision returns the StateDecision field if non-nil, zero value otherwise.

### GetStateDecisionOk

`func (o *WorkflowWorkerRpcResponse) GetStateDecisionOk() (*StateDecision, bool)`

GetStateDecisionOk returns a tuple with the StateDecision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateDecision

`func (o *WorkflowWorkerRpcResponse) SetStateDecision(v StateDecision)`

SetStateDecision sets StateDecision field to given value.

### HasStateDecision

`func (o *WorkflowWorkerRpcResponse) HasStateDecision() bool`

HasStateDecision returns a boolean if a field has been set.

### GetUpsertSearchAttributes

`func (o *WorkflowWorkerRpcResponse) GetUpsertSearchAttributes() []SearchAttribute`

GetUpsertSearchAttributes returns the UpsertSearchAttributes field if non-nil, zero value otherwise.

### GetUpsertSearchAttributesOk

`func (o *WorkflowWorkerRpcResponse) GetUpsertSearchAttributesOk() (*[]SearchAttribute, bool)`

GetUpsertSearchAttributesOk returns a tuple with the UpsertSearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertSearchAttributes

`func (o *WorkflowWorkerRpcResponse) SetUpsertSearchAttributes(v []SearchAttribute)`

SetUpsertSearchAttributes sets UpsertSearchAttributes field to given value.

### HasUpsertSearchAttributes

`func (o *WorkflowWorkerRpcResponse) HasUpsertSearchAttributes() bool`

HasUpsertSearchAttributes returns a boolean if a field has been set.

### GetUpsertDataAttributes

`func (o *WorkflowWorkerRpcResponse) GetUpsertDataAttributes() []KeyValue`

GetUpsertDataAttributes returns the UpsertDataAttributes field if non-nil, zero value otherwise.

### GetUpsertDataAttributesOk

`func (o *WorkflowWorkerRpcResponse) GetUpsertDataAttributesOk() (*[]KeyValue, bool)`

GetUpsertDataAttributesOk returns a tuple with the UpsertDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertDataAttributes

`func (o *WorkflowWorkerRpcResponse) SetUpsertDataAttributes(v []KeyValue)`

SetUpsertDataAttributes sets UpsertDataAttributes field to given value.

### HasUpsertDataAttributes

`func (o *WorkflowWorkerRpcResponse) HasUpsertDataAttributes() bool`

HasUpsertDataAttributes returns a boolean if a field has been set.

### GetRecordEvents

`func (o *WorkflowWorkerRpcResponse) GetRecordEvents() []KeyValue`

GetRecordEvents returns the RecordEvents field if non-nil, zero value otherwise.

### GetRecordEventsOk

`func (o *WorkflowWorkerRpcResponse) GetRecordEventsOk() (*[]KeyValue, bool)`

GetRecordEventsOk returns a tuple with the RecordEvents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecordEvents

`func (o *WorkflowWorkerRpcResponse) SetRecordEvents(v []KeyValue)`

SetRecordEvents sets RecordEvents field to given value.

### HasRecordEvents

`func (o *WorkflowWorkerRpcResponse) HasRecordEvents() bool`

HasRecordEvents returns a boolean if a field has been set.

### GetUpsertStateLocals

`func (o *WorkflowWorkerRpcResponse) GetUpsertStateLocals() []KeyValue`

GetUpsertStateLocals returns the UpsertStateLocals field if non-nil, zero value otherwise.

### GetUpsertStateLocalsOk

`func (o *WorkflowWorkerRpcResponse) GetUpsertStateLocalsOk() (*[]KeyValue, bool)`

GetUpsertStateLocalsOk returns a tuple with the UpsertStateLocals field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpsertStateLocals

`func (o *WorkflowWorkerRpcResponse) SetUpsertStateLocals(v []KeyValue)`

SetUpsertStateLocals sets UpsertStateLocals field to given value.

### HasUpsertStateLocals

`func (o *WorkflowWorkerRpcResponse) HasUpsertStateLocals() bool`

HasUpsertStateLocals returns a boolean if a field has been set.

### GetPublishToInterStateChannel

`func (o *WorkflowWorkerRpcResponse) GetPublishToInterStateChannel() []InterStateChannelPublishing`

GetPublishToInterStateChannel returns the PublishToInterStateChannel field if non-nil, zero value otherwise.

### GetPublishToInterStateChannelOk

`func (o *WorkflowWorkerRpcResponse) GetPublishToInterStateChannelOk() (*[]InterStateChannelPublishing, bool)`

GetPublishToInterStateChannelOk returns a tuple with the PublishToInterStateChannel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublishToInterStateChannel

`func (o *WorkflowWorkerRpcResponse) SetPublishToInterStateChannel(v []InterStateChannelPublishing)`

SetPublishToInterStateChannel sets PublishToInterStateChannel field to given value.

### HasPublishToInterStateChannel

`func (o *WorkflowWorkerRpcResponse) HasPublishToInterStateChannel() bool`

HasPublishToInterStateChannel returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


