package cli

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestHasNoData(t *testing.T) {
	tests := []struct {
		name     string
		data     []any
		expected bool
	}{
		{
			name:     "empty slice",
			data:     []any{},
			expected: true,
		},
		{
			name:     "nil slice",
			data:     nil,
			expected: true,
		},
		{
			name:     "single item",
			data:     []any{"item"},
			expected: false,
		},
		{
			name:     "multiple items",
			data:     []any{"item1", "item2", "item3"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNoData(tt.data)
			if result != tt.expected {
				t.Errorf("hasNoData(%v) = %v, expected %v", tt.data, result, tt.expected)
			}
		})
	}
}

func TestIsJSONFormat(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected bool
	}{
		{
			name:     "json format",
			format:   config.PrintFormatJSON,
			expected: true,
		},
		{
			name:     "table format",
			format:   config.PrintFormatTable,
			expected: false,
		},
		{
			name:     "empty format",
			format:   "",
			expected: false,
		},
		{
			name:     "unknown format",
			format:   "unknown",
			expected: false,
		},
		{
			name:     "case sensitive test",
			format:   "JSON",
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

func TestIsAPMisconfigured(t *testing.T) {
	tests := []struct {
		name            string
		isMisconfigured bool
		expected        bool
	}{
		{
			name:            "misconfigured true",
			isMisconfigured: true,
			expected:        true,
		},
		{
			name:            "misconfigured false",
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

func TestUtilsJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "function results",
			data: map[string]interface{}{
				"hasNoData":         hasNoData([]any{}),
				"isJSONFormat":      isJSONFormat(config.PrintFormatJSON),
				"isAPMisconfigured": isAPMisconfigured(true),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("Failed to marshal data: %v", err)
			}

			var unmarshaled map[string]interface{}
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal data: %v", err)
			}
		})
	}
}

func TestUtilsFailFast(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "hasNoData should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("hasNoData panicked: %v", r)
					}
				}()
				_ = hasNoData(nil)
				_ = hasNoData([]any{})
				_ = hasNoData([]any{"item"})
			},
		},
		{
			name: "isJSONFormat should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("isJSONFormat panicked: %v", r)
					}
				}()
				_ = isJSONFormat("")
				_ = isJSONFormat(config.PrintFormatJSON)
				_ = isJSONFormat("invalid")
			},
		},
		{
			name: "isAPMisconfigured should not panic",
			test: func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("isAPMisconfigured panicked: %v", r)
					}
				}()
				_ = isAPMisconfigured(true)
				_ = isAPMisconfigured(false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestUtilsTableDriven(t *testing.T) {
	// Test multiple combinations of utility functions
	tests := []struct {
		name            string
		dataLen         int
		format          string
		misconfigured   bool
		expectEmpty     bool
		expectJSON      bool
		expectMisconfig bool
	}{
		{
			name:            "empty data, json format, misconfigured",
			dataLen:         0,
			format:          config.PrintFormatJSON,
			misconfigured:   true,
			expectEmpty:     true,
			expectJSON:      true,
			expectMisconfig: true,
		},
		{
			name:            "with data, table format, not misconfigured",
			dataLen:         3,
			format:          config.PrintFormatTable,
			misconfigured:   false,
			expectEmpty:     false,
			expectJSON:      false,
			expectMisconfig: false,
		},
		{
			name:            "with data, json format, not misconfigured",
			dataLen:         1,
			format:          config.PrintFormatJSON,
			misconfigured:   false,
			expectEmpty:     false,
			expectJSON:      true,
			expectMisconfig: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test data based on length
			var data []any
			for i := 0; i < tt.dataLen; i++ {
				data = append(data, "item")
			}

			// Test all utility functions
			isEmpty := hasNoData(data)
			isJSON := isJSONFormat(tt.format)
			isMisconfig := isAPMisconfigured(tt.misconfigured)

			if isEmpty != tt.expectEmpty {
				t.Errorf("hasNoData = %v, expected %v", isEmpty, tt.expectEmpty)
			}
			if isJSON != tt.expectJSON {
				t.Errorf("isJSONFormat = %v, expected %v", isJSON, tt.expectJSON)
			}
			if isMisconfig != tt.expectMisconfig {
				t.Errorf("isAPMisconfigured = %v, expected %v", isMisconfig, tt.expectMisconfig)
			}
		})
	}
}
