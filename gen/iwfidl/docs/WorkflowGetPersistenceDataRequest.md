# WorkflowGetPersistenceDataRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**WorkflowId** | **string** |  | 
**WorkflowRunId** | Pointer to **string** |  | [optional] 
**DataAttributes** | Pointer to **[]string** |  | [optional] 
**CachedDataAttributes** | Pointer to **[]string** |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttributeKeyAndType**](SearchAttributeKeyAndType.md) |  | [optional] 

## Methods

### NewWorkflowGetPersistenceDataRequest

`func NewWorkflowGetPersistenceDataRequest(workflowId string, ) *WorkflowGetPersistenceDataRequest`

NewWorkflowGetPersistenceDataRequest instantiates a new WorkflowGetPersistenceDataRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetPersistenceDataRequestWithDefaults

`func NewWorkflowGetPersistenceDataRequestWithDefaults() *WorkflowGetPersistenceDataRequest`

NewWorkflowGetPersistenceDataRequestWithDefaults instantiates a new WorkflowGetPersistenceDataRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWorkflowId

`func (o *WorkflowGetPersistenceDataRequest) GetWorkflowId() string`

GetWorkflowId returns the WorkflowId field if non-nil, zero value otherwise.

### GetWorkflowIdOk

`func (o *WorkflowGetPersistenceDataRequest) GetWorkflowIdOk() (*string, bool)`

GetWorkflowIdOk returns a tuple with the WorkflowId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowId

`func (o *WorkflowGetPersistenceDataRequest) SetWorkflowId(v string)`

SetWorkflowId sets WorkflowId field to given value.


### GetWorkflowRunId

`func (o *WorkflowGetPersistenceDataRequest) GetWorkflowRunId() string`

GetWorkflowRunId returns the WorkflowRunId field if non-nil, zero value otherwise.

### GetWorkflowRunIdOk

`func (o *WorkflowGetPersistenceDataRequest) GetWorkflowRunIdOk() (*string, bool)`

GetWorkflowRunIdOk returns a tuple with the WorkflowRunId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkflowRunId

`func (o *WorkflowGetPersistenceDataRequest) SetWorkflowRunId(v string)`

SetWorkflowRunId sets WorkflowRunId field to given value.

### HasWorkflowRunId

`func (o *WorkflowGetPersistenceDataRequest) HasWorkflowRunId() bool`

HasWorkflowRunId returns a boolean if a field has been set.

### GetDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) GetDataAttributes() []string`

GetDataAttributes returns the DataAttributes field if non-nil, zero value otherwise.

### GetDataAttributesOk

`func (o *WorkflowGetPersistenceDataRequest) GetDataAttributesOk() (*[]string, bool)`

GetDataAttributesOk returns a tuple with the DataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) SetDataAttributes(v []string)`

SetDataAttributes sets DataAttributes field to given value.

### HasDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) HasDataAttributes() bool`

HasDataAttributes returns a boolean if a field has been set.

### GetCachedDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) GetCachedDataAttributes() []string`

GetCachedDataAttributes returns the CachedDataAttributes field if non-nil, zero value otherwise.

### GetCachedDataAttributesOk

`func (o *WorkflowGetPersistenceDataRequest) GetCachedDataAttributesOk() (*[]string, bool)`

GetCachedDataAttributesOk returns a tuple with the CachedDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCachedDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) SetCachedDataAttributes(v []string)`

SetCachedDataAttributes sets CachedDataAttributes field to given value.

### HasCachedDataAttributes

`func (o *WorkflowGetPersistenceDataRequest) HasCachedDataAttributes() bool`

HasCachedDataAttributes returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowGetPersistenceDataRequest) GetSearchAttributes() []SearchAttributeKeyAndType`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowGetPersistenceDataRequest) GetSearchAttributesOk() (*[]SearchAttributeKeyAndType, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowGetPersistenceDataRequest) SetSearchAttributes(v []SearchAttributeKeyAndType)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowGetPersistenceDataRequest) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


