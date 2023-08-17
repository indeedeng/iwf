# HealthInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Condition** | Pointer to **string** |  | [optional] 
**Hostname** | Pointer to **string** |  | [optional] 
**Duration** | Pointer to **int32** |  | [optional] 

## Methods

### NewHealthInfo

`func NewHealthInfo() *HealthInfo`

NewHealthInfo instantiates a new HealthInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHealthInfoWithDefaults

`func NewHealthInfoWithDefaults() *HealthInfo`

NewHealthInfoWithDefaults instantiates a new HealthInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCondition

`func (o *HealthInfo) GetCondition() string`

GetCondition returns the Condition field if non-nil, zero value otherwise.

### GetConditionOk

`func (o *HealthInfo) GetConditionOk() (*string, bool)`

GetConditionOk returns a tuple with the Condition field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCondition

`func (o *HealthInfo) SetCondition(v string)`

SetCondition sets Condition field to given value.

### HasCondition

`func (o *HealthInfo) HasCondition() bool`

HasCondition returns a boolean if a field has been set.

### GetHostname

`func (o *HealthInfo) GetHostname() string`

GetHostname returns the Hostname field if non-nil, zero value otherwise.

### GetHostnameOk

`func (o *HealthInfo) GetHostnameOk() (*string, bool)`

GetHostnameOk returns a tuple with the Hostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostname

`func (o *HealthInfo) SetHostname(v string)`

SetHostname sets Hostname field to given value.

### HasHostname

`func (o *HealthInfo) HasHostname() bool`

HasHostname returns a boolean if a field has been set.

### GetDuration

`func (o *HealthInfo) GetDuration() int32`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *HealthInfo) GetDurationOk() (*int32, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *HealthInfo) SetDuration(v int32)`

SetDuration sets Duration field to given value.

### HasDuration

`func (o *HealthInfo) HasDuration() bool`

HasDuration returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


