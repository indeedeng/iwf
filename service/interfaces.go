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
		RunId            string
		StartedTimestamp int64
	}

	GetDataObjectsQueryRequest struct {
		Keys []string
	}

	GetDataObjectsQueryResponse struct {
		DataObjects []iwfidl.KeyValue
	}
)
