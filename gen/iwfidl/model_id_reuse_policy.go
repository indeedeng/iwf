/*
Workflow APIs

This APIs for iwf SDKs to operate workflows

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package iwfidl

import (
	"encoding/json"
	"fmt"
)

// IDReusePolicy the model 'IDReusePolicy'
type IDReusePolicy string

// List of IDReusePolicy
const (
	ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY IDReusePolicy = "ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY"
	ALLOW_IF_NO_RUNNING                 IDReusePolicy = "ALLOW_IF_NO_RUNNING"
	DISALLOW_REUSE                      IDReusePolicy = "DISALLOW_REUSE"
	ALLOW_TERMINATE_IF_RUNNING          IDReusePolicy = "ALLOW_TERMINATE_IF_RUNNING"
)

// All allowed values of IDReusePolicy enum
var AllowedIDReusePolicyEnumValues = []IDReusePolicy{
	"ALLOW_IF_PREVIOUS_EXISTS_ABNORMALLY",
	"ALLOW_IF_NO_RUNNING",
	"DISALLOW_REUSE",
	"ALLOW_TERMINATE_IF_RUNNING",
}

func (v *IDReusePolicy) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := IDReusePolicy(value)
	for _, existing := range AllowedIDReusePolicyEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid IDReusePolicy", value)
}

// NewIDReusePolicyFromValue returns a pointer to a valid IDReusePolicy
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewIDReusePolicyFromValue(v string) (*IDReusePolicy, error) {
	ev := IDReusePolicy(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for IDReusePolicy: valid values are %v", v, AllowedIDReusePolicyEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v IDReusePolicy) IsValid() bool {
	for _, existing := range AllowedIDReusePolicyEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to IDReusePolicy value
func (v IDReusePolicy) Ptr() *IDReusePolicy {
	return &v
}

type NullableIDReusePolicy struct {
	value *IDReusePolicy
	isSet bool
}

func (v NullableIDReusePolicy) Get() *IDReusePolicy {
	return v.value
}

func (v *NullableIDReusePolicy) Set(val *IDReusePolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableIDReusePolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableIDReusePolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIDReusePolicy(val *IDReusePolicy) *NullableIDReusePolicy {
	return &NullableIDReusePolicy{value: val, isSet: true}
}

func (v NullableIDReusePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIDReusePolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
