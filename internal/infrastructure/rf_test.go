package infrastructure

import (
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
)

// TestRfRepository tests the RfRepository structure
func TestRfRepository(t *testing.T) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 1, // Use short timeout for unit tests
		},
	}

	repo := &RfRepository{
		Config: cfg,
	}

	// Test that repository is properly initialized
	if repo.Config == nil {
		t.Error("Expected config to be set")
	}

	if repo.Config.ShowCmdConfig.Timeout != 1 {
		t.Errorf("Expected timeout to be 1, got %d", repo.Config.ShowCmdConfig.Timeout)
	}
}

// TestRfRepository_GetRfTags tests the GetRfTags method
func TestRfRepository_GetRfTags(t *testing.T) {
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
					Timeout: 1,
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
			name: "empty_apikey",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RfRepository{
				Config: tt.config,
			}

			result := repo.GetRfTags(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRfRepository_GetRfTagsTimeout tests timeout configuration
func TestRfRepository_GetRfTagsTimeout(t *testing.T) {
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
			name:             "extended_timeout",
			timeoutSeconds:   180,
			expectedDuration: 180 * time.Second,
		},
		{
			name:             "minimal_timeout",
			timeoutSeconds:   1,
			expectedDuration: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: tt.timeoutSeconds,
				},
			}

			repo := &RfRepository{
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

// TestRfRepository_EdgeCases tests edge cases and error scenarios
func TestRfRepository_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.Config
		setupFunc func() *RfRepository
	}{
		{
			name:   "nil_config",
			config: nil,
			setupFunc: func() *RfRepository {
				return &RfRepository{Config: nil}
			},
		},
		{
			name: "negative_timeout",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: -1,
				},
			},
			setupFunc: func() *RfRepository {
				return &RfRepository{
					Config: &config.Config{
						ShowCmdConfig: config.ShowCmdConfig{
							Timeout: -1,
						},
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupFunc()

			// Test error handling for edge cases
			defer func() {
				if r := recover(); r != nil && tt.config == nil {
					// Expected panic for nil config
					t.Logf("Expected panic for nil config: %v", r)
				}
			}()

			result := repo.GetRfTags("test", "test", nil)
			if result != nil {
				t.Error("Expected nil result for invalid configuration")
			}
		})
	}
}
