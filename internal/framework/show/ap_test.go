package cli

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestApCli_JSON tests JSON serialization and deserialization
func TestApCli_JSON(t *testing.T) {
	tests := []struct {
		name string
		data ApCli
	}{
		{
			name: "valid ApCli struct",
			data: ApCli{
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
			var unmarshaled ApCli
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}

// TestApCli_GetShowApTableHeaders tests the getShowApTableHeaders method
func TestApCli_GetShowApTableHeaders(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ApCli
		expected []string
	}{
		{
			name: "valid headers",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			expected: []string{
				"AP Name", "Slots", "Model", "Serial", "Ethernet MAC", "Radio MAC",
				"Country Code", "Domain", "IP Address", "OS Version",
				"State", "LLDP Neighbor", "Power Type", "Power Mode", "Controller",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("getShowApTableHeaders panicked: %v", r)
				}
			}()

			result := tt.cli.getShowApTableHeaders()
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

// TestApCli_FormatShowApRow tests the formatShowApRow method
func TestApCli_FormatShowApRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *ApCli
		ap   *application.ShowApData
	}{
		{
			name: "nil ap data",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			ap: nil,
		},
		{
			name: "empty ap data",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			ap: &application.ShowApData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("formatShowApRow panicked: %v", r)
				}
			}()

			if tt.ap == nil {
				return // Skip nil test case
			}

			row, err := tt.cli.formatShowApRow(tt.ap)
			if err != nil {
				t.Errorf("formatShowApRow returned error: %v", err)
				return
			}

			if len(row) == 0 {
				t.Error("formatShowApRow returned empty row")
			}
		})
	}
}

// TestApCli_SortShowClientRow tests the sortShowClientRow method
func TestApCli_SortShowClientRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *ApCli
		aps  []*application.ShowApData
	}{
		{
			name: "nil slice",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			aps: nil,
		},
		{
			name: "empty slice",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			aps: []*application.ShowApData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sortShowClientRow panicked: %v", r)
				}
			}()

			tt.cli.sortShowClientRow(tt.aps)
		})
	}
}

// TestApCli_ConvertApOperDataApPowPowerMode tests the convertApOperDataApPowPowerMode method
func TestApCli_ConvertApOperDataApPowPowerMode(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ApCli
		input    string
		expected string
	}{
		{
			name: "dot11-default-low-pwr",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "dot11-default-low-pwr",
			expected: "Default Low",
		},
		{
			name: "dot11-set-high-pwr",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "dot11-set-high-pwr",
			expected: "High",
		},
		{
			name: "unknown value",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "unknown-value",
			expected: "unknown-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertApOperDataApPowPowerMode panicked: %v", r)
				}
			}()

			result := tt.cli.convertApOperDataApPowPowerMode(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestApCli_ConvertApOperDataApPowPowerType tests the convertApOperDataApPowPowerType method
func TestApCli_ConvertApOperDataApPowPowerType(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ApCli
		input    string
		expected string
	}{
		{
			name: "pwr-src-poe-plus",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "pwr-src-poe-plus",
			expected: "Advanced PoE",
		},
		{
			name: "pwr-src-unknown",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "pwr-src-unknown",
			expected: "Unknown",
		},
		{
			name: "unknown value",
			cli: &ApCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "unknown-value",
			expected: "unknown-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertApOperDataApPowPowerType panicked: %v", r)
				}
			}()

			result := tt.cli.convertApOperDataApPowPowerType(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
