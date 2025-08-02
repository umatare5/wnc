package integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/umatare5/wnc/pkg/cisco"
)

// TestWNCConnectivity tests basic connectivity to WNC controllers (Integration test)
func TestWNCConnectivity(t *testing.T) {
	controllers := os.Getenv("WNC_CONTROLLERS")
	if controllers == "" {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	// Parse controller configuration
	for _, controllerPair := range strings.Split(controllers, ",") {
		parts := strings.Split(controllerPair, ":")
		if len(parts) != 2 {
			t.Errorf("Invalid controller format: %s", controllerPair)
			continue
		}

		hostname := parts[0]
		token := parts[1]

		t.Run("test_connectivity_"+hostname, func(t *testing.T) {
			// Create client
			isSecure := false
			client, err := cisco.NewClient(hostname, token, &isSecure)
			if err != nil {
				t.Fatalf("Failed to create client for %s: %v", hostname, err)
			}

			// Test basic AP operation data retrieval
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			_, err = cisco.GetApOper(client, ctx)
			if err != nil {
				t.Errorf("Failed to get AP operational data from %s: %v", hostname, err)
			} else {
				t.Logf("Successfully retrieved AP data from %s", hostname)
			}

			// Test WLAN configuration retrieval
			_, err = cisco.GetWlanCfg(client, ctx)
			if err != nil {
				t.Errorf("Failed to get WLAN configuration from %s: %v", hostname, err)
			} else {
				t.Logf("Successfully retrieved WLAN configuration from %s", hostname)
			}
		})
	}
}

// TestWNCOperations tests comprehensive WNC operations (Integration test)
func TestWNCOperations(t *testing.T) {
	controllers := os.Getenv("WNC_CONTROLLERS")
	if controllers == "" {
		t.Skip("Skipping integration test: WNC_CONTROLLERS not set")
	}

	// Parse first controller for detailed testing
	controllerPair := strings.Split(controllers, ",")[0]
	parts := strings.Split(controllerPair, ":")
	if len(parts) != 2 {
		t.Fatalf("Invalid controller format: %s", controllerPair)
	}

	hostname := parts[0]
	token := parts[1]

	// Create client
	isSecure := false
	client, err := cisco.NewClient(hostname, token, &isSecure)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "get_ap_operational_data",
			test: func(t *testing.T) {
				response, err := cisco.GetApOper(client, ctx)
				if err != nil {
					t.Errorf("GetApOper failed: %v", err)
					return
				}
				if response == nil {
					t.Error("GetApOper returned nil response")
				}
				t.Logf("AP Operational data retrieved successfully")
			},
		},
		{
			name: "get_wlan_configuration",
			test: func(t *testing.T) {
				response, err := cisco.GetWlanCfg(client, ctx)
				if err != nil {
					t.Errorf("GetWlanCfg failed: %v", err)
					return
				}
				if response == nil {
					t.Error("GetWlanCfg returned nil response")
				}
				t.Logf("WLAN configuration retrieved successfully")
			},
		},
		{
			name: "get_client_operational_data",
			test: func(t *testing.T) {
				response, err := cisco.GetClientOper(client, ctx)
				if err != nil {
					t.Errorf("GetClientOper failed: %v", err)
					return
				}
				if response == nil {
					t.Error("GetClientOper returned nil response")
				}
				t.Logf("Client operational data retrieved successfully")
			},
		},
		{
			name: "get_radio_configuration",
			test: func(t *testing.T) {
				response, err := cisco.GetRadioCfg(client, ctx)
				if err != nil {
					t.Errorf("GetRadioCfg failed: %v", err)
					return
				}
				if response == nil {
					t.Error("GetRadioCfg returned nil response")
				}
				t.Logf("Radio configuration retrieved successfully")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}
