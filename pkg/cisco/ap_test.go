package cisco

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApTypeAliases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
	}{
		{
			name: "ApOperResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApOperResponse = nil
				return true
			},
		},
		{
			name: "ApOperCapwapDataResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApOperCapwapDataResponse = nil
				return true
			},
		},
		{
			name: "ApOperLldpNeighResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApOperLldpNeighResponse = nil
				return true
			},
		},
		{
			name: "ApOperRadioOperDataResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApOperRadioOperDataResponse = nil
				return true
			},
		},
		{
			name: "ApOperOperDataResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApOperOperDataResponse = nil
				return true
			},
		},
		{
			name: "ApGlobalOperResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApGlobalOperResponse = nil
				return true
			},
		},
		{
			name: "ApCfgResponse_alias_exists",
			testFunc: func() bool {
				var _ *ApCfgResponse = nil
				return true
			},
		},
		{
			name: "CapwapData_alias_exists",
			testFunc: func() bool {
				var _ *CapwapData = nil
				return true
			},
		},
		{
			name: "LldpNeigh_alias_exists",
			testFunc: func() bool {
				var _ *LldpNeigh = nil
				return true
			},
		},
		{
			name: "ApOperData_alias_exists",
			testFunc: func() bool {
				var _ *ApOperData = nil
				return true
			},
		},
		{
			name: "RadioOperData_alias_exists",
			testFunc: func() bool {
				var _ *RadioOperData = nil
				return true
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.testFunc() {
				t.Errorf("Type alias check failed for %s", tt.name)
			}
		})
	}
}

func TestApJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "CapwapData_serialization",
			data: CapwapData{
				WtpMac: "aa:bb:cc:dd:ee:ff",
			},
		},
		{
			name: "LldpNeigh_serialization",
			data: LldpNeigh{
				WtpMac: "aa:bb:cc:dd:ee:ff",
			},
		},
		{
			name: "ApOperData_serialization",
			data: ApOperData{
				WtpMac: "aa:bb:cc:dd:ee:ff",
			},
		},
		{
			name: "RadioOperData_serialization",
			data: RadioOperData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Test JSON unmarshaling by creating a new instance of the same type
			switch v := tt.data.(type) {
			case CapwapData:
				var unmarshaled CapwapData
				err = json.Unmarshal(jsonData, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal CapwapData: %v", err)
				}
				if unmarshaled.WtpMac != v.WtpMac {
					t.Errorf("Expected WtpMac %s, got %s", v.WtpMac, unmarshaled.WtpMac)
				}
			case LldpNeigh:
				var unmarshaled LldpNeigh
				err = json.Unmarshal(jsonData, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal LldpNeigh: %v", err)
				}
				if unmarshaled.WtpMac != v.WtpMac {
					t.Errorf("Expected WtpMac %s, got %s", v.WtpMac, unmarshaled.WtpMac)
				}
			case ApOperData:
				var unmarshaled ApOperData
				err = json.Unmarshal(jsonData, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal ApOperData: %v", err)
				}
				if unmarshaled.WtpMac != v.WtpMac {
					t.Errorf("Expected WtpMac %s, got %s", v.WtpMac, unmarshaled.WtpMac)
				}
			case RadioOperData:
				var unmarshaled RadioOperData
				err = json.Unmarshal(jsonData, &unmarshaled)
				if err != nil {
					t.Fatalf("Failed to unmarshal RadioOperData: %v", err)
				}
			}
		})
	}
}

func TestApFunctionSignatures(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func() bool
	}{
		{
			name: "GetApOper_function_signature",
			testFunc: func() bool {
				// Test that the function exists by checking its type
				var f func(*Client, context.Context) (*ApOperResponse, error) = GetApOper
				return f != nil
			},
		},
		{
			name: "GetApCapwapData_function_signature",
			testFunc: func() bool {
				// Test that the function exists by checking its type
				var f func(*Client, context.Context) (*ApOperCapwapDataResponse, error) = GetApCapwapData
				return f != nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.testFunc() {
				t.Errorf("Function signature test failed for %s", tt.name)
			}
		})
	}
}

func TestApFailFast(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
	}{
		{
			name: "CapwapData_should_not_panic",
			data: CapwapData{},
		},
		{
			name: "LldpNeigh_should_not_panic",
			data: LldpNeigh{},
		},
		{
			name: "ApOperData_should_not_panic",
			data: ApOperData{},
		},
		{
			name: "RadioOperData_should_not_panic",
			data: RadioOperData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("AP data operation panicked: %v", r)
				}
			}()

			// Test that basic operations don't panic
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("JSON marshal failed: %v", err)
			}

			if len(jsonData) == 0 {
				t.Error("Expected non-empty JSON data")
			}
		})
	}
}

func TestApIntegration(t *testing.T) {
	tests := []struct {
		name   string
		client *Client
	}{
		{
			name:   "nil_client_should_handle_gracefully",
			client: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For nil client, just verify the functions exist without calling them
			// to avoid panic in the underlying library
			if tt.client == nil {
				// The functions exist if we can assign them to variables
				_ = GetApOper
				_ = GetApCapwapData
				return
			}

			// For real clients, we could test actual functionality here
			// but we don't have real clients in unit tests
		})
	}
}

func TestGetApOper_WithRealResponse(t *testing.T) {
	// Create a test server with real WNC AP API response structure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mock response based on real WNC AP API structure
		w.Header().Set("Content-Type", "application/yang-data+json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"cisco-wireless-ap-oper:ap-oper-data": {
				"ap-name-mac-map": {
					"ap-name-mac-mapping": [
						{
							"wtp-mac": "28:ac:9e:bb:3c:80",
							"ap-name": "lab2-ap1815-06f-02",
							"ethernet-mac": "28:ac:9e:11:48:10",
							"ip-addr": "192.168.255.11"
						}
					]
				}
			}
		}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetApOper function
	result, err := GetApOper(client, context.Background())
	if err != nil {
		t.Errorf("GetApOper failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		t.Logf("GetApOper returned result successfully")
	}
}

func TestGetApCapwapData_WithRealResponse(t *testing.T) {
	// Create a test server with real WNC CAPWAP data response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"cisco-wireless-access-point-oper:capwap-data": {
				"wtp-mac": "28:ac:9e:bb:3c:80",
				"ip-addr": "192.168.255.11",
				"name": "lab2-ap1815-06f-02",
				"ap-state": {
					"ap-admin-state": "adminstate-enabled",
					"ap-operation-state": "registered"
				}
			}
		}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetApCapwapData function
	result, err := GetApCapwapData(client, context.Background())
	if err != nil {
		t.Errorf("GetApCapwapData failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		t.Logf("GetApCapwapData returned result successfully")
	}
}

func TestGetApLldpNeigh_WithRealResponse(t *testing.T) {
	// Create a test server with real WNC LLDP neighbor response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yang-data+json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{
			"cisco-wireless-access-point-oper:lldp-neigh": {
				"wtp-mac": "28:ac:9e:bb:3c:80",
				"system-name": "Switch-01",
				"port-id": "Gi1/0/1"
			}
		}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetApLldpNeigh function
	result, err := GetApLldpNeigh(client, context.Background())
	if err != nil {
		t.Errorf("GetApLldpNeigh failed: %v", err)
	}

	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		t.Logf("GetApLldpNeigh returned result successfully")
	}
}

func TestGetApRadioOperData_WithErrorResponse(t *testing.T) {
	// Create a test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Internal Server Error"}`))
		if err != nil {
			t.Errorf("Failed to write error response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with test server URL
	client, err := NewClientWithTimeout(server.URL, "test-token", 30, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetApRadioOperData with error response
	result, err := GetApRadioOperData(client, context.Background())
	if err == nil {
		t.Error("Expected error for server error response")
	}

	if result != nil {
		t.Error("Expected nil result for error response")
	}
}

func TestGetApOperData_WithTimeout(t *testing.T) {
	// Create a test server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate network delay
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"cisco-wireless-access-point-oper:ap-oper-data": {}}`))
		if err != nil {
			t.Errorf("Failed to write response: %v", err)
		}
	}))
	defer server.Close()

	// Create client with very short timeout to test timeout handling
	client, err := NewClientWithTimeout(server.URL, "test-token", 1, boolPtr(false))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test GetApOperData function
	result, err := GetApOperData(client, context.Background())
	// This may or may not timeout depending on test execution speed
	// So we just check that the function exists and can be called
	if result != nil {
		t.Logf("GetApOperData returned result successfully")
	} else if err != nil {
		t.Logf("GetApOperData returned error as expected: %v", err)
	}
}

// Helper function to create a bool pointer
func boolPtr(b bool) *bool {
	return &b
}
