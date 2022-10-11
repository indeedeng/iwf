# WorkflowResetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**ResetType** | **string** |  | 
**HistoryEventId** | Pointer to **int32** |  | [optional] 
**Reason** | Pointer to **string** |  | [optional] 
**DecisionOffset** | Pointer to **int32** |  | [optional] 
**ResetBadBinaryChecksum** | Pointer to **string** |  | [optional] 
**EarliestTime** | Pointer to **string** |  | [optional] 
**SkipSignalReapply** | Pointer to **bool** |  | [optional] 

## Methods

### NewWorkflowResetRequest

`func NewWorkflowResetRequest(workflowId string, resetType string, ) *WorkflowResetRequest`

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

`func (o *WorkflowResetRequest) GetResetType() string`

GetResetType returns the ResetType field if non-nil, zero value otherwise.

### GetResetTypeOk

`func (o *WorkflowResetRequest) GetResetTypeOk() (*string, bool)`

GetResetTypeOk returns a tuple with the ResetType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResetType

`func (o *WorkflowResetRequest) SetResetType(v string)`

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

### GetDecisionOffset

`func (o *WorkflowResetRequest) GetDecisionOffset() int32`

GetDecisionOffset returns the DecisionOffset field if non-nil, zero value otherwise.

### GetDecisionOffsetOk

`func (o *WorkflowResetRequest) GetDecisionOffsetOk() (*int32, bool)`

GetDecisionOffsetOk returns a tuple with the DecisionOffset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecisionOffset

`func (o *WorkflowResetRequest) SetDecisionOffset(v int32)`

SetDecisionOffset sets DecisionOffset field to given value.

### HasDecisionOffset

`func (o *WorkflowResetRequest) HasDecisionOffset() bool`

HasDecisionOffset returns a boolean if a field has been set.

### GetResetBadBinaryChecksum

`func (o *WorkflowResetRequest) GetResetBadBinaryChecksum() string`

GetResetBadBinaryChecksum returns the ResetBadBinaryChecksum field if non-nil, zero value otherwise.

### GetResetBadBinaryChecksumOk

`func (o *WorkflowResetRequest) GetResetBadBinaryChecksumOk() (*string, bool)`

GetResetBadBinaryChecksumOk returns a tuple with the ResetBadBinaryChecksum field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResetBadBinaryChecksum

`func (o *WorkflowResetRequest) SetResetBadBinaryChecksum(v string)`

SetResetBadBinaryChecksum sets ResetBadBinaryChecksum field to given value.

### HasResetBadBinaryChecksum

`func (o *WorkflowResetRequest) HasResetBadBinaryChecksum() bool`

HasResetBadBinaryChecksum returns a boolean if a field has been set.

### GetEarliestTime

`func (o *WorkflowResetRequest) GetEarliestTime() string`

GetEarliestTime returns the EarliestTime field if non-nil, zero value otherwise.

### GetEarliestTimeOk

`func (o *WorkflowResetRequest) GetEarliestTimeOk() (*string, bool)`

GetEarliestTimeOk returns a tuple with the EarliestTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEarliestTime

`func (o *WorkflowResetRequest) SetEarliestTime(v string)`

SetEarliestTime sets EarliestTime field to given value.

### HasEarliestTime

`func (o *WorkflowResetRequest) HasEarliestTime() bool`

HasEarliestTime returns a boolean if a field has been set.

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


