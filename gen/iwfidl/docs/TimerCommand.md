# TimerCommand

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**FiringUnixTimestampSeconds** | **int64** |  | 

## Methods

### NewTimerCommand

`func NewTimerCommand(commandId string, firingUnixTimestampSeconds int64, ) *TimerCommand`

NewTimerCommand instantiates a new TimerCommand object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimerCommandWithDefaults

`func NewTimerCommandWithDefaults() *TimerCommand`

NewTimerCommandWithDefaults instantiates a new TimerCommand object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *TimerCommand) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *TimerCommand) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *TimerCommand) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetFiringUnixTimestampSeconds

`func (o *TimerCommand) GetFiringUnixTimestampSeconds() int64`

GetFiringUnixTimestampSeconds returns the FiringUnixTimestampSeconds field if non-nil, zero value otherwise.

### GetFiringUnixTimestampSecondsOk

`func (o *TimerCommand) GetFiringUnixTimestampSecondsOk() (*int64, bool)`

GetFiringUnixTimestampSecondsOk returns a tuple with the FiringUnixTimestampSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFiringUnixTimestampSeconds

`func (o *TimerCommand) SetFiringUnixTimestampSeconds(v int64)`

SetFiringUnixTimestampSeconds sets FiringUnixTimestampSeconds field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


