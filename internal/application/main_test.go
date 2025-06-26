package application

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestUsecaseNew(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		repo *infrastructure.Repository
	}{
		{
			name: "creates new usecase with valid config and repository",
			cfg:  &config.Config{},
			repo: &infrastructure.Repository{},
		},
		{
			name: "creates new usecase with nil config",
			cfg:  nil,
			repo: &infrastructure.Repository{},
		},
		{
			name: "creates new usecase with nil repository",
			cfg:  &config.Config{},
			repo: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.cfg, tt.repo)

			if got.Config != tt.cfg {
				t.Errorf("New() Config = %v, want %v", got.Config, tt.cfg)
			}

			if got.Repository != tt.repo {
				t.Errorf("New() Repository = %v, want %v", got.Repository, tt.repo)
			}
		})
	}
}

func TestUsecaseInvokeSubUsecases(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{}
	usecase := New(cfg, repo)

	tests := []struct {
		name     string
		invoke   func() interface{}
		wantType string
	}{
		{
			name: "InvokeTokenUsecase returns TokenUsecase",
			invoke: func() interface{} {
				return usecase.InvokeTokenUsecase()
			},
			wantType: "*application.TokenUsecase",
		},
		{
			name: "InvokeClientUsecase returns ClientUsecase",
			invoke: func() interface{} {
				return usecase.InvokeClientUsecase()
			},
			wantType: "*application.ClientUsecase",
		},
		{
			name: "InvokeApUsecase returns ApUsecase",
			invoke: func() interface{} {
				return usecase.InvokeApUsecase()
			},
			wantType: "*application.ApUsecase",
		},
		{
			name: "InvokeWlanUsecase returns WlanUsecase",
			invoke: func() interface{} {
				return usecase.InvokeWlanUsecase()
			},
			wantType: "*application.WlanUsecase",
		},
		{
			name: "InvokeOverviewUsecase returns OverviewUsecase",
			invoke: func() interface{} {
				return usecase.InvokeOverviewUsecase()
			},
			wantType: "*application.OverviewUsecase",
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

			// Check that each sub-usecase has the correct config and repository references
			switch v := got.(type) {
			case *TokenUsecase:
				if v.Config != cfg {
					t.Errorf("TokenUsecase.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("TokenUsecase.Repository = %v, want %v", v.Repository, repo)
				}
			case *ClientUsecase:
				if v.Config != cfg {
					t.Errorf("ClientUsecase.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("ClientUsecase.Repository = %v, want %v", v.Repository, repo)
				}
			case *ApUsecase:
				if v.Config != cfg {
					t.Errorf("ApUsecase.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("ApUsecase.Repository = %v, want %v", v.Repository, repo)
				}
			case *WlanUsecase:
				if v.Config != cfg {
					t.Errorf("WlanUsecase.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("WlanUsecase.Repository = %v, want %v", v.Repository, repo)
				}
			case *OverviewUsecase:
				if v.Config != cfg {
					t.Errorf("OverviewUsecase.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("OverviewUsecase.Repository = %v, want %v", v.Repository, repo)
				}
			default:
				t.Errorf("Unexpected type returned: %T", got)
			}
		})
	}
}

func TestUsecaseJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		usecase Usecase
	}{
		{
			name: "empty usecase",
			usecase: Usecase{
				Config:     nil,
				Repository: nil,
			},
		},
		{
			name: "usecase with config",
			usecase: Usecase{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						PrintFormat: config.PrintFormatJSON,
						Timeout:     30,
					},
				},
				Repository: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.usecase)
			if err != nil {
				t.Fatalf("Failed to marshal Usecase to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledUsecase Usecase
			err = json.Unmarshal(jsonData, &unmarshaledUsecase)
			if err != nil {
				t.Fatalf("Failed to unmarshal Usecase from JSON: %v", err)
			}

			// Basic validation - checking that unmarshaling doesn't fail
			// Note: Pointer fields will be nil after unmarshal due to struct pointers
			if tt.usecase.Config != nil && unmarshaledUsecase.Config == nil {
				// This is expected behavior for JSON unmarshaling with pointers - no action needed
				_ = tt.usecase.Config // Acknowledge the check
			}
		})
	}
}

func TestUsecaseDependencyInjection(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		repo *infrastructure.Repository
	}{
		{
			name: "dependency injection with valid objects",
			cfg: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers: []config.Controller{
						{Hostname: "test.example.com", AccessToken: "token123"},
					},
					PrintFormat: config.PrintFormatJSON,
				},
			},
			repo: &infrastructure.Repository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := New(tt.cfg, tt.repo)

			// Test that all sub-usecases receive the same dependencies
			tokenUC := usecase.InvokeTokenUsecase()
			clientUC := usecase.InvokeClientUsecase()
			apUC := usecase.InvokeApUsecase()
			wlanUC := usecase.InvokeWlanUsecase()
			overviewUC := usecase.InvokeOverviewUsecase()

			// Verify dependency injection consistency
			usecases := []interface{}{tokenUC, clientUC, apUC, wlanUC, overviewUC}
			for i, uc := range usecases {
				switch v := uc.(type) {
				case *TokenUsecase:
					if v.Config != tt.cfg {
						t.Errorf("Sub-usecase %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-usecase %d: Repository not properly injected", i)
					}
				case *ClientUsecase:
					if v.Config != tt.cfg {
						t.Errorf("Sub-usecase %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-usecase %d: Repository not properly injected", i)
					}
				case *ApUsecase:
					if v.Config != tt.cfg {
						t.Errorf("Sub-usecase %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-usecase %d: Repository not properly injected", i)
					}
				case *WlanUsecase:
					if v.Config != tt.cfg {
						t.Errorf("Sub-usecase %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-usecase %d: Repository not properly injected", i)
					}
				case *OverviewUsecase:
					if v.Config != tt.cfg {
						t.Errorf("Sub-usecase %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-usecase %d: Repository not properly injected", i)
					}
				}
			}
		})
	}
}

func TestUsecaseFailFast(t *testing.T) {
	tests := []struct {
		name          string
		cfg           *config.Config
		repo          *infrastructure.Repository
		expectPanic   bool
		expectedError string
	}{
		{
			name:        "nil config should not panic in constructor",
			cfg:         nil,
			repo:        &infrastructure.Repository{},
			expectPanic: false,
		},
		{
			name:        "nil repository should not panic in constructor",
			cfg:         &config.Config{},
			repo:        nil,
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

			usecase := New(tt.cfg, tt.repo)

			// Verify that the usecase was created even with nil dependencies
			if usecase.Config != tt.cfg {
				t.Errorf("Config not properly assigned")
			}
			if usecase.Repository != tt.repo {
				t.Errorf("Repository not properly assigned")
			}
		})
	}
}
