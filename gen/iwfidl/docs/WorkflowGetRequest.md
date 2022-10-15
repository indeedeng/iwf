# WorkflowGetRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**NeedsResults** | Pointer to **bool** |  | [optional] 
**WaitTimeSeconds** | Pointer to **int32** |  | [optional] 

## Methods

### NewWorkflowGetRequest

`func NewWorkflowGetRequest(workflowId string, ) *WorkflowGetRequest`

NewWorkflowGetRequest instantiates a new WorkflowGetRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetRequestWithDefaults

`func NewWorkflowGetRequestWithDefaults() *WorkflowGetRequest`

NewWorkflowGetRequestWithDefaults instantiates a new WorkflowGetRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowGetRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowGetRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowGetRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowGetRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowGetRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetNeedsResults

`func (o *WorkflowGetRequest) GetNeedsResults() bool`

GetNeedsResults returns the NeedsResults field if non-nil, zero value otherwise.

### GetNeedsResultsOk

`func (o *WorkflowGetRequest) GetNeedsResultsOk() (*bool, bool)`

GetNeedsResultsOk returns a tuple with the NeedsResults field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNeedsResults

`func (o *WorkflowGetRequest) SetNeedsResults(v bool)`

SetNeedsResults sets NeedsResults field to given value.

### HasNeedsResults

`func (o *WorkflowGetRequest) HasNeedsResults() bool`

HasNeedsResults returns a boolean if a field has been set.

### GetWaitTimeSeconds

`func (o *WorkflowGetRequest) GetWaitTimeSeconds() int32`

GetWaitTimeSeconds returns the WaitTimeSeconds field if non-nil, zero value otherwise.

### GetWaitTimeSecondsOk

`func (o *WorkflowGetRequest) GetWaitTimeSecondsOk() (*int32, bool)`

GetWaitTimeSecondsOk returns a tuple with the WaitTimeSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitTimeSeconds

`func (o *WorkflowGetRequest) SetWaitTimeSeconds(v int32)`

SetWaitTimeSeconds sets WaitTimeSeconds field to given value.

### HasWaitTimeSeconds

`func (o *WorkflowGetRequest) HasWaitTimeSeconds() bool`

HasWaitTimeSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


