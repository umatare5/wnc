package config

import (
	"encoding/json"
	"testing"
)

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
		{
			name: "config with show cmd data",
			config: Config{
				GenerateCmdConfig: GenerateCmdConfig{},
				ShowCmdConfig: ShowCmdConfig{
					Controllers: []Controller{
						{Hostname: "test-wnc.example.com", AccessToken: "test-token"},
					},
					PrintFormat: PrintFormatJSON,
					Timeout:     30,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Fatalf("Failed to marshal config to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledConfig Config
			err = json.Unmarshal(jsonData, &unmarshaledConfig)
			if err != nil {
				t.Fatalf("Failed to unmarshal config from JSON: %v", err)
			}

			// Verify round-trip integrity
			if len(unmarshaledConfig.ShowCmdConfig.Controllers) != len(tt.config.ShowCmdConfig.Controllers) {
				t.Errorf("Controller count mismatch after JSON round-trip")
			}

			if unmarshaledConfig.ShowCmdConfig.PrintFormat != tt.config.ShowCmdConfig.PrintFormat {
				t.Errorf("PrintFormat mismatch after JSON round-trip")
			}
		})
	}
}

func TestConfigConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"ControllersFlagName", ControllersFlagName, "controllers"},
		{"AllowInsecureAccessFlagName", AllowInsecureAccessFlagName, "insecure"},
		{"PrintFormatFlagName", PrintFormatFlagName, "format"},
		{"TimeoutFlagName", TimeoutFlagName, "timeout"},
		{"RadioFlagName", RadioFlagName, "radio"},
		{"SSIDFlagName", SSIDFlagName, "ssid"},
		{"SortByFlagName", SortByFlagName, "sort-by"},
		{"SortOrderFlagName", SortOrderFlagName, "sort-order"},
		{"APNameFlagName", APNameFlagName, "ap-name"},
		{"PrintFormatJSON", PrintFormatJSON, "json"},
		{"PrintFormatTable", PrintFormatTable, "table"},
		{"OrderByAscending", OrderByAscending, "asc"},
		{"OrderByDescending", OrderByDescending, "desc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Constant %s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestRadioSlotConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{"RadioSlotNumSlot0ID", RadioSlotNumSlot0ID, 0},
		{"RadioSlotNumSlot1ID", RadioSlotNumSlot1ID, 1},
		{"RadioSlotNumSlot2ID", RadioSlotNumSlot2ID, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Radio slot constant %s = %d, want %d", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
