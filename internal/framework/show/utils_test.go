package show

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestHasNoData tests the hasNoData helper function
func TestHasNoData(t *testing.T) {
	tests := []struct {
		name     string
		data     []any
		expected bool
	}{
		{
			name:     "empty_slice",
			data:     []any{},
			expected: true,
		},
		{
			name:     "nil_slice",
			data:     nil,
			expected: true,
		},
		{
			name:     "single_element",
			data:     []any{"test"},
			expected: false,
		},
		{
			name:     "multiple_elements",
			data:     []any{"test1", "test2", "test3"},
			expected: false,
		},
		{
			name:     "mixed_types",
			data:     []any{"string", 123, true},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNoData(tt.data)
			if result != tt.expected {
				t.Errorf("hasNoData() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestIsJSONFormat tests the isJSONFormat helper function
func TestIsJSONFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "json_format",
			format:   config.PrintFormatJSON,
			expected: true,
		},
		{
			name:     "table_format",
			format:   config.PrintFormatTable,
			expected: false,
		},
		{
			name:     "empty_format",
			format:   "",
			expected: false,
		},
		{
			name:     "invalid_format",
			format:   "invalid",
			expected: false,
		},
		{
			name:     "case_sensitive_json",
			format:   "JSON",
			expected: false,
		},
		{
			name:     "partial_match",
			format:   "json_extended",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isJSONFormat(tt.format)
			if result != tt.expected {
				t.Errorf("isJSONFormat(%q) = %v, expected %v", tt.format, result, tt.expected)
			}
		})
	}
}

// TestIsAPMisconfigured tests the isAPMisconfigured helper function
func TestIsAPMisconfigured(t *testing.T) {
	tests := []struct {
		name            string
		isMisconfigured bool
		expected        bool
	}{
		{
			name:            "ap_misconfigured_true",
			isMisconfigured: true,
			expected:        true,
		},
		{
			name:            "ap_misconfigured_false",
			isMisconfigured: false,
			expected:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAPMisconfigured(tt.isMisconfigured)
			if result != tt.expected {
				t.Errorf("isAPMisconfigured(%v) = %v, expected %v", tt.isMisconfigured, result, tt.expected)
			}
		})
	}
}

// TestUtilityFunctionsIntegration tests the utility functions working together
func TestUtilityFunctionsIntegration(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "utility_functions_integration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test empty data with JSON format
			emptyData := []any{}
			if !hasNoData(emptyData) {
				t.Error("Expected hasNoData to return true for empty slice")
			}

			// Test JSON format detection
			if !isJSONFormat(config.PrintFormatJSON) {
				t.Error("Expected isJSONFormat to return true for JSON format")
			}

			// Test AP misconfiguration detection
			if !isAPMisconfigured(true) {
				t.Error("Expected isAPMisconfigured to return true for misconfigured AP")
			}

			// Test combination scenarios
			nonEmptyData := []any{"test"}
			if hasNoData(nonEmptyData) {
				t.Error("Expected hasNoData to return false for non-empty slice")
			}

			if isJSONFormat(config.PrintFormatTable) {
				t.Error("Expected isJSONFormat to return false for table format")
			}

			if isAPMisconfigured(false) {
				t.Error("Expected isAPMisconfigured to return false for properly configured AP")
			}
		})
	}
}

// TestUtilityFunctionEdgeCases tests edge cases for utility functions
func TestUtilityFunctionEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "hasNoData_with_nil_elements",
			test: func(t *testing.T) {
				data := []any{nil, nil}
				if hasNoData(data) {
					t.Error("Expected hasNoData to return false for slice with nil elements")
				}
			},
		},
		{
			name: "isJSONFormat_with_whitespace",
			test: func(t *testing.T) {
				formats := []string{" json", "json ", " json ", "\tjson\n"}
				for _, format := range formats {
					if isJSONFormat(format) {
						t.Errorf("Expected isJSONFormat to return false for format with whitespace: %q", format)
					}
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
