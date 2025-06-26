package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	radioLogPrefix = "radio: "
)

// RadioRepository handles operations related to radio data retrieval.
type RadioRepository struct {
	Config *config.Config
}

// GetRadioCfg retrieves configuration data for radios from the specified controller.
func (r *RadioRepository) GetRadioCfg(controller, apikey string, isSecure *bool) *cisco.RadioCfgResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(radioLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRadioCfg(wncClient, context.Background())
	if err != nil {
		log.Errorf(radioLogPrefix+"%v", err)
		return nil
	}

	return resp
}
