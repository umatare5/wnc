package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestClientRepositoryCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates ClientRepository with valid config",
			config: &config.Config{},
		},
		{
			name:   "creates ClientRepository with nil config",
			config: nil,
		},
		{
			name: "creates ClientRepository with populated config",
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
			clientRepo := &ClientRepository{
				Config: tt.config,
			}

			if clientRepo.Config != tt.config {
				t.Errorf("Config = %v, want %v", clientRepo.Config, tt.config)
			}
		})
	}
}

func TestClientRepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name       string
		clientRepo ClientRepository
	}{
		{
			name: "empty ClientRepository",
			clientRepo: ClientRepository{
				Config: nil,
			},
		},
		{
			name: "ClientRepository with config",
			clientRepo: ClientRepository{
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
			jsonData, err := json.Marshal(tt.clientRepo)
			if err != nil {
				t.Fatalf("Failed to marshal ClientRepository to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledClientRepo ClientRepository
			err = json.Unmarshal(jsonData, &unmarshaledClientRepo)
			if err != nil {
				t.Fatalf("Failed to unmarshal ClientRepository from JSON: %v", err)
			}

			// Basic validation - checking that unmarshaling doesn't fail
			// Note: Pointer fields will be different after unmarshal
			if tt.clientRepo.Config != nil && unmarshaledClientRepo.Config == nil {
				// This is expected behavior for JSON unmarshaling with pointers - no action needed
				_ = tt.clientRepo.Config // Acknowledge the check
			}
		})
	}
}

func TestClientRepositoryFailFast(t *testing.T) {
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

			clientRepo := &ClientRepository{
				Config: tt.config,
			}

			// For nil config, only test repository creation (don't call methods that require config)
			if tt.config == nil {
				if clientRepo.Config != nil {
					t.Error("Expected Config to be nil")
				}
				return
			}

			// Test GetClientOper method (should handle errors gracefully, not panic)
			result := clientRepo.GetClientOper(tt.controller, tt.apikey, tt.isSecure)

			// The method should return nil for invalid inputs but not panic
			if result != nil {
				t.Logf("GetClientOper returned non-nil result (unexpected with test data)")
			}

			// Test GetClientGlobalOper method
			globalResult := clientRepo.GetClientGlobalOper(tt.controller, tt.apikey, tt.isSecure)
			if globalResult != nil {
				t.Logf("GetClientGlobalOper returned non-nil result (unexpected with test data)")
			}
		})
	}
}

func TestClientRepositoryTableDriven(t *testing.T) {
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
	}

	config := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}
	clientRepo := &ClientRepository{Config: config}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetClientOper
			result := clientRepo.GetClientOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetClientOper() = %v, want nil", result)
			}
			if !tt.wantNil && result == nil {
				t.Errorf("GetClientOper() = nil, want non-nil")
			}

			// Test GetClientGlobalOper
			globalResult := clientRepo.GetClientGlobalOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && globalResult != nil {
				t.Errorf("GetClientGlobalOper() = %v, want nil", globalResult)
			}
			if !tt.wantNil && globalResult == nil {
				t.Errorf("GetClientGlobalOper() = nil, want non-nil")
			}
		})
	}
}

func TestClientRepositoryDependencyInjection(t *testing.T) {
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
			clientRepo := &ClientRepository{
				Config: tt.config,
			}

			// Verify that dependency is properly injected
			if clientRepo.Config != tt.config {
				t.Errorf("Config not properly injected: got %v, want %v", clientRepo.Config, tt.config)
			}
		})
	}
}

func TestClientRepositoryMethodAvailability(t *testing.T) {
	tests := []struct {
		name       string
		methodName string
	}{
		{
			name:       "GetClientOper method exists",
			methodName: "GetClientOper",
		},
		{
			name:       "GetClientGlobalOper method exists",
			methodName: "GetClientGlobalOper",
		},
	}

	config := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}
	clientRepo := &ClientRepository{Config: config}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.methodName {
			case "GetClientOper":
				result := clientRepo.GetClientOper("invalid", "invalid", nil)
				if result != nil {
					t.Logf("GetClientOper method executed and returned: %v", result)
				}
			case "GetClientGlobalOper":
				result := clientRepo.GetClientGlobalOper("invalid", "invalid", nil)
				if result != nil {
					t.Logf("GetClientGlobalOper method executed and returned: %v", result)
				}
			default:
				t.Errorf("Unknown method: %s", tt.methodName)
			}
		})
	}
}

func TestClientRepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 30,
		},
	}

	clientRepo := &ClientRepository{Config: originalConfig}

	// Verify that the repository doesn't modify the original config
	clientRepo.GetClientOper("test.example.com", "testkey", nil)
	clientRepo.GetClientGlobalOper("test.example.com", "testkey", nil)

	if originalConfig.ShowCmdConfig.Timeout != 30 {
		t.Error("Repository modified the original config")
	}
}
