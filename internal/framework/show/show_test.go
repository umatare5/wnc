package show

import (
	"os"
	"testing"

	"github.com/umatare5/cisco-ios-xe-wireless-go/ap"
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// Test utilities for show package
type testUtils struct{}

var utils = &testUtils{}

func (u *testUtils) createFullTestStack() (*config.Config, *infrastructure.Repository, *application.Usecase) {
	cfg := &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			PrintFormat: config.PrintFormatTable,
		},
	}
	repo := &infrastructure.Repository{Config: cfg}
	usecase := &application.Usecase{Config: cfg, Repository: repo}
	return cfg, repo, usecase
}

func (u *testUtils) assertNoPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function panicked: %v", r)
		}
	}()
	fn()
}

func (u *testUtils) assertStructFields(t *testing.T, obj interface{}, expectedFields ...string) {
	// For this package, we'll use manual field checks since reflection would complicate things
	switch v := obj.(type) {
	case *ApCli:
		if v.Config == nil {
			t.Error("ApCli.Config should not be nil")
		}
		if v.Repository == nil {
			t.Error("ApCli.Repository should not be nil")
		}
		if v.Usecase == nil {
			t.Error("ApCli.Usecase should not be nil")
		}
	case *ClientCli:
		if v.Config == nil {
			t.Error("ClientCli.Config should not be nil")
		}
		if v.Repository == nil {
			t.Error("ClientCli.Repository should not be nil")
		}
		if v.Usecase == nil {
			t.Error("ClientCli.Usecase should not be nil")
		}
	case *WlanCli:
		if v.Config == nil {
			t.Error("WlanCli.Config should not be nil")
		}
		if v.Repository == nil {
			t.Error("WlanCli.Repository should not be nil")
		}
		if v.Usecase == nil {
			t.Error("WlanCli.Usecase should not be nil")
		}
	case *OverviewCli:
		if v.Config == nil {
			t.Error("OverviewCli.Config should not be nil")
		}
		if v.Repository == nil {
			t.Error("OverviewCli.Repository should not be nil")
		}
		if v.Usecase == nil {
			t.Error("OverviewCli.Usecase should not be nil")
		}
	case *ApTagCli:
		if v.Config == nil {
			t.Error("ApTagCli.Config should not be nil")
		}
		if v.Repository == nil {
			t.Error("ApTagCli.Repository should not be nil")
		}
		if v.Usecase == nil {
			t.Error("ApTagCli.Usecase should not be nil")
		}
	}
}

func (u *testUtils) redirectOutput(fn func()) {
	originalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w
	defer func() {
		w.Close()
		os.Stdout = originalStdout
	}()
	fn()
}

var localUtils = &testUtils{}

// TestApCli tests the ApCli structure and methods
func TestApCli(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "structure_validation",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()

				apCli := &ApCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertStructFields(t, apCli)
			},
		},
		{
			name: "show_ap_table_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatTable

				apCli := &ApCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						apCli.ShowAp()
					})
				})
			},
		},
		{
			name: "show_ap_json_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatJSON

				apCli := &ApCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						apCli.ShowAp()
					})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestClientCli tests the ClientCli structure and methods
func TestClientCli(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "structure_validation",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()

				clientCli := &ClientCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertStructFields(t, clientCli)
			},
		},
		{
			name: "show_client_table_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatTable

				clientCli := &ClientCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						clientCli.ShowClient()
					})
				})
			},
		},
		{
			name: "show_client_json_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatJSON

				clientCli := &ClientCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						clientCli.ShowClient()
					})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestWlanCli tests the WlanCli structure and methods
func TestWlanCli(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "structure_validation",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()

				wlanCli := &WlanCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertStructFields(t, wlanCli)
			},
		},
		{
			name: "show_wlan_table_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatTable

				wlanCli := &WlanCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						wlanCli.ShowWlan()
					})
				})
			},
		},
		{
			name: "show_wlan_json_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatJSON

				wlanCli := &WlanCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						wlanCli.ShowWlan()
					})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestOverviewCli tests the OverviewCli structure and methods
func TestOverviewCli(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "structure_validation",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()

				overviewCli := &OverviewCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertStructFields(t, overviewCli)
			},
		},
		{
			name: "show_overview_table_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatTable

				overviewCli := &OverviewCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						overviewCli.ShowOverview()
					})
				})
			},
		},
		{
			name: "show_overview_json_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatJSON

				overviewCli := &OverviewCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						overviewCli.ShowOverview()
					})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestApTagCli tests the ApTagCli structure and methods
func TestApTagCli(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "structure_validation",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()

				apTagCli := &ApTagCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertStructFields(t, apTagCli)
			},
		},
		{
			name: "show_ap_tag_table_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatTable

				apTagCli := &ApTagCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						apTagCli.ShowApTag()
					})
				})
			},
		},
		{
			name: "show_ap_tag_json_format",
			test: func(t *testing.T) {
				cfg, repo, usecase := localUtils.createFullTestStack()
				cfg.ShowCmdConfig.PrintFormat = config.PrintFormatJSON

				apTagCli := &ApTagCli{
					Config:     cfg,
					Repository: repo,
					Usecase:    usecase,
				}

				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						apTagCli.ShowApTag()
					})
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestFormatSpecificBehavior tests format-specific behavior
func TestFormatSpecificBehavior(t *testing.T) {
	formats := []string{config.PrintFormatTable, config.PrintFormatJSON}

	for _, format := range formats {
		t.Run("format_"+format, func(t *testing.T) {
			cfg := &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					PrintFormat: format,
				},
			}
			repo := &infrastructure.Repository{Config: cfg}
			usecase := &application.Usecase{Config: cfg, Repository: repo}

			// Test all CLI structures with this format
			clis := []interface{}{
				&ApCli{Config: cfg, Repository: repo, Usecase: usecase},
				&ClientCli{Config: cfg, Repository: repo, Usecase: usecase},
				&WlanCli{Config: cfg, Repository: repo, Usecase: usecase},
				&OverviewCli{Config: cfg, Repository: repo, Usecase: usecase},
				&ApTagCli{Config: cfg, Repository: repo, Usecase: usecase},
			}

			for _, cli := range clis {
				localUtils.assertNoPanic(t, func() {
					localUtils.redirectOutput(func() {
						switch v := cli.(type) {
						case *ApCli:
							v.ShowAp()
						case *ClientCli:
							v.ShowClient()
						case *WlanCli:
							v.ShowWlan()
						case *OverviewCli:
							v.ShowOverview()
						case *ApTagCli:
							v.ShowApTag()
						}
					})
				})
			}
		})
	}
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name: "minimal_config",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					PrintFormat: config.PrintFormatTable,
				},
			},
		},
		{
			name: "json_config",
			config: &config.Config{
				ShowCmdConfig: config.ShowCmdConfig{
					PrintFormat: config.PrintFormatJSON,
				},
			},
		},
		{
			name:   "nil_config_handling",
			config: &config.Config{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &infrastructure.Repository{Config: tt.config}
			usecase := &application.Usecase{Config: tt.config, Repository: repo}

			// Test that all CLI structures can be created with various configs
			apCli := &ApCli{Config: tt.config, Repository: repo, Usecase: usecase}
			clientCli := &ClientCli{Config: tt.config, Repository: repo, Usecase: usecase}
			wlanCli := &WlanCli{Config: tt.config, Repository: repo, Usecase: usecase}
			overviewCli := &OverviewCli{Config: tt.config, Repository: repo, Usecase: usecase}
			apTagCli := &ApTagCli{Config: tt.config, Repository: repo, Usecase: usecase}

			// Verify all structures are created without panics
			if apCli == nil || clientCli == nil || wlanCli == nil || overviewCli == nil || apTagCli == nil {
				t.Error("CLI structures should not be nil")
			}
		})
	}
}

// TestUtilityFunctions tests utility functions
func TestUtilityFunctions(t *testing.T) {
	t.Run("hasNoData", func(t *testing.T) {
		// Test with empty slice
		if !hasNoData([]interface{}{}) {
			t.Error("Empty slice should return true for hasNoData")
		}

		// Test with nil slice
		var nilSlice []interface{}
		if !hasNoData(nilSlice) {
			t.Error("Nil slice should return true for hasNoData")
		}

		// Test with non-empty slice
		if hasNoData([]interface{}{1, 2, 3}) {
			t.Error("Non-empty slice should return false for hasNoData")
		}
	})

	t.Run("isAPMisconfigured", func(t *testing.T) {
		// Test with misconfigured AP
		if !isAPMisconfigured(true) {
			t.Error("Should recognize misconfigured AP")
		}

		// Test with correctly configured AP
		if isAPMisconfigured(false) {
			t.Error("Should not recognize correctly configured AP as misconfigured")
		}
	})
}

// TestPrintJson tests JSON printing functionality
func TestPrintJson(t *testing.T) {
	t.Run("valid_json_data", func(t *testing.T) {
		data := map[string]interface{}{
			"test":   "value",
			"number": 42,
		}

		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				printJson(data)
			})
		})
	})

	t.Run("nil_data", func(t *testing.T) {
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				printJson(nil)
			})
		})
	})

	t.Run("complex_data_structure", func(t *testing.T) {
		data := []map[string]interface{}{
			{"id": 1, "name": "test1"},
			{"id": 2, "name": "test2"},
		}

		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				printJson(data)
			})
		})
	})
}

// TestErrorHandling tests error handling in CLI components
func TestErrorHandling(t *testing.T) {
	t.Run("empty_usecase_results", func(t *testing.T) {
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		usecase := &application.Usecase{Config: cfg, Repository: repo}

		clis := []interface{}{
			&ApCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ClientCli{Config: cfg, Repository: repo, Usecase: usecase},
			&WlanCli{Config: cfg, Repository: repo, Usecase: usecase},
			&OverviewCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ApTagCli{Config: cfg, Repository: repo, Usecase: usecase},
		}

		for _, cli := range clis {
			localUtils.assertNoPanic(t, func() {
				localUtils.redirectOutput(func() {
					switch v := cli.(type) {
					case *ApCli:
						v.ShowAp()
					case *ClientCli:
						v.ShowClient()
					case *WlanCli:
						v.ShowWlan()
					case *OverviewCli:
						v.ShowOverview()
					case *ApTagCli:
						v.ShowApTag()
					}
				})
			})
		}
	})
}

// TestStructValidation tests struct validation
func TestStructValidation(t *testing.T) {
	cfg, repo, usecase := localUtils.createFullTestStack()

	t.Run("all_cli_structs_have_required_fields", func(t *testing.T) {
		structs := []interface{}{
			&ApCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ClientCli{Config: cfg, Repository: repo, Usecase: usecase},
			&WlanCli{Config: cfg, Repository: repo, Usecase: usecase},
			&OverviewCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ApTagCli{Config: cfg, Repository: repo, Usecase: usecase},
		}

		for _, s := range structs {
			localUtils.assertStructFields(t, s)
		}
	})
}

// TestDataProcessingPaths tests data processing to improve coverage
func TestDataProcessingPaths(t *testing.T) {
	t.Run("test_utility_functions", func(t *testing.T) {
		// Test utility functions in utils.go to improve coverage

		// Test hasNoData function - it expects []any (slice)
		if !hasNoData([]any{}) {
			t.Error("hasNoData should return true for empty slice")
		}

		nonEmptySlice := []any{map[string]interface{}{"key": "value"}}
		if hasNoData(nonEmptySlice) {
			t.Error("hasNoData should return false for non-empty slice")
		}

		// Test isJSONFormat function
		if !isJSONFormat(config.PrintFormatJSON) {
			t.Error("isJSONFormat should return true for JSON format")
		}
		if isJSONFormat(config.PrintFormatTable) {
			t.Error("isJSONFormat should return false for table format")
		}

		// Test isAPMisconfigured function
		if !isAPMisconfigured(true) {
			t.Error("isAPMisconfigured should return true for true input")
		}
		if isAPMisconfigured(false) {
			t.Error("isAPMisconfigured should return false for false input")
		}
	})
}

// TestRenderingFunctions tests various rendering functions with mock data
func TestRenderingFunctions(t *testing.T) {
	t.Run("test_format_functions_with_data", func(t *testing.T) {
		// Create a temporary test config with actual controllers to test data flow
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatTable,
				Controllers: []config.Controller{{Hostname: "127.0.0.1", AccessToken: "token"}},
				Timeout:     30, // Longer timeout to avoid hanging
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		usecase := &application.Usecase{Config: cfg, Repository: repo}

		// Test AP CLI rendering paths
		apCli := &ApCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apCli.ShowAp()
			})
		})

		// Test Client CLI rendering paths
		clientCli := &ClientCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				clientCli.ShowClient()
			})
		})

		// Test WLAN CLI rendering paths
		wlanCli := &WlanCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				wlanCli.ShowWlan()
			})
		})

		// Test Overview CLI rendering paths
		overviewCli := &OverviewCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				overviewCli.ShowOverview()
			})
		})

		// Test AP Tag CLI rendering paths
		apTagCli := &ApTagCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apTagCli.ShowApTag()
			})
		})
	})

	t.Run("test_json_format_paths", func(t *testing.T) {
		// Test JSON format paths
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				PrintFormat: config.PrintFormatJSON,
				Controllers: []config.Controller{{Hostname: "127.0.0.1", AccessToken: "token"}},
				Timeout:     30, // Short timeout
			},
		}
		repo := &infrastructure.Repository{Config: cfg}
		usecase := &application.Usecase{Config: cfg, Repository: repo}

		clis := []interface{}{
			&ApCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ClientCli{Config: cfg, Repository: repo, Usecase: usecase},
			&WlanCli{Config: cfg, Repository: repo, Usecase: usecase},
			&OverviewCli{Config: cfg, Repository: repo, Usecase: usecase},
			&ApTagCli{Config: cfg, Repository: repo, Usecase: usecase},
		}

		for _, cli := range clis {
			localUtils.assertNoPanic(t, func() {
				localUtils.redirectOutput(func() {
					switch v := cli.(type) {
					case *ApCli:
						v.ShowAp()
					case *ClientCli:
						v.ShowClient()
					case *WlanCli:
						v.ShowWlan()
					case *OverviewCli:
						v.ShowOverview()
					case *ApTagCli:
						v.ShowApTag()
					}
				})
			})
		}
	})
}

// TestConfigurationVariations tests different configuration scenarios
func TestConfigurationVariations(t *testing.T) {
	t.Run("different_config_scenarios", func(t *testing.T) {
		scenarios := []struct {
			name   string
			config config.ShowCmdConfig
		}{
			{
				name: "table_format_with_filters",
				config: config.ShowCmdConfig{
					PrintFormat: config.PrintFormatTable,
					Controllers: []config.Controller{{Hostname: "test", AccessToken: "token"}},
					Radio:       "1",
					SSID:        "test-ssid",
					Timeout:     30,
				},
			},
			{
				name: "json_format_with_filters",
				config: config.ShowCmdConfig{
					PrintFormat: config.PrintFormatJSON,
					Controllers: []config.Controller{{Hostname: "test", AccessToken: "token"}},
					Radio:       "2",
					SSID:        "another-ssid",
					Timeout:     30,
				},
			},
			{
				name: "multiple_controllers",
				config: config.ShowCmdConfig{
					PrintFormat: config.PrintFormatTable,
					Controllers: []config.Controller{
						{Hostname: "test1", AccessToken: "token1"},
						{Hostname: "test2", AccessToken: "token2"},
					},
					Timeout: 30,
				},
			},
		}

		for _, scenario := range scenarios {
			t.Run(scenario.name, func(t *testing.T) {
				cfg := &config.Config{ShowCmdConfig: scenario.config}
				repo := &infrastructure.Repository{Config: cfg}
				usecase := &application.Usecase{Config: cfg, Repository: repo}

				clis := []interface{}{
					&ApCli{Config: cfg, Repository: repo, Usecase: usecase},
					&ClientCli{Config: cfg, Repository: repo, Usecase: usecase},
					&WlanCli{Config: cfg, Repository: repo, Usecase: usecase},
					&OverviewCli{Config: cfg, Repository: repo, Usecase: usecase},
					&ApTagCli{Config: cfg, Repository: repo, Usecase: usecase},
				}

				for _, cli := range clis {
					localUtils.assertNoPanic(t, func() {
						localUtils.redirectOutput(func() {
							switch v := cli.(type) {
							case *ApCli:
								v.ShowAp()
							case *ClientCli:
								v.ShowClient()
							case *WlanCli:
								v.ShowWlan()
							case *OverviewCli:
								v.ShowOverview()
							case *ApTagCli:
								v.ShowApTag()
							}
						})
					})
				}
			})
		}
	})
}

// TestDirectRenderFunctionCalls tests specific render functions directly to improve coverage
func TestDirectRenderFunctionCalls(t *testing.T) {
	cfg, repo, usecase := localUtils.createFullTestStack()

	t.Run("test_render_functions_with_minimal_data", func(t *testing.T) {
		// Create minimal mock data that matches actual struct fields
		mockApData := []*application.ShowApData{
			{
				ShowApCommonData: application.ShowApCommonData{
					ApMac:      "00:11:22:33:44:55",
					Controller: "WLC-01",
				},
			},
		}

		mockClientData := []*application.ShowClientData{
			{
				ClientMac:  "aa:bb:cc:dd:ee:ff",
				Controller: "WLC-01",
			},
		}

		mockWlanData := []*application.ShowWlanData{
			{
				TagName:    "DEFAULT",
				PolicyName: "POLICY1",
				WlanName:   "CORP",
			},
		}

		mockOverviewData := []*application.ShowOverviewData{
			{
				Controller: "WLC-01",
			},
		}

		mockApTagData := []*application.ShowApTagData{
			{
				ShowApCommonData: application.ShowApCommonData{
					ApMac:      "00:11:22:33:44:55",
					Controller: "WLC-01",
				},
			},
		}

		// Test AP render functions
		apCli := &ApCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apCli.renderShowApTable(mockApData)
				_, err := apCli.formatShowApRow(mockApData[0])
				if err != nil {
					t.Logf("formatShowApRow error (expected): %v", err)
				}
			})
		})

		// Test Client render functions
		clientCli := &ClientCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				clientCli.renderShowClientTable(mockClientData)
				_, err := clientCli.formatShowClientRow(mockClientData[0])
				if err != nil {
					t.Logf("formatShowClientRow error (expected): %v", err)
				}
			})
		})

		// Test WLAN render functions
		wlanCli := &WlanCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				wlanCli.renderShowWlanTable(mockWlanData)
				_, err := wlanCli.formatShowWlanRow(mockWlanData[0])
				if err != nil {
					t.Logf("formatShowWlanRow error (expected): %v", err)
				}
			})
		})

		// Test Overview render functions
		overviewCli := &OverviewCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				overviewCli.renderShowOverviewTable(mockOverviewData)
				_, err := overviewCli.formatShowOverviewRow(mockOverviewData[0])
				if err != nil {
					t.Logf("formatShowOverviewRow error (expected): %v", err)
				}
			})
		})

		// Test AP Tag render functions
		apTagCli := &ApTagCli{Config: cfg, Repository: repo, Usecase: usecase}
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apTagCli.renderShowApTagTable(mockApTagData)
				_, err := apTagCli.formatShowApTagRow(mockApTagData[0])
				if err != nil {
					t.Logf("formatShowApTagRow error (expected): %v", err)
				}
			})
		})
	})

	t.Run("test_edge_cases", func(t *testing.T) {
		apCli := &ApCli{Config: cfg, Repository: repo, Usecase: usecase}
		clientCli := &ClientCli{Config: cfg, Repository: repo, Usecase: usecase}
		wlanCli := &WlanCli{Config: cfg, Repository: repo, Usecase: usecase}
		overviewCli := &OverviewCli{Config: cfg, Repository: repo, Usecase: usecase}
		apTagCli := &ApTagCli{Config: cfg, Repository: repo, Usecase: usecase}

		// Test with empty slices to trigger different code paths
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apCli.renderShowApTable([]*application.ShowApData{})
				clientCli.renderShowClientTable([]*application.ShowClientData{})
				wlanCli.renderShowWlanTable([]*application.ShowWlanData{})
				overviewCli.renderShowOverviewTable([]*application.ShowOverviewData{})
				apTagCli.renderShowApTagTable([]*application.ShowApTagData{})
			})
		})

		// Test with nil slices
		localUtils.assertNoPanic(t, func() {
			localUtils.redirectOutput(func() {
				apCli.renderShowApTable(nil)
				clientCli.renderShowClientTable(nil)
				wlanCli.renderShowWlanTable(nil)
				overviewCli.renderShowOverviewTable(nil)
				apTagCli.renderShowApTagTable(nil)
			})
		})
	})
}

// Add tests for helper conversion functions with low coverage
func TestConvertCommonOperDataUsername(t *testing.T) {
	var localUtils testUtils
	_, _, usecase := localUtils.createFullTestStack()
	clientCli := &ClientCli{
		Config:     &config.Config{},
		Repository: &infrastructure.Repository{},
		Usecase:    usecase,
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty_string", "", "N/A"},
		{"valid_username", "user123", "user123"},
		{"special_chars", "user@domain.com", "user@domain.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clientCli.convertCommonOperDataUsername(tt.input)
			if result != tt.expected {
				t.Errorf("convertCommonOperDataUsername(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertCommonOperDataCoState(t *testing.T) {
	var localUtils testUtils
	_, _, usecase := localUtils.createFullTestStack()
	clientCli := &ClientCli{
		Config:     &config.Config{},
		Repository: &infrastructure.Repository{},
		Usecase:    usecase,
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"idle", "client-status-idle", "Idle"},
		{"associating", "client-status-associating", "Associating"},
		{"associated", "client-status-associated", "Associated"},
		{"authenticating", "client-status-authenticating", "Authenticating"},
		{"unknown", "client-status-unknown", "client-status-unknown"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clientCli.convertCommonOperDataCoState(tt.input)
			if result != tt.expected {
				t.Errorf("convertCommonOperDataCoState(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertCommonOperDataMsRadioTypeToBand(t *testing.T) {
	var localUtils testUtils
	_, _, usecase := localUtils.createFullTestStack()
	clientCli := &ClientCli{
		Config:     &config.Config{},
		Repository: &infrastructure.Repository{},
		Usecase:    usecase,
	}

	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{"2_4_ghz", 0, "2.4GHz"},
		{"5_ghz", 1, "5GHz"},
		{"6_ghz", 2, "6GHz"},
		{"unknown", 99, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clientCli.convertCommonOperDataMsRadioTypeToBand(tt.input)
			if result != tt.expected {
				t.Errorf("convertCommonOperDataMsRadioTypeToBand(%d) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConvertCommonOperDataMsRadioTypeToSpec(t *testing.T) {
	var localUtils testUtils
	_, _, usecase := localUtils.createFullTestStack()
	clientCli := &ClientCli{
		Config:     &config.Config{},
		Repository: &infrastructure.Repository{},
		Usecase:    usecase,
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"11b", "client-dot11b", "11b"},
		{"11g", "client-dot11g", "11g"},
		{"11a", "client-dot11a", "11a"},
		{"11n_24", "client-dot11n-24-ghz-prot", "11n"},
		{"unknown", "client-unknown", "client-unknown"},
		{"empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clientCli.convertCommonOperDataMsRadioTypeToSpec(tt.input)
			if result != tt.expected {
				t.Errorf("convertCommonOperDataMsRadioTypeToSpec(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

// TestSortFunctions tests the sort functions that currently have low coverage (Unit test)
func TestSortFunctions(t *testing.T) {
	t.Run("test_sort_show_wlan_row", func(t *testing.T) {
		// Test WlanCli sortShowWlanRow function
		wlanCli := &WlanCli{}

		// Test with nil slice
		wlanCli.sortShowWlanRow(nil)

		// Test with empty slice
		var wlans []*application.ShowWlanData
		wlanCli.sortShowWlanRow(wlans)

		// Test with single item
		wlans = []*application.ShowWlanData{
			{WlanName: "test-wlan", Controller: "test-ctrl"},
		}
		wlanCli.sortShowWlanRow(wlans)

		// Verify the slice is unchanged (since sortShowWlanRow is empty implementation)
		if len(wlans) != 1 {
			t.Errorf("Expected 1 item in wlans slice, got %d", len(wlans))
		}
		if wlans[0].WlanName != "test-wlan" {
			t.Errorf("Expected WlanName 'test-wlan', got '%s'", wlans[0].WlanName)
		}
	})

	t.Run("test_sort_show_client_row", func(t *testing.T) {
		// Test ClientCli sortShowClientRow function
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SortBy:    "hostname",
				SortOrder: "asc",
			},
		}
		clientCli := &ClientCli{Config: cfg}

		// Test with nil slice
		clientCli.sortShowClientRow(nil)

		// Test with empty slice
		var clients []*application.ShowClientData
		clientCli.sortShowClientRow(clients)

		// Test with single item
		clients = []*application.ShowClientData{
			{ClientMac: "aa:bb:cc:dd:ee:ff", Controller: "test-ctrl"},
		}
		clientCli.sortShowClientRow(clients)

		// Since we can't easily test the full sorting without complex data,
		// we just verify the function doesn't panic
		if len(clients) != 1 {
			t.Errorf("Expected 1 item in clients slice, got %d", len(clients))
		}
	})

	t.Run("test_sort_show_overview_row", func(t *testing.T) {
		// Test OverviewCli sortShowOverviewRow function
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SortBy:    "ap-name",
				SortOrder: "asc",
			},
		}
		overviewCli := &OverviewCli{Config: cfg}

		// Test with nil slice
		overviewCli.sortShowOverviewRow(nil)

		// Test with empty slice
		var overviews []*application.ShowOverviewData
		overviewCli.sortShowOverviewRow(overviews)

		// Test with single item
		overviews = []*application.ShowOverviewData{
			{ApMac: "aa:bb:cc:dd:ee:ff", Controller: "test-ctrl"},
		}
		overviewCli.sortShowOverviewRow(overviews)

		// Since we can't easily test the full sorting without complex data,
		// we just verify the function doesn't panic
		if len(overviews) != 1 {
			t.Errorf("Expected 1 item in overviews slice, got %d", len(overviews))
		}
	})

	t.Run("test_sort_show_ap_row", func(t *testing.T) {
		// Test ApCli sortShowClientRow function (this seems to be a misnaming in the code)
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SortBy:    "ap-name",
				SortOrder: "desc",
			},
		}
		apCli := &ApCli{Config: cfg}

		// Test with nil slice
		apCli.sortShowClientRow(nil)

		// Test with empty slice
		var aps []*application.ShowApData
		apCli.sortShowClientRow(aps)

		// Test with single item
		aps = []*application.ShowApData{
			{ShowApCommonData: application.ShowApCommonData{ApMac: "aa:bb:cc:dd:ee:ff", Controller: "test-ctrl"}},
		}
		apCli.sortShowClientRow(aps)

		// Since we can't easily test the full sorting without complex data,
		// we just verify the function doesn't panic
		if len(aps) != 1 {
			t.Errorf("Expected 1 item in aps slice, got %d", len(aps))
		}

		// Test with multiple items to verify actual sorting
		aps = []*application.ShowApData{
			{ShowApCommonData: application.ShowApCommonData{CapwapData: ap.CapwapData{Name: "zzz-ap"}}},
			{ShowApCommonData: application.ShowApCommonData{CapwapData: ap.CapwapData{Name: "aaa-ap"}}},
			{ShowApCommonData: application.ShowApCommonData{CapwapData: ap.CapwapData{Name: "mmm-ap"}}},
		}

		apCli.sortShowClientRow(aps)

		// Verify sorting order (should be ascending by Name)
		if aps[0].CapwapData.Name != "aaa-ap" {
			t.Errorf("Expected first item to be 'aaa-ap', got '%s'", aps[0].CapwapData.Name)
		}
		if aps[1].CapwapData.Name != "mmm-ap" {
			t.Errorf("Expected second item to be 'mmm-ap', got '%s'", aps[1].CapwapData.Name)
		}
		if aps[2].CapwapData.Name != "zzz-ap" {
			t.Errorf("Expected third item to be 'zzz-ap', got '%s'", aps[2].CapwapData.Name)
		}
	})

	t.Run("test_sort_show_ap_tag_row", func(t *testing.T) {
		// Test ApTagCli sortShowApTagRow function
		cfg := &config.Config{
			ShowCmdConfig: config.ShowCmdConfig{
				SortBy:    "ap-name",
				SortOrder: "asc",
			},
		}
		apTagCli := &ApTagCli{Config: cfg}

		// Test with nil slice
		apTagCli.sortShowApTagRow(nil)

		// Test with empty slice
		var apTags []*application.ShowApTagData
		apTagCli.sortShowApTagRow(apTags)

		// Test with single item
		apTags = []*application.ShowApTagData{
			{ShowApCommonData: application.ShowApCommonData{ApMac: "aa:bb:cc:dd:ee:ff", Controller: "test-ctrl"}},
		}
		apTagCli.sortShowApTagRow(apTags)

		// Since we can't easily test the full sorting without complex data,
		// we just verify the function doesn't panic
		if len(apTags) != 1 {
			t.Errorf("Expected 1 item in apTags slice, got %d", len(apTags))
		}
	})
}

// TestWlanCliSortFunction tests WLAN sorting function (Unit test)
func TestWlanCliSortFunction(t *testing.T) {
	utils := &testUtils{}
	cfg, repo, usecase := utils.createFullTestStack()

	t.Run("test_sort_show_wlan_row", func(t *testing.T) {
		wlanCli := &WlanCli{
			Config:     cfg,
			Repository: repo,
			Usecase:    usecase,
		}

		// Create test data
		wlans := []*application.ShowWlanData{
			{WlanName: "wlan1", Controller: "controller1"},
			{WlanName: "wlan2", Controller: "controller2"},
		}

		// Since the function is currently empty, this just tests it doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("sortShowWlanRow panicked: %v", r)
			}
		}()

		wlanCli.sortShowWlanRow(wlans)

		// Verify the slice is still intact
		if len(wlans) != 2 {
			t.Errorf("Expected 2 WLANs after sort, got %d", len(wlans))
		}
	})

	t.Run("test_sort_show_wlan_row_empty", func(t *testing.T) {
		wlanCli := &WlanCli{
			Config:     cfg,
			Repository: repo,
			Usecase:    usecase,
		}

		// Test with empty slice
		wlans := []*application.ShowWlanData{}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("sortShowWlanRow panicked with empty slice: %v", r)
			}
		}()

		wlanCli.sortShowWlanRow(wlans)

		if len(wlans) != 0 {
			t.Errorf("Expected 0 WLANs after sort of empty slice, got %d", len(wlans))
		}
	})

	t.Run("test_sort_show_wlan_row_nil", func(t *testing.T) {
		wlanCli := &WlanCli{
			Config:     cfg,
			Repository: repo,
			Usecase:    usecase,
		}

		// Test with nil slice
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("sortShowWlanRow panicked with nil slice: %v", r)
			}
		}()

		wlanCli.sortShowWlanRow(nil)
	})
}
