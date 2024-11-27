# WorkflowResetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**ResetType** | [**WorkflowResetType**](WorkflowResetType.md) |  | 
**HistoryEventId** | Pointer to **int32** |  | [optional] 
**Reason** | Pointer to **string** |  | [optional] 
**HistoryEventTime** | Pointer to **string** |  | [optional] 
**StateId** | Pointer to **string** |  | [optional] 
**StateExecutionId** | Pointer to **string** |  | [optional] 
**SkipSignalReapply** | Pointer to **bool** |  | [optional] 
**SkipUpdateReapply** | Pointer to **bool** |  | [optional] 

## Methods

### NewWorkflowResetRequest

`func NewWorkflowResetRequest(workflowId string, resetType WorkflowResetType, ) *WorkflowResetRequest`

NewWorkflowResetRequest instantiates a new WorkflowResetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowResetRequestWithDefaults

`func NewWorkflowResetRequestWithDefaults() *WorkflowResetRequest`

NewWorkflowResetRequestWithDefaults instantiates a new WorkflowResetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowResetRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowResetRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowResetRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowResetRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowResetRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowResetRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowResetRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetResetType

`func (o *WorkflowResetRequest) GetResetType() WorkflowResetType`

GetResetType returns the ResetType field if non-nil, zero value otherwise.

### GetResetTypeOk

`func (o *WorkflowResetRequest) GetResetTypeOk() (*WorkflowResetType, bool)`

GetResetTypeOk returns a tuple with the ResetType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResetType

`func (o *WorkflowResetRequest) SetResetType(v WorkflowResetType)`

SetResetType sets ResetType field to given value.


### GetHistoryEventId

`func (o *WorkflowResetRequest) GetHistoryEventId() int32`

GetHistoryEventId returns the HistoryEventId field if non-nil, zero value otherwise.

### GetHistoryEventIdOk

`func (o *WorkflowResetRequest) GetHistoryEventIdOk() (*int32, bool)`

GetHistoryEventIdOk returns a tuple with the HistoryEventId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHistoryEventId

`func (o *WorkflowResetRequest) SetHistoryEventId(v int32)`

SetHistoryEventId sets HistoryEventId field to given value.

### HasHistoryEventId

`func (o *WorkflowResetRequest) HasHistoryEventId() bool`

HasHistoryEventId returns a boolean if a field has been set.

### GetReason

`func (o *WorkflowResetRequest) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *WorkflowResetRequest) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *WorkflowResetRequest) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *WorkflowResetRequest) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetHistoryEventTime

`func (o *WorkflowResetRequest) GetHistoryEventTime() string`

GetHistoryEventTime returns the HistoryEventTime field if non-nil, zero value otherwise.

### GetHistoryEventTimeOk

`func (o *WorkflowResetRequest) GetHistoryEventTimeOk() (*string, bool)`

GetHistoryEventTimeOk returns a tuple with the HistoryEventTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHistoryEventTime

`func (o *WorkflowResetRequest) SetHistoryEventTime(v string)`

SetHistoryEventTime sets HistoryEventTime field to given value.

### HasHistoryEventTime

`func (o *WorkflowResetRequest) HasHistoryEventTime() bool`

HasHistoryEventTime returns a boolean if a field has been set.

### GetStateId

`func (o *WorkflowResetRequest) GetStateId() string`

GetStateId returns the StateId field if non-nil, zero value otherwise.

### GetStateIdOk

`func (o *WorkflowResetRequest) GetStateIdOk() (*string, bool)`

GetStateIdOk returns a tuple with the StateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateId

`func (o *WorkflowResetRequest) SetStateId(v string)`

SetStateId sets StateId field to given value.

### HasStateId

`func (o *WorkflowResetRequest) HasStateId() bool`

HasStateId returns a boolean if a field has been set.

### GetStateExecutionId

`func (o *WorkflowResetRequest) GetStateExecutionId() string`

GetStateExecutionId returns the StateExecutionId field if non-nil, zero value otherwise.

### GetStateExecutionIdOk

`func (o *WorkflowResetRequest) GetStateExecutionIdOk() (*string, bool)`

GetStateExecutionIdOk returns a tuple with the StateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateExecutionId

`func (o *WorkflowResetRequest) SetStateExecutionId(v string)`

SetStateExecutionId sets StateExecutionId field to given value.

### HasStateExecutionId

`func (o *WorkflowResetRequest) HasStateExecutionId() bool`

HasStateExecutionId returns a boolean if a field has been set.

### GetSkipSignalReapply

`func (o *WorkflowResetRequest) GetSkipSignalReapply() bool`

GetSkipSignalReapply returns the SkipSignalReapply field if non-nil, zero value otherwise.

### GetSkipSignalReapplyOk

`func (o *WorkflowResetRequest) GetSkipSignalReapplyOk() (*bool, bool)`

GetSkipSignalReapplyOk returns a tuple with the SkipSignalReapply field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipSignalReapply

`func (o *WorkflowResetRequest) SetSkipSignalReapply(v bool)`

SetSkipSignalReapply sets SkipSignalReapply field to given value.

### HasSkipSignalReapply

`func (o *WorkflowResetRequest) HasSkipSignalReapply() bool`

HasSkipSignalReapply returns a boolean if a field has been set.

### GetSkipUpdateReapply

`func (o *WorkflowResetRequest) GetSkipUpdateReapply() bool`

GetSkipUpdateReapply returns the SkipUpdateReapply field if non-nil, zero value otherwise.

### GetSkipUpdateReapplyOk

`func (o *WorkflowResetRequest) GetSkipUpdateReapplyOk() (*bool, bool)`

GetSkipUpdateReapplyOk returns a tuple with the SkipUpdateReapply field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipUpdateReapply

`func (o *WorkflowResetRequest) SetSkipUpdateReapply(v bool)`

SetSkipUpdateReapply sets SkipUpdateReapply field to given value.

### HasSkipUpdateReapply

`func (o *WorkflowResetRequest) HasSkipUpdateReapply() bool`

HasSkipUpdateReapply returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


