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

		Config WorkflowConfig `json:"config,omitempty"`

		// ContinueAsNew indicate this is input for continueAsNew, when true, all fields can be blank except for ContinueAsNewInput
		ContinueAsNew bool `json:"continueAsNew"`

		ContinueAsNewInput ContinueAsNewInput `json:"continueAsNewInput"`
	}

	ContinueAsNewInput struct {
		Config WorkflowConfig `json:"config,omitempty"`

		IwfWorkflowExecution IwfWorkflowExecution `json:"iwfWorkflowExecution"`
	}

	WorkflowConfig struct {
		DisableSystemSearchAttributes                bool `json:"disableSystemSearchAttributes,omitempty"`
		ContinueAsNewThresholdExecutedStateExecution int  `json:"continueAsNewThresholdExecutedStateExecution"`
		ContinueAsNewThresholdSignalsReceived        int  `json:"continueAsNewThresholdSignalsReceived"`
	}

	InterpreterWorkflowOutput struct {
		StateCompletionOutputs []iwfidl.StateCompletionOutput `json:"stateCompletionOutputs,omitempty"`
	}

	StateStartActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateStartRequest
	}

	StateDecideActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateDecideRequest
	}

	IwfWorkflowExecution struct {
		IwfWorkerUrl     string
		WorkflowType     string
		WorkflowId       string
		RunId            string // this is the first runId including reset & continueAsNew
		StartedTimestamp int64
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

	DumpAllInternalResponse struct {
		InterStateChannelReceived               map[string][]*iwfidl.EncodedObject
		SignalChannelReceived                   map[string][]*iwfidl.EncodedObject
		StateExecutionCounterInfo               StateExecutionCounterInfo
		PendingStateExecutionsCompletedCommands map[string]PendingStateExecutionCompletedCommands
		PendingStateExecutionsRequestCommands   map[string]PendingStateExecutionRequestCommands
		DataObjects                             []iwfidl.KeyValue
		SearchAttributes                        []iwfidl.SearchAttribute
	}

	StateExecutionCounterInfo struct {
		ExecutedStateIdCount      map[string]int // for stateExecutionId
		PendingStateIdCount       map[string]int // for sys search attribute
		TotalPendingStateExeCount int            // for "dead end"
	}

	PendingStateExecution struct {
		State                iwfidl.StateMovement
		DeciderTriggerType   iwfidl.DeciderTriggerType
		StateExecutionStatus StateExecutionStatus
	}

	PendingStateExecutionRequestCommands struct {
		TimerCommands             []iwfidl.TimerCommand
		SignalCommands            []iwfidl.SignalCommand
		InterStateChannelCommands []iwfidl.InterStateChannelCommand
	}

	PendingStateExecutionCompletedCommands struct {
		CompletedTimerCommands             map[int]bool
		CompletedSignalCommands            map[int]*iwfidl.EncodedObject
		CompletedInterStateChannelCommands map[int]*iwfidl.EncodedObject
	}
)

type StateExecutionStatus string

const FailureStateExecutionStatus StateExecutionStatus = "FailureStateExecutionStatus"                       // this will fail the workflow, no continueAsNew
const StartApiAbortedStateExecutionStatus StateExecutionStatus = "StartApiAbortedStateExecutionStatus"       // this will put the state into a special pending queue for continueAsNew from startApi
const WaitingCommandsStateExecutionStatus StateExecutionStatus = "WaitingCommandsStateExecutionStatus"       // this will put the state into a special pending queue for continueAsNew from waiting command
const WaitingCompletedStateExecutionStatus StateExecutionStatus = "WaitingCompletedStateExecutionStatus"     // this will put the state into a special pending queue for continueAsNew from decideApi
const DecideApiCompletedStateExecutionStatus StateExecutionStatus = "DecideApiCompletedStateExecutionStatus" // this will process as normal

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
