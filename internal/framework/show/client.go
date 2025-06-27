package show

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
	"github.com/umatare5/wnc/pkg/tablewriter"
)

// ClientCli struct
type ClientCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowClient retrieves the list of clients from the controllers
func (cc *ClientCli) ShowClient() {
	isSecure := !cc.Config.ShowCmdConfig.AllowInsecureAccess
	res := cc.Usecase.InvokeClientUsecase().ShowClient(
		&cc.Config.ShowCmdConfig.Controllers,
		&isSecure,
	)

	if cc.Config.ShowCmdConfig.PrintFormat == config.PrintFormatJSON {
		printJson(res)
		return
	}

	// Skip table rendering if no data is available
	if len(res) == 0 {
		return
	}

	cc.renderShowClientTable(res)
}

// renderShowClientTable renders the client data in a table format
func (cc *ClientCli) renderShowClientTable(clients []*application.ShowClientData) {
	table := tablewriter.NewTable(os.Stdout)

	// Set table headers
	headers := cc.getShowClientTableHeaders()
	table.Header(headers)

	// Set table rows
	cc.sortShowClientRow(clients)
	for _, client := range clients {
		row, _ := cc.formatShowClientRow(client)
		table.Append(row)
	}
	// Render the table
	_ = table.Render()
}

// getShowClientTableHeaders returns the headers for the client table
func (cc *ClientCli) getShowClientTableHeaders() []string {
	return []string{
		config.ShowClientHeaderMacAddress,
		config.ShowClientHeaderIP,
		config.ShowClientHeaderHostname,
		config.ShowClientHeaderUsername,
		config.ShowClientHeaderSSID,
		config.ShowClientHeaderProtocol,
		config.ShowClientHeaderBand,
		config.ShowClientHeaderState,
		config.ShowClientHeaderThroughput,
		config.ShowClientHeaderRSSI,
		config.ShowClientHeaderSNR,
		config.ShowClientHeaderStream,
		config.ShowClientHeaderRxTraffic,
		config.ShowClientHeaderTxTraffic,
		config.ShowCommonHeaderApName,
		config.ShowCommonHeaderController,
	}
}

// formatShowClientRow formats a single client's data into a table row
func (cc *ClientCli) formatShowClientRow(client *application.ShowClientData) ([]string, error) {
	bytesRx, err := strconv.ParseInt(client.TrafficStats.BytesRx, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Rx bytes: %v", err)
	}
	bytesTx, err := strconv.ParseInt(client.TrafficStats.BytesTx, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid Tx bytes: %v", err)
	}

	return []string{
		client.ClientMac,
		client.SisfDbMac.Ipv4Binding.IPKey.IPAddr,
		client.DcInfo.DeviceName,
		cc.convertCommonOperDataUsername(client.CommonOperData.Username),
		client.Dot11OperData.VapSsid,
		cc.convertCommonOperDataMsRadioTypeToSpec(client.CommonOperData.MsRadioType),
		cc.convertCommonOperDataMsRadioTypeToBand(client.CommonOperData.MsApSlotID),
		cc.convertCommonOperDataCoState(client.CommonOperData.CoState),
		fmt.Sprintf("%d Mbps", client.TrafficStats.Speed),
		fmt.Sprintf("%d dBm", client.TrafficStats.MostRecentRssi),
		fmt.Sprintf("%d dB", client.TrafficStats.MostRecentSnr),
		fmt.Sprintf("%d Streams", client.TrafficStats.SpatialStream),
		fmt.Sprintf("%s KB", humanize.Comma(bytesRx/1024)),
		fmt.Sprintf("%s KB", humanize.Comma(bytesTx/1024)),
		client.CommonOperData.ApName,
		client.Controller,
	}, nil
}

func (cc *ClientCli) sortShowClientRow(clients []*application.ShowClientData) {
	sort.Slice(clients, func(i, j int) bool {
		sortBy := cc.Config.ShowCmdConfig.SortBy
		sortOrder := cc.Config.ShowCmdConfig.SortOrder

		var data bool
		switch sortBy {
		case config.ShowClientHeaderIP:
			data = clients[i].SisfDbMac.Ipv4Binding.IPKey.IPAddr < clients[j].SisfDbMac.Ipv4Binding.IPKey.IPAddr
		case config.ShowClientHeaderHostname:
			data = clients[i].DcInfo.DeviceName < clients[j].DcInfo.DeviceName
		case config.ShowClientHeaderThroughput:
			data = clients[i].TrafficStats.Speed < clients[j].TrafficStats.Speed
		case config.ShowClientHeaderRSSI:
			data = clients[i].TrafficStats.MostRecentRssi < clients[j].TrafficStats.MostRecentRssi
		case config.ShowClientHeaderSNR:
			data = clients[i].TrafficStats.MostRecentSnr < clients[j].TrafficStats.MostRecentSnr
		case config.ShowClientHeaderTxTraffic:
			bytesTx1, _ := strconv.ParseInt(clients[i].TrafficStats.BytesTx, 10, 64)
			bytesTx2, _ := strconv.ParseInt(clients[j].TrafficStats.BytesTx, 10, 64)
			data = bytesTx1 < bytesTx2
		case config.ShowClientHeaderRxTraffic:
			bytesRx1, _ := strconv.ParseInt(clients[i].TrafficStats.BytesRx, 10, 64)
			bytesRx2, _ := strconv.ParseInt(clients[j].TrafficStats.BytesRx, 10, 64)
			data = bytesRx1 < bytesRx2
		default:
			data = false
		}

		if sortOrder == config.OrderByDescending {
			return !data
		}
		return data
	})
}

func (cc *ClientCli) convertCommonOperDataUsername(v string) string {
	if v == "" {
		return "N/A"
	}
	return v
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-client-types.yang#L136-L211
func (cc *ClientCli) convertCommonOperDataCoState(v string) string {
	if v == "client-status-idle" {
		return "Idle"
	}
	if v == "client-status-associating" {
		return "Associating"
	}
	if v == "client-status-associated" {
		return "Associated"
	}
	if v == "client-status-authenticating" {
		return "Authenticating"
	}
	if v == "client-status-authenticated" {
		return "Authenticated"
	}
	if v == "client-status-mobility-discovery" {
		return "Mobility Discovery"
	}
	if v == "client-status-mobility-complete" {
		return "Mobility Completed"
	}
	if v == "client-status-ip-learning" {
		return "IP Learning"
	}
	if v == "client-status-ip-learn-complete" {
		return "IP Learned"
	}
	if v == "client-status-webauth-required" {
		return "WebAuth Required"
	}
	if v == "client-status-static-ip-anchor-discovery" {
		return "Anchor Discovery"
	}
	if v == "client-status-run" {
		return "Run"
	}
	if v == "client-status-delete-in-progress" {
		return "In Progress"
	}
	if v == "client-status-deleted" {
		return "Deleted"
	}
	return v
}

func (cc *ClientCli) convertCommonOperDataMsRadioTypeToBand(v int) string {
	if v == 0 {
		return "2.4GHz"
	}
	if v == 1 {
		return "5GHz"
	}
	if v == 2 {
		return "6GHz"
	}
	return "Unknown"
}

// Reference: https://github.com/YangModels/yang/blob/d0fc4d40ae414990cc0858c60446b67069b95173/vendor/cisco/xe/17121/Cisco-IOS-XE-wireless-client-types.yang#L240-L310
func (cc *ClientCli) convertCommonOperDataMsRadioTypeToSpec(v string) string {
	if v == "client-dot11b" {
		return "11b"
	}
	if v == "client-dot11g" {
		return "11g"
	}
	if v == "client-dot11a" {
		return "11a"
	}
	if v == "client-dot11n-24-ghz-prot" {
		return "11n"
	}
	if v == "client-dot11n-5-ghz-prot" {
		return "11n"
	}
	if v == "client-dot11ac" {
		return "11ac"
	}
	if v == "client-phy-type-notappl" {
		return "Not Applicable"
	}
	if v == "client-ethernet" {
		return "Ethernet"
	}
	if v == "client-dot11ax-5ghz-prot" {
		return "dot11ax"
	}
	if v == "client-dot11ax-24ghz-prot" {
		return "dot11ax"
	}
	if v == "client-802-3" {
		return "802.3"
	}
	if v == "client-dot11ax-6ghz-prot" {
		return "dot11ax"
	}
	if v == "client-unknown-prot" {
		return "Unknown"
	}
	return v
}
