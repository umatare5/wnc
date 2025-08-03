package application

import (
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// testUtils provides common utilities for application layer tests
type testUtils struct{}

func (tu *testUtils) createMockConfig() *config.Config {
	return &config.Config{
		ShowCmdConfig: config.ShowCmdConfig{
			PrintFormat: config.PrintFormatTable,
			Controllers: []config.Controller{
				{Hostname: "test-controller", AccessToken: "test-token"},
			},
			Timeout: 30,
		},
		GenerateCmdConfig: config.GenerateCmdConfig{
			Username: "testuser",
			Password: "testpass",
		},
	}
}

func (tu *testUtils) createMockRepository(cfg *config.Config) *infrastructure.Repository {
	return &infrastructure.Repository{Config: cfg}
}

func (tu *testUtils) assertNoPanic(t *testing.T, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Function panicked: %v", r)
		}
	}()
	fn()
}

// Shared test utilities instance
var testUtilsInstance = &testUtils{}
