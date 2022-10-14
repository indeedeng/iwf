# WorkflowGetSearchAttributesRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**AttributeKeys** | Pointer to [**[]SearchAttributeKeyAndType**](SearchAttributeKeyAndType.md) |  | [optional] 

## Methods

### NewWorkflowGetSearchAttributesRequest

`func NewWorkflowGetSearchAttributesRequest(workflowId string, ) *WorkflowGetSearchAttributesRequest`

NewWorkflowGetSearchAttributesRequest instantiates a new WorkflowGetSearchAttributesRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetSearchAttributesRequestWithDefaults

`func NewWorkflowGetSearchAttributesRequestWithDefaults() *WorkflowGetSearchAttributesRequest`

NewWorkflowGetSearchAttributesRequestWithDefaults instantiates a new WorkflowGetSearchAttributesRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowGetSearchAttributesRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowGetSearchAttributesRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowGetSearchAttributesRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowGetSearchAttributesRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetSearchAttributesRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetSearchAttributesRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowGetSearchAttributesRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetAttributeKeys

`func (o *WorkflowGetSearchAttributesRequest) GetAttributeKeys() []SearchAttributeKeyAndType`

GetAttributeKeys returns the AttributeKeys field if non-nil, zero value otherwise.

### GetAttributeKeysOk

`func (o *WorkflowGetSearchAttributesRequest) GetAttributeKeysOk() (*[]SearchAttributeKeyAndType, bool)`

GetAttributeKeysOk returns a tuple with the AttributeKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributeKeys

`func (o *WorkflowGetSearchAttributesRequest) SetAttributeKeys(v []SearchAttributeKeyAndType)`

SetAttributeKeys sets AttributeKeys field to given value.

### HasAttributeKeys

`func (o *WorkflowGetSearchAttributesRequest) HasAttributeKeys() bool`

HasAttributeKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


