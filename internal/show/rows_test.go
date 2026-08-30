package show

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter/pkg/twwidth"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

var target = config.Target{Name: "lab-wlc"}

// cellsOf renders one row through its own columns, so an assertion sees exactly what
// an operator would.
func cellsOf[R any](cols []render.Column[R], row R) map[string]string {
	out := make(map[string]string, len(cols))
	for _, c := range cols {
		out[c.Key] = c.Cell(row)
	}

	return out
}

func TestClientRowsAbsenceRules(t *testing.T) {
	t.Parallel()

	clients := []wnc.WirelessClient{
		{
			MAC: "00:00:5e:00:53:11", APName: "lab-ap-1", Slot: 0, State: "client-status-run",
			// The controller sends this key with an empty value on most rows, so a
			// presence check would not catch it.
			Username: "",
			IPv4:     "0.0.0.0",
			Band:     "dot11-radio-type-bg", PHY: "client-dot11ax-24ghz-prot",
			// Zero is impossible for a channel, a rate and a stream count on an
			// associated client, and impossible for an RSSI. Zero SNR is a real reading.
			Channel: 0, Speed: 0, Streams: 0, RSSI: 0, SNR: ptr(0),
			HasDot11: true,
		},
	}

	rows := clientRows(clients, ClientFilter{}, target, &Reporter{})
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}

	cells := cellsOf(ClientColumns(), rows[0])

	for _, key := range []string{"username", "ipv4", "channel", "speed_mbps", "spatial_streams", "rssi_dbm"} {
		if cells[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, cells[key], render.Absent)
		}
	}

	// Slot 0 is a real radio and 0 dB of SNR is a real margin, so neither is suppressed.
	if cells["slot"] != "0" || cells["snr_db"] != "0dB" {
		t.Errorf("slot = %q, snr = %q; both should read as reported zeros", cells["slot"], cells["snr_db"])
	}

	// The band and the protocol come from different leaves and must not be the same
	// value restated.
	if cells["band"] != "2.4" || cells["protocol"] != "11ax" {
		t.Errorf("band = %q, protocol = %q", cells["band"], cells["protocol"])
	}
}

// The measured columns carry their unit in the cell only. The JSON has to keep the bare
// number a script parses, and the ordering has to stay numeric: "11ch" sorts before
// "6ch" as text, which would put channel 11 above channel 6.
func TestClientCellsCarryTheirUnitsAndNothingElseDoes(t *testing.T) {
	t.Parallel()

	cols := ClientColumns()

	reported := clientRows([]wnc.WirelessClient{{
		MAC: "00:00:5e:00:53:11", Channel: 6, RSSI: -21, SNR: ptr(78), Speed: 143, Streams: 1,
	}}, ClientFilter{}, target, &Reporter{})

	cells := cellsOf(cols, reported[0])
	for key, want := range map[string]string{
		"channel": "6ch", "rssi_dbm": "-21dBm", "snr_db": "78dB",
		"speed_mbps": "143Mbps", "spatial_streams": "1ss",
	} {
		if cells[key] != want {
			t.Errorf("%s = %q, want %q", key, cells[key], want)
		}
	}

	// A unit on an unreported value would read as a measurement rather than the lack of one.
	// snr_db is in this list because 0 dB is a real margin: it is the one cell here whose
	// value cannot be told from its absence without the pointer the fetch layer sets.
	absent := clientRows([]wnc.WirelessClient{{MAC: "00:00:5e:00:53:12"}}, ClientFilter{}, target, &Reporter{})
	for _, key := range []string{"channel", "rssi_dbm", "snr_db", "speed_mbps", "spatial_streams"} {
		if got := cellsOf(cols, absent[0])[key]; got != render.Absent {
			t.Errorf("%s = %q, want %q", key, got, render.Absent)
		}
	}

	// The JSON is built from the row fields, so no suffix can reach a consumer.
	var buf bytes.Buffer
	if err := render.JSON(&buf, reported); err != nil {
		t.Fatalf("JSON: %v", err)
	}

	for _, suffixed := range []string{`"6ch"`, `"-21dBm"`, `"78dB"`, `"143Mbps"`, `"1ss"`} {
		if strings.Contains(buf.String(), suffixed) {
			t.Errorf("the JSON carries %s:\n%s", suffixed, buf.String())
		}
	}

	for _, bare := range []string{`"channel":6`, `"rssi_dbm":-21`, `"snr_db":78`, `"speed_mbps":143`} {
		if !strings.Contains(buf.String(), bare) {
			t.Errorf("the JSON lost %s:\n%s", bare, buf.String())
		}
	}

	// Text order would be 11, 1, 6; numeric order is 1, 6, 11.
	rows := clientRows([]wnc.WirelessClient{
		{MAC: "b", Channel: 6}, {MAC: "c", Channel: 11}, {MAC: "a", Channel: 1},
	}, ClientFilter{}, target, &Reporter{})

	if err := render.Sort(rows, cols, keyChannel, false); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	got := make([]string, 0, len(rows))
	for _, r := range rows {
		got = append(got, cellsOf(cols, r)[keyChannel])
	}

	if strings.Join(got, ",") != "1ch,6ch,11ch" {
		t.Errorf("channels sorted %v, want the numeric order", got)
	}
}

func TestClientRowsFilters(t *testing.T) {
	t.Parallel()

	clients := []wnc.WirelessClient{
		{MAC: "a", SSID: "labo-a", APName: "ap-1", Band: "dot11-radio-type-bg", HasDot11: true},
		{MAC: "b", SSID: "labo-b", APName: "ap-1", Band: "dot11-radio-type-a", HasDot11: true},
		{MAC: "c", SSID: "labo-a", APName: "ap-2", Band: "dot11-radio-type-6ghz", HasDot11: true},
		// No 802.11 facts at all, so the band cannot be honestly compared.
		{MAC: "d", APName: "ap-1"},
	}

	tests := []struct {
		name   string
		filter ClientFilter
		want   string
	}{
		{name: "no filter", filter: ClientFilter{}, want: "a,b,c,d"},
		{name: "band", filter: ClientFilter{Band: "5"}, want: "b"},
		{name: "ssid", filter: ClientFilter{SSID: "labo-a"}, want: "a,c"},
		{name: "ap name", filter: ClientFilter{APName: "ap-2"}, want: "c"},
		{name: "band and ssid", filter: ClientFilter{Band: "2.4", SSID: "labo-a"}, want: "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rows := clientRows(clients, tt.filter, target, &Reporter{})

			got := make([]string, 0, len(rows))
			for _, r := range rows {
				got = append(got, deref(r.MAC))
			}

			if strings.Join(got, ",") != tt.want {
				t.Errorf("rows = %v, want %s", got, tt.want)
			}
		})
	}
}

// A row dropped because the band was never reported is not the same as a row that did
// not match, so the count is reported rather than silently absorbed.
func TestClientRowsReportsRowsExcludedForAnUnreportedBand(t *testing.T) {
	t.Parallel()

	clients := []wnc.WirelessClient{
		{MAC: "a", Band: "dot11-radio-type-a", HasDot11: true},
		{MAC: "b"},
		{MAC: "c"},
	}

	rep := &Reporter{}
	rows := clientRows(clients, ClientFilter{Band: "5"}, target, rep)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}

	if len(rep.notes) != 1 || !strings.Contains(rep.notes[0], "2 rows excluded") {
		t.Errorf("notes = %#v", rep.notes)
	}
}

func TestClientRowsRejectsTheEpochAssociationTime(t *testing.T) {
	t.Parallel()

	clients := []wnc.WirelessClient{
		{MAC: "a", AssocTime: time.Unix(0, 0)},
		{MAC: "b", AssocTime: time.Now().Add(-time.Hour)},
	}

	rows := clientRows(clients, ClientFilter{}, target, &Reporter{})

	if rows[0].Assoc != nil {
		t.Errorf("the epoch produced an age of %v", *rows[0].Assoc)
	}

	if rows[1].Assoc == nil || *rows[1].Assoc < 3500 {
		t.Errorf("a real instant produced %v", rows[1].Assoc)
	}
}

func TestOverviewRows(t *testing.T) {
	t.Parallel()

	power := int8(19)
	util := 28
	clients := 0

	radios := []wnc.Radio{
		{
			APName: "lab-ap-1", APMAC: "00:00:5e:00:53:01", Slot: 0,
			Mode: "radio-mode-flex-connect", Band: "dot11-2-dot-4-ghz-band",
			AdminState: "enabled", OperState: "radio-up",
			Channel: ptr(11), Width: ptr(20), TxPowerDBm: &power, Clients: &clients, ChUtil: &util,
			RFProfile: "labo-rf-24gh",
		},
		{
			// Everything the controller guards behind a "when" is absent here.
			APName: "lab-ap-2", APMAC: "00:00:5e:00:53:02", Slot: 1,
		},
	}

	rows := overviewRows(radios, "", target, &Reporter{})
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	first := cellsOf(OverviewColumns(), rows[0])
	if first["mode"] != "FlexConnect" || first["band"] != "2.4" || first["oper"] != "Up" {
		t.Errorf("first row = %#v", first)
	}

	// Each measured cell carries its unit; the JSON keeps the bare number.
	for key, want := range map[string]string{
		"channel": "11ch", "channel_width_mhz": "20MHz", "tx_power_dbm": "19dBm",
		"ch_util_percent": "28%",
	} {
		if first[key] != want {
			t.Errorf("%s = %q, want %q", key, first[key], want)
		}
	}

	// A radio the client list could reach but with no client is a reported zero.
	if first["clients"] != "0clients" {
		t.Errorf("clients = %q, want a reported zero", first["clients"])
	}

	second := cellsOf(OverviewColumns(), rows[1])
	// An absent oper state must not be folded into Down: that would report an outage
	// the controller never described.
	for _, key := range []string{"mode", "band", "admin", "oper", "channel", "channel_width_mhz", "tx_power_dbm", "clients", "ch_util_percent", "rf_profile"} {
		if second[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, second[key], render.Absent)
		}
	}

	// The slot is still a reading even when everything else is absent.
	if second["slot"] != "1" {
		t.Errorf("slot = %q", second["slot"])
	}
}

// The overview's Admin column takes the same two glyphs as show ap's, so a reader
// moving between the two views is not asked to learn a second vocabulary. The radio
// domain has exactly two members, and an absent or unknown one must reach neither glyph.
func TestOverviewAdminGlyphs(t *testing.T) {
	t.Parallel()

	cols := OverviewColumns()

	tests := []struct{ admin, want string }{
		{"enabled", glyphOK},
		{"disabled", glyphNo},
		{"quiescing", "quiescing"},
		{"", render.Absent},
	}

	for _, tc := range tests {
		rows := overviewRows([]wnc.Radio{{APMAC: "a", AdminState: tc.admin}}, "", target, &Reporter{})
		if got := prettyOf(cols, rows[0])[keyAdmin]; got != tc.want {
			t.Errorf("admin %q rendered %q, want %q", tc.admin, got, tc.want)
		}
	}
}

func TestOverviewRowsBandFilter(t *testing.T) {
	t.Parallel()

	radios := []wnc.Radio{
		{APMAC: "a", Band: "dot11-2-dot-4-ghz-band"},
		{APMAC: "b", Band: "dot11-5-ghz-band"},
		{APMAC: "c", Band: "dot11-6-ghz-band"},
		{APMAC: "d"},
	}

	rep := &Reporter{}
	rows := overviewRows(radios, "6", target, rep)

	if len(rows) != 1 || deref(rows[0].APMAC) != "c" {
		t.Fatalf("rows = %+v", rows)
	}

	if len(rep.notes) != 1 || !strings.Contains(rep.notes[0], "1 row excluded") {
		t.Errorf("notes = %#v", rep.notes)
	}
}

func TestAPRowsAbsenceRules(t *testing.T) {
	t.Parallel()

	boot := time.Now().Add(-48 * time.Hour)
	join := time.Now().Add(-time.Hour)

	aps := []wnc.AP{
		{
			Name: "lab-ap-1", WtpMAC: "00:00:5e:00:53:01",
			// The controller pads the regulatory code to a fixed width.
			Country: "J4 ", Slots: 2,
			Mode: "mode-monitor", SubMode: "wips-mode",
			AdminState: "adminstate-enabled", OperState: "registered",
			PowerType: "pwr-src-poe-lgcy", PowerMode: "dot11-default-high-pwr",
			BootTime: boot, JoinTime: join,
			Neighbors: []string{"lab-sw-1:Gi0/2", "lab-sw-2:Gi0/3"},
		},
		{Name: "lab-ap-2"},
	}

	rows := apRows(aps, target)

	first := cellsOf(APColumns(), rows[0])

	if first["country"] != "J4" {
		t.Errorf("country = %q, want the padding trimmed", first["country"])
	}

	// The sub-mode is what says a Monitor access point with no clients is healthy.
	if first["mode"] != "Monitor (WIPS)" {
		t.Errorf("mode = %q", first["mode"])
	}

	if first["power_type"] != "PoE (legacy)" || first["power_mode"] != "Full Power (default)" {
		t.Errorf("power = %q / %q", first["power_type"], first["power_mode"])
	}

	// Two different quantities: a controller switchover renews only the second.
	if first["uptime_seconds"] == first["assoc_uptime_seconds"] {
		t.Errorf("uptime and assoc rendered the same value %q", first["uptime_seconds"])
	}

	if first["lldp_neighbor"] != "lab-sw-1:Gi0/2, lab-sw-2:Gi0/3" {
		t.Errorf("neighbors = %q", first["lldp_neighbor"])
	}

	second := cellsOf(APColumns(), rows[1])
	for _, key := range []string{"slots", "country", "mode", "admin", "state", "lldp_neighbor", "power_type", "uptime_seconds"} {
		if second[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, second[key], render.Absent)
		}
	}
}

func TestAPTagRowsAbsenceRules(t *testing.T) {
	t.Parallel()

	no := false
	tags := []wnc.APTag{
		{Name: "lab-ap-1", WtpMAC: "00:00:5e:00:53:01", Misconfigured: &no, TagSource: "tag-source-ap-pnp"},
		{Name: "lab-ap-2", WtpMAC: "00:00:5e:00:53:02"},
	}

	rows := apTagRows(tags, target)

	first := cellsOf(APTagColumns(), rows[0])
	// The controller sends an explicit false on a healthy access point on 17.12, so
	// "No" is a reading and not a substitute for silence.
	if first["misconfigured"] != "No" {
		t.Errorf("misconfigured = %q", first["misconfigured"])
	}

	if first["tag_source"] != "AP-PnP" {
		t.Errorf("tag_source = %q", first["tag_source"])
	}

	second := cellsOf(APTagColumns(), rows[1])
	if second["misconfigured"] != render.Absent {
		t.Errorf("an omitted flag rendered %q", second["misconfigured"])
	}

	for _, key := range []string{
		"tag_source", "misconfig_reason", "filter_name",
		"policy_tag", "site_tag", "rf_tag", "ap_profile", "flex_profile",
	} {
		if second[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, second[key], render.Absent)
		}
	}
}

// The reason column has two ways of saying nothing and they are not the same. The enum
// carries its own "no misconfiguration" member, which is a reading and must render; a
// release that does not declare the leaf sends nothing, which must render as a dash.
func TestAPTagMisconfigReasonKeepsItsNoneApartFromAbsence(t *testing.T) {
	t.Parallel()

	none := "apmgr-no-misconfig"
	country := "country-misconfig"
	future := "some-member-a-later-release-adds"

	tags := []wnc.APTag{
		{Name: "lab-ap-1", MisconfigReason: &none, FilterName: ""},
		{Name: "lab-ap-2", MisconfigReason: &country, FilterName: "labo-filter"},
		{Name: "lab-ap-3", MisconfigReason: &future},
		{Name: "lab-ap-4"},
	}

	rows := apTagRows(tags, target)

	for i, want := range []string{"None", "Country", future, render.Absent} {
		if got := cellsOf(APTagColumns(), rows[i])["misconfig_reason"]; got != want {
			t.Errorf("row %d misconfig_reason = %q, want %q", i, got, want)
		}
	}

	// The controller sends an empty string where no filter exists, so the cell is a dash
	// for the same reason an absent leaf is: there is no filter name to report.
	if got := cellsOf(APTagColumns(), rows[0])["filter_name"]; got != render.Absent {
		t.Errorf("an empty filter name rendered %q", got)
	}

	if got := cellsOf(APTagColumns(), rows[1])["filter_name"]; got != "labo-filter" {
		t.Errorf("filter_name = %q, want labo-filter", got)
	}
}

// An unmapped spelling passes through so a release that adds an enum member shows the
// raw value rather than an empty cell.
func TestEnumPassesUnknownSpellingsThrough(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func(string) string
	}{
		{name: "radio band", fn: showRadioBand},
		{name: "client band", fn: showClientBand},
		{name: "radio admin", fn: showRadioAdmin},
		{name: "radio oper", fn: showRadioOper},
		{name: "ap admin", fn: showAPAdmin},
		{name: "ap state", fn: showAPState},
		{name: "power type", fn: showPowerType},
		{name: "power mode", fn: showPowerMode},
		{name: "tag source", fn: showTagSource},
		{name: "radio mode", fn: showRadioMode},
		{name: "client phy", fn: showClientPHY},
		{name: "client state", fn: showClientState},
		{name: "p2p block", fn: showP2PBlock},
		{name: "ft mode", fn: showFTMode},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.fn("some-future-member"); got != "some-future-member" {
				t.Errorf("an unknown spelling became %q", got)
			}

			if got := tt.fn(""); got != "" {
				t.Errorf("an empty value became %q", got)
			}
		})
	}
}

// 17.18 adds four radio modes and 17.15 the three 11be client PHY types. Each must render as
// something, and the display must not collide with an older member's.
func TestEnumCoversTheNewerReleaseMembers(t *testing.T) {
	t.Parallel()

	added := map[string]func(string) string{
		"radio-mode-wgb-uplink":     showRadioMode,
		"radio-mode-uwgb":           showRadioMode,
		"radio-mode-urwb":           showRadioMode,
		"radio-mode-wgb-scan":       showRadioMode,
		"client-dot11be-5ghz-prot":  showClientPHY,
		"client-dot11be-24ghz-prot": showClientPHY,
		"client-dot11be-6ghz-prot":  showClientPHY,
	}

	for member, fn := range added {
		if got := fn(member); got == member {
			t.Errorf("%s is not in the display table", member)
		}
	}
}

func TestAPModeWithSub(t *testing.T) {
	t.Parallel()

	tests := []struct{ mode, sub, want string }{
		// Member zero of the mode domain is spelt local-mode while the rest are mode-*,
		// so no prefix rule produces this.
		{"local-mode", "ap-sub-mode-none", "Local"},
		{"mode-monitor", "wips-mode", "Monitor (WIPS)"},
		{"mode-flex-connect", "ap-sub-mode-none", "FlexConnect"},
		{"mode-monitor", "", "Monitor"},
		{"", "wips-mode", ""},
		{"mode-sensor", "local-network", "Sensor (Local-Network)"},
	}

	for _, tt := range tests {
		t.Run(tt.mode+"/"+tt.sub, func(t *testing.T) {
			t.Parallel()

			if got := showAPModeWithSub(tt.mode, tt.sub); got != tt.want {
				t.Errorf("showAPModeWithSub(%q,%q) = %q, want %q", tt.mode, tt.sub, got, tt.want)
			}
		})
	}
}

func TestOutcome(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		total, failed, degraded int
		want                    error
	}{
		{name: "clean", total: 2, want: nil},
		{name: "one of two failed", total: 2, failed: 1, want: ErrPartial},
		{name: "a degraded read", total: 1, degraded: 1, want: ErrPartial},
		{name: "every controller failed", total: 2, failed: 2, want: ErrAllFailed},
		{name: "the only controller failed", total: 1, failed: 1, want: ErrAllFailed},
		{name: "no controller at all", total: 0, want: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := outcome(tt.total, tt.failed, tt.degraded)
			if (tt.want == nil) != (got == nil) || (tt.want != nil && !errors.Is(got, tt.want)) {
				t.Errorf("outcome(%d,%d,%d) = %v, want %v", tt.total, tt.failed, tt.degraded, got, tt.want)
			}
		})
	}
}

func TestReporter(t *testing.T) {
	t.Parallel()

	rep := &Reporter{target: target}

	// A nil error and a zero count are not events.
	rep.Degraded("endpoint", nil)
	rep.Excluded(0, "nothing")
	rep.Note("")

	if rep.Degradations() != 0 || len(rep.notes) != 0 {
		t.Errorf("a non-event was recorded: %d faults, %#v notes", rep.Degradations(), rep.notes)
	}

	rep.Excluded(1, "the band was not reported")
	rep.Excluded(2, "the band was not reported")

	if len(rep.notes) != 2 || !strings.HasPrefix(rep.notes[0], "1 row ") ||
		!strings.HasPrefix(rep.notes[1], "2 rows ") {
		t.Errorf("notes = %#v", rep.notes)
	}
}

// A Monitor or Sniffer radio must read as a radio doing its job and not as a broken serving one.
// The channel is absent because the controller guards that leaf on the mode, and the width and
// power are values because it does not; the device's own summary renders all three as N/A, and
// discarding a reported value is the mirror image of inventing one.
func TestOverviewRowsMonitorAndSnifferRadios(t *testing.T) {
	t.Parallel()

	power := int8(20)
	zero := 0

	radios := []wnc.Radio{
		{
			APName: "lab-ap-1", APMAC: "00:00:5e:00:53:01", Slot: 0,
			Mode: "radio-mode-sniffer", Band: "dot11-2-dot-4-ghz-band",
			AdminState: "enabled", OperState: "radio-up",
			Width: ptr(20), TxPowerDBm: &power, Clients: &zero, ChUtil: &zero,
		},
		{
			APName: "lab-ap-2", APMAC: "00:00:5e:00:53:02", Slot: 1,
			Mode: "radio-mode-monitor", Band: "dot11-5-ghz-band",
			AdminState: "enabled", OperState: "radio-up",
			Width: ptr(40), TxPowerDBm: &power, Clients: &zero, ChUtil: &zero,
		},
	}

	rows := overviewRows(radios, "", target, &Reporter{})
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	for i, want := range []string{"Sniffer", "Monitor"} {
		cells := cellsOf(OverviewColumns(), rows[i])

		if cells["mode"] != want {
			t.Errorf("row %d mode = %q, want %q", i, cells["mode"], want)
		}

		// Without this the row would claim channel 0, which is not a channel.
		if cells["channel"] != render.Absent {
			t.Errorf("row %d channel = %q, want %q", i, cells["channel"], render.Absent)
		}

		if rows[i].Channel != nil {
			t.Errorf("row %d kept a channel of %d", i, *rows[i].Channel)
		}

		// These two carry no guard in the schema, so suppressing them would discard a
		// value the controller sent.
		if cells["channel_width_mhz"] == render.Absent || cells["tx_power_dbm"] == render.Absent {
			t.Errorf("row %d dropped a reported value: width=%q power=%q",
				i, cells["channel_width_mhz"], cells["tx_power_dbm"])
		}

		// A serving radio with no clients and an idle channel reads the same as one of
		// these, which is exactly why the mode column exists.
		if cells["clients"] != "0clients" || cells["ch_util_percent"] != "0%" {
			t.Errorf("row %d lost a reported zero: clients=%q util=%q",
				i, cells["clients"], cells["ch_util_percent"])
		}
	}
}

// The row builder must keep an unjoined access point and must not turn the epoch into
// an age: those are the two readings that make the join view worth having.
func TestAPJoinRows(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC)
	joins := []wnc.APJoin{
		{
			Name: "lab-ap-1", WtpMAC: "00:00:5e:00:53:01", EthernetMAC: "00:00:5e:00:53:11",
			IPAddr: "192.168.255.11", Joined: ptr(true),
			LastFailurePhase: "ap-con-failure-run", LastJoinFailure: "jf-none",
			LastConfigFailure: "cf-none", LastDiscFailure: "disc-fail-none",
			DisconnectReason: "DTLS close alert from peer", RebootReason: "ap-reboot-reason-none",
			LastJoin: now.Add(-2 * time.Hour), LastConfig: now.Add(-2 * time.Hour),
		},
		{
			Name: "lab-ap-2", WtpMAC: "00:00:5e:00:53:02", Joined: ptr(false),
			LastFailurePhase: "ap-con-failure-imgdwnld",
			DisconnectReason: "Mode change to sniffer",
		},
	}

	rows := apJoinRows(joins, now, config.Target{Name: "lab"})
	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	first, second := rows[0], rows[1]

	if got := render.StrPtr(first.Status); got != "Joined" {
		t.Errorf("status = %q, want the device's own word", got)
	}

	if got := render.StrPtr(second.Status); got != "Not Joined" {
		t.Errorf("status = %q, want the device's own word", got)
	}

	// The healthy members render as one word rather than as their raw spellings.
	if got := render.StrPtr(first.LastFailurePhase); got != "Run" {
		t.Errorf("last_failure_phase = %q, want Run", got)
	}

	if got := render.StrPtr(first.LastJoinFailure); got != "None" {
		t.Errorf("last_join_failure = %q, want None", got)
	}

	// "Image-Download" is an expansion of the spelling, which is why no prefix rule
	// can produce this table.
	if got := render.StrPtr(second.LastFailurePhase); got != "Image-Download" {
		t.Errorf("last_failure_phase = %q, want Image-Download", got)
	}

	// The free text is the device's own and is never mapped.
	if got := render.StrPtr(first.DisconnectReason); got != "DTLS close alert from peer" {
		t.Errorf("disconnect_reason = %q, want the free text verbatim", got)
	}

	if got := render.Duration(first.LastJoin); got != "2h" {
		t.Errorf("last_join_seconds = %q, want 2h", got)
	}

	// The zero instant is the never-happened reading, and it must stay absent.
	for name, got := range map[string]string{
		"last_join_seconds":      render.Duration(second.LastJoin),
		"last_config_seconds":    render.Duration(second.LastConfig),
		"last_discovery_seconds": render.Duration(second.LastDiscovery),
		"last_error_seconds":     render.Duration(second.LastError),
	} {
		if got != render.Absent {
			t.Errorf("%s = %q, want %q", name, got, render.Absent)
		}
	}

	// An unknown spelling passes through so a release that adds a member stays readable.
	// The discovery spelling is one of the two 17.15 appended to a domain 17.12 ends at
	// 16, so this case is a release fact rather than an invented member.
	future := apJoinRows([]wnc.APJoin{{
		LastFailurePhase: "ap-con-failure-quantum",
		LastDiscFailure:  "disc-fail-req-migr-off-disabled",
	}}, now, config.Target{})
	if got := render.StrPtr(future[0].LastFailurePhase); got != "ap-con-failure-quantum" {
		t.Errorf("an unknown member rendered as %q, want it passed through", got)
	}

	if got := render.StrPtr(future[0].LastDiscFailure); got != "disc-fail-req-migr-off-disabled" {
		t.Errorf("the 17.15 discovery member rendered as %q, want it passed through", got)
	}
}

// The glyphs are a table style, so they must not assert a state the controller never
// gave, and they must not reach the sort key or the JSON.
func TestPrettyCellsNeverInventAState(t *testing.T) {
	t.Parallel()

	t.Run("nil and the empty string stay absent", func(t *testing.T) {
		t.Parallel()

		if got := prettyBool(nil, glyphOK, glyphOff); got != render.Absent {
			t.Errorf("prettyBool(nil) = %q, want %q", got, render.Absent)
		}

		if got := prettyState(nil, dispUp, dispDown, glyphBad); got != render.Absent {
			t.Errorf("prettyState(nil) = %q, want %q", got, render.Absent)
		}

		// A non-nil pointer to "" is what the plain cell renders as absent, so the
		// glyph path has to agree with it rather than print a blank cell.
		if got := prettyState(ptr(""), dispUp, dispDown, glyphBad); got != render.Absent {
			t.Errorf("prettyState(ptr(%q)) = %q, want %q", "", got, render.Absent)
		}

		if got := prettyOtherwise(nil, dispRegistered, glyphWarn); got != render.Absent {
			t.Errorf("prettyOtherwise(nil) = %q, want %q", got, render.Absent)
		}

		// This one folds every reported value it does not recognize, so the empty string
		// is the case that would turn "not reported" into a warning.
		if got := prettyOtherwise(ptr(""), dispRegistered, glyphWarn); got != render.Absent {
			t.Errorf("prettyOtherwise(ptr(%q)) = %q, want %q", "", got, render.Absent)
		}
	})

	t.Run("an unknown spelling passes through", func(t *testing.T) {
		t.Parallel()

		if got := prettyState(ptr("Suspended"), dispUp, dispDown, glyphBad); got != "Suspended" {
			t.Errorf("prettyState = %q, want the spelling passed through", got)
		}
	})

	// The overview's Oper column is the one with an inverted-free mapping, and ap-tag's
	// Misconfigured is the one where true is the fault.
	t.Run("the polarity is per column", func(t *testing.T) {
		t.Parallel()

		if got := prettyBool(ptr(true), glyphBad, glyphOK); got != glyphBad {
			t.Errorf("a flagged misconfiguration rendered %q, want the cross", got)
		}

		if got := prettyBool(ptr(false), glyphBad, glyphOK); got != glyphOK {
			t.Errorf("a healthy access point rendered %q, want the check", got)
		}
	})
}

// prettyOf renders one row the way the bordered table does, falling back to the plain
// cell exactly as render.PrettyTable does for a column that declares no glyph.
func prettyOf[R any](cols []render.Column[R], row R) map[string]string {
	out := make(map[string]string, len(cols))

	for _, c := range cols {
		if c.Pretty != nil {
			out[c.Key] = c.Pretty(row)

			continue
		}

		out[c.Key] = c.Cell(row)
	}

	return out
}

// show ap's two glyph columns are asymmetric: Admin has two members and names both,
// while State has six and names only the serving one. The raw spellings go in so the
// enum lookup and the glyph mapping are exercised as one path.
func TestAPAdminAndStateGlyphs(t *testing.T) {
	t.Parallel()

	cols := APColumns()

	tests := []struct {
		name      string
		admin     string
		operState string
		wantAdmin string
		wantState string
	}{
		{"a serving access point", "adminstate-enabled", "registered", glyphOK, glyphOK},
		{"administratively disabled", "adminstate-disabled", "registered", glyphNo, glyphOK},
		{"joined but not registered", "adminstate-enabled", "ap-up", glyphOK, glyphWarn},
		{"down", "adminstate-enabled", "ap-down", glyphOK, glyphWarn},
		// Downloading is expected during an image upgrade and Unregistered is not, but
		// the rule folds both: --pretty says the row wants a look, and dropping the flag
		// says why.
		{"mid image download", "adminstate-enabled", "downloading", glyphOK, glyphWarn},
		{"unregistered", "adminstate-enabled", "unregistered", glyphOK, glyphWarn},
		// A member neither table knows keeps the raw spelling under Admin, because that
		// column names both of its members and cannot tell which side a third belongs
		// to. State has no such doubt: anything but Registered is not serving.
		{"a future release", "adminstate-suspended", "ap-quiescing", "adminstate-suspended", glyphWarn},
		{"nothing reported", "", "", render.Absent, render.Absent},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rows := apRows([]wnc.AP{{AdminState: tc.admin, OperState: tc.operState}}, target)
			cells := prettyOf(cols, rows[0])

			if cells[keyAdmin] != tc.wantAdmin {
				t.Errorf("admin = %q, want %q", cells[keyAdmin], tc.wantAdmin)
			}

			if cells[keyState] != tc.wantState {
				t.Errorf("state = %q, want %q", cells[keyState], tc.wantState)
			}
		})
	}
}

// A glyph has to measure for tablewriter what the terminal draws, or every bordered
// row drifts past the cell holding it. Two properties give that, and this asserts both
// with the same function that sizes the column: one code point, because a U+FE0F
// selector asks for a two-column emoji rendering the measurement does not see, and a
// width that does not move with the reader's locale, because tablewriter widens the
// East Asian ambiguous range when it detects one.
func TestGlyphsMeasureTheSameWidthEverywhere(t *testing.T) {
	t.Parallel()

	for _, g := range []string{glyphOK, glyphBad, glyphOff, glyphNo, glyphWarn} {
		if n := utf8.RuneCountInString(g); n != 1 {
			t.Errorf("glyph %q is %d code points, want 1", g, n)
		}

		narrow := twwidth.WidthWithOptions(g, twwidth.Options{EastAsianWidth: false})
		wide := twwidth.WidthWithOptions(g, twwidth.Options{EastAsianWidth: true})

		if narrow != wide {
			t.Errorf("glyph %q measures %d columns narrow and %d wide", g, narrow, wide)
		}
	}
}

// Sort reads Cell and never Pretty, or --pretty would reorder the rows: the glyphs
// sort by code point, which is nothing like the words they replace.
func TestSortIgnoresThePrettyRendering(t *testing.T) {
	t.Parallel()

	rows := []OverviewRow{
		{APName: ptr("a"), Oper: ptr(dispUp)},
		{APName: ptr("b"), Oper: ptr(dispDown)},
		{APName: ptr("c"), Oper: nil},
	}

	cols := OverviewColumns()

	var declaresPretty bool

	for _, c := range cols {
		if c.Key == "oper" && c.Pretty != nil {
			declaresPretty = true
		}
	}

	if !declaresPretty {
		t.Fatal("the oper column declares no Pretty rendering")
	}

	if err := render.Sort(rows, cols, "oper", false); err != nil {
		t.Fatalf("Sort: %v", err)
	}

	got := make([]string, 0, len(rows))
	for _, r := range rows {
		got = append(got, render.StrPtr(r.APName))
	}

	// Down before Up alphabetically, absence last. By glyph it would be a, b, c.
	if want := []string{"b", "a", "c"}; !slicesEqual(got, want) {
		t.Errorf("order = %v, want %v", got, want)
	}
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
