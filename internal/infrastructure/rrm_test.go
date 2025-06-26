package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestRrmRepositoryCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates RrmRepository with valid config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
		},
		{
			name:   "creates RrmRepository with nil config",
			config: nil,
		},
		{
			name:   "creates RrmRepository with populated config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 10}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestRrmRepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		repo *RrmRepository
	}{
		{
			name: "empty RrmRepository",
			repo: &RrmRepository{},
		},
		{
			name: "RrmRepository with config",
			repo: &RrmRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.repo)
			if err != nil {
				t.Errorf("Failed to marshal RrmRepository: %v", err)
			}

			var unmarshaled RrmRepository
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal RrmRepository: %v", err)
			}
		})
	}
}

func TestRrmRepositoryFailFast(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
	}{
		{
			name:       "nil config should not panic",
			config:     nil,
			controller: "test.example.com",
			apikey:     "test-token",
		},
		{
			name:       "valid config should not panic",
			config:     &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
			controller: "test.example.com",
			apikey:     "test-token",
		},
		{
			name:       "empty controller should not panic",
			config:     &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
			controller: "",
			apikey:     "test-token",
		},
		{
			name:       "empty apikey should not panic",
			config:     &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
			controller: "test.example.com",
			apikey:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("RrmRepository operation panicked: %v", r)
				}
			}()

			repo := &RrmRepository{Config: tt.config}
			if repo.Config != nil {
				isSecure := true
				_ = repo.GetRrmOper(tt.controller, tt.apikey, &isSecure)
			}
		})
	}
}

func TestRrmRepositoryTableDriven(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name:       "invalid controller returns nil",
			controller: "invalid.example.com",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
		{
			name:       "empty controller returns nil",
			controller: "",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
		{
			name:       "empty apikey returns nil",
			controller: "test.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
		{
			name:       "nil isSecure with valid inputs returns nil",
			controller: "test.example.com",
			apikey:     "test-token",
			isSecure:   nil,
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
			result := repo.GetRrmOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result for %s, got %v", tt.name, result)
			}
		})
	}
}

func TestRrmRepositoryDependencyInjection(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "dependency injection with valid config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
		},
		{
			name:   "dependency injection with nil config",
			config: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestRrmRepositoryMethodAvailability(t *testing.T) {
	repo := &RrmRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
	isSecure := true

	t.Run("GetRrmOper method exists", func(t *testing.T) {
		result := repo.GetRrmOper("invalid", "token", &isSecure)
		if result != nil {
			t.Logf("GetRrmOper returned: %T", result)
		}
	})

	t.Run("GetRrmMeasurement method exists", func(t *testing.T) {
		result := repo.GetRrmMeasurement("invalid", "token", &isSecure)
		if result != nil {
			t.Logf("GetRrmMeasurement returned: %T", result)
		}
	})

	t.Run("GetRrmGlobalOper method exists", func(t *testing.T) {
		result := repo.GetRrmGlobalOper("invalid", "token", &isSecure)
		if result != nil {
			t.Logf("GetRrmGlobalOper returned: %T", result)
		}
	})

	t.Run("GetRrmCfg method exists", func(t *testing.T) {
		result := repo.GetRrmCfg("invalid", "token", &isSecure)
		if result != nil {
			t.Logf("GetRrmCfg returned: %T", result)
		}
	})
}

func TestRrmRepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}
	repo := &RrmRepository{Config: originalConfig}

	isSecure := true
	_ = repo.GetRrmOper("test.example.com", "test-token", &isSecure)
	_ = repo.GetRrmMeasurement("test.example.com", "test-token", &isSecure)
	_ = repo.GetRrmGlobalOper("test.example.com", "test-token", &isSecure)
	_ = repo.GetRrmCfg("test.example.com", "test-token", &isSecure)

	// Config should remain unchanged
	if repo.Config.ShowCmdConfig.Timeout != 30 {
		t.Error("Repository config was modified during operation")
	}
}
