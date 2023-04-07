package service

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
)

type (
	InterpreterWorkflowInput struct {
		IwfWorkflowType string `json:"iwfWorkflowType,omitempty"`

		IwfWorkerUrl string `json:"iwfWorkerUrl,omitempty"`

		StartStateId string `json:"startStateId,omitempty"`

		StateInput iwfidl.EncodedObject `json:"stateInput,omitempty"`

		StateOptions iwfidl.WorkflowStateOptions `json:"stateOptions,omitempty"`

		InitSearchAttributes []iwfidl.SearchAttribute `json:"initSearchAttributes,omitempty"`

		Config iwfidl.WorkflowConfig `json:"config,omitempty"`

		// IsResumeFromContinueAsNew indicate this is input for continueAsNew
		// when true, will ignore StartStateId, StateInput, StateOptions, InitSearchAttributes
		IsResumeFromContinueAsNew bool `json:"isResumeFromContinueAsNew"`

		ContinueAsNewInput ContinueAsNewInput `json:"continueAsNewInput"`
	}

	ContinueAsNewInput struct {
		PreviousInternalRunId string `json:"previousInternalRunId"` // for loading from previous run
	}

	InterpreterWorkflowOutput struct {
		StateCompletionOutputs []iwfidl.StateCompletionOutput `json:"stateCompletionOutputs,omitempty"`
	}

	BasicInfo struct {
		IwfWorkflowType string `json:"iwfWorkflowType,omitempty"`

		IwfWorkerUrl string `json:"iwfWorkerUrl,omitempty"`
	}

	StateStartActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateStartRequest
	}

	StateDecideActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateDecideRequest
	}

	GetDataObjectsQueryRequest struct {
		Keys []string
	}

	GetDataObjectsQueryResponse struct {
		DataObjects []iwfidl.KeyValue
	}

	GetCurrentTimerInfosQueryResponse struct {
		StateExecutionCurrentTimerInfos map[string][]*TimerInfo // key is stateExecutionId
	}

	TimerInfo struct {
		CommandId                  string
		FiringUnixTimestampSeconds int64
		Status                     InternalTimerStatus
	}

	SkipTimerSignalRequest struct {
		StateExecutionId string
		CommandId        string
		CommandIndex     int
	}

	FailWorkflowSignalRequest struct {
		Reason string
	}

	InternalTimerStatus string

	ContinueAsNewDumpResponse struct {
		StatesToStartFromBeginning []iwfidl.StateMovement              // StatesToStartFromBeginning means they haven't started in the previous run
		StateExecutionsToResume    map[string]StateExecutionResumeInfo // stateExeId to StateExecutionResumeInfo
		InterStateChannelReceived  map[string][]*iwfidl.EncodedObject
		SignalsReceived            map[string][]*iwfidl.EncodedObject
		StateExecutionCounterInfo  StateExecutionCounterInfo
		StateOutputs               []iwfidl.StateCompletionOutput
		StaleSkipTimerSignals      []StaleSkipTimerSignal

		DataObjects      []iwfidl.KeyValue
		SearchAttributes []iwfidl.SearchAttribute
	}

	DebugDumpResponse struct {
		Config iwfidl.WorkflowConfig
	}

	StateExecutionCounterInfo struct {
		StateIdStartedCount            map[string]int // for stateExecutionId
		StateIdCurrentlyExecutingCount map[string]int // for sys search attribute ExecutingStateIds
		TotalCurrentlyExecutingCount   int            // for "dead end"
	}

	StateExecutionResumeInfo struct {
		StateExecutionId                string                          `json:"stateExecutionId"`
		State                           iwfidl.StateMovement            `json:"state"`
		StateExecutionCompletedCommands StateExecutionCompletedCommands `json:"stateExecutionCompletedCommands"`
		CommandRequest                  iwfidl.CommandRequest           `json:"commandRequest"`
		StateExecutionLocals            []iwfidl.KeyValue               `json:"stateExecutionLocals"`
	}

	StateExecutionCompletedCommands struct {
		CompletedTimerCommands             map[int]InternalTimerStatus   `json:"completedTimerCommands"`
		CompletedSignalCommands            map[int]*iwfidl.EncodedObject `json:"completedSignalCommands"`
		CompletedInterStateChannelCommands map[int]*iwfidl.EncodedObject `json:"completedInterStateChannelCommands"`
	}

	StaleSkipTimerSignal struct {
		StateExecutionId  string
		TimerCommandId    string
		TimerCommandIndex int
	}
)

type StateExecutionStatus string

const FailureStateExecutionStatus StateExecutionStatus = "Failure"
const WaitingCommandsStateExecutionStatus StateExecutionStatus = "WaitingCommands" // this will put the state into a special pending queue for continueAsNew from waiting command
const CompletedStateExecutionStatus StateExecutionStatus = "Completed"             // this will process as normal

const (
	TimerPending InternalTimerStatus = "Pending"
	TimerFired   InternalTimerStatus = "Fired"
	TimerSkipped InternalTimerStatus = "Skipped"
)

// ValidateTimerSkipRequest validates if the skip timer request is valid
// return true if it's valid, along with the timer pointer
// use timerIdx if timerId is not empty
func ValidateTimerSkipRequest(stateExeTimerInfos map[string][]*TimerInfo, stateExeId, timerId string, timerIdx int) (*TimerInfo, bool) {
	timerInfos := stateExeTimerInfos[stateExeId]
	if len(timerInfos) == 0 {
		return nil, false
	}
	if timerId != "" {
		for _, t := range timerInfos {
			if t.CommandId == timerId {
				return t, true
			}
		}
		return nil, false
	}
	if timerIdx >= 0 && timerIdx < len(timerInfos) {
		t := timerInfos[timerIdx]
		if t.Status == TimerPending {
			return t, true
		}
	}
	return nil, false
}
