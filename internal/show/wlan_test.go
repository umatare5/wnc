package show

import (
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

func yes() *bool { v := true; return &v }
func no() *bool  { v := false; return &v }

// The security column's failure modes are the worst in this view: calling a WPA3
// enterprise WLAN "Open", or calling a WEP WLAN open because every AKM leaf is false.
func TestSecurity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		entry wnc.WLANEntry
		want  string
	}{
		{
			name:  "the master switch off wins over everything",
			entry: wnc.WLANEntry{SecurityWPA: no(), WPA2Enabled: yes(), AKMPSK: yes()},
			want:  securityOpen,
		},
		{
			// Every AKM leaf is false here, which is exactly why an AKM-only rule would
			// report this as open.
			name:  "static WEP",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WEPEnabled: yes()},
			want:  securityWEP,
		},
		{
			name:  "OSEN",
			entry: wnc.WLANEntry{SecurityWPA: yes(), OSEN: yes()},
			want:  securityOSEN,
		},
		{
			name:  "shared-key authentication with no AKM",
			entry: wnc.WLANEntry{SecurityWPA: yes(), DOT11AuthType: strptr("apf-vap-80211-auth-shared-key")},
			want:  securityWEP,
		},
		{
			name:  "WPA2 personal",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WPA2Enabled: yes(), AKMPSK: yes()},
			want:  "WPA2 PSK",
		},
		{
			name:  "WPA3 enterprise",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WPA3Enabled: yes(), AKMDot1xSHA256: yes()},
			want:  "WPA3 802.1X-SHA256",
		},
		{
			name:  "WPA3 personal with fast transition",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WPA3Enabled: yes(), AKMSAE: yes(), AKMFTSAE: yes()},
			want:  "WPA3 SAE+FT",
		},
		{
			// Added after 17.12. Reading only the older AKM leaves would make this Open.
			name:  "WPA3 with the extended-key AKM alone",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WPA3Enabled: yes(), AKMSAEExtKey: yes()},
			want:  "WPA3 SAE-EXT-KEY",
		},
		{
			name:  "Suite-B 192",
			entry: wnc.WLANEntry{SecurityWPA: yes(), WPA3Enabled: yes(), AKMSuiteB192: yes()},
			want:  "WPA3 SuiteB192-1X",
		},
		{
			name: "a transition WLAN carrying two generations",
			entry: wnc.WLANEntry{
				SecurityWPA: yes(),
				WPA2Enabled: yes(),
				WPA3Enabled: yes(),
				AKMPSK:      yes(),
				AKMSAE:      yes(),
			},
			want: "WPA2/WPA3 SAE+PSK",
		},
		{
			name:  "nothing reported",
			entry: wnc.WLANEntry{},
			want:  securityOpen,
		},
		{
			name: "every input present and false",
			entry: wnc.WLANEntry{
				SecurityWPA: yes(),
				WPA1Enabled: no(),
				WPA2Enabled: no(),
				WPA3Enabled: no(),
				WEPEnabled:  no(),
				OSEN:        no(),
			},
			want: securityOpen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := security(tt.entry); got != tt.want {
				t.Errorf("security = %q, want %q", got, tt.want)
			}
		})
	}
}

// ft-mode defaults to the adaptive member, so an FT suffix derived from it would land
// on nearly every WLAN. Only the FT-specific AKM leaves may produce it.
func TestFTSuffixComesFromTheAKMLeaves(t *testing.T) {
	t.Parallel()

	base := wnc.WLANEntry{SecurityWPA: yes(), WPA2Enabled: yes(), AKMPSK: yes()}

	if got := security(base); strings.Contains(got, "FT") {
		t.Errorf("security = %q, want no FT suffix", got)
	}

	withFT := base
	withFT.AKMFTPSK = yes()

	if got := security(withFT); !strings.HasSuffix(got, "+FT") {
		t.Errorf("security = %q, want an FT suffix", got)
	}
}

func TestWLANRowsPairsEachWLANWithItsProfiles(t *testing.T) {
	t.Parallel()

	view := wnc.WLANView{
		Entries: []wnc.WLANEntry{
			entryNamed(5, "wlan-a"),
			entryNamed(6, "wlan-b"),
			entryNamed(7, "wlan-unbound"),
		},
		Profiles: map[string]wnc.PolicyProfile{
			"prof-1": {Name: "prof-1", InterfaceName: "LAB-A", CentralSwitching: no()},
			"prof-2": {Name: "prof-2", Shutdown: true, InterfaceName: "LAB-B"},
		},
		Bindings: []wnc.Binding{
			{Tag: "tag-1", WLANProfile: "wlan-a", Policy: "prof-1"},
			// The same WLAN under a second tag and a second profile is a second row.
			{Tag: "tag-2", WLANProfile: "wlan-a", Policy: "prof-2"},
			// Two tags naming the same pair collapse into one row listing both tags.
			{Tag: "tag-3", WLANProfile: "wlan-b", Policy: "prof-1"},
			{Tag: "tag-4", WLANProfile: "wlan-b", Policy: "prof-1"},
			// A binding naming a WLAN that does not exist produces no row.
			{Tag: "tag-5", WLANProfile: "wlan-missing", Policy: "prof-1"},
		},
	}

	rep := &Reporter{target: config.Target{Name: "test"}}
	rows := wlanRows(view, config.Target{Name: "test"}, rep)

	if len(rows) != 4 {
		t.Fatalf("got %d rows, want 4", len(rows))
	}

	// wlan-a appears twice, once per profile, in profile-name order.
	if got := deref(rows[0].PolicyProfile); got != "prof-1" {
		t.Errorf("first row profile = %q", got)
	}

	if got := deref(rows[1].PolicyProfile); got != "prof-2" {
		t.Errorf("second row profile = %q", got)
	}

	if got := deref(rows[1].PolicyStatus); got != "Shutdown" {
		t.Errorf("a shut profile rendered %q", got)
	}

	// Two tags on one pair are listed together on the single row.
	if got := deref(rows[2].Tags); got != "tag-3,tag-4" {
		t.Errorf("tags = %q", got)
	}

	// The unbound WLAN keeps its row and leaves the policy half unreported.
	last := rows[3]
	if deref(last.Profile) != "wlan-unbound" {
		t.Fatalf("last row = %+v", last)
	}

	if last.PolicyProfile != nil || last.Interface != nil || last.PolicyStatus != nil {
		t.Errorf("an unbound WLAN reported policy values: %+v", last)
	}

	// The dangling binding is reported rather than silently dropped.
	if len(rep.notes) != 1 || !strings.Contains(rep.notes[0], "does not exist") {
		t.Errorf("notes = %#v, want the dangling binding reported", rep.notes)
	}
}

// A binding may name a profile the profile list did not return. The policy cells then
// stay unreported instead of showing an active profile with no interface.
func TestWLANRowsHandlesAnUnknownProfile(t *testing.T) {
	t.Parallel()

	view := wnc.WLANView{
		Entries:  []wnc.WLANEntry{entryNamed(5, "wlan-a")},
		Profiles: map[string]wnc.PolicyProfile{},
		Bindings: []wnc.Binding{{Tag: "tag-1", WLANProfile: "wlan-a", Policy: "prof-x"}},
	}

	rows := wlanRows(view, config.Target{Name: "test"}, &Reporter{})
	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}

	if deref(rows[0].PolicyProfile) != "prof-x" {
		t.Errorf("the bound profile name was lost: %+v", rows[0])
	}

	if rows[0].PolicyStatus != nil || rows[0].Interface != nil {
		t.Errorf("an unread profile produced values: %+v", rows[0])
	}
}

func TestSwitchingAndEnabledDisabled(t *testing.T) {
	t.Parallel()

	if got := switching(nil); got != nil {
		t.Errorf("switching(nil) = %q, want nil", *got)
	}

	if got := deref(switching(yes())); got != "Central" {
		t.Errorf("switching(true) = %q", got)
	}

	if got := deref(switching(no())); got != "Local" {
		t.Errorf("switching(false) = %q", got)
	}

	if got := enabledDisabled(nil); got != nil {
		t.Errorf("enabledDisabled(nil) = %q, want nil", *got)
	}

	if got := deref(enabledDisabled(yes())); got != "Enabled" {
		t.Errorf("enabledDisabled(true) = %q", got)
	}
}

func TestBandList(t *testing.T) {
	t.Parallel()

	if got := bandList(nil); got != "" {
		t.Errorf("bandList(nil) = %q", got)
	}

	got := bandList([]string{"dot11-5-ghz-band", "dot11-6-ghz-band"})
	if got != "5/6" {
		t.Errorf("bandList = %q, want 5/6", got)
	}

	// An unknown spelling passes through so a new release shows the raw value rather
	// than an empty cell.
	if got := bandList([]string{"dot11-8-ghz-band"}); got != "dot11-8-ghz-band" {
		t.Errorf("bandList = %q, want the raw value", got)
	}
}

func entryNamed(id int, profile string) wnc.WLANEntry {
	e := wnc.WLANEntry{WLANID: id, ProfileName: profile}
	e.APFVAPIDData.SSID = profile

	return e
}

func strptr(s string) *string { return &s }

func deref(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}
