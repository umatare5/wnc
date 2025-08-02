package cisco

import (
	"testing"
	"time"

	wnc "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// TestNewClient tests client creation (Unit test)
func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantErr    bool
	}{
		{
			name:       "valid_secure_client",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   boolPtr(true),
			wantErr:    false,
		},
		{
			name:       "valid_insecure_client",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   boolPtr(false),
			wantErr:    false,
		},
		{
			name:       "nil_secure_option",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   nil,
			wantErr:    false,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-token",
			isSecure:   boolPtr(true),
			wantErr:    true,
		},
		{
			name:       "empty_apikey",
			controller: "192.168.1.1:443",
			apikey:     "",
			isSecure:   boolPtr(true),
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
				t.Error("expected non-nil client when no error")
			}
		})
	}
}

// TestNewClientWithTimeout tests client creation with timeout (Unit test)
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
			name:       "valid_timeout",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			timeout:    30 * time.Second,
			isSecure:   boolPtr(true),
			wantErr:    false,
		},
		{
			name:       "zero_timeout",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			timeout:    0,
			isSecure:   boolPtr(true),
			wantErr:    true,
		},
		{
			name:       "negative_timeout",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			timeout:    -5 * time.Second,
			isSecure:   boolPtr(true),
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
				t.Error("expected non-nil client when no error")
			}
		})
	}
}

// TestNewClientWithOptions tests client creation with options (Unit test)
func TestNewClientWithOptions(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		wantErr    bool
	}{
		{
			name:       "valid_options",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			wantErr:    false,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-token",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithOptions(tt.controller, tt.apikey)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("expected non-nil client when no error")
			}
		})
	}
}

// TestNewClientWithConfig tests client creation with config (Unit test)
func TestNewClientWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  wnc.Config
		wantErr bool
	}{
		{
			name: "valid_config",
			config: wnc.Config{
				Controller:  "192.168.1.1:443",
				AccessToken: "test-token",
			},
			wantErr: false,
		},
		{
			name: "empty_controller_config",
			config: wnc.Config{
				Controller:  "",
				AccessToken: "test-token",
			},
			wantErr: true,
		},
		{
			name: "insecure_config",
			config: wnc.Config{
				Controller:         "192.168.1.1:443",
				AccessToken:        "test-token",
				InsecureSkipVerify: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClientWithConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClientWithConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("expected non-nil client when no error")
			}
		})
	}
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}
