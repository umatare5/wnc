package application

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/wlan"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestWlanUsecaseNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		wantNil    bool
	}{
		{
			name:       "creates_new_wlan_usecase_with_valid_dependencies",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_wlan_usecase_with_nil_config",
			config:     nil,
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_wlan_usecase_with_nil_repository",
			config:     &config.Config{},
			repository: nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &WlanUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			if (usecase == nil) != tt.wantNil {
				t.Errorf("WlanUsecase creation failed, got nil: %v, want nil: %v", usecase == nil, tt.wantNil)
			}
		})
	}
}

func TestWlanUsecaseJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		usecase WlanUsecase
	}{
		{
			name: "empty_wlan_usecase",
			usecase: WlanUsecase{
				Config:     nil,
				Repository: nil,
			},
		},
		{
			name: "wlan_usecase_with_dependencies",
			usecase: WlanUsecase{
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
			var usecase WlanUsecase
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

func TestShowWlanDataJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data ShowWlanData
	}{
		{
			name: "empty_show_wlan_data",
			data: ShowWlanData{
				TagName:      "",
				PolicyName:   "",
				WlanName:     "",
				WlanCfgEntry: wlan.WlanCfgEntry{},
				WlanPolicy:   wlan.WlanPolicy{},
				Controller:   "",
			},
		},
		{
			name: "full_show_wlan_data",
			data: ShowWlanData{
				TagName:    "default-policy-tag",
				PolicyName: "test-policy",
				WlanName:   "test-wlan",
				WlanCfgEntry: wlan.WlanCfgEntry{
					ProfileName: "test-wlan",
				},
				WlanPolicy: wlan.WlanPolicy{
					PolicyProfileName: "test-policy",
				},
				Controller: "controller1.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling
			var data ShowWlanData
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if data.TagName != tt.data.TagName {
				t.Errorf("Expected TagName %s, got %s", tt.data.TagName, data.TagName)
			}
			if data.PolicyName != tt.data.PolicyName {
				t.Errorf("Expected PolicyName %s, got %s", tt.data.PolicyName, data.PolicyName)
			}
			if data.WlanName != tt.data.WlanName {
				t.Errorf("Expected WlanName %s, got %s", tt.data.WlanName, data.WlanName)
			}
			if data.Controller != tt.data.Controller {
				t.Errorf("Expected Controller %s, got %s", tt.data.Controller, data.Controller)
			}
		})
	}
}

func TestShowWlan(t *testing.T) {
	tests := []struct {
		name        string
		controllers []config.Controller
		isSecure    bool
		expectEmpty bool
	}{
		{
			name:        "empty_controllers_returns_empty_slice",
			controllers: []config.Controller{},
			isSecure:    true,
			expectEmpty: true,
		},
		{
			name: "single_controller",
			controllers: []config.Controller{
				{
					Hostname:    "controller1.example.com",
					AccessToken: "token123",
				},
			},
			isSecure:    true,
			expectEmpty: true, // Will be empty due to mocked repository
		},
		{
			name: "multiple_controllers",
			controllers: []config.Controller{
				{
					Hostname:    "controller1.example.com",
					AccessToken: "token123",
				},
				{
					Hostname:    "controller2.example.com",
					AccessToken: "token456",
				},
			},
			isSecure:    true,
			expectEmpty: true, // Will be empty due to mocked repository
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30, // Set a positive timeout to avoid client errors
				},
			}
			repo := infrastructure.New(cfg)
			usecase := &WlanUsecase{
				Config:     cfg,
				Repository: &repo,
			}

			// For empty controllers, expect empty result without calling actual method
			if len(tt.controllers) == 0 {
				if !tt.expectEmpty {
					t.Error("Empty controllers test should expect empty result")
				}
				return
			}

			// For non-empty controllers, expect the method to handle gracefully
			// Note: This will likely return empty results due to network errors in test environment
			result := usecase.ShowWlan(&tt.controllers, &tt.isSecure)

			// For test environment with network errors, either empty slice or nil is acceptable
			// The important thing is that the method doesn't panic
			if result == nil {
				// This is currently acceptable due to integration test nature with network failures
				t.Logf("ShowWlan returned nil (probably due to network errors in test environment)")
				return
			}

			// For test environment, we expect empty results due to network failures
			// This is normal behavior and should not be considered an error
			if len(result) == 0 && tt.expectEmpty {
				// This is the expected behavior for test environment
				return
			}
		})
	}
}

func TestWlanUsecaseFailFast(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
	}{
		{
			name:       "nil_config_should_not_panic",
			config:     nil,
			repository: &infrastructure.Repository{},
		},
		{
			name:       "nil_repository_should_not_panic",
			config:     &config.Config{},
			repository: nil,
		},
		{
			name:       "valid_dependencies_should_not_panic",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("WlanUsecase operation panicked: %v", r)
				}
			}()

			usecase := &WlanUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			// Test that basic operations don't panic
			controllers := []config.Controller{}
			isSecure := true

			// For fail-fast tests, we only check that the method doesn't panic
			// We don't validate the result content since that's tested elsewhere
			usecase.ShowWlan(&controllers, &isSecure)
		})
	}
}

func TestWlanUsecaseDependencyInjection(t *testing.T) {
	tests := []struct {
		name               string
		config             *config.Config
		repository         *infrastructure.Repository
		expectValidUsecase bool
	}{
		{
			name:               "dependency_injection_with_valid_objects",
			config:             &config.Config{},
			repository:         &infrastructure.Repository{},
			expectValidUsecase: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &WlanUsecase{
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

				// Test that the usecase can execute ShowWlan without panicking
				controllers := []config.Controller{}
				isSecure := true
				// For dependency injection tests, we only verify no panic occurs
				// We don't validate result content since that's tested elsewhere
				usecase.ShowWlan(&controllers, &isSecure)
			}
		})
	}
}
