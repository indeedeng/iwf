# WorkflowQueryResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**QueryAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 

## Methods

### NewWorkflowQueryResponse

`func NewWorkflowQueryResponse() *WorkflowQueryResponse`

NewWorkflowQueryResponse instantiates a new WorkflowQueryResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowQueryResponseWithDefaults

`func NewWorkflowQueryResponseWithDefaults() *WorkflowQueryResponse`

NewWorkflowQueryResponseWithDefaults instantiates a new WorkflowQueryResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowRunId

`func (o *WorkflowQueryResponse) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowQueryResponse) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowQueryResponse) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowQueryResponse) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetQueryAttributes

`func (o *WorkflowQueryResponse) GetQueryAttributes() []KeyValue`

GetQueryAttributes returns the QueryAttributes field if non-nil, zero value otherwise.

### GetQueryAttributesOk

`func (o *WorkflowQueryResponse) GetQueryAttributesOk() (*[]KeyValue, bool)`

GetQueryAttributesOk returns a tuple with the QueryAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryAttributes

`func (o *WorkflowQueryResponse) SetQueryAttributes(v []KeyValue)`

SetQueryAttributes sets QueryAttributes field to given value.

### HasQueryAttributes

`func (o *WorkflowQueryResponse) HasQueryAttributes() bool`

HasQueryAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


