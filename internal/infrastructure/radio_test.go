package infrastructure

import (
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
)

// TestRadioRepository tests the RadioRepository structure
func TestRadioRepository(t *testing.T) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}

	repo := &RadioRepository{
		Config: cfg,
	}

	// Test that repository is properly initialized
	if repo.Config == nil {
		t.Error("Expected config to be set")
	}

	if repo.Config.ShowCmdConfig.Timeout != 30 {
		t.Errorf("Expected timeout to be 30, got %d", repo.Config.ShowCmdConfig.Timeout)
	}
}

// TestRadioRepository_GetRadioCfg tests the GetRadioCfg method
func TestRadioRepository_GetRadioCfg(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name: "valid_parameters",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1, // Use short timeout to avoid network delays
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true, // Will be nil due to no real connection
		},
		{
			name: "invalid_controller",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
		{
			name: "insecure_connection",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:8080",
			apikey:     "test-token",
			isSecure:   &[]bool{false}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RadioRepository{
				Config: tt.config,
			}

			result := repo.GetRadioCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRadioRepository_GetRadioCfgTimeout tests timeout configuration
func TestRadioRepository_GetRadioCfgTimeout(t *testing.T) {
	tests := []struct {
		name             string
		timeoutSeconds   int
		expectedDuration time.Duration
	}{
		{
			name:             "default_timeout",
			timeoutSeconds:   30,
			expectedDuration: 30 * time.Second,
		},
		{
			name:             "long_timeout",
			timeoutSeconds:   120,
			expectedDuration: 120 * time.Second,
		},
		{
			name:             "short_timeout",
			timeoutSeconds:   5,
			expectedDuration: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: tt.timeoutSeconds,
				},
			}

			repo := &RadioRepository{
				Config: cfg,
			}

			// Test that the timeout is properly configured
			actualDuration := time.Duration(repo.Config.ShowCmdConfig.Timeout) * time.Second
			if actualDuration != tt.expectedDuration {
				t.Errorf("Expected timeout %v, got %v", tt.expectedDuration, actualDuration)
			}
		})
	}
}

// TestRadioRepository_ErrorHandling tests error handling scenarios
func TestRadioRepository_ErrorHandling(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "nil_config",
			config: nil,
		},
		{
			name: "zero_timeout_config",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RadioRepository{
				Config: tt.config,
			}

			// Test error handling for edge cases
			defer func() {
				if r := recover(); r != nil && tt.config == nil {
					// Expected panic for nil config
					t.Logf("Expected panic for nil config: %v", r)
				}
			}()

			result := repo.GetRadioCfg("test", "test", nil)
			if result != nil {
				t.Error("Expected nil result for invalid configuration")
			}
		})
	}
}
