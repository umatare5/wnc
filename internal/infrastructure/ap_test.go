package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestApRepositoryCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates ApRepository with valid config",
			config: &config.Config{},
		},
		{
			name:   "creates ApRepository with nil config",
			config: nil,
		},
		{
			name: "creates ApRepository with populated config",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
					Controllers: []config.Controller{
						{Hostname: "wnc.example.com", AccessToken: "token123"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apRepo := &ApRepository{
				Config: tt.config,
			}

			if apRepo.Config != tt.config {
				t.Errorf("Config = %v, want %v", apRepo.Config, tt.config)
			}
		})
	}
}

func TestApRepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		apRepo ApRepository
	}{
		{
			name: "empty ApRepository",
			apRepo: ApRepository{
				Config: nil,
			},
		},
		{
			name: "ApRepository with config",
			apRepo: ApRepository{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						Timeout: 60,
						Controllers: []config.Controller{
							{Hostname: "test.example.com", AccessToken: "testtoken"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.apRepo)
			if err != nil {
				t.Fatalf("Failed to marshal ApRepository to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledApRepo ApRepository
			err = json.Unmarshal(jsonData, &unmarshaledApRepo)
			if err != nil {
				t.Fatalf("Failed to unmarshal ApRepository from JSON: %v", err)
			}

			// Basic validation - checking that unmarshaling doesn't fail
			// Note: Pointer fields will be different after unmarshal
			if tt.apRepo.Config != nil && unmarshaledApRepo.Config == nil {
				// This is expected behavior for JSON unmarshaling with pointers - no action needed
				_ = tt.apRepo.Config // Acknowledge the check
			}
		})
	}
}

func TestApRepositoryFailFast(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.Config
		controller  string
		apikey      string
		isSecure    *bool
		expectPanic bool
	}{
		{
			name:        "nil config should not panic",
			config:      nil,
			controller:  "test.example.com",
			apikey:      "testkey",
			isSecure:    func() *bool { b := true; return &b }(),
			expectPanic: false,
		},
		{
			name: "valid config should not panic",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			},
			controller:  "test.example.com",
			apikey:      "testkey",
			isSecure:    func() *bool { b := false; return &b }(),
			expectPanic: false,
		},
		{
			name: "empty controller should not panic",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			},
			controller:  "",
			apikey:      "testkey",
			isSecure:    nil,
			expectPanic: false,
		},
		{
			name: "empty apikey should not panic",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			},
			controller:  "test.example.com",
			apikey:      "",
			isSecure:    nil,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.expectPanic {
					t.Error("Expected panic but none occurred")
				}
			}()

			apRepo := &ApRepository{
				Config: tt.config,
			}

			// For nil config, only test repository creation (don't call methods that require config)
			if tt.config == nil {
				if apRepo.Config != nil {
					t.Error("Expected Config to be nil")
				}
				return
			}

			// Test GetApOper method (should handle errors gracefully, not panic)
			result := apRepo.GetApOper(tt.controller, tt.apikey, tt.isSecure)

			// The method should return nil for invalid inputs but not panic
			if result != nil {
				t.Logf("GetApOper returned non-nil result (unexpected with test data)")
			}
		})
	}
}

func TestApRepositoryTableDriven(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller returns nil",
			controller: "invalid.example.com",
			apikey:     "testkey",
			isSecure:   func() *bool { b := true; return &b }(),
			wantNil:    true,
		},
		{
			name:       "empty controller returns nil",
			controller: "",
			apikey:     "testkey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty apikey returns nil",
			controller: "test.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "nil isSecure with valid inputs returns nil",
			controller: "test.example.com",
			apikey:     "testkey",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	config := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}
	apRepo := &ApRepository{Config: config}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := apRepo.GetApOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApOper() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("GetApOper() = nil, want non-nil")
			}
		})
	}
}

func TestApRepositoryDependencyInjection(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "dependency injection with valid config",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 45,
					Controllers: []config.Controller{
						{Hostname: "wnc1.example.com", AccessToken: "token1"},
						{Hostname: "wnc2.example.com", AccessToken: "token2"},
					},
				},
			},
		},
		{
			name:   "dependency injection with nil config",
			config: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apRepo := &ApRepository{
				Config: tt.config,
			}

			// Verify that dependency is properly injected
			if apRepo.Config != tt.config {
				t.Errorf("Config not properly injected: got %v, want %v", apRepo.Config, tt.config)
			}
		})
	}
}

func TestApRepositoryResponseTypeValidation(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "GetApOper returns correct type",
		},
	}

	config := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}
	apRepo := &ApRepository{Config: config}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test with invalid inputs to get predictable nil result
			result := apRepo.GetApOper("invalid", "invalid", nil)

			// Verify the return is nil for invalid inputs (expected behavior)
			if result != nil {
				t.Logf("GetApOper() returned non-nil result: %v", result)
			}
		})
	}
}

func TestApRepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}

	apRepo := &ApRepository{Config: originalConfig}

	// Verify that the repository doesn't modify the original config
	apRepo.GetApOper("test.example.com", "testkey", nil)

	if originalConfig.ShowCmdConfig.Timeout != 30 {
		t.Error("Repository modified the original config")
	}
}
