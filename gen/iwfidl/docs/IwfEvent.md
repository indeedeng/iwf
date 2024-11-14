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
**StartTimestamp** | Pointer to **int64** |  | [optional] 
**EndTimestamp** | Pointer to **int64** |  | [optional] 

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

### GetStartTimestamp

`func (o *IwfEvent) GetStartTimestamp() int64`

GetStartTimestamp returns the StartTimestamp field if non-nil, zero value otherwise.

### GetStartTimestampOk

`func (o *IwfEvent) GetStartTimestampOk() (*int64, bool)`

GetStartTimestampOk returns a tuple with the StartTimestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartTimestamp

`func (o *IwfEvent) SetStartTimestamp(v int64)`

SetStartTimestamp sets StartTimestamp field to given value.

### HasStartTimestamp

`func (o *IwfEvent) HasStartTimestamp() bool`

HasStartTimestamp returns a boolean if a field has been set.

### GetEndTimestamp

`func (o *IwfEvent) GetEndTimestamp() int64`

GetEndTimestamp returns the EndTimestamp field if non-nil, zero value otherwise.

### GetEndTimestampOk

`func (o *IwfEvent) GetEndTimestampOk() (*int64, bool)`

GetEndTimestampOk returns a tuple with the EndTimestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndTimestamp

`func (o *IwfEvent) SetEndTimestamp(v int64)`

SetEndTimestamp sets EndTimestamp field to given value.

### HasEndTimestamp

`func (o *IwfEvent) HasEndTimestamp() bool`

HasEndTimestamp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


