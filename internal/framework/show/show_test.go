package show

import (
	"testing"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TestApCli tests the ApCli structure
func TestApCli(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{Config: cfg}
	usecase := &application.Usecase{Config: cfg, Repository: repo}

	apCli := &ApCli{
		Config:     cfg,
		Repository: repo,
		Usecase:    usecase,
	}

	if apCli.Config != cfg {
		t.Error("Expected config to be set")
	}
}

// TestClientCli tests the ClientCli structure
func TestClientCli(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{Config: cfg}
	usecase := &application.Usecase{Config: cfg, Repository: repo}

	clientCli := &ClientCli{
		Config:     cfg,
		Repository: repo,
		Usecase:    usecase,
	}

	if clientCli.Config != cfg {
		t.Error("Expected config to be set")
	}
}

// TestOverviewCli tests the OverviewCli structure
func TestOverviewCli(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{Config: cfg}
	usecase := &application.Usecase{Config: cfg, Repository: repo}

	overviewCli := &OverviewCli{
		Config:     cfg,
		Repository: repo,
		Usecase:    usecase,
	}

	if overviewCli.Config != cfg {
		t.Error("Expected config to be set")
	}
}

// TestWlanCli tests the WlanCli structure
func TestWlanCli(t *testing.T) {
	cfg := &config.Config{}
	repo := &infrastructure.Repository{Config: cfg}
	usecase := &application.Usecase{Config: cfg, Repository: repo}

	wlanCli := &WlanCli{
		Config:     cfg,
		Repository: repo,
		Usecase:    usecase,
	}

	if wlanCli.Config != cfg {
		t.Error("Expected config to be set")
	}
}
