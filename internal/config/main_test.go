package config

import (
	"encoding/json"
	"testing"
)

// TestNew tests config creation (Unit test)
func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "creates new config with initialized sub-configs",
			want: Config{
				GenerateCmdConfig: GenerateCmdConfig{},
				ShowCmdConfig:     ShowCmdConfig{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()

			// Check that GenerateCmdConfig is properly initialized (zero value)
			if got.GenerateCmdConfig != (GenerateCmdConfig{}) {
				t.Errorf("New() GenerateCmdConfig not properly initialized")
			}

			// Check that ShowCmdConfig is properly initialized (zero value for slice-containing struct)
			if got.ShowCmdConfig.Controllers != nil {
				t.Errorf("New() ShowCmdConfig.Controllers should be nil, got %v", got.ShowCmdConfig.Controllers)
			}
			if got.ShowCmdConfig.PrintFormat != "" {
				t.Errorf("New() ShowCmdConfig.PrintFormat should be empty, got %v", got.ShowCmdConfig.PrintFormat)
			}
		})
	}
}

// TestConfigJSONSerialization tests JSON serialization (Unit test)
func TestConfigJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "empty config serialization",
			config: Config{
				GenerateCmdConfig: GenerateCmdConfig{},
				ShowCmdConfig:     ShowCmdConfig{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Errorf("JSON marshaling failed: %v", err)
				return
			}

			// Test JSON unmarshaling
			var unmarshaled Config
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}
