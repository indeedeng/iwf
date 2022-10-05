# WorkflowSearchResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**QueryAttributes** | Pointer to [**[]WorkflowSearchResponseEntry**](WorkflowSearchResponseEntry.md) |  | [optional] 

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

### GetQueryAttributes

`func (o *WorkflowSearchResponse) GetQueryAttributes() []WorkflowSearchResponseEntry`

GetQueryAttributes returns the QueryAttributes field if non-nil, zero value otherwise.

### GetQueryAttributesOk

`func (o *WorkflowSearchResponse) GetQueryAttributesOk() (*[]WorkflowSearchResponseEntry, bool)`

GetQueryAttributesOk returns a tuple with the QueryAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryAttributes

`func (o *WorkflowSearchResponse) SetQueryAttributes(v []WorkflowSearchResponseEntry)`

SetQueryAttributes sets QueryAttributes field to given value.

### HasQueryAttributes

`func (o *WorkflowSearchResponse) HasQueryAttributes() bool`

HasQueryAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


