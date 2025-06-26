package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	rrmLogPrefix = "rrm: "
)

// RrmRepository handles operations related to rrm data retrieval.
type RrmRepository struct {
	Config *config.Config
}

// GetRrmOper retrieves rrm operational data from the specified controller.
func (r *RrmRepository) GetRrmOper(controller, apikey string, isSecure *bool) *cisco.RrmOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(rrmLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRrmOper(wncClient, context.Background())
	if err != nil {
		log.Errorf(rrmLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetRrmMeasurement retrieves rrm measurement data from the specified controller.
func (r *RrmRepository) GetRrmMeasurement(controller, apikey string, isSecure *bool) *cisco.RrmMeasurementResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(rrmLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRrmMeasurement(wncClient, context.Background())
	if err != nil {
		log.Errorf(rrmLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetRrmGlobalOper retrieves global operational data for rrms from the specified controller.
func (r *RrmRepository) GetRrmGlobalOper(controller, apikey string, isSecure *bool) *cisco.RrmGlobalOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(rrmLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRrmGlobalOper(wncClient, context.Background())
	if err != nil {
		log.Errorf(rrmLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetRrmCfg retrieves configuration data for rrms from the specified controller.
func (r *RrmRepository) GetRrmCfg(controller, apikey string, isSecure *bool) *cisco.RrmCfgResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(rrmLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetRrmCfg(wncClient, context.Background())
	if err != nil {
		log.Errorf(rrmLogPrefix+"%v", err)
		return nil
	}

	return resp
}
