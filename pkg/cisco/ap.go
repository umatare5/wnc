// Package cisco provides AP-related operations for Cisco WNC
package cisco

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/ap"
)

// AP-related type aliases
type (
	ApOperResponse              = ap.ApOperResponse
	ApOperCapwapDataResponse    = ap.ApOperCapwapDataResponse
	ApOperLldpNeighResponse     = ap.ApOperLldpNeighResponse
	ApOperRadioOperDataResponse = ap.ApOperRadioOperDataResponse
	ApOperOperDataResponse      = ap.ApOperOperDataResponse
	ApGlobalOperResponse        = ap.ApGlobalOperResponse
	ApCfgResponse               = ap.ApCfgResponse
	CapwapData                  = ap.CapwapData
	LldpNeigh                   = ap.LldpNeigh
	ApOperData                  = ap.ApOperData
	RadioOperData               = ap.RadioOperData
)

// GetApOper retrieves access point operational data
func GetApOper(client *Client, ctx context.Context) (*ApOperResponse, error) {
	return ap.GetApOper(client, ctx)
}

// GetApCapwapData retrieves CAPWAP data for access points
func GetApCapwapData(client *Client, ctx context.Context) (*ApOperCapwapDataResponse, error) {
	return ap.GetApCapwapData(client, ctx)
}

// GetApLldpNeigh retrieves LLDP neighbor information for access points
func GetApLldpNeigh(client *Client, ctx context.Context) (*ApOperLldpNeighResponse, error) {
	return ap.GetApLldpNeigh(client, ctx)
}

// GetApRadioOperData retrieves radio operational data for access points
func GetApRadioOperData(client *Client, ctx context.Context) (*ApOperRadioOperDataResponse, error) {
	return ap.GetApRadioOperData(client, ctx)
}

// GetApOperData retrieves operational data for access points
func GetApOperData(client *Client, ctx context.Context) (*ApOperOperDataResponse, error) {
	return ap.GetApOperData(client, ctx)
}

// GetApGlobalOper retrieves global operational data for access points
func GetApGlobalOper(client *Client, ctx context.Context) (*ApGlobalOperResponse, error) {
	return ap.GetApGlobalOper(client, ctx)
}

// GetApCfg retrieves configuration data for access points
func GetApCfg(client *Client, ctx context.Context) (*ApCfgResponse, error) {
	return ap.GetApCfg(client, ctx)
}
