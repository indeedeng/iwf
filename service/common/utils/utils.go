package utils

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"net/http"
	"time"
)

const (
	defaultMaxApiTimeoutSeconds = 60
)

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

func TrimContextByTimeoutWithCappedDDL(parent context.Context, waitSeconds *int32, configuredMaxSeconds int64) (context.Context, context.CancelFunc) {
	maxWaitSeconds := configuredMaxSeconds
	if waitSeconds != nil {
		maxWaitSeconds = int64(*waitSeconds)
	}
	if maxWaitSeconds == 0 {
		maxWaitSeconds = defaultMaxApiTimeoutSeconds
	}
	
	newDdlUnix := time.Now().Unix() + maxWaitSeconds

	// then capped by context
	ddl, ok := parent.Deadline()
	if ok {
		maxDdlUnix := ddl.Unix()
		if maxDdlUnix < newDdlUnix {
			newDdlUnix = maxDdlUnix
		}
	}

	newDdl := time.Unix(newDdlUnix, 0)
	return context.WithDeadline(parent, newDdl)
}

func CheckHttpError(err error, httpResp *http.Response) bool {
	if err != nil || (httpResp != nil && httpResp.StatusCode != http.StatusOK) {
		return true
	}
	return false
}
