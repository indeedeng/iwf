package handleevent

import "github.com/indeedeng/iwf/gen/iwfidl"

// The implementation must be lightweight, reliable and fast (less than 1s)
type HandleEventFunc func(event iwfidl.IwfEvent)

var Handle HandleEventFunc = DefaultHandleEventFunc

func SetHandleEventFunc(handler HandleEventFunc) {
	Handle = handler
}

func DefaultHandleEventFunc(event iwfidl.IwfEvent) {
	// Noop by default
}
