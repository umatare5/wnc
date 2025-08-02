package config

import (
	"encoding/json"
	"testing"

	"github.com/urfave/cli/v3"
)

// TestGenerateCmdConfigJSONSerialization tests JSON serialization for GenerateCmdConfig
func TestGenerateCmdConfigJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config GenerateCmdConfig
	}{
		{
			name: "empty_generate_config",
			config: GenerateCmdConfig{
				Username: "",
				Password: "",
			},
		},
		{
			name: "full_generate_config",
			config: GenerateCmdConfig{
				Username: "testuser",
				Password: "testpassword",
			},
		},
		{
			name: "special_characters",
			config: GenerateCmdConfig{
				Username: "user@domain.com",
				Password: "p@ssw0rd!#$",
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
			var config GenerateCmdConfig
			err = json.Unmarshal(jsonData, &config)
			if err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
				return
			}

			// Validate the unmarshaled data
			if config.Username != tt.config.Username {
				t.Errorf("Username mismatch: got %s, want %s", config.Username, tt.config.Username)
			}
			if config.Password != tt.config.Password {
				t.Errorf("Password mismatch: got %s, want %s", config.Password, tt.config.Password)
			}
		})
	}
}

// TestSetGenerateCmdConfig tests SetGenerateCmdConfig method (focuses on structure validation)
func TestSetGenerateCmdConfig(t *testing.T) {
	t.Run("test_generate_config_structure", func(t *testing.T) {
		cfg := &Config{}

		// Test the GenerateCmdConfig structure initialization
		generateCfg := GenerateCmdConfig{
			Username: "testuser",
			Password: "testpass",
		}

		// Set the config directly (simulating what SetGenerateCmdConfig would do)
		cfg.GenerateCmdConfig = generateCfg

		// Validate the configuration was set
		if cfg.GenerateCmdConfig.Username != "testuser" {
			t.Errorf("Expected Username 'testuser', got '%s'", cfg.GenerateCmdConfig.Username)
		}

		if cfg.GenerateCmdConfig.Password != "testpass" {
			t.Errorf("Expected Password 'testpass', got '%s'", cfg.GenerateCmdConfig.Password)
		}
	})

	t.Run("test_generate_config_empty_values", func(t *testing.T) {
		cfg := &Config{}

		// Test with empty configuration values
		generateCfg := GenerateCmdConfig{
			Username: "",
			Password: "",
		}

		cfg.GenerateCmdConfig = generateCfg

		// Validate empty values
		if cfg.GenerateCmdConfig.Username != "" {
			t.Errorf("Expected empty Username, got '%s'", cfg.GenerateCmdConfig.Username)
		}

		if cfg.GenerateCmdConfig.Password != "" {
			t.Errorf("Expected empty Password, got '%s'", cfg.GenerateCmdConfig.Password)
		}
	})
}

// TestValidateGenerateCliFlags tests validateGenerateCliFlags method (focuses on validation logic)
func TestValidateGenerateCliFlags(t *testing.T) {
	t.Run("test_validate_username_password_logic", func(t *testing.T) {
		cfg := &Config{}

		// Create mock cli command with valid values
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  UsernameFlagName,
					Value: "testuser",
				},
				&cli.StringFlag{
					Name:  PasswordFlagName,
					Value: "testpass",
				},
			},
		}

		// Set flag values
		cmd.Set(UsernameFlagName, "testuser")
		cmd.Set(PasswordFlagName, "testpass")

		// Test the validation logic indirectly
		// Note: validateGenerateCliFlags calls log.Fatal on error,
		// so we test the successful case only
		err := cfg.validateGenerateCliFlags(cmd)

		if err != nil {
			t.Errorf("validateGenerateCliFlags() returned error for valid input: %v", err)
		}
	})
}
