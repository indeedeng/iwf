package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
)

type StateRequest struct {
	newRequest               iwfidl.StateMovement
	pendingFromContinueAsNew bool
	pendingRequest           service.PendingStateExecution
}

func NewStateRequest(movement iwfidl.StateMovement) StateRequest {
	return StateRequest{
		newRequest: movement,
	}
}

func NewPendingStateExecutionRequest(pendingRequest service.PendingStateExecution) StateRequest {
	return StateRequest{
		pendingRequest:           pendingRequest,
		pendingFromContinueAsNew: true,
	}
}

func (sq StateRequest) GetNewRequest() iwfidl.StateMovement {
	return sq.newRequest
}

func (sq StateRequest) GetPendingRequest() service.PendingStateExecution {
	return sq.pendingRequest
}

func (sq StateRequest) IsPendingFromContinueAsNew() bool {
	return sq.pendingFromContinueAsNew
}

func (sq StateRequest) GetStateId() string {
	if sq.IsPendingFromContinueAsNew() {
		return sq.pendingRequest.State.StateId
	}
	return sq.newRequest.StateId
}
