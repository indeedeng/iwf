package logevent

import "github.com/indeedeng/iwf/gen/iwfidl"

// The implementation must be lightweight, reliable and fast (less than 1s)
type LogEventFunc func(event iwfidl.IwfEvent)

var Log LogEventFunc = DefaultLogEventFunc

func SetLogEventFunc(logger LogEventFunc) {
	Log = logger
}

func DefaultLogEventFunc(event iwfidl.IwfEvent) {
	// Noop by default
}
