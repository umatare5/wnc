package wnc

import (
	"context"

	"github.com/umatare5/cisco-ios-xe-wireless-go/service/ap"

	"github.com/umatare5/wnc/internal/config"
)

// RadioAdmin is one radio's identity and admin state. BandWire is the number
// set-ap-slot-admin-state takes and follows radio-type, BandLabel is the band the radio is
// serving and follows current-active-band, and a spelling neither table holds leaves its own
// field empty for the caller to refuse.
type RadioAdmin struct {
	Slot      int
	Type      string
	Band      string
	BandLabel string
	BandWire  uint32
	Admin     string
}

// IsRadio reports whether this record is a radio at all. A remote-LAN port is listed in
// the same list and carries neither a band nor an admin state, so there is nothing to set.
func (r RadioAdmin) IsRadio() bool {
	return r.Type != radioTypeRemoteLAN
}

// The band numbers set-ap-slot-admin-state declares. Dual band names a kind of radio rather than a
// frequency, so no table keyed on the band a radio is serving can produce it.
const (
	bandWire24   uint32 = 1
	bandWire5    uint32 = 2
	bandWireDual uint32 = 3
	bandWire6    uint32 = 4
)

// radioTypeBand maps enm-radio-type to the band number the RPC takes, following what occupies the
// slot rather than the band the radio is serving. Four spellings are absent: the RPC's 1-to-4
// domain names no invalid, UWB or remote-LAN radio, and a 2.4/6 GHz XOR radio fits both the
// dual-band 3 and the 6 GHz 4 with no controller asked which. A spelling this table does not hold
// leaves BandWire zero, which is the caller's refusal signal rather than a default.
var radioTypeBand = map[string]uint32{
	"radio-80211bg":          bandWire24,
	"radio-80211a":           bandWire5,
	"radio-80211abgn":        bandWireDual,
	"radio-80211-xor-5-6ghz": bandWireDual,
	"radio-80211-6ghz":       bandWire6,
}

// activeBandLabel maps current-active-band to the band the operator is shown, and stays a table
// of its own because neither leaf determines the other: dot11-6-ghz-band takes band 3 on an XOR
// radio where a dedicated 6 GHz radio takes 4. dot11-invalid-band is absent because it is an
// explicit member of the read domain and names no band a prompt may claim.
var activeBandLabel = map[string]string{
	"dot11-2-dot-4-ghz-band": config.Band24,
	"dot11-5-ghz-band":       config.Band5,
	"dot11-6-ghz-band":       config.Band6,
}

// slotAllowedForBand mirrors the range the RPC's own must statement permits, so a pair the
// controller would answer 400 on is refused here instead.
var slotAllowedForBand = map[uint32][]int{
	bandWire24:   {0},
	bandWire5:    {1, 2},
	bandWireDual: {0, 2},
	bandWire6:    {2, 3},
}

// MaxRadioSlot is the highest slot any band the RPC accepts is permitted on. It is derived
// from the must table rather than declared beside it, so the bound the CLI checks an
// operator's --slot against and the pairs a band is checked against cannot disagree.
func MaxRadioSlot() int {
	highest := 0

	for _, slots := range slotAllowedForBand {
		for _, s := range slots {
			highest = max(highest, s)
		}
	}

	return highest
}

// SlotAllowed reports whether the RPC accepts this radio's slot for the band it reports.
func (r RadioAdmin) SlotAllowed() bool {
	for _, s := range slotAllowedForBand[r.BandWire] {
		if s == r.Slot {
			return true
		}
	}

	return false
}

// AllowedSlots lists the slots the RPC accepts for the band this radio reports, so a
// refusal can name them.
func (r RadioAdmin) AllowedSlots() []int {
	return slotAllowedForBand[r.BandWire]
}

// RadioBySlot reads one radio by the list's own key, wtp-mac and radio-slot-id together, so the
// read returns one record or none. It is not pruned: a fields expression naming a node the release
// does not declare answers 200 with a body that stops mid-object, and every leaf below is
// when-guarded on the radio type.
func (c *Client) RadioBySlot(ctx context.Context, radioMAC string, slot int) (*RadioAdmin, error) {
	resp, err := c.sdk.AP().GetRadioStatusByWTPMACAndSlot(ctx, radioMAC, slot)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.RadioOperData) == 0 {
		return nil, nil
	}

	r := resp.RadioOperData[0]

	// A spelling one table does not hold leaves only its own field empty, which is what
	// radioForAdmin refuses on.
	return &RadioAdmin{
		Slot:      r.RadioSlotID,
		Type:      string(r.RadioType),
		Band:      r.CurrentActiveBand,
		BandLabel: activeBandLabel[r.CurrentActiveBand],
		BandWire:  radioTypeBand[string(r.RadioType)],
		Admin:     r.AdminState,
	}, nil
}

// SetAPAdminState enables or disables one access point. The access point stays registered and does
// not reboot, and both radios go on reporting their own admin-state as enabled while their
// oper-state goes down, so show ap's Admin column is the authority here and show overview's is not.
func (c *Client) SetAPAdminState(ctx context.Context, apName string, on bool) error {
	if on {
		return c.sdk.AP().EnableAPByName(ctx, apName)
	}

	return c.sdk.AP().DisableAPByName(ctx, apName)
}

// SetRadioAdminState enables or disables one radio, named by its slot and by the radio type the
// controller reported for it. The SDK derives the band number from that type and refuses no slot,
// so radioTypeBand reaches no wire and stays only to feed slotAllowedForBand;
// TestRadioTypeBandMatchesTheSDKsOwnDerivation keeps the two derivations from drifting apart.
func (c *Client) SetRadioAdminState(
	ctx context.Context, radioMAC string, slot int, radioType string, on bool,
) error {
	if on {
		return c.sdk.AP().EnableRadioByMAC(ctx, radioMAC, slot, ap.RadioType(radioType))
	}

	return c.sdk.AP().DisableRadioByMAC(ctx, radioMAC, slot, ap.RadioType(radioType))
}
