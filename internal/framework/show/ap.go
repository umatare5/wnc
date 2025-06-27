package show

import (
	"fmt"
	"os"
	"sort"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
	"github.com/umatare5/wnc/pkg/tablewriter"
)

// ApCli struct
type ApCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowAp retrieves the list of access points from the controllers
func (ac *ApCli) ShowAp() {
	isSecure := !ac.Config.ShowCmdConfig.AllowInsecureAccess
	aps := ac.Usecase.InvokeApUsecase().ShowAp(
		&ac.Config.ShowCmdConfig.Controllers,
		&isSecure,
	)

	if ac.Config.ShowCmdConfig.PrintFormat == config.PrintFormatJSON {
		printJson(aps)
		return
	}

	// Skip table rendering if no data is available
	if len(aps) == 0 {
		return
	}

	ac.renderShowApTable(aps)
}

// renderShowApTable renders the access point data in a table format
func (ac *ApCli) renderShowApTable(aps []*application.ShowApData) {
	table := tablewriter.NewTable(os.Stdout)

	// Set table headers
	headers := ac.getShowApTableHeaders()
	table.Header(headers)
	// Set table rows
	ac.sortShowClientRow(aps)
	for _, ap := range aps {
		row, _ := ac.formatShowApRow(ap)
		table.Append(row)
	}
	// Render the table
	_ = table.Render()
}

func (ac *ApCli) getShowApTableHeaders() []string {
	return []string{
		"AP Name", "Slots", "Model", "Serial", "Ethernet MAC", "Radio MAC",
		"Country Code", "Domain", "IP Address", "OS Version",
		"State", "LLDP Neighbor", "Power Type", "Power Mode", "Controller",
	}
}

func (ac *ApCli) formatShowApRow(ap *application.ShowApData) ([]string, error) {
	row := []string{
		ap.CapwapData.Name,
		fmt.Sprintf("%d", ap.CapwapData.NumRadioSlots),
		ap.CapwapData.DeviceDetail.StaticInfo.ApModels.Model,
		ap.CapwapData.DeviceDetail.StaticInfo.BoardData.WtpSerialNum,
		ap.CapwapData.DeviceDetail.StaticInfo.BoardData.WtpEnetMac,
		ap.CapwapData.WtpMac,
		ap.CapwapData.CountryCode,
		ap.CapwapData.RegDomain,
		ap.CapwapData.IPAddr,
		ap.CapwapData.DeviceDetail.WtpVersion.SwVersion,
		ap.CapwapData.ApState.ApOperationState,
		ap.LLDPnei.SystemName + " " + ap.LLDPnei.PortID,
		ac.convertApOperDataApPowPowerType(ap.ApOperData.ApPow.PowerType),
		ac.convertApOperDataApPowPowerMode(ap.ApOperData.ApPow.PowerMode),
		ap.Controller,
	}
	return row, nil
}

// sortShowClientRow sorts the access point data by name
func (ac *ApCli) sortShowClientRow(aps []*application.ShowApData) {
	// Sort the access points by name
	sort.Slice(aps, func(i, j int) bool {
		return aps[i].CapwapData.Name < aps[j].CapwapData.Name
	})
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-enum-types.yang#L2445-L2501
func (ac *ApCli) convertApOperDataApPowPowerMode(v string) string {
	if v == "dot11-default-low-pwr" {
		return "Default Low"
	}
	if v == "dot11-set-low-pwr" {
		return "Low"
	}
	if v == "dot11-set-15-4-pwr" {
		return "15.4W"
	}
	if v == "dot11-set-16-8-pwr" {
		return "16.8W"
	}
	if v == "dot11-default-high-pwr" {
		return "Default High"
	}
	if v == "dot11-set-high-pwr" {
		return "High"
	}
	if v == "dot11-set-no-pwr" {
		return "No Power"
	}
	if v == "dot11-set-25-5-pwr" {
		return "25.5W"
	}
	if v == "unknown-pwr" {
		return "Unknown"
	}
	return v
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-enum-types.yang#L2503-L2538
func (ac *ApCli) convertApOperDataApPowPowerType(v string) string {
	if v == "pwr-src-brick-old" {
		return "Power Supply"
	}
	if v == "pwr-src-brick-new" {
		return "Power Supply"
	}
	if v == "pwr-src-inj" {
		return "PoE Injector"
	}
	if v == "pwr-src-poe-lgcy" {
		return "Legacy PoE"
	}
	if v == "pwr-src-poe-plus" {
		return "Advanced PoE"
	}
	if v == "pwr-src-unknown" {
		return "Unknown"
	}
	return v
}
