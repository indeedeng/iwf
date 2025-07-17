package utils

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

const (
	defaultMaxApiTimeoutSeconds = 60
	contextBufferSeconds        = 1
)

func MergeStringSlice(first, second []string) []string {
	exists := map[string]bool{}
	var out []string
	for _, k := range first {
		if !exists[k] {
			exists[k] = true
			out = append(out, k)
		}
	}
	for _, k := range second {
		if !exists[k] {
			exists[k] = true
			out = append(out, k)
		}
	}
	return out
}

func MergeMap(first map[string]interface{}, second map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(first))
	for k, v := range first {
		out[k] = v
	}

	for k, v := range second {
		out[k] = v
	}
	return out
}

func TrimRpcTimeoutSeconds(ctx context.Context, req iwfidl.WorkflowRpcRequest) int32 {
	secondsRemaining := int32(defaultMaxApiTimeoutSeconds)
	ddl, ok := ctx.Deadline()
	if ok {
		timeRemaining := ddl.Sub(time.Now())
		if int32(timeRemaining.Seconds()) < secondsRemaining {
			secondsRemaining = int32(timeRemaining.Seconds())
		}
	}
	if req.TimeoutSeconds == nil && req.GetTimeoutSeconds() > 0 && req.GetTimeoutSeconds() < secondsRemaining {
		secondsRemaining = req.GetTimeoutSeconds()
	}
	return secondsRemaining
}

func TrimContextByTimeoutWithCappedDDL(parent context.Context, reqWaitSeconds *int32, configuredMaxSeconds int64) (context.Context, context.CancelFunc) {
	maxWaitSeconds := configuredMaxSeconds
	if maxWaitSeconds == 0 {
		maxWaitSeconds = defaultMaxApiTimeoutSeconds
	}

	if reqWaitSeconds != nil && *reqWaitSeconds > 0 && int64(*reqWaitSeconds) < maxWaitSeconds {
		maxWaitSeconds = int64(*reqWaitSeconds)
	}

	maxWaitTimestamp := time.Now().Unix() + maxWaitSeconds

	// then capped by context minus contextBufferSeconds(so that the context doesn't timeout)
	ddl, ok := parent.Deadline()
	if ok {
		maxDdlUnix := ddl.Unix() - contextBufferSeconds
		if maxDdlUnix < maxWaitTimestamp {
			maxWaitTimestamp = maxDdlUnix
		}
	}

	newDdl := time.Unix(maxWaitTimestamp, 0)
	return context.WithDeadline(parent, newDdl)
}

func CheckHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}

func ToNanoSeconds(e *timestamppb.Timestamp) int64 {
	return e.GetSeconds()*1000*1000*1000 + int64(e.GetNanos())
}

func GetWorkflowIdForWaitForStateExecution(parentId string, stateExeId *string, waitForKey *string, stateId *string) string {
	if waitForKey != nil && *waitForKey != "" {
		return service.IwfSystemConstPrefix + parentId + "_" + *stateId + "_" + *waitForKey
	} else {
		return service.IwfSystemConstPrefix + parentId + "_" + *stateExeId
	}
}
