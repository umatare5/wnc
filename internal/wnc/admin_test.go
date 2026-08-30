package wnc

import (
	"net/http"
	"slices"
	"strconv"
	"testing"
)

// The paths the harness routes on. The keyed read carries the list's whole key in the last element,
// wtp-mac and radio-slot-id joined by a comma, so a read naming another slot routes nowhere. Both
// admin-state paths are spelt out because the CLI declares neither and the SDK's route constants
// are internal.
const (
	apAdminStateRPC  = "Cisco-IOS-XE-wireless-access-point-cfg-rpc:set-ap-admin-state"
	slotAdminRPCPath = "Cisco-IOS-XE-wireless-access-point-cfg-rpc:set-ap-slot-admin-state"
	keyedSlot0AP1    = "radio-oper-data=" + macAP1 + ",0"
	keyedSlot1AP1    = "radio-oper-data=" + macAP1 + ",1"
	keyedSlot2AP1    = "radio-oper-data=" + macAP1 + ",2"
	keyedSlot3AP1    = "radio-oper-data=" + macAP1 + ",3"
	keyedSlot2AP2    = "radio-oper-data=" + macAP2 + ",2"
)

// radioOperData wraps rows in the envelope the keyed read answers with. No row carries a
// slot-id leaf: the record declares one beside radio-slot-id, and radio-slot-id is the list
// key and so the slot a write names.
func radioOperData(rows string) string {
	return `{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[` + rows + `]}`
}

// The wire number comes from the radio type and the label from the band the radio serves, and
// neither leaf determines the other: one type reports two served bands over its life and takes one
// number either way, while dot11-6-ghz-band takes 3 on an XOR radio and 4 on a dedicated one. The
// two rows carrying only one of the pair are the assertion that a missing key leaves its own field
// empty rather than zeroing the other's.
func TestRadioBySlotTakesTheWireFromTheTypeAndTheLabelFromTheBand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		route   string
		mac     string
		slot    int
		row     string
		want    RadioAdmin
		isRadio bool
	}{
		{
			name:  "an XOR radio the controller reports on 6 GHz",
			route: keyedSlot2AP1, mac: macAP1, slot: 2,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":2,"radio-type":"radio-80211-xor-5-6ghz",` +
				`"current-active-band":"dot11-6-ghz-band","current-band-id":2,` +
				`"admin-state":"enabled","oper-state":"radio-up"}`,
			want: RadioAdmin{
				Slot: 2, Type: "radio-80211-xor-5-6ghz", Band: "dot11-6-ghz-band",
				BandLabel: "6", BandWire: 3, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// The label moves with the served band and the wire number does not: the radio is
			// the same piece of hardware either way.
			name:  "the same XOR radio moved to 5 GHz",
			route: keyedSlot2AP1, mac: macAP1, slot: 2,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":2,"radio-type":"radio-80211-xor-5-6ghz",` +
				`"current-active-band":"dot11-5-ghz-band","current-band-id":1,` +
				`"admin-state":"enabled","oper-state":"radio-up"}`,
			want: RadioAdmin{
				Slot: 2, Type: "radio-80211-xor-5-6ghz", Band: "dot11-5-ghz-band",
				BandLabel: "5", BandWire: 3, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// A dual-band radio serving 2.4 GHz takes 3 and not 1 — measured on 17.12, where
			// band 1 against this type answered 400 and band 3 answered 204.
			name:  "a dual-band radio a write has already disabled",
			route: keyedSlot0AP1, mac: macAP1, slot: 0,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"radio-type":"radio-80211abgn",` +
				`"current-active-band":"dot11-2-dot-4-ghz-band","current-band-id":0,` +
				`"admin-state":"disabled","oper-state":"radio-down"}`,
			want: RadioAdmin{
				Slot: 0, Type: "radio-80211abgn", Band: "dot11-2-dot-4-ghz-band",
				BandLabel: "2.4", BandWire: 3, Admin: "disabled",
			},
			isRadio: true,
		},
		{
			name:  "a dedicated 2.4 GHz radio",
			route: keyedSlot0AP1, mac: macAP1, slot: 0,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"radio-type":"radio-80211bg",` +
				`"current-active-band":"dot11-2-dot-4-ghz-band","admin-state":"enabled"}`,
			want: RadioAdmin{
				Slot: 0, Type: "radio-80211bg", Band: "dot11-2-dot-4-ghz-band",
				BandLabel: "2.4", BandWire: 1, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// Composed rather than observed: the lab holds no dedicated 5 GHz radio in slot 2
			// and no dedicated 6 GHz radio at all. Each pins what the CLI would send, and the
			// must clause is the only thing behind the pair.
			name:  "a dedicated 5 GHz radio in slot 2",
			route: keyedSlot2AP1, mac: macAP1, slot: 2,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":2,"radio-type":"radio-80211a",` +
				`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`,
			want: RadioAdmin{
				Slot: 2, Type: "radio-80211a", Band: "dot11-5-ghz-band",
				BandLabel: "5", BandWire: 2, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			name:  "a dedicated 6 GHz radio in slot 3",
			route: keyedSlot3AP1, mac: macAP1, slot: 3,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":3,"radio-type":"radio-80211-6ghz",` +
				`"current-active-band":"dot11-6-ghz-band","admin-state":"enabled"}`,
			want: RadioAdmin{
				Slot: 3, Type: "radio-80211-6ghz", Band: "dot11-6-ghz-band",
				BandLabel: "6", BandWire: 4, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// A label with no wire: the RPC's band domain holds no number for this type.
			name:  "a radio type the RPC has no band number for",
			route: keyedSlot1AP1, mac: macAP1, slot: 1,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":1,"radio-type":"radio-uwb",` +
				`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`,
			want: RadioAdmin{
				Slot: 1, Type: "radio-uwb", Band: "dot11-5-ghz-band",
				BandLabel: "5", Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// A wire with no label: the served band names nothing a prompt may claim, while the
			// type still names a number. The caller refuses on the label.
			name:  "a served band the prompt has no label for",
			route: keyedSlot2AP1, mac: macAP1, slot: 2,
			row: `{"wtp-mac":"` + macAP1 + `","radio-slot-id":2,"radio-type":"radio-80211-xor-5-6ghz",` +
				`"current-active-band":"dot11-invalid-band","admin-state":"enabled"}`,
			want: RadioAdmin{
				Slot: 2, Type: "radio-80211-xor-5-6ghz", Band: "dot11-invalid-band",
				BandWire: 3, Admin: "enabled",
			},
			isRadio: true,
		},
		{
			// The remote-LAN port is listed in the same list and carries neither a band
			// nor an admin state, so there is nothing for a write to set.
			name:  "a remote-LAN pseudo-radio",
			route: keyedSlot2AP2, mac: macAP2, slot: 2,
			row:  `{"wtp-mac":"` + macAP2 + `","radio-slot-id":2,"radio-type":"radio-remote-lan"}`,
			want: RadioAdmin{Slot: 2, Type: "radio-remote-lan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{tt.route: {body: radioOperData(tt.row)}})

			got, err := c.RadioBySlot(t.Context(), tt.mac, tt.slot)
			if err != nil {
				t.Fatalf("RadioBySlot: %v", err)
			}

			if got == nil {
				t.Fatal("a record decoded to no radio")
			}

			if *got != tt.want {
				t.Errorf("radio = %+v, want %+v", *got, tt.want)
			}

			if got.IsRadio() != tt.isRadio {
				t.Errorf("IsRadio = %v, want %v", got.IsRadio(), tt.isRadio)
			}
		})
	}
}

// A slot the access point does not hold and a controller that failed must not arrive as the
// same thing: the CLI reports the first as holding no radio in the slot and the second as a
// read failure naming the controller.
func TestRadioBySlotSeparatesAnAbsentRadioFromAFailedRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		served     answer
		wantCause  Cause
		wantStatus int
	}{
		{name: "an empty list", served: answer{body: radioOperData("")}},
		{name: "a 204 with no body", served: answer{status: http.StatusNoContent}},
		{
			name: "a 404", served: answer{status: http.StatusNotFound},
			wantCause: CauseNotFound, wantStatus: http.StatusNotFound,
		},
		{
			name: "a 500", served: answer{status: http.StatusInternalServerError},
			wantCause: CauseHTTP, wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{keyedSlot2AP1: tt.served})

			got, err := c.RadioBySlot(t.Context(), macAP1, 2)

			if tt.wantCause == "" {
				if err != nil {
					t.Fatalf("an absence produced an error: %v", err)
				}

				if got != nil {
					t.Errorf("radio = %+v, want none", *got)
				}

				return
			}

			if err == nil {
				t.Fatalf("a %d was reported as an absence", tt.wantStatus)
			}

			if cause, status := Classify(err); cause != tt.wantCause || status != tt.wantStatus {
				t.Errorf("Classify = %q/%d, want %q/%d", cause, status, tt.wantCause, tt.wantStatus)
			}
		})
	}
}

// This read is deliberately not pruned, and nothing in the answer shows that: a fields
// expression naming a node the release does not declare answers 200 with a body that stops
// mid-object, and every leaf this decodes is when-guarded on the radio type. The query is
// asserted so an optimisation cannot add one silently.
func TestRadioBySlotReadsTheWholeRecordRatherThanPruningIt(t *testing.T) {
	t.Parallel()

	var got string

	c := newClientWithQuery(t, &got, radioOperData(
		`{"wtp-mac":"`+macAP1+`","radio-slot-id":2,"current-active-band":"dot11-6-ghz-band"}`))

	if _, err := c.RadioBySlot(t.Context(), macAP1, 2); err != nil {
		t.Fatalf("RadioBySlot: %v", err)
	}

	if got != "" {
		t.Errorf("query = %q, want none", got)
	}
}

// The band on the wire is the one the radio type takes, which for an XOR radio serving 6 GHz is 3
// and not the 4 its served band would suggest. The band asserted here is the one the keyed read
// produced, so a table keyed on the served band fails this test rather than passing through it.
func TestSetRadioAdminStateSendsTheBandTheRadioTypeTakes(t *testing.T) {
	t.Parallel()

	const sixGHzSlot2 = `{"wtp-mac":"` + macAP1 + `","radio-slot-id":2,` +
		`"radio-type":"radio-80211-xor-5-6ghz","current-active-band":"dot11-6-ghz-band",` +
		`"admin-state":"enabled","oper-state":"radio-up"}`

	const fiveGHzSlot1 = `{"wtp-mac":"` + macAP1 + `","radio-slot-id":1,` +
		`"radio-type":"radio-80211a","current-active-band":"dot11-5-ghz-band",` +
		`"admin-state":"enabled","oper-state":"radio-up"}`

	tests := []struct {
		name  string
		route string
		row   string
		slot  int
		on    bool
		want  string
	}{
		{
			name: "disable an XOR radio", route: keyedSlot2AP1, row: sixGHzSlot2, slot: 2, on: false,
			want: `{"Cisco-IOS-XE-wireless-access-point-cfg-rpc:input":{"mode":"admin-state-disabled",` +
				`"slot-id":2,"band":"3","mac-addr":"` + macAP1 + `"}}`,
		},
		{
			name: "enable an XOR radio", route: keyedSlot2AP1, row: sixGHzSlot2, slot: 2, on: true,
			want: `{"Cisco-IOS-XE-wireless-access-point-cfg-rpc:input":{"mode":"admin-state-enabled",` +
				`"slot-id":2,"band":"3","mac-addr":"` + macAP1 + `"}}`,
		},
		{
			// The second row is what stops the table answering dual band unconditionally.
			name: "disable a dedicated radio", route: keyedSlot1AP1, row: fiveGHzSlot1, slot: 1, on: false,
			want: `{"Cisco-IOS-XE-wireless-access-point-cfg-rpc:input":{"mode":"admin-state-disabled",` +
				`"slot-id":1,"band":"2","mac-addr":"` + macAP1 + `"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			read := newClient(t, routes{tt.route: {body: radioOperData(tt.row)}})

			radio, err := read.RadioBySlot(t.Context(), macAP1, tt.slot)
			if err != nil {
				t.Fatalf("RadioBySlot: %v", err)
			}

			r := newRecorder(t, http.StatusNoContent, "")

			if err := r.client.SetRadioAdminState(t.Context(), macAP1, radio.Slot, radio.Type, tt.on); err != nil {
				t.Fatalf("SetRadioAdminState: %v", err)
			}

			got := r.last(t)
			if !contains(got.path, slotAdminRPCPath) {
				t.Errorf("path = %s, want the slot operation", got.path)
			}

			if got.body != tt.want {
				t.Errorf("payload =\n  %s\nwant\n  %s", got.body, tt.want)
			}
		})
	}
}

// The mode spellings are the RPC's write family, admin-state-*, where the read side reports
// adminstate-* for the same two states, so the two domains must never converge. The arm is
// asserted beside the mode: the input is a mandatory choice and an arm that changed would still be
// accepted by the controller while needing an address this leaf does not carry.
func TestSetAPAdminStateSendsTheWriteFamilyMode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		on   bool
		want string
	}{
		{
			name: "enable", on: true,
			want: `{"Cisco-IOS-XE-wireless-access-point-cfg-rpc:input":` +
				`{"mode":"admin-state-enabled","ap-name":"` + nameAP1 + `"}}`,
		},
		{
			name: "disable", on: false,
			want: `{"Cisco-IOS-XE-wireless-access-point-cfg-rpc:input":` +
				`{"mode":"admin-state-disabled","ap-name":"` + nameAP1 + `"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusNoContent, "")

			if err := r.client.SetAPAdminState(t.Context(), nameAP1, tt.on); err != nil {
				t.Fatalf("SetAPAdminState: %v", err)
			}

			// One exchange: the name arm goes straight to the RPC, so the CLI's own keyed
			// resolve is the only read in front of it.
			if n := len(r.all()); n != 1 {
				t.Errorf("exchanges = %d, want the write alone", n)
			}

			got := r.last(t)
			if !contains(got.path, apAdminStateRPC) {
				t.Errorf("path = %s, want the access-point operation", got.path)
			}

			if got.body != tt.want {
				t.Errorf("payload =\n  %s\nwant\n  %s", got.body, tt.want)
			}
		})
	}
}

// A pair the must clause refuses is answered with a bare 400, and both writes have to carry
// that to the caller as a failure rather than as a sent instruction. The harness routes on
// the operation, so a write posted elsewhere fails here too.
func TestAdminWritesReportTheControllerRefusal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rpc  string
		call func(*Client) error
	}{
		{
			name: "the access point", rpc: apAdminStateRPC,
			call: func(c *Client) error { return c.SetAPAdminState(t.Context(), nameAP1, false) },
		},
		{
			name: "one radio", rpc: slotAdminRPCPath,
			call: func(c *Client) error {
				return c.SetRadioAdminState(t.Context(), macAP1, 1, "radio-80211a", false)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{tt.rpc: {status: http.StatusBadRequest}})

			err := tt.call(c)
			if err == nil {
				t.Fatal("a 400 did not produce an error")
			}

			if cause, status := Classify(err); cause != CauseHTTP || status != http.StatusBadRequest {
				t.Errorf("Classify = %q/%d, want http/400", cause, status)
			}
		})
	}
}

// slotAllowedForBand mirrors the RPC's own must clause, so a pair the controller would refuse
// is refused before a socket opens. Measured on 17.15.6: band 1 with slot 1 answers 400
// invalid-value, naming the must. Band 2 on slot 2 is the pair the SDK cannot express and
// must stay accepted here.
func TestSlotAllowedMirrorsTheRPCMustClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		radio     RadioAdmin
		want      bool
		wantSlots []int
	}{
		{
			name:  "2.4 GHz on slot 0",
			radio: RadioAdmin{Slot: 0, BandWire: 1}, want: true, wantSlots: []int{0},
		},
		{
			name:  "2.4 GHz on slot 1",
			radio: RadioAdmin{Slot: 1, BandWire: 1}, wantSlots: []int{0},
		},
		{
			name:  "5 GHz on slot 2",
			radio: RadioAdmin{Slot: 2, BandWire: 2}, want: true, wantSlots: []int{1, 2},
		},
		{
			name:  "5 GHz on slot 0",
			radio: RadioAdmin{Slot: 0, BandWire: 2}, wantSlots: []int{1, 2},
		},
		{
			name:  "6 GHz on slot 3",
			radio: RadioAdmin{Slot: 3, BandWire: 4}, want: true, wantSlots: []int{2, 3},
		},
		{
			name:  "6 GHz on slot 1",
			radio: RadioAdmin{Slot: 1, BandWire: 4}, wantSlots: []int{2, 3},
		},
		{
			name:  "dual band on slot 0",
			radio: RadioAdmin{Slot: 0, BandWire: 3}, want: true, wantSlots: []int{0, 2},
		},
		{
			name:  "dual band on slot 2",
			radio: RadioAdmin{Slot: 2, BandWire: 3}, want: true, wantSlots: []int{0, 2},
		},
		{
			name:  "dual band on slot 1",
			radio: RadioAdmin{Slot: 1, BandWire: 3}, wantSlots: []int{0, 2},
		},
		{
			// A radio type with no number names no slot either, or a refusal would list one.
			name:  "a radio type the table has no number for",
			radio: RadioAdmin{Slot: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.radio.SlotAllowed(); got != tt.want {
				t.Errorf("SlotAllowed = %v, want %v", got, tt.want)
			}

			if got := tt.radio.AllowedSlots(); !slices.Equal(got, tt.wantSlots) {
				t.Errorf("AllowedSlots = %v, want %v", got, tt.wantSlots)
			}
		})
	}
}

// MaxRadioSlot is derived from the must table rather than declared beside it, so a slot added
// to a band there moves the bound the CLI reports without a second edit. The number is the
// table's own maximum and 3 is what the CLI's two --slot messages read.
func TestMaxRadioSlotIsTheMustTablesOwnMaximum(t *testing.T) {
	t.Parallel()

	want := 0

	for _, slots := range slotAllowedForBand {
		for _, s := range slots {
			want = max(want, s)
		}
	}

	if got := MaxRadioSlot(); got != want {
		t.Errorf("MaxRadioSlot = %d, want the table maximum %d", got, want)
	}

	if got := MaxRadioSlot(); got != 3 {
		t.Errorf("MaxRadioSlot = %d, want 3", got)
	}
}

// The eight members of enm-radio-type that 17.12 and 17.15 declare, against the band number the
// RPC takes for each; the ninth 17.18 adds is covered below. Three carry no number: the RPC's
// domain is 1 to 4 and nothing there names an invalid radio, a UWB radio or a remote-LAN port.
func TestRadioTypeBandHoldsEveryMemberOfTheRadioTypeEnum(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		spelling string
		want     uint32
	}{
		{spelling: "radio-invalid", want: 0},
		{spelling: "radio-80211bg", want: 1},
		{spelling: "radio-80211a", want: 2},
		{spelling: "radio-80211abgn", want: 3},
		{spelling: "radio-uwb", want: 0},
		{spelling: "radio-remote-lan", want: 0},
		{spelling: "radio-80211-6ghz", want: 4},
		{spelling: "radio-80211-xor-5-6ghz", want: 3},
	} {
		t.Run(tt.spelling, func(t *testing.T) {
			t.Parallel()

			if got := radioTypeBand[tt.spelling]; got != tt.want {
				t.Errorf("radioTypeBand[%q] = %d, want %d", tt.spelling, got, tt.want)
			}
		})
	}

	if len(radioTypeBand) != 5 {
		t.Errorf("radioTypeBand holds %d types, want the 5 the RPC has a number for", len(radioTypeBand))
	}
}

// The must clause the controller serves declares exactly four bands, and the shipped table omitted
// one of them while nothing derived it. This is the direct regression for that omission and for a
// fifth row arriving without the schema saying so.
func TestSlotAllowedForBandCoversEveryBandTheRPCDeclares(t *testing.T) {
	t.Parallel()

	want := []uint32{bandWire24, bandWire5, bandWireDual, bandWire6}

	got := make([]uint32, 0, len(slotAllowedForBand))
	for band := range slotAllowedForBand {
		got = append(got, band)
	}

	slices.Sort(got)

	if !slices.Equal(got, want) {
		t.Errorf("slotAllowedForBand keys = %v, want %v", got, want)
	}
}

// radioTypeBand reaches no wire: the SDK derives the band from the radio type and this table only
// feeds slotAllowedForBand, which the SDK has no equivalent of. The same derivation therefore
// lives in two places, so this reads the band off the request rather than trusting either table.
func TestRadioTypeBandMatchesTheSDKsOwnDerivation(t *testing.T) {
	t.Parallel()

	if len(radioTypeBand) == 0 {
		t.Fatal("the table is empty, so this asserts nothing")
	}

	for spelling, band := range radioTypeBand {
		t.Run(spelling, func(t *testing.T) {
			t.Parallel()

			r := newRecorder(t, http.StatusNoContent, "")

			if err := r.client.SetRadioAdminState(t.Context(), macAP1, 2, spelling, true); err != nil {
				t.Fatalf("SetRadioAdminState: %v", err)
			}

			if want := `"band":"` + strconv.FormatUint(uint64(band), 10) + `"`; !contains(r.last(t).body, want) {
				t.Errorf("the SDK sent %s, but this tree's table says %s", r.last(t).body, want)
			}
		})
	}
}

// The other half of the same contract: a spelling this tree has no number for must be one the SDK
// refuses too. If the SDK started numbering one of them, the CLI would refuse a write the SDK
// would have completed, and the refusal would name a radio the controller does support.
func TestTheSDKRefusesEveryRadioTypeThisTreeHasNoBandFor(t *testing.T) {
	t.Parallel()

	// Read from enm-radio-type as the controllers serve it: 17.12.8 and 17.15.6 declare eight
	// members and 17.18.4a nine, so radio-80211-xor-24-6ghz is the one below the older releases do
	// not declare. not-a-radio-type is declared by none.
	for _, spelling := range []string{
		"radio-80211-xor-24-6ghz", "radio-remote-lan", "radio-uwb", "radio-invalid", "not-a-radio-type",
	} {
		t.Run(spelling, func(t *testing.T) {
			t.Parallel()

			if _, ok := radioTypeBand[spelling]; ok {
				t.Fatalf("%s is in the table, so this row is stale", spelling)
			}

			r := newRecorder(t, http.StatusNoContent, "")

			if err := r.client.SetRadioAdminState(t.Context(), macAP1, 2, spelling, true); err == nil {
				t.Errorf("the SDK accepted %s, which this tree refuses before the write", spelling)
			}
		})
	}
}
