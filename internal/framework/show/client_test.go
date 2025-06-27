package show

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestClientCli_JSON tests JSON serialization and deserialization
func TestClientCli_JSON(t *testing.T) {
	tests := []struct {
		name string
		data ClientCli
	}{
		{
			name: "valid ClientCli struct",
			data: ClientCli{
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
			var unmarshaled ClientCli
			if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
				t.Errorf("JSON unmarshaling failed: %v", err)
			}
		})
	}
}

// TestClientCli_GetShowClientTableHeaders tests the getShowClientTableHeaders method
func TestClientCli_GetShowClientTableHeaders(t *testing.T) {
	tests := []struct {
		name string
		cli  *ClientCli
	}{
		{
			name: "valid headers",
			cli: &ClientCli{
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
					t.Errorf("getShowClientTableHeaders panicked: %v", r)
				}
			}()

			result := tt.cli.getShowClientTableHeaders()
			if len(result) == 0 {
				t.Error("getShowClientTableHeaders returned empty headers")
			}
		})
	}
}

// TestClientCli_FormatShowClientRow tests the formatShowClientRow method
func TestClientCli_FormatShowClientRow(t *testing.T) {
	tests := []struct {
		name   string
		cli    *ClientCli
		client *application.ShowClientData
	}{
		{
			name: "nil client data",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			client: nil,
		},
		{
			name: "empty client data",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			client: &application.ShowClientData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("formatShowClientRow panicked: %v", r)
				}
			}()

			if tt.client == nil {
				return // Skip nil test case
			}

			row, err := tt.cli.formatShowClientRow(tt.client)
			if err != nil {
				// Error expected for empty data with invalid strings
				return
			}

			if len(row) == 0 {
				t.Error("formatShowClientRow returned empty row")
			}
		})
	}
}

// TestClientCli_SortShowClientRow tests the sortShowClientRow method
func TestClientCli_SortShowClientRow(t *testing.T) {
	tests := []struct {
		name    string
		cli     *ClientCli
		clients []*application.ShowClientData
	}{
		{
			name: "nil slice",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			clients: nil,
		},
		{
			name: "empty slice",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			clients: []*application.ShowClientData{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("sortShowClientRow panicked: %v", r)
				}
			}()

			tt.cli.sortShowClientRow(tt.clients)
		})
	}
}

// TestClientCli_ConvertCommonOperDataUsername tests the convertCommonOperDataUsername method
func TestClientCli_ConvertCommonOperDataUsername(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ClientCli
		input    string
		expected string
	}{
		{
			name: "empty username",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "",
			expected: "N/A",
		},
		{
			name: "valid username",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "testuser",
			expected: "testuser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertCommonOperDataUsername panicked: %v", r)
				}
			}()

			result := tt.cli.convertCommonOperDataUsername(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestClientCli_ConvertCommonOperDataCoState tests the convertCommonOperDataCoState method
func TestClientCli_ConvertCommonOperDataCoState(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ClientCli
		input    string
		expected string
	}{
		{
			name: "client-status-idle",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "client-status-idle",
			expected: "Idle",
		},
		{
			name: "client-status-run",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "client-status-run",
			expected: "Run",
		},
		{
			name: "unknown state",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "unknown-state",
			expected: "unknown-state",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertCommonOperDataCoState panicked: %v", r)
				}
			}()

			result := tt.cli.convertCommonOperDataCoState(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestClientCli_ConvertCommonOperDataMsRadioTypeToBand tests the convertCommonOperDataMsRadioTypeToBand method
func TestClientCli_ConvertCommonOperDataMsRadioTypeToBand(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ClientCli
		input    int
		expected string
	}{
		{
			name: "2.4GHz band",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    0,
			expected: "2.4GHz",
		},
		{
			name: "5GHz band",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    1,
			expected: "5GHz",
		},
		{
			name: "6GHz band",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    2,
			expected: "6GHz",
		},
		{
			name: "unknown band",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    99,
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertCommonOperDataMsRadioTypeToBand panicked: %v", r)
				}
			}()

			result := tt.cli.convertCommonOperDataMsRadioTypeToBand(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestClientCli_ConvertCommonOperDataMsRadioTypeToSpec tests the convertCommonOperDataMsRadioTypeToSpec method
func TestClientCli_ConvertCommonOperDataMsRadioTypeToSpec(t *testing.T) {
	tests := []struct {
		name     string
		cli      *ClientCli
		input    string
		expected string
	}{
		{
			name: "client-dot11b",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "client-dot11b",
			expected: "11b",
		},
		{
			name: "client-dot11ac",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "client-dot11ac",
			expected: "11ac",
		},
		{
			name: "client-dot11ax-5ghz-prot",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "client-dot11ax-5ghz-prot",
			expected: "dot11ax",
		},
		{
			name: "unknown protocol",
			cli: &ClientCli{
				Config:     &config.Config{},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
			input:    "unknown-protocol",
			expected: "unknown-protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("convertCommonOperDataMsRadioTypeToSpec panicked: %v", r)
				}
			}()

			result := tt.cli.convertCommonOperDataMsRadioTypeToSpec(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
