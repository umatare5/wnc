package wnc

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
	"time"
)

// WirelessClient is one associated client, joined across the five collections that
// describe it. Every field carries what the controller sent; the zero values are
// resolved into absences one layer up, where the rule per field is documented.
type WirelessClient struct {
	MAC      string
	APName   string
	Slot     int
	State    string
	Username string
	IPv4     string
	IPv6     string
	Device   string
	SSID     string
	Band     string
	PHY      string
	Channel  int
	RSSI     int
	// SNR is a pointer where its three siblings are not, because 0 dB is a real margin and their
	// zeros are not readings. A client the traffic counters carry no row for would otherwise reach
	// the row builder as a fabricated 0 dB.
	SNR       *int
	Speed     int
	Streams   int
	AssocTime time.Time
	RxBytes   *uint64
	TxBytes   *uint64

	// HasDot11 records whether the 802.11 collection supplied this client. The SSID
	// and band filters read from it, so a filter cannot be honestly applied without it.
	HasDot11 bool
}

// ClientReads reports which of the secondary reads failed. A filter whose source is
// among them must not be applied: reporting zero matches would claim the fleet has no
// such client, when the truth is that nothing was read.
type ClientReads struct {
	Dot11 error
	Stats error
	SISF  error
	DC    error
}

// Clients reads the client view. The common collection drives the rows.
func (c *Client) Clients(ctx context.Context) ([]WirelessClient, ClientReads, error) {
	resp, err := c.sdk.Client().ListCommonInfo(ctx)
	if err != nil {
		return nil, ClientReads{}, fmt.Errorf("reading common-oper-data: %w", err)
	}

	if resp == nil {
		return nil, ClientReads{}, nil
	}

	dot11, dot11Err := c.clientDot11(ctx)
	stats, statsErr := c.clientStats(ctx)
	addrs, sisfErr := c.clientAddresses(ctx)
	devices, dcErr := c.clientDevices(ctx)

	reads := ClientReads{Dot11: dot11Err, Stats: statsErr, SISF: sisfErr, DC: dcErr}

	out := make([]WirelessClient, 0, len(resp.CommonOperData))

	for _, cl := range resp.CommonOperData {
		row := WirelessClient{
			MAC:      cl.ClientMAC,
			APName:   cl.ApName,
			Slot:     cl.MsApSlotID,
			State:    cl.CoState,
			Username: cl.Username,
			Device:   devices[cl.ClientMAC],
		}

		if d, ok := dot11[cl.ClientMAC]; ok {
			row.HasDot11 = true
			row.SSID, row.Band, row.PHY = d.ssid, d.band, d.phy
			row.Channel, row.AssocTime = d.channel, d.assocTime
		}

		if s, ok := stats[cl.ClientMAC]; ok {
			row.RSSI, row.Speed, row.Streams = s.rssi, s.speed, s.streams
			row.SNR, row.RxBytes, row.TxBytes = &s.snr, s.rx, s.tx
		}

		if a, ok := addrs[cl.ClientMAC]; ok {
			row.IPv4, row.IPv6 = a.v4, a.v6
		}

		out = append(out, row)
	}

	return out, reads, nil
}

type dot11Facts struct {
	ssid      string
	band      string
	phy       string
	channel   int
	assocTime time.Time
}

// clientDot11 indexes the 802.11 facts by client. The band comes from radio-type, whose typedef
// genuinely denotes a band; the similarly named ms-radio-type leaf on the common collection is
// typed with the PHY-generation typedef instead and restates the protocol.
func (c *Client) clientDot11(ctx context.Context) (map[string]dot11Facts, error) {
	resp, err := c.sdk.Client().ListDot11Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading dot11-oper-data: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string]dot11Facts, len(resp.Dot11OperData))

	for _, d := range resp.Dot11OperData {
		out[d.MsMACAddress] = dot11Facts{
			ssid:      d.VapSsid,
			band:      string(d.RadioType),
			phy:       string(d.EwlcMsPhyType),
			channel:   d.CurrentChannel,
			assocTime: d.MsAssocTime,
		}
	}

	return out, nil
}

type statFacts struct {
	rssi    int
	snr     int
	speed   int
	streams int
	rx      *uint64
	tx      *uint64
}

// clientStats indexes the traffic counters by client. The octet counts arrive as JSON
// strings because they are 64-bit, so an unparseable one stays absent rather than
// becoming a zero.
func (c *Client) clientStats(ctx context.Context) (map[string]statFacts, error) {
	resp, err := c.sdk.Client().ListTrafficStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading traffic-stats: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string]statFacts, len(resp.TrafficStats))

	for _, s := range resp.TrafficStats {
		out[s.MsMACAddress] = statFacts{
			rssi:    s.MostRecentRSSI,
			snr:     s.MostRecentSNR,
			speed:   s.Speed,
			streams: s.SpatialStream,
			rx:      parseCounter(s.BytesRx),
			tx:      parseCounter(s.BytesTx),
		}
	}

	return out, nil
}

type addrFacts struct {
	v4 string
	v6 string
}

func (c *Client) clientAddresses(ctx context.Context) (map[string]addrFacts, error) {
	resp, err := c.sdk.Client().ListSISFDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading sisf-db-mac: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string]addrFacts, len(resp.SisfDBMac))

	for _, b := range resp.SisfDBMac {
		v6 := make([]string, 0, len(b.Ipv6Binding))
		for _, entry := range b.Ipv6Binding {
			v6 = append(v6, entry.Ipv6BindingIPKey.IPAddr)
		}

		out[b.MACAddr] = addrFacts{v4: b.Ipv4Binding.IPKey.IPAddr, v6: pickIPv6(v6)}
	}

	return out, nil
}

func (c *Client) clientDevices(ctx context.Context) (map[string]string, error) {
	resp, err := c.sdk.Client().ListDCInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading dc-info: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[string]string, len(resp.DcInfo))
	for _, d := range resp.DcInfo {
		out[d.ClientMAC] = d.DeviceName
	}

	return out, nil
}

// pickIPv6 chooses one address from a binding list the schema bounds at eight entries. Link-local
// addresses are dropped because every client has one, and the rest are compared as parsed values
// rather than as text, because the controller compresses some entries with "::" and not others.
func pickIPv6(addrs []string) string {
	var best netip.Addr

	for _, s := range addrs {
		a, err := netip.ParseAddr(s)
		if err != nil || !a.Is6() || a.IsLinkLocalUnicast() {
			continue
		}

		if !best.IsValid() || a.Compare(best) < 0 {
			best = a
		}
	}

	if !best.IsValid() {
		return ""
	}

	return best.String()
}

func parseCounter(s string) *uint64 {
	if s == "" {
		return nil
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil
	}

	return &v
}
