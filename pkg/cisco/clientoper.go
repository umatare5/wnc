// Package cisco provides Client-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/client"
)

// Client-related type aliases
type (
	ClientOperResponse       = client.ClientOperResponse
	ClientGlobalOperResponse = client.ClientGlobalOperResponse
)

// GetClientOper retrieves client operational data
func GetClientOper(c *Client, ctx context.Context) (*ClientOperResponse, error) {
	return client.GetClientOper(c, ctx)
}

// GetClientGlobalOper retrieves client global operational data
func GetClientGlobalOper(c *Client, ctx context.Context) (*ClientGlobalOperResponse, error) {
	return client.GetClientGlobalOper(c, ctx)
}
