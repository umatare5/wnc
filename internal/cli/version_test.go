package cli

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/pkg/version"
)

func TestVersionFunction(t *testing.T) {
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
			result := getVersion()
			if result != tt.expected {
				t.Errorf("Expected version %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestVersionConstants(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		expected string
	}{
		{
			name:     "version_variable_is_dev",
			variable: version.Version,
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

func TestVersionJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "dev_version",
			version: "dev",
		},
		{
			name:    "release_version",
			version: "1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a simple struct to test version serialization
			versionInfo := struct {
				Version string `json:"version"`
			}{
				Version: tt.version,
			}

			// Test JSON marshaling
			jsonData, err := json.Marshal(versionInfo)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaled struct {
				Version string `json:"version"`
			}
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			if unmarshaled.Version != tt.version {
				t.Errorf("Expected version %s, got %s", tt.version, unmarshaled.Version)
			}
		})
	}
}

func TestVersionFailFast(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "getVersion_should_not_panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("getVersion() panicked: %v", r)
				}
			}()

			result := getVersion()
			if result == "" {
				t.Error("Expected non-empty version string")
			}
		})
	}
}

func TestVersionTableDriven(t *testing.T) {
	tests := []struct {
		name           string
		currentVersion string
		expectEmpty    bool
		expectDev      bool
	}{
		{
			name:           "current_version_is_dev",
			currentVersion: version.Version,
			expectEmpty:    false,
			expectDev:      true,
		},
		{
			name:           "version_function_returns_consistent_value",
			currentVersion: getVersion(),
			expectEmpty:    false,
			expectDev:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectEmpty && tt.currentVersion != "" {
				t.Errorf("Expected empty version, got %s", tt.currentVersion)
			}
			if !tt.expectEmpty && tt.currentVersion == "" {
				t.Error("Expected non-empty version")
			}
			if tt.expectDev && tt.currentVersion != "dev" {
				t.Errorf("Expected dev version, got %s", tt.currentVersion)
			}
		})
	}
}
