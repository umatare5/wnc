package framework

import (
	"encoding/json"
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework/show"
	"github.com/umatare5/wnc/internal/infrastructure"
)

func TestNewShowCli(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		repo *infrastructure.Repository
		uc   *application.Usecase
	}{
		{
			name: "creates new ShowCli with valid dependencies",
			cfg:  &config.Config{},
			repo: &infrastructure.Repository{},
			uc:   &application.Usecase{},
		},
		{
			name: "creates new ShowCli with nil config",
			cfg:  nil,
			repo: &infrastructure.Repository{},
			uc:   &application.Usecase{},
		},
		{
			name: "creates new ShowCli with nil repository",
			cfg:  &config.Config{},
			repo: nil,
			uc:   &application.Usecase{},
		},
		{
			name: "creates new ShowCli with nil usecase",
			cfg:  &config.Config{},
			repo: &infrastructure.Repository{},
			uc:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewShowCli(tt.cfg, tt.repo, tt.uc)

			if got.Config != tt.cfg {
				t.Errorf("NewShowCli() Config = %v, want %v", got.Config, tt.cfg)
			}

			if got.Repository != tt.repo {
				t.Errorf("NewShowCli() Repository = %v, want %v", got.Repository, tt.repo)
			}

			if got.Usecase != tt.uc {
				t.Errorf("NewShowCli() Usecase = %v, want %v", got.Usecase, tt.uc)
			}
		})
	}
}

func TestShowCliInvokeSubClis(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{}
	uc := &application.Usecase{}
	showCli := NewShowCli(cfg, repo, uc)

	tests := []struct {
		name     string
		invoke   func() interface{}
		wantType string
	}{
		{
			name: "InvokeClientCli returns ClientCli",
			invoke: func() interface{} {
				return showCli.InvokeClientCli()
			},
			wantType: "*show.ClientCli",
		},
		{
			name: "InvokeApCli returns ApCli",
			invoke: func() interface{} {
				return showCli.InvokeApCli()
			},
			wantType: "*show.ApCli",
		},
		{
			name: "InvokeApTagCli returns ApTagCli",
			invoke: func() interface{} {
				return showCli.InvokeApTagCli()
			},
			wantType: "*show.ApTagCli",
		},
		{
			name: "InvokeWlanCli returns WlanCli",
			invoke: func() interface{} {
				return showCli.InvokeWlanCli()
			},
			wantType: "*show.WlanCli",
		},
		{
			name: "InvokeOverviewCli returns OverviewCli",
			invoke: func() interface{} {
				return showCli.InvokeOverviewCli()
			},
			wantType: "*show.OverviewCli",
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

			// Check that each sub-CLI has the correct dependency references
			switch v := got.(type) {
			case *show.ClientCli:
				if v.Config != cfg {
					t.Errorf("ClientCli.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("ClientCli.Repository = %v, want %v", v.Repository, repo)
				}
				if v.Usecase != uc {
					t.Errorf("ClientCli.Usecase = %v, want %v", v.Usecase, uc)
				}
			case *show.ApCli:
				if v.Config != cfg {
					t.Errorf("ApCli.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("ApCli.Repository = %v, want %v", v.Repository, repo)
				}
				if v.Usecase != uc {
					t.Errorf("ApCli.Usecase = %v, want %v", v.Usecase, uc)
				}
			case *show.ApTagCli:
				if v.Config != cfg {
					t.Errorf("ApTagCli.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("ApTagCli.Repository = %v, want %v", v.Repository, repo)
				}
				if v.Usecase != uc {
					t.Errorf("ApTagCli.Usecase = %v, want %v", v.Usecase, uc)
				}
			case *show.WlanCli:
				if v.Config != cfg {
					t.Errorf("WlanCli.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("WlanCli.Repository = %v, want %v", v.Repository, repo)
				}
				if v.Usecase != uc {
					t.Errorf("WlanCli.Usecase = %v, want %v", v.Usecase, uc)
				}
			case *show.OverviewCli:
				if v.Config != cfg {
					t.Errorf("OverviewCli.Config = %v, want %v", v.Config, cfg)
				}
				if v.Repository != repo {
					t.Errorf("OverviewCli.Repository = %v, want %v", v.Repository, repo)
				}
				if v.Usecase != uc {
					t.Errorf("OverviewCli.Usecase = %v, want %v", v.Usecase, uc)
				}
			default:
				t.Errorf("Unexpected type returned: %T", got)
			}
		})
	}
}

func TestShowCliJSONSerialization(t *testing.T) {
	tests := []struct {
		name    string
		showCli ShowCli
	}{
		{
			name: "empty ShowCli",
			showCli: ShowCli{
				Config:     nil,
				Repository: nil,
				Usecase:    nil,
			},
		},
		{
			name: "ShowCli with dependencies",
			showCli: ShowCli{
				Config: &config.Config{
					ShowCmdConfig: config.ShowCmdConfig{
						PrintFormat: config.PrintFormatJSON,
						Timeout:     30,
					},
				},
				Repository: &infrastructure.Repository{},
				Usecase:    &application.Usecase{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			jsonData, err := json.Marshal(tt.showCli)
			if err != nil {
				t.Fatalf("Failed to marshal ShowCli to JSON: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaledShowCli ShowCli
			err = json.Unmarshal(jsonData, &unmarshaledShowCli)
			if err != nil {
				t.Fatalf("Failed to unmarshal ShowCli from JSON: %v", err)
			}

			// Basic validation - checking that unmarshaling doesn't fail
			// Note: Pointer fields will be different after unmarshal
			if tt.showCli.Config != nil && unmarshaledShowCli.Config == nil {
				// This is expected behavior for JSON unmarshaling with pointers
				t.Logf("Config field correctly handled during JSON unmarshal")
			}
		})
	}
}

func TestShowCliDependencyInjection(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		repo *infrastructure.Repository
		uc   *application.Usecase
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
			uc:   &application.Usecase{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			showCli := NewShowCli(tt.cfg, tt.repo, tt.uc)

			// Test that all sub-CLIs receive the same dependencies
			clientCli := showCli.InvokeClientCli()
			apCli := showCli.InvokeApCli()
			apTagCli := showCli.InvokeApTagCli()
			wlanCli := showCli.InvokeWlanCli()
			overviewCli := showCli.InvokeOverviewCli()

			// Verify dependency injection consistency
			clis := []interface{}{clientCli, apCli, apTagCli, wlanCli, overviewCli}
			for i, c := range clis {
				switch v := c.(type) {
				case *show.ClientCli:
					if v.Config != tt.cfg {
						t.Errorf("Sub-CLI %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-CLI %d: Repository not properly injected", i)
					}
					if v.Usecase != tt.uc {
						t.Errorf("Sub-CLI %d: Usecase not properly injected", i)
					}
				case *show.ApCli:
					if v.Config != tt.cfg {
						t.Errorf("Sub-CLI %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-CLI %d: Repository not properly injected", i)
					}
					if v.Usecase != tt.uc {
						t.Errorf("Sub-CLI %d: Usecase not properly injected", i)
					}
				case *show.ApTagCli:
					if v.Config != tt.cfg {
						t.Errorf("Sub-CLI %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-CLI %d: Repository not properly injected", i)
					}
					if v.Usecase != tt.uc {
						t.Errorf("Sub-CLI %d: Usecase not properly injected", i)
					}
				case *show.WlanCli:
					if v.Config != tt.cfg {
						t.Errorf("Sub-CLI %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-CLI %d: Repository not properly injected", i)
					}
					if v.Usecase != tt.uc {
						t.Errorf("Sub-CLI %d: Usecase not properly injected", i)
					}
				case *show.OverviewCli:
					if v.Config != tt.cfg {
						t.Errorf("Sub-CLI %d: Config not properly injected", i)
					}
					if v.Repository != tt.repo {
						t.Errorf("Sub-CLI %d: Repository not properly injected", i)
					}
					if v.Usecase != tt.uc {
						t.Errorf("Sub-CLI %d: Usecase not properly injected", i)
					}
				}
			}
		})
	}
}

func TestShowCliFailFast(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.Config
		repo        *infrastructure.Repository
		uc          *application.Usecase
		expectPanic bool
	}{
		{
			name:        "nil dependencies should not panic in constructor",
			cfg:         nil,
			repo:        nil,
			uc:          nil,
			expectPanic: false,
		},
		{
			name: "valid dependencies should not panic",
			cfg: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					Controllers: []config.Controller{
						{Hostname: "test.example.com", AccessToken: "token123"},
					},
				},
			},
			repo:        &infrastructure.Repository{},
			uc:          &application.Usecase{},
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

			showCli := NewShowCli(tt.cfg, tt.repo, tt.uc)

			// Verify that the ShowCli was created even with nil dependencies
			if showCli.Config != tt.cfg {
				t.Errorf("Config not properly assigned")
			}
			if showCli.Repository != tt.repo {
				t.Errorf("Repository not properly assigned")
			}
			if showCli.Usecase != tt.uc {
				t.Errorf("Usecase not properly assigned")
			}

			// Test that sub-CLI creation doesn't panic
			_ = showCli.InvokeClientCli()
			_ = showCli.InvokeApCli()
			_ = showCli.InvokeApTagCli()
			_ = showCli.InvokeWlanCli()
			_ = showCli.InvokeOverviewCli()
		})
	}
}
