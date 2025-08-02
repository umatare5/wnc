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
	}{
		{
			name: "create_application_layer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}

			app := New(cfg, repo)
			if app.Config != cfg {
				t.Error("New() Config not set correctly")
			}
			if app.Repository != repo {
				t.Error("New() Repository not set correctly")
			}
		})
	}
}

// TestInvokeTokenUsecase tests token usecase invocation (Unit test)
func TestInvokeTokenUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_token_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			tokenUsecase := app.InvokeTokenUsecase()
			if tokenUsecase == nil {
				t.Error("InvokeTokenUsecase returned nil")
			}
			if tokenUsecase.Config != cfg {
				t.Error("TokenUsecase Config not set correctly")
			}
			if tokenUsecase.Repository != repo {
				t.Error("TokenUsecase Repository not set correctly")
			}
		})
	}
}

// TestGenerateBasicAuthToken tests token generation functionality (Unit test)
func TestGenerateBasicAuthToken(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		want     string
	}{
		{
			name:     "basic_auth_token",
			username: "admin",
			password: "password",
			want:     "YWRtaW46cGFzc3dvcmQ=", // base64("admin:password")
		},
		{
			name:     "empty_credentials",
			username: "",
			password: "",
			want:     "Og==", // base64(":")
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				GenerateCmdConfig: config.GenerateCmdConfig{
					Username: tt.username,
					Password: tt.password,
				},
			}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			tokenUsecase := app.InvokeTokenUsecase()
			got := tokenUsecase.GenerateBasicAuthToken()

			if got != tt.want {
				t.Errorf("GenerateBasicAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestInvokeApUsecase tests AP usecase invocation (Unit test)
func TestInvokeApUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_ap_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			apUsecase := app.InvokeApUsecase()
			if apUsecase == nil {
				t.Error("InvokeApUsecase returned nil")
			}
			if apUsecase.Config != cfg {
				t.Error("ApUsecase Config not set correctly")
			}
			if apUsecase.Repository != repo {
				t.Error("ApUsecase Repository not set correctly")
			}
		})
	}
}

// TestInvokeWlanUsecase tests WLAN usecase invocation (Unit test)
func TestInvokeWlanUsecase(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "invoke_wlan_usecase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			repo := &infrastructure.Repository{}
			app := New(cfg, repo)

			wlanUsecase := app.InvokeWlanUsecase()
			if wlanUsecase == nil {
				t.Error("InvokeWlanUsecase returned nil")
			}
			if wlanUsecase.Config != cfg {
				t.Error("WlanUsecase Config not set correctly")
			}
			if wlanUsecase.Repository != repo {
				t.Error("WlanUsecase Repository not set correctly")
			}
		})
	}
}
