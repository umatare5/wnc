package show

// The key lists below are each command's column order and, at the same time, the vocabulary of its
// --sort-by flag and the field names of its JSON output. json/v2 rejects a duplicated json tag
// outright, so an invariant test asserts a Row's tags are exactly these, in this order.

// Keys more than one command declares. The invariant test compares these against the struct tags,
// which must stay literal.
const (
	keyAPName        = "ap_name"
	keyAPMAC         = "ap_mac"
	keySlot          = "slot"
	keyBand          = "band"
	keyMode          = "mode"
	keyAdmin         = "admin"
	keyChannel       = "channel"
	keySSID          = "ssid"
	keyStatus        = "status"
	keyRadioMAC      = "radio_mac"
	keyEthernetMAC   = "ethernet_mac"
	keyIPAddress     = "ip_address"
	keyState         = "state"
	keyController    = "controller"
	keyDescription   = "description"
	keyPolicyTag     = "policy_tag"
	keySiteTag       = "site_tag"
	keyRFTag         = "rf_tag"
	keyFlexProfile   = "flex_profile"
	keyPolicyProfile = "policy_profile"
)

const (
	headAPName        = "AP Name"
	headRadioMAC      = "Radio MAC"
	headEthernetMAC   = "Ethernet MAC"
	headIPAddress     = "IP Address"
	headController    = "Controller"
	headDescription   = "Description"
	headPolicyTag     = "Policy Tag"
	headSiteTag       = "Site Tag"
	headRFTag         = "RF Tag"
	headFlexProfile   = "Flex Profile"
	headPolicyProfile = "Policy Profile"
)

// Default sort keys, one per command. The three tag views sort by their own key leaf, the one
// column of each the controller cannot omit.
const (
	DefaultSortAPName    = keyAPName
	DefaultSortMAC       = "mac"
	DefaultSortWLANID    = "wlan_id"
	DefaultSortPolicyTag = keyPolicyTag
	DefaultSortSiteTag   = keySiteTag
	DefaultSortRFTag     = keyRFTag
)

// OverviewKeys are the columns of show overview, one row per AP radio.
func OverviewKeys() []string {
	return []string{
		keyAPName, keyAPMAC, keySlot, keyMode, keyBand, keyAdmin, "oper",
		keyChannel, "channel_width_mhz", "tx_power_dbm", "clients",
		"ch_util_percent", "rf_profile", keyController,
	}
}

// APKeys are the columns of show ap, one row per access point.
func APKeys() []string {
	return []string{
		keyAPName, "model", "serial", keyEthernetMAC, keyRadioMAC, keyIPAddress,
		"sw_version", "slots", "country", keyMode, keyAdmin, keyState,
		"lldp_neighbor", "power_type", "power_mode",
		"uptime_seconds", "assoc_uptime_seconds", keyController,
	}
}

// APJoinKeys are the columns of show ap-join, one row per access point the controller remembers,
// joined or not. No counter is a column: nothing on the list declares a clear time, so a total
// since boot has no window.
func APJoinKeys() []string {
	return []string{
		keyAPName, keyRadioMAC, keyEthernetMAC, keyIPAddress, keyStatus,
		"last_failure_phase", "last_join_failure", "last_config_failure", "last_disc_failure",
		"disconnect_reason", "reboot_reason",
		"last_join_seconds", "last_config_seconds", "last_discovery_seconds", "last_error_seconds",
		keyController,
	}
}

// APTagKeys are the columns of show ap-tag, one row per access point. The MAC is named as "show ap
// tag summary" names it, where show ap keeps "radio_mac" beside the Ethernet MAC.
func APTagKeys() []string {
	return []string{
		keyAPName, keyAPMAC, "misconfigured", "misconfig_reason", "tag_source", "filter_name",
		keyPolicyTag, keySiteTag, keyRFTag, "ap_profile", keyFlexProfile, keyController,
	}
}

// PolicyTagKeys are the columns of show policy-tag, one row per WLAN binding the tag carries and
// one row for a tag that carries none.
func PolicyTagKeys() []string {
	return []string{
		DefaultSortPolicyTag, keyDescription, "wlan", keyPolicyProfile, keyController,
	}
}

// SiteTagKeys are the columns of show site-tag, one row per site tag. ap_join_profile is spelt as
// the leaf and the --ap-join-profile flag spell it, where show ap-tag calls it ap_profile.
func SiteTagKeys() []string {
	return []string{
		DefaultSortSiteTag, keyDescription, "ap_join_profile", keyFlexProfile,
		"local_site", keyController,
	}
}

// RFTagKeys are the columns of show rf-tag, one row per RF tag. The three profile keys
// are the --profile-24ghz, --profile-5ghz and --profile-6ghz flags of set rf-tag.
func RFTagKeys() []string {
	return []string{
		DefaultSortRFTag, keyDescription, "profile_24ghz", "profile_5ghz", "profile_6ghz",
		keyController,
	}
}

// ClientKeys are the columns of show client, one row per associated client.
func ClientKeys() []string {
	return []string{
		DefaultSortMAC, "ipv4", "ipv6", "device", "username", keySSID, keyAPName, keySlot,
		keyBand, "protocol", keyChannel, keyState, "rssi_dbm", "snr_db",
		"speed_mbps", "spatial_streams", "assoc_seconds", "rx_bytes", "tx_bytes", keyController,
	}
}

// WLANKeys are the columns of show wlan, one row per WLAN and bound policy profile. interface is
// the policy profile's interface name, which is not a VLAN id.
func WLANKeys() []string {
	return []string{
		DefaultSortWLANID, "profile", keySSID, keyStatus, "security", "bands", "broadcast",
		"p2p_block", "policy_status", "switching", "interface", "session_timeout_seconds",
		"dhcp_required", "policy_profile", "tags", keyController,
	}
}
