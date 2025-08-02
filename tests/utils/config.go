package utils

import (
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// CreateTestConfig creates a minimal test configuration
func CreateTestConfig() *config.Config {
	cfg := config.New()
	return &cfg
}

// CreateTestConfigWithControllers creates a test config with sample controllers
func CreateTestConfigWithControllers() *config.Config {
	cfg := CreateTestConfig()
	cfg.ShowCmdConfig.Controllers = []config.Controller{
		{
			Hostname:    "controller1.example.com",
			AccessToken: "token123",
		},
		{
			Hostname:    "controller2.example.com",
			AccessToken: "token456",
		},
	}
	return cfg
}

// CreateTestRepository creates a test repository with the given config
func CreateTestRepository(cfg *config.Config) *infrastructure.Repository {
	if cfg == nil {
		cfg = CreateTestConfig()
	}
	repo := infrastructure.New(cfg)
	return &repo
}

// CreateTestUsecase creates a test usecase with the given config and repository
func CreateTestUsecase(cfg *config.Config, repo *infrastructure.Repository) *application.Usecase {
	if cfg == nil {
		cfg = CreateTestConfig()
	}
	if repo == nil {
		repo = CreateTestRepository(cfg)
	}
	usecase := application.New(cfg, repo)
	return &usecase
}

// CreateFullTestStack creates a complete test stack (config, repository, usecase)
func CreateFullTestStack() (*config.Config, *infrastructure.Repository, *application.Usecase) {
	cfg := CreateTestConfig()
	repo := CreateTestRepository(cfg)
	usecase := CreateTestUsecase(cfg, repo)
	return cfg, repo, usecase
}

// CreateFullTestStackWithControllers creates a complete test stack with sample controllers
func CreateFullTestStackWithControllers() (*config.Config, *infrastructure.Repository, *application.Usecase) {
	cfg := CreateTestConfigWithControllers()
	repo := CreateTestRepository(cfg)
	usecase := CreateTestUsecase(cfg, repo)
	return cfg, repo, usecase
}

// ValidateConfigFields validates that a config has all required fields
func ValidateConfigFields(t *testing.T, cfg *config.Config) {
	t.Helper()

	if cfg == nil {
		t.Fatal("Config is nil")
	}

	// Validate ShowCmdConfig
	AssertStructFields(t, &cfg.ShowCmdConfig,
		"Controllers", "AllowInsecureAccess", "PrintFormat", "Timeout")

	// Validate GenerateCmdConfig
	AssertStructFields(t, &cfg.GenerateCmdConfig,
		"Username", "Password")
}

// ValidateRepositoryFields validates that a repository has all required fields
func ValidateRepositoryFields(t *testing.T, repo *infrastructure.Repository) {
	t.Helper()

	if repo == nil {
		t.Fatal("Repository is nil")
	}

	AssertStructFields(t, repo, "Config")
	AssertNonNilFields(t, repo, "Config")
}

// ValidateUsecaseFields validates that a usecase has all required fields
func ValidateUsecaseFields(t *testing.T, usecase *application.Usecase) {
	t.Helper()

	if usecase == nil {
		t.Fatal("Usecase is nil")
	}

	AssertStructFields(t, usecase, "Config", "Repository")
	AssertNonNilFields(t, usecase, "Config", "Repository")
}

// TestPrintFormats provides all valid print formats for testing
func TestPrintFormats() []string {
	return []string{
		config.PrintFormatTable,
		config.PrintFormatJSON,
	}
}

// TestTimeouts provides various timeout values for testing
func TestTimeouts() []int {
	return []int{
		1,   // Minimum
		30,  // Default
		60,  // Extended
		300, // Maximum
	}
}

// SetupTestController creates a test controller with given parameters
func SetupTestController(hostname, token string) config.Controller {
	return config.Controller{
		Hostname:    hostname,
		AccessToken: token,
	}
}

// CreateControllersFromEnv creates controllers from environment variable format
func CreateControllersFromEnv(envValue string) []config.Controller {
	controllers := []config.Controller{}

	if envValue == "" {
		return controllers
	}

	// Parse environment variable format: "host1:token1,host2:token2"
	pairs := strings.Split(envValue, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			controllers = append(controllers, config.Controller{
				Hostname:    strings.TrimSpace(parts[0]),
				AccessToken: strings.TrimSpace(parts[1]),
			})
		}
	}

	return controllers
}

// MockRepositoryMethods contains common mock setup for repository methods
type MockRepositoryMethods struct {
	GetApOper          func() interface{}
	GetApCapwapData    func() interface{}
	GetApLldpNeigh     func() interface{}
	GetApRadioOperData func() interface{}
	GetApOperData      func() interface{}
	GetClientOper      func() interface{}
	GetWlanCfg         func() interface{}
}

// CreateMockRepository creates a mock repository for testing
func CreateMockRepository(methods *MockRepositoryMethods) *infrastructure.Repository {
	cfg := CreateTestConfig()
	repo := infrastructure.New(cfg)

	// In a real implementation, this would set up actual mocks
	// For now, return a basic repository
	return &repo
}

// ValidateControllerConnection validates that a controller config is valid
func ValidateControllerConnection(t *testing.T, controller config.Controller) {
	t.Helper()

	if controller.Hostname == "" {
		t.Error("Controller hostname should not be empty")
	}

	if controller.AccessToken == "" {
		t.Error("Controller access token should not be empty")
	}

	// Basic URL format validation
	if !strings.Contains(controller.Hostname, ".") {
		t.Errorf("Controller hostname %q should be a valid domain", controller.Hostname)
	}
}
