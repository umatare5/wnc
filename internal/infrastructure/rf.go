package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	rfLogPrefix = "rf: "
)

// RfRepository handles operations related to rf data retrieval.
type RfRepository struct {
	Config *config.Config
}

// GetRfCfg retrieves configuration data for rf from the specified controller.
func (r *RfRepository) GetRfTags(controller, apikey string, isSecure *bool) *cisco.RfTagsResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(rfLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRfTags(wncClient, context.Background())
	if err != nil {
		log.Errorf(rfLogPrefix+"%v", err)
		return nil
	}

	return resp
}
