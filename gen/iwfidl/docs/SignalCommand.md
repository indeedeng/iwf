# SignalCommand

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | Pointer to **string** |  | [optional] 
**SignalChannelName** | **string** |  | 
**AtLeast** | Pointer to **int32** |  | [optional] 
**AtMost** | Pointer to **int32** |  | [optional] 

## Methods

### NewSignalCommand

`func NewSignalCommand(signalChannelName string, ) *SignalCommand`

NewSignalCommand instantiates a new SignalCommand object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignalCommandWithDefaults

`func NewSignalCommandWithDefaults() *SignalCommand`

NewSignalCommandWithDefaults instantiates a new SignalCommand object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *SignalCommand) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *SignalCommand) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *SignalCommand) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.

### HasCommandId

`func (o *SignalCommand) HasCommandId() bool`

HasCommandId returns a boolean if a field has been set.

### GetSignalChannelName

`func (o *SignalCommand) GetSignalChannelName() string`

GetSignalChannelName returns the SignalChannelName field if non-nil, zero value otherwise.

### GetSignalChannelNameOk

`func (o *SignalCommand) GetSignalChannelNameOk() (*string, bool)`

GetSignalChannelNameOk returns a tuple with the SignalChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalChannelName

`func (o *SignalCommand) SetSignalChannelName(v string)`

SetSignalChannelName sets SignalChannelName field to given value.


### GetAtLeast

`func (o *SignalCommand) GetAtLeast() int32`

GetAtLeast returns the AtLeast field if non-nil, zero value otherwise.

### GetAtLeastOk

`func (o *SignalCommand) GetAtLeastOk() (*int32, bool)`

GetAtLeastOk returns a tuple with the AtLeast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAtLeast

`func (o *SignalCommand) SetAtLeast(v int32)`

SetAtLeast sets AtLeast field to given value.

### HasAtLeast

`func (o *SignalCommand) HasAtLeast() bool`

HasAtLeast returns a boolean if a field has been set.

### GetAtMost

`func (o *SignalCommand) GetAtMost() int32`

GetAtMost returns the AtMost field if non-nil, zero value otherwise.

### GetAtMostOk

`func (o *SignalCommand) GetAtMostOk() (*int32, bool)`

GetAtMostOk returns a tuple with the AtMost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAtMost

`func (o *SignalCommand) SetAtMost(v int32)`

SetAtMost sets AtMost field to given value.

### HasAtMost

`func (o *SignalCommand) HasAtMost() bool`

HasAtMost returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


