package application

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
)

// TestOverviewUsecaseInitialization tests OverviewUsecase initialization (Unit test)
func TestOverviewUsecaseInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "invoke_overview_usecase",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)

				overviewUsecase := app.InvokeOverviewUsecase()
				if overviewUsecase == nil {
					t.Error("InvokeOverviewUsecase returned nil")
				}
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

// TestOverviewUsecaseShowOverview tests ShowOverview functionality (Unit test)
func TestOverviewUsecaseShowOverview(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "show_overview_basic",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				// Test ShowOverview method with controllers and isSecure parameters
				controllers := []config.Controller{{Hostname: "controller1", AccessToken: "token1"}}
				isSecure := true

				testUtilsInstance.assertNoPanic(t, func() {
					_ = overviewUsecase.ShowOverview(&controllers, &isSecure)
				})
			},
		},
		{
			name: "show_overview_insecure",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				controllers := []config.Controller{{Hostname: "controller1", AccessToken: "token1"}}
				isSecure := false

				testUtilsInstance.assertNoPanic(t, func() {
					_ = overviewUsecase.ShowOverview(&controllers, &isSecure)
				})
			},
		},
		{
			name: "show_overview_multiple_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				controllers := []config.Controller{
					{Hostname: "controller1", AccessToken: "token1"},
					{Hostname: "controller2", AccessToken: "token2"},
					{Hostname: "controller3", AccessToken: "token3"},
				}
				isSecure := true

				testUtilsInstance.assertNoPanic(t, func() {
					_ = overviewUsecase.ShowOverview(&controllers, &isSecure)
				})
			},
		},
		{
			name: "show_overview_empty_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				var controllers []config.Controller
				isSecure := true

				testUtilsInstance.assertNoPanic(t, func() {
					_ = overviewUsecase.ShowOverview(&controllers, &isSecure)
				})
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

// TestOverviewUsecaseShowOverviewDataProcessing tests data processing and filtering (Unit test)
func TestOverviewUsecaseShowOverviewDataProcessing(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "show_overview_with_radio_filter",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "1" // Set radio filter
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				controllers := []config.Controller{{Hostname: "controller1", AccessToken: "token1"}}
				isSecure := true

				result := overviewUsecase.ShowOverview(&controllers, &isSecure)
				// Should return empty slice when no data matches filter
				if result == nil {
					t.Error("ShowOverview should return empty slice, not nil")
				}
			},
		},
		{
			name: "show_overview_with_empty_radio_filter",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "" // Empty radio filter
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				controllers := []config.Controller{{Hostname: "controller1", AccessToken: "token1"}}
				isSecure := true

				result := overviewUsecase.ShowOverview(&controllers, &isSecure)
				if result == nil {
					t.Error("ShowOverview should return empty slice, not nil")
				}
			},
		},
		{
			name: "show_overview_nil_controllers",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				isSecure := true

				result := overviewUsecase.ShowOverview(nil, &isSecure)
				if result == nil {
					t.Error("ShowOverview should return empty slice for nil controllers")
				}
				if len(result) != 0 {
					t.Error("ShowOverview should return empty slice for nil controllers")
				}
			},
		},
		{
			name: "show_overview_nil_repository",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				app := New(cfg, nil) // Nil repository
				overviewUsecase := app.InvokeOverviewUsecase()

				controllers := []config.Controller{{Hostname: "controller1", AccessToken: "token1"}}
				isSecure := true

				result := overviewUsecase.ShowOverview(&controllers, &isSecure)
				if result == nil {
					t.Error("ShowOverview should return empty slice for nil repository")
				}
				if len(result) != 0 {
					t.Error("ShowOverview should return empty slice for nil repository")
				}
			},
		},
		{
			name: "filter_by_radio_with_nil_config",
			test: func(t *testing.T) {
				repo := testUtilsInstance.createMockRepository(nil)
				app := New(nil, repo) // Nil config
				overviewUsecase := app.InvokeOverviewUsecase()

				// Test filterByRadio method directly with nil config
				testData := []*ShowOverviewData{
					{SlotID: 0, ApMac: "test-mac"},
				}

				result := overviewUsecase.filterByRadio(testData)
				if result == nil {
					t.Error("filterByRadio should return data when config is nil")
				}
				if len(result) != len(testData) {
					t.Error("filterByRadio should return all data when config is nil")
				}
			},
		},
		{
			name: "filter_by_radio_with_matching_slot",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "1" // Filter for slot 1
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				testData := []*ShowOverviewData{
					{SlotID: 0, ApMac: "test-mac-0"},
					{SlotID: 1, ApMac: "test-mac-1"},
					{SlotID: 2, ApMac: "test-mac-2"},
				}

				result := overviewUsecase.filterByRadio(testData)
				if result == nil {
					t.Error("filterByRadio should return filtered data")
				}
				if len(result) != 1 {
					t.Errorf("filterByRadio should return 1 item, got %d", len(result))
				}
				if len(result) > 0 && result[0].SlotID != 1 {
					t.Errorf("filterByRadio should return item with SlotID 1, got %d", result[0].SlotID)
				}
			},
		},
		{
			name: "filter_by_radio_with_no_matching_slot",
			test: func(t *testing.T) {
				cfg := testUtilsInstance.createMockConfig()
				cfg.ShowCmdConfig.Radio = "5" // Filter for slot 5 (non-existent)
				repo := testUtilsInstance.createMockRepository(cfg)
				app := New(cfg, repo)
				overviewUsecase := app.InvokeOverviewUsecase()

				testData := []*ShowOverviewData{
					{SlotID: 0, ApMac: "test-mac-0"},
					{SlotID: 1, ApMac: "test-mac-1"},
				}

				result := overviewUsecase.filterByRadio(testData)
				if result == nil {
					t.Error("filterByRadio should return empty slice for no matches")
				}
				if len(result) != 0 {
					t.Errorf("filterByRadio should return empty slice, got %d items", len(result))
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
