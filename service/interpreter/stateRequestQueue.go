package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"sort"
)

type StateRequestQueue struct {
	queue []StateRequest
}

func NewStateRequestQueue() *StateRequestQueue {
	return &StateRequestQueue{}
}

func NewStateRequestQueueWithResumeRequests(startReqs []iwfidl.StateMovement, resumeReqs map[string]service.StateExecutionResumeInfo) *StateRequestQueue {
	var queue []StateRequest
	for _, r := range startReqs {
		queue = append(queue, NewStateStartRequest(r))
	}

	var keys []string
	for k := range resumeReqs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		r := resumeReqs[k]
		queue = append(queue, NewStateResumeRequest(r))
	}

	return &StateRequestQueue{
		queue: queue,
	}
}

func (srq *StateRequestQueue) IsEmpty() bool {
	return len(srq.queue) == 0
}

func (srq *StateRequestQueue) TakeAll() []StateRequest {
	// copy the whole slice(pointer)
	res := srq.queue
	//reset to empty slice since each iteration will process all current states in the queue
	srq.queue = nil
	return res
}

func (srq *StateRequestQueue) GetAllStateStartRequests() []iwfidl.StateMovement {
	var res []iwfidl.StateMovement
	for _, r := range srq.queue {
		if r.IsResumeRequest() {
			continue
		}
		res = append(res, r.GetStateStartRequest())
	}
	return res
}

func (srq *StateRequestQueue) GetAllStateResumeRequests() []service.StateExecutionResumeInfo {
	var res []service.StateExecutionResumeInfo
	for _, r := range srq.queue {
		if !r.IsResumeRequest() {
			continue
		}
		res = append(res, r.GetStateResumeRequest())
	}
	return res
}

func (srq *StateRequestQueue) AddStateStartRequests(reqs []iwfidl.StateMovement) {
	for _, r := range reqs {
		srq.queue = append(srq.queue, NewStateStartRequest(r))
	}
}

func (srq *StateRequestQueue) AddSingleStateStartRequest(stateId string, input *iwfidl.EncodedObject, options *iwfidl.WorkflowStateOptions) {
	srq.queue = append(srq.queue, NewStateStartRequest(
		iwfidl.StateMovement{
			StateId:      stateId,
			StateInput:   input,
			StateOptions: options,
		},
	))
}
