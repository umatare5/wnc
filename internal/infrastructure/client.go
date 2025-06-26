package infrastructure

import (
	"context"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/pkg/cisco"
	"github.com/umatare5/wnc/pkg/log"
)

const (
	clientLogPrefix = "client: "
)

// ClientRepository handles operations related to client data retrieval.
type ClientRepository struct {
	Config *config.Config
}

// GetClientOper retrieves client operational data from the specified controller.
func (r *ClientRepository) GetClientOper(controller, apikey string, isSecure *bool) *cisco.ClientOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(clientLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetClientOper(wncClient, context.Background())
	if err != nil {
		log.Errorf(clientLogPrefix+"%v", err)
		return nil
	}

	return resp
}

// GetClientGlobalOper retrieves global operational data for clients from the specified controller.
func (r *ClientRepository) GetClientGlobalOper(controller, apikey string, isSecure *bool) *cisco.ClientGlobalOperResponse {
	timeout := time.Duration(r.Config.ShowCmdConfig.Timeout) * time.Second

	wncClient, err := cisco.NewClientWithTimeout(controller, apikey, timeout, isSecure)
	if err != nil {
		log.Errorf(clientLogPrefix+"failed to create client: %v", err)
		return nil
	}

	resp, err := cisco.GetClientGlobalOper(wncClient, context.Background())
	if err != nil {
		log.Errorf(clientLogPrefix+"%v", err)
		return nil
	}

	return resp
}
