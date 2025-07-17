package utils

import (
	"context"
	"testing"
	"time"
)

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
			name:                 "Parent context deadline caps the timeout when sooner",
			reqWaitSeconds:       intPtr(50),
			configuredMaxSeconds: 60,
			parentContextTimeout: durationPtr(10 * time.Second),
			expectedTimeout:      10,
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

		// The deadline should be very close to or before now
		if deadline.After(time.Now().Add(1 * time.Second)) {
			t.Error("Expected deadline to be very soon or in the past")
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
