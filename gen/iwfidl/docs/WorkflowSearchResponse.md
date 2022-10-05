# WorkflowSearchResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowExecutions** | Pointer to [**[]WorkflowSearchResponseEntry**](WorkflowSearchResponseEntry.md) |  | [optional] 

## Methods

### NewWorkflowSearchResponse

`func NewWorkflowSearchResponse() *WorkflowSearchResponse`

NewWorkflowSearchResponse instantiates a new WorkflowSearchResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowSearchResponseWithDefaults

`func NewWorkflowSearchResponseWithDefaults() *WorkflowSearchResponse`

NewWorkflowSearchResponseWithDefaults instantiates a new WorkflowSearchResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowExecutions

`func (o *WorkflowSearchResponse) GetWorkflowExecutions() []WorkflowSearchResponseEntry`

GetWorkflowExecutions returns the WorkflowExecutions field if non-nil, zero value otherwise.

### GetWorkflowExecutionsOk

`func (o *WorkflowSearchResponse) GetWorkflowExecutionsOk() (*[]WorkflowSearchResponseEntry, bool)`

GetWorkflowExecutionsOk returns a tuple with the WorkflowExecutions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowExecutions

`func (o *WorkflowSearchResponse) SetWorkflowExecutions(v []WorkflowSearchResponseEntry)`

SetWorkflowExecutions sets WorkflowExecutions field to given value.

### HasWorkflowExecutions

`func (o *WorkflowSearchResponse) HasWorkflowExecutions() bool`

HasWorkflowExecutions returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


