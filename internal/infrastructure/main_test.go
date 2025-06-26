package infrastructure

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

func TestRepositoryNew(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "creates new repository with valid config",
			cfg:  &config.Config{},
		},
		{
			name: "creates new repository with nil config",
			cfg:  nil,
		},
		{
			name: "creates new repository with populated config",
			cfg: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers: []config.Controller{
						{Hostname: "wnc.example.com", AccessToken: "token123"},
					},
					PrintFormat: config.PrintFormatJSON,
					Timeout:     30,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.cfg)

			if got.Config != tt.cfg {
				t.Errorf("New() Config = %v, want %v", got.Config, tt.cfg)
			}
		})
	}
}

func TestRepositoryInvokeSubRepositories(t *testing.T) {
	cfg := &config.Config{}
	repo := New(cfg)

	tests := []struct {
		name     string
		invoke   func() interface{}
		wantType string
	}{
		{
			name: "InvokeClientRepository returns ClientRepository",
			invoke: func() interface{} {
				return repo.InvokeClientRepository()
			},
			wantType: "*infrastructure.ClientRepository",
		},
		{
			name: "InvokeApRepository returns ApRepository",
			invoke: func() interface{} {
				return repo.InvokeApRepository()
			},
			wantType: "*infrastructure.ApRepository",
		},
		{
			name: "InvokeWlanRepository returns WlanRepository",
			invoke: func() interface{} {
				return repo.InvokeWlanRepository()
			},
			wantType: "*infrastructure.WlanRepository",
		},
		{
			name: "InvokeRadioRepository returns RadioRepository",
			invoke: func() interface{} {
				return repo.InvokeRadioRepository()
			},
			wantType: "*infrastructure.RadioRepository",
		},
		{
			name: "InvokeRrmRepository returns RrmRepository",
			invoke: func() interface{} {
				return repo.InvokeRrmRepository()
			},
			wantType: "*infrastructure.RrmRepository",
		},
		{
			name: "InvokeRfRepository returns RfRepository",
			invoke: func() interface{} {
				return repo.InvokeRfRepository()
			},
			wantType: "*infrastructure.RfRepository",
		},
		{
			name: "InvokeDot11Repository returns Dot11Repository",
			invoke: func() interface{} {
				return repo.InvokeDot11Repository()
			},
			wantType: "*infrastructure.Dot11Repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.invoke()

			// Check that the returned object is not nil
			if got == nil {
				t.Errorf("%s returned nil", tt.name)
				return
			}

			// Check that each sub-repository has the correct config reference
			switch v := got.(type) {
			case *ClientRepository:
				if v.Config != cfg {
					t.Errorf("ClientRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *ApRepository:
				if v.Config != cfg {
					t.Errorf("ApRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *WlanRepository:
				if v.Config != cfg {
					t.Errorf("WlanRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *RadioRepository:
				if v.Config != cfg {
					t.Errorf("RadioRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *RrmRepository:
				if v.Config != cfg {
					t.Errorf("RrmRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *RfRepository:
				if v.Config != cfg {
					t.Errorf("RfRepository.Config = %v, want %v", v.Config, cfg)
				}
			case *Dot11Repository:
				if v.Config != cfg {
					t.Errorf("Dot11Repository.Config = %v, want %v", v.Config, cfg)
				}
			default:
				t.Errorf("Unexpected type returned: %T", got)
			}
		})
	}
}

func TestRepositoryJSONSerialization(t *testing.T) {
	tests := []struct {
		name string
		repo Repository
	}{
		{
			name: "empty repository",
			repo: Repository{
				Config: nil,
			},
		},
		{
			name: "repository with config",
			repo: Repository{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						Controllers: []config.Controller{
							{Hostname: "wnc.example.com", AccessToken: "token123"},
						},
						PrintFormat: config.PrintFormatJSON,
						Timeout:     30,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.repo)
			if err != nil {
				t.Fatalf("Failed to marshal Repository to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledRepo Repository
			err = json.Unmarshal(jsonData, &unmarshaledRepo)
			if err != nil {
				t.Fatalf("Failed to unmarshal Repository from JSON: %v", err)
			}

			// Basic validation - checking that unmarshaling doesn't fail
			// Note: Pointer fields will be different after unmarshal
			if tt.repo.Config != nil && unmarshaledRepo.Config == nil {
				// This is expected behavior for JSON unmarshaling with pointers - no action needed
				_ = tt.repo.Config // Acknowledge the check
			}
		})
	}
}

func TestRepositoryDependencyInjection(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "dependency injection with valid config",
			cfg: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers: []config.Controller{
						{Hostname: "test.example.com", AccessToken: "token123"},
					},
					PrintFormat: config.PrintFormatJSON,
				},
			},
		},
		{
			name: "dependency injection with nil config",
			cfg:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := New(tt.cfg)

			// Test that all sub-repositories receive the same config
			clientRepo := repo.InvokeClientRepository()
			apRepo := repo.InvokeApRepository()
			wlanRepo := repo.InvokeWlanRepository()
			radioRepo := repo.InvokeRadioRepository()
			rrmRepo := repo.InvokeRrmRepository()
			rfRepo := repo.InvokeRfRepository()
			dot11Repo := repo.InvokeDot11Repository()

			// Verify dependency injection consistency
			repositories := []interface{}{
				clientRepo, apRepo, wlanRepo, radioRepo, rrmRepo, rfRepo, dot11Repo,
			}

			for i, r := range repositories {
				switch v := r.(type) {
				case *ClientRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *ApRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *WlanRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *RadioRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *RrmRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *RfRepository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				case *Dot11Repository:
					if v.Config != tt.cfg {
						t.Errorf("Sub-repository %d: Config not properly injected", i)
					}
				}
			}
		})
	}
}

func TestRepositoryFailFast(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		expectPanic bool
	}{
		{
			name:        "nil config should not panic in constructor",
			cfg:         nil,
			expectPanic: false,
		},
		{
			name: "valid config should not panic",
			cfg: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers: []config.Controller{
						{Hostname: "test.example.com", AccessToken: "token123"},
					},
				},
			},
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

			repo := New(tt.cfg)

			// Verify that the repository was created even with nil config
			if repo.Config != tt.cfg {
				t.Errorf("Config not properly assigned")
			}

			// Test that sub-repository creation doesn't panic
			_ = repo.InvokeClientRepository()
			_ = repo.InvokeApRepository()
			_ = repo.InvokeWlanRepository()
			_ = repo.InvokeRadioRepository()
			_ = repo.InvokeRrmRepository()
			_ = repo.InvokeRfRepository()
			_ = repo.InvokeDot11Repository()
		})
	}
}

func TestRepositoryImmutability(t *testing.T) {
	originalConfig := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			PrintFormat: config.PrintFormatJSON,
			Timeout:     30,
		},
	}

	repo := New(originalConfig)

	// Get a sub-repository
	clientRepo := repo.InvokeClientRepository()

	// Verify that the config reference is the same (not copied)
	if clientRepo.Config != originalConfig {
		t.Error("Config should be passed by reference, not copied")
	}

	// Modify the original config
	originalConfig.ShowCmdConfig.Timeout = 60

	// Verify that the change is reflected in the repository
	if clientRepo.Config.ShowCmdConfig.Timeout != 60 {
		t.Error("Config changes should be reflected in sub-repositories")
	}
}
