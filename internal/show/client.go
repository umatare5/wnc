package show

import (
	"context"
	"fmt"
	"time"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// ClientFilter narrows the rows. An empty field selects everything.
type ClientFilter struct {
	Band   string
	SSID   string
	APName string
}

// needsDot11 reports whether a filter reads from the 802.11 collection, which promotes that
// collection to primary: zero matches would claim the fleet holds no such client when the truth
// is that nothing was read.
func (f ClientFilter) needsDot11() bool {
	return f.Band != "" || f.SSID != ""
}

// ClientRow is one row of show client.
type ClientRow struct {
	MAC        *string `json:"mac,omitzero"`
	IPv4       *string `json:"ipv4,omitzero"`
	IPv6       *string `json:"ipv6,omitzero"`
	Device     *string `json:"device,omitzero"`
	Username   *string `json:"username,omitzero"`
	SSID       *string `json:"ssid,omitzero"`
	APName     *string `json:"ap_name,omitzero"`
	Slot       *int    `json:"slot,omitzero"`
	Band       *string `json:"band,omitzero"`
	Protocol   *string `json:"protocol,omitzero"`
	Channel    *int    `json:"channel,omitzero"`
	State      *string `json:"state,omitzero"`
	RSSI       *int    `json:"rssi_dbm,omitzero"`
	SNR        *int    `json:"snr_db,omitzero"`
	Speed      *int    `json:"speed_mbps,omitzero"`
	Streams    *int    `json:"spatial_streams,omitzero"`
	Assoc      *int64  `json:"assoc_seconds,omitzero"`
	RxBytes    *uint64 `json:"rx_bytes,omitzero"`
	TxBytes    *uint64 `json:"tx_bytes,omitzero"`
	Controller string  `json:"controller"`
}

// ClientColumns describes the client view. Device is the controller's own classification label
// and not a hostname: the leaf that does carry one sits on a list this SDK cannot reach.
func ClientColumns() []render.Column[ClientRow] {
	return []render.Column[ClientRow]{
		{Key: DefaultSortMAC, Header: "MAC", Cell: func(r ClientRow) string { return render.StrPtr(r.MAC) }},
		{Key: "ipv4", Header: "IPv4", Cell: func(r ClientRow) string { return render.StrPtr(r.IPv4) }},
		{Key: "ipv6", Header: "IPv6", Cell: func(r ClientRow) string { return render.StrPtr(r.IPv6) }},
		{Key: "device", Header: "Device", Cell: func(r ClientRow) string { return render.StrPtr(r.Device) }},
		{Key: "username", Header: "Username", Cell: func(r ClientRow) string { return render.StrPtr(r.Username) }},
		{Key: keySSID, Header: "SSID", Cell: func(r ClientRow) string { return render.StrPtr(r.SSID) }},
		{Key: keyAPName, Header: headAPName, Cell: func(r ClientRow) string { return render.StrPtr(r.APName) }},
		{
			Key: keySlot, Header: "Slot",
			Cell: func(r ClientRow) string { return render.IntPtr(r.Slot) },
			Sort: func(r ClientRow) any { return render.SortValue(r.Slot) },
		},
		{Key: keyBand, Header: "Band", Cell: func(r ClientRow) string { return render.StrPtr(r.Band) }},
		{Key: "protocol", Header: "Protocol", Cell: func(r ClientRow) string { return render.StrPtr(r.Protocol) }},
		{
			Key: keyChannel, Header: "Channel",
			Cell: func(r ClientRow) string { return render.UnitPtr(r.Channel, "ch") },
			Sort: func(r ClientRow) any { return render.SortValue(r.Channel) },
		},
		{Key: keyState, Header: "State", Cell: func(r ClientRow) string { return render.StrPtr(r.State) }},
		{
			Key: "rssi_dbm", Header: "RSSI",
			Cell: func(r ClientRow) string { return render.UnitPtr(r.RSSI, "dBm") },
			Sort: func(r ClientRow) any { return render.SortValue(r.RSSI) },
		},
		{
			Key: "snr_db", Header: "SNR",
			Cell: func(r ClientRow) string { return render.UnitPtr(r.SNR, "dB") },
			Sort: func(r ClientRow) any { return render.SortValue(r.SNR) },
		},
		{
			Key: "speed_mbps", Header: "Rate",
			Cell: func(r ClientRow) string { return render.UnitPtr(r.Speed, "Mbps") },
			Sort: func(r ClientRow) any { return render.SortValue(r.Speed) },
		},
		{
			Key: "spatial_streams", Header: "Streams",
			Cell: func(r ClientRow) string { return render.UnitPtr(r.Streams, "ss") },
			Sort: func(r ClientRow) any { return render.SortValue(r.Streams) },
		},
		{
			Key: "assoc_seconds", Header: "Assoc",
			Cell: func(r ClientRow) string { return render.Duration(r.Assoc) },
			Sort: func(r ClientRow) any { return render.SortValue(r.Assoc) },
		},
		{
			Key: "rx_bytes", Header: "Rx",
			Cell: func(r ClientRow) string { return render.IEC(r.RxBytes) },
			Sort: func(r ClientRow) any { return render.SortValue(r.RxBytes) },
		},
		{
			Key: "tx_bytes", Header: "Tx",
			Cell: func(r ClientRow) string { return render.IEC(r.TxBytes) },
			Sort: func(r ClientRow) any { return render.SortValue(r.TxBytes) },
		},
		{
			Key:    keyController,
			Header: headController,
			Cell:   func(r ClientRow) string { return render.Str(r.Controller) },
		},
	}
}

// FetchClients builds the fetcher for one filter.
func FetchClients(filter ClientFilter) Fetcher[ClientRow] {
	return func(ctx context.Context, c *wnc.Client, t config.Target, rep *Reporter) ([]ClientRow, error) {
		clients, reads, err := c.Clients(ctx)
		if err != nil {
			return nil, err
		}

		if reads.Dot11 != nil && filter.needsDot11() {
			return nil, fmt.Errorf("a filter reads from this collection: %w", reads.Dot11)
		}

		rep.Degraded("dot11-oper-data", reads.Dot11)
		rep.Degraded("traffic-stats", reads.Stats)
		rep.Degraded("sisf-db-mac", reads.SISF)
		rep.Degraded("dc-info", reads.DC)

		return clientRows(clients, filter, t, rep), nil
	}
}

// clientRows filters and converts the joined clients.
func clientRows(
	clients []wnc.WirelessClient, filter ClientFilter, t config.Target, rep *Reporter,
) []ClientRow {
	now := time.Now()
	rows := make([]ClientRow, 0, len(clients))
	unreported := 0

	for _, cl := range clients {
		band := showClientBand(cl.Band)

		if filter.Band != "" && band == "" {
			unreported++

			continue
		}

		if !matches(filter, cl, band) {
			continue
		}

		rows = append(rows, ClientRow{
			MAC:      optional(cl.MAC),
			IPv4:     optional(usableIPv4(cl.IPv4)),
			IPv6:     optional(cl.IPv6),
			Device:   optional(cl.Device),
			Username: optional(cl.Username),
			SSID:     optional(cl.SSID),
			APName:   optional(cl.APName),
			// A client is on a real radio slot and zero is one of them, so zero is a
			// reading here and the slot is not zero-sentinelled.
			Slot:     ptr(cl.Slot),
			Band:     optional(band),
			Protocol: optional(showClientPHY(cl.PHY)),
			// A channel is never zero on an associated client, so zero is an omission.
			Channel: zeroAbsent(cl.Channel),
			State:   optional(showClientState(cl.State)),
			// An RSSI of exactly 0 dBm is not a reading a controller reports; an SNR of
			// 0 dB is, so the first is zero-sentinelled and the second arrives already
			// nil where the traffic counters carried no row for this client.
			RSSI:       zeroAbsent(cl.RSSI),
			SNR:        cl.SNR,
			Speed:      zeroAbsent(cl.Speed),
			Streams:    zeroAbsent(cl.Streams),
			Assoc:      render.SecondsSince(now, cl.AssocTime),
			RxBytes:    cl.RxBytes,
			TxBytes:    cl.TxBytes,
			Controller: t.Name,
		})
	}

	rep.Excluded(unreported, "the band was not reported for them")

	return rows
}

func matches(filter ClientFilter, cl wnc.WirelessClient, band string) bool {
	switch {
	case filter.Band != "" && band != filter.Band:
		return false
	case filter.SSID != "" && cl.SSID != filter.SSID:
		return false
	case filter.APName != "" && cl.APName != filter.APName:
		return false
	}

	return true
}

// usableIPv4 rejects the two values that mean "no address". The binding table reports
// an unassigned client as the empty string on some releases and as the all-zeros
// address on others.
func usableIPv4(s string) string {
	if s == "0.0.0.0" {
		return ""
	}

	return s
}
