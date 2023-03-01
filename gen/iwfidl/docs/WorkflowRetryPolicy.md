# WorkflowRetryPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InitialIntervalSeconds** | Pointer to **int32** |  | [optional] 
**BackoffCoefficient** | Pointer to **float32** |  | [optional] 
**MaximumIntervalSeconds** | Pointer to **int32** |  | [optional] 
**MaximumAttempts** | Pointer to **int32** |  | [optional] 

## Methods

### NewWorkflowRetryPolicy

`func NewWorkflowRetryPolicy() *WorkflowRetryPolicy`

NewWorkflowRetryPolicy instantiates a new WorkflowRetryPolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWorkflowRetryPolicyWithDefaults

`func NewWorkflowRetryPolicyWithDefaults() *WorkflowRetryPolicy`

NewWorkflowRetryPolicyWithDefaults instantiates a new WorkflowRetryPolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInitialIntervalSeconds

`func (o *WorkflowRetryPolicy) GetInitialIntervalSeconds() int32`

GetInitialIntervalSeconds returns the InitialIntervalSeconds field if non-nil, zero value otherwise.

### GetInitialIntervalSecondsOk

`func (o *WorkflowRetryPolicy) GetInitialIntervalSecondsOk() (*int32, bool)`

GetInitialIntervalSecondsOk returns a tuple with the InitialIntervalSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialIntervalSeconds

`func (o *WorkflowRetryPolicy) SetInitialIntervalSeconds(v int32)`

SetInitialIntervalSeconds sets InitialIntervalSeconds field to given value.

### HasInitialIntervalSeconds

`func (o *WorkflowRetryPolicy) HasInitialIntervalSeconds() bool`

HasInitialIntervalSeconds returns a boolean if a field has been set.

### GetBackoffCoefficient

`func (o *WorkflowRetryPolicy) GetBackoffCoefficient() float32`

GetBackoffCoefficient returns the BackoffCoefficient field if non-nil, zero value otherwise.

### GetBackoffCoefficientOk

`func (o *WorkflowRetryPolicy) GetBackoffCoefficientOk() (*float32, bool)`

GetBackoffCoefficientOk returns a tuple with the BackoffCoefficient field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackoffCoefficient

`func (o *WorkflowRetryPolicy) SetBackoffCoefficient(v float32)`

SetBackoffCoefficient sets BackoffCoefficient field to given value.

### HasBackoffCoefficient

`func (o *WorkflowRetryPolicy) HasBackoffCoefficient() bool`

HasBackoffCoefficient returns a boolean if a field has been set.

### GetMaximumIntervalSeconds

`func (o *WorkflowRetryPolicy) GetMaximumIntervalSeconds() int32`

GetMaximumIntervalSeconds returns the MaximumIntervalSeconds field if non-nil, zero value otherwise.

### GetMaximumIntervalSecondsOk

`func (o *WorkflowRetryPolicy) GetMaximumIntervalSecondsOk() (*int32, bool)`

GetMaximumIntervalSecondsOk returns a tuple with the MaximumIntervalSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumIntervalSeconds

`func (o *WorkflowRetryPolicy) SetMaximumIntervalSeconds(v int32)`

SetMaximumIntervalSeconds sets MaximumIntervalSeconds field to given value.

### HasMaximumIntervalSeconds

`func (o *WorkflowRetryPolicy) HasMaximumIntervalSeconds() bool`

HasMaximumIntervalSeconds returns a boolean if a field has been set.

### GetMaximumAttempts

`func (o *WorkflowRetryPolicy) GetMaximumAttempts() int32`

GetMaximumAttempts returns the MaximumAttempts field if non-nil, zero value otherwise.

### GetMaximumAttemptsOk

`func (o *WorkflowRetryPolicy) GetMaximumAttemptsOk() (*int32, bool)`

GetMaximumAttemptsOk returns a tuple with the MaximumAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumAttempts

`func (o *WorkflowRetryPolicy) SetMaximumAttempts(v int32)`

SetMaximumAttempts sets MaximumAttempts field to given value.

### HasMaximumAttempts

`func (o *WorkflowRetryPolicy) HasMaximumAttempts() bool`

HasMaximumAttempts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


