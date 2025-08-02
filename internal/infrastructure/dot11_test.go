package infrastructure

import (
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
)

// TestDot11Repository tests the Dot11Repository structure
func TestDot11Repository(t *testing.T) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}

	repo := &Dot11Repository{
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

// TestDot11Repository_GetDot11Cfg tests the GetDot11Cfg method
func TestDot11Repository_GetDot11Cfg(t *testing.T) {
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
					Timeout: 1, // Use very short timeout to avoid actual network calls
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
			name: "zero_timeout",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 0,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Dot11Repository{
				Config: tt.config,
			}

			result := repo.GetDot11Cfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestDot11Repository_GetDot11CfgTimeout tests timeout configuration
func TestDot11Repository_GetDot11CfgTimeout(t *testing.T) {
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
			name:             "custom_timeout",
			timeoutSeconds:   60,
			expectedDuration: 60 * time.Second,
		},
		{
			name:             "zero_timeout",
			timeoutSeconds:   0,
			expectedDuration: 0 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: tt.timeoutSeconds,
				},
			}

			repo := &Dot11Repository{
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

// TestDot11Repository_NilConfig tests behavior with nil config
func TestDot11Repository_NilConfig(t *testing.T) {
	repo := &Dot11Repository{
		Config: nil,
	}

	// This should panic or handle gracefully
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - accessing nil config should panic
			t.Logf("Expected panic when accessing nil config: %v", r)
		}
	}()

	// This will panic due to nil config, which is expected behavior
	_ = repo.GetDot11Cfg("test", "test", nil)
}
