# WorkflowQueryRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | Pointer to **string** |  | [optional] 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**AttributeKeys** | Pointer to **[]string** |  | [optional] 

## Methods

### NewWorkflowQueryRequest

`func NewWorkflowQueryRequest() *WorkflowQueryRequest`

NewWorkflowQueryRequest instantiates a new WorkflowQueryRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowQueryRequestWithDefaults

`func NewWorkflowQueryRequestWithDefaults() *WorkflowQueryRequest`

NewWorkflowQueryRequestWithDefaults instantiates a new WorkflowQueryRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowQueryRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowQueryRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowQueryRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.

### HasWorkflowId

`func (o *WorkflowQueryRequest) HasWorkflowId() bool`

HasWorkflowId returns a boolean if a field has been set.

### GetWorkflowRunId

`func (o *WorkflowQueryRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowQueryRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowQueryRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowQueryRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetAttributeKeys

`func (o *WorkflowQueryRequest) GetAttributeKeys() []string`

GetAttributeKeys returns the AttributeKeys field if non-nil, zero value otherwise.

### GetAttributeKeysOk

`func (o *WorkflowQueryRequest) GetAttributeKeysOk() (*[]string, bool)`

GetAttributeKeysOk returns a tuple with the AttributeKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributeKeys

`func (o *WorkflowQueryRequest) SetAttributeKeys(v []string)`

SetAttributeKeys sets AttributeKeys field to given value.

### HasAttributeKeys

`func (o *WorkflowQueryRequest) HasAttributeKeys() bool`

HasAttributeKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


