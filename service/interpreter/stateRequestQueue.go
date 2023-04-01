package interpreter

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"sort"
)

type StateRequestQueue struct {
	queue []StateRequest
}

func NewStateRequestQueue(initReq iwfidl.StateMovement) *StateRequestQueue {
	return &StateRequestQueue{
		queue: []StateRequest{
			NewStateRequest(initReq),
		},
	}
}

func NewStateRequestQueueForContinueAsNew(initReqs []iwfidl.StateMovement, initResumeReqs map[string]service.PendingStateExecution) *StateRequestQueue {
	var queue []StateRequest
	for _, r := range initReqs {
		queue = append(queue, NewStateRequest(r))
	}

	var keys []string
	for k := range initResumeReqs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		r := initResumeReqs[k]
		queue = append(queue, NewPendingStateExecutionRequest(r))
	}

	return &StateRequestQueue{}
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

func (srq *StateRequestQueue) GetAllNonStartedRequest() []iwfidl.StateMovement {
	var res []iwfidl.StateMovement
	for _, r := range srq.queue {
		if r.IsPendingFromContinueAsNew() {
			continue
		}
		res = append(res, r.GetNewRequest())
	}
	return res
}

func (srq *StateRequestQueue) AddNewRequests(reqs []iwfidl.StateMovement) {
	for _, r := range reqs {
		srq.queue = append(srq.queue, NewStateRequest(r))
	}
}
