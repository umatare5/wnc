package application

import (
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/client"
	"github.com/umatare5/wnc/internal/config"
)

// TestClientUsecaseInitialization tests ClientUsecase initialization (Unit test)
func TestClientUsecaseInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "invoke_client_usecase",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
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

// TestClientUsecaseShowClient tests ClientUsecase ShowClient method (Unit test)
func TestClientUsecaseShowClient(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "ShowClient_with_nil_repository",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				app := New(cfg, nil)
				clientUsecase := app.InvokeClientUsecase()

				controllers := &[]config.Controller{{Hostname: "test", AccessToken: "token"}}
				isSecure := false
				result := clientUsecase.ShowClient(controllers, &isSecure)

				if result == nil {
					t.Error("ShowClient should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowClient should return empty slice when repository is nil")
				}
			},
		},
		{
			name: "ShowClient_with_nil_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				isSecure := false
				result := clientUsecase.ShowClient(nil, &isSecure)

				if result == nil {
					t.Error("ShowClient should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowClient should return empty slice when controllers is nil")
				}
			},
		},
		{
			name: "ShowClient_with_empty_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				controllers := &[]config.Controller{}
				isSecure := false
				result := clientUsecase.ShowClient(controllers, &isSecure)

				if result == nil {
					t.Error("ShowClient should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowClient should return empty slice when controllers is empty")
				}
			},
		},
		{
			name: "ShowClient_with_valid_controllers_but_nil_data",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				controllers := &[]config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}
				isSecure := false
				result := clientUsecase.ShowClient(controllers, &isSecure)

				// Should not panic and should return empty slice when API returns nil
				if result == nil {
					t.Error("ShowClient should return empty slice, not nil")
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

// TestClientFilterFunctions tests client filter helper functions (Unit test)
func TestClientFilterFunctions(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "filterBySSID_function",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "TestSSID1"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				// Create test data
				clients := []*ShowClientData{
					{Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID1"}},
					{Dot11OperData: client.Dot11OperData{VapSsid: "TestSSID2"}},
					{Dot11OperData: client.Dot11OperData{VapSsid: "OtherSSID"}},
				}

				// Test filtering
				filtered := clientUsecase.filterBySSID(clients)
				if len(filtered) != 1 {
					t.Errorf("Expected 1 client, got %d", len(filtered))
				}
				if len(filtered) > 0 && filtered[0].Dot11OperData.VapSsid != "TestSSID1" {
					t.Error("Filtered client has wrong SSID")
				}

				// Test empty filter
				cfg.ShowCmdConfig.SSID = ""
				all := clientUsecase.filterBySSID(clients)
				if len(all) != len(clients) {
					t.Error("Empty filter should return all clients")
				}

				// Test with nil config
				clientUsecase.Config = nil
				nilConfig := clientUsecase.filterBySSID(clients)
				if len(nilConfig) != len(clients) {
					t.Error("Nil config should return all clients")
				}
			},
		},
		{
			name: "filterByRadio_function",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "1"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				// Create test data with CommonOperData containing MsApSlotID
				clients := []*ShowClientData{
					{CommonOperData: client.CommonOperData{MsApSlotID: 1}},
					{CommonOperData: client.CommonOperData{MsApSlotID: 2}},
					{CommonOperData: client.CommonOperData{MsApSlotID: 3}},
				}

				// Test filtering
				filtered := clientUsecase.filterByRadio(clients)
				if len(filtered) != 1 {
					t.Errorf("Expected 1 client, got %d", len(filtered))
				}
				if len(filtered) > 0 && filtered[0].CommonOperData.MsApSlotID != 1 {
					t.Error("Filtered client has wrong radio slot ID")
				}

				// Test empty filter
				cfg.ShowCmdConfig.Radio = ""
				all := clientUsecase.filterByRadio(clients)
				if len(all) != len(clients) {
					t.Error("Empty filter should return all clients")
				}

				// Test with nil config
				clientUsecase.Config = nil
				nilConfig := clientUsecase.filterByRadio(clients)
				if len(nilConfig) != len(clients) {
					t.Error("Nil config should return all clients")
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

// TestClientUsecaseDataProcessingAndFiltering tests comprehensive data processing and filtering (Unit test)
func TestClientUsecaseDataProcessingAndFiltering(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "complex_filtering_with_ssid_and_radio",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "Guest-Network"
				cfg.ShowCmdConfig.Radio = "1"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						Dot11OperData:  client.Dot11OperData{VapSsid: "Guest-Network"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
					{
						ClientMac:      "00:11:22:33:44:66",
						Dot11OperData:  client.Dot11OperData{VapSsid: "Corporate-Network"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
					{
						ClientMac:      "00:11:22:33:44:77",
						Dot11OperData:  client.Dot11OperData{VapSsid: "Guest-Network"},
						CommonOperData: client.CommonOperData{MsApSlotID: 0},
					},
				}

				// Apply both filters
				ssidFiltered := clientUsecase.filterBySSID(testData)
				radioFiltered := clientUsecase.filterByRadio(ssidFiltered)

				if len(radioFiltered) != 1 {
					t.Errorf("Expected 1 client after both filters, got %d", len(radioFiltered))
				}
				if len(radioFiltered) > 0 {
					result := radioFiltered[0]
					if result.Dot11OperData.VapSsid != "Guest-Network" {
						t.Error("Filtered client should have Guest-Network SSID")
					}
					if result.CommonOperData.MsApSlotID != 1 {
						t.Error("Filtered client should have radio slot 1")
					}
				}
			},
		},
		{
			name: "empty_data_handling",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				emptyData := []*ShowClientData{}

				// Test all filters with empty data
				ssidResult := clientUsecase.filterBySSID(emptyData)
				if ssidResult == nil || len(ssidResult) != 0 {
					t.Error("filterBySSID should return empty slice for empty input")
				}

				radioResult := clientUsecase.filterByRadio(emptyData)
				if radioResult == nil || len(radioResult) != 0 {
					t.Error("filterByRadio should return empty slice for empty input")
				}
			},
		},
		{
			name: "nil_data_handling",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				// Test all filters with nil data
				ssidResult := clientUsecase.filterBySSID(nil)
				if ssidResult != nil && len(ssidResult) != 0 {
					t.Error("filterBySSID should handle nil input gracefully")
				}

				radioResult := clientUsecase.filterByRadio(nil)
				if radioResult != nil && len(radioResult) != 0 {
					t.Error("filterByRadio should handle nil input gracefully")
				}
			},
		},
		{
			name: "edge_case_radio_slot_values",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "999" // Large slot ID
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						CommonOperData: client.CommonOperData{MsApSlotID: 999},
					},
					{
						ClientMac:      "00:11:22:33:44:66",
						CommonOperData: client.CommonOperData{MsApSlotID: 0},
					},
				}

				result := clientUsecase.filterByRadio(testData)
				if len(result) != 1 {
					t.Errorf("Expected 1 client with slot 999, got %d", len(result))
				}
				if len(result) > 0 && result[0].CommonOperData.MsApSlotID != 999 {
					t.Error("Filtered client should have slot ID 999")
				}
			},
		},
		{
			name: "ssid_case_sensitivity",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "TestNetwork" // Exact case
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:     "00:11:22:33:44:55",
						Dot11OperData: client.Dot11OperData{VapSsid: "TestNetwork"},
					},
					{
						ClientMac:     "00:11:22:33:44:66",
						Dot11OperData: client.Dot11OperData{VapSsid: "testnetwork"},
					},
					{
						ClientMac:     "00:11:22:33:44:77",
						Dot11OperData: client.Dot11OperData{VapSsid: "TESTNETWORK"},
					},
				}

				result := clientUsecase.filterBySSID(testData)
				if len(result) != 1 {
					t.Errorf("SSID filter should be case-sensitive, expected 1 match, got %d", len(result))
				}
				if len(result) > 0 && result[0].Dot11OperData.VapSsid != "TestNetwork" {
					t.Error("Filtered client should have exact case match")
				}
			},
		},
		{
			name: "filter_chain_preserves_data_integrity",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "TestSSID"
				cfg.ShowCmdConfig.Radio = "2"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				originalData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						Dot11OperData:  client.Dot11OperData{VapSsid: "TestSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 2},
					},
				}

				// Apply filters and verify data integrity
				ssidFiltered := clientUsecase.filterBySSID(originalData)
				radioFiltered := clientUsecase.filterByRadio(ssidFiltered)

				if len(radioFiltered) != 1 {
					t.Error("Filter chain should preserve matching data")
				}
				if len(radioFiltered) > 0 {
					result := radioFiltered[0]
					if result.ClientMac != "00:11:22:33:44:55" {
						t.Error("Filter chain should preserve all original data fields")
					}
				}
			},
		},
		{
			name: "data_processing_with_various_client_macs",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						Dot11OperData:  client.Dot11OperData{VapSsid: "TestSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
					{
						ClientMac:      "00:11:22:33:44:66",
						Dot11OperData:  client.Dot11OperData{VapSsid: "TestSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
					{
						ClientMac:      "00:11:22:33:44:77",
						Dot11OperData:  client.Dot11OperData{VapSsid: "TestSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
				}

				// Test that all data is preserved with same SSID
				result := clientUsecase.filterBySSID(testData)
				if len(result) != 3 {
					t.Errorf("Expected all 3 clients with TestSSID, got %d", len(result))
				}

				// Verify each client MAC is preserved
				clientMacs := make(map[string]bool)
				for _, client := range result {
					clientMacs[client.ClientMac] = true
				}

				expectedMacs := []string{"00:11:22:33:44:55", "00:11:22:33:44:66", "00:11:22:33:44:77"}
				for _, mac := range expectedMacs {
					if !clientMacs[mac] {
						t.Errorf("Expected to find client MAC %s", mac)
					}
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

// TestClientUsecaseAdvancedFiltering tests advanced filtering scenarios (Unit test)
func TestClientUsecaseAdvancedFiltering(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "filter_by_ssid_with_special_characters",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "Guest-WiFi_2024!" // SSID with special characters
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:     "00:11:22:33:44:55",
						Dot11OperData: client.Dot11OperData{VapSsid: "Guest-WiFi_2024!"},
					},
					{
						ClientMac:     "00:11:22:33:44:66",
						Dot11OperData: client.Dot11OperData{VapSsid: "Guest-WiFi_2024"},
					},
				}

				result := clientUsecase.filterBySSID(testData)
				if len(result) != 1 {
					t.Errorf("Expected 1 client with special characters SSID, got %d", len(result))
				}
				if len(result) > 0 && result[0].Dot11OperData.VapSsid != "Guest-WiFi_2024!" {
					t.Error("Filtered client should have exact SSID with special characters")
				}
			},
		},
		{
			name: "filter_by_radio_with_invalid_slot_id",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "invalid" // Invalid radio slot
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
				}

				result := clientUsecase.filterByRadio(testData)
				// Should return empty slice when radio filter is invalid/non-matching
				if len(result) != 0 {
					t.Errorf("Expected no clients for invalid radio filter, got %d", len(result))
				}
			},
		},
		{
			name: "filter_chain_order_independence",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.SSID = "TestSSID"
				cfg.ShowCmdConfig.Radio = "1"
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				clientUsecase := app.InvokeClientUsecase()

				testData := []*ShowClientData{
					{
						ClientMac:      "00:11:22:33:44:55",
						Dot11OperData:  client.Dot11OperData{VapSsid: "TestSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
					{
						ClientMac:      "00:11:22:33:44:66",
						Dot11OperData:  client.Dot11OperData{VapSsid: "OtherSSID"},
						CommonOperData: client.CommonOperData{MsApSlotID: 1},
					},
				}

				// Apply filters in different orders
				order1 := clientUsecase.filterByRadio(clientUsecase.filterBySSID(testData))
				order2 := clientUsecase.filterBySSID(clientUsecase.filterByRadio(testData))

				if len(order1) != len(order2) {
					t.Error("Filter order should not affect final result count")
				}
				if len(order1) != 1 || len(order2) != 1 {
					t.Error("Both filter orders should result in 1 client")
				}
				if len(order1) > 0 && len(order2) > 0 {
					if order1[0].ClientMac != order2[0].ClientMac {
						t.Error("Filter order should not affect which client is returned")
					}
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
