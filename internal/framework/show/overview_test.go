package show

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestOverviewCli_JSON tests JSON serialization and deserialization
func TestOverviewCli_JSON(t *testing.T) {
	tests := []struct {
		name string
		data OverviewCli
	}{
		{
			name: "valid OverviewCli struct",
			data: OverviewCli{
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
			var unmarshaled OverviewCli
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}

// TestOverviewCli_GetShowOverviewTableHeaders tests the getShowOverviewTableHeaders method
func TestOverviewCli_GetShowOverviewTableHeaders(t *testing.T) {
	tests := []struct {
		name string
		cli  *OverviewCli
	}{
		{
			name: "valid headers",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("getShowOverviewTableHeaders panicked: %v", r)
				}
			}()

			result := tt.cli.getShowOverviewTableHeaders()
			if len(result) == 0 {
				t.Error("getShowOverviewTableHeaders returned empty headers")
			}
		})
	}
}

// TestOverviewCli_FormatShowOverviewRow tests the formatShowOverviewRow method
func TestOverviewCli_FormatShowOverviewRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *OverviewCli
		data *application.ShowOverviewData
	}{
		{
			name: "nil overview data",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			data: nil,
		},
		{
			name: "empty overview data",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			data: &application.ShowOverviewData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("formatShowOverviewRow panicked: %v", r)
				}
			}()

			if tt.data == nil {
				return // Skip nil test case
			}

			row, err := tt.cli.formatShowOverviewRow(tt.data)
			if err != nil {
				// Error expected for empty data with access to slice index
				return
			}

			if len(row) == 0 {
				t.Error("formatShowOverviewRow returned empty row")
			}
		})
	}
}

// TestOverviewCli_SortShowOverviewRow tests the sortShowOverviewRow method
func TestOverviewCli_SortShowOverviewRow(t *testing.T) {
	tests := []struct {
		name string
		cli  *OverviewCli
		data []*application.ShowOverviewData
	}{
		{
			name: "nil slice",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			data: nil,
		},
		{
			name: "empty slice",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			data: []*application.ShowOverviewData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sortShowOverviewRow panicked: %v", r)
				}
			}()

			tt.cli.sortShowOverviewRow(tt.data)
		})
	}
}

// TestOverviewCli_ConvertUtilizationsToIndicator tests the convertUtilizationsToIndicator method
func TestOverviewCli_ConvertUtilizationsToIndicator(t *testing.T) {
	tests := []struct {
		name     string
		cli      *OverviewCli
		rx       int
		tx       int
		noise    int
		expected string
	}{
		{
			name: "zero utilization",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			rx:       0,
			tx:       0,
			noise:    0,
			expected: "[          ] 0%",
		},
		{
			name: "normal utilization",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			rx:       10,
			tx:       10,
			noise:    30,
			expected: "[#####     ] 50%",
		},
		{
			name: "over 100% utilization",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			rx:       60,
			tx:       60,
			noise:    30,
			expected: "[##########] 100%",
		},
		{
			name: "negative utilization",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			rx:       -10,
			tx:       -10,
			noise:    -10,
			expected: "[          ] 0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertUtilizationsToIndicator panicked: %v", r)
				}
			}()

			result := tt.cli.convertUtilizationsToIndicator(tt.rx, tt.tx, tt.noise)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestOverviewCli_ConvertRadioOperDataOperState tests the convertRadioOperDataOperState method
func TestOverviewCli_ConvertRadioOperDataOperState(t *testing.T) {
	tests := []struct {
		name     string
		cli      *OverviewCli
		input    string
		expected string
	}{
		{
			name: "radio-up state",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "radio-up",
			expected: "  ✅️",
		},
		{
			name: "radio-down state",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "radio-down",
			expected: "  ❌️",
		},
		{
			name: "unknown state",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "unknown",
			expected: "  ❌️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertRadioOperDataOperState panicked: %v", r)
				}
			}()

			result := tt.cli.convertRadioOperDataOperState(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestOverviewCli_ConvertRrmMeasurementLoadStations tests the convertRrmMeasurementLoadStations method
func TestOverviewCli_ConvertRrmMeasurementLoadStations(t *testing.T) {
	tests := []struct {
		name     string
		cli      *OverviewCli
		input    int
		expected string
	}{
		{
			name: "zero clients",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    0,
			expected: "0 clients",
		},
		{
			name: "multiple clients",
			cli: &OverviewCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    25,
			expected: "25 clients",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertRrmMeasurementLoadStations panicked: %v", r)
				}
			}()

			result := tt.cli.convertRrmMeasurementLoadStations(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
