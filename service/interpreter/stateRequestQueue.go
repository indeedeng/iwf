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
			CreateNewStateRequest(initReq),
		},
	}
}

func NewStateRequestQueueWithResumeRequests(newReqs []iwfidl.StateMovement, resumeReqs map[string]service.StateExecutionResumeInfo) *StateRequestQueue {
	var queue []StateRequest
	for _, r := range newReqs {
		queue = append(queue, CreateNewStateRequest(r))
	}

	var keys []string
	for k := range resumeReqs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		r := resumeReqs[k]
		queue = append(queue, CreateResumeStateExecutionRequest(r))
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

func (srq *StateRequestQueue) GetAllNewStateRequests() []iwfidl.StateMovement {
	var res []iwfidl.StateMovement
	for _, r := range srq.queue {
		if r.IsResumeFromContinueAsNew() {
			continue
		}
		res = append(res, r.GetNewStateRequest())
	}
	return res
}

func (srq *StateRequestQueue) AddNewStateRequests(reqs []iwfidl.StateMovement) {
	for _, r := range reqs {
		srq.queue = append(srq.queue, CreateNewStateRequest(r))
	}
}
