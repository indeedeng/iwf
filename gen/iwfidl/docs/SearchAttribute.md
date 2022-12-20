# SearchAttribute

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Key** | Pointer to **string** |  | [optional] 
**StringValue** | Pointer to **string** |  | [optional] 
**IntegerValue** | Pointer to **int64** |  | [optional] 
**ValueType** | Pointer to [**SearchAttributeValueType**](SearchAttributeValueType.md) |  | [optional] 

## Methods

### NewSearchAttribute

`func NewSearchAttribute() *SearchAttribute`

NewSearchAttribute instantiates a new SearchAttribute object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSearchAttributeWithDefaults

`func NewSearchAttributeWithDefaults() *SearchAttribute`

NewSearchAttributeWithDefaults instantiates a new SearchAttribute object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKey

`func (o *SearchAttribute) GetKey() string`

GetKey returns the Key field if non-nil, zero value otherwise.

### GetKeyOk

`func (o *SearchAttribute) GetKeyOk() (*string, bool)`

GetKeyOk returns a tuple with the Key field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKey

`func (o *SearchAttribute) SetKey(v string)`

SetKey sets Key field to given value.

### HasKey

`func (o *SearchAttribute) HasKey() bool`

HasKey returns a boolean if a field has been set.

### GetStringValue

`func (o *SearchAttribute) GetStringValue() string`

GetStringValue returns the StringValue field if non-nil, zero value otherwise.

### GetStringValueOk

`func (o *SearchAttribute) GetStringValueOk() (*string, bool)`

GetStringValueOk returns a tuple with the StringValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStringValue

`func (o *SearchAttribute) SetStringValue(v string)`

SetStringValue sets StringValue field to given value.

### HasStringValue

`func (o *SearchAttribute) HasStringValue() bool`

HasStringValue returns a boolean if a field has been set.

### GetIntegerValue

`func (o *SearchAttribute) GetIntegerValue() int64`

GetIntegerValue returns the IntegerValue field if non-nil, zero value otherwise.

### GetIntegerValueOk

`func (o *SearchAttribute) GetIntegerValueOk() (*int64, bool)`

GetIntegerValueOk returns a tuple with the IntegerValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntegerValue

`func (o *SearchAttribute) SetIntegerValue(v int64)`

SetIntegerValue sets IntegerValue field to given value.

### HasIntegerValue

`func (o *SearchAttribute) HasIntegerValue() bool`

HasIntegerValue returns a boolean if a field has been set.

### GetValueType

`func (o *SearchAttribute) GetValueType() SearchAttributeValueType`

GetValueType returns the ValueType field if non-nil, zero value otherwise.

### GetValueTypeOk

`func (o *SearchAttribute) GetValueTypeOk() (*SearchAttributeValueType, bool)`

GetValueTypeOk returns a tuple with the ValueType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValueType

`func (o *SearchAttribute) SetValueType(v SearchAttributeValueType)`

SetValueType sets ValueType field to given value.

### HasValueType

`func (o *SearchAttribute) HasValueType() bool`

HasValueType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


