package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	wlanLogPrefix = "wlan: "
)

// WlanRepository handles operations related to wlan data retrieval.
type WlanRepository struct {
	Config *config.Config
}

// GetWlanCfg retrieves the WLAN configuration from the specified controller using the provided wlanikey.
func (r *WlanRepository) GetWlanCfg(controller, apikey string, isSecure *bool) *cisco.WlanCfgResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(wlanLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetWlanCfg(wncClient, context.Background())
	if err != nil {
		log.Errorf(wlanLogPrefix+"%v", err)
		return nil
	}

	return resp
}
