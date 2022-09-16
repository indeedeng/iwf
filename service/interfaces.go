package service

import (
	"github.com/cadence-oss/iwf-server/gen/iwfidl"
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
		CompletedStateId string               `json:"completedStateId,omitempty"`
		StateOutput      iwfidl.EncodedObject `json:"stateOutput,omitempty"`
	}

	StateStartActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateStartRequest
	}

	StateDecideActivityInput struct {
		IwfWorkerUrl string
		Request      iwfidl.WorkflowStateDecideRequest
	}
)