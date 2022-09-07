package service

import (
	"github.com/cadence-oss/iwf-server/gen/server/workflow"
)

type InterpreterWorkflowInput struct {
	IwfWorkflowType string `json:"iwfWorkflowType,omitempty"`

	IwfWorkerUrl string `json:"iwfWorkerUrl,omitempty"`

	StartStateId string `json:"startStateId,omitempty"`

	StateInput workflow.EncodedObject `json:"stateInput,omitempty"`

	StateOptions workflow.WorkflowStateOptions `json:"stateOptions,omitempty"`
}

type InterpreterWorkflowOutput struct {
	CompletedStateId string                 `json:"completedStateId,omitempty"`
	StateOutput      workflow.EncodedObject `json:"stateOutput,omitempty"`
}
