package timeparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime_EmptyString(t *testing.T) {
	result, err := ParseTime("")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), result)
}

func TestParseTime_DateTimeFormat(t *testing.T) {
	timeStr := "2023-06-15T10:30:00-07:00"
	parsedTime, err := time.Parse(DateTimeFormat, timeStr)
	assert.NoError(t, err)

	result, err := ParseTime(timeStr)
	assert.NoError(t, err)
	assert.Equal(t, parsedTime.UnixNano(), result)
}

func TestParseTime_RawUnixNano(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "Simple integer",
			input:    "1687000000000000000",
			expected: 1687000000000000000,
		},
		{
			name:     "Small integer",
			input:    "12345",
			expected: 12345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTime_TimeRangeShort(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		duration time.Duration
	}{
		{
			name:     "Seconds short form",
			input:    "30s",
			duration: 30 * time.Second,
		},
		{
			name:     "Minutes short form",
			input:    "5m",
			duration: 5 * time.Minute,
		},
		{
			name:     "Hours short form",
			input:    "2h",
			duration: 2 * time.Hour,
		},
		{
			name:     "Days short form",
			input:    "3d",
			duration: 3 * 24 * time.Hour,
		},
		{
			name:     "Weeks short form",
			input:    "1w",
			duration: 7 * 24 * time.Hour,
		},
		{
			name:     "Months short form",
			input:    "2M",
			duration: 2 * 30 * 24 * time.Hour,
		},
		{
			name:     "Years short form",
			input:    "1y",
			duration: 365 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			result, err := ParseTime(tt.input)
			after := time.Now()

			assert.NoError(t, err)

			expectedLow := before.Add(-tt.duration).UnixNano()
			expectedHigh := after.Add(-tt.duration).UnixNano()

			assert.True(t, result >= expectedLow && result <= expectedHigh,
				"expected result %d to be between %d and %d", result, expectedLow, expectedHigh)
		})
	}
}

func TestParseTime_TimeRangeLong(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		duration time.Duration
	}{
		{
			name:     "Seconds long form",
			input:    "30second",
			duration: 30 * time.Second,
		},
		{
			name:     "Minutes long form",
			input:    "5minute",
			duration: 5 * time.Minute,
		},
		{
			name:     "Hours long form",
			input:    "2hour",
			duration: 2 * time.Hour,
		},
		{
			name:     "Days long form",
			input:    "3day",
			duration: 3 * 24 * time.Hour,
		},
		{
			name:     "Weeks long form",
			input:    "1week",
			duration: 7 * 24 * time.Hour,
		},
		{
			name:     "Months long form",
			input:    "2month",
			duration: 2 * 30 * 24 * time.Hour,
		},
		{
			name:     "Years long form",
			input:    "1year",
			duration: 365 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			result, err := ParseTime(tt.input)
			after := time.Now()

			assert.NoError(t, err)

			expectedLow := before.Add(-tt.duration).UnixNano()
			expectedHigh := after.Add(-tt.duration).UnixNano()

			assert.True(t, result >= expectedLow && result <= expectedHigh,
				"expected result %d to be between %d and %d", result, expectedLow, expectedHigh)
		})
	}
}

func TestParseTime_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid string",
			input: "not-a-time",
		},
		{
			name:  "Invalid time range suffix",
			input: "5x",
		},
		{
			name:  "Zero prefix not allowed",
			input: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.input)
			assert.Error(t, err)
		})
	}
}

func TestParseTimeDuration(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Duration
		hasError bool
	}{
		{"second short", "s", time.Second, false},
		{"second long", "second", time.Second, false},
		{"minute short", "m", time.Minute, false},
		{"minute long", "minute", time.Minute, false},
		{"hour short", "h", time.Hour, false},
		{"hour long", "hour", time.Hour, false},
		{"day short", "d", day, false},
		{"day long", "day", day, false},
		{"week short", "w", week, false},
		{"week long", "week", week, false},
		{"month short", "M", month, false},
		{"month long", "month", month, false},
		{"year short", "y", year, false},
		{"year long", "year", year, false},
		{"unknown duration", "x", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimeDuration(tt.input)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestParseTimeRange_MultiplierLimit(t *testing.T) {
	// Multiplier must be less than 1e6
	_, err := parseTimeRange("1000000s")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid time-duation multiplier")

	// Just under the limit should work
	_, err = parseTimeRange("999999s")
	assert.NoError(t, err)
}

func TestParseTimeRange_EpochFloor(t *testing.T) {
	// 60 years in the past from 2026 is ~1966, which is before epoch (1970).
	// The function should clamp the result to epoch.
	result, err := parseTimeRange("60y")
	assert.NoError(t, err)

	epochTime := time.Unix(0, 0)
	assert.Equal(t, epochTime, result)
}
