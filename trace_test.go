package toolbox

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestSetTracing(t *testing.T) {
	tests := []struct {
		name string
		on   bool
		want bool
	}{
		{
			name: "enable tracing",
			on:   true,
			want: true,
		},
		{
			name: "disable tracing",
			on:   false,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTracing(tt.on)
			if isTracingOn != tt.want {
				t.Errorf("SetTracing(%v): isTracingOn = %v, want %v", tt.on, isTracingOn, tt.want)
			}
		})
	}

	// Cleanup
	SetTracing(false)
}

func TestTracef(t *testing.T) {
	tests := []struct {
		name           string
		tracingEnabled bool
		format         string
		args           []any
		expectOutput   bool
		expectContent  string
	}{
		{
			name:           "tracing disabled",
			tracingEnabled: false,
			format:         "test message",
			args:           []any{},
			expectOutput:   false,
			expectContent:  "",
		},
		{
			name:           "tracing enabled with simple message",
			tracingEnabled: true,
			format:         "simple message",
			args:           []any{},
			expectOutput:   true,
			expectContent:  "simple message",
		},
		{
			name:           "tracing enabled with format",
			tracingEnabled: true,
			format:         "message with %s",
			args:           []any{"value"},
			expectOutput:   true,
			expectContent:  "message with value",
		},
		{
			name:           "tracing enabled without newline",
			tracingEnabled: true,
			format:         "no newline",
			args:           []any{},
			expectOutput:   true,
			expectContent:  "no newline",
		},
		{
			name:           "tracing enabled with newline",
			tracingEnabled: true,
			format:         "with newline\n",
			args:           []any{},
			expectOutput:   true,
			expectContent:  "with newline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			// Set tracing state
			SetTracing(tt.tracingEnabled)

			// Call Tracef
			Tracef(tt.format, tt.args...)

			// Restore stderr and read output
			_ = w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			if tt.expectOutput {
				if output == "" {
					t.Errorf("Tracef(%q, %v): expected output but got none", tt.format, tt.args)
				}
				if !strings.Contains(output, tt.expectContent) {
					t.Errorf("Tracef(%q, %v): expected %q in output, got %q", tt.format, tt.args, tt.expectContent, output)
				}
				// Verify trace prefix is present
				if !strings.Contains(output, "## TOOLBOX TRACE") {
					t.Errorf("Tracef(%q, %v): expected trace prefix in output, got %q", tt.format, tt.args, output)
				}
			} else {
				if output != "" {
					t.Errorf("Tracef(%q, %v): expected no output but got %q", tt.format, tt.args, output)
				}
			}
		})
	}

	// Cleanup
	SetTracing(false)
}

func TestTracefMultipleArgs(t *testing.T) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	SetTracing(true)
	Tracef("test %s %d %v", "string", 42, true)

	_ = w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "test string 42 true") {
		t.Errorf("Tracef with multiple args: expected formatted output, got %q", output)
	}

	SetTracing(false)
}
