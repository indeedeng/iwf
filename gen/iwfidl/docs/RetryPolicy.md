# RetryPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InitialIntervalSeconds** | Pointer to **int32** |  | [optional] 
**BackoffCoefficient** | Pointer to **float32** |  | [optional] 
**MaximumIntervalSeconds** | Pointer to **int32** |  | [optional] 
**MaximumAttempts** | Pointer to **int32** |  | [optional] 

## Methods

### NewRetryPolicy

`func NewRetryPolicy() *RetryPolicy`

NewRetryPolicy instantiates a new RetryPolicy object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRetryPolicyWithDefaults

`func NewRetryPolicyWithDefaults() *RetryPolicy`

NewRetryPolicyWithDefaults instantiates a new RetryPolicy object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetInitialIntervalSeconds

`func (o *RetryPolicy) GetInitialIntervalSeconds() int32`

GetInitialIntervalSeconds returns the InitialIntervalSeconds field if non-nil, zero value otherwise.

### GetInitialIntervalSecondsOk

`func (o *RetryPolicy) GetInitialIntervalSecondsOk() (*int32, bool)`

GetInitialIntervalSecondsOk returns a tuple with the InitialIntervalSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialIntervalSeconds

`func (o *RetryPolicy) SetInitialIntervalSeconds(v int32)`

SetInitialIntervalSeconds sets InitialIntervalSeconds field to given value.

### HasInitialIntervalSeconds

`func (o *RetryPolicy) HasInitialIntervalSeconds() bool`

HasInitialIntervalSeconds returns a boolean if a field has been set.

### GetBackoffCoefficient

`func (o *RetryPolicy) GetBackoffCoefficient() float32`

GetBackoffCoefficient returns the BackoffCoefficient field if non-nil, zero value otherwise.

### GetBackoffCoefficientOk

`func (o *RetryPolicy) GetBackoffCoefficientOk() (*float32, bool)`

GetBackoffCoefficientOk returns a tuple with the BackoffCoefficient field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBackoffCoefficient

`func (o *RetryPolicy) SetBackoffCoefficient(v float32)`

SetBackoffCoefficient sets BackoffCoefficient field to given value.

### HasBackoffCoefficient

`func (o *RetryPolicy) HasBackoffCoefficient() bool`

HasBackoffCoefficient returns a boolean if a field has been set.

### GetMaximumIntervalSeconds

`func (o *RetryPolicy) GetMaximumIntervalSeconds() int32`

GetMaximumIntervalSeconds returns the MaximumIntervalSeconds field if non-nil, zero value otherwise.

### GetMaximumIntervalSecondsOk

`func (o *RetryPolicy) GetMaximumIntervalSecondsOk() (*int32, bool)`

GetMaximumIntervalSecondsOk returns a tuple with the MaximumIntervalSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumIntervalSeconds

`func (o *RetryPolicy) SetMaximumIntervalSeconds(v int32)`

SetMaximumIntervalSeconds sets MaximumIntervalSeconds field to given value.

### HasMaximumIntervalSeconds

`func (o *RetryPolicy) HasMaximumIntervalSeconds() bool`

HasMaximumIntervalSeconds returns a boolean if a field has been set.

### GetMaximumAttempts

`func (o *RetryPolicy) GetMaximumAttempts() int32`

GetMaximumAttempts returns the MaximumAttempts field if non-nil, zero value otherwise.

### GetMaximumAttemptsOk

`func (o *RetryPolicy) GetMaximumAttemptsOk() (*int32, bool)`

GetMaximumAttemptsOk returns a tuple with the MaximumAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumAttempts

`func (o *RetryPolicy) SetMaximumAttempts(v int32)`

SetMaximumAttempts sets MaximumAttempts field to given value.

### HasMaximumAttempts

`func (o *RetryPolicy) HasMaximumAttempts() bool`

HasMaximumAttempts returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


