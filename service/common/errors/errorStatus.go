package errors

import (
	"github.com/indeedeng/iwf/gen/iwfidl"
	"github.com/indeedeng/iwf/service/common/ptr"
)

type ErrorAndStatus struct {
	StatusCode int
	Error      iwfidl.ErrorResponse
}

func NewErrorAndStatus(statusCode int, subStatus iwfidl.ErrorSubStatus, details string) *ErrorAndStatus {
	return &ErrorAndStatus{
		StatusCode: statusCode,
		Error: iwfidl.ErrorResponse{
			SubStatus: ptr.Any(subStatus),
			Detail:    iwfidl.PtrString(details),
		},
	}
}
