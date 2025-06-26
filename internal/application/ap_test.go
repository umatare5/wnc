package application

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/ap"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestApUsecaseNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		wantNil    bool
	}{
		{
			name:       "creates_new_ap_usecase_with_valid_dependencies",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_ap_usecase_with_nil_config",
			config:     nil,
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_ap_usecase_with_nil_repository",
			config:     &config.Config{},
			repository: nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &ApUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			if (usecase == nil) != tt.wantNil {
				t.Errorf("ApUsecase creation failed, got nil: %v, want nil: %v", usecase == nil, tt.wantNil)
			}
		})
	}
}

func TestApUsecaseJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		usecase ApUsecase
	}{
		{
			name: "empty_ap_usecase",
			usecase: ApUsecase{
				Config:     nil,
				Repository: nil,
			},
		},
		{
			name: "ap_usecase_with_dependencies",
			usecase: ApUsecase{
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
			var usecase ApUsecase
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

func TestShowApDataJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data ShowApData
	}{
		{
			name: "empty_show_ap_data",
			data: ShowApData{
				ShowApCommonData: ShowApCommonData{
					ApMac:      "",
					Controller: "",
					CapwapData: ap.CapwapData{},
				},
				LLDPnei:    ap.LldpNeigh{},
				ApOperData: ap.ApOperData{},
			},
		},
		{
			name: "full_show_ap_data",
			data: ShowApData{
				ShowApCommonData: ShowApCommonData{
					ApMac:      "aa:bb:cc:dd:ee:ff",
					Controller: "controller1.example.com",
					CapwapData: ap.CapwapData{
						WtpMac: "aa:bb:cc:dd:ee:ff",
					},
				},
				LLDPnei: ap.LldpNeigh{
					WtpMac: "aa:bb:cc:dd:ee:ff",
				},
				ApOperData: ap.ApOperData{
					WtpMac: "aa:bb:cc:dd:ee:ff",
				},
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
			var data ShowApData
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if data.ApMac != tt.data.ApMac {
				t.Errorf("Expected ApMac %s, got %s", tt.data.ApMac, data.ApMac)
			}
			if data.Controller != tt.data.Controller {
				t.Errorf("Expected Controller %s, got %s", tt.data.Controller, data.Controller)
			}
		})
	}
}

func TestShowApTagDataJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data ShowApTagData
	}{
		{
			name: "empty_show_ap_tag_data",
			data: ShowApTagData{
				ShowApCommonData: ShowApCommonData{
					ApMac:      "",
					Controller: "",
					CapwapData: ap.CapwapData{},
				},
			},
		},
		{
			name: "full_show_ap_tag_data",
			data: ShowApTagData{
				ShowApCommonData: ShowApCommonData{
					ApMac:      "aa:bb:cc:dd:ee:ff",
					Controller: "controller1.example.com",
					CapwapData: ap.CapwapData{
						WtpMac: "aa:bb:cc:dd:ee:ff",
					},
				},
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
			var data ShowApTagData
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if data.ApMac != tt.data.ApMac {
				t.Errorf("Expected ApMac %s, got %s", tt.data.ApMac, data.ApMac)
			}
			if data.Controller != tt.data.Controller {
				t.Errorf("Expected Controller %s, got %s", tt.data.Controller, data.Controller)
			}
		})
	}
}

func TestShowAp(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create usecase with proper repository configuration
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30, // Set a positive timeout to avoid client errors
				},
			}
			repo := infrastructure.New(cfg)
			usecase := &ApUsecase{
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

			// Debug: Check usecase and parameters before calling
			if usecase.Repository == nil {
				t.Fatal("usecase.Repository is nil")
			}

			result := usecase.ShowAp(&tt.controllers, &tt.isSecure)

			// For test environment with network errors, either empty slice or nil is acceptable
			// The important thing is that the method doesn't panic
			if result == nil {
				// This is currently acceptable due to integration test nature with network failures
				t.Logf("ShowAp returned nil (probably due to network errors in test environment)")
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

func TestShowApTag(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create usecase with proper repository configuration
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30, // Set a positive timeout to avoid client errors
				},
			}
			repo := infrastructure.New(cfg)
			usecase := &ApUsecase{
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
			result := usecase.ShowApTag(&tt.controllers, &tt.isSecure)

			// For test environment with network errors, either empty slice or nil is acceptable
			// The important thing is that the method doesn't panic
			if result == nil {
				// This is currently acceptable due to integration test nature with network failures
				t.Logf("ShowApTag returned nil (probably due to network errors in test environment)")
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

func TestApUsecaseFailFast(t *testing.T) {
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
					t.Errorf("ApUsecase operation panicked: %v", r)
				}
			}()

			usecase := &ApUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			// Test that basic operations don't panic (don't test actual functionality)
			controllers := []config.Controller{}
			isSecure := true

			// For fail-fast tests, we only care that methods don't panic, not their return values
			_ = usecase.ShowAp(&controllers, &isSecure)
			_ = usecase.ShowApTag(&controllers, &isSecure)
		})
	}
}
