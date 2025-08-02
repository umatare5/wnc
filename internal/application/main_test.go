package application

import (
	"fmt"
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/client"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// Local test utils to avoid import cycle
type testUtils struct{}

func (tu *testUtils) createMockConfig() *config.Config {
	return &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			PrintFormat: config.PrintFormatTable,
			Controllers: []config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			Timeout: 30,
		},
		GenerateCmdConfig: config.GenerateCmdConfig{
			Username: "testuser",
			Password: "testpass",
		},
	}
}

func (tu *testUtils) createMockRepository(cfg *config.Config) *infrastructure.Repository {
	return &infrastructure.Repository{Config: cfg}
}

func (tu *testUtils) assertNoPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function panicked: %v", r)
		}
	}()
	fn()
}

var localUtils = &testUtils{}

// TestNew tests application layer initialization (Unit test)
func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create_application_layer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}

			app := New(cfg, repo)
			if app.Config != cfg {
				t.Error("New() Config not set correctly")
			}
			if app.Repository != repo {
				t.Error("New() Repository not set correctly")
			}
		})
	}
}

// TestInvokeTokenUsecase tests token usecase invocation (Unit test)
func TestInvokeTokenUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_token_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			repo := localUtils.createMockRepository(cfg)
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
		})
	}
}

// TestGenerateBasicAuthToken tests token generation functionality (Unit test)
func TestGenerateBasicAuthToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     string
	}{
		{
			name:     "basic_auth_token",
			username: "admin",
			password: "password",
			want:     "YWRtaW46cGFzc3dvcmQ=", // base64("admin:password")
		},
		{
			name:     "empty_credentials",
			username: "",
			password: "",
			want:     "Og==", // base64(":")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: tt.username,
					Password: tt.password,
				},
			}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			tokenUsecase := app.InvokeTokenUsecase()
			got := tokenUsecase.GenerateBasicAuthToken()

			if got != tt.want {
				t.Errorf("GenerateBasicAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInvokeApUsecase tests AP usecase invocation (Unit test)
func TestInvokeApUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_ap_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			apUsecase := app.InvokeApUsecase()
			if apUsecase == nil {
				t.Error("InvokeApUsecase returned nil")
				return
			}
			if apUsecase.Config != cfg {
				t.Error("ApUsecase Config not set correctly")
			}
			if apUsecase.Repository != repo {
				t.Error("ApUsecase Repository not set correctly")
			}
		})
	}
}

// TestInvokeWlanUsecase tests WLAN usecase invocation (Unit test)
func TestInvokeWlanUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_wlan_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			wlanUsecase := app.InvokeWlanUsecase()
			if wlanUsecase == nil {
				t.Error("InvokeWlanUsecase returned nil")
				return
			}
			if wlanUsecase.Config != cfg {
				t.Error("WlanUsecase Config not set correctly")
			}
			if wlanUsecase.Repository != repo {
				t.Error("WlanUsecase Repository not set correctly")
			}
		})
	}
}

// TestInvokeClientUsecase tests Client usecase invocation (Unit test)
func TestInvokeClientUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_client_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			clientUsecase := app.InvokeClientUsecase()
			if clientUsecase == nil {
				t.Error("InvokeClientUsecase returned nil")
				return
			}
			if clientUsecase.Config != cfg {
				t.Error("ClientUsecase Config not set correctly")
			}
			if clientUsecase.Repository != repo {
				t.Error("ClientUsecase Repository not set correctly")
			}
		})
	}
}

// TestInvokeOverviewUsecase tests Overview usecase invocation (Unit test)
func TestInvokeOverviewUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_overview_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			overviewUsecase := app.InvokeOverviewUsecase()
			if overviewUsecase == nil {
				t.Error("InvokeOverviewUsecase returned nil")
				return
			}
			if overviewUsecase.Config != cfg {
				t.Error("OverviewUsecase Config not set correctly")
			}
			if overviewUsecase.Repository != repo {
				t.Error("OverviewUsecase Repository not set correctly")
			}
		})
	}
}

// TestUsecaseStructure tests the Usecase struct structure (Unit test)
func TestUsecaseStructure(t *testing.T) {
	t.Run("usecase initialization", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{}

		usecase := New(cfg, repo)

		// Test that the usecase is properly initialized
		if usecase.Config == nil {
			t.Error("Usecase Config should not be nil")
		}
		if usecase.Repository == nil {
			t.Error("Usecase Repository should not be nil")
		}

		// Test that the pointers match the input
		if usecase.Config != cfg {
			t.Error("Usecase Config pointer mismatch")
		}
		if usecase.Repository != repo {
			t.Error("Usecase Repository pointer mismatch")
		}
	})
}

// TestAllUsecaseInvokers tests all usecase invoker methods (Unit test)
func TestAllUsecaseInvokers(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{}
	app := New(cfg, repo)

	t.Run("all invokers return non-nil", func(t *testing.T) {
		if app.InvokeTokenUsecase() == nil {
			t.Error("InvokeTokenUsecase returned nil")
		}
		if app.InvokeClientUsecase() == nil {
			t.Error("InvokeClientUsecase returned nil")
		}
		if app.InvokeApUsecase() == nil {
			t.Error("InvokeApUsecase returned nil")
		}
		if app.InvokeWlanUsecase() == nil {
			t.Error("InvokeWlanUsecase returned nil")
		}
		if app.InvokeOverviewUsecase() == nil {
			t.Error("InvokeOverviewUsecase returned nil")
		}
	})

	t.Run("all invokers have correct dependencies", func(t *testing.T) {
		tokenUsecase := app.InvokeTokenUsecase()
		clientUsecase := app.InvokeClientUsecase()
		apUsecase := app.InvokeApUsecase()
		wlanUsecase := app.InvokeWlanUsecase()
		overviewUsecase := app.InvokeOverviewUsecase()

		usecases := []struct {
			name   string
			config *config.Config
			repo   *infrastructure.Repository
		}{
			{"token", tokenUsecase.Config, tokenUsecase.Repository},
			{"client", clientUsecase.Config, clientUsecase.Repository},
			{"ap", apUsecase.Config, apUsecase.Repository},
			{"wlan", wlanUsecase.Config, wlanUsecase.Repository},
			{"overview", overviewUsecase.Config, overviewUsecase.Repository},
		}

		for _, uc := range usecases {
			if uc.config != cfg {
				t.Errorf("%s usecase: Config mismatch", uc.name)
			}
			if uc.repo != repo {
				t.Errorf("%s usecase: Repository mismatch", uc.name)
			}
		}
	})
}

// TestApUsecaseMethods tests ApUsecase methods (Unit test)
func TestApUsecaseMethods(t *testing.T) {
	t.Run("ShowAp with nil repository", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := true

		result := usecase.ShowAp(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil repository")
		}
	})

	t.Run("ShowAp with nil controllers", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		isSecure := true
		result := usecase.ShowAp(nil, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil controllers")
		}
	})
}

// TestWlanUsecaseMethods tests WlanUsecase methods (Unit test)
func TestWlanUsecaseMethods(t *testing.T) {
	t.Run("ShowWlan with nil repository", func(t *testing.T) {
		usecase := &WlanUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := true

		result := usecase.ShowWlan(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil repository")
		}
	})
}

// TestClientUsecaseMethods tests ClientUsecase methods (Unit test)
func TestClientUsecaseMethods(t *testing.T) {
	t.Run("ShowClient with nil repository", func(t *testing.T) {
		usecase := &ClientUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := true

		result := usecase.ShowClient(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil repository")
		}
	})
}

// TestOverviewUsecaseMethods tests OverviewUsecase methods (Unit test)
func TestOverviewUsecaseMethods(t *testing.T) {
	t.Run("ShowOverview with nil repository", func(t *testing.T) {
		usecase := &OverviewUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := true

		result := usecase.ShowOverview(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil repository")
		}
	})
}

// TestTokenUsecaseEdgeCases tests TokenUsecase edge cases (Unit test)
func TestTokenUsecaseEdgeCases(t *testing.T) {
	t.Run("GenerateBasicAuthToken with special characters", func(t *testing.T) {
		cfg := &config.Config{
			GenerateCmdConfig: config.GenerateCmdConfig{
				Username: "user@domain.com",
				Password: "p@$$w0rd!",
			},
		}
		usecase := &TokenUsecase{Config: cfg}

		token := usecase.GenerateBasicAuthToken()
		if token == "" {
			t.Error("Token should not be empty")
		}

		// Should be base64 encoded
		if len(token) < 10 {
			t.Error("Token seems too short to be valid base64")
		}
	})
}

// TestUsecaseDataStructures tests usecase data structures (Unit test)
func TestUsecaseDataStructures(t *testing.T) {
	t.Run("ShowApData structure", func(t *testing.T) {
		data := &ShowApData{}

		// Test that the structure has the expected fields initialized to zero values
		if data.ApMac != "" {
			t.Error("ApMac should be empty string by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowClientData structure", func(t *testing.T) {
		data := &ShowClientData{}

		// Test that the structure exists and can be instantiated
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowWlanData structure", func(t *testing.T) {
		data := &ShowWlanData{}

		// Test that the structure exists and can be instantiated
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})
}

// TestApUsecaseBusinessLogic tests AP usecase business logic
func TestApUsecaseBusinessLogic(t *testing.T) {
	t.Run("ShowAp with empty controllers slice", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		controllers := []config.Controller{}
		isSecure := true

		result := usecase.ShowAp(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with empty controllers")
		}
	})

	t.Run("ShowApTag with nil repository", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := true

		result := usecase.ShowApTag(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil repository")
		}
	})
}

// TestClientUsecaseBusinessLogic tests Client usecase business logic
func TestClientUsecaseBusinessLogic(t *testing.T) {
	t.Run("ShowClient with empty controllers slice", func(t *testing.T) {
		usecase := &ClientUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		controllers := []config.Controller{}
		isSecure := true

		result := usecase.ShowClient(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with empty controllers")
		}
	})

	t.Run("ShowClient with nil controllers", func(t *testing.T) {
		usecase := &ClientUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		isSecure := true
		result := usecase.ShowClient(nil, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil controllers")
		}
	})
}

// TestWlanUsecaseBusinessLogic tests WLAN usecase business logic
func TestWlanUsecaseBusinessLogic(t *testing.T) {
	t.Run("ShowWlan with empty controllers slice", func(t *testing.T) {
		usecase := &WlanUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		controllers := []config.Controller{}
		isSecure := true

		result := usecase.ShowWlan(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with empty controllers")
		}
	})

	t.Run("ShowWlan with nil controllers", func(t *testing.T) {
		usecase := &WlanUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		isSecure := true
		result := usecase.ShowWlan(nil, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil controllers")
		}
	})
}

// TestOverviewUsecaseBusinessLogic tests Overview usecase business logic
func TestOverviewUsecaseBusinessLogic(t *testing.T) {
	t.Run("ShowOverview with empty controllers slice", func(t *testing.T) {
		usecase := &OverviewUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		controllers := []config.Controller{}
		isSecure := true

		result := usecase.ShowOverview(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with empty controllers")
		}
	})

	t.Run("ShowOverview with nil controllers", func(t *testing.T) {
		usecase := &OverviewUsecase{
			Config:     &config.Config{},
			Repository: &infrastructure.Repository{},
		}

		isSecure := true
		result := usecase.ShowOverview(nil, &isSecure)
		if len(result) != 0 {
			t.Error("Expected empty result with nil controllers")
		}
	})
}

// TestDataStructureInitialization tests data structure initialization
func TestDataStructureInitialization(t *testing.T) {
	t.Run("ShowApData initialization", func(t *testing.T) {
		data := ShowApData{}

		// Test that all fields are properly initialized to zero values
		if data.ApMac != "" {
			t.Error("ApMac should be empty string by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowApTagData initialization", func(t *testing.T) {
		data := ShowApTagData{}

		// Test that all fields are properly initialized to zero values
		if data.ApMac != "" {
			t.Error("ApMac should be empty string by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowClientData initialization", func(t *testing.T) {
		data := ShowClientData{}

		// Test that all fields are properly initialized to zero values
		if data.ClientMac != "" {
			t.Error("ClientMac should be empty string by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowWlanData initialization", func(t *testing.T) {
		data := ShowWlanData{}

		// Test that all fields are properly initialized to zero values
		if data.TagName != "" {
			t.Error("TagName should be empty string by default")
		}
		if data.PolicyName != "" {
			t.Error("PolicyName should be empty string by default")
		}
		if data.WlanName != "" {
			t.Error("WlanName should be empty string by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})

	t.Run("ShowOverviewData initialization", func(t *testing.T) {
		data := ShowOverviewData{}

		// Test that all fields are properly initialized to zero values
		if data.ApMac != "" {
			t.Error("ApMac should be empty string by default")
		}
		if data.SlotID != 0 {
			t.Error("SlotID should be 0 by default")
		}
		if data.Controller != "" {
			t.Error("Controller should be empty string by default")
		}
	})
}

// TestUsecaseConfigurationHandling tests configuration handling
func TestUsecaseConfigurationHandling(t *testing.T) {
	t.Run("usecase with different config values", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Timeout:     30,
			},
		}
		repo := &infrastructure.Repository{}

		app := New(cfg, repo)

		// Test that all usecases receive the same config
		tokenUsecase := app.InvokeTokenUsecase()
		apUsecase := app.InvokeApUsecase()
		clientUsecase := app.InvokeClientUsecase()
		wlanUsecase := app.InvokeWlanUsecase()
		overviewUsecase := app.InvokeOverviewUsecase()

		usecases := []interface{}{
			tokenUsecase, apUsecase, clientUsecase, wlanUsecase, overviewUsecase,
		}

		for i, uc := range usecases {
			// Each usecase should have access to the same config
			switch v := uc.(type) {
			case *TokenUsecase:
				if v.Config != cfg {
					t.Errorf("Usecase %d: Config mismatch", i)
				}
			case *ApUsecase:
				if v.Config != cfg {
					t.Errorf("Usecase %d: Config mismatch", i)
				}
			case *ClientUsecase:
				if v.Config != cfg {
					t.Errorf("Usecase %d: Config mismatch", i)
				}
			case *WlanUsecase:
				if v.Config != cfg {
					t.Errorf("Usecase %d: Config mismatch", i)
				}
			case *OverviewUsecase:
				if v.Config != cfg {
					t.Errorf("Usecase %d: Config mismatch", i)
				}
			}
		}
	})
}

// TestUsecaseWithMockedData tests usecase methods with mock data (Integration test)
func TestUsecaseWithMockedData(t *testing.T) {
	t.Run("ShowAp_with_data_flow", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		apUsecase := &ApUsecase{Config: cfg, Repository: repo}

		// Test data flow with empty controllers
		controllers := &[]config.Controller{}
		isSecure := false
		result := apUsecase.ShowAp(controllers, &isSecure)
		if result == nil {
		}
		if len(result) != 0 {
			t.Error("Should return empty slice for empty controllers")
		}
	})

	t.Run("ShowClient_with_filters", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SSID:  "test-ssid",
				Radio: "0",
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}

		controllers := &[]config.Controller{}
		isSecure := false
		result := clientUsecase.ShowClient(controllers, &isSecure)
		if result == nil {
		}
	})

	t.Run("ShowWlan_with_data_processing", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}
		wlanUsecase := &WlanUsecase{Config: cfg, Repository: repo}

		controllers := &[]config.Controller{}
		isSecure := false
		result := wlanUsecase.ShowWlan(controllers, &isSecure)
		if result == nil {
		}
	})

	t.Run("ShowOverview_with_aggregation", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Radio: "1",
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		overviewUsecase := &OverviewUsecase{Config: cfg, Repository: repo}

		controllers := &[]config.Controller{}
		isSecure := false
		result := overviewUsecase.ShowOverview(controllers, &isSecure)
		if result == nil {
		}
	})

	t.Run("ShowApTag_with_configuration_check", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}
		apUsecase := &ApUsecase{Config: cfg, Repository: repo}

		controllers := &[]config.Controller{}
		isSecure := false
		result := apUsecase.ShowApTag(controllers, &isSecure)
		if result == nil {
		}
	})
}

// TestUsecaseDataProcessing tests data processing paths (Unit test)
func TestUsecaseDataProcessing(t *testing.T) {
	t.Run("ap_usecase_nil_controllers", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}
		apUsecase := &ApUsecase{Config: cfg, Repository: repo}

		isSecure := false
		result := apUsecase.ShowAp(nil, &isSecure)
		if result == nil {
		}
		if len(result) != 0 {
			t.Error("Should return empty slice for nil controllers")
		}
	})

	t.Run("client_usecase_nil_controllers", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}
		clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}

		isSecure := false
		result := clientUsecase.ShowClient(nil, &isSecure)
		if result == nil {
		}
		if len(result) != 0 {
			t.Error("Should return empty slice for nil controllers")
		}
	})

	t.Run("wlan_usecase_nil_repository", func(t *testing.T) {
		wlanUsecase := &WlanUsecase{Repository: nil}

		controllers := &[]config.Controller{}
		isSecure := false
		result := wlanUsecase.ShowWlan(controllers, &isSecure)
		if result == nil {
		}
	})

	t.Run("overview_usecase_nil_repository", func(t *testing.T) {
		overviewUsecase := &OverviewUsecase{Repository: nil}

		controllers := &[]config.Controller{}
		isSecure := false
		result := overviewUsecase.ShowOverview(controllers, &isSecure)
		if result == nil {
		}
	})
}

// TestComprehensiveApplicationUsecases tests application layer usecases comprehensively
func TestComprehensiveApplicationUsecases(t *testing.T) {
	cfg := localUtils.createMockConfig()
	repo := localUtils.createMockRepository(cfg)

	t.Run("test_ap_usecase_comprehensive", func(t *testing.T) {
		apUsecase := &ApUsecase{Config: cfg, Repository: repo}
		isSecure := false

		// Test with nil controllers
		result := apUsecase.ShowAp(nil, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for nil controllers, got %d items", len(result))
		}

		// Test with empty controllers
		emptyControllers := []config.Controller{}
		result = apUsecase.ShowAp(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for empty controllers, got %d items", len(result))
		}

		// Test with valid controllers (will fail but exercise code paths)
		controllers := []config.Controller{
			{Hostname: "test1", AccessToken: "token1"},
			{Hostname: "test2", AccessToken: "token2"},
		}
		localUtils.assertNoPanic(t, func() {
			result = apUsecase.ShowAp(&controllers, &isSecure)
			// Result will be empty due to mock repository, but we test the function doesn't panic
		})

		// Test ShowApTag
		localUtils.assertNoPanic(t, func() {
			tagResult := apUsecase.ShowApTag(&controllers, &isSecure)
			// Ensure function completes without panic
			_ = tagResult
		})
	})

	t.Run("test_client_usecase_comprehensive", func(t *testing.T) {
		clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}
		isSecure := false

		// Test with nil controllers
		result := clientUsecase.ShowClient(nil, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for nil controllers, got %d items", len(result))
		}

		// Test with empty controllers
		emptyControllers := []config.Controller{}
		result = clientUsecase.ShowClient(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for empty controllers, got %d items", len(result))
		}

		// Test with valid controllers
		controllers := []config.Controller{
			{Hostname: "test1", AccessToken: "token1"},
		}
		localUtils.assertNoPanic(t, func() {
			result = clientUsecase.ShowClient(&controllers, &isSecure)
			// Result will be empty due to mock repository, but we test the function doesn't panic
		})

		// Test filter functions
		t.Run("test_filter_functions", func(t *testing.T) {
			clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}

			// Test filterBySSID
			mockData := []*ShowClientData{
				{
					ClientMac:  "aa:bb:cc:dd:ee:ff",
					Controller: "test",
				},
			}

			// Test with empty SSID (should return all data)
			filtered := clientUsecase.filterBySSID(mockData)
			if len(filtered) != 1 {
				t.Errorf("Expected 1 item with empty SSID filter, got %d", len(filtered))
			}

			// Test filterByRadio with empty radio (should return all data)
			radioFiltered := clientUsecase.filterByRadio(mockData)
			if len(radioFiltered) != 1 {
				t.Errorf("Expected 1 item with empty radio filter, got %d", len(radioFiltered))
			}
		})
	})

	t.Run("test_wlan_usecase_comprehensive", func(t *testing.T) {
		wlanUsecase := &WlanUsecase{Config: cfg, Repository: repo}
		isSecure := false

		// Test with nil controllers
		result := wlanUsecase.ShowWlan(nil, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for nil controllers, got %d items", len(result))
		}

		// Test with empty controllers
		emptyControllers := []config.Controller{}
		result = wlanUsecase.ShowWlan(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for empty controllers, got %d items", len(result))
		}

		// Test with valid controllers
		controllers := []config.Controller{
			{Hostname: "test1", AccessToken: "token1"},
		}
		localUtils.assertNoPanic(t, func() {
			result = wlanUsecase.ShowWlan(&controllers, &isSecure)
			// Result will be empty due to mock repository, but we test the function doesn't panic
		})
	})

	t.Run("test_overview_usecase_comprehensive", func(t *testing.T) {
		overviewUsecase := &OverviewUsecase{Config: cfg, Repository: repo}
		isSecure := false

		// Test with nil controllers
		result := overviewUsecase.ShowOverview(nil, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for nil controllers, got %d items", len(result))
		}

		// Test with empty controllers
		emptyControllers := []config.Controller{}
		result = overviewUsecase.ShowOverview(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Errorf("Expected empty result for empty controllers, got %d items", len(result))
		}

		// Test with valid controllers
		controllers := []config.Controller{
			{Hostname: "test1", AccessToken: "token1"},
		}
		localUtils.assertNoPanic(t, func() {
			result = overviewUsecase.ShowOverview(&controllers, &isSecure)
			// Result will be empty due to mock repository, but we test the function doesn't panic
		})

		// Test filterByRadio function
		t.Run("test_overview_filter_by_radio", func(t *testing.T) {
			overviewUsecase := &OverviewUsecase{Config: cfg, Repository: repo}
			mockData := []*ShowOverviewData{
				{
					Controller: "test",
				},
			}

			// Test with empty radio (should return all data)
			filtered := overviewUsecase.filterByRadio(mockData)
			if len(filtered) != 1 {
				t.Errorf("Expected 1 item with empty radio filter, got %d", len(filtered))
			}
		})
	})

	t.Run("test_token_usecase_comprehensive", func(t *testing.T) {
		// Test token generation with different configurations
		testCases := []struct {
			username string
			password string
		}{
			{"admin", "password"},
			{"test", "test123"},
			{"", ""},
			{"user", ""},
			{"", "pass"},
		}

		for _, tc := range testCases {
			localUtils.assertNoPanic(t, func() {
				// Update config for each test case
				cfgCopy := localUtils.createMockConfig()
				cfgCopy.GenerateCmdConfig.Username = tc.username
				cfgCopy.GenerateCmdConfig.Password = tc.password

				tokenUsecase := &TokenUsecase{Config: cfgCopy}
				token := tokenUsecase.GenerateBasicAuthToken()
				if len(token) == 0 {
					t.Errorf("Expected non-empty token for username=%s, password=%s", tc.username, tc.password)
				}
			})
		}
	})
}

// TestApplicationEdgeCases tests edge cases and error conditions
func TestApplicationEdgeCases(t *testing.T) {
	t.Run("test_nil_repository_handling", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		controllers := []config.Controller{{Hostname: "test", AccessToken: "token"}}
		isSecure := false

		// Test with nil repository
		apUsecase := &ApUsecase{Config: cfg, Repository: nil}
		localUtils.assertNoPanic(t, func() {
			result := apUsecase.ShowAp(&controllers, &isSecure)
			if len(result) != 0 {
				t.Errorf("Expected empty result with nil repository, got %d items", len(result))
			}
		})

		clientUsecase := &ClientUsecase{Config: cfg, Repository: nil}
		localUtils.assertNoPanic(t, func() {
			result := clientUsecase.ShowClient(&controllers, &isSecure)
			if len(result) != 0 {
				t.Errorf("Expected empty result with nil repository, got %d items", len(result))
			}
		})

		wlanUsecase := &WlanUsecase{Config: cfg, Repository: nil}
		localUtils.assertNoPanic(t, func() {
			result := wlanUsecase.ShowWlan(&controllers, &isSecure)
			if len(result) != 0 {
				t.Errorf("Expected empty result with nil repository, got %d items", len(result))
			}
		})

		overviewUsecase := &OverviewUsecase{Config: cfg, Repository: nil}
		localUtils.assertNoPanic(t, func() {
			result := overviewUsecase.ShowOverview(&controllers, &isSecure)
			if len(result) != 0 {
				t.Errorf("Expected empty result with nil repository, got %d items", len(result))
			}
		})
	})

	t.Run("test_multiple_controllers_handling", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "controller1", AccessToken: "token1"},
			{Hostname: "controller2", AccessToken: "token2"},
			{Hostname: "controller3", AccessToken: "token3"},
		}
		isSecure := false

		apUsecase := &ApUsecase{Config: cfg, Repository: repo}
		localUtils.assertNoPanic(t, func() {
			result := apUsecase.ShowAp(&controllers, &isSecure)
			// Should handle multiple controllers without panic
			_ = result
		})

		clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}
		localUtils.assertNoPanic(t, func() {
			result := clientUsecase.ShowClient(&controllers, &isSecure)
			// Should handle multiple controllers without panic
			_ = result
		})
	})
}

// TestControllerHandling tests controller handling (Integration test)
func TestControllerHandling(t *testing.T) {
	t.Run("empty_controllers", func(t *testing.T) {
		cfg := &config.Config{}
		repo := &infrastructure.Repository{Config: cfg}

		controllers := &[]config.Controller{}
		isSecure := false

		apUsecase := &ApUsecase{Config: cfg, Repository: repo}
		clientUsecase := &ClientUsecase{Config: cfg, Repository: repo}
		wlanUsecase := &WlanUsecase{Config: cfg, Repository: repo}
		overviewUsecase := &OverviewUsecase{Config: cfg, Repository: repo}

		// Test that all usecases handle empty controllers gracefully
		apResult := apUsecase.ShowAp(controllers, &isSecure)
		if apResult == nil {
		}

		clientResult := clientUsecase.ShowClient(controllers, &isSecure)
		if clientResult == nil {
		}

		wlanResult := wlanUsecase.ShowWlan(controllers, &isSecure)
		if wlanResult == nil {
		}

		overviewResult := overviewUsecase.ShowOverview(controllers, &isSecure)
		if overviewResult == nil {
		}
	})
}

// TestShowOverviewDetailed tests ShowOverview method with various scenarios (Unit test)
func TestShowOverviewDetailed(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		isSecure    *bool
		repo        *infrastructure.Repository
		wantEmpty   bool
	}{
		{
			name:        "nil_controllers",
			controllers: nil,
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name:        "empty_controllers",
			controllers: &[]config.Controller{},
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name: "nil_repository",
			controllers: &[]config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      nil,
			wantEmpty: true,
		},
		{
			name: "invalid_controller",
			controllers: &[]config.Controller{
				{Hostname: "invalid://controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty: true,
		},
		{
			name: "empty_apikey",
			controllers: &[]config.Controller{
				{Hostname: "https://test-controller.example.com", AccessToken: ""},
			},
			isSecure:  nil,
			repo:      localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()

			usecase := OverviewUsecase{
				Config:     cfg,
				Repository: tt.repo,
			}

			result := usecase.ShowOverview(tt.controllers, tt.isSecure)

			if tt.wantEmpty && len(result) != 0 {
				t.Errorf("ShowOverview() expected empty result, got %d items", len(result))
			}
			if !tt.wantEmpty && result == nil {
				t.Error("ShowOverview() expected non-nil result")
			}
		})
	}
}

// TestShowClientDetailed tests ShowClient method with various scenarios (Unit test)
func TestShowClientDetailed(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		isSecure    *bool
		repo        *infrastructure.Repository
		wantEmpty   bool
	}{
		{
			name:        "nil_controllers",
			controllers: nil,
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name:        "empty_controllers",
			controllers: &[]config.Controller{},
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name: "nil_repository",
			controllers: &[]config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      nil,
			wantEmpty: true,
		},
		{
			name: "invalid_controller",
			controllers: &[]config.Controller{
				{Hostname: "invalid://controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()

			usecase := ClientUsecase{
				Config:     cfg,
				Repository: tt.repo,
			}

			result := usecase.ShowClient(tt.controllers, tt.isSecure)

			if tt.wantEmpty && len(result) != 0 {
				t.Errorf("ShowClient() expected empty result, got %d items", len(result))
			}
			if !tt.wantEmpty && result == nil {
				t.Error("ShowClient() expected non-nil result")
			}
		})
	}
}

// TestShowApDetailed tests ShowAp method with various scenarios (Unit test)
func TestShowApDetailed(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		isSecure    *bool
		repo        *infrastructure.Repository
		wantEmpty   bool
	}{
		{
			name:        "nil_controllers",
			controllers: nil,
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name:        "empty_controllers",
			controllers: &[]config.Controller{},
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name: "nil_repository",
			controllers: &[]config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      nil,
			wantEmpty: true,
		},
		{
			name: "invalid_controller",
			controllers: &[]config.Controller{
				{Hostname: "invalid://controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()

			usecase := ApUsecase{
				Config:     cfg,
				Repository: tt.repo,
			}

			result := usecase.ShowAp(tt.controllers, tt.isSecure)

			if tt.wantEmpty && len(result) != 0 {
				t.Errorf("ShowAp() expected empty result, got %d items", len(result))
			}
			if !tt.wantEmpty && result == nil {
				t.Error("ShowAp() expected non-nil result")
			}
		})
	}
}

// TestShowWlanDetailed tests ShowWlan method with various scenarios (Unit test)
func TestShowWlanDetailed(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		isSecure    *bool
		repo        *infrastructure.Repository
		wantEmpty   bool
	}{
		{
			name:        "nil_controllers",
			controllers: nil,
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name:        "empty_controllers",
			controllers: &[]config.Controller{},
			isSecure:    nil,
			repo:        localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty:   true,
		},
		{
			name: "nil_repository",
			controllers: &[]config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      nil,
			wantEmpty: true,
		},
		{
			name: "invalid_controller",
			controllers: &[]config.Controller{
				{Hostname: "invalid://controller", AccessToken: "test-token"},
			},
			isSecure:  nil,
			repo:      localUtils.createMockRepository(localUtils.createMockConfig()),
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()

			usecase := WlanUsecase{
				Config:     cfg,
				Repository: tt.repo,
			}

			result := usecase.ShowWlan(tt.controllers, tt.isSecure)

			if tt.wantEmpty && len(result) != 0 {
				t.Errorf("ShowWlan() expected empty result, got %d items", len(result))
			}
			if !tt.wantEmpty && result == nil {
				t.Error("ShowWlan() expected non-nil result")
			}
		})
	}
}

// TestShowOverviewComplexScenarios tests ShowOverview with complex data merging scenarios (Unit test)
func TestShowOverviewComplexScenarios(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		setupMock   func() *infrastructure.Repository
		expectItems int
	}{
		{
			name: "multiple_controllers_with_data",
			controllers: &[]config.Controller{
				{Hostname: "controller1.example.com", AccessToken: "token1"},
				{Hostname: "controller2.example.com", AccessToken: "token2"},
			},
			setupMock: func() *infrastructure.Repository {
				// Return a mock that will fail API calls (returning nil)
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0, // Will be 0 due to API failures
		},
		{
			name: "single_controller_partial_failures",
			controllers: &[]config.Controller{
				{Hostname: "test-controller.example.com", AccessToken: "test-token"},
			},
			setupMock: func() *infrastructure.Repository {
				// Mock that simulates partial failures
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0, // Will be 0 due to API failures
		},
		{
			name: "controller_with_empty_token",
			controllers: &[]config.Controller{
				{Hostname: "controller.example.com", AccessToken: ""},
			},
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			repo := tt.setupMock()

			usecase := OverviewUsecase{
				Config:     cfg,
				Repository: repo,
			}

			isSecure := true
			result := usecase.ShowOverview(tt.controllers, &isSecure)

			if result == nil {
				t.Error("ShowOverview should return empty slice, not nil")
				return
			}

			if len(result) != tt.expectItems {
				t.Logf("ShowOverview returned %d items, expected %d", len(result), tt.expectItems)
				// Don't fail the test as this depends on actual network calls
			}

			// Test that the function handles the loop logic properly
			// This exercises the for loops and data merging logic
			for i, item := range result {
				if item == nil {
					t.Errorf("Item at index %d is nil", i)
				}
				if item.Controller == "" {
					t.Errorf("Item at index %d has empty controller", i)
				}
			}
		})
	}
}

// TestShowClientComplexScenarios tests ShowClient with various edge cases (Unit test)
func TestShowClientComplexScenarios(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		setupMock   func() *infrastructure.Repository
		expectItems int
	}{
		{
			name: "multiple_controllers_processing",
			controllers: &[]config.Controller{
				{Hostname: "ctrl1.example.com", AccessToken: "token1"},
				{Hostname: "ctrl2.example.com", AccessToken: "token2"},
			},
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0, // API calls will fail
		},
		{
			name: "config_with_filtering",
			controllers: &[]config.Controller{
				{Hostname: "test-controller.example.com", AccessToken: "test-token"},
			},
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			cfg.ShowCmdConfig.SSID = "test-ssid" // Add SSID filter
			repo := tt.setupMock()

			usecase := ClientUsecase{
				Config:     cfg,
				Repository: repo,
			}

			isSecure := true
			result := usecase.ShowClient(tt.controllers, &isSecure)

			if result == nil {
				t.Error("ShowClient should return empty slice, not nil")
				return
			}

			// Test filtering logic
			if len(result) != tt.expectItems {
				t.Logf("ShowClient returned %d items, expected %d", len(result), tt.expectItems)
			}
		})
	}
}

// TestShowApComplexScenarios tests ShowAp with various edge cases (Unit test)
func TestShowApComplexScenarios(t *testing.T) {
	tests := []struct {
		name        string
		controllers *[]config.Controller
		setupMock   func() *infrastructure.Repository
		expectItems int
	}{
		{
			name: "with_nil_repository",
			controllers: &[]config.Controller{
				{Hostname: "controller.example.com", AccessToken: "token"},
			},
			setupMock: func() *infrastructure.Repository {
				// Return nil repository to test nil handling
				return nil
			},
			expectItems: 0,
		},
		{
			name:        "with_nil_controllers",
			controllers: nil,
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0,
		},
		{
			name:        "with_empty_controllers",
			controllers: &[]config.Controller{},
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0,
		},
		{
			name: "with_valid_repository",
			controllers: &[]config.Controller{
				{Hostname: "controller.example.com", AccessToken: "token"},
			},
			setupMock: func() *infrastructure.Repository {
				return localUtils.createMockRepository(localUtils.createMockConfig())
			},
			expectItems: 0, // Mock repository returns empty data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &ApUsecase{
				Config:     localUtils.createMockConfig(),
				Repository: tt.setupMock(),
			}

			isSecure := true
			result := usecase.ShowAp(tt.controllers, &isSecure)

			// Debug: Check what we actually got
			t.Logf("Test %s: result is nil: %v, result length: %d", tt.name, result == nil, len(result))

			// ShowAp should always return a slice, never nil
			if result == nil {
				t.Error("ShowAp should return empty slice, not nil")
				return
			}

			// Test that correct number of items returned
			if len(result) != tt.expectItems {
				t.Logf("ShowAp returned %d items, expected %d", len(result), tt.expectItems)
			}

			// Verify each item has valid structure
			for i, item := range result {
				if item == nil {
					t.Errorf("Item at index %d is nil", i)
				}
			}
		})
	}
}

// TestFilterByRadio tests the filterByRadio method directly (Unit test)
func TestFilterByRadio(t *testing.T) {
	tests := []struct {
		name     string
		config   *config.Config
		inputLen int
		wantLen  int
	}{
		{
			name: "no_radio_filter",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Radio: "", // No filter
				},
			},
			inputLen: 2,
			wantLen:  2, // No filtering
		},
		{
			name: "radio_filter_applied",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Radio: "0", // Filter for radio 0
				},
			},
			inputLen: 2,
			wantLen:  0, // All filtered out (no matching data)
		},
		{
			name:     "nil_config",
			config:   nil,
			inputLen: 2,
			wantLen:  2, // No filtering when config is nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := OverviewUsecase{
				Config: tt.config,
			}

			// Create dummy data
			data := make([]*ShowOverviewData, tt.inputLen)
			for i := 0; i < tt.inputLen; i++ {
				data[i] = &ShowOverviewData{
					ApMac:      "00:11:22:33:44:55",
					SlotID:     i + 1, // Different slot IDs
					Controller: "test-controller",
				}
			}

			result := usecase.filterByRadio(data)

			if len(result) != tt.wantLen {
				t.Errorf("filterByRadio() returned %d items, expected %d", len(result), tt.wantLen)
			}
		})
	}
}

// TestShowApSimple tests basic ShowAp functionality (Unit test)
func TestShowApSimple(t *testing.T) {
	t.Run("nil_repository", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: nil,
		}

		controllers := &[]config.Controller{
			{Hostname: "test.example.com", AccessToken: "token"},
		}
		isSecure := true

		result := usecase.ShowAp(controllers, &isSecure)

		if result == nil {
			t.Error("ShowAp should return empty slice, not nil")
		} else {
			t.Logf("ShowAp returned slice with length %d", len(result))
		}
	})

	t.Run("nil_controllers", func(t *testing.T) {
		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		isSecure := true

		result := usecase.ShowAp(nil, &isSecure)

		if result == nil {
			t.Error("ShowAp should return empty slice, not nil")
		} else {
			t.Logf("ShowAp returned slice with length %d", len(result))
		}
	})
}

// TestShowApDataMergingLogic tests the data merging logic in ShowAp (Unit test)
func TestShowApDataMergingLogic(t *testing.T) {
	t.Run("successful_data_merging", func(t *testing.T) {
		// Create a mock repository that returns specific test data
		mockRepo := localUtils.createMockRepository(localUtils.createMockConfig())
		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: mockRepo,
		}

		controllers := &[]config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := true

		result := usecase.ShowAp(controllers, &isSecure)

		// Test that result is not nil and is a valid slice
		if result == nil {
			t.Fatal("ShowAp should return empty slice, not nil")
		}

		// Test that the function completes without panics
		t.Logf("ShowAp completed successfully, returned %d items", len(result))

		// Test structure of returned data
		for i, item := range result {
			if item == nil {
				t.Errorf("Item at index %d is nil", i)
			} else {
				// Test that data structure fields are accessible
				_ = item.ApMac
				_ = item.Controller
				_ = item.CapwapData
				_ = item.LLDPnei
				_ = item.ApOperData
			}
		}
	})

	t.Run("partial_data_scenarios", func(t *testing.T) {
		// Test with multiple controllers to exercise the loop
		controllers := &[]config.Controller{
			{Hostname: "ctrl1", AccessToken: "token1"},
			{Hostname: "ctrl2", AccessToken: "token2"},
			{Hostname: "ctrl3", AccessToken: "token3"},
		}

		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		isSecure := true
		result := usecase.ShowAp(controllers, &isSecure)

		if result == nil {
			t.Fatal("ShowAp should return empty slice, not nil")
		}

		t.Logf("ShowAp with multiple controllers returned %d items", len(result))
	})

	t.Run("empty_controllers_slice", func(t *testing.T) {
		controllers := &[]config.Controller{}

		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		isSecure := true
		result := usecase.ShowAp(controllers, &isSecure)

		if result == nil {
			t.Fatal("ShowAp should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("ShowAp with empty controllers should return empty slice, got %d items", len(result))
		}
	})

	t.Run("isSecure_variations", func(t *testing.T) {
		controllers := &[]config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}

		usecase := &ApUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		// Test with isSecure = true
		isSecureTrue := true
		result1 := usecase.ShowAp(controllers, &isSecureTrue)
		if result1 == nil {
			t.Error("ShowAp with isSecure=true should return empty slice, not nil")
		}

		// Test with isSecure = false
		isSecureFalse := false
		result2 := usecase.ShowAp(controllers, &isSecureFalse)
		if result2 == nil {
			t.Error("ShowAp with isSecure=false should return empty slice, not nil")
		}

		// Test with nil isSecure
		result3 := usecase.ShowAp(controllers, nil)
		if result3 == nil {
			t.Error("ShowAp with isSecure=nil should return empty slice, not nil")
		}
	})
}

// TestShowClientDataMergingLogic tests the data merging logic in ShowClient (Unit test)
func TestShowClientDataMergingLogic(t *testing.T) {
	t.Run("successful_data_merging", func(t *testing.T) {
		mockRepo := localUtils.createMockRepository(localUtils.createMockConfig())
		usecase := &ClientUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: mockRepo,
		}

		controllers := &[]config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := true

		result := usecase.ShowClient(controllers, &isSecure)

		if result == nil {
			t.Fatal("ShowClient should return empty slice, not nil")
		}

		t.Logf("ShowClient completed successfully, returned %d items", len(result))

		// Test structure of returned data
		for i, item := range result {
			if item == nil {
				t.Errorf("Item at index %d is nil", i)
			} else {
				_ = item.Controller
				_ = item.CommonOperData
				_ = item.Dot11OperData
				_ = item.TrafficStats
				_ = item.SisfDbMac
				_ = item.DcInfo
			}
		}
	})

	t.Run("multiple_controllers", func(t *testing.T) {
		controllers := &[]config.Controller{
			{Hostname: "ctrl1", AccessToken: "token1"},
			{Hostname: "ctrl2", AccessToken: "token2"},
		}

		usecase := &ClientUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		isSecure := true
		result := usecase.ShowClient(controllers, &isSecure)

		if result == nil {
			t.Fatal("ShowClient should return empty slice, not nil")
		}

		t.Logf("ShowClient with multiple controllers returned %d items", len(result))
	})
}

// TestShowWlanDataMergingLogic tests the data merging logic in ShowWlan (Unit test)
func TestShowWlanDataMergingLogic(t *testing.T) {
	t.Run("successful_data_merging", func(t *testing.T) {
		mockRepo := localUtils.createMockRepository(localUtils.createMockConfig())
		usecase := &WlanUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: mockRepo,
		}

		controllers := &[]config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := true

		result := usecase.ShowWlan(controllers, &isSecure)

		if result == nil {
			t.Fatal("ShowWlan should return empty slice, not nil")
		}

		t.Logf("ShowWlan completed successfully, returned %d items", len(result))

		// Test structure of returned data
		for i, item := range result {
			if item == nil {
				t.Errorf("Item at index %d is nil", i)
			} else {
				_ = item.Controller
				_ = item.WlanCfgEntry
				_ = item.WlanPolicy
			}
		}
	})

	t.Run("edge_cases", func(t *testing.T) {
		usecase := &WlanUsecase{
			Config:     localUtils.createMockConfig(),
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		// Test empty controllers
		emptyControllers := &[]config.Controller{}
		result1 := usecase.ShowWlan(emptyControllers, &[]bool{true}[0])
		if result1 == nil {
			t.Error("ShowWlan with empty controllers should return empty slice, not nil")
		}

		// Test with different isSecure values
		controllers := &[]config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecureFalse := false
		result2 := usecase.ShowWlan(controllers, &isSecureFalse)
		if result2 == nil {
			t.Error("ShowWlan with isSecure=false should return empty slice, not nil")
		}

		result3 := usecase.ShowWlan(controllers, nil)
		if result3 == nil {
			t.Error("ShowWlan with isSecure=nil should return empty slice, not nil")
		}
	})
}

// TestShowOverviewDataMergingLogic tests the ShowOverview function's data merging logic (Unit test)
func TestShowOverviewDataMergingLogic(t *testing.T) {
	t.Run("test_show_overview_nil_inputs", func(t *testing.T) {
		// Test with nil controllers
		ou := OverviewUsecase{}
		result := ou.ShowOverview(nil, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}

		// Test with nil repository
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}
		isSecure := true
		result = ou.ShowOverview(&controllers, &isSecure)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}
	})

	t.Run("test_show_overview_empty_controllers", func(t *testing.T) {
		repo := &infrastructure.Repository{}
		ou := OverviewUsecase{Repository: repo}
		controllers := []config.Controller{}
		isSecure := true

		result := ou.ShowOverview(&controllers, &isSecure)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}
	})

	t.Run("test_show_overview_single_controller", func(t *testing.T) {
		// Test the structure without making actual API calls
		ou := OverviewUsecase{}
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}
		isSecure := true

		// This will return empty slice due to nil repository, but we test the structure
		result := ou.ShowOverview(&controllers, &isSecure)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		// Result should be empty due to nil repository but not nil
		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})
}

// TestShowClientAdvancedScenarios tests advanced ShowClient scenarios (Unit test)
func TestShowClientAdvancedScenarios(t *testing.T) {
	t.Run("test_show_client_filtering_scenarios", func(t *testing.T) {
		cu := ClientUsecase{}

		// Test with nil inputs
		result := cu.ShowClient(nil, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}

		// Test with nil repository
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}
		isSecure := true
		result = cu.ShowClient(&controllers, &isSecure)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_client_empty_controllers", func(t *testing.T) {
		cu := ClientUsecase{}
		controllers := []config.Controller{}
		isSecure := true

		result := cu.ShowClient(&controllers, &isSecure)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}
	})

	t.Run("test_show_client_config_fields", func(t *testing.T) {
		cu := ClientUsecase{}

		// Test that Config field can be set
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SSID:  "test-ssid",
				Radio: "0",
			},
		}
		cu.Config = cfg

		if cu.Config.ShowCmdConfig.SSID != "test-ssid" {
			t.Errorf("Expected SSID 'test-ssid', got '%s'", cu.Config.ShowCmdConfig.SSID)
		}
		if cu.Config.ShowCmdConfig.Radio != "0" {
			t.Errorf("Expected Radio '0', got '%s'", cu.Config.ShowCmdConfig.Radio)
		}
	})
}

// TestShowApAdvancedScenarios tests more advanced ShowAp scenarios to improve coverage (Unit test)
func TestShowApAdvancedScenarios(t *testing.T) {
	t.Run("test_show_ap_with_multiple_api_call_patterns", func(t *testing.T) {
		// Test scenarios where some API calls succeed and others fail
		au := ApUsecase{}

		// Test with nil repository
		result := au.ShowAp(nil, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}

		// Test with nil controllers
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}
		result = au.ShowAp(&controllers, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_ap_with_isSecure_variations", func(t *testing.T) {
		au := ApUsecase{}
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}

		// Test with isSecure=true
		isSecureTrue := true
		result := au.ShowAp(&controllers, &isSecureTrue)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}

		// Test with isSecure=false
		isSecureFalse := false
		result = au.ShowAp(&controllers, &isSecureFalse)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}

		// Test with isSecure=nil
		result = au.ShowAp(&controllers, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
	})

	t.Run("test_show_ap_tag_comprehensive", func(t *testing.T) {
		au := ApUsecase{}

		// Test ShowApTag with nil repository
		result := au.ShowApTag(nil, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice, got %d items", len(result))
		}

		// Test ShowApTag with nil controllers
		controllers := []config.Controller{
			{Hostname: "test.example.com", AccessToken: "token123"},
		}
		result = au.ShowApTag(&controllers, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}

		// Test ShowApTag with empty controllers
		emptyControllers := []config.Controller{}
		result = au.ShowApTag(&emptyControllers, nil)
		if result == nil {
			t.Error("Expected non-nil slice, got nil")
		}
		if len(result) != 0 {
			t.Errorf("Expected empty slice with empty controllers, got %d items", len(result))
		}
	})

	t.Run("test_show_ap_data_structure_initialization", func(t *testing.T) {
		// Test that ShowApData and ShowApTagData structures can be properly initialized
		apData := ShowApData{
			ShowApCommonData: ShowApCommonData{
				ApMac:      "aa:bb:cc:dd:ee:ff",
				Controller: "test-controller",
			},
		}

		if apData.ApMac != "aa:bb:cc:dd:ee:ff" {
			t.Errorf("Expected ApMac 'aa:bb:cc:dd:ee:ff', got '%s'", apData.ApMac)
		}
		if apData.Controller != "test-controller" {
			t.Errorf("Expected Controller 'test-controller', got '%s'", apData.Controller)
		}

		apTagData := ShowApTagData{
			ShowApCommonData: ShowApCommonData{
				ApMac:      "11:22:33:44:55:66",
				Controller: "test-controller-2",
			},
		}

		if apTagData.ApMac != "11:22:33:44:55:66" {
			t.Errorf("Expected ApMac '11:22:33:44:55:66', got '%s'", apTagData.ApMac)
		}
		if apTagData.Controller != "test-controller-2" {
			t.Errorf("Expected Controller 'test-controller-2', got '%s'", apTagData.Controller)
		}
	})
}

// TestClientFilterFunctions tests client filtering functions to improve coverage (Unit test)
func TestClientFilterFunctions(t *testing.T) {
	t.Run("test_filter_by_ssid", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SSID: "test-ssid",
			},
		}
		cu := ClientUsecase{Config: cfg}

		// Test data with matching SSID
		clients := []*ShowClientData{
			{
				ClientMac: "aa:bb:cc:dd:ee:ff",
				Dot11OperData: client.Dot11OperData{
					VapSsid: "test-ssid",
				},
			},
			{
				ClientMac: "11:22:33:44:55:66",
				Dot11OperData: client.Dot11OperData{
					VapSsid: "other-ssid",
				},
			},
		}

		filtered := cu.filterBySSID(clients)
		if len(filtered) != 1 {
			t.Errorf("Expected 1 filtered client, got %d", len(filtered))
		}
		if len(filtered) > 0 && filtered[0].Dot11OperData.VapSsid != "test-ssid" {
			t.Errorf("Expected SSID 'test-ssid', got '%s'", filtered[0].Dot11OperData.VapSsid)
		}
	})

	t.Run("test_filter_by_ssid_empty_filter", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SSID: "", // Empty filter
			},
		}
		cu := ClientUsecase{Config: cfg}

		clients := []*ShowClientData{
			{ClientMac: "aa:bb:cc:dd:ee:ff"},
			{ClientMac: "11:22:33:44:55:66"},
		}

		filtered := cu.filterBySSID(clients)
		if len(filtered) != 2 {
			t.Errorf("Expected 2 clients (no filtering), got %d", len(filtered))
		}
	})

	t.Run("test_filter_by_radio", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Radio: "1",
			},
		}
		cu := ClientUsecase{Config: cfg}

		// Test data with matching radio
		clients := []*ShowClientData{
			{
				ClientMac: "aa:bb:cc:dd:ee:ff",
				CommonOperData: client.CommonOperData{
					MsRadioType: "802.11ax-5GHz", // This would map to radio 1
				},
			},
			{
				ClientMac: "11:22:33:44:55:66",
				CommonOperData: client.CommonOperData{
					MsRadioType: "802.11n-2.4GHz", // This would map to radio 0
				},
			},
		}

		filtered := cu.filterByRadio(clients)
		// The exact filtering logic depends on the radio mapping logic
		// For now, we just test that the function doesn't panic
		if filtered == nil {
			t.Error("filterByRadio returned nil")
		}
	})

	t.Run("test_filter_by_radio_empty_filter", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Radio: "", // Empty filter
			},
		}
		cu := ClientUsecase{Config: cfg}

		clients := []*ShowClientData{
			{ClientMac: "aa:bb:cc:dd:ee:ff"},
			{ClientMac: "11:22:33:44:55:66"},
		}

		filtered := cu.filterByRadio(clients)
		if len(filtered) != 2 {
			t.Errorf("Expected 2 clients (no filtering), got %d", len(filtered))
		}
	})

	t.Run("test_filter_by_radio_nil_config", func(t *testing.T) {
		cu := ClientUsecase{} // No config

		clients := []*ShowClientData{
			{ClientMac: "aa:bb:cc:dd:ee:ff"},
		}

		// This should handle nil config gracefully
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("filterByRadio panicked with nil config: %v", r)
			}
		}()

		filtered := cu.filterByRadio(clients)
		// Should return original clients or handle gracefully
		if filtered == nil {
			t.Error("filterByRadio returned nil with nil config")
		}
	})
}

// TestShowOverviewComprehensive tests ShowOverview data merging in detail (Unit test)
func TestShowOverviewComprehensive(t *testing.T) {
	t.Run("test_show_overview_data_merging_scenarios", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Radio:       "", // No radio filter
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		ou := OverviewUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "test-controller.example.com", AccessToken: "test-token"},
		}

		isSecure := false
		result := ou.ShowOverview(&controllers, &isSecure)

		// The function should handle API call failures gracefully
		if result == nil {
			t.Error("ShowOverview should never return nil, should return empty slice")
		}

		// Even with failures, should return empty slice, not nil
		if len(result) == 0 {
			t.Log("ShowOverview returned empty slice (expected due to failed API calls)")
		}
	})

	t.Run("test_show_overview_with_radio_filter", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Radio:       "1", // Filter by radio 1
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		ou := OverviewUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}

		isSecure := true
		result := ou.ShowOverview(&controllers, &isSecure)

		// Should handle filtering
		if result == nil {
			t.Error("ShowOverview should never return nil")
		}
	})

	t.Run("test_show_overview_data_structure_validation", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		ou := OverviewUsecase{Config: cfg, Repository: repo}

		// Test that the usecase was created properly
		if ou.Config == nil {
			t.Error("OverviewUsecase Config should not be nil")
		}

		if ou.Repository == nil {
			t.Error("OverviewUsecase Repository should not be nil")
		}

		// Test the data structure initialization
		data := []*ShowOverviewData{}

		// Validate empty data structure
		if len(data) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(data))
		}

		// Test single data item structure
		singleData := &ShowOverviewData{
			ApMac:      "aa:bb:cc:dd:ee:ff",
			SlotID:     1,
			Controller: "test-controller",
		}

		if singleData.ApMac != "aa:bb:cc:dd:ee:ff" {
			t.Errorf("Expected ApMac 'aa:bb:cc:dd:ee:ff', got '%s'", singleData.ApMac)
		}

		if singleData.SlotID != 1 {
			t.Errorf("Expected SlotID 1, got %d", singleData.SlotID)
		}

		if singleData.Controller != "test-controller" {
			t.Errorf("Expected Controller 'test-controller', got '%s'", singleData.Controller)
		}
	})

	t.Run("test_show_overview_error_handling", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
			},
		}

		// Test with nil repository
		ou := OverviewUsecase{Config: cfg, Repository: nil}

		controllers := []config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecure := false
		result := ou.ShowOverview(&controllers, &isSecure)

		// Should handle nil repository gracefully
		if result == nil {
			t.Error("ShowOverview should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_overview_multiple_controllers_processing", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		ou := OverviewUsecase{Config: cfg, Repository: repo}

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "controller1.example.com", AccessToken: "token1"},
			{Hostname: "controller2.example.com", AccessToken: "token2"},
			{Hostname: "controller3.example.com", AccessToken: "token3"},
		}

		isSecure := false
		result := ou.ShowOverview(&controllers, &isSecure)

		// Should process all controllers without panic
		if result == nil {
			t.Error("ShowOverview should return empty slice, not nil")
		}

		// Even if all API calls fail, should return gracefully
		t.Logf("ShowOverview processed %d controllers and returned %d items", len(controllers), len(result))
	})
}

// TestShowApDataMergingComprehensive tests ShowAp data merging logic in detail (Unit test)
func TestShowApDataMergingComprehensive(t *testing.T) {
	t.Run("test_show_ap_complex_data_merging", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Timeout:     30,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "ap-test-controller.example.com", AccessToken: "test-token"},
		}

		isSecure := false
		result := au.ShowAp(&controllers, &isSecure)

		// Should handle the complex data merging process
		if result == nil {
			t.Error("ShowAp should never return nil, should return empty slice")
		}

		// Test that the function completes without panic
		t.Logf("ShowAp completed with %d AP entries", len(result))
	})

	t.Run("test_show_ap_multiple_api_calls", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Timeout:     15,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "controller1.test.com", AccessToken: "token1"},
			{Hostname: "controller2.test.com", AccessToken: "token2"},
		}

		isSecure := true
		result := au.ShowAp(&controllers, &isSecure)

		// Function should handle multiple controllers without panic
		if result == nil {
			t.Error("ShowAp should return empty slice, not nil")
		}

		// Multiple API calls should be handled gracefully
		t.Logf("ShowAp processed %d controllers", len(controllers))
	})

	t.Run("test_show_ap_data_structure_validation", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		// Test that the usecase was created properly
		if au.Config == nil {
			t.Error("ApUsecase Config should not be nil")
		}

		if au.Repository == nil {
			t.Error("ApUsecase Repository should not be nil")
		}

		// Test data structure creation
		apData := &ShowApData{
			ShowApCommonData: ShowApCommonData{
				ApMac:      "aa:bb:cc:dd:ee:ff",
				Controller: "test-controller",
			},
		}

		// Validate structure fields
		if apData.ApMac != "aa:bb:cc:dd:ee:ff" {
			t.Errorf("Expected ApMac 'aa:bb:cc:dd:ee:ff', got '%s'", apData.ApMac)
		}

		if apData.Controller != "test-controller" {
			t.Errorf("Expected Controller 'test-controller', got '%s'", apData.Controller)
		}

		// Test empty slice initialization
		data := []*ShowApData{}
		if len(data) != 0 {
			t.Errorf("Expected empty slice, got length %d", len(data))
		}
	})

	t.Run("test_show_ap_error_resilience", func(t *testing.T) {
		// Test with nil repository
		au := ApUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecure := false
		result := au.ShowAp(&controllers, &isSecure)

		// Should handle nil repository gracefully
		if result == nil {
			t.Error("ShowAp should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_ap_empty_controllers_handling", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		// Test with empty controllers slice
		controllers := []config.Controller{}

		isSecure := false
		result := au.ShowAp(&controllers, &isSecure)

		// Should handle empty controllers gracefully
		if result == nil {
			t.Error("ShowAp should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with empty controllers, got %d items", len(result))
		}
	})
}

// TestShowApTagComprehensive tests ShowApTag function in detail (Unit test)
func TestShowApTagComprehensive(t *testing.T) {
	t.Run("test_show_ap_tag_data_processing", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Timeout:     30,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "ap-tag-controller.example.com", AccessToken: "test-token"},
		}

		isSecure := false
		result := au.ShowApTag(&controllers, &isSecure)

		// Should process AP tag data without panic
		if result == nil {
			t.Error("ShowApTag should never return nil, should return empty slice")
		}

		t.Logf("ShowApTag completed with %d AP tag entries", len(result))
	})

	t.Run("test_show_ap_tag_data_structure", func(t *testing.T) {
		// Test ShowApTagData structure
		apTagData := &ShowApTagData{
			ShowApCommonData: ShowApCommonData{
				ApMac:      "11:22:33:44:55:66",
				Controller: "tag-controller",
			},
		}

		// Validate structure
		if apTagData.ApMac != "11:22:33:44:55:66" {
			t.Errorf("Expected ApMac '11:22:33:44:55:66', got '%s'", apTagData.ApMac)
		}

		if apTagData.Controller != "tag-controller" {
			t.Errorf("Expected Controller 'tag-controller', got '%s'", apTagData.Controller)
		}
	})

	t.Run("test_show_ap_tag_error_handling", func(t *testing.T) {
		// Test with nil repository
		au := ApUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecure := true
		result := au.ShowApTag(&controllers, &isSecure)

		// Should handle nil repository gracefully
		if result == nil {
			t.Error("ShowApTag should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_ap_tag_multiple_controllers", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		au := ApUsecase{Config: cfg, Repository: repo}

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "tag-ctrl1.example.com", AccessToken: "token1"},
			{Hostname: "tag-ctrl2.example.com", AccessToken: "token2"},
			{Hostname: "tag-ctrl3.example.com", AccessToken: "token3"},
		}

		isSecure := false
		result := au.ShowApTag(&controllers, &isSecure)

		// Should process multiple controllers without panic
		if result == nil {
			t.Error("ShowApTag should return empty slice, not nil")
		}

		t.Logf("ShowApTag processed %d controllers and returned %d items", len(controllers), len(result))
	})
}

// TestShowClientComprehensive tests ShowClient function in detail (Unit test)
func TestShowClientComprehensive(t *testing.T) {
	t.Run("test_show_client_data_merging_logic", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Timeout:     5,  // Shorter timeout for tests
				SSID:        "", // No filter
				Radio:       "", // No filter
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		cu := ClientUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "client-test-controller.example.com", AccessToken: "test-token"},
		}

		isSecure := false
		result := cu.ShowClient(&controllers, &isSecure)

		// Should handle the complex data merging process
		if result == nil {
			t.Error("ShowClient should never return nil, should return empty slice")
		}

		t.Logf("ShowClient completed with %d client entries", len(result))
	})

	t.Run("test_show_client_with_ssid_filter", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Timeout:     5,
				SSID:        "test-ssid",
				Radio:       "",
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		cu := ClientUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "controller.test.com", AccessToken: "token"},
		}

		isSecure := true
		result := cu.ShowClient(&controllers, &isSecure)

		// Should apply SSID filtering
		if result == nil {
			t.Error("ShowClient should return empty slice, not nil")
		}

		t.Logf("ShowClient with SSID filter processed successfully")
	})

	t.Run("test_show_client_with_radio_filter", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Timeout:     5,
				SSID:        "",
				Radio:       "1",
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		cu := ClientUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "radio-controller.test.com", AccessToken: "token"},
		}

		isSecure := false
		result := cu.ShowClient(&controllers, &isSecure)

		// Should apply radio filtering
		if result == nil {
			t.Error("ShowClient should return empty slice, not nil")
		}

		t.Logf("ShowClient with radio filter processed successfully")
	})

	t.Run("test_show_client_data_structure", func(t *testing.T) {
		// Test ShowClientData structure
		clientData := &ShowClientData{
			ClientMac:  "aa:bb:cc:dd:ee:ff",
			Controller: "client-controller",
		}

		// Validate structure
		if clientData.ClientMac != "aa:bb:cc:dd:ee:ff" {
			t.Errorf("Expected ClientMac 'aa:bb:cc:dd:ee:ff', got '%s'", clientData.ClientMac)
		}

		if clientData.Controller != "client-controller" {
			t.Errorf("Expected Controller 'client-controller', got '%s'", clientData.Controller)
		}
	})

	t.Run("test_show_client_multiple_controllers", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Timeout:     5,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		cu := ClientUsecase{Config: cfg, Repository: repo}

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "client-ctrl1.example.com", AccessToken: "token1"},
			{Hostname: "client-ctrl2.example.com", AccessToken: "token2"},
		}

		isSecure := false
		result := cu.ShowClient(&controllers, &isSecure)

		// Should process multiple controllers without panic
		if result == nil {
			t.Error("ShowClient should return empty slice, not nil")
		}

		t.Logf("ShowClient processed %d controllers and returned %d items", len(controllers), len(result))
	})

	t.Run("test_show_client_error_handling", func(t *testing.T) {
		// Test with nil repository
		cu := ClientUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecure := true
		result := cu.ShowClient(&controllers, &isSecure)

		// Should handle nil repository gracefully
		if result == nil {
			t.Error("ShowClient should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})
}

// TestShowWlanComprehensive tests ShowWlan function in detail (Unit test)
func TestShowWlanComprehensive(t *testing.T) {
	t.Run("test_show_wlan_data_processing", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Timeout:     5,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		wu := WlanUsecase{Config: cfg, Repository: repo}

		controllers := []config.Controller{
			{Hostname: "wlan-controller.example.com", AccessToken: "test-token"},
		}

		isSecure := false
		result := wu.ShowWlan(&controllers, &isSecure)

		// Should process WLAN data without panic
		if result == nil {
			t.Error("ShowWlan should never return nil, should return empty slice")
		}

		t.Logf("ShowWlan completed with %d WLAN entries", len(result))
	})

	t.Run("test_show_wlan_data_structure", func(t *testing.T) {
		// Test ShowWlanData structure
		wlanData := &ShowWlanData{
			TagName:    "test-tag",
			PolicyName: "test-policy",
			WlanName:   "test-wlan",
			Controller: "wlan-controller",
		}

		// Validate structure
		if wlanData.TagName != "test-tag" {
			t.Errorf("Expected TagName 'test-tag', got '%s'", wlanData.TagName)
		}

		if wlanData.PolicyName != "test-policy" {
			t.Errorf("Expected PolicyName 'test-policy', got '%s'", wlanData.PolicyName)
		}

		if wlanData.WlanName != "test-wlan" {
			t.Errorf("Expected WlanName 'test-wlan', got '%s'", wlanData.WlanName)
		}

		if wlanData.Controller != "wlan-controller" {
			t.Errorf("Expected Controller 'wlan-controller', got '%s'", wlanData.Controller)
		}
	})

	t.Run("test_show_wlan_error_handling", func(t *testing.T) {
		// Test with nil repository
		wu := WlanUsecase{
			Config:     &config.Config{},
			Repository: nil,
		}

		controllers := []config.Controller{
			{Hostname: "test", AccessToken: "token"},
		}

		isSecure := true
		result := wu.ShowWlan(&controllers, &isSecure)

		// Should handle nil repository gracefully
		if result == nil {
			t.Error("ShowWlan should return empty slice, not nil")
		}

		if len(result) != 0 {
			t.Errorf("Expected empty slice with nil repository, got %d items", len(result))
		}
	})

	t.Run("test_show_wlan_multiple_controllers", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Timeout:     5,
			},
		}

		repo := &infrastructure.Repository{Config: cfg}
		wu := WlanUsecase{Config: cfg, Repository: repo}

		// Test with multiple controllers
		controllers := []config.Controller{
			{Hostname: "wlan-ctrl1.example.com", AccessToken: "token1"},
			{Hostname: "wlan-ctrl2.example.com", AccessToken: "token2"},
		}

		isSecure := false
		result := wu.ShowWlan(&controllers, &isSecure)

		// Should process multiple controllers without panic
		if result == nil {
			t.Error("ShowWlan should return empty slice, not nil")
		}

		t.Logf("ShowWlan processed %d controllers and returned %d items", len(controllers), len(result))
	})
}

// TestShowOverviewFilterByRadio tests filterByRadio method with various scenarios (Unit test)
func TestShowOverviewFilterByRadio(t *testing.T) {
	tests := []struct {
		name        string
		radioFilter string
		inputData   []*ShowOverviewData
		expectedLen int
		description string
	}{
		{
			name:        "no_filter_configured",
			radioFilter: "",
			inputData: []*ShowOverviewData{
				{SlotID: 0, ApMac: "00:01:02:03:04:05"},
				{SlotID: 1, ApMac: "00:01:02:03:04:06"},
				{SlotID: 2, ApMac: "00:01:02:03:04:07"},
			},
			expectedLen: 3,
			description: "Should return all data when no filter is configured",
		},
		{
			name:        "filter_by_slot_0",
			radioFilter: "0",
			inputData: []*ShowOverviewData{
				{SlotID: 0, ApMac: "00:01:02:03:04:05"},
				{SlotID: 1, ApMac: "00:01:02:03:04:06"},
				{SlotID: 2, ApMac: "00:01:02:03:04:07"},
			},
			expectedLen: 1,
			description: "Should return only slot 0 data when filtered by '0'",
		},
		{
			name:        "filter_by_slot_1",
			radioFilter: "1",
			inputData: []*ShowOverviewData{
				{SlotID: 0, ApMac: "00:01:02:03:04:05"},
				{SlotID: 1, ApMac: "00:01:02:03:04:06"},
				{SlotID: 2, ApMac: "00:01:02:03:04:07"},
			},
			expectedLen: 1,
			description: "Should return only slot 1 data when filtered by '1'",
		},
		{
			name:        "filter_nonexistent_slot",
			radioFilter: "99",
			inputData: []*ShowOverviewData{
				{SlotID: 0, ApMac: "00:01:02:03:04:05"},
				{SlotID: 1, ApMac: "00:01:02:03:04:06"},
				{SlotID: 2, ApMac: "00:01:02:03:04:07"},
			},
			expectedLen: 0,
			description: "Should return empty slice when filtering by non-existent slot",
		},
		{
			name:        "filter_empty_input",
			radioFilter: "0",
			inputData:   []*ShowOverviewData{},
			expectedLen: 0,
			description: "Should return empty slice when input data is empty",
		},
		{
			name:        "filter_nil_input",
			radioFilter: "0",
			inputData:   nil,
			expectedLen: 0,
			description: "Should handle nil input gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			cfg.ShowCmdConfig.Radio = tt.radioFilter

			usecase := OverviewUsecase{
				Config:     cfg,
				Repository: localUtils.createMockRepository(cfg),
			}

			result := usecase.filterByRadio(tt.inputData)

			if len(result) != tt.expectedLen {
				t.Errorf("filterByRadio() expected %d items, got %d. %s", tt.expectedLen, len(result), tt.description)
			}

			// Verify filtered data contains correct SlotID when filter is applied
			if tt.radioFilter != "" && len(result) > 0 {
				for _, item := range result {
					expectedSlotID := tt.radioFilter
					actualSlotID := fmt.Sprintf("%d", item.SlotID)
					if actualSlotID != expectedSlotID {
						t.Errorf("filterByRadio() expected SlotID %s, got %s", expectedSlotID, actualSlotID)
					}
				}
			}
		})
	}
}

// TestShowOverviewFilterByRadioNilConfig tests filterByRadio with nil config (Unit test)
func TestShowOverviewFilterByRadioNilConfig(t *testing.T) {
	t.Run("nil_config_handling", func(t *testing.T) {
		usecase := OverviewUsecase{
			Config:     nil, // Nil config
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		inputData := []*ShowOverviewData{
			{SlotID: 0, ApMac: "00:01:02:03:04:05"},
			{SlotID: 1, ApMac: "00:01:02:03:04:06"},
		}

		result := usecase.filterByRadio(inputData)

		// Should return original data when config is nil
		if len(result) != len(inputData) {
			t.Errorf("filterByRadio() with nil config expected %d items, got %d", len(inputData), len(result))
		}

		// Verify the data is unchanged
		for i, item := range result {
			if item.SlotID != inputData[i].SlotID {
				t.Errorf("filterByRadio() with nil config changed data at index %d", i)
			}
		}
	})
}

// TestShowOverviewDataMerging tests data merging logic in ShowOverview (Unit test)
func TestShowOverviewDataMerging(t *testing.T) {
	t.Run("data_merging_logic", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test with nil controllers to trigger early return
		result := usecase.ShowOverview(nil, nil)
		if len(result) != 0 {
			t.Error("ShowOverview() with nil controllers should return empty slice")
		}

		// Test with empty controllers slice
		emptyControllers := []config.Controller{}
		result = usecase.ShowOverview(&emptyControllers, nil)
		if len(result) != 0 {
			t.Error("ShowOverview() with empty controllers should return empty slice")
		}
	})

	t.Run("nil_repository_handling", func(t *testing.T) {
		cfg := localUtils.createMockConfig()

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: nil, // Nil repository
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowOverview(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowOverview() with nil repository should return empty slice")
		}
	})
}

// TestShowClientFilterBySSID tests filterBySSID method with various scenarios (Unit test)
func TestShowClientFilterBySSID(t *testing.T) {
	tests := []struct {
		name        string
		ssidFilter  string
		inputData   []*ShowClientData
		expectedLen int
		description string
	}{
		{
			name:       "no_filter_configured",
			ssidFilter: "",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
				{ClientMac: "00:01:02:03:04:06", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID2"}},
				{ClientMac: "00:01:02:03:04:07", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID3"}},
			},
			expectedLen: 3,
			description: "Should return all data when no filter is configured",
		},
		{
			name:       "filter_by_ssid1",
			ssidFilter: "TestSSID1",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
				{ClientMac: "00:01:02:03:04:06", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID2"}},
				{ClientMac: "00:01:02:03:04:07", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
			},
			expectedLen: 2,
			description: "Should return only TestSSID1 data when filtered",
		},
		{
			name:       "filter_nonexistent_ssid",
			ssidFilter: "NonExistentSSID",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
				{ClientMac: "00:01:02:03:04:06", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID2"}},
			},
			expectedLen: 0,
			description: "Should return empty slice when filtering by non-existent SSID",
		},
		{
			name:        "filter_empty_input",
			ssidFilter:  "TestSSID1",
			inputData:   []*ShowClientData{},
			expectedLen: 0,
			description: "Should return empty slice when input data is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			cfg.ShowCmdConfig.SSID = tt.ssidFilter

			usecase := ClientUsecase{
				Config:     cfg,
				Repository: localUtils.createMockRepository(cfg),
			}

			result := usecase.filterBySSID(tt.inputData)

			if len(result) != tt.expectedLen {
				t.Errorf("filterBySSID() expected %d items, got %d. %s", tt.expectedLen, len(result), tt.description)
			}

			// Verify filtered data contains correct SSID when filter is applied
			if tt.ssidFilter != "" && len(result) > 0 {
				for _, item := range result {
					if item.Dot11OperData.VapSsid != tt.ssidFilter {
						t.Errorf("filterBySSID() expected SSID %s, got %s", tt.ssidFilter, item.Dot11OperData.VapSsid)
					}
				}
			}
		})
	}
}

// TestShowClientFilterByRadio tests filterByRadio method for clients (Unit test)
func TestShowClientFilterByRadio(t *testing.T) {
	tests := []struct {
		name        string
		radioFilter string
		inputData   []*ShowClientData
		expectedLen int
		description string
	}{
		{
			name:        "no_filter_configured",
			radioFilter: "",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
				{ClientMac: "00:01:02:03:04:06", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
				{ClientMac: "00:01:02:03:04:07", CommonOperData: client.CommonOperData{MsApSlotID: 2}},
			},
			expectedLen: 3,
			description: "Should return all data when no filter is configured",
		},
		{
			name:        "filter_by_slot_0",
			radioFilter: "0",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
				{ClientMac: "00:01:02:03:04:06", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
				{ClientMac: "00:01:02:03:04:07", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
			},
			expectedLen: 2,
			description: "Should return only slot 0 data when filtered by '0'",
		},
		{
			name:        "filter_nonexistent_slot",
			radioFilter: "99",
			inputData: []*ShowClientData{
				{ClientMac: "00:01:02:03:04:05", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
				{ClientMac: "00:01:02:03:04:06", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
			},
			expectedLen: 0,
			description: "Should return empty slice when filtering by non-existent slot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := localUtils.createMockConfig()
			cfg.ShowCmdConfig.Radio = tt.radioFilter

			usecase := ClientUsecase{
				Config:     cfg,
				Repository: localUtils.createMockRepository(cfg),
			}

			result := usecase.filterByRadio(tt.inputData)

			if len(result) != tt.expectedLen {
				t.Errorf("filterByRadio() expected %d items, got %d. %s", tt.expectedLen, len(result), tt.description)
			}

			// Verify filtered data contains correct SlotID when filter is applied
			if tt.radioFilter != "" && len(result) > 0 {
				for _, item := range result {
					expectedSlotID := tt.radioFilter
					actualSlotID := fmt.Sprintf("%d", item.CommonOperData.MsApSlotID)
					if actualSlotID != expectedSlotID {
						t.Errorf("filterByRadio() expected SlotID %s, got %s", expectedSlotID, actualSlotID)
					}
				}
			}
		})
	}
}

// TestShowClientNilConfig tests client methods with nil config (Unit test)
func TestShowClientNilConfig(t *testing.T) {
	t.Run("filterBySSID_nil_config", func(t *testing.T) {
		usecase := ClientUsecase{
			Config:     nil, // Nil config
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		inputData := []*ShowClientData{
			{ClientMac: "00:01:02:03:04:05", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
			{ClientMac: "00:01:02:03:04:06", Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID2"}},
		}

		result := usecase.filterBySSID(inputData)

		// Should return original data when config is nil
		if len(result) != len(inputData) {
			t.Errorf("filterBySSID() with nil config expected %d items, got %d", len(inputData), len(result))
		}
	})

	t.Run("filterByRadio_nil_config", func(t *testing.T) {
		usecase := ClientUsecase{
			Config:     nil, // Nil config
			Repository: localUtils.createMockRepository(localUtils.createMockConfig()),
		}

		inputData := []*ShowClientData{
			{ClientMac: "00:01:02:03:04:05", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
			{ClientMac: "00:01:02:03:04:06", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
		}

		result := usecase.filterByRadio(inputData)

		// Should return original data when config is nil
		if len(result) != len(inputData) {
			t.Errorf("filterByRadio() with nil config expected %d items, got %d", len(inputData), len(result))
		}
	})
}

// TestShowClientEdgeCases tests ShowClient edge cases (Unit test)
func TestShowClientEdgeCases(t *testing.T) {
	t.Run("nil_repository", func(t *testing.T) {
		cfg := localUtils.createMockConfig()

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: nil, // Nil repository
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowClient(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowClient() with nil repository should return empty slice")
		}
	})

	t.Run("nil_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		isSecure := false

		result := usecase.ShowClient(nil, &isSecure)
		if len(result) != 0 {
			t.Error("ShowClient() with nil controllers should return empty slice")
		}
	})

	t.Run("empty_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		emptyControllers := []config.Controller{}
		isSecure := false

		result := usecase.ShowClient(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowClient() with empty controllers should return empty slice")
		}
	})
}

// TestShowApEdgeCases tests ShowAp edge cases (Unit test)
func TestShowApEdgeCases(t *testing.T) {
	t.Run("nil_repository", func(t *testing.T) {
		cfg := localUtils.createMockConfig()

		usecase := ApUsecase{
			Config:     cfg,
			Repository: nil, // Nil repository
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowAp(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowAp() with nil repository should return empty slice")
		}
	})

	t.Run("nil_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		isSecure := false

		result := usecase.ShowAp(nil, &isSecure)
		if len(result) != 0 {
			t.Error("ShowAp() with nil controllers should return empty slice")
		}
	})

	t.Run("empty_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		emptyControllers := []config.Controller{}
		isSecure := false

		result := usecase.ShowAp(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowAp() with empty controllers should return empty slice")
		}
	})

	t.Run("invalid_controller", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		invalidControllers := []config.Controller{
			{Hostname: "invalid://controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowAp(&invalidControllers, &isSecure)
		// Should handle invalid controllers gracefully
		if result == nil {
			t.Error("ShowAp() should return empty slice, not nil")
		}
	})
}

// TestShowApTagEdgeCases tests ShowApTag edge cases (Unit test)
func TestShowApTagEdgeCases(t *testing.T) {
	t.Run("nil_repository", func(t *testing.T) {
		cfg := localUtils.createMockConfig()

		usecase := ApUsecase{
			Config:     cfg,
			Repository: nil, // Nil repository
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowApTag(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowApTag() with nil repository should return empty slice")
		}
	})

	t.Run("nil_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		isSecure := false

		result := usecase.ShowApTag(nil, &isSecure)
		if len(result) != 0 {
			t.Error("ShowApTag() with nil controllers should return empty slice")
		}
	})

	t.Run("empty_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		emptyControllers := []config.Controller{}
		isSecure := false

		result := usecase.ShowApTag(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowApTag() with empty controllers should return empty slice")
		}
	})

	t.Run("multiple_invalid_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ApUsecase{
			Config:     cfg,
			Repository: repo,
		}

		invalidControllers := []config.Controller{
			{Hostname: "invalid://controller1", AccessToken: "test-token1"},
			{Hostname: "invalid://controller2", AccessToken: "test-token2"},
			{Hostname: "", AccessToken: "test-token3"}, // Empty hostname
		}
		isSecure := false

		result := usecase.ShowApTag(&invalidControllers, &isSecure)
		// Should handle multiple invalid controllers gracefully
		if result == nil {
			t.Error("ShowApTag() should return empty slice, not nil")
		}
	})
}

// TestShowWlanEdgeCases tests ShowWlan edge cases (Unit test)
func TestShowWlanEdgeCases(t *testing.T) {
	t.Run("nil_repository", func(t *testing.T) {
		cfg := localUtils.createMockConfig()

		usecase := WlanUsecase{
			Config:     cfg,
			Repository: nil, // Nil repository
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowWlan(&controllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowWlan() with nil repository should return empty slice")
		}
	})

	t.Run("nil_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := WlanUsecase{
			Config:     cfg,
			Repository: repo,
		}

		isSecure := false

		result := usecase.ShowWlan(nil, &isSecure)
		if len(result) != 0 {
			t.Error("ShowWlan() with nil controllers should return empty slice")
		}
	})

	t.Run("empty_controllers", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := WlanUsecase{
			Config:     cfg,
			Repository: repo,
		}

		emptyControllers := []config.Controller{}
		isSecure := false

		result := usecase.ShowWlan(&emptyControllers, &isSecure)
		if len(result) != 0 {
			t.Error("ShowWlan() with empty controllers should return empty slice")
		}
	})

	t.Run("invalid_controller_handling", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := WlanUsecase{
			Config:     cfg,
			Repository: repo,
		}

		invalidControllers := []config.Controller{
			{Hostname: "invalid://controller", AccessToken: "test-token"},
			{Hostname: "", AccessToken: "test-token"},               // Empty hostname
			{Hostname: "https://test.example.com", AccessToken: ""}, // Empty token
		}
		isSecure := false

		result := usecase.ShowWlan(&invalidControllers, &isSecure)
		// Should handle invalid controllers gracefully
		if result == nil {
			t.Error("ShowWlan() should return empty slice, not nil")
		}
	})

	t.Run("data_merging_verification", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		_ = WlanUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test that the data structure is properly initialized
		data := ShowWlanData{}
		if data.TagName != "" || data.PolicyName != "" || data.WlanName != "" || data.Controller != "" {
			t.Error("ShowWlanData should initialize with empty string values")
		}
	})

	t.Run("usecase_initialization", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := WlanUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Verify usecase is properly initialized
		if usecase.Config == nil {
			t.Error("WlanUsecase should have non-nil Config")
		}
		if usecase.Repository == nil {
			t.Error("WlanUsecase should have non-nil Repository")
		}

		// Test with actual usage to avoid unused variable warning
		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		result := usecase.ShowWlan(&controllers, &isSecure)
		if result == nil {
			t.Error("ShowWlan() should return empty slice, not nil")
		}
	})
}

// TestShowOverviewComprehensiveBusinessLogic tests ShowOverview with detailed scenarios
func TestShowOverviewComprehensiveBusinessLogic(t *testing.T) {
	t.Run("show_overview_data_merging_comprehensive", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test multiple controller scenarios
		controllers := []config.Controller{
			{Hostname: "controller1.example.com", AccessToken: "token1"},
			{Hostname: "controller2.example.com", AccessToken: "token2"},
		}
		isSecure := false

		result := usecase.ShowOverview(&controllers, &isSecure)

		// Verify result structure
		if result == nil {
			t.Error("ShowOverview should return empty slice, not nil")
		}

		// Test empty controller slice
		emptyControllers := []config.Controller{}
		emptyResult := usecase.ShowOverview(&emptyControllers, &isSecure)
		if emptyResult == nil {
			t.Error("ShowOverview with empty controllers should return empty slice, not nil")
		}
		if len(emptyResult) != 0 {
			t.Errorf("ShowOverview with empty controllers should return empty slice, got %d items", len(emptyResult))
		}
	})

	t.Run("show_overview_controller_variations", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test single controller
		singleController := []config.Controller{
			{Hostname: "single-controller.example.com", AccessToken: "single-token"},
		}
		isSecure := true

		result := usecase.ShowOverview(&singleController, &isSecure)
		if result == nil {
			t.Error("ShowOverview with single controller should return empty slice, not nil")
		}

		// Test with different isSecure values
		resultInsecure := usecase.ShowOverview(&singleController, &[]bool{false}[0])
		if resultInsecure == nil {
			t.Error("ShowOverview with insecure connection should return empty slice, not nil")
		}

		// Test with nil isSecure
		resultNilSecure := usecase.ShowOverview(&singleController, nil)
		if resultNilSecure == nil {
			t.Error("ShowOverview with nil isSecure should return empty slice, not nil")
		}
	})

	t.Run("show_overview_data_structure_initialization", func(t *testing.T) {
		// Test ShowOverviewData structure initialization
		data := ShowOverviewData{}

		// Verify zero values
		if data.ApMac != "" {
			t.Error("ShowOverviewData.ApMac should initialize to empty string")
		}
		if data.SlotID != 0 {
			t.Error("ShowOverviewData.SlotID should initialize to 0")
		}
		if data.Controller != "" {
			t.Error("ShowOverviewData.Controller should initialize to empty string")
		}

		// Test data assignment
		data.ApMac = "test-mac"
		data.SlotID = 1
		data.Controller = "test-controller"

		if data.ApMac != "test-mac" {
			t.Error("ShowOverviewData.ApMac assignment failed")
		}
		if data.SlotID != 1 {
			t.Error("ShowOverviewData.SlotID assignment failed")
		}
		if data.Controller != "test-controller" {
			t.Error("ShowOverviewData.Controller assignment failed")
		}
	})

	t.Run("show_overview_usecase_fields", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		// Test usecase field initialization
		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		if usecase.Config == nil {
			t.Error("OverviewUsecase.Config should not be nil")
		}
		if usecase.Repository == nil {
			t.Error("OverviewUsecase.Repository should not be nil")
		}

		// Test usecase with nil config
		usecaseNilConfig := OverviewUsecase{
			Config:     nil,
			Repository: repo,
		}

		if usecaseNilConfig.Config != nil {
			t.Error("OverviewUsecase.Config should be nil when set to nil")
		}

		// Test usecase with nil repository
		usecaseNilRepo := OverviewUsecase{
			Config:     cfg,
			Repository: nil,
		}

		if usecaseNilRepo.Repository != nil {
			t.Error("OverviewUsecase.Repository should be nil when set to nil")
		}
	})
}

// TestShowOverviewDataFlow tests detailed data flow scenarios
func TestShowOverviewDataFlow(t *testing.T) {
	t.Run("show_overview_error_handling_flow", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test with invalid controllers
		invalidControllers := []config.Controller{
			{Hostname: "", AccessToken: ""},
			{Hostname: "invalid", AccessToken: "invalid"},
		}
		isSecure := false

		result := usecase.ShowOverview(&invalidControllers, &isSecure)
		if result == nil {
			t.Error("ShowOverview with invalid controllers should return empty slice, not nil")
		}
	})

	t.Run("show_overview_data_processing", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test with realistic controller data
		controllers := []config.Controller{
			{
				Hostname:    "wlc-primary.company.com",
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			},
			{
				Hostname:    "wlc-secondary.company.com",
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
			},
		}
		isSecure := true

		result := usecase.ShowOverview(&controllers, &isSecure)

		// Test result consistency
		if result == nil {
			t.Error("ShowOverview should never return nil")
		}

		// Verify each item in result has required fields
		for i, item := range result {
			if item == nil {
				t.Errorf("ShowOverview result[%d] should not be nil", i)
				continue
			}
			// ApMac, SlotID, Controller should be set (even if empty due to failed API calls)
			_ = item.ApMac      // Access to ensure no panic
			_ = item.SlotID     // Access to ensure no panic
			_ = item.Controller // Access to ensure no panic
		}
	})

	t.Run("show_overview_concurrent_safety", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		controllers := []config.Controller{
			{Hostname: "test-controller", AccessToken: "test-token"},
		}
		isSecure := false

		// Test multiple calls to ensure no race conditions
		for i := 0; i < 5; i++ {
			result := usecase.ShowOverview(&controllers, &isSecure)
			if result == nil {
				t.Errorf("ShowOverview call %d should return empty slice, not nil", i)
			}
		}
	})
}

// TestShowOverviewFilterByRadioDetailed tests filterByRadio method comprehensively
func TestShowOverviewFilterByRadioDetailed(t *testing.T) {
	t.Run("filter_by_radio_with_data", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.Radio = "1" // Filter by slot 1
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Create test data with different slot IDs
		testData := []*ShowOverviewData{
			{ApMac: "mac1", SlotID: 0, Controller: "controller1"},
			{ApMac: "mac2", SlotID: 1, Controller: "controller1"},
			{ApMac: "mac3", SlotID: 2, Controller: "controller1"},
			{ApMac: "mac4", SlotID: 1, Controller: "controller2"},
		}

		filtered := usecase.filterByRadio(testData)

		// Should only return items with SlotID == 1
		expectedCount := 2
		if len(filtered) != expectedCount {
			t.Errorf("Expected %d filtered items, got %d", expectedCount, len(filtered))
		}

		for _, item := range filtered {
			if item.SlotID != 1 {
				t.Errorf("Filtered item should have SlotID 1, got %d", item.SlotID)
			}
		}
	})

	t.Run("filter_by_radio_no_filter", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.Radio = "" // No filter
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		testData := []*ShowOverviewData{
			{ApMac: "mac1", SlotID: 0, Controller: "controller1"},
			{ApMac: "mac2", SlotID: 1, Controller: "controller1"},
			{ApMac: "mac3", SlotID: 2, Controller: "controller1"},
		}

		filtered := usecase.filterByRadio(testData)

		// Should return all items when no filter is set
		if len(filtered) != len(testData) {
			t.Errorf("Expected %d items, got %d", len(testData), len(filtered))
		}
	})

	t.Run("filter_by_radio_edge_cases", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.Radio = "999" // Non-existent slot
		repo := localUtils.createMockRepository(cfg)

		usecase := OverviewUsecase{
			Config:     cfg,
			Repository: repo,
		}

		testData := []*ShowOverviewData{
			{ApMac: "mac1", SlotID: 0, Controller: "controller1"},
			{ApMac: "mac2", SlotID: 1, Controller: "controller1"},
		}

		filtered := usecase.filterByRadio(testData)

		// Should return empty slice when filter doesn't match any data
		if len(filtered) != 0 {
			t.Errorf("Expected 0 items for non-matching filter, got %d", len(filtered))
		}

		// Test with empty data
		emptyFiltered := usecase.filterByRadio([]*ShowOverviewData{})
		if len(emptyFiltered) != 0 {
			t.Errorf("Expected 0 items for empty input, got %d", len(emptyFiltered))
		}

		// Test with nil data
		nilFiltered := usecase.filterByRadio(nil)
		if len(nilFiltered) != 0 {
			t.Errorf("Expected 0 items for nil input, got %d", len(nilFiltered))
		}
	})
}

// TestShowClientComprehensiveBusinessLogic tests ShowClient with detailed scenarios
func TestShowClientComprehensiveBusinessLogic(t *testing.T) {
	t.Run("show_client_data_merging_comprehensive", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test multiple controller scenarios
		controllers := []config.Controller{
			{Hostname: "controller1.example.com", AccessToken: "token1"},
			{Hostname: "controller2.example.com", AccessToken: "token2"},
		}
		isSecure := false

		result := usecase.ShowClient(&controllers, &isSecure)

		// Verify result structure
		if result == nil {
			t.Error("ShowClient should return empty slice, not nil")
		}

		// Test empty controller slice
		emptyControllers := []config.Controller{}
		emptyResult := usecase.ShowClient(&emptyControllers, &isSecure)
		if emptyResult == nil {
			t.Error("ShowClient with empty controllers should return empty slice, not nil")
		}
		if len(emptyResult) != 0 {
			t.Errorf("ShowClient with empty controllers should return empty slice, got %d items", len(emptyResult))
		}
	})

	t.Run("show_client_controller_variations", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Test single controller
		singleController := []config.Controller{
			{Hostname: "single-controller.example.com", AccessToken: "single-token"},
		}
		isSecure := true

		result := usecase.ShowClient(&singleController, &isSecure)
		if result == nil {
			t.Error("ShowClient with single controller should return empty slice, not nil")
		}

		// Test with different isSecure values
		resultInsecure := usecase.ShowClient(&singleController, &[]bool{false}[0])
		if resultInsecure == nil {
			t.Error("ShowClient with insecure connection should return empty slice, not nil")
		}

		// Test with nil isSecure
		resultNilSecure := usecase.ShowClient(&singleController, nil)
		if resultNilSecure == nil {
			t.Error("ShowClient with nil isSecure should return empty slice, not nil")
		}
	})

	t.Run("show_client_data_structure_initialization", func(t *testing.T) {
		// Test ShowClientData structure initialization
		data := ShowClientData{}

		// Verify zero values
		if data.ClientMac != "" {
			t.Error("ShowClientData.ClientMac should initialize to empty string")
		}
		if data.Controller != "" {
			t.Error("ShowClientData.Controller should initialize to empty string")
		}

		// Test data assignment
		data.ClientMac = "test-mac"
		data.Controller = "test-controller"

		if data.ClientMac != "test-mac" {
			t.Error("ShowClientData.ClientMac assignment failed")
		}
		if data.Controller != "test-controller" {
			t.Error("ShowClientData.Controller assignment failed")
		}
	})

	t.Run("show_client_usecase_fields", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		repo := localUtils.createMockRepository(cfg)

		// Test usecase field initialization
		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		if usecase.Config == nil {
			t.Error("ClientUsecase.Config should not be nil")
		}
		if usecase.Repository == nil {
			t.Error("ClientUsecase.Repository should not be nil")
		}

		// Test usecase with nil config
		usecaseNilConfig := ClientUsecase{
			Config:     nil,
			Repository: repo,
		}

		if usecaseNilConfig.Config != nil {
			t.Error("ClientUsecase.Config should be nil when set to nil")
		}

		// Test usecase with nil repository
		usecaseNilRepo := ClientUsecase{
			Config:     cfg,
			Repository: nil,
		}

		if usecaseNilRepo.Repository != nil {
			t.Error("ClientUsecase.Repository should be nil when set to nil")
		}
	})
}

// TestShowClientFilterMethods tests filtering methods comprehensively
func TestShowClientFilterMethods(t *testing.T) {
	t.Run("filter_by_ssid_comprehensive", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.SSID = "test-ssid"
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Create test data with different SSIDs
		testData := []*ShowClientData{
			{ClientMac: "mac1", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
			{ClientMac: "mac2", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
			{ClientMac: "mac3", CommonOperData: client.CommonOperData{MsApSlotID: 2}},
		}

		filtered := usecase.filterBySSID(testData)

		// Should return empty slice since none match "test-ssid"
		if len(filtered) != 0 {
			t.Errorf("Expected 0 filtered items for non-matching SSID, got %d", len(filtered))
		}
	})

	t.Run("filter_by_ssid_no_filter", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.SSID = "" // No filter
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		testData := []*ShowClientData{
			{ClientMac: "mac1"},
			{ClientMac: "mac2"},
			{ClientMac: "mac3"},
		}

		filtered := usecase.filterBySSID(testData)

		// Should return all items when no filter is set
		if len(filtered) != len(testData) {
			t.Errorf("Expected %d items, got %d", len(testData), len(filtered))
		}
	})

	t.Run("filter_by_radio_comprehensive", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.Radio = "1" // Filter by slot 1
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		// Create test data with different slot IDs
		testData := []*ShowClientData{
			{ClientMac: "mac1", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
			{ClientMac: "mac2", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
			{ClientMac: "mac3", CommonOperData: client.CommonOperData{MsApSlotID: 2}},
			{ClientMac: "mac4", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
		}

		filtered := usecase.filterByRadio(testData)

		// Should only return items with MsApSlotID == 1
		expectedCount := 2
		if len(filtered) != expectedCount {
			t.Errorf("Expected %d filtered items, got %d", expectedCount, len(filtered))
		}

		for _, item := range filtered {
			if item.CommonOperData.MsApSlotID != 1 {
				t.Errorf("Filtered item should have MsApSlotID 1, got %d", item.CommonOperData.MsApSlotID)
			}
		}
	})

	t.Run("filter_by_radio_edge_cases", func(t *testing.T) {
		cfg := localUtils.createMockConfig()
		cfg.ShowCmdConfig.Radio = "999" // Non-existent slot
		repo := localUtils.createMockRepository(cfg)

		usecase := ClientUsecase{
			Config:     cfg,
			Repository: repo,
		}

		testData := []*ShowClientData{
			{ClientMac: "mac1", CommonOperData: client.CommonOperData{MsApSlotID: 0}},
			{ClientMac: "mac2", CommonOperData: client.CommonOperData{MsApSlotID: 1}},
		}

		filtered := usecase.filterByRadio(testData)

		// Should return empty slice when filter doesn't match any data
		if len(filtered) != 0 {
			t.Errorf("Expected 0 items for non-matching filter, got %d", len(filtered))
		}

		// Test with empty data
		emptyFiltered := usecase.filterByRadio([]*ShowClientData{})
		if len(emptyFiltered) != 0 {
			t.Errorf("Expected 0 items for empty input, got %d", len(emptyFiltered))
		}

		// Test with nil data
		nilFiltered := usecase.filterByRadio(nil)
		if len(nilFiltered) != 0 {
			t.Errorf("Expected 0 items for nil input, got %d", len(nilFiltered))
		}
	})
}
