package show

import (
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
	"github.com/umatare5/wnc/pkg/tablewriter"
)

// WlanCli struct
type WlanCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowWlan retrives the list of WLANs from the controllers
func (wc *WlanCli) ShowWlan() {
	isSecure := !wc.Config.ShowCmdConfig.AllowInsecureAccess
	wlans := wc.Usecase.InvokeWlanUsecase().ShowWlan(
		&wc.Config.ShowCmdConfig.Controllers,
		&isSecure,
	)

	if wc.Config.ShowCmdConfig.PrintFormat == config.PrintFormatJSON {
		printJson(wlans)
		return
	}

	// Skip table rendering if no data is available
	if len(wlans) == 0 {
		return
	}

	wc.renderShowWlanTable(wlans)
}

// renderShowWlanTable renders the WLAN data in a table format
func (wc *WlanCli) renderShowWlanTable(wlans []*application.ShowWlanData) {
	table := tablewriter.NewTable(os.Stdout)

	// Set table headers
	headers := wc.getShowWlanTableHeaders()
	table.Header(headers)

	// Set table rows
	wc.sortShowWlanRow(wlans)
	for _, wlan := range wlans {
		row, _ := wc.formatShowWlanRow(wlan)
		table.Append(row)
	}
	// Render the table
	_ = table.Render()
}

// getShowWlanTableHeaders returns the headers for the WLAN table
func (wc *WlanCli) getShowWlanTableHeaders() []string {
	return []string{
		"Status", "ESSID", "ID", "Profile Name", "VLAN", "Session Timeout",
		"DHCP Required", "Egress QoS", "Ingress QoS", "ATF Policies",
		"Auth Key Management", "mDNS Forwarding", "P2P Blocking", "Loadbalance",
		"Broadcast", "Tag Name", "Controller",
	}
}

// formatShowWlanRow formats a row of WLAN data
func (wc *WlanCli) formatShowWlanRow(wlan *application.ShowWlanData) ([]string, error) {
	// Check if essential data is available to avoid index out of range panics
	atfPolicyName := "N/A"
	if len(wlan.WlanPolicy.AtfPolicyMapEntries.Entries) > 0 {
		atfPolicyName = wlan.WlanPolicy.AtfPolicyMapEntries.Entries[0].AtfPolicyName
	}

	row := []string{
		wc.convertWlanPolicyStatus(wlan.WlanPolicy.Status),
		wlan.WlanName,
		fmt.Sprintf("%d", wlan.WlanCfgEntry.WlanID),
		wlan.PolicyName,
		wlan.WlanPolicy.InterfaceName,
		fmt.Sprintf("%ss", humanize.Comma(int64(wlan.WlanPolicy.WlanTimeout.SessionTimeout))),
		wc.convertDhcpParamsIsDhcpEnabled(wlan.WlanPolicy.DhcpParams.IsDhcpEnabled),
		wlan.WlanPolicy.PerSsidQos.EgressServiceName,
		wlan.WlanPolicy.PerSsidQos.IngressServiceName,
		atfPolicyName,
		wc.convertWlanCfgEntryAuthKeyMgmt(
			wlan.WlanCfgEntry.AuthKeyMgmtDot1x,
			wlan.WlanCfgEntry.AuthKeyMgmtPsk,
			wlan.WlanCfgEntry.AuthKeyMgmtSae,
		),
		wc.convertWlanCfgEntryMdnsSdMode(wlan.WlanCfgEntry.MdnsSdMode),
		wc.convertWlanCfgEntryP2PBlockAction(wlan.WlanCfgEntry.ApfVapIDData.P2PBlockAction),
		wc.convertWlanCfgEntryLoadBalance(wlan.WlanCfgEntry.LoadBalance),
		wc.convertWlanCfgEntryBroadcastSsid(wlan.WlanCfgEntry.ApfVapIDData.BroadcastSsid),
		wlan.TagName,
		wlan.Controller,
	}

	return row, nil
}

// sortShowWlanRow sorts the WLAN data by SSID name
func (wc *WlanCli) sortShowWlanRow(wlans []*application.ShowWlanData) {
}

func (wc *WlanCli) convertWlanPolicyStatus(v bool) string {
	if v {
		return "  ✅️"
	}
	return "  ❌️"
}

func (wc *WlanCli) convertDhcpParamsIsDhcpEnabled(v bool) string {
	if v {
		return "      ✅️"
	}
	return "      ⬜️"
}

func (wc *WlanCli) convertWlanCfgEntryLoadBalance(v bool) string {
	if v {
		return "     ✅️"
	}
	return "     ⬜️"
}

func (wc *WlanCli) convertWlanCfgEntryBroadcastSsid(v bool) string {
	if v {
		return "    ✅️"
	}
	return "    ⬜️"
}

func (wc *WlanCli) convertWlanCfgEntryAuthKeyMgmt(dot1x, psk, sae bool) string {
	if dot1x {
		return "Dot1x"
	}
	if psk {
		return "PSK"
	}
	if sae {
		return "SAE"
	}
	return "Unknown"
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-enum-types.yang#L3804-L3824
func (wc *WlanCli) convertWlanCfgEntryMdnsSdMode(v string) string {
	if v == "mdns-sd-drop" {
		return "Drop"
	}
	if v == "mdns-sd-gateway" {
		return "Gateway"
	}
	return "Bridging"
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-enum-types.yang#L223-L248
func (wc *WlanCli) convertWlanCfgEntryP2PBlockAction(v string) string {
	if v == "p2p-blocking-action-fwdup" {
		return "Forward-UpStream"
	}
	if v == "p2p-blocking-action-drop" {
		return "Drop"
	}
	if v == "p2p-blocking-action-allow-private-group" {
		return "Allow Private Group"
	}
	return "Disabled"
}
