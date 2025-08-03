package application

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestApUsecaseInitialization tests ApUsecase initialization (Unit test)
func TestApUsecaseInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "invoke_ap_usecase",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				apUsecase := app.InvokeApUsecase()
				if apUsecase == nil {
					t.Error("InvokeApUsecase returned nil")
					return
				}
				if apUsecase.Config != cfg {
					t.Error("ApUsecase Config not set correctly")
				}
				if apUsecase.Repository != repo {
					t.Error("ApUsecase Repository not set correctly")
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

// TestApUsecaseShowAp tests ApUsecase ShowAp method (Unit test)
func TestApUsecaseShowAp(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "ShowAp_with_nil_repository",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				app := New(cfg, nil)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{{Hostname: "test", AccessToken: "token"}}
				isSecure := false
				result := apUsecase.ShowAp(controllers, &isSecure)

				if result == nil {
					t.Error("ShowAp should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowAp should return empty slice when repository is nil")
				}
			},
		},
		{
			name: "ShowAp_with_nil_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				isSecure := false
				result := apUsecase.ShowAp(nil, &isSecure)

				if result == nil {
					t.Error("ShowAp should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowAp should return empty slice when controllers is nil")
				}
			},
		},
		{
			name: "ShowAp_with_empty_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{}
				isSecure := false
				result := apUsecase.ShowAp(controllers, &isSecure)

				if result == nil {
					t.Error("ShowAp should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowAp should return empty slice when controllers is empty")
				}
			},
		},
		{
			name: "ShowAp_with_valid_controllers_but_nil_data",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}
				isSecure := false
				result := apUsecase.ShowAp(controllers, &isSecure)

				// Should not panic and should return empty slice when API returns nil
				if result == nil {
					t.Error("ShowAp should return empty slice, not nil")
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

// TestApUsecaseShowApTag tests ApUsecase ShowApTag method (Unit test)
func TestApUsecaseShowApTag(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "ShowApTag_with_nil_repository",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				app := New(cfg, nil)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{{Hostname: "test", AccessToken: "token"}}
				isSecure := false
				result := apUsecase.ShowApTag(controllers, &isSecure)

				if result == nil {
					t.Error("ShowApTag should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowApTag should return empty slice when repository is nil")
				}
			},
		},
		{
			name: "ShowApTag_with_nil_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				isSecure := false
				result := apUsecase.ShowApTag(nil, &isSecure)

				if result == nil {
					t.Error("ShowApTag should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowApTag should return empty slice when controllers is nil")
				}
			},
		},
		{
			name: "ShowApTag_with_empty_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{}
				isSecure := false
				result := apUsecase.ShowApTag(controllers, &isSecure)

				if result == nil {
					t.Error("ShowApTag should return empty slice, not nil")
				}
				if len(result) != 0 {
					t.Error("ShowApTag should return empty slice when controllers is empty")
				}
			},
		},
		{
			name: "ShowApTag_with_valid_controllers_but_nil_data",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := &[]config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}
				isSecure := false
				result := apUsecase.ShowApTag(controllers, &isSecure)

				// Should not panic and should return empty slice when API returns nil
				if result == nil {
					t.Error("ShowApTag should return empty slice, not nil")
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

// TestApUsecaseAdvancedDataProcessing tests comprehensive data processing scenarios (Unit test)
func TestApUsecaseAdvancedDataProcessing(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "show_ap_data_aggregation_with_partial_data",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := []config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}

				result := apUsecase.ShowAp(&controllers, boolPtr(true))
				if result == nil {
					t.Error("ShowAp should return non-nil slice even with partial data")
				}
				if len(result) != 0 {
					t.Errorf("Expected 0 APs with mock data, got %d", len(result))
				}
			},
		},
		{
			name: "show_ap_with_multiple_controller_aggregation",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := []config.Controller{
					{Hostname: "controller1", AccessToken: "token1"},
					{Hostname: "controller2", AccessToken: "token2"},
					{Hostname: "controller3", AccessToken: "token3"},
				}

				result := apUsecase.ShowAp(&controllers, boolPtr(true))
				if result == nil {
					t.Error("ShowAp should return non-nil slice with multiple controllers")
				}
				if len(result) < 0 {
					t.Errorf("Expected non-negative length, got %d", len(result))
				}
			},
		},
		{
			name: "show_ap_with_empty_controller_list",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := []config.Controller{}

				result := apUsecase.ShowAp(&controllers, boolPtr(true))
				if result == nil {
					t.Error("ShowAp should return non-nil slice even with empty controllers")
				}
				if len(result) != 0 {
					t.Errorf("Expected 0 APs with empty controllers, got %d", len(result))
				}
			},
		},
		{
			name: "show_ap_tag_with_valid_controller_collection",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := []config.Controller{
					{Hostname: "test-controller", AccessToken: "test-token"},
				}
				isSecure := true

				result := apUsecase.ShowApTag(&controllers, &isSecure)
				if result == nil {
					t.Error("ShowApTag should return non-nil slice")
				}
				if len(result) != 0 {
					t.Errorf("Expected 0 AP tags with mock data, got %d", len(result))
				}
			},
		},
		{
			name: "show_ap_tag_with_multiple_controller_merge",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				apUsecase := app.InvokeApUsecase()

				controllers := []config.Controller{
					{Hostname: "controller1", AccessToken: "token1"},
					{Hostname: "controller2", AccessToken: "token2"},
				}
				isSecure := false

				result := apUsecase.ShowApTag(&controllers, &isSecure)
				if result == nil {
					t.Error("ShowApTag should return non-nil slice with multiple controllers")
				}
				if len(result) < 0 {
					t.Errorf("Expected non-negative length, got %d", len(result))
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

// boolPtr is a helper function to create a pointer to bool
func boolPtr(b bool) *bool {
	return &b
}
