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

// SignalCommand struct for SignalCommand
type SignalCommand struct {
	CommandId string `json:"commandId"`
	SignalChannelName string `json:"signalChannelName"`
}

// NewSignalCommand instantiates a new SignalCommand object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSignalCommand(commandId string, signalChannelName string) *SignalCommand {
	this := SignalCommand{}
	this.CommandId = commandId
	this.SignalChannelName = signalChannelName
	return &this
}

// NewSignalCommandWithDefaults instantiates a new SignalCommand object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSignalCommandWithDefaults() *SignalCommand {
	this := SignalCommand{}
	return &this
}

// GetCommandId returns the CommandId field value
func (o *SignalCommand) GetCommandId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CommandId
}

// GetCommandIdOk returns a tuple with the CommandId field value
// and a boolean to check if the value has been set.
func (o *SignalCommand) GetCommandIdOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.CommandId, true
}

// SetCommandId sets field value
func (o *SignalCommand) SetCommandId(v string) {
	o.CommandId = v
}

// GetSignalChannelName returns the SignalChannelName field value
func (o *SignalCommand) GetSignalChannelName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SignalChannelName
}

// GetSignalChannelNameOk returns a tuple with the SignalChannelName field value
// and a boolean to check if the value has been set.
func (o *SignalCommand) GetSignalChannelNameOk() (*string, bool) {
	if o == nil {
    return nil, false
	}
	return &o.SignalChannelName, true
}

// SetSignalChannelName sets field value
func (o *SignalCommand) SetSignalChannelName(v string) {
	o.SignalChannelName = v
}

func (o SignalCommand) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["commandId"] = o.CommandId
	}
	if true {
		toSerialize["signalChannelName"] = o.SignalChannelName
	}
	return json.Marshal(toSerialize)
}

type NullableSignalCommand struct {
	value *SignalCommand
	isSet bool
}

func (v NullableSignalCommand) Get() *SignalCommand {
	return v.value
}

func (v *NullableSignalCommand) Set(val *SignalCommand) {
	v.value = val
	v.isSet = true
}

func (v NullableSignalCommand) IsSet() bool {
	return v.isSet
}

func (v *NullableSignalCommand) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSignalCommand(val *SignalCommand) *NullableSignalCommand {
	return &NullableSignalCommand{value: val, isSet: true}
}

func (v NullableSignalCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSignalCommand) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


