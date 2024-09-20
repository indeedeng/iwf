package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInvalidAnyCommandCombination(t *testing.T) {
	validTimerCommands, validSignalCommands, internalCommands := createCommands()

	resp := iwfidl.WorkflowStateStartResponse{
		CommandRequest: &iwfidl.CommandRequest{
			SignalCommands:            validSignalCommands,
			TimerCommands:             validTimerCommands,
			InterStateChannelCommands: internalCommands,
			DeciderTriggerType:        iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
			CommandCombinations: []iwfidl.CommandCombination{
				{
					CommandIds: []string{
						"timer-cmd1", "signal-cmd1",
					},
				},
				{
					CommandIds: []string{
						"timer-cmd1", "invalid",
					},
				},
			},
		},
	}

	err := checkCommandRequestFromWaitUntilResponse(&resp)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ANY_COMMAND_COMBINATION_COMPLETED can only be used when every command has an commandId that is found in TimerCommands, SignalCommands or InternalChannelCommand")
}

func TestValidAnyCommandCombination(t *testing.T) {
	validTimerCommands, validSignalCommands, internalCommands := createCommands()

	resp := iwfidl.WorkflowStateStartResponse{
		CommandRequest: &iwfidl.CommandRequest{
			SignalCommands:            validSignalCommands,
			TimerCommands:             validTimerCommands,
			InterStateChannelCommands: internalCommands,
			DeciderTriggerType:        iwfidl.ANY_COMMAND_COMBINATION_COMPLETED.Ptr(),
			CommandCombinations: []iwfidl.CommandCombination{
				{
					CommandIds: []string{
						"timer-cmd1", "signal-cmd1",
					},
				},
				{
					CommandIds: []string{
						"timer-cmd1", "internal-cmd1",
					},
				},
			},
		},
	}

	err := checkCommandRequestFromWaitUntilResponse(&resp)

	assert.NoError(t, err)
}

func createCommands() ([]iwfidl.TimerCommand, []iwfidl.SignalCommand, []iwfidl.InterStateChannelCommand) {
	validTimerCommands := []iwfidl.TimerCommand{
		{
			CommandId:                  ptr.Any("timer-cmd1"),
			FiringUnixTimestampSeconds: iwfidl.PtrInt64(time.Now().Unix() + 86400*365), // one year later
		},
	}
	validSignalCommands := []iwfidl.SignalCommand{
		{
			CommandId:         ptr.Any("signal-cmd1"),
			SignalChannelName: "test-signal-name1",
		},
	}
	internalCommands := []iwfidl.InterStateChannelCommand{
		{
			CommandId:   ptr.Any("internal-cmd1"),
			ChannelName: "test-internal-name1",
		},
	}
	return validTimerCommands, validSignalCommands, internalCommands
}
