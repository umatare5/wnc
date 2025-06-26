// Package cisco provides WLAN-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/wlan"
)

// WLAN-related type aliases
type (
	WlanCfgResponse = wlan.WlanCfgResponse
)

// GetWlanCfg retrieves WLAN configuration data
func GetWlanCfg(c *Client, ctx context.Context) (*WlanCfgResponse, error) {
	return wlan.GetWlanCfg(c, ctx)
}
