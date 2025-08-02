package utils

import (
	"os"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// RequireControllerEnv checks for required environment variables for integration tests
func RequireControllerEnv(t *testing.T) {
	t.Helper()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if !HasTestControllers() {
		t.Skip("Skipping integration test - WNC_CONTROLLERS not set or WNC_BASIC_AUTH_TOKEN not set")
	}
}

// HasTestControllers checks if test controllers are configured
func HasTestControllers() bool {
	controllers := os.Getenv("WNC_CONTROLLERS")
	token := os.Getenv("WNC_BASIC_AUTH_TOKEN")
	return controllers != "" && token != ""
}

// GetTestControllers returns test controllers from environment variables
func GetTestControllers() []config.Controller {
	controllersEnv := os.Getenv("WNC_CONTROLLERS")
	if controllersEnv == "" {
		return []config.Controller{}
	}

	var controllers []config.Controller
	for _, controllerStr := range strings.Split(controllersEnv, ",") {
		parts := strings.Split(strings.TrimSpace(controllerStr), ":")
		if len(parts) >= 2 {
			controllers = append(controllers, config.Controller{
				Hostname:    parts[0] + ":" + parts[1],
				AccessToken: os.Getenv("WNC_BASIC_AUTH_TOKEN"),
			})
		}
	}

	return controllers
}

// GetTestToken returns the test token from environment variables
func GetTestToken() string {
	return os.Getenv("WNC_BASIC_AUTH_TOKEN")
}

// SetupTestEnv sets up common test environment variables for unit tests
func SetupTestEnv(t *testing.T) func() {
	t.Helper()

	// Save original values
	originalToken := os.Getenv("WNC_BASIC_AUTH_TOKEN")
	originalControllers := os.Getenv("WNC_CONTROLLERS")

	// Set test values if not already set
	if originalToken == "" {
		os.Setenv("WNC_BASIC_AUTH_TOKEN", "test-token")
	}
	if originalControllers == "" {
		os.Setenv("WNC_CONTROLLERS", "test.example.com:443")
	}

	// Return cleanup function
	return func() {
		if originalToken == "" {
			os.Unsetenv("WNC_BASIC_AUTH_TOKEN")
		} else {
			os.Setenv("WNC_BASIC_AUTH_TOKEN", originalToken)
		}

		if originalControllers == "" {
			os.Unsetenv("WNC_CONTROLLERS")
		} else {
			os.Setenv("WNC_CONTROLLERS", originalControllers)
		}
	}
}

// IsIntegrationTest returns true if integration test environment is available
func IsIntegrationTest() bool {
	return os.Getenv("WNC_BASIC_AUTH_TOKEN") != "" && os.Getenv("WNC_CONTROLLERS") != ""
}

// RequireEnvVar requires a specific environment variable to be set
func RequireEnvVar(t *testing.T, envVar string) string {
	t.Helper()

	value := os.Getenv(envVar)
	if value == "" {
		t.Skipf("Environment variable %s is required for this test", envVar)
	}

	return value
}
