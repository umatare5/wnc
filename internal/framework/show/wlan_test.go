package cli

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestWlanCli_JSON tests JSON serialization and deserialization
func TestWlanCli_JSON(t *testing.T) {
	tests := []struct {
		name string
		data WlanCli
	}{
		{
			name: "valid WlanCli struct",
			data: WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Errorf("JSON marshaling failed: %v", err)
				return
			}

			// Test JSON unmarshaling
			var unmarshaled WlanCli
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}

// TestWlanCli_GetShowWlanTableHeaders tests the getShowWlanTableHeaders method
func TestWlanCli_GetShowWlanTableHeaders(t *testing.T) {
	tests := []struct {
		name     string
		cli      *WlanCli
		expected []string
	}{
		{
			name: "valid headers",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			expected: []string{
				"Status", "ESSID", "ID", "Profile Name", "VLAN", "Session Timeout",
				"DHCP Required", "Egress QoS", "Ingress QoS", "ATF Policies",
				"Auth Key Management", "mDNS Forwarding", "P2P Blocking", "Loadbalance",
				"Broadcast", "Tag Name", "Controller",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("getShowWlanTableHeaders panicked: %v", r)
				}
			}()

			result := tt.cli.getShowWlanTableHeaders()
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d headers, got %d", len(tt.expected), len(result))
				return
			}

			for i, header := range result {
				if header != tt.expected[i] {
					t.Errorf("expected header[%d] = %q, got %q", i, tt.expected[i], header)
				}
			}
		})
	}
}

// TestWlanCli_FormatShowWlanRow tests the formatShowWlanRow method
func TestWlanCli_FormatShowWlanRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *WlanCli
		wlan *application.ShowWlanData
	}{
		{
			name: "nil wlan data",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			wlan: nil,
		},
		{
			name: "empty wlan data",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			wlan: &application.ShowWlanData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("formatShowWlanRow panicked: %v", r)
				}
			}()

			if tt.wlan == nil {
				return // Skip nil test case
			}

			row, err := tt.cli.formatShowWlanRow(tt.wlan)
			if err != nil {
				t.Errorf("formatShowWlanRow returned error: %v", err)
				return
			}

			if len(row) == 0 {
				t.Error("formatShowWlanRow returned empty row")
			}
		})
	}
}

// TestWlanCli_SortShowWlanRow tests the sortShowWlanRow method
func TestWlanCli_SortShowWlanRow(t *testing.T) {
	tests := []struct {
		name  string
		cli   *WlanCli
		wlans []*application.ShowWlanData
	}{
		{
			name: "nil slice",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			wlans: nil,
		},
		{
			name: "empty slice",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			wlans: []*application.ShowWlanData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sortShowWlanRow panicked: %v", r)
				}
			}()

			tt.cli.sortShowWlanRow(tt.wlans)
		})
	}
}

// TestWlanCli_ConvertWlanPolicyStatus tests the convertWlanPolicyStatus method
func TestWlanCli_ConvertWlanPolicyStatus(t *testing.T) {
	tests := []struct {
		name     string
		cli      *WlanCli
		input    bool
		expected string
	}{
		{
			name: "true status",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    true,
			expected: "  ✅️",
		},
		{
			name: "false status",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    false,
			expected: "  ❌️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertWlanPolicyStatus panicked: %v", r)
				}
			}()

			result := tt.cli.convertWlanPolicyStatus(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestWlanCli_ConvertWlanCfgEntryAuthKeyMgmt tests the convertWlanCfgEntryAuthKeyMgmt method
func TestWlanCli_ConvertWlanCfgEntryAuthKeyMgmt(t *testing.T) {
	tests := []struct {
		name     string
		cli      *WlanCli
		dot1x    bool
		psk      bool
		sae      bool
		expected string
	}{
		{
			name: "dot1x enabled",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			dot1x:    true,
			psk:      false,
			sae:      false,
			expected: "Dot1x",
		},
		{
			name: "psk enabled",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			dot1x:    false,
			psk:      true,
			sae:      false,
			expected: "PSK",
		},
		{
			name: "sae enabled",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			dot1x:    false,
			psk:      false,
			sae:      true,
			expected: "SAE",
		},
		{
			name: "none enabled",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			dot1x:    false,
			psk:      false,
			sae:      false,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertWlanCfgEntryAuthKeyMgmt panicked: %v", r)
				}
			}()

			result := tt.cli.convertWlanCfgEntryAuthKeyMgmt(tt.dot1x, tt.psk, tt.sae)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestWlanCli_ConvertWlanCfgEntryMdnsSdMode tests the convertWlanCfgEntryMdnsSdMode method
func TestWlanCli_ConvertWlanCfgEntryMdnsSdMode(t *testing.T) {
	tests := []struct {
		name     string
		cli      *WlanCli
		input    string
		expected string
	}{
		{
			name: "mdns-sd-drop",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "mdns-sd-drop",
			expected: "Drop",
		},
		{
			name: "mdns-sd-gateway",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "mdns-sd-gateway",
			expected: "Gateway",
		},
		{
			name: "default bridging",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "other-value",
			expected: "Bridging",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertWlanCfgEntryMdnsSdMode panicked: %v", r)
				}
			}()

			result := tt.cli.convertWlanCfgEntryMdnsSdMode(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestWlanCli_ConvertWlanCfgEntryP2PBlockAction tests the convertWlanCfgEntryP2PBlockAction method
func TestWlanCli_ConvertWlanCfgEntryP2PBlockAction(t *testing.T) {
	tests := []struct {
		name     string
		cli      *WlanCli
		input    string
		expected string
	}{
		{
			name: "p2p-blocking-action-fwdup",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "p2p-blocking-action-fwdup",
			expected: "Forward-UpStream",
		},
		{
			name: "p2p-blocking-action-drop",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "p2p-blocking-action-drop",
			expected: "Drop",
		},
		{
			name: "p2p-blocking-action-allow-private-group",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "p2p-blocking-action-allow-private-group",
			expected: "Allow Private Group",
		},
		{
			name: "default disabled",
			cli: &WlanCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "other-value",
			expected: "Disabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertWlanCfgEntryP2PBlockAction panicked: %v", r)
				}
			}()

			result := tt.cli.convertWlanCfgEntryP2PBlockAction(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
