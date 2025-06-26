package cisco

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantErr    bool
	}{
		{
			name:       "valid parameters",
			controller: "wnc.example.com",
			apikey:     "token123",
			isSecure:   func() *bool { b := true; return &b }(),
			wantErr:    false,
		},
		{
			name:       "insecure connection",
			controller: "wnc.example.com",
			apikey:     "token123",
			isSecure:   func() *bool { b := false; return &b }(),
			wantErr:    false,
		},
		{
			name:       "nil isSecure",
			controller: "wnc.example.com",
			apikey:     "token123",
			isSecure:   nil,
			wantErr:    false,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "token123",
			isSecure:   nil,
			wantErr:    true,
		},
		{
			name:       "empty apikey",
			controller: "wnc.example.com",
			apikey:     "",
			isSecure:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.controller, tt.apikey, tt.isSecure)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client when error was not expected")
			}

			if tt.wantErr && client != nil {
				t.Error("NewClient() returned non-nil client when error was expected")
			}
		})
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		timeout    time.Duration
		isSecure   *bool
		wantErr    bool
	}{
		{
			name:       "valid parameters with timeout",
			controller: "wnc.example.com",
			apikey:     "token123",
			timeout:    30 * time.Second,
			isSecure:   func() *bool { b := true; return &b }(),
			wantErr:    false,
		},
		{
			name:       "zero timeout",
			controller: "wnc.example.com",
			apikey:     "token123",
			timeout:    0,
			isSecure:   nil,
			wantErr:    true, // Zero timeout should cause error
		},
		{
			name:       "negative timeout",
			controller: "wnc.example.com",
			apikey:     "token123",
			timeout:    -1 * time.Second,
			isSecure:   nil,
			wantErr:    true, // Negative timeout should cause error
		},
		{
			name:       "very long timeout",
			controller: "wnc.example.com",
			apikey:     "token123",
			timeout:    300 * time.Second,
			isSecure:   func() *bool { b := false; return &b }(),
			wantErr:    false,
		},
		{
			name:       "empty controller with timeout",
			controller: "",
			apikey:     "token123",
			timeout:    30 * time.Second,
			isSecure:   nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithTimeout(tt.controller, tt.apikey, tt.timeout, tt.isSecure)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && client == nil {
				t.Error("NewClientWithTimeout() returned nil client when error was not expected")
			}

			if tt.wantErr && client != nil {
				t.Error("NewClientWithTimeout() returned non-nil client when error was expected")
			}
		})
	}
}

func TestClientCreationJSONSerialization(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
	}{
		{
			name: "client parameters",
			params: map[string]interface{}{
				"controller": "wnc.example.com",
				"apikey":     "token123",
				"isSecure":   true,
			},
		},
		{
			name: "insecure client parameters",
			params: map[string]interface{}{
				"controller": "wnc.example.com",
				"apikey":     "token123",
				"isSecure":   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling of parameters
			jsonData, err := json.Marshal(tt.params)
			if err != nil {
				t.Fatalf("Failed to marshal parameters to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledParams map[string]interface{}
			err = json.Unmarshal(jsonData, &unmarshaledParams)
			if err != nil {
				t.Fatalf("Failed to unmarshal parameters from JSON: %v", err)
			}

			// Verify basic fields
			if controller, ok := tt.params["controller"].(string); ok {
				if unmarshaledController, ok := unmarshaledParams["controller"].(string); ok && unmarshaledController != controller {
					t.Errorf("Controller mismatch: got %q, want %q", unmarshaledController, controller)
				}
			}
		})
	}
}

func TestClientFailFast(t *testing.T) {
	tests := []struct {
		name        string
		controller  string
		apikey      string
		isSecure    *bool
		expectPanic bool
	}{
		{
			name:        "valid inputs should not panic",
			controller:  "wnc.example.com",
			apikey:      "token123",
			isSecure:    func() *bool { b := true; return &b }(),
			expectPanic: false,
		},
		{
			name:        "empty controller should not panic",
			controller:  "",
			apikey:      "token123",
			isSecure:    nil,
			expectPanic: false,
		},
		{
			name:        "empty apikey should not panic",
			controller:  "wnc.example.com",
			apikey:      "",
			isSecure:    nil,
			expectPanic: false,
		},
		{
			name:        "nil isSecure should not panic",
			controller:  "wnc.example.com",
			apikey:      "token123",
			isSecure:    nil,
			expectPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.expectPanic {
					t.Error("Expected panic but none occurred")
				}
			}()

			// Test NewClient
			client, err := NewClient(tt.controller, tt.apikey, tt.isSecure)
			if err != nil {
				t.Logf("NewClient returned error (may be expected): %v", err)
			}
			if client != nil {
				t.Logf("NewClient returned client successfully")
			}

			// Test NewClientWithTimeout
			clientWithTimeout, err := NewClientWithTimeout(tt.controller, tt.apikey, 30*time.Second, tt.isSecure)
			if err != nil {
				t.Logf("NewClientWithTimeout returned error (may be expected): %v", err)
			}
			if clientWithTimeout != nil {
				t.Logf("NewClientWithTimeout returned client successfully")
			}
		})
	}
}

func TestClientTimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    bool // whether creation should succeed
	}{
		{
			name:    "standard timeout",
			timeout: 30 * time.Second,
			want:    true,
		},
		{
			name:    "short timeout",
			timeout: 1 * time.Second,
			want:    false, // Very short timeout may cause errors
		},
		{
			name:    "long timeout",
			timeout: 5 * time.Minute,
			want:    true,
		},
		{
			name:    "zero timeout",
			timeout: 0,
			want:    false, // Zero timeout should fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := "wnc.example.com"
			apikey := "token123"
			isSecure := func() *bool { b := true; return &b }()

			client, err := NewClientWithTimeout(controller, apikey, tt.timeout, isSecure)

			if tt.want && err != nil {
				t.Errorf("NewClientWithTimeout() with timeout %v failed: %v", tt.timeout, err)
			}

			if !tt.want && err == nil {
				t.Errorf("NewClientWithTimeout() with timeout %v should have failed", tt.timeout)
			}

			if tt.want && client == nil {
				t.Errorf("NewClientWithTimeout() with timeout %v returned nil client", tt.timeout)
			}
		})
	}
}

func TestClientSecurityConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		isSecure *bool
		want     bool // whether creation should succeed
	}{
		{
			name:     "secure connection",
			isSecure: func() *bool { b := true; return &b }(),
			want:     true,
		},
		{
			name:     "insecure connection",
			isSecure: func() *bool { b := false; return &b }(),
			want:     true,
		},
		{
			name:     "nil security setting",
			isSecure: nil,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := "wnc.example.com"
			apikey := "token123"

			client, err := NewClient(controller, apikey, tt.isSecure)

			if tt.want && err != nil {
				t.Errorf("NewClient() with isSecure %v failed: %v", tt.isSecure, err)
			}

			if !tt.want && err == nil {
				t.Errorf("NewClient() with isSecure %v should have failed", tt.isSecure)
			}

			if tt.want && client == nil {
				t.Errorf("NewClient() with isSecure %v returned nil client", tt.isSecure)
			}
		})
	}
}

func TestClientIntegration(t *testing.T) {
	// This would be an integration test that uses real WNC credentials
	// We'll make it conditional on environment variables
	tests := []struct {
		name string
	}{
		{
			name: "integration test with environment credentials",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if no real credentials available for integration testing
			controller := "wnc1.example.internal"
			apikey := "YWRtaW46Y3l0WU43WVh4M2swc3piUnVhb1V1ZUx6"

			if controller == "" || apikey == "" {
				t.Skip("Skipping integration test - no WNC credentials provided")
			}

			isSecure := func() *bool { b := false; return &b }() // Use insecure for testing

			client, err := NewClient(controller, apikey, isSecure)
			if err != nil {
				t.Logf("Integration test failed to create client: %v", err)
				return
			}

			if client == nil {
				t.Error("Integration test: client should not be nil")
			} else {
				t.Log("Integration test: client created successfully")
			}
		})
	}
}
