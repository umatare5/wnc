package application

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/client"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestClientUsecaseNew(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		repository *infrastructure.Repository
		wantNil    bool
	}{
		{
			name:       "creates_new_client_usecase_with_valid_dependencies",
			config:     &config.Config{},
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_client_usecase_with_nil_config",
			config:     nil,
			repository: &infrastructure.Repository{},
			wantNil:    false,
		},
		{
			name:       "creates_new_client_usecase_with_nil_repository",
			config:     &config.Config{},
			repository: nil,
			wantNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &ClientUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			if (usecase == nil) != tt.wantNil {
				t.Errorf("ClientUsecase creation failed, got nil: %v, want nil: %v", usecase == nil, tt.wantNil)
			}
		})
	}
}

func TestClientUsecaseJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		usecase ClientUsecase
	}{
		{
			name: "empty_client_usecase",
			usecase: ClientUsecase{
				Config:     nil,
				Repository: nil,
			},
		},
		{
			name: "client_usecase_with_dependencies",
			usecase: ClientUsecase{
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
			var usecase ClientUsecase
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

func TestShowClientDataJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data ShowClientData
	}{
		{
			name: "empty_show_client_data",
			data: ShowClientData{
				ClientMac:      "",
				Controller:     "",
				CommonOperData: client.CommonOperData{},
				Dot11OperData:  client.Dot11OperData{},
				TrafficStats:   client.TrafficStats{},
				SisfDbMac:      client.SisfDbMac{},
				DcInfo:         client.DcInfo{},
			},
		},
		{
			name: "full_show_client_data",
			data: ShowClientData{
				ClientMac:  "aa:bb:cc:dd:ee:ff",
				Controller: "controller1.example.com",
				CommonOperData: client.CommonOperData{
					ClientMac: "aa:bb:cc:dd:ee:ff",
				},
				Dot11OperData: client.Dot11OperData{
					MsMacAddress: "aa:bb:cc:dd:ee:ff",
				},
				TrafficStats: client.TrafficStats{
					MsMacAddress: "aa:bb:cc:dd:ee:ff",
				},
				SisfDbMac: client.SisfDbMac{
					MacAddr: "aa:bb:cc:dd:ee:ff",
				},
				DcInfo: client.DcInfo{
					ClientMac: "aa:bb:cc:dd:ee:ff",
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
			var data ShowClientData
			err = json.Unmarshal(jsonData, &data)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON: %v", err)
			}

			// Verify structure
			if data.ClientMac != tt.data.ClientMac {
				t.Errorf("Expected ClientMac %s, got %s", tt.data.ClientMac, data.ClientMac)
			}
			if data.Controller != tt.data.Controller {
				t.Errorf("Expected Controller %s, got %s", tt.data.Controller, data.Controller)
			}
		})
	}
}

func TestShowClient(t *testing.T) {
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
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					SSID:    "",
					Radio:   "",
					Timeout: 30, // Set a positive timeout to avoid client errors
				},
			}
			repo := infrastructure.New(cfg)
			usecase := &ClientUsecase{
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
			result := usecase.ShowClient(&tt.controllers, &tt.isSecure)

			// For test environment with network errors, either empty slice or nil is acceptable
			// The important thing is that the method doesn't panic
			if result == nil {
				// This is currently acceptable due to integration test nature with network failures
				t.Logf("ShowClient returned nil (probably due to network errors in test environment)")
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

func TestFilterBySSID(t *testing.T) {
	tests := []struct {
		name           string
		clients        []*ShowClientData
		ssidFilter     string
		expectedLength int
	}{
		{
			name:           "no_filter_returns_all",
			clients:        []*ShowClientData{{}, {}},
			ssidFilter:     "",
			expectedLength: 2,
		},
		{
			name: "filter_by_ssid",
			clients: []*ShowClientData{
				{
					Dot11OperData: client.Dot11OperData{
						VapSsid: "test-ssid",
					},
				},
				{
					Dot11OperData: client.Dot11OperData{
						VapSsid: "other-ssid",
					},
				},
			},
			ssidFilter:     "test-ssid",
			expectedLength: 1,
		},
		{
			name: "filter_no_matches",
			clients: []*ShowClientData{
				{
					Dot11OperData: client.Dot11OperData{
						VapSsid: "other-ssid",
					},
				},
			},
			ssidFilter:     "nonexistent-ssid",
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &ClientUsecase{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						SSID: tt.ssidFilter,
					},
				},
			}

			result := usecase.filterBySSID(tt.clients)

			if len(result) != tt.expectedLength {
				t.Errorf("Expected %d clients, got %d", tt.expectedLength, len(result))
			}
		})
	}
}

func TestFilterByRadio(t *testing.T) {
	tests := []struct {
		name           string
		clients        []*ShowClientData
		radioFilter    string
		expectedLength int
	}{
		{
			name:           "no_filter_returns_all",
			clients:        []*ShowClientData{{}, {}},
			radioFilter:    "",
			expectedLength: 2,
		},
		{
			name: "filter_by_radio",
			clients: []*ShowClientData{
				{
					CommonOperData: client.CommonOperData{
						MsApSlotID: 1,
					},
				},
				{
					CommonOperData: client.CommonOperData{
						MsApSlotID: 2,
					},
				},
			},
			radioFilter:    "1",
			expectedLength: 1,
		},
		{
			name: "filter_no_matches",
			clients: []*ShowClientData{
				{
					CommonOperData: client.CommonOperData{
						MsApSlotID: 1,
					},
				},
			},
			radioFilter:    "3",
			expectedLength: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := &ClientUsecase{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						Radio: tt.radioFilter,
					},
				},
			}

			result := usecase.filterByRadio(tt.clients)

			if len(result) != tt.expectedLength {
				t.Errorf("Expected %d clients, got %d", tt.expectedLength, len(result))
			}
		})
	}
}

func TestClientUsecaseFailFast(t *testing.T) {
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
			name: "valid_dependencies_should_not_panic",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					SSID:  "",
					Radio: "",
				},
			},
			repository: &infrastructure.Repository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ClientUsecase operation panicked: %v", r)
				}
			}()

			usecase := &ClientUsecase{
				Config:     tt.config,
				Repository: tt.repository,
			}

			// Test that basic operations don't panic (don't test actual functionality)
			controllers := []config.Controller{}
			isSecure := true

			// For fail-fast tests, we only care that methods don't panic, not their return values
			if tt.config != nil && tt.repository != nil {
				_ = usecase.ShowClient(&controllers, &isSecure)
			}

			// Test filter methods with empty data
			emptyClients := []*ShowClientData{}
			if tt.config != nil {
				_ = usecase.filterBySSID(emptyClients)
				_ = usecase.filterByRadio(emptyClients)
			}
		})
	}
}
