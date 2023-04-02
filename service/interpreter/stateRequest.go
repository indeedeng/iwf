package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateRequest struct {
	newStateRequest           iwfidl.StateMovement
	isResumeFromContinueAsNew bool
	resumeStateRequest        service.StateExecutionResumeInfo
}

func CreateNewStateRequest(movement iwfidl.StateMovement) StateRequest {
	return StateRequest{
		newStateRequest: movement,
	}
}

func CreateResumeStateExecutionRequest(pendingRequest service.StateExecutionResumeInfo) StateRequest {
	return StateRequest{
		resumeStateRequest:        pendingRequest,
		isResumeFromContinueAsNew: true,
	}
}

func (sq StateRequest) GetNewStateRequest() iwfidl.StateMovement {
	return sq.newStateRequest
}

func (sq StateRequest) GetResumeStateRequest() service.StateExecutionResumeInfo {
	return sq.resumeStateRequest
}

func (sq StateRequest) IsResumeFromContinueAsNew() bool {
	return sq.isResumeFromContinueAsNew
}

func (sq StateRequest) GetStateId() string {
	if sq.IsResumeFromContinueAsNew() {
		return sq.resumeStateRequest.State.StateId
	}
	return sq.newStateRequest.StateId
}
