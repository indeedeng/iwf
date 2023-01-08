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

// PersistenceLoadingPolicy struct for PersistenceLoadingPolicy
type PersistenceLoadingPolicy struct {
	PersistenceLoadingType *PersistenceLoadingType `json:"persistenceLoadingType,omitempty"`
	PartialLoadingKeys []string `json:"partialLoadingKeys,omitempty"`
}

// NewPersistenceLoadingPolicy instantiates a new PersistenceLoadingPolicy object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPersistenceLoadingPolicy() *PersistenceLoadingPolicy {
	this := PersistenceLoadingPolicy{}
	return &this
}

// NewPersistenceLoadingPolicyWithDefaults instantiates a new PersistenceLoadingPolicy object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPersistenceLoadingPolicyWithDefaults() *PersistenceLoadingPolicy {
	this := PersistenceLoadingPolicy{}
	return &this
}

// GetPersistenceLoadingType returns the PersistenceLoadingType field value if set, zero value otherwise.
func (o *PersistenceLoadingPolicy) GetPersistenceLoadingType() PersistenceLoadingType {
	if o == nil || isNil(o.PersistenceLoadingType) {
		var ret PersistenceLoadingType
		return ret
	}
	return *o.PersistenceLoadingType
}

// GetPersistenceLoadingTypeOk returns a tuple with the PersistenceLoadingType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetPersistenceLoadingTypeOk() (*PersistenceLoadingType, bool) {
	if o == nil || isNil(o.PersistenceLoadingType) {
    return nil, false
	}
	return o.PersistenceLoadingType, true
}

// HasPersistenceLoadingType returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasPersistenceLoadingType() bool {
	if o != nil && !isNil(o.PersistenceLoadingType) {
		return true
	}

	return false
}

// SetPersistenceLoadingType gets a reference to the given PersistenceLoadingType and assigns it to the PersistenceLoadingType field.
func (o *PersistenceLoadingPolicy) SetPersistenceLoadingType(v PersistenceLoadingType) {
	o.PersistenceLoadingType = &v
}

// GetPartialLoadingKeys returns the PartialLoadingKeys field value if set, zero value otherwise.
func (o *PersistenceLoadingPolicy) GetPartialLoadingKeys() []string {
	if o == nil || isNil(o.PartialLoadingKeys) {
		var ret []string
		return ret
	}
	return o.PartialLoadingKeys
}

// GetPartialLoadingKeysOk returns a tuple with the PartialLoadingKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetPartialLoadingKeysOk() ([]string, bool) {
	if o == nil || isNil(o.PartialLoadingKeys) {
    return nil, false
	}
	return o.PartialLoadingKeys, true
}

// HasPartialLoadingKeys returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasPartialLoadingKeys() bool {
	if o != nil && !isNil(o.PartialLoadingKeys) {
		return true
	}

	return false
}

// SetPartialLoadingKeys gets a reference to the given []string and assigns it to the PartialLoadingKeys field.
func (o *PersistenceLoadingPolicy) SetPartialLoadingKeys(v []string) {
	o.PartialLoadingKeys = v
}

func (o PersistenceLoadingPolicy) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if !isNil(o.PersistenceLoadingType) {
		toSerialize["persistenceLoadingType"] = o.PersistenceLoadingType
	}
	if !isNil(o.PartialLoadingKeys) {
		toSerialize["partialLoadingKeys"] = o.PartialLoadingKeys
	}
	return json.Marshal(toSerialize)
}

type NullablePersistenceLoadingPolicy struct {
	value *PersistenceLoadingPolicy
	isSet bool
}

func (v NullablePersistenceLoadingPolicy) Get() *PersistenceLoadingPolicy {
	return v.value
}

func (v *NullablePersistenceLoadingPolicy) Set(val *PersistenceLoadingPolicy) {
	v.value = val
	v.isSet = true
}

func (v NullablePersistenceLoadingPolicy) IsSet() bool {
	return v.isSet
}

func (v *NullablePersistenceLoadingPolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePersistenceLoadingPolicy(val *PersistenceLoadingPolicy) *NullablePersistenceLoadingPolicy {
	return &NullablePersistenceLoadingPolicy{value: val, isSet: true}
}

func (v NullablePersistenceLoadingPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePersistenceLoadingPolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


