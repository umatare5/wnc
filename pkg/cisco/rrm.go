// Package cisco provides RRM-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/rrm"
)

// RRM-related type aliases
type (
	RrmOperResponse        = rrm.RrmOperResponse
	RrmMeasurementResponse = rrm.RrmMeasurementResponse
	RrmGlobalOperResponse  = rrm.RrmGlobalOperResponse
	RrmCfgResponse         = rrm.RrmCfgResponse
)

// GetRrmOper retrieves RRM operational data
func GetRrmOper(c *Client, ctx context.Context) (*RrmOperResponse, error) {
	return rrm.GetRrmOper(c, ctx)
}

// GetRrmMeasurement retrieves RRM measurement data
func GetRrmMeasurement(c *Client, ctx context.Context) (*RrmMeasurementResponse, error) {
	return rrm.GetRrmMeasurement(c, ctx)
}

// GetRrmGlobalOper retrieves RRM global operational data
func GetRrmGlobalOper(c *Client, ctx context.Context) (*RrmGlobalOperResponse, error) {
	return rrm.GetRrmGlobalOper(c, ctx)
}

// GetRrmCfg retrieves RRM configuration data
func GetRrmCfg(c *Client, ctx context.Context) (*RrmCfgResponse, error) {
	return rrm.GetRrmCfg(c, ctx)
}
