# IwfEvent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EventType** | [**EventType**](EventType.md) |  | 
**WorkflowType** | **string** |  | 
**WorkflowId** | **string** |  | 
**WorkflowRunId** | **string** |  | 
**StateId** | Pointer to **string** |  | [optional] 
**StateExecutionId** | Pointer to **string** |  | [optional] 
**RpcName** | Pointer to **string** |  | [optional] 
**StartTimestampInMs** | Pointer to **int64** |  | [optional] 
**EndTimestampInMs** | Pointer to **int64** |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 

## Methods

### NewIwfEvent

`func NewIwfEvent(eventType EventType, workflowType string, workflowId string, workflowRunId string, ) *IwfEvent`

NewIwfEvent instantiates a new IwfEvent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIwfEventWithDefaults

`func NewIwfEventWithDefaults() *IwfEvent`

NewIwfEventWithDefaults instantiates a new IwfEvent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEventType

`func (o *IwfEvent) GetEventType() EventType`

GetEventType returns the EventType field if non-nil, zero value otherwise.

### GetEventTypeOk

`func (o *IwfEvent) GetEventTypeOk() (*EventType, bool)`

GetEventTypeOk returns a tuple with the EventType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEventType

`func (o *IwfEvent) SetEventType(v EventType)`

SetEventType sets EventType field to given value.


### GetWorkflowType

`func (o *IwfEvent) GetWorkflowType() string`

GetWorkflowType returns the WorkflowType field if non-nil, zero value otherwise.

### GetWorkflowTypeOk

`func (o *IwfEvent) GetWorkflowTypeOk() (*string, bool)`

GetWorkflowTypeOk returns a tuple with the WorkflowType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowType

`func (o *IwfEvent) SetWorkflowType(v string)`

SetWorkflowType sets WorkflowType field to given value.


### GetWorkflowId

`func (o *IwfEvent) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *IwfEvent) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *IwfEvent) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *IwfEvent) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *IwfEvent) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *IwfEvent) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.


### GetStateId

`func (o *IwfEvent) GetStateId() string`

GetStateId returns the StateId field if non-nil, zero value otherwise.

### GetStateIdOk

`func (o *IwfEvent) GetStateIdOk() (*string, bool)`

GetStateIdOk returns a tuple with the StateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateId

`func (o *IwfEvent) SetStateId(v string)`

SetStateId sets StateId field to given value.

### HasStateId

`func (o *IwfEvent) HasStateId() bool`

HasStateId returns a boolean if a field has been set.

### GetStateExecutionId

`func (o *IwfEvent) GetStateExecutionId() string`

GetStateExecutionId returns the StateExecutionId field if non-nil, zero value otherwise.

### GetStateExecutionIdOk

`func (o *IwfEvent) GetStateExecutionIdOk() (*string, bool)`

GetStateExecutionIdOk returns a tuple with the StateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateExecutionId

`func (o *IwfEvent) SetStateExecutionId(v string)`

SetStateExecutionId sets StateExecutionId field to given value.

### HasStateExecutionId

`func (o *IwfEvent) HasStateExecutionId() bool`

HasStateExecutionId returns a boolean if a field has been set.

### GetRpcName

`func (o *IwfEvent) GetRpcName() string`

GetRpcName returns the RpcName field if non-nil, zero value otherwise.

### GetRpcNameOk

`func (o *IwfEvent) GetRpcNameOk() (*string, bool)`

GetRpcNameOk returns a tuple with the RpcName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRpcName

`func (o *IwfEvent) SetRpcName(v string)`

SetRpcName sets RpcName field to given value.

### HasRpcName

`func (o *IwfEvent) HasRpcName() bool`

HasRpcName returns a boolean if a field has been set.

### GetStartTimestampInMs

`func (o *IwfEvent) GetStartTimestampInMs() int64`

GetStartTimestampInMs returns the StartTimestampInMs field if non-nil, zero value otherwise.

### GetStartTimestampInMsOk

`func (o *IwfEvent) GetStartTimestampInMsOk() (*int64, bool)`

GetStartTimestampInMsOk returns a tuple with the StartTimestampInMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartTimestampInMs

`func (o *IwfEvent) SetStartTimestampInMs(v int64)`

SetStartTimestampInMs sets StartTimestampInMs field to given value.

### HasStartTimestampInMs

`func (o *IwfEvent) HasStartTimestampInMs() bool`

HasStartTimestampInMs returns a boolean if a field has been set.

### GetEndTimestampInMs

`func (o *IwfEvent) GetEndTimestampInMs() int64`

GetEndTimestampInMs returns the EndTimestampInMs field if non-nil, zero value otherwise.

### GetEndTimestampInMsOk

`func (o *IwfEvent) GetEndTimestampInMsOk() (*int64, bool)`

GetEndTimestampInMsOk returns a tuple with the EndTimestampInMs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndTimestampInMs

`func (o *IwfEvent) SetEndTimestampInMs(v int64)`

SetEndTimestampInMs sets EndTimestampInMs field to given value.

### HasEndTimestampInMs

`func (o *IwfEvent) HasEndTimestampInMs() bool`

HasEndTimestampInMs returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *IwfEvent) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *IwfEvent) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *IwfEvent) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *IwfEvent) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


