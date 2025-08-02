package humanize

import (
	"testing"
)

// TestFormatComma tests the FormatComma function
func TestFormatComma(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "small_number",
			input:    123,
			expected: "123",
		},
		{
			name:     "thousand",
			input:    1000,
			expected: "1,000",
		},
		{
			name:     "million",
			input:    1000000,
			expected: "1,000,000",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "negative",
			input:    -1000,
			expected: "-1,000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatComma(tt.input)
			if result != tt.expected {
				t.Errorf("FormatComma(%d) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatBytes tests the FormatBytes function
func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "small_bytes",
			input:    500,
			expected: "500 B",
		},
		{
			name:     "zero_bytes",
			input:    0,
			expected: "0 B",
		},
		{
			name:     "kilobytes",
			input:    1024,
			expected: "1 KB",
		},
		{
			name:     "large_kilobytes",
			input:    1536,
			expected: "1 KB",
		},
		{
			name:     "multiple_kilobytes",
			input:    5120,
			expected: "5 KB",
		},
		{
			name:     "large_bytes",
			input:    1048576,
			expected: "1,024 KB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestFormatTimeoutSeconds tests the FormatTimeoutSeconds function
func TestFormatTimeoutSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "small_timeout",
			input:    30,
			expected: "30s",
		},
		{
			name:     "zero_timeout",
			input:    0,
			expected: "0s",
		},
		{
			name:     "large_timeout",
			input:    1000,
			expected: "1,000s",
		},
		{
			name:     "very_large_timeout",
			input:    123456,
			expected: "123,456s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeoutSeconds(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTimeoutSeconds(%d) = %s, expected %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHumanizePackageStructure tests the overall package structure
func TestHumanizePackageStructure(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "format_comma_function_exists",
			test: func(t *testing.T) {
				// Test that FormatComma function exists and is callable
				result := FormatComma(1000)
				if result == "" {
					t.Error("FormatComma should return non-empty string")
				}
			},
		},
		{
			name: "format_bytes_function_exists",
			test: func(t *testing.T) {
				// Test that FormatBytes function exists and is callable
				result := FormatBytes(1024)
				if result == "" {
					t.Error("FormatBytes should return non-empty string")
				}
			},
		},
		{
			name: "format_timeout_seconds_function_exists",
			test: func(t *testing.T) {
				// Test that FormatTimeoutSeconds function exists and is callable
				result := FormatTimeoutSeconds(60)
				if result == "" {
					t.Error("FormatTimeoutSeconds should return non-empty string")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestHumanizeEdgeCases tests edge cases and boundary conditions
func TestHumanizeEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "max_int64_comma",
			test: func(t *testing.T) {
				result := FormatComma(9223372036854775807) // max int64
				if result == "" {
					t.Error("FormatComma should handle max int64")
				}
			},
		},
		{
			name: "min_int64_comma",
			test: func(t *testing.T) {
				result := FormatComma(-9223372036854775808) // min int64
				if result == "" {
					t.Error("FormatComma should handle min int64")
				}
			},
		},
		{
			name: "large_bytes_formatting",
			test: func(t *testing.T) {
				result := FormatBytes(9223372036854775807) // max int64
				if result == "" {
					t.Error("FormatBytes should handle large values")
				}
			},
		},
		{
			name: "large_timeout_formatting",
			test: func(t *testing.T) {
				result := FormatTimeoutSeconds(9223372036854775807) // max int64
				if result == "" {
					t.Error("FormatTimeoutSeconds should handle large values")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}
