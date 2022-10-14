# SignalResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**SignalRequestStatus** | **string** |  | 
**SignalChannelName** | **string** |  | 
**SignalValue** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewSignalResult

`func NewSignalResult(commandId string, signalRequestStatus string, signalChannelName string, ) *SignalResult`

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


### GetSignalRequestStatus

`func (o *SignalResult) GetSignalRequestStatus() string`

GetSignalRequestStatus returns the SignalRequestStatus field if non-nil, zero value otherwise.

### GetSignalRequestStatusOk

`func (o *SignalResult) GetSignalRequestStatusOk() (*string, bool)`

GetSignalRequestStatusOk returns a tuple with the SignalRequestStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalRequestStatus

`func (o *SignalResult) SetSignalRequestStatus(v string)`

SetSignalRequestStatus sets SignalRequestStatus field to given value.


### GetSignalChannelName

`func (o *SignalResult) GetSignalChannelName() string`

GetSignalChannelName returns the SignalChannelName field if non-nil, zero value otherwise.

### GetSignalChannelNameOk

`func (o *SignalResult) GetSignalChannelNameOk() (*string, bool)`

GetSignalChannelNameOk returns a tuple with the SignalChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalChannelName

`func (o *SignalResult) SetSignalChannelName(v string)`

SetSignalChannelName sets SignalChannelName field to given value.


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


