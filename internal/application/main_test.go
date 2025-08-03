package application

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestNew tests application layer initialization (Unit test)
func TestNew(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "create_application_layer",
			test: func(t *testing.T) {
				cfg := &config.Config{}
				repo := &infrastructure.Repository{}

				app := New(cfg, repo)
				if app.Config != cfg {
					t.Error("New() Config not set correctly")
				}
				if app.Repository != repo {
					t.Error("New() Repository not set correctly")
				}
			},
		},
		{
			name: "create_with_nil_config",
			test: func(t *testing.T) {
				repo := &infrastructure.Repository{}

				app := New(nil, repo)
				if app.Config != nil {
					t.Error("New() should accept nil Config")
				}
				if app.Repository != repo {
					t.Error("New() Repository not set correctly")
				}
			},
		},
		{
			name: "create_with_nil_repository",
			test: func(t *testing.T) {
				cfg := &config.Config{}

				app := New(cfg, nil)
				if app.Config != cfg {
					t.Error("New() Config not set correctly")
				}
				if app.Repository != nil {
					t.Error("New() should accept nil Repository")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtilsInstance.assertNoPanic(t, func() {
				tt.test(t)
			})
		})
	}
}

// TestUsecaseStructure tests usecase structure and initialization (Unit test)
func TestUsecaseStructure(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "usecase_initialization",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)

				app := New(cfg, repo)
				if app.Config == nil {
					t.Error("Usecase Config should not be nil")
				}
				if app.Repository == nil {
					t.Error("Usecase Repository should not be nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtilsInstance.assertNoPanic(t, func() {
				tt.test(t)
			})
		})
	}
}

// TestAllUsecaseInvokers tests all usecase invoker methods (Unit test)
func TestAllUsecaseInvokers(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "all_invokers_return_non-nil",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				// Test all invoker methods
				tokenUsecase := app.InvokeTokenUsecase()
				if tokenUsecase == nil {
					t.Error("InvokeTokenUsecase returned nil")
				}

				clientUsecase := app.InvokeClientUsecase()
				if clientUsecase == nil {
					t.Error("InvokeClientUsecase returned nil")
				}

				apUsecase := app.InvokeApUsecase()
				if apUsecase == nil {
					t.Error("InvokeApUsecase returned nil")
				}

				wlanUsecase := app.InvokeWlanUsecase()
				if wlanUsecase == nil {
					t.Error("InvokeWlanUsecase returned nil")
				}

				overviewUsecase := app.InvokeOverviewUsecase()
				if overviewUsecase == nil {
					t.Error("InvokeOverviewUsecase returned nil")
				}
			},
		},
		{
			name: "all_invokers_have_correct_dependencies",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				// Test that all usecases have correct dependencies
				tokenUsecase := app.InvokeTokenUsecase()
				if tokenUsecase.Config != cfg {
					t.Error("TokenUsecase Config not set correctly")
				}
				if tokenUsecase.Repository != repo {
					t.Error("TokenUsecase Repository not set correctly")
				}

				clientUsecase := app.InvokeClientUsecase()
				if clientUsecase.Config != cfg {
					t.Error("ClientUsecase Config not set correctly")
				}
				if clientUsecase.Repository != repo {
					t.Error("ClientUsecase Repository not set correctly")
				}

				apUsecase := app.InvokeApUsecase()
				if apUsecase.Config != cfg {
					t.Error("ApUsecase Config not set correctly")
				}
				if apUsecase.Repository != repo {
					t.Error("ApUsecase Repository not set correctly")
				}

				wlanUsecase := app.InvokeWlanUsecase()
				if wlanUsecase.Config != cfg {
					t.Error("WlanUsecase Config not set correctly")
				}
				if wlanUsecase.Repository != repo {
					t.Error("WlanUsecase Repository not set correctly")
				}

				overviewUsecase := app.InvokeOverviewUsecase()
				if overviewUsecase.Config != cfg {
					t.Error("OverviewUsecase Config not set correctly")
				}
				if overviewUsecase.Repository != repo {
					t.Error("OverviewUsecase Repository not set correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUtilsInstance.assertNoPanic(t, func() {
				tt.test(t)
			})
		})
	}
}
