package wnc

import (
	"context"
	"fmt"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
	"github.com/umatare5/cisco-ios-xe-wireless-go/service/ap"
)

// radioTypeRemoteLAN is the pseudo-radio the controller lists for a remote-LAN port. A "when
// radio-type != 'radio-remote-lan'" guard withholds every leaf but the list key and this one on
// all three releases in scope, so no radio quantity is readable and Radios drops the record.
const radioTypeRemoteLAN = "radio-remote-lan"

// Radio is one access point radio. The identity is the pair (APMAC, Slot), which is the YANG list
// key: APMAC alone is the access point's base radio address and every radio of one access point
// reports the same value.
type Radio struct {
	APName     string
	APMAC      string
	Slot       int
	RadioType  string
	Mode       string
	Band       string
	AdminState string
	OperState  string

	// Channel and Width are nil when the controller withheld the leaf. A Monitor or
	// Sniffer radio is the case that matters: its container arrives with the width and
	// omits the channel, so the two cannot share one guard.
	Channel *int
	Width   *int

	TxPowerDBm *int8
	RFProfile  string

	// Clients is nil when the count could not be established, which is not the same
	// as a radio with no clients.
	Clients *int

	// ChUtil is nil when the radio has no measurement row at all.
	ChUtil *int
}

// OverviewReads reports which secondary read failed.
type OverviewReads struct {
	CAPWAP  error
	Clients error
	RRM     error
	RFTags  error
}

// radioKey identifies one radio, matching the YANG list key.
type radioKey struct {
	mac  string
	slot int
}

// Radios reads the per-radio view. The radio collection drives the rows; the access
// point names, the client tally, the channel utilization and the RF profile names are
// each secondary and cost only their own cells.
func (c *Client) Radios(ctx context.Context) ([]Radio, OverviewReads, error) {
	resp, err := c.sdk.AP().ListRadioData(ctx)
	if err != nil {
		return nil, OverviewReads{}, fmt.Errorf("reading radio-oper-data: %w", err)
	}

	if resp == nil {
		return nil, OverviewReads{}, nil
	}

	names, tags, capwapErr := c.radioAPInfo(ctx)
	clients, clientsErr := c.radioClientCounts(ctx, names)
	util, rrmErr := c.radioUtilization(ctx)
	profiles, rfErr := c.rfTagProfiles(ctx)

	reads := OverviewReads{CAPWAP: capwapErr, Clients: clientsErr, RRM: rrmErr, RFTags: rfErr}
	radios := make([]Radio, 0, len(resp.RadioOperData))

	for _, r := range resp.RadioOperData {
		if r.RadioType == radioTypeRemoteLAN {
			continue
		}

		key := radioKey{mac: r.WtpMAC, slot: r.RadioSlotID}
		radio := Radio{
			APName:     names[r.WtpMAC],
			APMAC:      r.WtpMAC,
			Slot:       r.RadioSlotID,
			RadioType:  string(r.RadioType),
			Mode:       r.RadioMode,
			Band:       r.CurrentActiveBand,
			AdminState: r.AdminState,
			OperState:  r.OperState,
			TxPowerDBm: txPower(r.RadioBandInfo, r.CurrentBandID),
			RFProfile:  profiles[bandProfileKey{tag: tags[r.WtpMAC], band: r.CurrentActiveBand}],
		}

		if r.PhyHtCfg != nil {
			radio.Channel = r.PhyHtCfg.CfgData.CurrFreq
			radio.Width = r.PhyHtCfg.CfgData.ChanWidth
		}

		if clients != nil {
			radio.Clients = ptrTo(clients[key])
		}

		if u, ok := util[key]; ok {
			radio.ChUtil = ptrTo(u)
		}

		radios = append(radios, radio)
	}

	return radios, reads, nil
}

// radioAPInfo maps each access point's address to its name and to the RF tag in force. It takes
// apTagFields, the same three nodes the tag view reads, and names the whole tag-info container for
// the reason stated there.
func (c *Client) radioAPInfo(ctx context.Context) (names, tags map[string]string, err error) {
	resp, err := c.sdk.AP().ListCAPWAPData(ctx, sdk.WithFields(apTagFields))
	if err != nil {
		return nil, nil, fmt.Errorf("reading capwap-data: %w", err)
	}

	if resp == nil {
		return nil, nil, nil
	}

	names = make(map[string]string, len(resp.CAPWAPData))
	tags = make(map[string]string, len(resp.CAPWAPData))

	for _, ap := range resp.CAPWAPData {
		names[ap.WtpMAC] = ap.Name
		tags[ap.WtpMAC] = ap.TagInfo.ResolvedTagInfo.ResolvedRFTag
	}

	return names, tags, nil
}

// radioClientCounts tallies the clients the controller has in the run state, per radio; a nil
// result means the count could not be established at all. The tally comes from the client list
// rather than from the measurement row's stations leaf, because that list is shorter than the radio
// list and a radio without a row would report zero through a non-pointer int; the client records
// carry the access point's name and no address, so without the name map there is no count.
func (c *Client) radioClientCounts(ctx context.Context, names map[string]string) (map[radioKey]int, error) {
	resp, err := c.sdk.Client().ListCommonInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading common-oper-data: %w", err)
	}

	if resp == nil || names == nil {
		return nil, nil
	}

	macOf := make(map[string]string, len(names))
	for mac, name := range names {
		macOf[name] = mac
	}

	counts := make(map[radioKey]int, len(names))

	for _, cl := range resp.CommonOperData {
		if cl.CoState != clientStateRun {
			continue
		}

		if mac, ok := macOf[cl.ApName]; ok {
			counts[radioKey{mac: mac, slot: cl.MsApSlotID}]++
		}
	}

	return counts, nil
}

// clientStateRun is the association state that counts as a forwarding client.
const clientStateRun = "client-status-run"

// radioUtilization indexes the channel utilization by radio. The load container is a pointer and
// its presence is what marks the reading as real, but the percentage inside it is not, so a present
// container with the leaf omitted would read as an idle channel and nothing can tell those apart.
func (c *Client) radioUtilization(ctx context.Context) (map[radioKey]int, error) {
	resp, err := c.sdk.RRM().ListRRMMeasurement(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading rrm-measurement: %w", err)
	}

	if resp == nil {
		return nil, nil
	}

	out := make(map[radioKey]int, len(resp.RRMMeasurement))

	for _, m := range resp.RRMMeasurement {
		if m.Load == nil {
			continue
		}

		out[radioKey{mac: m.WtpMAC, slot: m.RadioSlotID}] = m.Load.CcaUtilPercentage
	}

	return out, nil
}

// bandProfileKey selects one RF profile name: the tag in force and the band the radio
// is on.
type bandProfileKey struct {
	tag  string
	band string
}

const bandsPerTag = 3

// Band spellings of the radio band enum, used to pick the matching per-band leaf.
const (
	bandEnum24 = "dot11-2-dot-4-ghz-band"
	bandEnum5  = "dot11-5-ghz-band"
	bandEnum6  = "dot11-6-ghz-band"
)

// rfTagProfiles reads the RF tags and indexes each per-band profile name. It is the one read in
// this view that takes with-defaults=report-all, needed because the built-in default tag omits all
// three per-band names on a plain read, and the operational reads must never take it because
// absence there is structural. The profile is selected by the radio's band and not its slot,
// because an XOR radio on slot 2 can be operating in 5 or 6 GHz.
func (c *Client) rfTagProfiles(ctx context.Context) (map[bandProfileKey]string, error) {
	resp, err := c.sdk.RF().ListRFTags(ctx, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return nil, fmt.Errorf("reading rf-tags: %w", err)
	}

	if resp == nil || resp.RFTags == nil {
		return nil, nil
	}

	out := make(map[bandProfileKey]string, len(resp.RFTags.RFTags)*bandsPerTag)

	for _, tag := range resp.RFTags.RFTags {
		out[bandProfileKey{tag: tag.TagName, band: bandEnum24}] = deref(tag.Dot11BRfProfileName)
		out[bandProfileKey{tag: tag.TagName, band: bandEnum5}] = deref(tag.Dot11ARfProfileName)
		out[bandProfileKey{tag: tag.TagName, band: bandEnum6}] = deref(tag.Dot116GhzRFProfName)
	}

	return out, nil
}

// txPower picks the transmit power of the band the radio is actually on, matching band-id against
// current-band-id. Those are bare integers on a domain of 0 for 2.4, 1 for 5 and 2 for 6 GHz, a
// different numbering from the band enum, so the two must never be converted into one another;
// neither taking the first record nor indexing by the band id is correct, and a nil
// current-band-id is refused rather than read as band 0.
func txPower(bands []ap.RadioBandInfo, currentBandID *int) *int8 {
	if currentBandID == nil {
		return nil
	}

	for _, b := range bands {
		if int(b.BandID) == *currentBandID {
			return b.PhyTxPwrLvlCfg.CfgData.CurrTxPowerInDbm
		}
	}

	return nil
}

// deref reads through a pointer the SDK declares for a leaf whose absence and whose empty
// value render alike downstream, so the two need not be kept apart here.
func deref[T any](p *T) T {
	if p == nil {
		var zero T

		return zero
	}

	return *p
}

func ptrTo[T any](v T) *T {
	return &v
}
