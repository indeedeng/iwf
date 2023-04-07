# WorkflowConfigUpdateRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**WorkflowConfig** | [**WorkflowConfig**](WorkflowConfig.md) |  | 

## Methods

### NewWorkflowConfigUpdateRequest

`func NewWorkflowConfigUpdateRequest(workflowId string, workflowConfig WorkflowConfig, ) *WorkflowConfigUpdateRequest`

NewWorkflowConfigUpdateRequest instantiates a new WorkflowConfigUpdateRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowConfigUpdateRequestWithDefaults

`func NewWorkflowConfigUpdateRequestWithDefaults() *WorkflowConfigUpdateRequest`

NewWorkflowConfigUpdateRequestWithDefaults instantiates a new WorkflowConfigUpdateRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowConfigUpdateRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowConfigUpdateRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowConfigUpdateRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowConfigUpdateRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowConfigUpdateRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowConfigUpdateRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowConfigUpdateRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetWorkflowConfig

`func (o *WorkflowConfigUpdateRequest) GetWorkflowConfig() WorkflowConfig`

GetWorkflowConfig returns the WorkflowConfig field if non-nil, zero value otherwise.

### GetWorkflowConfigOk

`func (o *WorkflowConfigUpdateRequest) GetWorkflowConfigOk() (*WorkflowConfig, bool)`

GetWorkflowConfigOk returns a tuple with the WorkflowConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowConfig

`func (o *WorkflowConfigUpdateRequest) SetWorkflowConfig(v WorkflowConfig)`

SetWorkflowConfig sets WorkflowConfig field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


