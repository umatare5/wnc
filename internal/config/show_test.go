package config

import (
	"encoding/json"
	"testing"
)

// TestShowCmdConfigJSONSerialization tests JSON serialization for ShowCmdConfig (Unit test)
func TestShowCmdConfigJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config ShowCmdConfig
	}{
		{
			name: "empty_show_config",
			config: ShowCmdConfig{
				Controllers: nil,
				PrintFormat: "",
			},
		},
		{
			name: "full_show_config",
			config: ShowCmdConfig{
				Controllers: []Controller{
					{Hostname: "controller1", AccessToken: "token1"},
					{Hostname: "controller2", AccessToken: "token2"},
				},
				PrintFormat: "table",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var config ShowCmdConfig
			err = json.Unmarshal(jsonData, &config)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Basic validation
			if config.PrintFormat != tt.config.PrintFormat {
				t.Errorf("PrintFormat mismatch: got %s, want %s", config.PrintFormat, tt.config.PrintFormat)
			}
			if len(config.Controllers) != len(tt.config.Controllers) {
				t.Errorf("Controllers length mismatch: got %d, want %d", len(config.Controllers), len(tt.config.Controllers))
			}
		})
	}
}
