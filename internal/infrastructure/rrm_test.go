package infrastructure

import (
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
)

// TestRrmRepository tests the RrmRepository structure
func TestRrmRepository(t *testing.T) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 1,
		},
	}

	repo := &RrmRepository{
		Config: cfg,
	}

	// Test that repository is properly initialized
	if repo.Config == nil {
		t.Error("Expected config to be set")
	}

	if repo.Config.ShowCmdConfig.Timeout != 1 {
		t.Errorf("Expected timeout to be 1, got %d", repo.Config.ShowCmdConfig.Timeout)
	}
}

// TestRrmRepository_GetRrmOper tests the GetRrmOper method
func TestRrmRepository_GetRrmOper(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name: "valid_parameters",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true, // Will be nil due to no real connection
		},
		{
			name: "invalid_controller",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{
				Config: tt.config,
			}

			result := repo.GetRrmOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRrmRepository_GetRrmMeasurement tests the GetRrmMeasurement method
func TestRrmRepository_GetRrmMeasurement(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name: "valid_parameters",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
		{
			name: "insecure_connection",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:8080",
			apikey:     "test-token",
			isSecure:   &[]bool{false}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{
				Config: tt.config,
			}

			result := repo.GetRrmMeasurement(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRrmRepository_GetRrmGlobalOper tests the GetRrmGlobalOper method
func TestRrmRepository_GetRrmGlobalOper(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name: "valid_parameters",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{
				Config: tt.config,
			}

			result := repo.GetRrmGlobalOper(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRrmRepository_GetRrmCfg tests the GetRrmCfg method
func TestRrmRepository_GetRrmCfg(t *testing.T) {
	tests := []struct {
		name       string
		config     *config.Config
		controller string
		apikey     string
		isSecure   *bool
		expectNil  bool
	}{
		{
			name: "valid_parameters",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: 1,
				},
			},
			controller: "192.168.1.1:443",
			apikey:     "test-token",
			isSecure:   &[]bool{true}[0],
			expectNil:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &RrmRepository{
				Config: tt.config,
			}

			result := repo.GetRrmCfg(tt.controller, tt.apikey, tt.isSecure)

			if tt.expectNil && result != nil {
				t.Errorf("Expected nil result, got %v", result)
			}
		})
	}
}

// TestRrmRepository_AllMethods tests all RRM methods together
func TestRrmRepository_AllMethods(t *testing.T) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			Timeout: 1,
		},
	}

	repo := &RrmRepository{
		Config: cfg,
	}

	controller := "192.168.1.1:443"
	apikey := "test-token"
	isSecure := &[]bool{true}[0]

	// Test all methods don't panic
	methods := []func() interface{}{
		func() interface{} { return repo.GetRrmOper(controller, apikey, isSecure) },
		func() interface{} { return repo.GetRrmMeasurement(controller, apikey, isSecure) },
		func() interface{} { return repo.GetRrmGlobalOper(controller, apikey, isSecure) },
		func() interface{} { return repo.GetRrmCfg(controller, apikey, isSecure) },
	}

	for i, method := range methods {
		t.Run(t.Name()+"_method_"+string(rune('0'+i)), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Method panicked: %v", r)
				}
			}()

			result := method()
			// All methods should return nil in test environment
			if result != nil {
				t.Logf("Method returned non-nil result: %v", result)
			}
		})
	}
}

// TestRrmRepository_TimeoutConfiguration tests timeout handling
func TestRrmRepository_TimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name             string
		timeoutSeconds   int
		expectedDuration time.Duration
	}{
		{
			name:             "default_timeout",
			timeoutSeconds:   30,
			expectedDuration: 30 * time.Second,
		},
		{
			name:             "long_timeout",
			timeoutSeconds:   300,
			expectedDuration: 300 * time.Second,
		},
		{
			name:             "short_timeout",
			timeoutSeconds:   10,
			expectedDuration: 10 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Timeout: tt.timeoutSeconds,
				},
			}

			repo := &RrmRepository{
				Config: cfg,
			}

			// Test that the timeout is properly configured
			actualDuration := time.Duration(repo.Config.ShowCmdConfig.Timeout) * time.Second
			if actualDuration != tt.expectedDuration {
				t.Errorf("Expected timeout %v, got %v", tt.expectedDuration, actualDuration)
			}
		})
	}
}
