package utils

import (
	"context"
	"testing"
	"time"
)

// Note: TrimContextByTimeoutWithCappedDDL applies a 1-second buffer (contextBufferSeconds)
// when using parent context deadlines to ensure the new context times out before the parent

func TestTrimContextByTimeoutWithCappedDDL(t *testing.T) {
	tests := []struct {
		name                 string
		reqWaitSeconds       *int32
		configuredMaxSeconds int64
		parentContextTimeout *time.Duration // nil means no deadline
		expectedTimeout      int64          // expected timeout in seconds (approximate)
	}{
		{
			name:                 "ConfiguredMaxSeconds is 0, should use default",
			reqWaitSeconds:       nil,
			configuredMaxSeconds: 0,
			parentContextTimeout: nil,
			expectedTimeout:      defaultMaxApiTimeoutSeconds,
		},
		{
			name:                 "Use configuredMaxSeconds when reqWaitSeconds is nil",
			reqWaitSeconds:       nil,
			configuredMaxSeconds: 30,
			parentContextTimeout: nil,
			expectedTimeout:      30,
		},
		{
			name:                 "Use reqWaitSeconds when smaller than configuredMaxSeconds",
			reqWaitSeconds:       intPtr(20),
			configuredMaxSeconds: 30,
			parentContextTimeout: nil,
			expectedTimeout:      20,
		},
		{
			name:                 "Use configuredMaxSeconds when reqWaitSeconds is larger",
			reqWaitSeconds:       intPtr(50),
			configuredMaxSeconds: 30,
			parentContextTimeout: nil,
			expectedTimeout:      30,
		},
		{
			name:                 "Ignore reqWaitSeconds when it's 0",
			reqWaitSeconds:       intPtr(0),
			configuredMaxSeconds: 30,
			parentContextTimeout: nil,
			expectedTimeout:      30,
		},
		{
			name:                 "Ignore reqWaitSeconds when it's negative",
			reqWaitSeconds:       intPtr(-10),
			configuredMaxSeconds: 30,
			parentContextTimeout: nil,
			expectedTimeout:      30,
		},
		{
			name:                 "Parent context deadline caps the timeout when sooner (with buffer)",
			reqWaitSeconds:       intPtr(50),
			configuredMaxSeconds: 60,
			parentContextTimeout: durationPtr(10 * time.Second),
			expectedTimeout:      9, // 10 - 1 second buffer
		},
		{
			name:                 "Parent context deadline doesn't affect when later",
			reqWaitSeconds:       intPtr(20),
			configuredMaxSeconds: 30,
			parentContextTimeout: durationPtr(60 * time.Second),
			expectedTimeout:      20,
		},
		{
			name:                 "Complex case: reqWaitSeconds < configuredMaxSeconds < parent deadline",
			reqWaitSeconds:       intPtr(15),
			configuredMaxSeconds: 25,
			parentContextTimeout: durationPtr(40 * time.Second),
			expectedTimeout:      15,
		},
		{
			name:                 "Parent context deadline with small buffer becomes limiting factor",
			reqWaitSeconds:       intPtr(10),
			configuredMaxSeconds: 15,
			parentContextTimeout: durationPtr(8 * time.Second),
			expectedTimeout:      7, // 8 - 1 second buffer = 7, which is less than reqWaitSeconds(10)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create parent context
			var parentCtx context.Context
			var cancel context.CancelFunc

			if tt.parentContextTimeout != nil {
				parentCtx, cancel = context.WithTimeout(context.Background(), *tt.parentContextTimeout)
				defer cancel()
			} else {
				parentCtx = context.Background()
			}

			// Record the time before calling the function
			startTime := time.Now()

			// Call the function
			newCtx, newCancel := TrimContextByTimeoutWithCappedDDL(parentCtx, tt.reqWaitSeconds, tt.configuredMaxSeconds)
			defer newCancel()

			// Check that the new context has a deadline
			deadline, ok := newCtx.Deadline()
			if !ok {
				t.Fatal("Expected new context to have a deadline")
			}

			// Calculate the actual timeout
			actualTimeout := deadline.Sub(startTime).Seconds()

			// Allow for some tolerance due to timing differences (1 second)
			tolerance := 1.0
			if actualTimeout < float64(tt.expectedTimeout)-tolerance || actualTimeout > float64(tt.expectedTimeout)+tolerance {
				t.Errorf("Expected timeout ~%d seconds, got %.2f seconds", tt.expectedTimeout, actualTimeout)
			}
		})
	}
}

func TestTrimContextByTimeoutWithCappedDDL_EdgeCases(t *testing.T) {
	t.Run("Very short parent context deadline", func(t *testing.T) {
		// Create a context that expires in 1 millisecond
		parentCtx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Wait a bit to ensure the context is expired or very close to expiry
		time.Sleep(2 * time.Millisecond)

		newCtx, newCancel := TrimContextByTimeoutWithCappedDDL(parentCtx, intPtr(30), 60)
		defer newCancel()

		// The new context should still have a deadline, even if it's in the past
		deadline, ok := newCtx.Deadline()
		if !ok {
			t.Fatal("Expected new context to have a deadline")
		}

		// The deadline should be in the past due to the 1-second buffer
		if deadline.After(time.Now()) {
			t.Error("Expected deadline to be in the past due to buffer")
		}
	})

	t.Run("Context buffer seconds behavior", func(t *testing.T) {
		// Test specifically that the 1-second buffer is applied correctly
		parentCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		startTime := time.Now()
		newCtx, newCancel := TrimContextByTimeoutWithCappedDDL(parentCtx, intPtr(10), 60)
		defer newCancel()

		deadline, ok := newCtx.Deadline()
		if !ok {
			t.Fatal("Expected new context to have a deadline")
		}

		// The deadline should be approximately 4 seconds from now (5 - 1 second buffer)
		actualTimeout := deadline.Sub(startTime).Seconds()
		expectedTimeout := 4.0 // 5 seconds parent timeout - 1 second buffer

		if actualTimeout < expectedTimeout-1.0 || actualTimeout > expectedTimeout+1.0 {
			t.Errorf("Expected timeout ~%.0f seconds (parent deadline minus buffer), got %.2f seconds", expectedTimeout, actualTimeout)
		}
	})

	t.Run("Buffer results in very small or negative timeout", func(t *testing.T) {
		// Test when parent deadline minus buffer is very small
		parentCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()

		startTime := time.Now()
		newCtx, newCancel := TrimContextByTimeoutWithCappedDDL(parentCtx, intPtr(10), 60)
		defer newCancel()

		deadline, ok := newCtx.Deadline()
		if !ok {
			t.Fatal("Expected new context to have a deadline")
		}

		// The deadline should be in the past or very soon (0.5 - 1 = -0.5 seconds)
		actualTimeout := deadline.Sub(startTime).Seconds()
		if actualTimeout > 0.5 {
			t.Errorf("Expected timeout to be very small or negative due to buffer, got %.2f seconds", actualTimeout)
		}
	})

	t.Run("Nil parent context", func(t *testing.T) {
		// This should not panic and should work with background context
		newCtx, newCancel := TrimContextByTimeoutWithCappedDDL(context.Background(), intPtr(30), 60)
		defer newCancel()

		deadline, ok := newCtx.Deadline()
		if !ok {
			t.Fatal("Expected new context to have a deadline")
		}

		expectedDeadline := time.Now().Add(30 * time.Second)
		if deadline.Before(expectedDeadline.Add(-1*time.Second)) || deadline.After(expectedDeadline.Add(1*time.Second)) {
			t.Error("Deadline not set correctly")
		}
	})
}

// Helper functions
func intPtr(i int32) *int32 {
	return &i
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}
