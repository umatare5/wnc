package config

import (
	"encoding/json"
	"reflect"
	"testing"
)

// validateConfigStructFields validates that a config has all required fields
func validateConfigStructFields(t *testing.T, cfg *Config) {
	t.Helper()

	if cfg == nil {
		t.Fatal("Config is nil")
	}

	// Check that struct has expected fields using reflection
	configType := reflect.TypeOf(*cfg)
	expectedFields := []string{"GenerateCmdConfig", "ShowCmdConfig"}

	for _, fieldName := range expectedFields {
		if _, found := configType.FieldByName(fieldName); !found {
			t.Errorf("Expected field %q not found in Config struct", fieldName)
		}
	}
}

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

			// Use helper for validation
			validateConfigStructFields(t, &got)

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

// TestConstants tests all package constants (Unit test)
func TestConstants(t *testing.T) {
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
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// TestRadioConstants tests radio slot constants (Unit test)
func TestRadioConstants(t *testing.T) {
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
				t.Errorf("Expected %s to be %d, got %d", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

// TestConfigStructure tests the Config struct structure (Unit test)
func TestConfigStructure(t *testing.T) {
	t.Run("config has required fields", func(t *testing.T) {
		config := Config{}

		// Test that fields exist and are of correct type
		var _ GenerateCmdConfig = config.GenerateCmdConfig
		var _ ShowCmdConfig = config.ShowCmdConfig
	})

	t.Run("config initialization", func(t *testing.T) {
		config := New()

		// Verify the structure is properly initialized
		if config.GenerateCmdConfig != (GenerateCmdConfig{}) {
			t.Error("GenerateCmdConfig should be zero value")
		}
		if config.ShowCmdConfig.Controllers != nil {
			t.Error("ShowCmdConfig.Controllers should be nil")
		}
	})
}

// TestConfigEdgeCases tests edge cases for the Config package (Unit test)
func TestConfigEdgeCases(t *testing.T) {
	t.Run("multiple new instances", func(t *testing.T) {
		config1 := New()
		config2 := New()

		// Each call to New() should return independent instances
		if &config1 == &config2 {
			t.Error("New() should return different instances")
		}

		// But they should have the same values
		if config1.GenerateCmdConfig != config2.GenerateCmdConfig {
			t.Error("New() instances should have identical GenerateCmdConfig")
		}
		if config1.ShowCmdConfig.PrintFormat != config2.ShowCmdConfig.PrintFormat {
			t.Error("New() instances should have identical ShowCmdConfig fields")
		}
	})

	t.Run("config modification isolation", func(t *testing.T) {
		config1 := New()
		config2 := New()

		// Modify one config
		config1.ShowCmdConfig.PrintFormat = "modified"

		// The other should remain unchanged
		if config2.ShowCmdConfig.PrintFormat == "modified" {
			t.Error("Config instances should be independent")
		}
	})
}

// TestControllerParsing tests controller string parsing (Unit test)
func TestControllerParsing(t *testing.T) {
	cfg := New()

	tests := []struct {
		name     string
		input    string
		expected int
		wantErr  bool
	}{
		{
			name:     "single controller",
			input:    "host1:token1",
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "multiple controllers",
			input:    "host1:token1,host2:token2",
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "empty input",
			input:    "",
			expected: 0,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controllers := cfg.parseControllers(tt.input)
			if len(controllers) != tt.expected {
				t.Errorf("Expected %d controllers, got %d", tt.expected, len(controllers))
			}
		})
	}
}

// TestPrintFormatValidation tests print format validation (Unit test)
func TestPrintFormatValidation(t *testing.T) {
	cfg := New()

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{
			name:    "valid json format",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "valid table format",
			format:  "table",
			wantErr: false,
		},
		{
			name:    "invalid format",
			format:  "invalid",
			wantErr: true,
		},
		{
			name:    "empty format",
			format:  "",
			wantErr: true, // empty format should be invalid in validation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.validatePrintFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("validatePrintFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestControllerValidation tests controller format validation (Unit test)
func TestControllerValidation(t *testing.T) {
	cfg := New()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid single controller",
			input:   "host1:token1",
			wantErr: false,
		},
		{
			name:    "valid multiple controllers",
			input:   "host1:token1,host2:token2",
			wantErr: false,
		},
		{
			name:    "invalid format no colon",
			input:   "host1token1",
			wantErr: true,
		},
		{
			name:    "invalid format empty hostname",
			input:   ":token1",
			wantErr: true,
		},
		{
			name:    "invalid format empty token",
			input:   "host1:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cfg.validateControllersFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateControllersFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigConstants tests additional constants (Unit test)
func TestConfigConstants(t *testing.T) {
	headerTests := []struct {
		name     string
		constant string
		expected string
	}{
		{"OverviewHeaderApMac", OverviewHeaderApMac, "APMac"},
		{"OverviewHeaderApRadioID", OverviewHeaderApRadioID, "Radio"},
		{"OverviewHeaderApOperStatus", OverviewHeaderApOperStatus, "Status"},
		{"OverviewHeaderChannelNumber", OverviewHeaderChannelNumber, "Channel"},
		{"ShowClientHeaderBand", ShowClientHeaderBand, "Band"},
		{"ShowCommonHeaderApName", ShowCommonHeaderApName, "APName"},
		{"ShowCommonHeaderController", ShowCommonHeaderController, "Controller"},
	}

	for _, tt := range headerTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.constant)
			}
		})
	}

	cmdConstTests := []struct {
		name     string
		constant string
		expected string
	}{
		{"UsernameFlagName", UsernameFlagName, "username"},
		{"PasswordFlagName", PasswordFlagName, "password"},
	}

	for _, tt := range cmdConstTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s to be %q, got %q", tt.name, tt.expected, tt.constant)
			}
		})
	}
}
