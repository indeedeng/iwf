# WorkflowStateOptions

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SearchAttributesLoadingPolicy** | Pointer to [**PersistenceLoadingPolicy**](PersistenceLoadingPolicy.md) |  | [optional] 
**DataObjectsLoadingPolicy** | Pointer to [**PersistenceLoadingPolicy**](PersistenceLoadingPolicy.md) |  | [optional] 
**CommandCarryOverPolicy** | Pointer to [**CommandCarryOverPolicy**](CommandCarryOverPolicy.md) |  | [optional] 
**StartApiTimeoutSeconds** | Pointer to **int32** |  | [optional] 
**DecideApiTimeoutSeconds** | Pointer to **int32** |  | [optional] 
**StartApiRetryPolicy** | Pointer to [**RetryPolicy**](RetryPolicy.md) |  | [optional] 
**DecideApiRetryPolicy** | Pointer to [**RetryPolicy**](RetryPolicy.md) |  | [optional] 

## Methods

### NewWorkflowStateOptions

`func NewWorkflowStateOptions() *WorkflowStateOptions`

NewWorkflowStateOptions instantiates a new WorkflowStateOptions object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowStateOptionsWithDefaults

`func NewWorkflowStateOptionsWithDefaults() *WorkflowStateOptions`

NewWorkflowStateOptionsWithDefaults instantiates a new WorkflowStateOptions object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) GetSearchAttributesLoadingPolicy() PersistenceLoadingPolicy`

GetSearchAttributesLoadingPolicy returns the SearchAttributesLoadingPolicy field if non-nil, zero value otherwise.

### GetSearchAttributesLoadingPolicyOk

`func (o *WorkflowStateOptions) GetSearchAttributesLoadingPolicyOk() (*PersistenceLoadingPolicy, bool)`

GetSearchAttributesLoadingPolicyOk returns a tuple with the SearchAttributesLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) SetSearchAttributesLoadingPolicy(v PersistenceLoadingPolicy)`

SetSearchAttributesLoadingPolicy sets SearchAttributesLoadingPolicy field to given value.

### HasSearchAttributesLoadingPolicy

`func (o *WorkflowStateOptions) HasSearchAttributesLoadingPolicy() bool`

HasSearchAttributesLoadingPolicy returns a boolean if a field has been set.

### GetDataObjectsLoadingPolicy

`func (o *WorkflowStateOptions) GetDataObjectsLoadingPolicy() PersistenceLoadingPolicy`

GetDataObjectsLoadingPolicy returns the DataObjectsLoadingPolicy field if non-nil, zero value otherwise.

### GetDataObjectsLoadingPolicyOk

`func (o *WorkflowStateOptions) GetDataObjectsLoadingPolicyOk() (*PersistenceLoadingPolicy, bool)`

GetDataObjectsLoadingPolicyOk returns a tuple with the DataObjectsLoadingPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDataObjectsLoadingPolicy

`func (o *WorkflowStateOptions) SetDataObjectsLoadingPolicy(v PersistenceLoadingPolicy)`

SetDataObjectsLoadingPolicy sets DataObjectsLoadingPolicy field to given value.

### HasDataObjectsLoadingPolicy

`func (o *WorkflowStateOptions) HasDataObjectsLoadingPolicy() bool`

HasDataObjectsLoadingPolicy returns a boolean if a field has been set.

### GetCommandCarryOverPolicy

`func (o *WorkflowStateOptions) GetCommandCarryOverPolicy() CommandCarryOverPolicy`

GetCommandCarryOverPolicy returns the CommandCarryOverPolicy field if non-nil, zero value otherwise.

### GetCommandCarryOverPolicyOk

`func (o *WorkflowStateOptions) GetCommandCarryOverPolicyOk() (*CommandCarryOverPolicy, bool)`

GetCommandCarryOverPolicyOk returns a tuple with the CommandCarryOverPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandCarryOverPolicy

`func (o *WorkflowStateOptions) SetCommandCarryOverPolicy(v CommandCarryOverPolicy)`

SetCommandCarryOverPolicy sets CommandCarryOverPolicy field to given value.

### HasCommandCarryOverPolicy

`func (o *WorkflowStateOptions) HasCommandCarryOverPolicy() bool`

HasCommandCarryOverPolicy returns a boolean if a field has been set.

### GetStartApiTimeoutSeconds

`func (o *WorkflowStateOptions) GetStartApiTimeoutSeconds() int32`

GetStartApiTimeoutSeconds returns the StartApiTimeoutSeconds field if non-nil, zero value otherwise.

### GetStartApiTimeoutSecondsOk

`func (o *WorkflowStateOptions) GetStartApiTimeoutSecondsOk() (*int32, bool)`

GetStartApiTimeoutSecondsOk returns a tuple with the StartApiTimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartApiTimeoutSeconds

`func (o *WorkflowStateOptions) SetStartApiTimeoutSeconds(v int32)`

SetStartApiTimeoutSeconds sets StartApiTimeoutSeconds field to given value.

### HasStartApiTimeoutSeconds

`func (o *WorkflowStateOptions) HasStartApiTimeoutSeconds() bool`

HasStartApiTimeoutSeconds returns a boolean if a field has been set.

### GetDecideApiTimeoutSeconds

`func (o *WorkflowStateOptions) GetDecideApiTimeoutSeconds() int32`

GetDecideApiTimeoutSeconds returns the DecideApiTimeoutSeconds field if non-nil, zero value otherwise.

### GetDecideApiTimeoutSecondsOk

`func (o *WorkflowStateOptions) GetDecideApiTimeoutSecondsOk() (*int32, bool)`

GetDecideApiTimeoutSecondsOk returns a tuple with the DecideApiTimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecideApiTimeoutSeconds

`func (o *WorkflowStateOptions) SetDecideApiTimeoutSeconds(v int32)`

SetDecideApiTimeoutSeconds sets DecideApiTimeoutSeconds field to given value.

### HasDecideApiTimeoutSeconds

`func (o *WorkflowStateOptions) HasDecideApiTimeoutSeconds() bool`

HasDecideApiTimeoutSeconds returns a boolean if a field has been set.

### GetStartApiRetryPolicy

`func (o *WorkflowStateOptions) GetStartApiRetryPolicy() RetryPolicy`

GetStartApiRetryPolicy returns the StartApiRetryPolicy field if non-nil, zero value otherwise.

### GetStartApiRetryPolicyOk

`func (o *WorkflowStateOptions) GetStartApiRetryPolicyOk() (*RetryPolicy, bool)`

GetStartApiRetryPolicyOk returns a tuple with the StartApiRetryPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartApiRetryPolicy

`func (o *WorkflowStateOptions) SetStartApiRetryPolicy(v RetryPolicy)`

SetStartApiRetryPolicy sets StartApiRetryPolicy field to given value.

### HasStartApiRetryPolicy

`func (o *WorkflowStateOptions) HasStartApiRetryPolicy() bool`

HasStartApiRetryPolicy returns a boolean if a field has been set.

### GetDecideApiRetryPolicy

`func (o *WorkflowStateOptions) GetDecideApiRetryPolicy() RetryPolicy`

GetDecideApiRetryPolicy returns the DecideApiRetryPolicy field if non-nil, zero value otherwise.

### GetDecideApiRetryPolicyOk

`func (o *WorkflowStateOptions) GetDecideApiRetryPolicyOk() (*RetryPolicy, bool)`

GetDecideApiRetryPolicyOk returns a tuple with the DecideApiRetryPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDecideApiRetryPolicy

`func (o *WorkflowStateOptions) SetDecideApiRetryPolicy(v RetryPolicy)`

SetDecideApiRetryPolicy sets DecideApiRetryPolicy field to given value.

### HasDecideApiRetryPolicy

`func (o *WorkflowStateOptions) HasDecideApiRetryPolicy() bool`

HasDecideApiRetryPolicy returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


