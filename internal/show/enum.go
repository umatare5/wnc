package show

// Every domain below was read from the module the controller serves at
// /restconf/tailf/modules/<module>/<revision>, never from documentation. A spelling this file
// does not know passes through unchanged, so a release that adds a member needs no code change.

// Display strings shared by more than one domain.
const (
	dispEnabled   = "Enabled"
	dispDisabled  = "Disabled"
	dispUnknown   = "Unknown"
	dispUp        = "Up"
	dispDown      = "Down"
	dispLocal     = "Local"
	dispMonitor   = "Monitor"
	dispFlex      = "FlexConnect"
	dispSniffer   = "Sniffer"
	dispBridge    = "Bridge"
	dispSensor    = "Sensor"
	dispRogueDet  = "Rogue-Detector"
	dispSEConnect = "SE-Connect"
	dispRemBridge = "Remote-Bridge"
	dispHybFlex   = "Hybrid-FlexConnect"
)

// Band display strings. They must match config.Band24, Band5 and Band6, the accepted values of
// --radio, so a filter and the cell it selects on read alike.
const (
	dispBand24 = "2.4"
	dispBand5  = "5"
	dispBand6  = "6"
)

// radioBand is enm-ewlc-dot11-radio-band, the AP radio's operating band.
// dot11-invalid-band is an explicit member, so it maps to a value of its own and
// must stay distinct from the leaf being absent.
var radioBand = map[string]string{
	"dot11-2-dot-4-ghz-band": dispBand24,
	"dot11-5-ghz-band":       dispBand5,
	"dot11-6-ghz-band":       dispBand6,
	"dot11-invalid-band":     dispUnknown,
}

// clientBand is ms-radio-type, carried by dot11-oper-data/radio-type. The similarly named
// common-oper-data/ms-radio-type leaf is typed with the PHY-generation typedef instead.
var clientBand = map[string]string{
	"dot11-radio-type-bg":   dispBand24,
	"dot11-radio-type-a":    dispBand5,
	"dot11-radio-type-6ghz": dispBand6,
	"dot11-radio-type-none": dispUnknown,
}

// radioAdmin and radioOper are enum-radio-admin-state and enum-radio-oper-state,
// both local to the access-point-oper module. Each holds exactly two members, so an
// absent leaf is the only third state and must not be folded into the negative one.
var (
	radioAdmin = map[string]string{"enabled": dispEnabled, "disabled": dispDisabled}
	radioOper  = map[string]string{"radio-up": dispUp, "radio-down": dispDown}
)

// apAdmin is the read side of the AP admin state, wireless-enum-types:admin-state.
// The RPC that sets it uses enm-admin-status, spelt "admin-state-enabled", with the
// same numbers — the two domains must never share a table.
var apAdmin = map[string]string{
	"adminstate-enabled":  dispEnabled,
	"adminstate-disabled": dispDisabled,
}

// dispRegistered is the one serving member of apState. show ap's glyph mapping names
// it so the cell and the glyph cannot drift apart.
const dispRegistered = "Registered"

// apState is enum-ap-state, local to the access-point-oper module.
var apState = map[string]string{
	"ap-down":         dispDown,
	"ap-up":           dispUp,
	"unregistered":    "Unregistered",
	"registered":      dispRegistered,
	"downloading":     "Downloading",
	"pre-downloading": "Pre-downloading",
}

// powerType is wireless-enum-types:power-type. The display strings name no 802.3 clause because
// the YANG names none, and "show ap config general" collapses the two PoE members to a bare
// "PoE", so keeping them apart shows which one was negotiated.
var powerType = map[string]string{
	"pwr-src-brick-old": "AC adapter",
	"pwr-src-brick-new": "AC adapter (new)",
	"pwr-src-inj":       "PoE injector",
	"pwr-src-poe-lgcy":  "PoE (legacy)",
	"pwr-src-poe-plus":  "PoE+",
	"pwr-src-unknown":   dispUnknown,
}

// powerMode is wireless-enum-types:power-mode-type. "(default)" marks the members
// whose name says the value was never set explicitly.
var powerMode = map[string]string{
	"dot11-default-low-pwr":  "Low Power (default)",
	"dot11-set-low-pwr":      "Low Power",
	"dot11-set-15-4-pwr":     "15.4W",
	"dot11-set-16-8-pwr":     "16.8W",
	"dot11-set-high-pwr":     "Full Power",
	"dot11-default-high-pwr": "Full Power (default)",
	"dot11-set-no-pwr":       "Off",
	"dot11-set-25-5-pwr":     "25.5W",
	"unknown-pwr":            dispUnknown,
}

// tagSource is wireless-ap-types:enm-ap-tag-source.
var tagSource = map[string]string{
	"tag-source-none":     "None",
	"tag-source-static":   "Static",
	"tag-source-filter":   "Filter",
	"tag-source-ap-pnp":   "AP-PnP",
	"tag-source-default":  "Default",
	"tag-source-location": "Location",
}

// apMisconfig is wireless-access-point-oper:ap-misconfig, a module-local typedef the 17.12
// releases do not declare at all. The device publishes no heading for it, so these strings come
// from the YANG descriptions.
var apMisconfig = map[string]string{
	"apmgr-no-misconfig": "None",
	"country-misconfig":  "Country",
	"world-wide-mode":    "World Wide Mode",
	"lic-no-comp":        "World Wide Mode (License)",
}

// radioMode is wireless-types:enm-ewlc-ap-radio-modes. The sibling xor-radio-mode leaf uses a
// different, four-member typedef with lookalike spellings.
var radioMode = map[string]string{
	"radio-mode-invalid":            dispUnknown,
	"radio-mode-local":              dispLocal,
	"radio-mode-monitor":            dispMonitor,
	"radio-mode-flex-connect":       dispFlex,
	"radio-mode-rogue-detector":     dispRogueDet,
	"radio-mode-sniffer":            dispSniffer,
	"radio-mode-bridge":             dispBridge,
	"radio-mode-se-connect":         dispSEConnect,
	"radio-mode-remote-bridge":      dispRemBridge,
	"radio-mode-hybrid-flexconnect": dispHybFlex,
	"radio-mode-sensor":             dispSensor,
	"radio-mode-wgb-uplink":         "WGB-Uplink",
	"radio-mode-uwgb":               "uWGB",
	"radio-mode-urwb":               "URWB",
	"radio-mode-wgb-scan":           "WGB-Scan",
}

// apMode is wireless-types:enm-ewlc-spam-ap-modes. Member zero is spelt local-mode
// while the rest are mode-*, so no prefix rule produces this table.
var apMode = map[string]string{
	"local-mode":              dispLocal,
	"mode-monitor":            dispMonitor,
	"mode-flex-connect":       dispFlex,
	"mode-rogue-detector":     dispRogueDet,
	"mode-sniffer":            dispSniffer,
	"mode-bridge":             dispBridge,
	"mode-se-connect":         dispSEConnect,
	"mode-remote-bridge":      dispRemBridge,
	"mode-hybrid-flexconnect": dispHybFlex,
	"mode-sensor":             dispSensor,
}

// apSubMode is wireless-ap-types:ap-sub-mode-type. Only the none member is
// suppressed in the composed Mode cell.
var apSubMode = map[string]string{
	"ap-sub-mode-none":    "",
	"wips-mode":           "WIPS",
	"non-local-network":   "Non-Local",
	"local-network":       "Local-Network",
	"forensic-awips-mode": "Forensic-aWIPS",
}

// clientPHY is wireless-client-types:ms-phy-radio-type. Four of its members denote no band at
// all, which is why this domain cannot supply the Band column.
var clientPHY = map[string]string{
	"client-unknown-prot":       dispUnknown,
	"client-dot11b":             "11b",
	"client-dot11g":             "11g",
	"client-dot11a":             "11a",
	"client-dot11n-24-ghz-prot": "11n",
	"client-dot11n-5-ghz-prot":  "11n",
	"client-dot11ac":            "11ac",
	"client-phy-type-notappl":   dispUnknown,
	"client-ethernet":           "Ethernet",
	"client-dot11ax-5ghz-prot":  "11ax",
	"client-dot11ax-24ghz-prot": "11ax",
	"client-802-3":              "Ethernet",
	"client-dot11ax-6ghz-prot":  "11ax",
	"client-dot11be-5ghz-prot":  "11be",
	"client-dot11be-24ghz-prot": "11be",
	"client-dot11be-6ghz-prot":  "11be",
}

// clientState is wireless-client-types:client-co-state.
var clientState = map[string]string{
	"client-status-idle":                       "Idle",
	"client-status-associating":                "Associating",
	"client-status-associated":                 "Associated",
	"client-status-authenticating":             "Authenticating",
	"client-status-authenticated":              "Authenticated",
	"client-status-mobility-discovery":         "Mobility Discovery",
	"client-status-mobility-complete":          "Mobility Complete",
	"client-status-ip-learning":                "IP Learning",
	"client-status-ip-learn-complete":          "IP Learned",
	"client-status-webauth-required":           "WebAuth Pending",
	"client-status-static-ip-anchor-discovery": "Anchor Discovery",
	"client-status-run":                        "Run",
	"client-status-delete-in-progress":         "Deleting",
	"client-status-deleted":                    "Deleted",
}

// p2pBlock is wireless-enum-types:apf-vap-p2p-blocking-action.
var p2pBlock = map[string]string{
	"p2p-blocking-action-none":                dispDisabled,
	"p2p-blocking-action-fwdup":               "Forward-Upstream",
	"p2p-blocking-action-drop":                "Drop",
	"p2p-blocking-action-allow-private-group": "Allow-Private-Group",
}

// ftMode is wireless-enum-types:ft-dot11r-mode.
var ftMode = map[string]string{
	"dot11r-disabled":         dispDisabled,
	"dot11r-enabled":          dispEnabled,
	"dot11r-adaptive-enabled": "Adaptive",
}

// joinFailurePhase is typedef last-failure-phase in Cisco-IOS-XE-wireless-ap-global-oper,
// carried by the last-error-type leaf. "Image-Download" expands the spelling rather than
// transforming it, so no prefix rule produces this table.
var joinFailurePhase = map[string]string{
	"ap-con-failure-unknown":   "Unknown",
	"ap-con-failure-discovery": "Discovery",
	"ap-con-failure-dtls":      "DTLS",
	"ap-con-failure-join":      "Join",
	"ap-con-failure-config":    "Config",
	"ap-con-failure-imgdwnld":  "Image-Download",
	"ap-con-failure-run":       "Run",
}

// joinNoFault holds the healthy member of the join, config, discovery and reboot failure
// domains. The device prints a display string for none of them, so only these four spellings are
// mapped and every other one passes through.
var joinNoFault = map[string]string{
	"jf-none":               "None",
	"cf-none":               "None",
	"disc-fail-none":        "None",
	"ap-reboot-reason-none": "None",
}

// lookup applies one domain, passing an unknown spelling through and leaving an
// empty input empty.
func lookup(table map[string]string, v string) string {
	if v == "" {
		return ""
	}

	if display, ok := table[v]; ok {
		return display
	}

	return v
}

func showRadioBand(v string) string   { return lookup(radioBand, v) }
func showClientBand(v string) string  { return lookup(clientBand, v) }
func showRadioAdmin(v string) string  { return lookup(radioAdmin, v) }
func showRadioOper(v string) string   { return lookup(radioOper, v) }
func showAPAdmin(v string) string     { return lookup(apAdmin, v) }
func showAPState(v string) string     { return lookup(apState, v) }
func showPowerType(v string) string   { return lookup(powerType, v) }
func showPowerMode(v string) string   { return lookup(powerMode, v) }
func showTagSource(v string) string   { return lookup(tagSource, v) }
func showAPMisconfig(v string) string { return lookup(apMisconfig, v) }
func showJoinPhase(v string) string   { return lookup(joinFailurePhase, v) }
func showJoinFault(v string) string   { return lookup(joinNoFault, v) }
func showRadioMode(v string) string   { return lookup(radioMode, v) }
func showClientPHY(v string) string   { return lookup(clientPHY, v) }
func showClientState(v string) string { return lookup(clientState, v) }
func showP2PBlock(v string) string    { return lookup(p2pBlock, v) }
func showFTMode(v string) string      { return lookup(ftMode, v) }

func showAPModeWithSub(mode, sub string) string {
	m := lookup(apMode, mode)

	s := lookup(apSubMode, sub)
	if s == "" || m == "" {
		return m
	}

	return m + " (" + s + ")"
}
