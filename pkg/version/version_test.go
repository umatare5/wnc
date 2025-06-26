package version

import (
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{
			name:     "returns_version_string",
			expected: "dev", // Current default version
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get()
			if result != tt.expected {
				t.Errorf("Expected version %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestVersionVariable(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		expected string
	}{
		{
			name:     "version_variable_is_dev",
			variable: Version,
			expected: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.variable != tt.expected {
				t.Errorf("Expected version %s, got %s", tt.expected, tt.variable)
			}
		})
	}
}

func TestVersionConsistency(t *testing.T) {
	// Test that Get() returns the same value as Version variable
	if Get() != Version {
		t.Errorf("Get() returned %s, but Version variable is %s", Get(), Version)
	}
}
