# WorkflowStopRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**Reason** | Pointer to **string** |  | [optional] 
**StopType** | Pointer to **string** |  | [optional] 

## Methods

### NewWorkflowStopRequest

`func NewWorkflowStopRequest(workflowId string, ) *WorkflowStopRequest`

NewWorkflowStopRequest instantiates a new WorkflowStopRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStopRequestWithDefaults

`func NewWorkflowStopRequestWithDefaults() *WorkflowStopRequest`

NewWorkflowStopRequestWithDefaults instantiates a new WorkflowStopRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowStopRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowStopRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowStopRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowStopRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowStopRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowStopRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowStopRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetReason

`func (o *WorkflowStopRequest) GetReason() string`

GetReason returns the Reason field if non-nil, zero value otherwise.

### GetReasonOk

`func (o *WorkflowStopRequest) GetReasonOk() (*string, bool)`

GetReasonOk returns a tuple with the Reason field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReason

`func (o *WorkflowStopRequest) SetReason(v string)`

SetReason sets Reason field to given value.

### HasReason

`func (o *WorkflowStopRequest) HasReason() bool`

HasReason returns a boolean if a field has been set.

### GetStopType

`func (o *WorkflowStopRequest) GetStopType() string`

GetStopType returns the StopType field if non-nil, zero value otherwise.

### GetStopTypeOk

`func (o *WorkflowStopRequest) GetStopTypeOk() (*string, bool)`

GetStopTypeOk returns a tuple with the StopType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStopType

`func (o *WorkflowStopRequest) SetStopType(v string)`

SetStopType sets StopType field to given value.

### HasStopType

`func (o *WorkflowStopRequest) HasStopType() bool`

HasStopType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


