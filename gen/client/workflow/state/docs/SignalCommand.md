# SignalCommand

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | Pointer to **string** |  | [optional] 
**SignalName** | Pointer to **string** |  | [optional] 

## Methods

### NewSignalCommand

`func NewSignalCommand() *SignalCommand`

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

### GetSignalName

`func (o *SignalCommand) GetSignalName() string`

GetSignalName returns the SignalName field if non-nil, zero value otherwise.

### GetSignalNameOk

`func (o *SignalCommand) GetSignalNameOk() (*string, bool)`

GetSignalNameOk returns a tuple with the SignalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalName

`func (o *SignalCommand) SetSignalName(v string)`

SetSignalName sets SignalName field to given value.

### HasSignalName

`func (o *SignalCommand) HasSignalName() bool`

HasSignalName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


