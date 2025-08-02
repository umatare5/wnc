package generate

import (
	"encoding/base64"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// isValidBase64 checks if a string is valid base64
func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

// TestTokenCli tests the TokenCli structure
func TestTokenCli(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create_token_cli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{Config: cfg}
			usecase := &application.Usecase{Config: cfg, Repository: repo}

			tokenCli := &TokenCli{
				Config:     cfg,
				Repository: repo,
				Usecase:    usecase,
			}

			// Test that TokenCli is properly initialized
			if tokenCli.Config != cfg {
				t.Error("Expected config to be set")
			}
			if tokenCli.Repository != repo {
				t.Error("Expected repository to be set")
			}
			if tokenCli.Usecase != usecase {
				t.Error("Expected usecase to be set")
			}
		})
	}
}

// TestGenerateToken tests the GenerateToken method
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		expectContains string
	}{
		{
			name: "generate_basic_auth_token",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: "admin",
					Password: "password",
				},
			},
			expectContains: "YWRtaW46", // Base64 encoded "admin:"
		},
		{
			name: "generate_token_empty_config",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{},
			},
			expectContains: "Og==", // Base64 encoded ":"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create TokenCli with test config
			repo := &infrastructure.Repository{Config: tt.config}
			usecase := &application.Usecase{Config: tt.config, Repository: repo}
			tokenCli := &TokenCli{
				Config:     tt.config,
				Repository: repo,
				Usecase:    usecase,
			}

			// Execute GenerateToken
			tokenCli.GenerateToken()

			// Restore stdout and read output
			w.Close()
			os.Stdout = originalStdout

			output, _ := io.ReadAll(r)
			outputStr := string(output)

			// Verify output contains expected content
			if !strings.Contains(outputStr, tt.expectContains) {
				t.Errorf("Expected output to contain '%s', got: %s", tt.expectContains, outputStr)
			}
		})
	}
}

// TestGenerateTokenOutput tests the token generation output format
func TestGenerateTokenOutput(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
	}{
		{
			name:     "admin_credentials",
			username: "admin",
			password: "admin123",
		},
		{
			name:     "test_credentials",
			username: "test",
			password: "test123",
		},
		{
			name:     "empty_credentials",
			username: "",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Setup test configuration
			cfg := &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: tt.username,
					Password: tt.password,
				},
			}

			repo := &infrastructure.Repository{Config: cfg}
			usecase := &application.Usecase{Config: cfg, Repository: repo}
			tokenCli := &TokenCli{
				Config:     cfg,
				Repository: repo,
				Usecase:    usecase,
			}

			// Execute token generation
			tokenCli.GenerateToken()

			// Restore stdout and capture output
			w.Close()
			os.Stdout = originalStdout

			output, _ := io.ReadAll(r)
			outputStr := strings.TrimSpace(string(output))

			// Verify output format - should be base64 encoded
			if len(outputStr) == 0 {
				t.Error("Expected non-empty output")
			}

			// Verify output is valid base64
			if !isValidBase64(outputStr) {
				t.Errorf("Expected valid base64 output, got: %s", outputStr)
			}
		})
	}
}

// TestTokenCliStructure tests the overall TokenCli structure
func TestTokenCliStructure(t *testing.T) {
	t.Run("token_cli_fields", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{}
		usecase := &application.Usecase{}

		tokenCli := &TokenCli{
			Config:     cfg,
			Repository: repo,
			Usecase:    usecase,
		}

		// Test field assignments
		if tokenCli.Config == nil {
			t.Error("TokenCli.Config should not be nil")
		}
		if tokenCli.Repository == nil {
			t.Error("TokenCli.Repository should not be nil")
		}
		if tokenCli.Usecase == nil {
			t.Error("TokenCli.Usecase should not be nil")
		}
	})

	t.Run("token_cli_method_exists", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}
		usecase := &application.Usecase{Config: cfg, Repository: repo}
		tokenCli := &TokenCli{
			Config:     cfg,
			Repository: repo,
			Usecase:    usecase,
		}

		// Test that GenerateToken method exists and can be called
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GenerateToken method panicked: %v", r)
			}
		}()

		// Redirect stdout to avoid polluting test output
		originalStdout := os.Stdout
		os.Stdout, _ = os.Open(os.DevNull)
		defer func() { os.Stdout = originalStdout }()

		tokenCli.GenerateToken()
	})
}

// TestTokenCliDependencyInjection tests dependency injection
func TestTokenCliDependencyInjection(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "config_dependency",
			test: func(t *testing.T) {
				cfg := &config.Config{
					GenerateCmdConfig: config.GenerateCmdConfig{
						Username: "test",
						Password: "test",
					},
				}
				repo := &infrastructure.Repository{Config: cfg}
				usecase := &application.Usecase{Config: cfg, Repository: repo}
				tokenCli := &TokenCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				if tokenCli.Config.GenerateCmdConfig.Username != "test" {
					t.Error("Config dependency not properly injected")
				}
			},
		},
		{
			name: "repository_dependency",
			test: func(t *testing.T) {
				cfg := &config.Config{}
				repo := &infrastructure.Repository{Config: cfg}
				usecase := &application.Usecase{Config: cfg, Repository: repo}
				tokenCli := &TokenCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				if tokenCli.Repository.Config != cfg {
					t.Error("Repository dependency not properly injected")
				}
			},
		},
		{
			name: "usecase_dependency",
			test: func(t *testing.T) {
				cfg := &config.Config{}
				repo := &infrastructure.Repository{Config: cfg}
				usecase := &application.Usecase{Config: cfg, Repository: repo}
				tokenCli := &TokenCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				if tokenCli.Usecase.Config != cfg {
					t.Error("Usecase dependency not properly injected")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}
