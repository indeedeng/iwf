package utils

import (
	"context"
	"github.com/indeedeng/iwf/gen/iwfidl"
	"net/http"
	"time"
)

const (
	defaultMaxApiTimeoutSeconds = 60
	waitBufferSeconds           = 2
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

	ddl, ok := parent.Deadline()

	var newDdlUnix int64
	maxDdlUnix := time.Now().Unix() + maxWaitSeconds
	if ok {
		ddlUnix := ddl.Unix()
		if ddlUnix < maxDdlUnix {
			newDdlUnix = ddlUnix - waitBufferSeconds
		} else {
			newDdlUnix = maxDdlUnix - waitBufferSeconds
		}
	} else {
		newDdlUnix = maxDdlUnix - waitBufferSeconds
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
