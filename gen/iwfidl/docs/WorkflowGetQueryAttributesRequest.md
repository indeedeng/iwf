# WorkflowGetQueryAttributesRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**AttributeKeys** | Pointer to **[]string** |  | [optional] 

## Methods

### NewWorkflowGetQueryAttributesRequest

`func NewWorkflowGetQueryAttributesRequest(workflowId string, ) *WorkflowGetQueryAttributesRequest`

NewWorkflowGetQueryAttributesRequest instantiates a new WorkflowGetQueryAttributesRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetQueryAttributesRequestWithDefaults

`func NewWorkflowGetQueryAttributesRequestWithDefaults() *WorkflowGetQueryAttributesRequest`

NewWorkflowGetQueryAttributesRequestWithDefaults instantiates a new WorkflowGetQueryAttributesRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowGetQueryAttributesRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowGetQueryAttributesRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowGetQueryAttributesRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowGetQueryAttributesRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetQueryAttributesRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetQueryAttributesRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowGetQueryAttributesRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetAttributeKeys

`func (o *WorkflowGetQueryAttributesRequest) GetAttributeKeys() []string`

GetAttributeKeys returns the AttributeKeys field if non-nil, zero value otherwise.

### GetAttributeKeysOk

`func (o *WorkflowGetQueryAttributesRequest) GetAttributeKeysOk() (*[]string, bool)`

GetAttributeKeysOk returns a tuple with the AttributeKeys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAttributeKeys

`func (o *WorkflowGetQueryAttributesRequest) SetAttributeKeys(v []string)`

SetAttributeKeys sets AttributeKeys field to given value.

### HasAttributeKeys

`func (o *WorkflowGetQueryAttributesRequest) HasAttributeKeys() bool`

HasAttributeKeys returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


