package show

import (
	"encoding/json/v2"
	"strings"
	"testing"

	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// A tag binding three WLANs is three rows and a tag binding none is one, so a tag that
// exists and binds nothing is still visible to the operator about to delete it.
func TestPolicyTagRowsExpandOverTheBindings(t *testing.T) {
	t.Parallel()

	desc := "lab flex"
	tags := []wnc.PolicyTag{
		{
			Name:        "labo-wlan-flex",
			Description: &desc,
			Bindings: []wnc.PolicyBinding{
				{WLANProfile: "labo-p736b2", PolicyProfile: "labo-wlan-profile"},
				{WLANProfile: "labo-p736b5", PolicyProfile: "labo-wlan-profile"},
			},
		},
		{Name: "default-policy-tag"},
	}

	rows := policyTagRows(tags, target)

	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}

	first := cellsOf(PolicyTagColumns(), rows[0])
	if first["policy_tag"] != "labo-wlan-flex" || first["wlan"] != "labo-p736b2" {
		t.Errorf("first row = %+v", first)
	}

	if first["policy_profile"] != "labo-wlan-profile" || first["description"] != "lab flex" {
		t.Errorf("first row = %+v", first)
	}

	// The description repeats on every row of one tag: each row is a whole reading of the
	// tag, so a consumer filtering the JSON on one binding keeps it.
	if got := cellsOf(PolicyTagColumns(), rows[1])["description"]; got != "lab flex" {
		t.Errorf("second row description = %q", got)
	}

	unbound := cellsOf(PolicyTagColumns(), rows[2])
	if unbound["policy_tag"] != "default-policy-tag" {
		t.Errorf("the unbound tag lost its row: %+v", unbound)
	}

	for _, key := range []string{"description", "wlan", "policy_profile"} {
		if unbound[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, unbound[key], render.Absent)
		}
	}
}

// A leaf the controller sent empty must not keep the JSON key while the table cell reads
// as a dash: the two outputs would then disagree about whether anything was reported.
func TestTagRowsCollapseAReportedEmptyString(t *testing.T) {
	t.Parallel()

	empty := ""
	rows := policyTagRows([]wnc.PolicyTag{{Name: "labo-empty", Description: &empty}}, target)

	if got := cellsOf(PolicyTagColumns(), rows[0])["description"]; got != render.Absent {
		t.Errorf("an empty description rendered %q, want %q", got, render.Absent)
	}

	body, err := json.Marshal(rows)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	if strings.Contains(string(body), "description") {
		t.Errorf("the JSON kept a key the table showed as absent: %s", body)
	}
}

// Local Site has three readings and they must stay three: a site tag that is local, one
// that is not, and one the controller said nothing about. The last is why the read asks
// for the defaults in force.
func TestSiteTagRowsKeepTheThreeLocalSiteReadings(t *testing.T) {
	t.Parallel()

	yes, no := true, false
	flex := "labo-flex"
	join := "labo-common"

	tags := []wnc.SiteTag{
		{Name: "labo-site-flex", APJoinProfile: &join, FlexProfile: &flex, LocalSite: &no},
		{Name: "default-site-tag", LocalSite: &yes},
		{Name: "labo-bare"},
	}

	rows := siteTagRows(tags, target)

	for i, want := range []string{"No", "Yes", render.Absent} {
		if got := cellsOf(SiteTagColumns(), rows[i])["local_site"]; got != want {
			t.Errorf("row %d local_site = %q, want %q", i, got, want)
		}
	}

	first := cellsOf(SiteTagColumns(), rows[0])
	if first["ap_join_profile"] != join || first["flex_profile"] != flex {
		t.Errorf("first row = %+v", first)
	}

	bare := cellsOf(SiteTagColumns(), rows[2])
	for _, key := range []string{"description", "ap_join_profile", "flex_profile"} {
		if bare[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, bare[key], render.Absent)
		}
	}
}

// Local Site declares no glyph on purpose: neither reading is a fault or a feature
// switched off, so the bordered table shows the same Yes or No the plain one does.
func TestSiteTagLocalSiteTakesNoGlyph(t *testing.T) {
	t.Parallel()

	yes := true
	rows := siteTagRows([]wnc.SiteTag{{Name: "labo-site", LocalSite: &yes}}, target)

	if got := prettyOf(SiteTagColumns(), rows[0])["local_site"]; got != "Yes" {
		t.Errorf("bordered local_site = %q, want the plain cell", got)
	}
}

// The SDK types every RF tag leaf but the name as a pointer, so an omitted profile arrives nil
// rather than as "", and the row layer is what keeps it out of the output.
func TestRFTagRowsAbsenceRules(t *testing.T) {
	t.Parallel()

	tags := []wnc.RFTag{
		{
			Name: "labo-inside", Description: ptr("lab inside"),
			Profile24GHz: ptr("labo-rf-24gh"), Profile5GHz: ptr("labo-rf-5gh"),
			Profile6GHz: ptr("labo-rf-6gh"),
		},
		{Name: "default-rf-tag"},
	}

	rows := rfTagRows(tags, target)

	first := cellsOf(RFTagColumns(), rows[0])
	if first["profile_24ghz"] != "labo-rf-24gh" || first["profile_5ghz"] != "labo-rf-5gh" {
		t.Errorf("first row = %+v", first)
	}

	second := cellsOf(RFTagColumns(), rows[1])
	if second["rf_tag"] != "default-rf-tag" {
		t.Errorf("the tag lost its name: %+v", second)
	}

	for _, key := range []string{"description", "profile_24ghz", "profile_5ghz", "profile_6ghz"} {
		if second[key] != render.Absent {
			t.Errorf("%s = %q, want %q", key, second[key], render.Absent)
		}
	}
}

// The controller column is the CLI's own and not a leaf, so it is the one cell of these
// views that is never absent.
func TestTagRowsCarryTheController(t *testing.T) {
	t.Parallel()

	policy := cellsOf(PolicyTagColumns(), policyTagRows([]wnc.PolicyTag{{Name: "a"}}, target)[0])
	if policy[keyController] != target.Name {
		t.Errorf("policy-tag controller = %q", policy[keyController])
	}

	site := cellsOf(SiteTagColumns(), siteTagRows([]wnc.SiteTag{{Name: "a"}}, target)[0])
	if site[keyController] != target.Name {
		t.Errorf("site-tag controller = %q", site[keyController])
	}

	rf := cellsOf(RFTagColumns(), rfTagRows([]wnc.RFTag{{Name: "a"}}, target)[0])
	if rf[keyController] != target.Name {
		t.Errorf("rf-tag controller = %q", rf[keyController])
	}
}

// The one absence rule the cells cannot pin. render.Str renders a reported empty string as
// Absent too, so a leaf handed straight through instead of collapsed reads identically in a
// cell and differs only in the JSON — which is where a consumer would read it as a value the
// controller sent. Both shapes are fed in: a leaf the controller omitted, and one it reported
// empty.
func TestTagRowsOmitEveryAbsentColumnFromTheJSON(t *testing.T) {
	t.Parallel()

	empty := ""

	tests := []struct {
		name string
		rows any
		keep []string
		drop []string
	}{
		{
			name: "policy-tag, everything omitted",
			rows: policyTagRows([]wnc.PolicyTag{{Name: "default-policy-tag"}}, target),
			keep: []string{"policy_tag", "controller"},
			drop: []string{"description", "wlan", "policy_profile"},
		},
		{
			name: "policy-tag, everything reported empty",
			rows: policyTagRows([]wnc.PolicyTag{{
				Name: "labo-empty", Description: &empty,
				Bindings: []wnc.PolicyBinding{{}},
			}}, target),
			keep: []string{"policy_tag", "controller"},
			drop: []string{"description", "wlan", "policy_profile"},
		},
		{
			name: "site-tag, everything omitted",
			rows: siteTagRows([]wnc.SiteTag{{Name: "labo-bare"}}, target),
			keep: []string{"site_tag", "controller"},
			drop: []string{"description", "ap_join_profile", "flex_profile", "local_site"},
		},
		{
			name: "site-tag, every string reported empty",
			rows: siteTagRows([]wnc.SiteTag{{
				Name: "labo-empty", Description: &empty, APJoinProfile: &empty, FlexProfile: &empty,
			}}, target),
			keep: []string{"site_tag", "controller"},
			drop: []string{"description", "ap_join_profile", "flex_profile"},
		},
		{
			name: "rf-tag, everything omitted",
			rows: rfTagRows([]wnc.RFTag{{Name: "default-rf-tag"}}, target),
			keep: []string{"rf_tag", "controller"},
			drop: []string{"description", "profile_24ghz", "profile_5ghz", "profile_6ghz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			body, err := json.Marshal(tt.rows)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}

			for _, k := range tt.keep {
				if !strings.Contains(string(body), `"`+k+`"`) {
					t.Errorf("the %s key was dropped: %s", k, body)
				}
			}

			for _, k := range tt.drop {
				if strings.Contains(string(body), `"`+k+`"`) {
					t.Errorf("a value was fabricated for %s: %s", k, body)
				}
			}
		})
	}
}

// The positive control for the test above: an assertion that only ever demands an absence
// would pass on rows that reported nothing at all.
func TestTagRowsKeepEveryReportedColumnInTheJSON(t *testing.T) {
	t.Parallel()

	flex, join := "labo-flex", "labo-common"
	local := false

	policy := policyTagRows([]wnc.PolicyTag{{
		Name:     "labo-wlan-flex",
		Bindings: []wnc.PolicyBinding{{WLANProfile: "labo-p736b2", PolicyProfile: "labo-wlan-profile"}},
	}}, target)

	site := siteTagRows([]wnc.SiteTag{{
		Name: "labo-site-flex", APJoinProfile: &join, FlexProfile: &flex, LocalSite: &local,
	}}, target)

	rf := rfTagRows([]wnc.RFTag{{Name: "labo-inside", Profile5GHz: ptr("labo-rf-5gh")}}, target)

	tests := map[string][]string{
		mustMarshal(t, policy): {`"wlan":"labo-p736b2"`, `"policy_profile":"labo-wlan-profile"`},
		mustMarshal(t, site):   {`"ap_join_profile":"labo-common"`, `"flex_profile":"labo-flex"`, `"local_site":false`},
		mustMarshal(t, rf):     {`"profile_5ghz":"labo-rf-5gh"`},
	}

	for body, wants := range tests {
		for _, want := range wants {
			if !strings.Contains(body, want) {
				t.Errorf("the JSON dropped %s: %s", want, body)
			}
		}
	}
}

func mustMarshal(t *testing.T, v any) string {
	t.Helper()

	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}

	return string(b)
}
