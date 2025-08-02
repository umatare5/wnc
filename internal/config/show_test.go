package config

import (
	"encoding/json"
	"testing"

	"github.com/urfave/cli/v3"
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

// TestParseControllerPair tests parseControllerPair function (Unit test)
func TestParseControllerPair(t *testing.T) {
	cfg := &Config{}

	t.Run("test_basic_controller_pair", func(t *testing.T) {
		hostname, token, err := cfg.parseControllerPair("controller1:token123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hostname != "controller1" {
			t.Errorf("Expected hostname 'controller1', got '%s'", hostname)
		}

		if token != "token123" {
			t.Errorf("Expected token 'token123', got '%s'", token)
		}
	})

	t.Run("test_https_url_controller_pair", func(t *testing.T) {
		hostname, token, err := cfg.parseControllerPair("https://controller.example.com:8443:token456")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hostname != "controller.example.com:8443" {
			t.Errorf("Expected hostname 'controller.example.com:8443', got '%s'", hostname)
		}

		if token != "token456" {
			t.Errorf("Expected token 'token456', got '%s'", token)
		}
	})

	t.Run("test_http_url_controller_pair", func(t *testing.T) {
		hostname, token, err := cfg.parseControllerPair("http://controller.local:9090:abc123")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hostname != "controller.local:9090" {
			t.Errorf("Expected hostname 'controller.local:9090', got '%s'", hostname)
		}

		if token != "abc123" {
			t.Errorf("Expected token 'abc123', got '%s'", token)
		}
	})

	t.Run("test_invalid_format_no_colon", func(t *testing.T) {
		_, _, err := cfg.parseControllerPair("controller_without_colon")

		if err == nil {
			t.Error("Expected error for invalid format, got nil")
		}

		expectedMsg := "invalid controllers format: controllers does not contain ':'."
		if err.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	t.Run("test_whitespace_handling", func(t *testing.T) {
		hostname, token, err := cfg.parseControllerPair("  controller1  :  token123  ")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hostname != "controller1" {
			t.Errorf("Expected trimmed hostname 'controller1', got '%s'", hostname)
		}

		if token != "token123" {
			t.Errorf("Expected trimmed token 'token123', got '%s'", token)
		}
	})

	t.Run("test_complex_url_with_multiple_colons", func(t *testing.T) {
		hostname, token, err := cfg.parseControllerPair("https://user:pass@controller.com:8443:finaltoken")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if hostname != "user:pass@controller.com:8443" {
			t.Errorf("Expected hostname 'user:pass@controller.com:8443', got '%s'", hostname)
		}

		if token != "finaltoken" {
			t.Errorf("Expected token 'finaltoken', got '%s'", token)
		}
	})
}

// TestSetShowCmdConfig tests SetShowCmdConfig method with CLI
func TestSetShowCmdConfig(t *testing.T) {
	t.Run("valid_config", func(t *testing.T) {
		// Create mock cli command
		cmd := &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  ControllersFlagName,
					Value: "controller1.com:token1",
				},
				&cli.StringFlag{
					Name:  PrintFormatFlagName,
					Value: "table",
				},
				&cli.BoolFlag{
					Name:  AllowInsecureAccessFlagName,
					Value: false,
				},
				&cli.IntFlag{
					Name:  TimeoutFlagName,
					Value: 30,
				},
			},
		}

		// Set flag values
		cmd.Set(ControllersFlagName, "controller1.com:token1")
		cmd.Set(PrintFormatFlagName, "table")
		cmd.Set(AllowInsecureAccessFlagName, "false")
		cmd.Set(TimeoutFlagName, "30")

		cfg := &Config{}
		cfg.SetShowCmdConfig(cmd)

		if cfg.ShowCmdConfig.PrintFormat != "table" {
			t.Errorf("Expected PrintFormat 'table', got '%s'", cfg.ShowCmdConfig.PrintFormat)
		}
	})
}
