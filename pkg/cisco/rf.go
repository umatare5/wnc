// Package cisco provides RF-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/rf"
)

// RF-related type aliases
type (
	RfTagsResponse = rf.RfTagsResponse
)

// GetRfTags retrieves RF tags configuration data
func GetRfTags(c *Client, ctx context.Context) (*RfTagsResponse, error) {
	return rf.GetRfTags(c, ctx)
}
