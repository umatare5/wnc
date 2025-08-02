package infrastructure

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestNew tests infrastructure layer initialization (Unit test)
func TestNew(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "create_infrastructure_layer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}

			repo := New(cfg)
			if repo.Config != cfg {
				t.Error("New() Config not set correctly")
			}
		})
	}
}

// TestInvokeClientRepository tests client repository invocation (Unit test)
func TestInvokeClientRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_client_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			clientRepo := repo.InvokeClientRepository()
			if clientRepo == nil {
				t.Error("InvokeClientRepository returned nil")
			}
			if clientRepo.Config != cfg {
				t.Error("ClientRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeApRepository tests AP repository invocation (Unit test)
func TestInvokeApRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_ap_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			apRepo := repo.InvokeApRepository()
			if apRepo == nil {
				t.Error("InvokeApRepository returned nil")
			}
			if apRepo.Config != cfg {
				t.Error("ApRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeWlanRepository tests WLAN repository invocation (Unit test)
func TestInvokeWlanRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_wlan_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			wlanRepo := repo.InvokeWlanRepository()
			if wlanRepo == nil {
				t.Error("InvokeWlanRepository returned nil")
			}
			if wlanRepo.Config != cfg {
				t.Error("WlanRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeRadioRepository tests radio repository invocation (Unit test)
func TestInvokeRadioRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_radio_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			radioRepo := repo.InvokeRadioRepository()
			if radioRepo == nil {
				t.Error("InvokeRadioRepository returned nil")
			}
			if radioRepo.Config != cfg {
				t.Error("RadioRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeRrmRepository tests RRM repository invocation (Unit test)
func TestInvokeRrmRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_rrm_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			rrmRepo := repo.InvokeRrmRepository()
			if rrmRepo == nil {
				t.Error("InvokeRrmRepository returned nil")
			}
			if rrmRepo.Config != cfg {
				t.Error("RrmRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeRfRepository tests RF repository invocation (Unit test)
func TestInvokeRfRepository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_rf_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			rfRepo := repo.InvokeRfRepository()
			if rfRepo == nil {
				t.Error("InvokeRfRepository returned nil")
			}
			if rfRepo.Config != cfg {
				t.Error("RfRepository Config not set correctly")
			}
		})
	}
}

// TestInvokeDot11Repository tests 802.11 repository invocation (Unit test)
func TestInvokeDot11Repository(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_dot11_repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)

			dot11Repo := repo.InvokeDot11Repository()
			if dot11Repo == nil {
				t.Error("InvokeDot11Repository returned nil")
			}
			if dot11Repo.Config != cfg {
				t.Error("Dot11Repository Config not set correctly")
			}
		})
	}
}

// TestGetClientOper tests GetClientOper method (Unit test)
func TestGetClientOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			clientRepo := repo.InvokeClientRepository()

			result := clientRepo.GetClientOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetClientOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetClientOper() expected non-nil result")
			}
		})
	}
}

// TestGetClientGlobalOper tests GetClientGlobalOper method (Unit test)
func TestGetClientGlobalOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			clientRepo := repo.InvokeClientRepository()

			result := clientRepo.GetClientGlobalOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetClientGlobalOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetClientGlobalOper() expected non-nil result")
			}
		})
	}
}

// TestGetWlanCfg tests GetWlanCfg method (Unit test)
func TestGetWlanCfg(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			wlanRepo := repo.InvokeWlanRepository()

			result := wlanRepo.GetWlanCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetWlanCfg() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetWlanCfg() expected non-nil result")
			}
		})
	}
}

// TestGetApOper tests GetApOper method (Unit test)
func TestGetApOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApOper() expected non-nil result")
			}
		})
	}
}

// TestGetApCapwapData tests GetApCapwapData method (Unit test)
func TestGetApCapwapData(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApCapwapData(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApCapwapData() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApCapwapData() expected non-nil result")
			}
		})
	}
}

// TestGetApLldpNeigh tests GetApLldpNeigh method (Unit test)
func TestGetApLldpNeigh(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApLldpNeigh(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApLldpNeigh() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApLldpNeigh() expected non-nil result")
			}
		})
	}
}

// TestGetApRadioOperData tests GetApRadioOperData method (Unit test)
func TestGetApRadioOperData(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApRadioOperData(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApRadioOperData() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApRadioOperData() expected non-nil result")
			}
		})
	}
}

// TestGetApOperData tests GetApOperData method (Unit test)
func TestGetApOperData(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApOperData(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApOperData() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApOperData() expected non-nil result")
			}
		})
	}
}

// TestGetApGlobalOper tests GetApGlobalOper method (Unit test)
func TestGetApGlobalOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApGlobalOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApGlobalOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApGlobalOper() expected non-nil result")
			}
		})
	}
}

// TestGetApCfg tests GetApCfg method (Unit test)
func TestGetApCfg(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid_controller",
			controller: "invalid://controller",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_controller",
			controller: "",
			apikey:     "test-apikey",
			isSecure:   nil,
			wantNil:    true,
		},
		{
			name:       "empty_apikey",
			controller: "https://test-controller.example.com",
			apikey:     "",
			isSecure:   nil,
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 30,
				},
			}
			repo := New(cfg)
			apRepo := repo.InvokeApRepository()

			result := apRepo.GetApCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetApCfg() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetApCfg() expected non-nil result")
			}
		})
	}
}

// TestGetRadioCfg tests GetRadioCfg method (Unit test)
func TestGetRadioCfg(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			radioRepo := repo.InvokeRadioRepository()

			result := radioRepo.GetRadioCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRadioCfg() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRadioCfg() expected non-nil result")
			}
		})
	}
}

// TestGetRfTags tests GetRfTags method (Unit test)
func TestGetRfTags(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			rfRepo := repo.InvokeRfRepository()

			result := rfRepo.GetRfTags(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRfTags() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRfTags() expected non-nil result")
			}
		})
	}
}

// TestGetRrmOper tests GetRrmOper method (Unit test)
func TestGetRrmOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			rrmRepo := repo.InvokeRrmRepository()

			result := rrmRepo.GetRrmOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRrmOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRrmOper() expected non-nil result")
			}
		})
	}
}

// TestGetRrmMeasurement tests GetRrmMeasurement method (Unit test)
func TestGetRrmMeasurement(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			rrmRepo := repo.InvokeRrmRepository()

			result := rrmRepo.GetRrmMeasurement(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRrmMeasurement() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRrmMeasurement() expected non-nil result")
			}
		})
	}
}

// TestGetRrmGlobalOper tests GetRrmGlobalOper method (Unit test)
func TestGetRrmGlobalOper(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			rrmRepo := repo.InvokeRrmRepository()

			result := rrmRepo.GetRrmGlobalOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRrmGlobalOper() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRrmGlobalOper() expected non-nil result")
			}
		})
	}
}

// TestGetRrmCfg tests GetRrmCfg method (Unit test)
func TestGetRrmCfg(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			rrmRepo := repo.InvokeRrmRepository()

			result := rrmRepo.GetRrmCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetRrmCfg() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetRrmCfg() expected non-nil result")
			}
		})
	}
}

// TestGetDot11Cfg tests GetDot11Cfg method (Unit test)
func TestGetDot11Cfg(t *testing.T) {
	tests := []struct {
		name       string
		controller string
		apikey     string
		isSecure   *bool
		wantNil    bool
	}{
		{
			name:       "invalid controller URL",
			controller: "invalid://controller",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty controller",
			controller: "",
			apikey:     "test-api-key",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
		{
			name:       "empty api key",
			controller: "https://controller.example.com",
			apikey:     "",
			isSecure:   &[]bool{true}[0],
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := New(cfg)
			dot11Repo := repo.InvokeDot11Repository()

			result := dot11Repo.GetDot11Cfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.wantNil && result != nil {
				t.Errorf("GetDot11Cfg() expected nil, got %v", result)
			}
			if !tt.wantNil && result == nil {
				t.Error("GetDot11Cfg() expected non-nil result")
			}
		})
	}
}
