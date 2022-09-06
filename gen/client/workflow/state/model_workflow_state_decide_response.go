/*
WorkflowState APIs

This APIs for iwf-server to invoke user workflow code defined in WorkflowState using any iwf SDKs

API version: 1.0.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package state

import (
	"encoding/json"
)

// WorkflowStateDecideResponse struct for WorkflowStateDecideResponse
type WorkflowStateDecideResponse struct {
	StateDecision []StateDecision `json:"stateDecision,omitempty"`
}

// NewWorkflowStateDecideResponse instantiates a new WorkflowStateDecideResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWorkflowStateDecideResponse() *WorkflowStateDecideResponse {
	this := WorkflowStateDecideResponse{}
	return &this
}

// NewWorkflowStateDecideResponseWithDefaults instantiates a new WorkflowStateDecideResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWorkflowStateDecideResponseWithDefaults() *WorkflowStateDecideResponse {
	this := WorkflowStateDecideResponse{}
	return &this
}

// GetStateDecision returns the StateDecision field value if set, zero value otherwise.
func (o *WorkflowStateDecideResponse) GetStateDecision() []StateDecision {
	if o == nil || o.StateDecision == nil {
		var ret []StateDecision
		return ret
	}
	return o.StateDecision
}

// GetStateDecisionOk returns a tuple with the StateDecision field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WorkflowStateDecideResponse) GetStateDecisionOk() ([]StateDecision, bool) {
	if o == nil || o.StateDecision == nil {
		return nil, false
	}
	return o.StateDecision, true
}

// HasStateDecision returns a boolean if a field has been set.
func (o *WorkflowStateDecideResponse) HasStateDecision() bool {
	if o != nil && o.StateDecision != nil {
		return true
	}

	return false
}

// SetStateDecision gets a reference to the given []StateDecision and assigns it to the StateDecision field.
func (o *WorkflowStateDecideResponse) SetStateDecision(v []StateDecision) {
	o.StateDecision = v
}

func (o WorkflowStateDecideResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.StateDecision != nil {
		toSerialize["stateDecision"] = o.StateDecision
	}
	return json.Marshal(toSerialize)
}

type NullableWorkflowStateDecideResponse struct {
	value *WorkflowStateDecideResponse
	isSet bool
}

func (v NullableWorkflowStateDecideResponse) Get() *WorkflowStateDecideResponse {
	return v.value
}

func (v *NullableWorkflowStateDecideResponse) Set(val *WorkflowStateDecideResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableWorkflowStateDecideResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableWorkflowStateDecideResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWorkflowStateDecideResponse(val *WorkflowStateDecideResponse) *NullableWorkflowStateDecideResponse {
	return &NullableWorkflowStateDecideResponse{value: val, isSet: true}
}

func (v NullableWorkflowStateDecideResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWorkflowStateDecideResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

