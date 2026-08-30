package show

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// Security display strings that are not built from an AKM name.
const (
	securityOpen = "Open"
	securityWEP  = "WEP"
	securityOSEN = "OSEN"
)

// Policy profile status display strings, named so the Pretty mapping and the cell cannot drift.
const (
	policyActive   = "Active"
	policyShutdown = "Shutdown"
)

// WLANRow is one row of show wlan: one WLAN paired with one policy profile it is bound to, so the
// same WLAN bound under two tags to two profiles is two rows.
type WLANRow struct {
	WLANID         *int    `json:"wlan_id,omitzero"`
	Profile        *string `json:"profile,omitzero"`
	SSID           *string `json:"ssid,omitzero"`
	Status         *string `json:"status,omitzero"`
	Security       *string `json:"security,omitzero"`
	Bands          *string `json:"bands,omitzero"`
	Broadcast      *string `json:"broadcast,omitzero"`
	P2PBlock       *string `json:"p2p_block,omitzero"`
	PolicyStatus   *string `json:"policy_status,omitzero"`
	Switching      *string `json:"switching,omitzero"`
	Interface      *string `json:"interface,omitzero"`
	SessionTimeout *int    `json:"session_timeout_seconds,omitzero"`
	DHCPRequired   *bool   `json:"dhcp_required,omitzero"`
	PolicyProfile  *string `json:"policy_profile,omitzero"`
	Tags           *string `json:"tags,omitzero"`
	Controller     string  `json:"controller"`
}

// WLANColumns describes the WLAN view. Policy Status is separate from Status because a WLAN can be
// enabled while the profile bound to it is shut, in which case no radio carries it.
func WLANColumns() []render.Column[WLANRow] {
	return []render.Column[WLANRow]{
		{
			Key: DefaultSortWLANID, Header: "ID",
			Cell: func(r WLANRow) string { return render.IntPtr(r.WLANID) },
			Sort: func(r WLANRow) any { return render.SortValue(r.WLANID) },
		},
		{Key: "profile", Header: "Profile", Cell: func(r WLANRow) string { return render.StrPtr(r.Profile) }},
		{Key: keySSID, Header: "SSID", Cell: func(r WLANRow) string { return render.StrPtr(r.SSID) }},
		{
			Key: keyStatus, Header: "Status",
			Cell: func(r WLANRow) string { return render.StrPtr(r.Status) },
			Pretty: func(r WLANRow) string {
				return prettyState(r.Status, dispEnabled, dispDisabled, glyphBad)
			},
		},
		{Key: "security", Header: "Security", Cell: func(r WLANRow) string { return render.StrPtr(r.Security) }},
		{Key: "bands", Header: "Bands", Cell: func(r WLANRow) string { return render.StrPtr(r.Bands) }},
		{
			Key: "broadcast", Header: "Broadcast",
			Cell: func(r WLANRow) string { return render.StrPtr(r.Broadcast) },
			// A hidden SSID is a design choice rather than a fault, so it takes the square.
			Pretty: func(r WLANRow) string {
				return prettyState(r.Broadcast, dispEnabled, dispDisabled, glyphOff)
			},
		},
		{Key: "p2p_block", Header: "P2P Block", Cell: func(r WLANRow) string { return render.StrPtr(r.P2PBlock) }},
		{
			Key:    "policy_status",
			Header: "Policy Status",
			Cell:   func(r WLANRow) string { return render.StrPtr(r.PolicyStatus) },
			Pretty: func(r WLANRow) string {
				return prettyState(r.PolicyStatus, policyActive, policyShutdown, glyphBad)
			},
		},
		{Key: "switching", Header: "Switching", Cell: func(r WLANRow) string { return render.StrPtr(r.Switching) }},
		{Key: "interface", Header: "Interface", Cell: func(r WLANRow) string { return render.StrPtr(r.Interface) }},
		{
			Key: "session_timeout_seconds", Header: "Session TO",
			Cell: func(r WLANRow) string { return render.IntPtr(r.SessionTimeout) },
			Sort: func(r WLANRow) any { return render.SortValue(r.SessionTimeout) },
		},
		{
			Key: "dhcp_required", Header: "DHCP Required",
			Cell: func(r WLANRow) string { return render.Bool(r.DHCPRequired) },
			Pretty: func(r WLANRow) string {
				return prettyBool(r.DHCPRequired, glyphOK, glyphOff)
			},
			Sort: func(r WLANRow) any {
				if r.DHCPRequired == nil {
					return nil
				}

				return *r.DHCPRequired
			},
		},
		{
			Key:    keyPolicyProfile,
			Header: headPolicyProfile,
			Cell:   func(r WLANRow) string { return render.StrPtr(r.PolicyProfile) },
		},
		{Key: "tags", Header: "Tags", Cell: func(r WLANRow) string { return render.StrPtr(r.Tags) }},
		{Key: keyController, Header: headController, Cell: func(r WLANRow) string { return render.Str(r.Controller) }},
	}
}

func FetchWLANs(ctx context.Context, c *wnc.Client, t config.Target, rep *Reporter) ([]WLANRow, error) {
	view, reads, err := c.WLANs(ctx)
	if err != nil {
		return nil, err
	}

	rep.Degraded("wlan-policies", reads.Profiles)
	rep.Degraded("policy-list-entries", reads.Bindings)

	return wlanRows(view, t, rep), nil
}

// wlanRows pairs each WLAN with the policy profiles bound to it. A binding may name a WLAN profile
// that has no configuration entry, which the model permits; those bindings produce no row, so
// their count is reported rather than hidden.
func wlanRows(view wnc.WLANView, t config.Target, rep *Reporter) []WLANRow {
	bound, dangling := groupBindings(view)

	rows := make([]WLANRow, 0, len(view.Entries))

	for _, entry := range view.Entries {
		pairs := bound[entry.ProfileName]
		if len(pairs) == 0 {
			// An unbound WLAN still exists; only the policy half of the row is unreported.
			rows = append(rows, wlanRow(entry, wnc.PolicyProfile{}, false, nil, t))

			continue
		}

		for _, policy := range sortedKeys(pairs) {
			profile, known := view.Profiles[policy]
			profile.Name = policy

			rows = append(rows, wlanRow(entry, profile, known, pairs[policy], t))
		}
	}

	if dangling > 0 {
		rep.Note(fmt.Sprintf(
			"%d policy-tag binding(s) name a WLAN profile that does not exist on this controller", dangling))
	}

	return rows
}

// groupBindings indexes the bindings by WLAN profile and then by policy profile, collecting the
// tag names that produced each pair, and counts the bindings whose WLAN profile has no entry.
func groupBindings(view wnc.WLANView) (bound map[string]map[string][]string, dangling int) {
	known := make(map[string]bool, len(view.Entries))
	for _, e := range view.Entries {
		known[e.ProfileName] = true
	}

	bound = make(map[string]map[string][]string, len(view.Entries))

	for _, b := range view.Bindings {
		if !known[b.WLANProfile] {
			dangling++

			continue
		}

		if bound[b.WLANProfile] == nil {
			bound[b.WLANProfile] = make(map[string][]string, 1)
		}

		bound[b.WLANProfile][b.Policy] = append(bound[b.WLANProfile][b.Policy], b.Tag)
	}

	return bound, dangling
}

// sortedKeys orders the policy profiles of one WLAN so two runs produce the same rows.
func sortedKeys(m map[string][]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}

	slices.Sort(out)

	return out
}

// wlanRow builds one row. known says whether the policy profile was read, so a binding to one the
// profile list did not return leaves the policy cells unreported.
func wlanRow(
	e wnc.WLANEntry, p wnc.PolicyProfile, known bool, tags []string, t config.Target,
) WLANRow {
	row := WLANRow{
		WLANID:     ptr(e.WLANID),
		Profile:    optional(e.ProfileName),
		SSID:       optional(e.APFVAPIDData.SSID),
		Status:     enabledDisabled(e.APFVAPIDData.WLANStatus),
		Security:   optional(security(e)),
		Bands:      optional(bandList(e.Bands())),
		Broadcast:  enabledDisabled(e.APFVAPIDData.BroadcastSSID),
		P2PBlock:   optionalPtr(e.APFVAPIDData.P2PBlock, showP2PBlock),
		Controller: t.Name,
	}

	if p.Name != "" {
		row.PolicyProfile = optional(p.Name)
		row.Tags = optional(strings.Join(tags, ","))
	}

	if known {
		row.PolicyStatus = ptr(map[bool]string{true: policyShutdown, false: policyActive}[p.Shutdown])
		row.Interface = optional(p.InterfaceName)
		row.SessionTimeout = p.SessionTimeout
		row.DHCPRequired = p.DHCPRequired
		row.Switching = switching(p.CentralSwitching)
	}

	return row
}

// security names the WLAN's link-layer security, and the order is what makes it right.
// security-wpa is the master switch and defaults to true, so it is read first; WEP and OSEN each
// give a WLAN with every AKM leaf false that is nonetheless not open, so they precede the AKM set,
// and shared-key authentication decides one whose AKM set came out empty. The FT suffix comes from
// the FT-specific AKM leaves and never from ft-mode, whose default is the adaptive member.
func security(e wnc.WLANEntry) string {
	if isFalse(e.SecurityWPA) {
		return securityOpen
	}

	if isTrue(e.WEPEnabled) {
		return securityWEP
	}

	if isTrue(e.OSEN) {
		return securityOSEN
	}

	versions := wpaVersions(e)
	akms := akmList(e)

	if len(versions) == 0 && len(akms) == 0 {
		if e.DOT11AuthType != nil && *e.DOT11AuthType == "apf-vap-80211-auth-shared-key" {
			return securityWEP
		}

		return securityOpen
	}

	suite := strings.Join(versions, "/")
	if len(akms) > 0 {
		suite = strings.TrimPrefix(suite+" "+strings.Join(akms, "+"), " ")
	}

	if ftEnabled(e) {
		suite += "+FT"
	}

	return suite
}

func wpaVersions(e wnc.WLANEntry) []string {
	var out []string

	for _, v := range []struct {
		on   *bool
		name string
	}{
		{e.WPA1Enabled, "WPA"},
		{e.WPA2Enabled, "WPA2"},
		{e.WPA3Enabled, "WPA3"},
	} {
		if isTrue(v.on) {
			out = append(out, v.name)
		}
	}

	return out
}

// akmList names the key-management methods enabled, FT variants excluded because those become the
// suffix instead.
func akmList(e wnc.WLANEntry) []string {
	var out []string

	for _, a := range []struct {
		on   *bool
		name string
	}{
		{e.AKMSAE, "SAE"},
		{e.AKMSAEExtKey, "SAE-EXT-KEY"},
		{e.AKMOWE, "OWE"},
		{e.AKMDot1x, "802.1X"},
		{e.AKMDot1xSHA256, "802.1X-SHA256"},
		{e.AKMSuiteB, "SuiteB-1X"},
		{e.AKMSuiteB192, "SuiteB192-1X"},
		{e.AKMPSK, "PSK"},
		{e.AKMPSKSHA256, "PSK-SHA256"},
		{e.AKMCCKM, "CCKM"},
	} {
		if isTrue(a.on) {
			out = append(out, a.name)
		}
	}

	return out
}

func ftEnabled(e wnc.WLANEntry) bool {
	return isTrue(e.AKMFTPSK) || isTrue(e.AKMFTDot1x) || isTrue(e.AKMFTSAE) || isTrue(e.AKMFTSAEExtKey)
}

// isTrue and isFalse read an optional boolean without treating absence as either.
func isTrue(p *bool) bool  { return p != nil && *p }
func isFalse(p *bool) bool { return p != nil && !*p }

// enabledDisabled renders an optional boolean as a state rather than a yes or no, as the
// controller's own WLAN view does.
func enabledDisabled(p *bool) *string {
	if p == nil {
		return nil
	}

	if *p {
		return ptr(dispEnabled)
	}

	return ptr(dispDisabled)
}

// switching renders the central-switching flag. Its schema default is true, so an
// absent value is not Local.
func switching(p *bool) *string {
	if p == nil {
		return nil
	}

	if *p {
		return ptr("Central")
	}

	return ptr("Local")
}

func bandList(bands []string) string {
	if len(bands) == 0 {
		return ""
	}

	out := make([]string, 0, len(bands))
	for _, b := range bands {
		out = append(out, showRadioBand(b))
	}

	return strings.Join(out, "/")
}

func optionalPtr(p *string, display func(string) string) *string {
	if p == nil {
		return nil
	}

	return optional(display(*p))
}
