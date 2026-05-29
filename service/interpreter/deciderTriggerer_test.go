package interpreter

import (
	"testing"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
)

func TestIsDeciderTriggerConditionMetAllWithMultiInterStateChannel(t *testing.T) {
	req := iwfidl.CommandRequest{
		DeciderTriggerType: iwfidl.ALL_COMMAND_COMPLETED.Ptr(),
		InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
			{CommandId: ptr.Any("cmd-1"), ChannelName: "ch"},
			{CommandId: ptr.Any("cmd-2"), ChannelName: "ch"},
		},
	}

	assert.False(t, IsDeciderTriggerConditionMet(
		req, nil, nil,
		map[int]*iwfidl.EncodedObject{0: testEncodedObject("1")},
		map[int][]*iwfidl.EncodedObject{},
	))
	assert.True(t, IsDeciderTriggerConditionMet(
		req, nil, nil,
		map[int]*iwfidl.EncodedObject{0: testEncodedObject("1")},
		map[int][]*iwfidl.EncodedObject{1: {testEncodedObject("2"), testEncodedObject("3")}},
	))
}

func TestIsDeciderTriggerConditionMetAnyWithMultiInterStateChannel(t *testing.T) {
	req := iwfidl.CommandRequest{
		DeciderTriggerType: iwfidl.ANY_COMMAND_COMPLETED.Ptr(),
		TimerCommands: []iwfidl.TimerCommand{
			{CommandId: ptr.Any("timer-1"), DurationSeconds: iwfidl.PtrInt64(10)},
		},
		InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
			{CommandId: ptr.Any("cmd-1"), ChannelName: "ch"},
		},
	}

	assert.True(t, IsDeciderTriggerConditionMet(
		req,
		map[int]service.InternalTimerStatus{},
		nil,
		nil,
		map[int][]*iwfidl.EncodedObject{0: {testEncodedObject("1")}},
	))
}

func TestIsDeciderTriggerConditionMetAnyCombinationWithMultiInterStateChannel(t *testing.T) {
	req := iwfidl.CommandRequest{
		DeciderTriggerType: iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
		TimerCommands: []iwfidl.TimerCommand{
			{CommandId: ptr.Any("timer-1"), DurationSeconds: iwfidl.PtrInt64(10)},
		},
		InterStateChannelCommands: []iwfidl.InterStateChannelCommand{
			{CommandId: ptr.Any("cmd-1"), ChannelName: "ch"},
		},
		CommandCombinations: []iwfidl.CommandCombination{
			{CommandIds: []string{"timer-1"}},
			{CommandIds: []string{"cmd-1"}},
		},
	}

	assert.True(t, IsDeciderTriggerConditionMet(
		req,
		map[int]service.InternalTimerStatus{},
		nil,
		nil,
		map[int][]*iwfidl.EncodedObject{0: {testEncodedObject("1")}},
	))
}
