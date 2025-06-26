// Package cisco provides Radio-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/radio"
)

// Radio-related type aliases
type (
	RadioCfgResponse = radio.RadioCfgResponse
)

// GetRadioCfg retrieves radio configuration data
func GetRadioCfg(c *Client, ctx context.Context) (*RadioCfgResponse, error) {
	return radio.GetRadioCfg(c, ctx)
}
