package infrastructure

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestApRepository_GetApOper tests the GetApOper method (Unit test)
func TestApRepository_GetApOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		timeout    int
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "",
			apikey:     "test-token",
			isSecure:   boolPtr(true),
			timeout:    30,
			wantNil:    true,
		},
		{
			name:       "invalid_apikey",
			controller: "192.168.1.1:443",
			apikey:     "",
			isSecure:   boolPtr(true),
			timeout:    30,
			wantNil:    true,
		},
		{
			name:       "valid_insecure_request",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   boolPtr(false),
			timeout:    30,
			wantNil:    true, // Still nil due to network error, but tests the path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ApRepository{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						Timeout: tt.timeout,
					},
				},
			}

			result := repo.GetApOper(tt.controller, tt.apikey, tt.isSecure)
			if tt.wantNil && result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

// TestApRepository_GetApCapwapData tests the GetApCapwapData method (Unit test)
func TestApRepository_GetApCapwapData(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		timeout    int
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "",
			apikey:     "test-token",
			isSecure:   boolPtr(true),
			timeout:    30,
			wantNil:    true,
		},
		{
			name:       "valid_request",
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   boolPtr(true),
			timeout:    30,
			wantNil:    true, // Still nil due to network error, but tests the path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ApRepository{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						Timeout: tt.timeout,
					},
				},
			}

			result := repo.GetApCapwapData(tt.controller, tt.apikey, tt.isSecure)
			if tt.wantNil && result != nil {
				t.Errorf("expected nil result, got %v", result)
			}
		})
	}
}

// TestApRepository_GetApLldpNeigh tests the GetApLldpNeigh method (Unit test)
func TestApRepository_GetApLldpNeigh(t *testing.T) {
	repo := &ApRepository{
		Config: &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Timeout: 1,
			},
		},
	}

	result := repo.GetApLldpNeigh("", "test-token", boolPtr(true))
	if result != nil {
		t.Errorf("expected nil result for empty controller, got %v", result)
	}
}

// TestApRepository_GetApRadioOperData tests the GetApRadioOperData method (Unit test)
func TestApRepository_GetApRadioOperData(t *testing.T) {
	repo := &ApRepository{
		Config: &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Timeout: 1,
			},
		},
	}

	result := repo.GetApRadioOperData("", "test-token", boolPtr(true))
	if result != nil {
		t.Errorf("expected nil result for empty controller, got %v", result)
	}
}

// TestApRepository_GetApOperData tests the GetApOperData method (Unit test)
func TestApRepository_GetApOperData(t *testing.T) {
	repo := &ApRepository{
		Config: &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Timeout: 1,
			},
		},
	}

	result := repo.GetApOperData("", "test-token", boolPtr(true))
	if result != nil {
		t.Errorf("expected nil result for empty controller, got %v", result)
	}
}

// TestApRepository_GetApGlobalOper tests the GetApGlobalOper method (Unit test)
func TestApRepository_GetApGlobalOper(t *testing.T) {
	repo := &ApRepository{
		Config: &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Timeout: 1,
			},
		},
	}

	result := repo.GetApGlobalOper("", "test-token", boolPtr(true))
	if result != nil {
		t.Errorf("expected nil result for empty controller, got %v", result)
	}
}

// TestApRepository_GetApCfg tests the GetApCfg method (Unit test)
func TestApRepository_GetApCfg(t *testing.T) {
	repo := &ApRepository{
		Config: &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				Timeout: 1,
			},
		},
	}

	result := repo.GetApCfg("", "test-token", boolPtr(true))
	if result != nil {
		t.Errorf("expected nil result for empty controller, got %v", result)
	}
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}
