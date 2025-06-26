// Package cisco provides Dot11-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/dot11"
)

// Dot11-related type aliases
type (
	Dot11CfgResponse = dot11.Dot11CfgResponse
)

// GetDot11Cfg retrieves 802.11 configuration data
func GetDot11Cfg(c *Client, ctx context.Context) (*Dot11CfgResponse, error) {
	return dot11.GetDot11Cfg(c, ctx)
}
