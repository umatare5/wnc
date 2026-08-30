package wnc

import (
	"context"
	"fmt"
	"time"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// The fields expression names the nine nodes this view renders and no others. What it leaves on
// the controller is the point: the four certificate leaves, the two external-module serial numbers
// and proxy-info, whose username and password leaves the controller does send. The serial number
// survives under device-detail because the Serial column is it.
const apViewFields = "name;wtp-mac;ip-addr;num-radio-slots;country-code;" +
	"device-detail;ap-mode-data;ap-state;ap-time-info"

// AP is one access point's identity, state, power and uplink neighbor. BootTime is how long the
// access point itself has been up and JoinTime how long the current CAPWAP association has lasted,
// so a view carrying only one of them reads a controller switchover as a fleet reboot; a zero value
// means the controller reported no instant.
type AP struct {
	Name        string
	WtpMAC      string
	EthernetMAC string
	Model       string
	Serial      string
	IPAddr      string
	SWVersion   string
	Slots       uint8
	Country     string
	Mode        string
	SubMode     string
	AdminState  string
	OperState   string
	BootTime    time.Time
	JoinTime    time.Time
	PowerType   string
	PowerMode   string
	Neighbors   []string
}

// APReads reports which secondary read failed.
type APReads struct {
	Power error
	LLDP  error
}

// APs reads the access point view. The CAPWAP collection drives the rows, so its
// failure is returned as the error and costs the controller its rows; the power pair
// and the LLDP neighbors are secondary and their failures cost only those cells.
func (c *Client) APs(ctx context.Context) ([]AP, APReads, error) {
	resp, err := c.sdk.AP().ListCAPWAPData(ctx, sdk.WithFields(apViewFields))
	if err != nil {
		return nil, APReads{}, fmt.Errorf("reading capwap-data: %w", err)
	}

	if resp == nil {
		return nil, APReads{}, nil
	}

	power, powerErr := c.apPower(ctx)
	neighbors, lldpErr := c.apNeighbors(ctx)
	reads := APReads{Power: powerErr, LLDP: lldpErr}

	aps := make([]AP, 0, len(resp.CAPWAPData))

	for _, ap := range resp.CAPWAPData {
		static := ap.DeviceDetail.StaticInfo
		row := AP{
			Name:        ap.Name,
			WtpMAC:      ap.WtpMAC,
			EthernetMAC: static.BoardData.WtpEnetMAC,
			Model:       static.ApModels.Model,
			Serial:      static.BoardData.WtpSerialNum,
			IPAddr:      ap.IPAddr,
			SWVersion:   ap.DeviceDetail.WtpVersion.SwVersion,
			Slots:       ap.NumRadioSlots,
			Country:     ap.CountryCode,
			Mode:        string(ap.ApModeData.WtpMode),
			SubMode:     ap.ApModeData.ApSubMode,
			AdminState:  string(ap.ApState.ApAdminState),
			OperState:   ap.ApState.ApOperationState,
			BootTime:    parseInstant(ap.ApTimeInfo.BootTime),
			JoinTime:    parseInstant(ap.ApTimeInfo.JoinTime),
			Neighbors:   neighbors[ap.WtpMAC],
		}

		if p, ok := power[ap.WtpMAC]; ok {
			row.PowerType, row.PowerMode = p.kind, p.mode
		}

		aps = append(aps, row)
	}

	return aps, reads, nil
}

type powerPair struct {
	kind string
	mode string
}

// apPower indexes the power pair by access point. ap-pow is a pointer, so an access point missing
// from the map is reported as unreported rather than as unpowered.
func (c *Client) apPower(ctx context.Context) (map[string]powerPair, error) {
	resp, err := c.sdk.AP().ListApOperData(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading oper-data: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string]powerPair, len(resp.OperData))

	for _, op := range resp.OperData {
		if op.ApPow == nil {
			continue
		}

		out[op.WtpMAC] = powerPair{kind: op.ApPow.PowerType, mode: op.ApPow.PowerMode}
	}

	return out, nil
}

// apNeighbors indexes the LLDP neighbors by access point. The list is keyed on the access point
// and the neighbor together, so reading the whole list and grouping it is what keeps an access
// point's row from being duplicated or dropped; keying the read by access point alone answers 404,
// because that is a partial list key.
func (c *Client) apNeighbors(ctx context.Context) (map[string][]string, error) {
	resp, err := c.sdk.AP().ListLldpNeigh(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading lldp-neigh: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string][]string, len(resp.LldpNeigh))

	for _, n := range resp.LldpNeigh {
		if label := neighborLabel(n.SystemName, n.PortID); label != "" {
			out[n.WtpMAC] = append(out[n.WtpMAC], label)
		}
	}

	return out, nil
}

// neighborLabel names one neighbor. The port is appended because a system name alone does not say
// which switchport the access point is on.
func neighborLabel(system, port string) string {
	switch {
	case system == "" && port == "":
		return ""
	case port == "":
		return system
	case system == "":
		return port
	default:
		return system + ":" + port
	}
}

// parseInstant decodes a controller timestamp. The fraction is variable width, five digits on one
// access point and six on another in the same read, and an unparseable value yields the zero time,
// which the renderer treats as unreported.
func parseInstant(s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}

	return t
}
