package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	apLogPrefix = "ap: "
)

// ApRepository handles operations related to ap data retrieval.
type ApRepository struct {
	Config *config.Config
}

// GetApOper retrieves ap operational data from the specified controller.
func (r *ApRepository) GetApOper(controller, apikey string, isSecure *bool) *cisco.ApOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApOper(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApCapwapData retrieves capwap data from the specified controller.
func (r *ApRepository) GetApCapwapData(controller, apikey string, isSecure *bool) *cisco.ApOperCapwapDataResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApCapwapData(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApLldpNeigh retrieves LLDP neighbor data from the specified controller.
func (r *ApRepository) GetApLldpNeigh(controller, apikey string, isSecure *bool) *cisco.ApOperLldpNeighResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApLldpNeigh(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApRadioOperData retrieves radio operational data from the specified controller.
func (r *ApRepository) GetApRadioOperData(controller, apikey string, isSecure *bool) *cisco.ApOperRadioOperDataResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApRadioOperData(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApOperData retrieves operational data for aps from the specified controller.
func (r *ApRepository) GetApOperData(controller, apikey string, isSecure *bool) *cisco.ApOperOperDataResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApOperData(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApGlobalOper retrieves global operational data for aps from the specified controller.
func (r *ApRepository) GetApGlobalOper(controller, apikey string, isSecure *bool) *cisco.ApGlobalOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApGlobalOper(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetApCfg retrieves configuration data for aps from the specified controller.
func (r *ApRepository) GetApCfg(controller, apikey string, isSecure *bool) *cisco.ApCfgResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	client, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(apLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetApCfg(client, context.Background())
	if err != nil {
		log.Errorf(apLogPrefix+"%v", err)
		return nil
	}

	return resp
}
