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

		// ContinueAsNew indicate this is input for continueAsNew, when true, will ignore StartStateId, StateInput, StateOptions, InitSearchAttributes
		ContinueAsNew bool `json:"continueAsNew"`

		ContinueAsNewInput ContinueAsNewInput `json:"continueAsNewInput"`
	}

	ContinueAsNewInput struct {
		IwfWorkflowExecution  IwfWorkflowExecution `json:"iwfWorkflowExecution"`
		PreviousInternalRunId string               `json:"previousInternalRunId"` // for loading from previous run
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
		NonStartedStates                        []iwfidl.StateMovement  // NonStartedStates means they haven't started in the previous run
		PendingStateExecution                   []PendingStateExecution // PendingStateExecution means they have started in the previous run
		InterStateChannelReceived               map[string][]*iwfidl.EncodedObject
		SignalsReceived                         map[string][]*iwfidl.EncodedObject
		StateExecutionCounterInfo               StateExecutionCounterInfo
		PendingStateExecutionsCompletedCommands map[string]PendingStateExecutionCompletedCommands
		PendingStateExecutionsRequestCommands   map[string]PendingStateExecutionRequestCommands
		DataObjects                             []iwfidl.KeyValue
		SearchAttributes                        []iwfidl.SearchAttribute
	}

	DumpAllInternalWithPaginationRequest struct {
		// default to DefaultContinueAsNewPageSizeInBytes(1024 * 1024), means 1MB
		PageSizeInBytes int
		// default to zero, means the first page
		PageNum int
	}

	DumpAllInternalWithPaginationResponse struct {
		// start over if the checksum is not matched anymore
		Checksum   string `json:"checksum"`
		TotalPages int    `json:"totalPages"`
		// combine all the JsonData of all pages to deserialize into DumpAllInternalResponse
		JsonData string `json:"jsonData"`
	}

	StateExecutionCounterInfo struct {
		ExecutedStateIdCount      map[string]int // for stateExecutionId
		PendingStateIdCount       map[string]int // for sys search attribute
		TotalPendingStateExeCount int            // for "dead end"
	}

	PendingStateExecution struct {
		StateExecutionId     string
		State                iwfidl.StateMovement
		StateExecutionStatus StateExecutionStatus
	}

	PendingStateExecutionRequestCommands struct {
		TimerCommands             []iwfidl.TimerCommand
		SignalCommands            []iwfidl.SignalCommand
		InterStateChannelCommands []iwfidl.InterStateChannelCommand
		DeciderTriggerType        iwfidl.DeciderTriggerType
	}

	PendingStateExecutionCompletedCommands struct {
		CompletedTimerCommands             map[int]bool
		CompletedSignalCommands            map[int]*iwfidl.EncodedObject
		CompletedInterStateChannelCommands map[int]*iwfidl.EncodedObject
	}
)

type StateExecutionStatus struct {
	StateExecutionStatusEnum StateExecutionStatusEnum `json:"stateExecutionStatusEnum"`
	StateExecutionLocals     []iwfidl.KeyValue        `json:"stateExecutionLocals"`
}

type StateExecutionStatusEnum string

const WaitingCommandsStateExecutionStatus StateExecutionStatusEnum = "WaitingCommands" // this will put the state into a special pending queue for continueAsNew from waiting command
const CompletedStateExecutionStatus StateExecutionStatusEnum = "Completed"             // this will process as normal

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
