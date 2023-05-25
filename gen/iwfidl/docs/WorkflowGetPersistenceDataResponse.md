# WorkflowGetPersistenceDataResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DataAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**CachedDataAttributes** | Pointer to [**[]KeyValue**](KeyValue.md) |  | [optional] 
**SearchAttributes** | Pointer to [**[]SearchAttribute**](SearchAttribute.md) |  | [optional] 

## Methods

### NewWorkflowGetPersistenceDataResponse

`func NewWorkflowGetPersistenceDataResponse() *WorkflowGetPersistenceDataResponse`

NewWorkflowGetPersistenceDataResponse instantiates a new WorkflowGetPersistenceDataResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowGetPersistenceDataResponseWithDefaults

`func NewWorkflowGetPersistenceDataResponseWithDefaults() *WorkflowGetPersistenceDataResponse`

NewWorkflowGetPersistenceDataResponseWithDefaults instantiates a new WorkflowGetPersistenceDataResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) GetDataAttributes() []KeyValue`

GetDataAttributes returns the DataAttributes field if non-nil, zero value otherwise.

### GetDataAttributesOk

`func (o *WorkflowGetPersistenceDataResponse) GetDataAttributesOk() (*[]KeyValue, bool)`

GetDataAttributesOk returns a tuple with the DataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) SetDataAttributes(v []KeyValue)`

SetDataAttributes sets DataAttributes field to given value.

### HasDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) HasDataAttributes() bool`

HasDataAttributes returns a boolean if a field has been set.

### GetCachedDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) GetCachedDataAttributes() []KeyValue`

GetCachedDataAttributes returns the CachedDataAttributes field if non-nil, zero value otherwise.

### GetCachedDataAttributesOk

`func (o *WorkflowGetPersistenceDataResponse) GetCachedDataAttributesOk() (*[]KeyValue, bool)`

GetCachedDataAttributesOk returns a tuple with the CachedDataAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCachedDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) SetCachedDataAttributes(v []KeyValue)`

SetCachedDataAttributes sets CachedDataAttributes field to given value.

### HasCachedDataAttributes

`func (o *WorkflowGetPersistenceDataResponse) HasCachedDataAttributes() bool`

HasCachedDataAttributes returns a boolean if a field has been set.

### GetSearchAttributes

`func (o *WorkflowGetPersistenceDataResponse) GetSearchAttributes() []SearchAttribute`

GetSearchAttributes returns the SearchAttributes field if non-nil, zero value otherwise.

### GetSearchAttributesOk

`func (o *WorkflowGetPersistenceDataResponse) GetSearchAttributesOk() (*[]SearchAttribute, bool)`

GetSearchAttributesOk returns a tuple with the SearchAttributes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributes

`func (o *WorkflowGetPersistenceDataResponse) SetSearchAttributes(v []SearchAttribute)`

SetSearchAttributes sets SearchAttributes field to given value.

### HasSearchAttributes

`func (o *WorkflowGetPersistenceDataResponse) HasSearchAttributes() bool`

HasSearchAttributes returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


