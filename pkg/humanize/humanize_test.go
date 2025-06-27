package humanize

import (
	"testing"
)

func TestFormatComma(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "small number",
			input:    123,
			expected: "123",
		},
		{
			name:     "thousands",
			input:    1234,
			expected: "1,234",
		},
		{
			name:     "millions",
			input:    1234567,
			expected: "1,234,567",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatComma(tt.input)
			if result != tt.expected {
				t.Errorf("FormatComma(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "bytes less than 1KB",
			input:    512,
			expected: "512 B",
		},
		{
			name:     "exactly 1KB",
			input:    1024,
			expected: "1 KB",
		},
		{
			name:     "multiple KB",
			input:    2048,
			expected: "2 KB",
		},
		{
			name:     "large KB value",
			input:    1234567,
			expected: "1,205 KB",
		},
		{
			name:     "zero bytes",
			input:    0,
			expected: "0 B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.input)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatTimeoutSeconds(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected string
	}{
		{
			name:     "small timeout",
			input:    30,
			expected: "30s",
		},
		{
			name:     "large timeout",
			input:    1800,
			expected: "1,800s",
		},
		{
			name:     "zero timeout",
			input:    0,
			expected: "0s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTimeoutSeconds(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTimeoutSeconds(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
