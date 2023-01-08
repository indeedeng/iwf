# WorkflowSkipTimerRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**WorkflowStateExecutionId** | **string** |  | 
**TimerCommandId** | Pointer to **string** |  | [optional] 
**TimerCommandIndex** | Pointer to **int32** |  | [optional] 

## Methods

### NewWorkflowSkipTimerRequest

`func NewWorkflowSkipTimerRequest(workflowId string, workflowStateExecutionId string, ) *WorkflowSkipTimerRequest`

NewWorkflowSkipTimerRequest instantiates a new WorkflowSkipTimerRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSkipTimerRequestWithDefaults

`func NewWorkflowSkipTimerRequestWithDefaults() *WorkflowSkipTimerRequest`

NewWorkflowSkipTimerRequestWithDefaults instantiates a new WorkflowSkipTimerRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowSkipTimerRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowSkipTimerRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowSkipTimerRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowSkipTimerRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowSkipTimerRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowSkipTimerRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowSkipTimerRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetWorkflowStateExecutionId

`func (o *WorkflowSkipTimerRequest) GetWorkflowStateExecutionId() string`

GetWorkflowStateExecutionId returns the WorkflowStateExecutionId field if non-nil, zero value otherwise.

### GetWorkflowStateExecutionIdOk

`func (o *WorkflowSkipTimerRequest) GetWorkflowStateExecutionIdOk() (*string, bool)`

GetWorkflowStateExecutionIdOk returns a tuple with the WorkflowStateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStateExecutionId

`func (o *WorkflowSkipTimerRequest) SetWorkflowStateExecutionId(v string)`

SetWorkflowStateExecutionId sets WorkflowStateExecutionId field to given value.


### GetTimerCommandId

`func (o *WorkflowSkipTimerRequest) GetTimerCommandId() string`

GetTimerCommandId returns the TimerCommandId field if non-nil, zero value otherwise.

### GetTimerCommandIdOk

`func (o *WorkflowSkipTimerRequest) GetTimerCommandIdOk() (*string, bool)`

GetTimerCommandIdOk returns a tuple with the TimerCommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimerCommandId

`func (o *WorkflowSkipTimerRequest) SetTimerCommandId(v string)`

SetTimerCommandId sets TimerCommandId field to given value.

### HasTimerCommandId

`func (o *WorkflowSkipTimerRequest) HasTimerCommandId() bool`

HasTimerCommandId returns a boolean if a field has been set.

### GetTimerCommandIndex

`func (o *WorkflowSkipTimerRequest) GetTimerCommandIndex() int32`

GetTimerCommandIndex returns the TimerCommandIndex field if non-nil, zero value otherwise.

### GetTimerCommandIndexOk

`func (o *WorkflowSkipTimerRequest) GetTimerCommandIndexOk() (*int32, bool)`

GetTimerCommandIndexOk returns a tuple with the TimerCommandIndex field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimerCommandIndex

`func (o *WorkflowSkipTimerRequest) SetTimerCommandIndex(v int32)`

SetTimerCommandIndex sets TimerCommandIndex field to given value.

### HasTimerCommandIndex

`func (o *WorkflowSkipTimerRequest) HasTimerCommandIndex() bool`

HasTimerCommandIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


