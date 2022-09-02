# EncodedObject

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Encoding** | Pointer to **string** |  | [optional] 
**Data** | Pointer to **string** |  | [optional] 

## Methods

### NewEncodedObject

`func NewEncodedObject() *EncodedObject`

NewEncodedObject instantiates a new EncodedObject object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEncodedObjectWithDefaults

`func NewEncodedObjectWithDefaults() *EncodedObject`

NewEncodedObjectWithDefaults instantiates a new EncodedObject object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEncoding

`func (o *EncodedObject) GetEncoding() string`

GetEncoding returns the Encoding field if non-nil, zero value otherwise.

### GetEncodingOk

`func (o *EncodedObject) GetEncodingOk() (*string, bool)`

GetEncodingOk returns a tuple with the Encoding field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncoding

`func (o *EncodedObject) SetEncoding(v string)`

SetEncoding sets Encoding field to given value.

### HasEncoding

`func (o *EncodedObject) HasEncoding() bool`

HasEncoding returns a boolean if a field has been set.

### GetData

`func (o *EncodedObject) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *EncodedObject) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *EncodedObject) SetData(v string)`

SetData sets Data field to given value.

### HasData

`func (o *EncodedObject) HasData() bool`

HasData returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


