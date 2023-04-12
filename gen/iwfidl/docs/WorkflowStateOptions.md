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
**StartApiFailurePolicy** | Pointer to [**StartApiFailurePolicy**](StartApiFailurePolicy.md) |  | [optional] 
**SkipStartApi** | Pointer to **bool** |  | [optional] 
**WaitUntilApiTimeoutSeconds** | Pointer to **int32** |  | [optional] 
**ExecuteApiTimeoutSeconds** | Pointer to **int32** |  | [optional] 
**WaitUntilApiRetryPolicy** | Pointer to [**RetryPolicy**](RetryPolicy.md) |  | [optional] 
**ExecuteApiRetryPolicy** | Pointer to [**RetryPolicy**](RetryPolicy.md) |  | [optional] 
**WaitUntilApiFailurePolicy** | Pointer to [**WaitUntilApiFailurePolicy**](WaitUntilApiFailurePolicy.md) |  | [optional] 
**SkipWaitUntil** | Pointer to **bool** |  | [optional] 

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

### GetStartApiFailurePolicy

`func (o *WorkflowStateOptions) GetStartApiFailurePolicy() StartApiFailurePolicy`

GetStartApiFailurePolicy returns the StartApiFailurePolicy field if non-nil, zero value otherwise.

### GetStartApiFailurePolicyOk

`func (o *WorkflowStateOptions) GetStartApiFailurePolicyOk() (*StartApiFailurePolicy, bool)`

GetStartApiFailurePolicyOk returns a tuple with the StartApiFailurePolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartApiFailurePolicy

`func (o *WorkflowStateOptions) SetStartApiFailurePolicy(v StartApiFailurePolicy)`

SetStartApiFailurePolicy sets StartApiFailurePolicy field to given value.

### HasStartApiFailurePolicy

`func (o *WorkflowStateOptions) HasStartApiFailurePolicy() bool`

HasStartApiFailurePolicy returns a boolean if a field has been set.

### GetSkipStartApi

`func (o *WorkflowStateOptions) GetSkipStartApi() bool`

GetSkipStartApi returns the SkipStartApi field if non-nil, zero value otherwise.

### GetSkipStartApiOk

`func (o *WorkflowStateOptions) GetSkipStartApiOk() (*bool, bool)`

GetSkipStartApiOk returns a tuple with the SkipStartApi field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipStartApi

`func (o *WorkflowStateOptions) SetSkipStartApi(v bool)`

SetSkipStartApi sets SkipStartApi field to given value.

### HasSkipStartApi

`func (o *WorkflowStateOptions) HasSkipStartApi() bool`

HasSkipStartApi returns a boolean if a field has been set.

### GetWaitUntilApiTimeoutSeconds

`func (o *WorkflowStateOptions) GetWaitUntilApiTimeoutSeconds() int32`

GetWaitUntilApiTimeoutSeconds returns the WaitUntilApiTimeoutSeconds field if non-nil, zero value otherwise.

### GetWaitUntilApiTimeoutSecondsOk

`func (o *WorkflowStateOptions) GetWaitUntilApiTimeoutSecondsOk() (*int32, bool)`

GetWaitUntilApiTimeoutSecondsOk returns a tuple with the WaitUntilApiTimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitUntilApiTimeoutSeconds

`func (o *WorkflowStateOptions) SetWaitUntilApiTimeoutSeconds(v int32)`

SetWaitUntilApiTimeoutSeconds sets WaitUntilApiTimeoutSeconds field to given value.

### HasWaitUntilApiTimeoutSeconds

`func (o *WorkflowStateOptions) HasWaitUntilApiTimeoutSeconds() bool`

HasWaitUntilApiTimeoutSeconds returns a boolean if a field has been set.

### GetExecuteApiTimeoutSeconds

`func (o *WorkflowStateOptions) GetExecuteApiTimeoutSeconds() int32`

GetExecuteApiTimeoutSeconds returns the ExecuteApiTimeoutSeconds field if non-nil, zero value otherwise.

### GetExecuteApiTimeoutSecondsOk

`func (o *WorkflowStateOptions) GetExecuteApiTimeoutSecondsOk() (*int32, bool)`

GetExecuteApiTimeoutSecondsOk returns a tuple with the ExecuteApiTimeoutSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecuteApiTimeoutSeconds

`func (o *WorkflowStateOptions) SetExecuteApiTimeoutSeconds(v int32)`

SetExecuteApiTimeoutSeconds sets ExecuteApiTimeoutSeconds field to given value.

### HasExecuteApiTimeoutSeconds

`func (o *WorkflowStateOptions) HasExecuteApiTimeoutSeconds() bool`

HasExecuteApiTimeoutSeconds returns a boolean if a field has been set.

### GetWaitUntilApiRetryPolicy

`func (o *WorkflowStateOptions) GetWaitUntilApiRetryPolicy() RetryPolicy`

GetWaitUntilApiRetryPolicy returns the WaitUntilApiRetryPolicy field if non-nil, zero value otherwise.

### GetWaitUntilApiRetryPolicyOk

`func (o *WorkflowStateOptions) GetWaitUntilApiRetryPolicyOk() (*RetryPolicy, bool)`

GetWaitUntilApiRetryPolicyOk returns a tuple with the WaitUntilApiRetryPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitUntilApiRetryPolicy

`func (o *WorkflowStateOptions) SetWaitUntilApiRetryPolicy(v RetryPolicy)`

SetWaitUntilApiRetryPolicy sets WaitUntilApiRetryPolicy field to given value.

### HasWaitUntilApiRetryPolicy

`func (o *WorkflowStateOptions) HasWaitUntilApiRetryPolicy() bool`

HasWaitUntilApiRetryPolicy returns a boolean if a field has been set.

### GetExecuteApiRetryPolicy

`func (o *WorkflowStateOptions) GetExecuteApiRetryPolicy() RetryPolicy`

GetExecuteApiRetryPolicy returns the ExecuteApiRetryPolicy field if non-nil, zero value otherwise.

### GetExecuteApiRetryPolicyOk

`func (o *WorkflowStateOptions) GetExecuteApiRetryPolicyOk() (*RetryPolicy, bool)`

GetExecuteApiRetryPolicyOk returns a tuple with the ExecuteApiRetryPolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExecuteApiRetryPolicy

`func (o *WorkflowStateOptions) SetExecuteApiRetryPolicy(v RetryPolicy)`

SetExecuteApiRetryPolicy sets ExecuteApiRetryPolicy field to given value.

### HasExecuteApiRetryPolicy

`func (o *WorkflowStateOptions) HasExecuteApiRetryPolicy() bool`

HasExecuteApiRetryPolicy returns a boolean if a field has been set.

### GetWaitUntilApiFailurePolicy

`func (o *WorkflowStateOptions) GetWaitUntilApiFailurePolicy() WaitUntilApiFailurePolicy`

GetWaitUntilApiFailurePolicy returns the WaitUntilApiFailurePolicy field if non-nil, zero value otherwise.

### GetWaitUntilApiFailurePolicyOk

`func (o *WorkflowStateOptions) GetWaitUntilApiFailurePolicyOk() (*WaitUntilApiFailurePolicy, bool)`

GetWaitUntilApiFailurePolicyOk returns a tuple with the WaitUntilApiFailurePolicy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWaitUntilApiFailurePolicy

`func (o *WorkflowStateOptions) SetWaitUntilApiFailurePolicy(v WaitUntilApiFailurePolicy)`

SetWaitUntilApiFailurePolicy sets WaitUntilApiFailurePolicy field to given value.

### HasWaitUntilApiFailurePolicy

`func (o *WorkflowStateOptions) HasWaitUntilApiFailurePolicy() bool`

HasWaitUntilApiFailurePolicy returns a boolean if a field has been set.

### GetSkipWaitUntil

`func (o *WorkflowStateOptions) GetSkipWaitUntil() bool`

GetSkipWaitUntil returns the SkipWaitUntil field if non-nil, zero value otherwise.

### GetSkipWaitUntilOk

`func (o *WorkflowStateOptions) GetSkipWaitUntilOk() (*bool, bool)`

GetSkipWaitUntilOk returns a tuple with the SkipWaitUntil field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipWaitUntil

`func (o *WorkflowStateOptions) SetSkipWaitUntil(v bool)`

SetSkipWaitUntil sets SkipWaitUntil field to given value.

### HasSkipWaitUntil

`func (o *WorkflowStateOptions) HasSkipWaitUntil() bool`

HasSkipWaitUntil returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


