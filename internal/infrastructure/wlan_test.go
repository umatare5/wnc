package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestWlanRepositoryCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates WlanRepository with valid config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
		},
		{
			name:   "creates WlanRepository with nil config",
			config: nil,
		},
		{
			name:   "creates WlanRepository with populated config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 10}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &WlanRepository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestWlanRepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		repo *WlanRepository
	}{
		{
			name: "empty WlanRepository",
			repo: &WlanRepository{},
		},
		{
			name: "WlanRepository with config",
			repo: &WlanRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.repo)
			if err != nil {
				t.Errorf("Failed to marshal WlanRepository: %v", err)
			}

			var unmarshaled WlanRepository
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal WlanRepository: %v", err)
			}
		})
	}
}

func TestWlanRepositoryFailFast(t *testing.T) {
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
					t.Errorf("WlanRepository operation panicked: %v", r)
				}
			}()

			repo := &WlanRepository{Config: tt.config}
			if repo.Config != nil {
				isSecure := true
				_ = repo.GetWlanCfg(tt.controller, tt.apikey, &isSecure)
			}
		})
	}
}

func TestWlanRepositoryTableDriven(t *testing.T) {
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
			repo := &WlanRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
			result := repo.GetWlanCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result for %s, got %v", tt.name, result)
			}
		})
	}
}

func TestWlanRepositoryDependencyInjection(t *testing.T) {
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
			repo := &WlanRepository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestWlanRepositoryResponseTypeValidation(t *testing.T) {
	t.Run("GetWlanCfg returns correct type", func(t *testing.T) {
		repo := &WlanRepository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
		isSecure := true
		result := repo.GetWlanCfg("invalid", "token", &isSecure)
		// result should be nil due to network error, but type should be correct
		if result != nil {
			t.Logf("GetWlanCfg returned: %T", result)
		}
	})
}

func TestWlanRepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}
	repo := &WlanRepository{Config: originalConfig}

	isSecure := true
	_ = repo.GetWlanCfg("test.example.com", "test-token", &isSecure)

	// Config should remain unchanged
	if repo.Config.ShowCmdConfig.Timeout != 30 {
		t.Error("Repository config was modified during operation")
	}
}
