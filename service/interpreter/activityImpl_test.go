package interpreter

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
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

func TestComposeHttpError_LocalActivity_LongErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	longError := strings.Repeat("a", 1000)

	httpResp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(longError)),
	}
	errMsg := "original error message"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("1st-attempt-failure", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", httpResp.StatusCode, longError[:50]+"...", errors.New(errMsg[:5]+"..."))).Return(returnedError)

	err = composeHttpError(true, mockActivityProvider, err, httpResp, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}

func TestComposeHttpError_RegularActivity_LongErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	longError := strings.Repeat("a", 1000)

	httpResp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(longError)),
	}
	errMsg := "original error message which is very long like this"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("test-error-type", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", httpResp.StatusCode, longError[:500]+"...", errors.New(errMsg[:50]+"..."))).Return(returnedError)

	err = composeHttpError(false, mockActivityProvider, err, httpResp, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}

func TestComposeHttpError_LocalActivity_ShortErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	shortError := strings.Repeat("a", 40)

	httpResp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(shortError)),
	}
	errMsg := "OK"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("1st-attempt-failure", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", httpResp.StatusCode, shortError, errors.New(errMsg))).Return(returnedError)

	err = composeHttpError(true, mockActivityProvider, err, httpResp, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}

func TestComposeHttpError_RegularActivity_ShortErrorResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	shortError := strings.Repeat("a", 40)

	httpResp := &http.Response{
		StatusCode: 400,
		Body:       io.NopCloser(strings.NewReader(shortError)),
	}
	errMsg := "OK"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("test-error-type", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", httpResp.StatusCode, shortError, errors.New(errMsg))).Return(returnedError)

	err = composeHttpError(false, mockActivityProvider, err, httpResp, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}

func TestComposeHttpError_LocalActivity_NilResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	errMsg := "OK"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("1st-attempt-failure", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", 0, "None", errors.New(errMsg))).Return(returnedError)

	err = composeHttpError(true, mockActivityProvider, err, nil, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}

func TestComposeHttpError_RegularActivity_NilResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockActivityProvider := NewMockActivityProvider(ctrl)

	errMsg := "OK"
	err := errors.New(errMsg)

	returnedError := errors.New("test error msg")
	mockActivityProvider.EXPECT().NewApplicationError("test-error-type", fmt.Sprintf("statusCode: %v, responseBody: %v, errMsg: %v", 0, "None", errors.New(errMsg))).Return(returnedError)

	err = composeHttpError(false, mockActivityProvider, err, nil, "test-error-type")
	if err != nil {
		return
	}

	assert.Error(t, err)
	assert.Equal(t, returnedError, err)
}
