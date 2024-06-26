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

// checks if the PersistenceLoadingPolicy type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PersistenceLoadingPolicy{}

// PersistenceLoadingPolicy struct for PersistenceLoadingPolicy
type PersistenceLoadingPolicy struct {
	PersistenceLoadingType *PersistenceLoadingType `json:"persistenceLoadingType,omitempty"`
	PartialLoadingKeys     []string                `json:"partialLoadingKeys,omitempty"`
	LockingKeys            []string                `json:"lockingKeys,omitempty"`
	UseKeyAsPrefix         *bool                   `json:"useKeyAsPrefix,omitempty"`
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
	if o == nil || IsNil(o.PersistenceLoadingType) {
		var ret PersistenceLoadingType
		return ret
	}
	return *o.PersistenceLoadingType
}

// GetPersistenceLoadingTypeOk returns a tuple with the PersistenceLoadingType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetPersistenceLoadingTypeOk() (*PersistenceLoadingType, bool) {
	if o == nil || IsNil(o.PersistenceLoadingType) {
		return nil, false
	}
	return o.PersistenceLoadingType, true
}

// HasPersistenceLoadingType returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasPersistenceLoadingType() bool {
	if o != nil && !IsNil(o.PersistenceLoadingType) {
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
	if o == nil || IsNil(o.PartialLoadingKeys) {
		var ret []string
		return ret
	}
	return o.PartialLoadingKeys
}

// GetPartialLoadingKeysOk returns a tuple with the PartialLoadingKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetPartialLoadingKeysOk() ([]string, bool) {
	if o == nil || IsNil(o.PartialLoadingKeys) {
		return nil, false
	}
	return o.PartialLoadingKeys, true
}

// HasPartialLoadingKeys returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasPartialLoadingKeys() bool {
	if o != nil && !IsNil(o.PartialLoadingKeys) {
		return true
	}

	return false
}

// SetPartialLoadingKeys gets a reference to the given []string and assigns it to the PartialLoadingKeys field.
func (o *PersistenceLoadingPolicy) SetPartialLoadingKeys(v []string) {
	o.PartialLoadingKeys = v
}

// GetLockingKeys returns the LockingKeys field value if set, zero value otherwise.
func (o *PersistenceLoadingPolicy) GetLockingKeys() []string {
	if o == nil || IsNil(o.LockingKeys) {
		var ret []string
		return ret
	}
	return o.LockingKeys
}

// GetLockingKeysOk returns a tuple with the LockingKeys field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetLockingKeysOk() ([]string, bool) {
	if o == nil || IsNil(o.LockingKeys) {
		return nil, false
	}
	return o.LockingKeys, true
}

// HasLockingKeys returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasLockingKeys() bool {
	if o != nil && !IsNil(o.LockingKeys) {
		return true
	}

	return false
}

// SetLockingKeys gets a reference to the given []string and assigns it to the LockingKeys field.
func (o *PersistenceLoadingPolicy) SetLockingKeys(v []string) {
	o.LockingKeys = v
}

// GetUseKeyAsPrefix returns the UseKeyAsPrefix field value if set, zero value otherwise.
func (o *PersistenceLoadingPolicy) GetUseKeyAsPrefix() bool {
	if o == nil || IsNil(o.UseKeyAsPrefix) {
		var ret bool
		return ret
	}
	return *o.UseKeyAsPrefix
}

// GetUseKeyAsPrefixOk returns a tuple with the UseKeyAsPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PersistenceLoadingPolicy) GetUseKeyAsPrefixOk() (*bool, bool) {
	if o == nil || IsNil(o.UseKeyAsPrefix) {
		return nil, false
	}
	return o.UseKeyAsPrefix, true
}

// HasUseKeyAsPrefix returns a boolean if a field has been set.
func (o *PersistenceLoadingPolicy) HasUseKeyAsPrefix() bool {
	if o != nil && !IsNil(o.UseKeyAsPrefix) {
		return true
	}

	return false
}

// SetUseKeyAsPrefix gets a reference to the given bool and assigns it to the UseKeyAsPrefix field.
func (o *PersistenceLoadingPolicy) SetUseKeyAsPrefix(v bool) {
	o.UseKeyAsPrefix = &v
}

func (o PersistenceLoadingPolicy) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PersistenceLoadingPolicy) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.PersistenceLoadingType) {
		toSerialize["persistenceLoadingType"] = o.PersistenceLoadingType
	}
	if !IsNil(o.PartialLoadingKeys) {
		toSerialize["partialLoadingKeys"] = o.PartialLoadingKeys
	}
	if !IsNil(o.LockingKeys) {
		toSerialize["lockingKeys"] = o.LockingKeys
	}
	if !IsNil(o.UseKeyAsPrefix) {
		toSerialize["useKeyAsPrefix"] = o.UseKeyAsPrefix
	}
	return toSerialize, nil
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
