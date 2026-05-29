package interpreter

import (
	"testing"

	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/stretchr/testify/assert"
)

func TestGetChannelCommandLimits(t *testing.T) {
	versioner := &GlobalVersioner{version: StartingVersionChannelConsumeN}

	tests := []struct {
		name          string
		cmd           iwfidl.InterStateChannelCommand
		wantAtLeast   int
		wantAtMost    int
		legacyVersion bool
	}{
		{
			name:        "default single",
			cmd:         iwfidl.InterStateChannelCommand{ChannelName: "ch"},
			wantAtLeast: 1,
			wantAtMost:  1,
		},
		{
			name:        "exact n",
			cmd:         iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(3), AtMost: iwfidl.PtrInt32(3)},
			wantAtLeast: 3,
			wantAtMost:  3,
		},
		{
			name:        "one to all",
			cmd:         iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(1)},
			wantAtLeast: 1,
			wantAtMost:  int(^uint(0) >> 1),
		},
		{
			name:        "zero to all",
			cmd:         iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(0)},
			wantAtLeast: 0,
			wantAtMost:  int(^uint(0) >> 1),
		},
		{
			name:        "at most only means exact n",
			cmd:         iwfidl.InterStateChannelCommand{ChannelName: "ch", AtMost: iwfidl.PtrInt32(2)},
			wantAtLeast: 2,
			wantAtMost:  2,
		},
		{
			name:          "legacy ignores limits",
			cmd:           iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(3), AtMost: iwfidl.PtrInt32(3)},
			wantAtLeast:   1,
			wantAtMost:    1,
			legacyVersion: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := versioner
			if tt.legacyVersion {
				v = &GlobalVersioner{version: StartingVersionWaitingCommandThreads}
			}
			atLeast, atMost := getChannelCommandLimits(tt.cmd, v)
			assert.Equal(t, tt.wantAtLeast, atLeast)
			assert.Equal(t, tt.wantAtMost, atMost)
		})
	}
}

func TestValidateChannelCommandLimits(t *testing.T) {
	tests := []struct {
		name    string
		cmd     iwfidl.InterStateChannelCommand
		wantErr string
	}{
		{
			name: "valid unset",
			cmd:  iwfidl.InterStateChannelCommand{ChannelName: "ch"},
		},
		{
			name: "valid zero",
			cmd:  iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(0)},
		},
		{
			name:    "negative at least",
			cmd:     iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(-1)},
			wantErr: "InterStateChannelCommand atLeast cannot be negative",
		},
		{
			name:    "negative at most",
			cmd:     iwfidl.InterStateChannelCommand{ChannelName: "ch", AtMost: iwfidl.PtrInt32(-1)},
			wantErr: "InterStateChannelCommand atMost cannot be negative",
		},
		{
			name:    "at most less than at least",
			cmd:     iwfidl.InterStateChannelCommand{ChannelName: "ch", AtLeast: iwfidl.PtrInt32(3), AtMost: iwfidl.PtrInt32(2)},
			wantErr: "InterStateChannelCommand atMost cannot be less than atLeast",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChannelCommandLimits(tt.cmd)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
