package framework

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestNewGenerateCli(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		usecase    *application.Usecase
		wantNil    bool
	}{
		{
			name:       "creates_new_generate_cli_with_valid_dependencies",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			usecase:    &application.Usecase{},
			wantNil:    false,
		},
		{
			name:       "creates_new_generate_cli_with_nil_config",
			config:     nil,
			repository: &infrastructure.Repository{},
			usecase:    &application.Usecase{},
			wantNil:    false,
		},
		{
			name:       "creates_new_generate_cli_with_nil_repository",
			config:     &config.Config{},
			repository: nil,
			usecase:    &application.Usecase{},
			wantNil:    false,
		},
		{
			name:       "creates_new_generate_cli_with_nil_usecase",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			usecase:    nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewGenerateCli(tt.config, tt.repository, tt.usecase)

			// Check that the CLI was created (value type, so never nil)
			if tt.config != nil && cli.Config != tt.config {
				t.Error("Expected config to match provided config")
			}
			if tt.repository != nil && cli.Repository != tt.repository {
				t.Error("Expected repository to match provided repository")
			}
			if tt.usecase != nil && cli.Usecase != tt.usecase {
				t.Error("Expected usecase to match provided usecase")
			}
		})
	}
}

func TestGenerateCliInvokeSubClis(t *testing.T) {
	tests := []struct {
		name   string
		cli    GenerateCli
		method string
	}{
		{
			name: "InvokeTokenCli_returns_TokenCli",
			cli: GenerateCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			method: "InvokeTokenCli",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.method {
			case "InvokeTokenCli":
				tokenCli := tt.cli.InvokeTokenCli()
				if tokenCli == nil {
					t.Error("Expected non-nil TokenCli")
					return
				}
				if tokenCli.Config != tt.cli.Config {
					t.Error("Expected TokenCli config to match GenerateCli config")
				}
				if tokenCli.Repository != tt.cli.Repository {
					t.Error("Expected TokenCli repository to match GenerateCli repository")
				}
				if tokenCli.Usecase != tt.cli.Usecase {
					t.Error("Expected TokenCli usecase to match GenerateCli usecase")
				}
			}
		})
	}
}

func TestGenerateCliJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		cli  GenerateCli
	}{
		{
			name: "empty_generate_cli",
			cli: GenerateCli{
				Config:     nil,
				Repository: nil,
				Usecase:    nil,
			},
		},
		{
			name: "generate_cli_with_dependencies",
			cli: GenerateCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.cli)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var cli GenerateCli
			err = json.Unmarshal(jsonData, &cli)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if tt.cli.Config == nil && cli.Config != nil {
				t.Error("Expected nil config after unmarshal")
			}
			if tt.cli.Repository == nil && cli.Repository != nil {
				t.Error("Expected nil repository after unmarshal")
			}
			if tt.cli.Usecase == nil && cli.Usecase != nil {
				t.Error("Expected nil usecase after unmarshal")
			}
		})
	}
}

func TestGenerateCliDependencyInjection(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		repository     *infrastructure.Repository
		usecase        *application.Usecase
		expectValidCli bool
	}{
		{
			name: "dependency_injection_with_valid_objects",
			config: &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: "testuser",
					Password: "testpass",
				},
			},
			repository:     &infrastructure.Repository{},
			usecase:        &application.Usecase{},
			expectValidCli: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := NewGenerateCli(tt.config, tt.repository, tt.usecase)

			if tt.expectValidCli {
				if cli.Config == nil {
					t.Error("Expected valid config in CLI")
				}
				if cli.Repository == nil {
					t.Error("Expected valid repository in CLI")
				}
				if cli.Usecase == nil {
					t.Error("Expected valid usecase in CLI")
				}

				// Test that the CLI can invoke sub-CLIs
				tokenCli := cli.InvokeTokenCli()
				if tokenCli == nil {
					t.Error("Expected non-nil TokenCli from valid CLI")
				}
			}
		})
	}
}

func TestGenerateCliFailFast(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		usecase    *application.Usecase
	}{
		{
			name:       "nil_dependencies_should_not_panic_in_constructor",
			config:     nil,
			repository: nil,
			usecase:    nil,
		},
		{
			name:       "valid_dependencies_should_not_panic",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			usecase:    &application.Usecase{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("GenerateCli operation panicked: %v", r)
				}
			}()

			cli := NewGenerateCli(tt.config, tt.repository, tt.usecase)

			// Test that basic operations don't panic
			tokenCli := cli.InvokeTokenCli()
			if tokenCli == nil {
				t.Error("Expected non-nil TokenCli")
			}
		})
	}
}
