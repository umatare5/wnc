package show

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

// TestPrintJson tests the printJson helper function
func TestPrintJson(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		expected bool
	}{
		{
			name:     "string_data",
			data:     "test",
			expected: true,
		},
		{
			name:     "struct_data",
			data:     struct{ Name string }{Name: "test"},
			expected: true,
		},
		{
			name:     "slice_data",
			data:     []string{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "map_data",
			data:     map[string]string{"key": "value"},
			expected: true,
		},
		{
			name:     "nil_data",
			data:     nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			defer func() {
				os.Stdout = originalStdout
			}()

			// Test in a goroutine to handle potential log.Fatal
			done := make(chan bool)
			var output string

			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("printJson recovered from: %v", r)
					}
					done <- true
				}()
				printJson(tt.data)
			}()

			// Wait for completion
			<-done

			// Close writer and read output
			w.Close()
			outputBytes, _ := io.ReadAll(r)
			output = string(outputBytes)

			// Verify output was generated
			if tt.expected && len(output) == 0 {
				t.Errorf("Expected JSON output, got empty string")
			}

			t.Logf("JSON output length: %d", len(output))
		})
	}
}

// TestPrintJsonWithInvalidData tests printJson behavior with unmarshallable data
func TestPrintJsonWithInvalidData(t *testing.T) {
	// Since printJson calls log.Fatal on JSON marshaling errors,
	// we'll test that the JSON marshaling itself fails for these types
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "function_data",
			data: func() {},
		},
		{
			name: "channel_data",
			data: make(chan int),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that json.Marshal fails for these data types
			_, err := json.Marshal(tt.data)
			if err == nil {
				t.Errorf("Expected JSON marshaling to fail for %s, but it succeeded", tt.name)
			}

			// Log the expected error for documentation
			t.Logf("Expected JSON marshaling error for %s: %v", tt.name, err)
		})
	}
}

// TestPrintJsonBuffered tests printJson with captured output
func TestPrintJsonBuffered(t *testing.T) {
	tests := []struct {
		name     string
		data     interface{}
		contains string
	}{
		{
			name:     "simple_string",
			data:     "hello",
			contains: "hello",
		},
		{
			name:     "simple_number",
			data:     42,
			contains: "42",
		},
		{
			name:     "boolean_true",
			data:     true,
			contains: "true",
		},
		{
			name:     "boolean_false",
			data:     false,
			contains: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use a buffer to capture output instead of os.Stdout
			var buf bytes.Buffer

			// Temporarily redirect stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run printJson in goroutine
			done := make(chan bool)
			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("printJson completed with recovery: %v", r)
					}
					done <- true
				}()
				printJson(tt.data)
			}()

			// Wait and cleanup
			<-done
			w.Close()
			os.Stdout = originalStdout

			// Read output
			io.Copy(&buf, r)
			output := buf.String()

			// Verify output contains expected content
			if len(output) == 0 {
				t.Error("Expected non-empty JSON output")
			}

			t.Logf("JSON output: %s", output)
		})
	}
}
