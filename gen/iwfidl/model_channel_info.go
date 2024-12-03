/*
Workflow APIs

This APIs for iwf SDKs to operate workflows

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package iwfidl

import (
	"encoding/json"
)

// checks if the ChannelInfo type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ChannelInfo{}

// ChannelInfo struct for ChannelInfo
type ChannelInfo struct {
	Name *string `json:"name,omitempty"`
	Size *int32  `json:"size,omitempty"`
}

// NewChannelInfo instantiates a new ChannelInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewChannelInfo() *ChannelInfo {
	this := ChannelInfo{}
	return &this
}

// NewChannelInfoWithDefaults instantiates a new ChannelInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewChannelInfoWithDefaults() *ChannelInfo {
	this := ChannelInfo{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ChannelInfo) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ChannelInfo) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ChannelInfo) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ChannelInfo) SetName(v string) {
	o.Name = &v
}

// GetSize returns the Size field value if set, zero value otherwise.
func (o *ChannelInfo) GetSize() int32 {
	if o == nil || IsNil(o.Size) {
		var ret int32
		return ret
	}
	return *o.Size
}

// GetSizeOk returns a tuple with the Size field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ChannelInfo) GetSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.Size) {
		return nil, false
	}
	return o.Size, true
}

// HasSize returns a boolean if a field has been set.
func (o *ChannelInfo) HasSize() bool {
	if o != nil && !IsNil(o.Size) {
		return true
	}

	return false
}

// SetSize gets a reference to the given int32 and assigns it to the Size field.
func (o *ChannelInfo) SetSize(v int32) {
	o.Size = &v
}

func (o ChannelInfo) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ChannelInfo) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Size) {
		toSerialize["size"] = o.Size
	}
	return toSerialize, nil
}

type NullableChannelInfo struct {
	value *ChannelInfo
	isSet bool
}

func (v NullableChannelInfo) Get() *ChannelInfo {
	return v.value
}

func (v *NullableChannelInfo) Set(val *ChannelInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableChannelInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableChannelInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableChannelInfo(val *ChannelInfo) *NullableChannelInfo {
	return &NullableChannelInfo{value: val, isSet: true}
}

func (v NullableChannelInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableChannelInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
