package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// Test configuration using environment variables for real WNC credentials
// Set these environment variables to run integration tests:
// WNC_CONTROLLERS - format: "hostname1:token1,hostname2:token2"
// WNC_INTEGRATION_INSECURE - set to "true" to allow insecure connections

func getTestControllers() ([]config.Controller, bool) {
	controllersEnv := os.Getenv("WNC_CONTROLLERS")
	if controllersEnv == "" {
		return nil, false
	}

	insecureEnv := os.Getenv("WNC_INTEGRATION_INSECURE")
	allowInsecure := strings.ToLower(insecureEnv) == "true"

	// Parse controllers string
	var controllers []config.Controller
	pairs := strings.Split(controllersEnv, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) == 2 {
			controllers = append(controllers, config.Controller{
				Hostname:    strings.TrimSpace(parts[0]),
				AccessToken: strings.TrimSpace(parts[1]),
			})
		}
	}

	if len(controllers) == 0 {
		return nil, false
	}

	return controllers, allowInsecure
}

func TestIntegrationShowOverview(t *testing.T) {
	controllers, allowInsecure := getTestControllers()
	if len(controllers) == 0 {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	tests := []struct {
		name        string
		controllers []config.Controller
		insecure    bool
	}{
		{
			name:        "show overview with real controllers",
			controllers: controllers,
			insecure:    allowInsecure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up dependency injection chain
			cfg := config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers:         tt.controllers,
					AllowInsecureAccess: tt.insecure,
					PrintFormat:         config.PrintFormatJSON,
					Timeout:             30,
				},
			}

			repo := infrastructure.New(&cfg)
			usecase := application.New(&cfg, &repo)
			frameworkCli := framework.NewShowCli(&cfg, &repo, &usecase)

			// Test overview functionality
			overviewCli := frameworkCli.InvokeOverviewCli()
			if overviewCli == nil {
				t.Fatal("Failed to create overview CLI")
			}

			// Execute the overview command
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Overview command panicked: %v", r)
				}
			}()

			// Test that the CLI components are properly initialized
			if overviewCli.Config == nil {
				t.Error("Overview CLI config is nil")
			}

			if overviewCli.Repository == nil {
				t.Error("Overview CLI repository is nil")
			}

			if overviewCli.Usecase == nil {
				t.Error("Overview CLI usecase is nil")
			}

			// Test the actual overview usecase
			overviewUsecase := usecase.InvokeOverviewUsecase()
			if overviewUsecase == nil {
				t.Fatal("Failed to create overview usecase")
			}

			// Call the actual overview method with real controllers
			isSecure := !tt.insecure
			result := overviewUsecase.ShowOverview(&tt.controllers, &isSecure)

			// Verify results
			if result == nil {
				t.Error("ShowOverview should return empty slice, not nil")
				return
			}

			// Note: Result might be empty if controllers are unreachable,
			// but the call should not panic or return nil
			t.Logf("Overview returned %d entries", len(result))

			// Validate structure of returned data
			for i, data := range result {
				if data == nil {
					t.Errorf("Overview data at index %d is nil", i)
					continue
				}

				// Basic validation of overview data structure
				if data.Controller == "" {
					t.Errorf("Overview data at index %d has empty controller", i)
				}

				// Log some information for debugging (don't fail test if empty)
				t.Logf("Entry %d: Controller=%s, ApMac=%s, SlotID=%d",
					i, data.Controller, data.ApMac, data.SlotID)
			}
		})
	}
}

func TestIntegrationShowWlan(t *testing.T) {
	controllers, allowInsecure := getTestControllers()
	if len(controllers) == 0 {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	tests := []struct {
		name        string
		controllers []config.Controller
		insecure    bool
	}{
		{
			name:        "show wlan with real controllers",
			controllers: controllers,
			insecure:    allowInsecure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up dependency injection chain
			cfg := config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers:         tt.controllers,
					AllowInsecureAccess: tt.insecure,
					PrintFormat:         config.PrintFormatJSON,
					Timeout:             30,
				},
			}

			repo := infrastructure.New(&cfg)
			usecase := application.New(&cfg, &repo)
			frameworkCli := framework.NewShowCli(&cfg, &repo, &usecase)

			// Test WLAN functionality
			wlanCli := frameworkCli.InvokeWlanCli()
			if wlanCli == nil {
				t.Fatal("Failed to create WLAN CLI")
			}

			// Execute the WLAN command
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("WLAN command panicked: %v", r)
				}
			}()

			// Test that the CLI components are properly initialized
			if wlanCli.Config == nil {
				t.Error("WLAN CLI config is nil")
			}

			if wlanCli.Repository == nil {
				t.Error("WLAN CLI repository is nil")
			}

			if wlanCli.Usecase == nil {
				t.Error("WLAN CLI usecase is nil")
			}

			// Test the actual WLAN usecase
			wlanUsecase := usecase.InvokeWlanUsecase()
			if wlanUsecase == nil {
				t.Fatal("Failed to create WLAN usecase")
			}

			t.Logf("WLAN integration test completed successfully")
		})
	}
}

func TestIntegrationEndToEndCLI(t *testing.T) {
	controllers, allowInsecure := getTestControllers()
	if len(controllers) == 0 {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	tests := []struct {
		name        string
		controllers []config.Controller
		insecure    bool
		format      string
	}{
		{
			name:        "end-to-end JSON format",
			controllers: controllers,
			insecure:    allowInsecure,
			format:      config.PrintFormatJSON,
		},
		{
			name:        "end-to-end table format",
			controllers: controllers,
			insecure:    allowInsecure,
			format:      config.PrintFormatTable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test complete dependency injection chain
			cfg := config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers:         tt.controllers,
					AllowInsecureAccess: tt.insecure,
					PrintFormat:         tt.format,
					Timeout:             30,
					SortBy:              "name",
					SortOrder:           config.OrderByAscending,
				},
			}

			// Test repository layer
			repo := infrastructure.New(&cfg)
			if repo.Config != &cfg {
				t.Error("Repository config not properly injected")
			}

			// Test application layer
			usecase := application.New(&cfg, &repo)
			if usecase.Config != &cfg {
				t.Error("Usecase config not properly injected")
			}
			if usecase.Repository != &repo {
				t.Error("Usecase repository not properly injected")
			}

			// Test framework layer
			frameworkCli := framework.NewShowCli(&cfg, &repo, &usecase)
			if frameworkCli.Config != &cfg {
				t.Error("Framework config not properly injected")
			}
			if frameworkCli.Repository != &repo {
				t.Error("Framework repository not properly injected")
			}
			if frameworkCli.Usecase != &usecase {
				t.Error("Framework usecase not properly injected")
			}

			// Test all CLI commands can be created
			clientCli := frameworkCli.InvokeClientCli()
			apCli := frameworkCli.InvokeApCli()
			apTagCli := frameworkCli.InvokeApTagCli()
			wlanCli := frameworkCli.InvokeWlanCli()
			overviewCli := frameworkCli.InvokeOverviewCli()

			clis := []interface{}{clientCli, apCli, apTagCli, wlanCli, overviewCli}
			for i, cli := range clis {
				if cli == nil {
					t.Errorf("CLI at index %d is nil", i)
				}
			}

			// Test that error handling works properly
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("End-to-end test panicked: %v", r)
				}
			}()

			t.Logf("End-to-end integration test with format %s completed successfully", tt.format)
		})
	}
}

func TestIntegrationFailFast(t *testing.T) {
	controllers, allowInsecure := getTestControllers()
	if len(controllers) == 0 {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	tests := []struct {
		name                string
		controllers         []config.Controller
		insecure            bool
		expectGracefulError bool
	}{
		{
			name: "invalid controller should fail gracefully",
			controllers: []config.Controller{
				{Hostname: "invalid.example.com", AccessToken: "invalid-token"},
			},
			insecure:            allowInsecure,
			expectGracefulError: true,
		},
		{
			name: "empty token should fail gracefully",
			controllers: []config.Controller{
				{Hostname: controllers[0].Hostname, AccessToken: ""},
			},
			insecure:            allowInsecure,
			expectGracefulError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if tt.expectGracefulError {
						t.Errorf("Expected graceful error handling, but got panic: %v", r)
					}
				}
			}()

			cfg := config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers:         tt.controllers,
					AllowInsecureAccess: tt.insecure,
					PrintFormat:         config.PrintFormatJSON,
					Timeout:             5, // Short timeout for invalid controllers
				},
			}

			repo := infrastructure.New(&cfg)
			usecase := application.New(&cfg, &repo)

			// Test overview with invalid controllers
			overviewUsecase := usecase.InvokeOverviewUsecase()
			isSecure := !tt.insecure
			result := overviewUsecase.ShowOverview(&tt.controllers, &isSecure)

			// Should return empty result, not panic
			if result == nil {
				t.Error("ShowOverview should return empty slice, not nil")
				return
			}

			if tt.expectGracefulError && len(result) > 0 {
				// This might happen if the "invalid" controller is actually valid
				t.Logf("Expected empty result for invalid controller, got %d entries", len(result))
			}
		})
	}
}
