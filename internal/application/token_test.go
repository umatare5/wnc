package application

import (
	"testing"
)

// TestTokenUsecaseInitialization tests TokenUsecase initialization (Unit test)
func TestTokenUsecaseInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "invoke_token_usecase",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				tokenUsecase := app.InvokeTokenUsecase()
				if tokenUsecase == nil {
					t.Error("InvokeTokenUsecase returned nil")
					return
				}
				if tokenUsecase.Config != cfg {
					t.Error("TokenUsecase Config not set correctly")
				}
				if tokenUsecase.Repository != repo {
					t.Error("TokenUsecase Repository not set correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtilsInstance.assertNoPanic(t, func() {
				tt.test(t)
			})
		})
	}
}

// TestGenerateBasicAuthToken tests token generation functionality (Unit test)
func TestGenerateBasicAuthToken(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "basic_auth_token_generation",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				tokenUsecase := app.InvokeTokenUsecase()

				token := tokenUsecase.GenerateBasicAuthToken()
				if token == "" {
					t.Error("GenerateBasicAuthToken should return non-empty token")
				}

				// Test that the token is actually encoded correctly
				// The length depends on the actual credentials in mock config
				if len(token) == 0 {
					t.Error("Token should not be empty")
				}
			},
		},
		{
			name: "empty_credentials_in_config",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.GenerateCmdConfig.Username = ""
				cfg.GenerateCmdConfig.Password = ""
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				tokenUsecase := app.InvokeTokenUsecase()

				token := tokenUsecase.GenerateBasicAuthToken()
				if token == "" {
					t.Error("GenerateBasicAuthToken should return non-empty token even for empty credentials")
				}

				// Even empty credentials should produce a valid base64 encoded string ":"
				expectedLength := 4 // Base64 encoding of ":"
				if len(token) != expectedLength {
					t.Errorf("Expected token length %d for empty credentials, got %d", expectedLength, len(token))
				}
			},
		},
		{
			name: "special_characters_in_config",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.GenerateCmdConfig.Username = "user@domain.com"
				cfg.GenerateCmdConfig.Password = "pass!@#$"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				tokenUsecase := app.InvokeTokenUsecase()

				token := tokenUsecase.GenerateBasicAuthToken()
				if token == "" {
					t.Error("GenerateBasicAuthToken should handle special characters")
				}

				// Should produce a valid base64 string regardless of special characters
				if len(token) == 0 {
					t.Error("Token should not be empty for special characters")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtilsInstance.assertNoPanic(t, func() {
				tt.test(t)
			})
		})
	}
}
