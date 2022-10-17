# InterStateChannelResult

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CommandId** | **string** |  | 
**RequestStatus** | **string** |  | 
**ChannelName** | **string** |  | 
**Value** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewInterStateChannelResult

`func NewInterStateChannelResult(commandId string, requestStatus string, channelName string, ) *InterStateChannelResult`

NewInterStateChannelResult instantiates a new InterStateChannelResult object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInterStateChannelResultWithDefaults

`func NewInterStateChannelResultWithDefaults() *InterStateChannelResult`

NewInterStateChannelResultWithDefaults instantiates a new InterStateChannelResult object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCommandId

`func (o *InterStateChannelResult) GetCommandId() string`

GetCommandId returns the CommandId field if non-nil, zero value otherwise.

### GetCommandIdOk

`func (o *InterStateChannelResult) GetCommandIdOk() (*string, bool)`

GetCommandIdOk returns a tuple with the CommandId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommandId

`func (o *InterStateChannelResult) SetCommandId(v string)`

SetCommandId sets CommandId field to given value.


### GetRequestStatus

`func (o *InterStateChannelResult) GetRequestStatus() string`

GetRequestStatus returns the RequestStatus field if non-nil, zero value otherwise.

### GetRequestStatusOk

`func (o *InterStateChannelResult) GetRequestStatusOk() (*string, bool)`

GetRequestStatusOk returns a tuple with the RequestStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestStatus

`func (o *InterStateChannelResult) SetRequestStatus(v string)`

SetRequestStatus sets RequestStatus field to given value.


### GetChannelName

`func (o *InterStateChannelResult) GetChannelName() string`

GetChannelName returns the ChannelName field if non-nil, zero value otherwise.

### GetChannelNameOk

`func (o *InterStateChannelResult) GetChannelNameOk() (*string, bool)`

GetChannelNameOk returns a tuple with the ChannelName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelName

`func (o *InterStateChannelResult) SetChannelName(v string)`

SetChannelName sets ChannelName field to given value.


### GetValue

`func (o *InterStateChannelResult) GetValue() EncodedObject`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *InterStateChannelResult) GetValueOk() (*EncodedObject, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *InterStateChannelResult) SetValue(v EncodedObject)`

SetValue sets Value field to given value.

### HasValue

`func (o *InterStateChannelResult) HasValue() bool`

HasValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


