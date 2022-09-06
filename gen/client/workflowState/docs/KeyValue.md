# KeyValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | Pointer to **string** |  | [optional] 
**Value** | Pointer to [**EncodedObject**](EncodedObject.md) |  | [optional] 

## Methods

### NewKeyValue

`func NewKeyValue() *KeyValue`

NewKeyValue instantiates a new KeyValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewKeyValueWithDefaults

`func NewKeyValueWithDefaults() *KeyValue`

NewKeyValueWithDefaults instantiates a new KeyValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *KeyValue) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *KeyValue) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *KeyValue) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *KeyValue) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetValue

`func (o *KeyValue) GetValue() EncodedObject`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *KeyValue) GetValueOk() (*EncodedObject, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *KeyValue) SetValue(v EncodedObject)`

SetValue sets Value field to given value.

### HasValue

`func (o *KeyValue) HasValue() bool`

HasValue returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


