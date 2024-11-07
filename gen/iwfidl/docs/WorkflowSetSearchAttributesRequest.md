# WorkflowSetSearchAttributesRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 

## Methods

### NewWorkflowSetSearchAttributesRequest

`func NewWorkflowSetSearchAttributesRequest(workflowId string, ) *WorkflowSetSearchAttributesRequest`

NewWorkflowSetSearchAttributesRequest instantiates a new WorkflowSetSearchAttributesRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSetSearchAttributesRequestWithDefaults

`func NewWorkflowSetSearchAttributesRequestWithDefaults() *WorkflowSetSearchAttributesRequest`

NewWorkflowSetSearchAttributesRequestWithDefaults instantiates a new WorkflowSetSearchAttributesRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowSetSearchAttributesRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowSetSearchAttributesRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowSetSearchAttributesRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowSetSearchAttributesRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowSetSearchAttributesRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowSetSearchAttributesRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowSetSearchAttributesRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowSetSearchAttributesRequest) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowSetSearchAttributesRequest) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowSetSearchAttributesRequest) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowSetSearchAttributesRequest) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


