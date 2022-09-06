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

// CommandRequest struct for CommandRequest
type CommandRequest struct {
	DeciderTriggerType *string `json:"deciderTriggerType,omitempty"`
	ActivityCommands []ActivityCommand `json:"activityCommands,omitempty"`
	TimerCommands []TimerCommand `json:"timerCommands,omitempty"`
	SignalCommands []SignalCommand `json:"signalCommands,omitempty"`
}

// NewCommandRequest instantiates a new CommandRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommandRequest() *CommandRequest {
	this := CommandRequest{}
	return &this
}

// NewCommandRequestWithDefaults instantiates a new CommandRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommandRequestWithDefaults() *CommandRequest {
	this := CommandRequest{}
	return &this
}

// GetDeciderTriggerType returns the DeciderTriggerType field value if set, zero value otherwise.
func (o *CommandRequest) GetDeciderTriggerType() string {
	if o == nil || o.DeciderTriggerType == nil {
		var ret string
		return ret
	}
	return *o.DeciderTriggerType
}

// GetDeciderTriggerTypeOk returns a tuple with the DeciderTriggerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommandRequest) GetDeciderTriggerTypeOk() (*string, bool) {
	if o == nil || o.DeciderTriggerType == nil {
		return nil, false
	}
	return o.DeciderTriggerType, true
}

// HasDeciderTriggerType returns a boolean if a field has been set.
func (o *CommandRequest) HasDeciderTriggerType() bool {
	if o != nil && o.DeciderTriggerType != nil {
		return true
	}

	return false
}

// SetDeciderTriggerType gets a reference to the given string and assigns it to the DeciderTriggerType field.
func (o *CommandRequest) SetDeciderTriggerType(v string) {
	o.DeciderTriggerType = &v
}

// GetActivityCommands returns the ActivityCommands field value if set, zero value otherwise.
func (o *CommandRequest) GetActivityCommands() []ActivityCommand {
	if o == nil || o.ActivityCommands == nil {
		var ret []ActivityCommand
		return ret
	}
	return o.ActivityCommands
}

// GetActivityCommandsOk returns a tuple with the ActivityCommands field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommandRequest) GetActivityCommandsOk() ([]ActivityCommand, bool) {
	if o == nil || o.ActivityCommands == nil {
		return nil, false
	}
	return o.ActivityCommands, true
}

// HasActivityCommands returns a boolean if a field has been set.
func (o *CommandRequest) HasActivityCommands() bool {
	if o != nil && o.ActivityCommands != nil {
		return true
	}

	return false
}

// SetActivityCommands gets a reference to the given []ActivityCommand and assigns it to the ActivityCommands field.
func (o *CommandRequest) SetActivityCommands(v []ActivityCommand) {
	o.ActivityCommands = v
}

// GetTimerCommands returns the TimerCommands field value if set, zero value otherwise.
func (o *CommandRequest) GetTimerCommands() []TimerCommand {
	if o == nil || o.TimerCommands == nil {
		var ret []TimerCommand
		return ret
	}
	return o.TimerCommands
}

// GetTimerCommandsOk returns a tuple with the TimerCommands field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommandRequest) GetTimerCommandsOk() ([]TimerCommand, bool) {
	if o == nil || o.TimerCommands == nil {
		return nil, false
	}
	return o.TimerCommands, true
}

// HasTimerCommands returns a boolean if a field has been set.
func (o *CommandRequest) HasTimerCommands() bool {
	if o != nil && o.TimerCommands != nil {
		return true
	}

	return false
}

// SetTimerCommands gets a reference to the given []TimerCommand and assigns it to the TimerCommands field.
func (o *CommandRequest) SetTimerCommands(v []TimerCommand) {
	o.TimerCommands = v
}

// GetSignalCommands returns the SignalCommands field value if set, zero value otherwise.
func (o *CommandRequest) GetSignalCommands() []SignalCommand {
	if o == nil || o.SignalCommands == nil {
		var ret []SignalCommand
		return ret
	}
	return o.SignalCommands
}

// GetSignalCommandsOk returns a tuple with the SignalCommands field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CommandRequest) GetSignalCommandsOk() ([]SignalCommand, bool) {
	if o == nil || o.SignalCommands == nil {
		return nil, false
	}
	return o.SignalCommands, true
}

// HasSignalCommands returns a boolean if a field has been set.
func (o *CommandRequest) HasSignalCommands() bool {
	if o != nil && o.SignalCommands != nil {
		return true
	}

	return false
}

// SetSignalCommands gets a reference to the given []SignalCommand and assigns it to the SignalCommands field.
func (o *CommandRequest) SetSignalCommands(v []SignalCommand) {
	o.SignalCommands = v
}

func (o CommandRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DeciderTriggerType != nil {
		toSerialize["deciderTriggerType"] = o.DeciderTriggerType
	}
	if o.ActivityCommands != nil {
		toSerialize["activityCommands"] = o.ActivityCommands
	}
	if o.TimerCommands != nil {
		toSerialize["timerCommands"] = o.TimerCommands
	}
	if o.SignalCommands != nil {
		toSerialize["signalCommands"] = o.SignalCommands
	}
	return json.Marshal(toSerialize)
}

type NullableCommandRequest struct {
	value *CommandRequest
	isSet bool
}

func (v NullableCommandRequest) Get() *CommandRequest {
	return v.value
}

func (v *NullableCommandRequest) Set(val *CommandRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCommandRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCommandRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommandRequest(val *CommandRequest) *NullableCommandRequest {
	return &NullableCommandRequest{value: val, isSet: true}
}

func (v NullableCommandRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommandRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

