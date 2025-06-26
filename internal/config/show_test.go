package config

import (
	"encoding/json"
	"testing"
)

func TestShowCmdConfigJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		config ShowCmdConfig
	}{
		{
			name:   "empty show config",
			config: ShowCmdConfig{},
		},
		{
			name: "full show config",
			config: ShowCmdConfig{
				Controllers: []Controller{
					{Hostname: "wnc1.example.com", AccessToken: "token1"},
					{Hostname: "wnc2.example.com", AccessToken: "token2"},
				},
				AllowInsecureAccess: true,
				PrintFormat:         PrintFormatJSON,
				Timeout:             60,
				APName:              "test-ap",
				Radio:               "0",
				SSID:                "test-ssid",
				SortBy:              "name",
				SortOrder:           OrderByAscending,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.config)
			if err != nil {
				t.Fatalf("Failed to marshal ShowCmdConfig to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledConfig ShowCmdConfig
			err = json.Unmarshal(jsonData, &unmarshaledConfig)
			if err != nil {
				t.Fatalf("Failed to unmarshal ShowCmdConfig from JSON: %v", err)
			}

			// Verify key fields
			if len(unmarshaledConfig.Controllers) != len(tt.config.Controllers) {
				t.Errorf("Controller count mismatch: got %d, want %d",
					len(unmarshaledConfig.Controllers), len(tt.config.Controllers))
			}

			if unmarshaledConfig.PrintFormat != tt.config.PrintFormat {
				t.Errorf("PrintFormat mismatch: got %q, want %q",
					unmarshaledConfig.PrintFormat, tt.config.PrintFormat)
			}

			if unmarshaledConfig.AllowInsecureAccess != tt.config.AllowInsecureAccess {
				t.Errorf("AllowInsecureAccess mismatch: got %t, want %t",
					unmarshaledConfig.AllowInsecureAccess, tt.config.AllowInsecureAccess)
			}
		})
	}
}

func TestControllerJSONSerialization(t *testing.T) {
	tests := []struct {
		name       string
		controller Controller
	}{
		{
			name:       "empty controller",
			controller: Controller{},
		},
		{
			name: "full controller",
			controller: Controller{
				Hostname:    "wnc.example.com",
				AccessToken: "test-token-123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.controller)
			if err != nil {
				t.Fatalf("Failed to marshal Controller to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledController Controller
			err = json.Unmarshal(jsonData, &unmarshaledController)
			if err != nil {
				t.Fatalf("Failed to unmarshal Controller from JSON: %v", err)
			}

			// Verify fields
			if unmarshaledController.Hostname != tt.controller.Hostname {
				t.Errorf("Hostname mismatch: got %q, want %q",
					unmarshaledController.Hostname, tt.controller.Hostname)
			}

			if unmarshaledController.AccessToken != tt.controller.AccessToken {
				t.Errorf("AccessToken mismatch: got %q, want %q",
					unmarshaledController.AccessToken, tt.controller.AccessToken)
			}
		})
	}
}

func TestValidateControllersFormat(t *testing.T) {
	c := &Config{}

	tests := []struct {
		name      string
		input     string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid single controller",
			input:     "wnc.example.com:token123",
			wantError: false,
		},
		{
			name:      "valid multiple controllers",
			input:     "wnc1.example.com:token1,wnc2.example.com:token2",
			wantError: false,
		},
		{
			name:      "valid controller with https",
			input:     "https://wnc.example.com:443:token123",
			wantError: false,
		},
		{
			name:      "valid controller with http",
			input:     "http://wnc.example.com:80:token123",
			wantError: false,
		},
		{
			name:      "empty hostname",
			input:     ":token123",
			wantError: true,
			errorMsg:  "invalid controllers format: hostname is empty.",
		},
		{
			name:      "empty token",
			input:     "wnc.example.com:",
			wantError: true,
			errorMsg:  "invalid controllers format: access token is empty.",
		},
		{
			name:      "missing colon",
			input:     "wnc.example.com",
			wantError: true,
			errorMsg:  "invalid controllers format: controllers does not contain ':'.",
		},
		{
			name:      "empty input with colon",
			input:     ":",
			wantError: true,
			errorMsg:  "invalid controllers format: hostname is empty.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validateControllersFormat(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("validateControllersFormat() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("validateControllersFormat() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateControllersFormat() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestParseControllerPair(t *testing.T) {
	c := &Config{}

	tests := []struct {
		name            string
		input           string
		wantHostname    string
		wantAccessToken string
		wantError       bool
		errorMsg        string
	}{
		{
			name:            "simple hostname:token",
			input:           "wnc.example.com:token123",
			wantHostname:    "wnc.example.com",
			wantAccessToken: "token123",
			wantError:       false,
		},
		{
			name:            "https URL with port",
			input:           "https://wnc.example.com:443:token123",
			wantHostname:    "wnc.example.com:443",
			wantAccessToken: "token123",
			wantError:       false,
		},
		{
			name:            "http URL with port",
			input:           "http://wnc.example.com:80:token123",
			wantHostname:    "wnc.example.com:80",
			wantAccessToken: "token123",
			wantError:       false,
		},
		{
			name:            "hostname with multiple colons",
			input:           "wnc.example.com:8443:token:with:colons",
			wantHostname:    "wnc.example.com:8443:token:with",
			wantAccessToken: "colons",
			wantError:       false,
		},
		{
			name:      "no colon",
			input:     "wnc.example.com",
			wantError: true,
			errorMsg:  "invalid controllers format: controllers does not contain ':'.",
		},
		{
			name:            "whitespace handling",
			input:           "  wnc.example.com  :  token123  ",
			wantHostname:    "wnc.example.com",
			wantAccessToken: "token123",
			wantError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostname, accessToken, err := c.parseControllerPair(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("parseControllerPair() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("parseControllerPair() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("parseControllerPair() unexpected error = %v", err)
					return
				}
				if hostname != tt.wantHostname {
					t.Errorf("parseControllerPair() hostname = %q, want %q", hostname, tt.wantHostname)
				}
				if accessToken != tt.wantAccessToken {
					t.Errorf("parseControllerPair() accessToken = %q, want %q", accessToken, tt.wantAccessToken)
				}
			}
		})
	}
}

func TestValidatePrintFormat(t *testing.T) {
	c := &Config{}

	tests := []struct {
		name      string
		format    string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid json format",
			format:    PrintFormatJSON,
			wantError: false,
		},
		{
			name:      "valid table format",
			format:    PrintFormatTable,
			wantError: false,
		},
		{
			name:      "invalid format",
			format:    "xml",
			wantError: true,
			errorMsg:  `invalid format: must be "json" or "table"`,
		},
		{
			name:      "empty format",
			format:    "",
			wantError: true,
			errorMsg:  `invalid format: must be "json" or "table"`,
		},
		{
			name:      "case sensitive",
			format:    "JSON",
			wantError: true,
			errorMsg:  `invalid format: must be "json" or "table"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.validatePrintFormat(tt.format)

			if tt.wantError {
				if err == nil {
					t.Errorf("validatePrintFormat() expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("validatePrintFormat() error = %q, want %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validatePrintFormat() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestParseControllers(t *testing.T) {
	c := &Config{}

	tests := []struct {
		name  string
		input string
		want  []Controller
	}{
		{
			name:  "single controller",
			input: "wnc.example.com:token123",
			want: []Controller{
				{Hostname: "wnc.example.com", AccessToken: "token123"},
			},
		},
		{
			name:  "multiple controllers",
			input: "wnc1.example.com:token1,wnc2.example.com:token2",
			want: []Controller{
				{Hostname: "wnc1.example.com", AccessToken: "token1"},
				{Hostname: "wnc2.example.com", AccessToken: "token2"},
			},
		},
		{
			name:  "controllers with https",
			input: "https://wnc1.example.com:443:token1,https://wnc2.example.com:8443:token2",
			want: []Controller{
				{Hostname: "wnc1.example.com:443", AccessToken: "token1"},
				{Hostname: "wnc2.example.com:8443", AccessToken: "token2"},
			},
		},
		{
			name:  "mixed formats",
			input: "wnc1.example.com:token1,https://wnc2.example.com:8443:token2",
			want: []Controller{
				{Hostname: "wnc1.example.com", AccessToken: "token1"},
				{Hostname: "wnc2.example.com:8443", AccessToken: "token2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.parseControllers(tt.input)

			if len(got) != len(tt.want) {
				t.Errorf("parseControllers() length = %d, want %d", len(got), len(tt.want))
				return
			}

			for i, controller := range got {
				if controller.Hostname != tt.want[i].Hostname {
					t.Errorf("parseControllers()[%d].Hostname = %q, want %q",
						i, controller.Hostname, tt.want[i].Hostname)
				}
				if controller.AccessToken != tt.want[i].AccessToken {
					t.Errorf("parseControllers()[%d].AccessToken = %q, want %q",
						i, controller.AccessToken, tt.want[i].AccessToken)
				}
			}
		})
	}
}

func TestShowCmdConfigHeaderConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"OverviewHeaderApMac", OverviewHeaderApMac, "APMac"},
		{"OverviewHeaderApRadioID", OverviewHeaderApRadioID, "Radio"},
		{"OverviewHeaderApOperStatus", OverviewHeaderApOperStatus, "Status"},
		{"OverviewHeaderChannelNumber", OverviewHeaderChannelNumber, "Channel"},
		{"OverviewHeaderChannelUtilization", OverviewHeaderChannelUtilization, "ChannelUtilization"},
		{"OverviewHeaderClientCount", OverviewHeaderClientCount, "ClientCount"},
		{"OverviewHeaderRFTagName", OverviewHeaderRFTagName, "RFTagName"},
		{"OverviewHeaderTxPower", OverviewHeaderTxPower, "TxPower"},
		{"ShowClientHeaderBand", ShowClientHeaderBand, "Band"},
		{"ShowClientHeaderHostname", ShowClientHeaderHostname, "Hostname"},
		{"ShowClientHeaderIP", ShowClientHeaderIP, "IPAddress"},
		{"ShowClientHeaderMacAddress", ShowClientHeaderMacAddress, "MACAddress"},
		{"ShowClientHeaderProtocol", ShowClientHeaderProtocol, "Protocol"},
		{"ShowClientHeaderRSSI", ShowClientHeaderRSSI, "RSSI"},
		{"ShowClientHeaderRxTraffic", ShowClientHeaderRxTraffic, "RxTraffic"},
		{"ShowClientHeaderSNR", ShowClientHeaderSNR, "SNR"},
		{"ShowClientHeaderSSID", ShowClientHeaderSSID, "SSID"},
		{"ShowClientHeaderState", ShowClientHeaderState, "State"},
		{"ShowClientHeaderStream", ShowClientHeaderStream, "Stream"},
		{"ShowClientHeaderThroughput", ShowClientHeaderThroughput, "Throughput"},
		{"ShowClientHeaderTxTraffic", ShowClientHeaderTxTraffic, "TxTraffic"},
		{"ShowClientHeaderUsername", ShowClientHeaderUsername, "Username"},
		{"ShowCommonHeaderApName", ShowCommonHeaderApName, "APName"},
		{"ShowCommonHeaderController", ShowCommonHeaderController, "Controller"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Header constant %s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}
