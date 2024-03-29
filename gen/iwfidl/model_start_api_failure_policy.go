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

// StartApiFailurePolicy the model 'StartApiFailurePolicy'
type StartApiFailurePolicy string

// List of StartApiFailurePolicy
const (
	FAIL_WORKFLOW_ON_START_API_FAILURE     StartApiFailurePolicy = "FAIL_WORKFLOW_ON_START_API_FAILURE"
	PROCEED_TO_DECIDE_ON_START_API_FAILURE StartApiFailurePolicy = "PROCEED_TO_DECIDE_ON_START_API_FAILURE"
)

// All allowed values of StartApiFailurePolicy enum
var AllowedStartApiFailurePolicyEnumValues = []StartApiFailurePolicy{
	"FAIL_WORKFLOW_ON_START_API_FAILURE",
	"PROCEED_TO_DECIDE_ON_START_API_FAILURE",
}

func (v *StartApiFailurePolicy) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := StartApiFailurePolicy(value)
	for _, existing := range AllowedStartApiFailurePolicyEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid StartApiFailurePolicy", value)
}

// NewStartApiFailurePolicyFromValue returns a pointer to a valid StartApiFailurePolicy
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewStartApiFailurePolicyFromValue(v string) (*StartApiFailurePolicy, error) {
	ev := StartApiFailurePolicy(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for StartApiFailurePolicy: valid values are %v", v, AllowedStartApiFailurePolicyEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v StartApiFailurePolicy) IsValid() bool {
	for _, existing := range AllowedStartApiFailurePolicyEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to StartApiFailurePolicy value
func (v StartApiFailurePolicy) Ptr() *StartApiFailurePolicy {
	return &v
}

type NullableStartApiFailurePolicy struct {
	value *StartApiFailurePolicy
	isSet bool
}

func (v NullableStartApiFailurePolicy) Get() *StartApiFailurePolicy {
	return v.value
}

func (v *NullableStartApiFailurePolicy) Set(val *StartApiFailurePolicy) {
	v.value = val
	v.isSet = true
}

func (v NullableStartApiFailurePolicy) IsSet() bool {
	return v.isSet
}

func (v *NullableStartApiFailurePolicy) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStartApiFailurePolicy(val *StartApiFailurePolicy) *NullableStartApiFailurePolicy {
	return &NullableStartApiFailurePolicy{value: val, isSet: true}
}

func (v NullableStartApiFailurePolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStartApiFailurePolicy) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
