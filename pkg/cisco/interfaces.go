// Package cisco provides interfaces for WNC client operations
package cisco

import (
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=../../internal/mock/mock_interfaces.go -package=mock

// WNCClient interface defines all WNC client operations for mocking
type WNCClient interface {
	// Client operations
	GetClientOper(ctx context.Context) (interface{}, error)
	GetClientGlobalOper(ctx context.Context) (interface{}, error)

	// AP operations
	GetApOper(ctx context.Context) (interface{}, error)
	GetApCapwapData(ctx context.Context) (interface{}, error)
	GetApLldpNeigh(ctx context.Context) (interface{}, error)
	GetApRadioOperData(ctx context.Context) (interface{}, error)
	GetApOperData(ctx context.Context) (interface{}, error)
	GetApGlobalOper(ctx context.Context) (interface{}, error)
	GetApCfg(ctx context.Context) (interface{}, error)

	// WLAN operations
	GetWlanCfg(ctx context.Context) (interface{}, error)

	// RF operations
	GetRfTags(ctx context.Context) (interface{}, error)

	// RRM operations
	GetRrmOper(ctx context.Context) (interface{}, error)
	GetRrmMeasurement(ctx context.Context) (interface{}, error)
	GetRrmGlobalOper(ctx context.Context) (interface{}, error)
	GetRrmCfg(ctx context.Context) (interface{}, error)

	// DOT11 operations
	GetDot11Cfg(ctx context.Context) (interface{}, error)

	// Radio operations
	GetRadioCfg(ctx context.Context) (interface{}, error)
}

// ClientFactory interface for creating WNC clients
type ClientFactory interface {
	NewClient(controller, apikey string, isSecure *bool) (WNCClient, error)
	NewClientWithTimeout(controller, apikey string, timeout int, isSecure *bool) (WNCClient, error)
	NewClientWithOptions(controller, apikey string, options ...ClientOption) (WNCClient, error)
}
