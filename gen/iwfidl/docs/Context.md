# Context

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | **string** |  | 
**WorkflowStartedTimestamp** | **int64** |  | 
**StateExecutionId** | **string** |  | 
**FirstAttemptTimestamp** | Pointer to **int64** |  | [optional] 
**Attempt** | Pointer to **int32** |  | [optional] 

## Methods

### NewContext

`func NewContext(workflowId string, workflowRunId string, workflowStartedTimestamp int64, stateExecutionId string, ) *Context`

NewContext instantiates a new Context object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContextWithDefaults

`func NewContextWithDefaults() *Context`

NewContextWithDefaults instantiates a new Context object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *Context) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *Context) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *Context) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *Context) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *Context) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *Context) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.


### GetWorkflowStartedTimestamp

`func (o *Context) GetWorkflowStartedTimestamp() int64`

GetWorkflowStartedTimestamp returns the WorkflowStartedTimestamp field if non-nil, zero value otherwise.

### GetWorkflowStartedTimestampOk

`func (o *Context) GetWorkflowStartedTimestampOk() (*int64, bool)`

GetWorkflowStartedTimestampOk returns a tuple with the WorkflowStartedTimestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowStartedTimestamp

`func (o *Context) SetWorkflowStartedTimestamp(v int64)`

SetWorkflowStartedTimestamp sets WorkflowStartedTimestamp field to given value.


### GetStateExecutionId

`func (o *Context) GetStateExecutionId() string`

GetStateExecutionId returns the StateExecutionId field if non-nil, zero value otherwise.

### GetStateExecutionIdOk

`func (o *Context) GetStateExecutionIdOk() (*string, bool)`

GetStateExecutionIdOk returns a tuple with the StateExecutionId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateExecutionId

`func (o *Context) SetStateExecutionId(v string)`

SetStateExecutionId sets StateExecutionId field to given value.


### GetFirstAttemptTimestamp

`func (o *Context) GetFirstAttemptTimestamp() int64`

GetFirstAttemptTimestamp returns the FirstAttemptTimestamp field if non-nil, zero value otherwise.

### GetFirstAttemptTimestampOk

`func (o *Context) GetFirstAttemptTimestampOk() (*int64, bool)`

GetFirstAttemptTimestampOk returns a tuple with the FirstAttemptTimestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirstAttemptTimestamp

`func (o *Context) SetFirstAttemptTimestamp(v int64)`

SetFirstAttemptTimestamp sets FirstAttemptTimestamp field to given value.

### HasFirstAttemptTimestamp

`func (o *Context) HasFirstAttemptTimestamp() bool`

HasFirstAttemptTimestamp returns a boolean if a field has been set.

### GetAttempt

`func (o *Context) GetAttempt() int32`

GetAttempt returns the Attempt field if non-nil, zero value otherwise.

### GetAttemptOk

`func (o *Context) GetAttemptOk() (*int32, bool)`

GetAttemptOk returns a tuple with the Attempt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttempt

`func (o *Context) SetAttempt(v int32)`

SetAttempt sets Attempt field to given value.

### HasAttempt

`func (o *Context) HasAttempt() bool`

HasAttempt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


