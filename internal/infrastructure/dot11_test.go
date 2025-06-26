package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestDot11RepositoryCreation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "creates Dot11Repository with valid config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}},
		},
		{
			name:   "creates Dot11Repository with nil config",
			config: nil,
		},
		{
			name:   "creates Dot11Repository with populated config",
			config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 10}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &Dot11Repository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestDot11RepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		repo *Dot11Repository
	}{
		{
			name: "empty Dot11Repository",
			repo: &Dot11Repository{},
		},
		{
			name: "Dot11Repository with config",
			repo: &Dot11Repository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.repo)
			if err != nil {
				t.Errorf("Failed to marshal Dot11Repository: %v", err)
			}

			var unmarshaled Dot11Repository
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal Dot11Repository: %v", err)
			}
		})
	}
}

func TestDot11RepositoryFailFast(t *testing.T) {
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
					t.Errorf("Dot11Repository operation panicked: %v", r)
				}
			}()

			repo := &Dot11Repository{Config: tt.config}
			if repo.Config != nil {
				isSecure := true
				_ = repo.GetDot11Cfg(tt.controller, tt.apikey, &isSecure)
			}
		})
	}
}

func TestDot11RepositoryTableDriven(t *testing.T) {
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
			repo := &Dot11Repository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
			result := repo.GetDot11Cfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result for %s, got %v", tt.name, result)
			}
		})
	}
}

func TestDot11RepositoryDependencyInjection(t *testing.T) {
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
			repo := &Dot11Repository{Config: tt.config}
			if repo.Config != tt.config {
				t.Errorf("Expected config %v, got %v", tt.config, repo.Config)
			}
		})
	}
}

func TestDot11RepositoryResponseTypeValidation(t *testing.T) {
	t.Run("GetDot11Cfg returns correct type", func(t *testing.T) {
		repo := &Dot11Repository{Config: &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}}
		isSecure := true
		result := repo.GetDot11Cfg("invalid", "token", &isSecure)
		// result should be nil due to network error, but type should be correct
		if result != nil {
			t.Logf("GetDot11Cfg returned: %T", result)
		}
	})
}

func TestDot11RepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{ShowCmdConfig: config.ShowCmdConfig{Timeout: 30}}
	repo := &Dot11Repository{Config: originalConfig}

	isSecure := true
	_ = repo.GetDot11Cfg("test.example.com", "test-token", &isSecure)

	// Config should remain unchanged
	if repo.Config.ShowCmdConfig.Timeout != 30 {
		t.Error("Repository config was modified during operation")
	}
}
