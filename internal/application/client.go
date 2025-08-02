package application

import (
	"fmt"

	"github.com/umatare5/cisco-ios-xe-wireless-go/client"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// ClientUsecase handles client-related operations
type ClientUsecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// ShowClientData holds merged client data from various sources
type ShowClientData struct {
	ClientMac      string                `json:"client-mac"`
	Controller     string                `json:"controller"`
	CommonOperData client.CommonOperData `json:"common-oper-client"`
	Dot11OperData  client.Dot11OperData  `json:"dot11-oper-client"`
	TrafficStats   client.TrafficStats   `json:"traffic-stats"`
	SisfDbMac      client.SisfDbMac      `json:"sisf-db-mac"`
	DcInfo         client.DcInfo         `json:"dc-info"`
}

// ShowClient retrieves and merges client data from multiple controllers
func (u *ClientUsecase) ShowClient(controllers *[]config.Controller, isSecure *bool) []*ShowClientData {
	data := make([]*ShowClientData, 0)

	// Return empty slice if repository is nil
	if u.Repository == nil {
		return data
	}

	// Return empty slice if controllers is nil
	if controllers == nil {
		return data
	}

	for _, controller := range *controllers {
		result := u.Repository.InvokeClientRepository().GetClientOper(controller.Hostname, controller.AccessToken, isSecure)
		if result == nil {
			// Skip this controller if authentication failed or other error occurred
			continue
		}

		for _, client := range result.CiscoIOSXEWirelessClientOperClientOperData.CommonOperData {
			var merged ShowClientData
			merged.ClientMac = client.ClientMac
			merged.Controller = controller.Hostname
			merged.CommonOperData = client

			// Search Dot11OperData
			for _, d := range result.CiscoIOSXEWirelessClientOperClientOperData.Dot11OperData {
				if d.MsMacAddress == client.ClientMac {
					merged.Dot11OperData = d
					break
				}
			}

			// Search TrafficStats
			for _, d := range result.CiscoIOSXEWirelessClientOperClientOperData.TrafficStats {
				if d.MsMacAddress == client.ClientMac {
					merged.TrafficStats = d
					break
				}
			}

			// Search SisfDbMac
			for _, d := range result.CiscoIOSXEWirelessClientOperClientOperData.SisfDbMac {
				if d.MacAddr == client.ClientMac {
					merged.SisfDbMac = d
					break
				}
			}

			// Search DcInfo
			for _, d := range result.CiscoIOSXEWirelessClientOperClientOperData.DcInfo {
				if d.ClientMac == client.ClientMac {
					merged.DcInfo = d
					break
				}
			}

			data = append(data, &merged)
		}
	}

	return u.filterBySSID(u.filterByRadio(data))
}

func (u *ClientUsecase) filterBySSID(clients []*ShowClientData) []*ShowClientData {
	// Return all clients if config is nil
	if u.Config == nil {
		return clients
	}
	filter := u.Config.ShowCmdConfig.SSID
	if filter == "" {
		return clients
	}
	filteredClients := []*ShowClientData{}
	for _, client := range clients {
		if client.Dot11OperData.VapSsid == filter {
			filteredClients = append(filteredClients, client)
		}
	}
	return filteredClients
}

func (u *ClientUsecase) filterByRadio(clients []*ShowClientData) []*ShowClientData {
	// Return all clients if config is nil
	if u.Config == nil {
		return clients
	}
	filter := u.Config.ShowCmdConfig.Radio
	if filter == "" {
		return clients
	}
	filteredClients := []*ShowClientData{}
	for _, client := range clients {
		if fmt.Sprintf("%d", client.CommonOperData.MsApSlotID) == filter {
			filteredClients = append(filteredClients, client)
		}
	}
	return filteredClients
}
