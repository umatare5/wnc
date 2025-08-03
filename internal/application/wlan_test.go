package application

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestWlanUsecaseInitialization tests WlanUsecase initialization (Unit test)
func TestWlanUsecaseInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "invoke_wlan_usecase",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				wlanUsecase := app.InvokeWlanUsecase()
				if wlanUsecase == nil {
					t.Error("InvokeWlanUsecase returned nil")
					return
				}
				if wlanUsecase.Config != cfg {
					t.Error("WlanUsecase Config not set correctly")
				}
				if wlanUsecase.Repository != repo {
					t.Error("WlanUsecase Repository not set correctly")
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

// TestWlanUsecaseShowWlan tests WlanUsecase ShowWlan method (Unit test)
func TestWlanUsecaseShowWlan(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "ShowWlan_with_nil_repository",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				app := New(cfg, nil)
				wlanUsecase := app.InvokeWlanUsecase()

				controllers := &[]config.Controller{{Hostname: "test", AccessToken: "token"}}
				isSecure := false
				result := wlanUsecase.ShowWlan(controllers, &isSecure)

				if result == nil {
					t.Error("ShowWlan should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowWlan should return empty slice when repository is nil")
				}
			},
		},
		{
			name: "ShowWlan_with_nil_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				wlanUsecase := app.InvokeWlanUsecase()

				isSecure := false
				result := wlanUsecase.ShowWlan(nil, &isSecure)

				if result == nil {
					t.Error("ShowWlan should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowWlan should return empty slice when controllers is nil")
				}
			},
		},
		{
			name: "ShowWlan_with_empty_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				wlanUsecase := app.InvokeWlanUsecase()

				controllers := &[]config.Controller{}
				isSecure := false
				result := wlanUsecase.ShowWlan(controllers, &isSecure)

				if result == nil {
					t.Error("ShowWlan should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowWlan should return empty slice when controllers is empty")
				}
			},
		},
		{
			name: "ShowWlan_with_valid_controllers_but_nil_data",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				wlanUsecase := app.InvokeWlanUsecase()

				controllers := &[]config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}
				isSecure := false
				result := wlanUsecase.ShowWlan(controllers, &isSecure)

				// Should not panic and should return empty slice when API returns nil
				if result == nil {
					t.Error("ShowWlan should return empty slice, not nil")
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
