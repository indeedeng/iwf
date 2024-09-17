# InterStateChannelCommand

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | Pointer to **string** |  | [optional] 
**ChannelName** | **string** |  | 
**AtLeast** | Pointer to **int32** |  | [optional] 
**AtMost** | Pointer to **int32** |  | [optional] 

## Methods

### NewInterStateChannelCommand

`func NewInterStateChannelCommand(channelName string, ) *InterStateChannelCommand`

NewInterStateChannelCommand instantiates a new InterStateChannelCommand object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInterStateChannelCommandWithDefaults

`func NewInterStateChannelCommandWithDefaults() *InterStateChannelCommand`

NewInterStateChannelCommandWithDefaults instantiates a new InterStateChannelCommand object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *InterStateChannelCommand) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *InterStateChannelCommand) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *InterStateChannelCommand) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.

### HasCommandId

`func (o *InterStateChannelCommand) HasCommandId() bool`

HasCommandId returns a boolean if a field has been set.

### GetChannelName

`func (o *InterStateChannelCommand) GetChannelName() string`

GetChannelName returns the ChannelName field if non-nil, zero value otherwise.

### GetChannelNameOk

`func (o *InterStateChannelCommand) GetChannelNameOk() (*string, bool)`

GetChannelNameOk returns a tuple with the ChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelName

`func (o *InterStateChannelCommand) SetChannelName(v string)`

SetChannelName sets ChannelName field to given value.


### GetAtLeast

`func (o *InterStateChannelCommand) GetAtLeast() int32`

GetAtLeast returns the AtLeast field if non-nil, zero value otherwise.

### GetAtLeastOk

`func (o *InterStateChannelCommand) GetAtLeastOk() (*int32, bool)`

GetAtLeastOk returns a tuple with the AtLeast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAtLeast

`func (o *InterStateChannelCommand) SetAtLeast(v int32)`

SetAtLeast sets AtLeast field to given value.

### HasAtLeast

`func (o *InterStateChannelCommand) HasAtLeast() bool`

HasAtLeast returns a boolean if a field has been set.

### GetAtMost

`func (o *InterStateChannelCommand) GetAtMost() int32`

GetAtMost returns the AtMost field if non-nil, zero value otherwise.

### GetAtMostOk

`func (o *InterStateChannelCommand) GetAtMostOk() (*int32, bool)`

GetAtMostOk returns a tuple with the AtMost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAtMost

`func (o *InterStateChannelCommand) SetAtMost(v int32)`

SetAtMost sets AtMost field to given value.

### HasAtMost

`func (o *InterStateChannelCommand) HasAtMost() bool`

HasAtMost returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


