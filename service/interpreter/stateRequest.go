package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateRequest struct {
	stateStartRequest  iwfidl.StateMovement
	isResumeRequest    bool
	stateResumeRequest service.StateExecutionResumeInfo
}

func NewStateStartRequest(movement iwfidl.StateMovement) StateRequest {
	return StateRequest{
		stateStartRequest: movement,
	}
}

func NewStateResumeRequest(resumeRequest service.StateExecutionResumeInfo) StateRequest {
	return StateRequest{
		stateResumeRequest: resumeRequest,
		isResumeRequest:    true,
	}
}

func (sq StateRequest) GetStateStartRequest() iwfidl.StateMovement {
	return sq.stateStartRequest
}

func (sq StateRequest) GetStateResumeRequest() service.StateExecutionResumeInfo {
	return sq.stateResumeRequest
}

func (sq StateRequest) IsResumeRequest() bool {
	return sq.isResumeRequest
}

func (sq StateRequest) GetStateMovement() iwfidl.StateMovement {
	if sq.isResumeRequest {
		return sq.stateResumeRequest.State
	} else {
		return sq.stateStartRequest
	}
}

func (sq StateRequest) GetStateId() string {
	if sq.IsResumeRequest() {
		return sq.stateResumeRequest.State.StateId
	}
	return sq.stateStartRequest.StateId
}
