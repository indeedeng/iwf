package errors

import "github.com/indeedeng/iwf/gen/iwfidl"

type ErrorAndStatus struct {
	StatusCode int
	Error      iwfidl.ErrorResponse
}
