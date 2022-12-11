# RetryPolicy

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**InitialInterval** | Pointer to **float32** |  | [optional] 
**BackoffCoefficient** | Pointer to **float32** |  | [optional] 
**MaximumInterval** | Pointer to **float32** |  | [optional] 
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

### GetInitialInterval

`func (o *RetryPolicy) GetInitialInterval() float32`

GetInitialInterval returns the InitialInterval field if non-nil, zero value otherwise.

### GetInitialIntervalOk

`func (o *RetryPolicy) GetInitialIntervalOk() (*float32, bool)`

GetInitialIntervalOk returns a tuple with the InitialInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInitialInterval

`func (o *RetryPolicy) SetInitialInterval(v float32)`

SetInitialInterval sets InitialInterval field to given value.

### HasInitialInterval

`func (o *RetryPolicy) HasInitialInterval() bool`

HasInitialInterval returns a boolean if a field has been set.

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

### GetMaximumInterval

`func (o *RetryPolicy) GetMaximumInterval() float32`

GetMaximumInterval returns the MaximumInterval field if non-nil, zero value otherwise.

### GetMaximumIntervalOk

`func (o *RetryPolicy) GetMaximumIntervalOk() (*float32, bool)`

GetMaximumIntervalOk returns a tuple with the MaximumInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumInterval

`func (o *RetryPolicy) SetMaximumInterval(v float32)`

SetMaximumInterval sets MaximumInterval field to given value.

### HasMaximumInterval

`func (o *RetryPolicy) HasMaximumInterval() bool`

HasMaximumInterval returns a boolean if a field has been set.

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


