# SignalResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**SignalStatus** | **string** |  | 
**SignalName** | **string** |  | 
**SignalValue** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewSignalResult

`func NewSignalResult(commandId string, signalStatus string, signalName string, ) *SignalResult`

NewSignalResult instantiates a new SignalResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSignalResultWithDefaults

`func NewSignalResultWithDefaults() *SignalResult`

NewSignalResultWithDefaults instantiates a new SignalResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *SignalResult) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *SignalResult) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *SignalResult) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetSignalStatus

`func (o *SignalResult) GetSignalStatus() string`

GetSignalStatus returns the SignalStatus field if non-nil, zero value otherwise.

### GetSignalStatusOk

`func (o *SignalResult) GetSignalStatusOk() (*string, bool)`

GetSignalStatusOk returns a tuple with the SignalStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalStatus

`func (o *SignalResult) SetSignalStatus(v string)`

SetSignalStatus sets SignalStatus field to given value.


### GetSignalName

`func (o *SignalResult) GetSignalName() string`

GetSignalName returns the SignalName field if non-nil, zero value otherwise.

### GetSignalNameOk

`func (o *SignalResult) GetSignalNameOk() (*string, bool)`

GetSignalNameOk returns a tuple with the SignalName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalName

`func (o *SignalResult) SetSignalName(v string)`

SetSignalName sets SignalName field to given value.


### GetSignalValue

`func (o *SignalResult) GetSignalValue() EncodedObject`

GetSignalValue returns the SignalValue field if non-nil, zero value otherwise.

### GetSignalValueOk

`func (o *SignalResult) GetSignalValueOk() (*EncodedObject, bool)`

GetSignalValueOk returns a tuple with the SignalValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalValue

`func (o *SignalResult) SetSignalValue(v EncodedObject)`

SetSignalValue sets SignalValue field to given value.

### HasSignalValue

`func (o *SignalResult) HasSignalValue() bool`

HasSignalValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


