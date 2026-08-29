package urlautofix

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultFixWorkerUrlFunc_NoEnvVars(t *testing.T) {
	// Ensure env vars are not set
	os.Unsetenv("AUTO_FIX_WORKER_URL")
	os.Unsetenv("AUTO_FIX_WORKER_PORT_FROM_ENV")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "URL unchanged without env vars",
			input:    "http://localhost:8080/api",
			expected: "http://localhost:8080/api",
		},
		{
			name:     "Trailing slash removed",
			input:    "http://localhost:8080/api/",
			expected: "http://localhost:8080/api",
		},
		{
			name:     "Multiple trailing slashes removed",
			input:    "http://localhost:8080/api///",
			expected: "http://localhost:8080/api",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultFixWorkerUrlFunc(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFixWorkerUrlFunc_AutoFixWorkerUrl(t *testing.T) {
	os.Unsetenv("AUTO_FIX_WORKER_PORT_FROM_ENV")

	tests := []struct {
		name       string
		autoFixUrl string
		input      string
		expected   string
	}{
		{
			name:       "Replace localhost with host.docker.internal",
			autoFixUrl: "host.docker.internal",
			input:      "http://localhost:8080/api",
			expected:   "http://host.docker.internal:8080/api",
		},
		{
			name:       "Replace 127.0.0.1 with custom host",
			autoFixUrl: "host.docker.internal",
			input:      "http://127.0.0.1:8080/api",
			expected:   "http://host.docker.internal:8080/api",
		},
		{
			name:       "Only first occurrence of localhost replaced",
			autoFixUrl: "myhost",
			input:      "http://localhost:8080/localhost",
			expected:   "http://myhost:8080/localhost",
		},
		{
			name:       "No match leaves URL unchanged",
			autoFixUrl: "myhost",
			input:      "http://example.com:8080/api",
			expected:   "http://example.com:8080/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AUTO_FIX_WORKER_URL", tt.autoFixUrl)
			defer os.Unsetenv("AUTO_FIX_WORKER_URL")

			result := DefaultFixWorkerUrlFunc(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFixWorkerUrlFunc_AutoFixWorkerPort(t *testing.T) {
	os.Unsetenv("AUTO_FIX_WORKER_URL")

	tests := []struct {
		name       string
		portEnvVar string
		portValue  string
		input      string
		expected   string
	}{
		{
			name:       "Replace port placeholder with env var value",
			portEnvVar: "MY_WORKER_PORT",
			portValue:  "9090",
			input:      "http://localhost:$MY_WORKER_PORT$/api",
			expected:   "http://localhost:9090/api",
		},
		{
			name:       "Empty port env value replaces with empty",
			portEnvVar: "MY_WORKER_PORT",
			portValue:  "",
			input:      "http://localhost:$MY_WORKER_PORT$/api",
			expected:   "http://localhost:/api",
		},
		{
			name:       "No placeholder in URL leaves it unchanged",
			portEnvVar: "MY_WORKER_PORT",
			portValue:  "9090",
			input:      "http://localhost:8080/api",
			expected:   "http://localhost:8080/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("AUTO_FIX_WORKER_PORT_FROM_ENV", tt.portEnvVar)
			os.Setenv(tt.portEnvVar, tt.portValue)
			defer func() {
				os.Unsetenv("AUTO_FIX_WORKER_PORT_FROM_ENV")
				os.Unsetenv(tt.portEnvVar)
			}()

			result := DefaultFixWorkerUrlFunc(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultFixWorkerUrlFunc_BothEnvVars(t *testing.T) {
	os.Setenv("AUTO_FIX_WORKER_URL", "host.docker.internal")
	os.Setenv("AUTO_FIX_WORKER_PORT_FROM_ENV", "WORKER_PORT")
	os.Setenv("WORKER_PORT", "9090")
	defer func() {
		os.Unsetenv("AUTO_FIX_WORKER_URL")
		os.Unsetenv("AUTO_FIX_WORKER_PORT_FROM_ENV")
		os.Unsetenv("WORKER_PORT")
	}()

	result := DefaultFixWorkerUrlFunc("http://localhost:$WORKER_PORT$/api/")
	assert.Equal(t, "http://host.docker.internal:9090/api", result)
}

func TestFixWorkerUrl_UsesDefaultFixer(t *testing.T) {
	os.Unsetenv("AUTO_FIX_WORKER_URL")
	os.Unsetenv("AUTO_FIX_WORKER_PORT_FROM_ENV")

	// Reset to default fixer
	SetWorkerUrlFixer(DefaultFixWorkerUrlFunc)

	result := FixWorkerUrl("http://localhost:8080/api/")
	assert.Equal(t, "http://localhost:8080/api", result)
}

func TestSetWorkerUrlFixer_CustomFixer(t *testing.T) {
	originalFixer := workerUrlFixer
	defer SetWorkerUrlFixer(originalFixer)

	customFixer := func(url string) string {
		return "custom://" + url
	}
	SetWorkerUrlFixer(customFixer)

	result := FixWorkerUrl("test-url")
	assert.Equal(t, "custom://test-url", result)
}
