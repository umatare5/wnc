package wnc

import (
	"context"
	"fmt"

	sdk "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// The path the WLAN view reads with a struct of its own.
const wlanEntriesPath = "Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-data/wlan-cfg-entries"

// wlanEntries is the envelope GetDataInto decodes into. The outer tag is module-qualified because
// that is the sole top-level key the SDK holds the response to.
type wlanEntries struct {
	Entries struct {
		Entry []WLANEntry `json:"wlan-cfg-entry"`
	} `json:"Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-entries"`
}

// WLANEntry stands in for the SDK's WlanCfgEntry as an allow-list rather than for the types: this
// read asks for the defaults in force, which materializes the psk and the WEP key material, and a
// leaf this struct does not declare is dropped at decode.
type WLANEntry struct {
	WLANID      int    `json:"wlan-id"`
	ProfileName string `json:"profile-name"`

	APFVAPIDData struct {
		SSID          string  `json:"ssid"`
		WLANStatus    *bool   `json:"wlan-status"`
		BroadcastSSID *bool   `json:"broadcast-ssid"`
		P2PBlock      *string `json:"p2p-block-action"`
	} `json:"apf-vap-id-data"`

	SecurityWPA *bool `json:"security-wpa"`
	WPA1Enabled *bool `json:"wpa1-enabled"`
	WPA2Enabled *bool `json:"wpa2-enabled"`
	WPA3Enabled *bool `json:"wpa3-enabled"`
	WEPEnabled  *bool `json:"wep-enabled"`
	OSEN        *bool `json:"osen"`

	AKMPSK         *bool `json:"auth-key-mgmt-psk"`
	AKMPSKSHA256   *bool `json:"auth-key-mgmt-psk-sha256"`
	AKMDot1x       *bool `json:"auth-key-mgmt-dot1x"`
	AKMDot1xSHA256 *bool `json:"auth-key-mgmt-dot1x-sha256"`
	AKMSAE         *bool `json:"auth-key-mgmt-sae"`
	AKMOWE         *bool `json:"akm-owe"`
	AKMCCKM        *bool `json:"auth-key-mgmt-cckm"`

	// Added after 17.12. A WPA3-Enterprise or Suite-B WLAN sets one of these and none
	// of the leaves above, so a view that omits them renders such a WLAN as Open.
	AKMSAEExtKey *bool `json:"akm-sae-ext-key"`
	AKMSuiteB    *bool `json:"auth-key-mgmt-suite-b"`
	AKMSuiteB192 *bool `json:"auth-key-mgmt-suite-b-192"`

	AKMFTPSK       *bool `json:"auth-key-mgmt-ft-psk"`
	AKMFTDot1x     *bool `json:"auth-key-mgmt-ft-dot1x"`
	AKMFTSAE       *bool `json:"auth-key-mgmt-ft-sae"`
	AKMFTSAEExtKey *bool `json:"akm-ft-sae-ext-key"`

	DOT11AuthType *string `json:"dot11-auth-type"`

	// The bands are read from here alone. The legacy scalar radio-policy leaf is marked obsolete on
	// every release in scope and stopped arriving at 17.18 even under with-defaults; where it does
	// arrive it reports "all bands" on a WLAN this list confines to one.
	WLANRadioPolicies *struct {
		Policy []struct {
			Band string `json:"band"`
		} `json:"wlan-radio-policy"`
	} `json:"wlan-radio-policies"`
}

// Bands lists the bands the WLAN is enabled on, in the order the controller sent them.
func (e WLANEntry) Bands() []string {
	if e.WLANRadioPolicies == nil {
		return nil
	}

	out := make([]string, 0, len(e.WLANRadioPolicies.Policy))
	for _, p := range e.WLANRadioPolicies.Policy {
		out = append(out, p.Band)
	}

	return out
}

// PolicyProfile is one policy profile, as far as this view needs it.
type PolicyProfile struct {
	Name             string
	Shutdown         bool
	InterfaceName    string
	CentralSwitching *bool
	SessionTimeout   *int
	DHCPRequired     *bool
}

// Binding is one policy tag naming one WLAN and the profile it is bound to.
type Binding struct {
	Tag         string
	WLANProfile string
	Policy      string
}

type WLANView struct {
	Entries  []WLANEntry
	Profiles map[string]PolicyProfile
	Bindings []Binding
}

// WLANReads reports which secondary read failed.
type WLANReads struct {
	Profiles error
	Bindings error
}

// WLANs reads the WLAN view; the configuration entries drive the rows. GetDataInto carries the
// allow-list struct in place of the SDK's typed accessor, and asks for the defaults in force
// because security-wpa, wpa2-enabled and auth-key-mgmt-dot1x default to true, so a plain read
// would report them off.
func (c *Client) WLANs(ctx context.Context) (WLANView, WLANReads, error) {
	entries, err := sdk.GetDataInto[wlanEntries](ctx, c.sdk, wlanEntriesPath, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return WLANView{}, WLANReads{}, fmt.Errorf("reading wlan-cfg-entries: %w", err)
	}

	profiles, profErr := c.policyProfiles(ctx)
	bindings, bindErr := c.policyBindings(ctx)

	view := WLANView{Entries: entries.Entries.Entry, Profiles: profiles, Bindings: bindings}

	return view, WLANReads{Profiles: profErr, Bindings: bindErr}, nil
}

// policyProfiles reads the policy profiles through the typed accessor, every leaf the view needs
// being declared. It still asks for the defaults in force, because a plain read omits the interface
// name, the session timeout and the DHCP flag on any profile that never set them, and the interface
// name's default is the string "1" rather than the empty string a non-pointer field would yield.
func (c *Client) policyProfiles(ctx context.Context) (map[string]PolicyProfile, error) {
	resp, err := c.sdk.WLAN().ListWlanPolicies(ctx, sdk.WithDefaults(sdk.ReportAll))
	if err != nil {
		return nil, fmt.Errorf("reading wlan-policies: %w", err)
	}

	if resp == nil || resp.WlanPolicies == nil {
		return nil, nil
	}

	out := make(map[string]PolicyProfile, len(resp.WlanPolicies.WlanPolicy))

	for _, p := range resp.WlanPolicies.WlanPolicy {
		profile := PolicyProfile{
			Name:          p.PolicyProfileName,
			Shutdown:      !p.Status,
			InterfaceName: p.InterfaceName,
		}

		if p.WlanSwitchingPolicy != nil {
			profile.CentralSwitching = p.WlanSwitchingPolicy.CentralSwitching
		}

		if p.WlanTimeout != nil {
			profile.SessionTimeout = p.WlanTimeout.SessionTimeout
		}

		if p.DHCPParams != nil {
			profile.DHCPRequired = ptrTo(p.DHCPParams.IsDHCPEnabled)
		}

		out[p.PolicyProfileName] = profile
	}

	return out, nil
}

// policyBindings reads the policy tags and flattens their WLAN-to-profile bindings. A tag with no
// bindings sends no container at all.
func (c *Client) policyBindings(ctx context.Context) ([]Binding, error) {
	resp, err := c.sdk.WLAN().ListCfgPolicyListEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading policy-list-entries: %w", err)
	}

	if resp == nil || resp.PolicyListEntries == nil {
		return nil, nil
	}

	var out []Binding

	for _, tag := range resp.PolicyListEntries.PolicyListEntry {
		if tag.WLANPolicies == nil {
			continue
		}

		for _, b := range tag.WLANPolicies.WLANPolicy {
			out = append(out, Binding{
				Tag:         tag.TagName,
				WLANProfile: b.WLANProfileName,
				Policy:      b.PolicyProfileName,
			})
		}
	}

	return out, nil
}
