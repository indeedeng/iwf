package workflow

import "time"

var mockApiLatencyMs int = 0

func SetMockApiLatencyMs(ms int) {
	mockApiLatencyMs = ms
}

func ClearMockApiLatencyMs() {
	mockApiLatencyMs = 0
}

func WaitForMockLatency() {
	if mockApiLatencyMs > 0 {
		time.Sleep(time.Millisecond * time.Duration(mockApiLatencyMs))
	}
}
