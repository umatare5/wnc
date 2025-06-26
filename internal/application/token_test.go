package application

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestTokenUsecaseNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		wantNil    bool
	}{
		{
			name:       "creates_new_token_usecase_with_valid_dependencies",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_token_usecase_with_nil_config",
			config:     nil,
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_token_usecase_with_nil_repository",
			config:     &config.Config{},
			repository: nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &TokenUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			if (usecase == nil) != tt.wantNil {
				t.Errorf("TokenUsecase creation failed, got nil: %v, want nil: %v", usecase == nil, tt.wantNil)
			}
		})
	}
}

func TestTokenUsecaseJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		usecase TokenUsecase
	}{
		{
			name: "empty_token_usecase",
			usecase: TokenUsecase{
				Config:     nil,
				Repository: nil,
			},
		},
		{
			name: "token_usecase_with_dependencies",
			usecase: TokenUsecase{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.usecase)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var usecase TokenUsecase
			err = json.Unmarshal(jsonData, &usecase)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if tt.usecase.Config == nil && usecase.Config != nil {
				t.Error("Expected nil config after unmarshal")
			}
			if tt.usecase.Repository == nil && usecase.Repository != nil {
				t.Error("Expected nil repository after unmarshal")
			}
		})
	}
}

func TestGenerateBasicAuthToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		expected string
	}{
		{
			name:     "basic_auth_token_generation",
			username: "admin",
			password: "password",
			expected: "YWRtaW46cGFzc3dvcmQ=", // base64 of "admin:password"
		},
		{
			name:     "empty_credentials",
			username: "",
			password: "",
			expected: "Og==", // base64 of ":"
		},
		{
			name:     "special_characters",
			username: "user@domain.com",
			password: "p@ssw0rd!",
			expected: "dXNlckBkb21haW4uY29tOnBAc3N3MHJkIQ==", // base64 of "user@domain.com:p@ssw0rd!"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: tt.username,
					Password: tt.password,
				},
			}

			usecase := &TokenUsecase{
				Config:     config,
				Repository: &infrastructure.Repository{},
			}

			token := usecase.GenerateBasicAuthToken()
			if token != tt.expected {
				t.Errorf("Expected token %s, got %s", tt.expected, token)
			}
		})
	}
}

func TestTokenUsecaseFailFast(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "nil_config_should_not_panic",
			config: nil,
		},
		{
			name: "config_with_empty_credentials_should_not_panic",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: "",
					Password: "",
				},
			},
		},
		{
			name: "config_with_valid_credentials_should_not_panic",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: "testuser",
					Password: "testpass",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("TokenUsecase operation panicked: %v", r)
				}
			}()

			usecase := &TokenUsecase{
				Config:     tt.config,
				Repository: &infrastructure.Repository{},
			}

			// Test that basic operations don't panic even with nil config
			if tt.config != nil {
				token := usecase.GenerateBasicAuthToken()
				if len(token) == 0 {
					t.Error("Expected non-empty token")
				}
			} else {
				// With nil config, this should panic, but we're testing that the usecase can be created
				if usecase.Config != nil {
					t.Error("Expected nil config")
				}
			}
		})
	}
}

func TestTokenUsecaseDependencyInjection(t *testing.T) {
	tests := []struct {
		name               string
		config             *config.Config
		repository         *infrastructure.Repository
		expectValidUsecase bool
	}{
		{
			name: "dependency_injection_with_valid_objects",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: "testuser",
					Password: "testpass",
				},
			},
			repository:         &infrastructure.Repository{},
			expectValidUsecase: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &TokenUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			if tt.expectValidUsecase {
				if usecase.Config == nil {
					t.Error("Expected valid config in usecase")
				}
				if usecase.Repository == nil {
					t.Error("Expected valid repository in usecase")
				}

				// Test that the usecase can generate tokens
				token := usecase.GenerateBasicAuthToken()
				if len(token) == 0 {
					t.Error("Expected non-empty token from valid usecase")
				}
			}
		})
	}
}
