package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	dot11LogPrefix = "dot11: "
)

// Dot11Repository handles operations related to dot11 data retrieval.
type Dot11Repository struct {
	Config *config.Config
}

// GetDot11Cfg retrieves the dot11 configuration from the specified controller using the provided apikey.
func (r *Dot11Repository) GetDot11Cfg(controller, apikey string, isSecure *bool) *cisco.Dot11CfgResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(dot11LogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetDot11Cfg(wncClient, context.Background())
	if err != nil {
		log.Errorf(dot11LogPrefix+"%v", err)
		return nil
	}

	return resp
}
