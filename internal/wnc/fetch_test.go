package wnc

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/umatare5/wnc/internal/config"
)

// Addresses in fixtures come from the documentation range RFC 7042 reserves, so no
// fixture carries a value read off a real device.
const (
	macAP1 = "00:00:5e:00:53:01"
	macAP2 = "00:00:5e:00:53:02"
	macCl1 = "00:00:5e:00:53:a1"
	macCl2 = "00:00:5e:00:53:a2"
)

func TestAPTags(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
	  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01","tag-info":{
	    "tag-source":"tag-source-static","is-ap-misconfigured":false,"ap-misconfig":"apmgr-no-misconfig",
	    "resolved-tag-info":{"resolved-policy-tag":"test-wlan-flex","resolved-site-tag":"test-site-flex","resolved-rf-tag":"test-inside"},
	    "site-tag":{"site-tag-name":"test-site-flex","ap-profile":"test-ap-profile01","flex-profile":"test-flex-profile01"},
	    "filter-info":{"filter-name":""}}},
	  {"wtp-mac":"` + macAP2 + `","name":"TEST-AP02","tag-info":{"tag-source":"tag-source-filter",
	    "filter-info":{"filter-name":"test-filter"}}}
	]}`}})

	tags, err := c.APTags(t.Context())
	if err != nil {
		t.Fatalf("APTags: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("got %d rows, want 2", len(tags))
	}

	if tags[0].Misconfigured == nil || *tags[0].Misconfigured {
		t.Errorf("an explicit false decoded as %v, want a pointer to false", tags[0].Misconfigured)
	}

	// The enum's own "no misconfiguration" member is a reading, so it has to survive as a
	// value rather than collapsing into the absence the next access point reports.
	if tags[0].MisconfigReason == nil || *tags[0].MisconfigReason != "apmgr-no-misconfig" {
		t.Errorf("ap-misconfig decoded as %v, want a pointer to apmgr-no-misconfig", tags[0].MisconfigReason)
	}

	if tags[1].FilterName != "test-filter" {
		t.Errorf("filter-name = %q, want test-filter", tags[1].FilterName)
	}

	// The second access point sent no resolved container at all. Every field below it is
	// a non-pointer struct in the SDK, so the whole set arrives empty rather than absent.
	if tags[1].Misconfigured != nil {
		t.Errorf("an omitted leaf decoded as %v, want nil", *tags[1].Misconfigured)
	}

	if tags[1].MisconfigReason != nil {
		t.Errorf("an omitted ap-misconfig decoded as %q, want nil", *tags[1].MisconfigReason)
	}

	if tags[1].PolicyTag != "" || tags[1].APProfile != "" {
		t.Errorf("an omitted container produced values: %+v", tags[1])
	}
}

// The fields expression names the whole tag-info container rather than the leaves inside
// it, and that is not a style choice: measured on 17.12.8, naming a node the release does
// not declare answers 200 with a chunked body that stops mid-object. The stub cannot
// reproduce that, so the request itself is asserted.
func TestAPTagsPrunesTheRequestToTheTagContainer(t *testing.T) {
	t.Parallel()

	var got string

	c := newClientWithQuery(t, &got,
		`{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[{"wtp-mac":"`+macAP1+`"}]}`)

	if _, err := c.APTags(t.Context()); err != nil {
		t.Fatalf("APTags: %v", err)
	}

	// Spelt out rather than built from apTagFields: comparing the constant against itself
	// would pass whatever the constant became. The semicolons must also reach the
	// controller unescaped, which is why this is the raw query — percent-encoded, the
	// controller reads the whole expression as one node name and answers with nothing.
	if want := "fields=wtp-mac;name;tag-info"; got != want {
		t.Errorf("query = %q, want %q", got, want)
	}
}

func TestAPTagsReportsAFailedRead(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"capwap-data": {status: http.StatusUnauthorized}})

	if _, err := c.APTags(t.Context()); err == nil {
		t.Fatal("a 401 did not produce an error")
	} else if cause, status := Classify(err); cause != CauseAuth || status != http.StatusUnauthorized {
		t.Errorf("Classify = %q/%d, want auth/401", cause, status)
	}
}

func TestAPs(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
		  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01","ip-addr":"192.168.0.11","country-code":"J4 ",
		   "num-radio-slots":2,
		   "device-detail":{"static-info":{"board-data":{"wtp-serial-num":"TST0000AP01","wtp-enet-mac":"` + macCl1 + `"},
		     "ap-models":{"model":"AIR-AP1815I-Q-K9"}},"wtp-version":{"sw-version":"17.12.7.13"}},
		   "ap-state":{"ap-admin-state":"adminstate-enabled","ap-operation-state":"registered"},
		   "ap-mode-data":{"wtp-mode":"mode-flex-connect","ap-sub-mode":"ap-sub-mode-none"},
		   "ap-time-info":{"boot-time":"2026-08-20T00:00:00.55648+00:00","join-time":"2026-08-20T01:00:00.801621+00:00"}},
		  {"wtp-mac":"` + macAP2 + `","name":"TEST-AP02"}
		]}`},
		"oper-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:oper-data":[
		  {"wtp-mac":"` + macAP1 + `","ap-pow":{"power-type":"pwr-src-poe-plus","power-mode":"dot11-default-high-pwr"}},
		  {"wtp-mac":"` + macAP2 + `"}
		]}`},
		"lldp-neigh": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:lldp-neigh":[
		  {"wtp-mac":"` + macAP1 + `","neigh-mac":"` + macCl1 + `","system-name":"test-sw-1","port-id":"Gi0/2"},
		  {"wtp-mac":"` + macAP1 + `","neigh-mac":"` + macCl2 + `","system-name":"test-sw-2","port-id":"Gi0/3"}
		]}`},
	})

	aps, reads, err := c.APs(t.Context())
	if err != nil {
		t.Fatalf("driving read failed: %v", err)
	}

	if reads.Power != nil || reads.LLDP != nil {
		t.Fatalf("secondary errors: %v / %v", reads.Power, reads.LLDP)
	}

	if len(aps) != 2 {
		t.Fatalf("got %d rows, want 2", len(aps))
	}

	first := aps[0]
	if first.Slots != 2 || first.Country != "J4 " || first.SWVersion != "17.12.7.13" {
		t.Errorf("first row = %+v", first)
	}

	// Two neighbors on one access point must stay one row, not duplicate it.
	if len(first.Neighbors) != 2 {
		t.Errorf("neighbors = %v, want two entries on one row", first.Neighbors)
	}

	if first.Neighbors[0] != "test-sw-1:Gi0/2" {
		t.Errorf("neighbor label = %q", first.Neighbors[0])
	}

	// The fraction is five digits on boot-time and six on join-time in this fixture,
	// which is the width variation the live controller produces.
	if first.BootTime.IsZero() || first.JoinTime.IsZero() {
		t.Errorf("timestamps did not parse: boot=%v join=%v", first.BootTime, first.JoinTime)
	}

	if !first.JoinTime.After(first.BootTime) {
		t.Errorf("join %v is not after boot %v", first.JoinTime, first.BootTime)
	}

	// An access point absent from the LLDP list keeps its row and reports no neighbor.
	if len(aps[1].Neighbors) != 0 {
		t.Errorf("second row invented a neighbor: %v", aps[1].Neighbors)
	}

	// ap-pow is one of the two absence-preserving pointers on this path.
	if aps[1].PowerType != "" {
		t.Errorf("an omitted power container produced %q", aps[1].PowerType)
	}
}

// The prune is asserted on the request because a stub cannot reproduce what the controller does
// with an undeclared node, which answers 200 with a body that stops mid-object. It also keeps
// proxy-info's username and password leaves on the controller.
func TestAPsPrunesTheRequestToTheRenderedNodes(t *testing.T) {
	t.Parallel()

	var got string

	c := newClient(t, routes{
		"capwap-data": {query: &got, body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
		  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01","num-radio-slots":2,
		   "device-detail":{"static-info":{"board-data":{"wtp-serial-num":"TST0000AP01"}}},
		   "ap-state":{"ap-admin-state":"adminstate-enabled"}}
		]}`},
		"oper-data":  {},
		"lldp-neigh": {},
	})

	aps, _, err := c.APs(t.Context())
	if err != nil {
		t.Fatalf("APs: %v", err)
	}

	want := "fields=name;wtp-mac;ip-addr;num-radio-slots;country-code;" +
		"device-detail;ap-mode-data;ap-state;ap-time-info"
	if got != want {
		t.Errorf("query = %q, want %q", got, want)
	}

	// A prune that drops a rendered column saves bytes by losing data. The serial sits
	// two containers deep under device-detail, so it is the one the expression could
	// silently cut.
	if len(aps) != 1 {
		t.Fatalf("got %d rows, want 1", len(aps))
	}

	if aps[0].Serial != "TST0000AP01" {
		t.Errorf("serial = %q, want TST0000AP01", aps[0].Serial)
	}

	if aps[0].Slots != 2 {
		t.Errorf("slots = %d, want 2", aps[0].Slots)
	}
}

func TestAPsReturnsTheDrivingReadFailureAsTheError(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"capwap-data": {status: http.StatusInternalServerError},
		"oper-data":   {body: `{"Cisco-IOS-XE-wireless-access-point-oper:oper-data":[]}`},
		"lldp-neigh":  {body: `{"Cisco-IOS-XE-wireless-access-point-oper:lldp-neigh":[]}`},
	})

	aps, reads, err := c.APs(t.Context())
	if err == nil {
		t.Fatal("the driving read failed and APs returned no error")
	}

	if reads.Power != nil || reads.LLDP != nil {
		t.Errorf("the driving read's failure leaked into a secondary slot: %+v", reads)
	}

	if aps != nil {
		t.Errorf("got %d rows, want none", len(aps))
	}
}

// A failure in either secondary read costs its own cells and never the rows.
func TestAPsKeepsRowsWhenASecondaryFails(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"capwap-data": {
			body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[{"wtp-mac":"` + macAP1 + `","name":"TEST-AP01"}]}`,
		},
		"oper-data":  {status: http.StatusInternalServerError},
		"lldp-neigh": {status: http.StatusInternalServerError},
	})

	aps, reads, err := c.APs(t.Context())
	if err != nil {
		t.Fatalf("a secondary failure cost the rows: %v", err)
	}

	if reads.Power == nil || reads.LLDP == nil {
		t.Fatal("a failed secondary read was not reported")
	}

	if len(aps) != 1 {
		t.Fatalf("got %d rows, want the row to survive", len(aps))
	}

	if aps[0].PowerType != "" || len(aps[0].Neighbors) != 0 {
		t.Errorf("a failed read produced values: %+v", aps[0])
	}
}

func TestClients(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"common-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[
		  {"client-mac":"` + macCl1 + `","ap-name":"TEST-AP01","ms-ap-slot-id":0,"co-state":"client-status-run","username":""},
		  {"client-mac":"` + macCl2 + `","ap-name":"TEST-AP01","ms-ap-slot-id":1,"co-state":"client-status-run","username":"test-user"}
		]}`},
		"dot11-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:dot11-oper-data":[
		  {"ms-mac-address":"` + macCl1 + `","vap-ssid":"test-essid01","radio-type":"dot11-radio-type-bg",
		   "ewlc-ms-phy-type":"client-dot11n-24-ghz-prot","current-channel":6,"ms-assoc-time":"2026-08-24T10:00:00Z"},
		  {"ms-mac-address":"` + macCl2 + `","vap-ssid":"test-essid02","radio-type":"dot11-radio-type-a",
		   "ewlc-ms-phy-type":"client-dot11ax-5ghz-prot","current-channel":64,"ms-assoc-time":"1970-01-01T00:00:00Z"}
		]}`},
		"traffic-stats": {body: `{"Cisco-IOS-XE-wireless-client-oper:traffic-stats":[
		  {"ms-mac-address":"` + macCl1 + `","bytes-rx":"18446744073709551615","bytes-tx":"1024",
		   "most-recent-rssi":-45,"most-recent-snr":0,"speed":72,"spatial-stream":0},
		  {"ms-mac-address":"` + macCl2 + `","bytes-rx":"not-a-number","bytes-tx":"","most-recent-rssi":0,"most-recent-snr":40}
		]}`},
		"sisf-db-mac": {body: `{"Cisco-IOS-XE-wireless-client-oper:sisf-db-mac":[
		  {"mac-addr":"` + macCl1 + `","ipv4-binding":{"ip-key":{"zone-id":0,"ip-addr":"192.168.0.21"}},
		   "ipv6-binding":[{"ip-key":{"zone-id":0,"ip-addr":"fe80::1"}},
		                   {"ip-key":{"zone-id":0,"ip-addr":"2001:db8::20"}},
		                   {"ip-key":{"zone-id":0,"ip-addr":"2001:db8::3"}}]},
		  {"mac-addr":"` + macCl2 + `","ipv4-binding":{"ip-key":{"zone-id":0,"ip-addr":"0.0.0.0"}},
		   "ipv6-binding":[{"ip-key":{"zone-id":0,"ip-addr":"fe80::2"}}]}
		]}`},
		"dc-info": {body: `{"Cisco-IOS-XE-wireless-client-oper:dc-info":[
		  {"client-mac":"` + macCl1 + `","device-name":"Unknown Device"}
		]}`},
	})

	clients, reads, err := c.Clients(t.Context())
	if err != nil {
		t.Fatalf("Clients: %v", err)
	}

	if reads.Dot11 != nil || reads.Stats != nil || reads.SISF != nil || reads.DC != nil {
		t.Fatalf("secondary errors: %+v", reads)
	}

	if len(clients) != 2 {
		t.Fatalf("got %d rows, want 2", len(clients))
	}

	first, second := clients[0], clients[1]

	// The lowest-comparing global address wins, and the two link-local entries are
	// dropped. "2001:db8::3" sorts below "2001:db8::20" numerically and above it as text,
	// which is what makes the comparison method visible here.
	if first.IPv6 != "2001:db8::3" {
		t.Errorf("IPv6 = %q, want the numerically lowest global address", first.IPv6)
	}

	if second.IPv6 != "" {
		t.Errorf("a link-local-only client reported %q", second.IPv6)
	}

	// A 64-bit counter arrives as a JSON string; an unparseable one stays absent rather
	// than becoming zero.
	if first.RxBytes == nil || *first.RxBytes != 18446744073709551615 {
		t.Errorf("RxBytes = %v, want the full 64-bit value", first.RxBytes)
	}

	if second.RxBytes != nil || second.TxBytes != nil {
		t.Errorf("an unparseable counter produced %v/%v", second.RxBytes, second.TxBytes)
	}

	// SNR is the one pointer the join sets: a genuine 0 dB must survive as a value,
	// and each row must carry its own address or a later row would overwrite an
	// earlier reading through the shared pointer.
	if first.SNR == nil || *first.SNR != 0 {
		t.Errorf("SNR = %v, want a pointer to a genuine 0 dB reading", first.SNR)
	}

	if second.SNR == nil || *second.SNR != 40 {
		t.Errorf("SNR = %v, want 40", second.SNR)
	}

	if first.SNR == second.SNR {
		t.Error("two rows share one SNR address")
	}

	if first.Band != "dot11-radio-type-bg" || first.PHY != "client-dot11n-24-ghz-prot" {
		t.Errorf("band and protocol come from different leaves: %+v", first)
	}

	if second.Device != "" {
		t.Errorf("a client absent from the classification list reported %q", second.Device)
	}
}

// A filter that reads from the 802.11 collection cannot be honestly applied when that
// read failed, and the caller is told so through the per-read error.
func TestClientsReportsTheFailedCollection(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"common-oper-data": {
			body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[{"client-mac":"` + macCl1 + `"}]}`,
		},
		"dot11-oper-data": {status: http.StatusInternalServerError},
		"traffic-stats":   {body: `{"Cisco-IOS-XE-wireless-client-oper:traffic-stats":[]}`},
		"sisf-db-mac":     {body: `{"Cisco-IOS-XE-wireless-client-oper:sisf-db-mac":[]}`},
		"dc-info":         {body: `{"Cisco-IOS-XE-wireless-client-oper:dc-info":[]}`},
	})

	clients, reads, err := c.Clients(t.Context())
	if err != nil {
		t.Fatalf("Clients: %v", err)
	}

	if reads.Dot11 == nil {
		t.Fatal("the failed 802.11 read was not reported")
	}

	if len(clients) != 1 || clients[0].HasDot11 {
		t.Errorf("rows = %+v, want one row marked as lacking 802.11 facts", clients)
	}

	// The traffic collection answered with an empty list, so it supplied no row for this
	// client and the SNR is an absence. It is the one leaf of the four this join carries
	// whose zero would read as a real 0 dB margin rather than as an omission.
	if clients[0].SNR != nil {
		t.Errorf("SNR = %d, want absent", *clients[0].SNR)
	}
}

func TestRadios(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"radio-oper-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"radio-type":"radio-80211abgn",
		   "radio-mode":"radio-mode-flex-connect","current-active-band":"dot11-2-dot-4-ghz-band","current-band-id":0,
		   "admin-state":"enabled","oper-state":"radio-up",
		   "phy-ht-cfg":{"cfg-data":{"curr-freq":1,"chan-width":20}},
		   "radio-band-info":[
		     {"band-id":0,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":19}}},
		     {"band-id":1,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":16}}}]},
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":2,"radio-type":"radio-80211-xor-5-6ghz",
		   "radio-mode":"radio-mode-flex-connect","current-active-band":"dot11-6-ghz-band","current-band-id":2,
		   "admin-state":"enabled","oper-state":"radio-up",
		   "phy-ht-cfg":{"cfg-data":{"curr-freq":5,"chan-width":40}},
		   "radio-band-info":[
		     {"band-id":1,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":22}}},
		     {"band-id":2,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":18}}}]},
		  {"wtp-mac":"` + macAP2 + `","radio-slot-id":2,"radio-type":"radio-remote-lan"}
		]}`},
		"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
		  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01","tag-info":{"resolved-tag-info":{"resolved-rf-tag":"test-inside"}}},
		  {"wtp-mac":"` + macAP2 + `","name":"TEST-AP02"}
		]}`},
		"common-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[
		  {"client-mac":"` + macCl1 + `","ap-name":"TEST-AP01","ms-ap-slot-id":0,"co-state":"client-status-run"},
		  {"client-mac":"` + macCl2 + `","ap-name":"TEST-AP01","ms-ap-slot-id":0,"co-state":"client-status-idle"}
		]}`},
		"rrm-measurement": {body: `{"Cisco-IOS-XE-wireless-rrm-oper:rrm-measurement":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"load":{"cca-util-percentage":28,"stations":1}}
		]}`},
		"rf-tags": {body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[
		  {"tag-name":"test-inside","dot11b-rf-profile-name":"test-rf-profile01","dot11a-rf-profile-name":"test-rf-profile02","dot11-6ghz-rf-prof-name":"test-rf-profile05"}
		]}}`},
	})

	radios, reads, err := c.Radios(t.Context())
	if err != nil {
		t.Fatalf("Radios: %v", err)
	}

	if reads.CAPWAP != nil || reads.Clients != nil || reads.RRM != nil || reads.RFTags != nil {
		t.Fatalf("secondary errors: %+v", reads)
	}

	// The remote-LAN entry is not a radio and is dropped.
	if len(radios) != 2 {
		t.Fatalf("got %d rows, want the remote-LAN entry filtered out", len(radios))
	}

	// The record is chosen by band-id == current-band-id. Taking the first would give
	// 22 on the second radio and taking the last would give 16 on the first.
	if radios[0].TxPowerDBm == nil || *radios[0].TxPowerDBm != 19 {
		t.Errorf("2.4 GHz TxPower = %v, want 19", radios[0].TxPowerDBm)
	}

	if radios[1].TxPowerDBm == nil || *radios[1].TxPowerDBm != 18 {
		t.Errorf("6 GHz TxPower = %v, want 18", radios[1].TxPowerDBm)
	}

	// The profile follows the band, not the slot: a slot-2 radio in 6 GHz takes the
	// 6 GHz profile.
	if radios[1].RFProfile != "test-rf-profile05" {
		t.Errorf("RF profile = %q, want the 6 GHz one", radios[1].RFProfile)
	}

	// Only the client in the run state counts.
	if radios[0].Clients == nil || *radios[0].Clients != 1 {
		t.Errorf("clients = %v, want 1", radios[0].Clients)
	}

	// A radio with no measurement row reports no utilization rather than zero.
	if radios[0].ChUtil == nil || *radios[0].ChUtil != 28 {
		t.Errorf("ChUtil = %v, want 28", radios[0].ChUtil)
	}

	if radios[1].ChUtil != nil {
		t.Errorf("a radio with no measurement row reported %v", *radios[1].ChUtil)
	}

	// A radio the client list could reach but with no client is a genuine zero.
	if radios[1].Clients == nil || *radios[1].Clients != 0 {
		t.Errorf("clients on an idle radio = %v, want a reported zero", radios[1].Clients)
	}
}

// current-band-id is 0 for 2.4 GHz, so an absent one must not select the band-0 record.
// The SDK declares the leaf as a pointer for this reason: before v0.10.0 it decoded as
// zero and this radio would have reported the 2.4 GHz power it is not serving.
func TestRadiosRefuseTxPowerWithoutTheCurrentBandID(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"radio-oper-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":1,"radio-type":"radio-80211a",
		   "radio-mode":"radio-mode-flex-connect","current-active-band":"dot11-5-ghz-band",
		   "admin-state":"enabled","oper-state":"radio-up",
		   "phy-ht-cfg":{"cfg-data":{"curr-freq":64,"chan-width":40}},
		   "radio-band-info":[
		     {"band-id":0,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":20}}},
		     {"band-id":1,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":17}}}]}
		]}`},
		"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
		  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01"}
		]}`},
		"common-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[]}`},
		"rrm-measurement":  {body: `{"Cisco-IOS-XE-wireless-rrm-oper:rrm-measurement":[]}`},
		"rf-tags":          {body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[]}}`},
	})

	radios, _, err := c.Radios(t.Context())
	if err != nil {
		t.Fatalf("Radios: %v", err)
	}

	if len(radios) != 1 {
		t.Fatalf("got %d rows, want 1", len(radios))
	}

	if radios[0].TxPowerDBm != nil {
		t.Errorf("TxPower = %d, want it absent: band 0 is 2.4 GHz and this radio serves 5",
			*radios[0].TxPowerDBm)
	}
}

// Without the access point list there is no name to resolve a client's radio through,
// so the tally is absent rather than zero.
func TestRadiosWithoutTheAccessPointList(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"radio-oper-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"radio-type":"radio-80211bg"}]}`},
		"capwap-data":      {status: http.StatusInternalServerError},
		"common-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[]}`},
		"rrm-measurement":  {body: `{"Cisco-IOS-XE-wireless-rrm-oper:rrm-measurement":[]}`},
		"rf-tags":          {body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[]}}`},
	})

	radios, reads, err := c.Radios(t.Context())
	if err != nil {
		t.Fatalf("Radios: %v", err)
	}

	if reads.CAPWAP == nil {
		t.Fatal("the failed access point read was not reported")
	}

	if len(radios) != 1 {
		t.Fatalf("got %d rows, want the row to survive", len(radios))
	}

	if radios[0].Clients != nil {
		t.Errorf("clients = %v, want nil when the name could not be resolved", *radios[0].Clients)
	}

	if radios[0].APName != "" {
		t.Errorf("AP name = %q, want empty", radios[0].APName)
	}
}

func TestWLANs(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"wlan-cfg-entries": {body: `{"Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-entries":{"wlan-cfg-entry":[
		  {"wlan-id":5,"profile-name":"test-wlan-profile01","psk":"must-never-be-decoded",
		   "apf-vap-id-data":{"ssid":"test-essid01","wlan-status":true,"broadcast-ssid":true,"p2p-block-action":"p2p-blocking-action-none"},
		   "security-wpa":true,"wpa2-enabled":true,"auth-key-mgmt-psk":true,"auth-key-mgmt-dot1x":false,
		   "wlan-radio-policies":{"wlan-radio-policy":[{"band":"dot11-2-dot-4-ghz-band"}]}},
		  {"wlan-id":9,"profile-name":"test-unbound",
		   "apf-vap-id-data":{"ssid":"test-unbound","wlan-status":false},
		   "security-wpa":false,
		   "wlan-radio-policies":{"wlan-radio-policy":[{"band":"dot11-5-ghz-band"},{"band":"dot11-6-ghz-band"}]}}
		]}}`},
		"wlan-policies": {body: `{"Cisco-IOS-XE-wireless-wlan-cfg:wlan-policies":{"wlan-policy":[
		  {"policy-profile-name":"test-policy-profile01","status":true,"interface-name":"TEST-INTERNAL",
		   "wlan-switching-policy":{"central-switching":false},"wlan-timeout":{"session-timeout":43200},
		   "dhcp-params":{"is-dhcp-enabled":true}}
		]}}`},
		"policy-list-entries": {body: `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entries":{"policy-list-entry":[
		  {"tag-name":"test-wlan-flex","wlan-policies":{"wlan-policy":[
		     {"wlan-profile-name":"test-wlan-profile01","policy-profile-name":"test-policy-profile01"},
		     {"wlan-profile-name":"test-missing","policy-profile-name":"test-policy-profile01"}]}},
		  {"tag-name":"default-policy-tag","description":"no bindings"}
		]}}`},
	})

	view, reads, err := c.WLANs(t.Context())
	if err != nil {
		t.Fatalf("WLANs: %v", err)
	}

	if reads.Profiles != nil || reads.Bindings != nil {
		t.Fatalf("secondary errors: %+v", reads)
	}

	if len(view.Entries) != 2 {
		t.Fatalf("got %d WLANs, want 2", len(view.Entries))
	}

	if got := view.Entries[0].Bands(); len(got) != 1 || got[0] != "dot11-2-dot-4-ghz-band" {
		t.Errorf("bands = %v", got)
	}

	if got := view.Entries[1].Bands(); len(got) != 2 {
		t.Errorf("a dual-band WLAN reported %v", got)
	}

	// The three bindings include one naming a WLAN that does not exist, and the tag
	// with no bindings sends no container at all.
	if len(view.Bindings) != 2 {
		t.Fatalf("bindings = %+v, want 2", view.Bindings)
	}

	profile, ok := view.Profiles["test-policy-profile01"]
	if !ok {
		t.Fatal("the policy profile was not indexed")
	}

	if profile.Shutdown || profile.InterfaceName != "TEST-INTERNAL" {
		t.Errorf("profile = %+v", profile)
	}

	if profile.CentralSwitching == nil || *profile.CentralSwitching {
		t.Errorf("central switching = %v, want a pointer to false", profile.CentralSwitching)
	}
}

// The raw struct declares no credential leaf, so report-all materializing one puts it
// nowhere: an undeclared member is dropped at decode. This asserts the fixture's psk
// value cannot be found anywhere in what was decoded.
func TestWLANsDropsUndeclaredCredentialLeaves(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"wlan-cfg-entries": {body: `{"Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-entries":{"wlan-cfg-entry":[
		  {"wlan-id":5,"profile-name":"p","psk":"must-never-be-decoded","wep-key":"must-never-be-decoded",
		   "psk-type":"ascii","apf-vap-id-data":{"ssid":"p"}}]}}`},
		"wlan-policies": {body: `{"Cisco-IOS-XE-wireless-wlan-cfg:wlan-policies":{"wlan-policy":[]}}`},
		"policy-list-entries": {
			body: `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entries":{"policy-list-entry":[]}}`,
		},
	})

	view, _, err := c.WLANs(t.Context())
	if err != nil {
		t.Fatalf("WLANs: %v", err)
	}

	if got := marshalOf(t, view.Entries); contains(got, "must-never-be-decoded") {
		t.Errorf("a credential leaf survived the decode: %s", got)
	}
}

func TestWLANsRejectsAWrongEnvelope(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"two top-level keys": `{"a":{},"b":{}}`,
		"wrong key":          `{"Cisco-IOS-XE-wireless-wlan-cfg:something-else":{}}`,
	}

	for name, body := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{"wlan-cfg-entries": {body: body}})

			if _, _, err := c.WLANs(t.Context()); err == nil {
				t.Fatal("a wrong envelope was accepted")
			}
		})
	}
}

// An empty body is how the controller answers a collection holding nothing. The SDK
// turns it into a success, so the view is empty rather than failing.
func TestEmptyBodyIsAnEmptyResult(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"wlan-cfg-entries":    {},
		"wlan-policies":       {},
		"policy-list-entries": {},
	})

	view, reads, err := c.WLANs(t.Context())
	if err != nil {
		t.Fatalf("WLANs: %v", err)
	}

	if reads.Profiles != nil || reads.Bindings != nil {
		t.Fatalf("an empty body was reported as a failure: %+v", reads)
	}

	if len(view.Entries) != 0 {
		t.Errorf("entries = %+v, want none", view.Entries)
	}
}

// The WLAN read carries with-defaults as a call-site option rather than through a wrapper of
// its own, so this is the only thing holding it in place. Five of the security leaves default
// to true, which makes a plain read report 802.1X off on a WLAN where it is on.
func TestWLANReadAsksForTheDefaultsInForce(t *testing.T) {
	t.Parallel()

	var got string

	c := newClient(t, routes{
		"wlan-cfg-entries": {
			query: &got,
			body:  `{"Cisco-IOS-XE-wireless-wlan-cfg:wlan-cfg-entries":{"wlan-cfg-entry":[]}}`,
		},
		"wlan-policies":       {},
		"policy-list-entries": {},
	})

	if _, _, err := c.WLANs(t.Context()); err != nil {
		t.Fatalf("WLANs: %v", err)
	}

	// Spelt out rather than built from the option: comparing the call against itself
	// would pass whatever the call became.
	if want := "with-defaults=report-all"; got != want {
		t.Errorf("query = %q, want %q", got, want)
	}
}

func TestParseInstant(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		zero bool
	}{
		{name: "six fraction digits", in: "2026-08-24T11:22:33.801621+00:00"},
		{name: "five fraction digits", in: "2026-08-24T11:22:33.55648+00:00"},
		{name: "no fraction", in: "2026-08-24T11:22:33+00:00"},
		{name: "zulu", in: "2026-08-24T11:22:33Z"},
		{name: "empty", in: "", zero: true},
		{name: "not a timestamp", in: "never", zero: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseInstant(tt.in)
			if got.IsZero() != tt.zero {
				t.Errorf("parseInstant(%q) zero = %v, want %v", tt.in, got.IsZero(), tt.zero)
			}
		})
	}
}

func TestPickIPv6(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   []string
		want string
	}{
		{name: "none", in: nil, want: ""},
		{name: "link-local only", in: []string{"fe80::1", "fe80::2"}, want: ""},
		{
			// Compared as text "2001:db8::20" sorts below "2001:db8::3"; compared as
			// addresses it does not. The controller compresses some entries and not
			// others, so only the address comparison is stable.
			name: "numeric order, not text order",
			in:   []string{"2001:db8::20", "2001:db8::3"},
			want: "2001:db8::3",
		},
		{name: "skips a malformed entry", in: []string{"nonsense", "2001:db8::5"}, want: "2001:db8::5"},
		{name: "ignores an IPv4 entry", in: []string{"192.168.0.1"}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pickIPv6(tt.in); got != tt.want {
				t.Errorf("pickIPv6(%v) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseCounter(t *testing.T) {
	t.Parallel()

	if got := parseCounter(""); got != nil {
		t.Errorf("an empty counter gave %v", *got)
	}

	if got := parseCounter("-1"); got != nil {
		t.Errorf("a negative counter gave %v", *got)
	}

	got := parseCounter("18446744073709551615")
	if got == nil || *got != 18446744073709551615 {
		t.Errorf("parseCounter = %v, want the full 64-bit value", got)
	}
}

func TestNeighborLabel(t *testing.T) {
	t.Parallel()

	tests := []struct{ system, port, want string }{
		{"", "", ""},
		{"test-sw-1", "", "test-sw-1"},
		{"", "Gi0/2", "Gi0/2"},
		{"test-sw-1", "Gi0/2", "test-sw-1:Gi0/2"},
	}

	for _, tt := range tests {
		t.Run(tt.system+"/"+tt.port, func(t *testing.T) {
			t.Parallel()

			if got := neighborLabel(tt.system, tt.port); got != tt.want {
				t.Errorf("neighborLabel(%q,%q) = %q, want %q", tt.system, tt.port, got, tt.want)
			}
		})
	}
}

// A read the caller canceled must be told apart from one that timed out, or an
// interrupted run is reported as an unreachable controller.
func TestClassifyTimeoutBeatsTheStatusCheck(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[]}`}})

	ctx, cancel := contextWithDeadlinePassed(t)
	defer cancel()

	_, err := c.APTags(ctx)
	if err == nil {
		t.Fatal("an expired deadline did not produce an error")
	}

	if cause, _ := Classify(err); cause != CauseTimeout {
		t.Errorf("Classify = %q, want timeout", cause)
	}
}

func TestClassifyStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		want   Cause
	}{
		{http.StatusUnauthorized, CauseAuth},
		{http.StatusForbidden, CauseForbidden},
		{http.StatusNotFound, CauseNotFound},
		{http.StatusBadRequest, CauseHTTP},
		{http.StatusInternalServerError, CauseHTTP},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{"capwap-data": {status: tt.status}})

			_, err := c.APTags(t.Context())
			if err == nil {
				t.Fatalf("status %d did not produce an error", tt.status)
			}

			cause, status := Classify(err)
			if cause != tt.want || status != tt.status {
				t.Errorf("Classify = %q/%d, want %q/%d", cause, status, tt.want, tt.status)
			}
		})
	}
}

// The SDK puts up to 512 bytes of the controller's error document into the message it
// builds, so the operator-facing line is rebuilt from the status alone.
func TestMessageDropsTheResponseBody(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"capwap-data": {status: http.StatusForbidden}})

	_, err := c.APTags(t.Context())
	if err == nil {
		t.Fatal("a 403 did not produce an error")
	}

	if got := Message(err); got != "the controller answered 403 Forbidden" {
		t.Errorf("Message = %q", got)
	}
}

// Every failure this CLI reports has to read alike, so no class may quote the error it came
// from. What each of these would otherwise carry is named in its own row: a timeout and a
// connection fault each bring a *url.Error naming the whole request URL, and an SDK call wraps
// its own sentence around both, which would make a typed write and an untyped one distinguishable
// by their prose alone.
func TestMessageQuotesNoErrorItDidNotWrite(t *testing.T) {
	t.Parallel()

	t.Run("a timeout", func(t *testing.T) {
		t.Parallel()

		c := newClient(t, routes{"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[]}`}})

		ctx, cancel := contextWithDeadlinePassed(t)
		defer cancel()

		_, err := c.APTags(ctx)
		if err == nil {
			t.Fatal("an expired deadline did not produce an error")
		}

		assertSynthesized(t, err, "the controller did not answer in time")
	})

	t.Run("an unreachable controller", func(t *testing.T) {
		t.Parallel()

		logger, err := logForTest(t)
		if err != nil {
			t.Fatalf("building the logger: %v", err)
		}

		// Port 1 on the loopback refuses rather than hanging, so this is the connection class
		// and not the timeout one.
		c, err := newClientFor(t, "127.0.0.1:1", fakeToken, logger)
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}

		_, err = c.APTags(t.Context())
		if err == nil {
			t.Fatal("a refused connection did not produce an error")
		}

		assertSynthesized(t, err, "the controller could not be reached")
	})

	t.Run("a certificate that does not verify", func(t *testing.T) {
		t.Parallel()

		logger, err := logForTest(t)
		if err != nil {
			t.Fatalf("building the logger: %v", err)
		}

		srv := httptest.NewTestServer(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		srv.StartTLS()

		// Verification left on against the harness's own certificate, which no root signs.
		c, err := NewClient(
			config.Target{Name: "test", Host: srv.Listener.Addr().String(), Token: fakeToken},
			config.Settings{Timeout: 10 * time.Second},
			logger, "wnc/test",
		)
		if err != nil {
			t.Fatalf("NewClient: %v", err)
		}

		_, err = c.APTags(t.Context())
		if err == nil {
			t.Fatal("an unsigned certificate did not produce an error")
		}

		if cause, _ := Classify(err); cause != CauseTLS {
			t.Fatalf("Classify = %q, want tls", cause)
		}

		assertSynthesized(t, err, "the controller's certificate did not verify")
	})
}

// assertSynthesized checks the reported line against the phrase and, as the control, that it is
// not the error's own text. Equality alone would pass a phrase that happened to be a prefix of
// what the chain says.
func assertSynthesized(t *testing.T, err error, want string) {
	t.Helper()

	got := Message(err)
	if got != want {
		t.Errorf("Message = %q, want %q", got, want)
	}

	if got == err.Error() {
		t.Errorf("Message quoted the error: %q", got)
	}
}

func TestClassifyNil(t *testing.T) {
	t.Parallel()

	if cause, status := Classify(nil); cause != "" || status != 0 {
		t.Errorf("Classify(nil) = %q/%d", cause, status)
	}
}

func TestNewClientRejectsAnEmptyHost(t *testing.T) {
	t.Parallel()

	logger, err := logForTest(t)
	if err != nil {
		t.Fatalf("logger: %v", err)
	}

	if _, err := newClientFor(t, "", fakeToken, logger); err == nil {
		t.Error("an empty host was accepted")
	}

	if _, err := newClientFor(t, "192.0.2.1", "", logger); err == nil {
		t.Error("an empty token was accepted")
	}
}

// A negative timeout is refused by every SDK option that takes one, so the failure has
// to surface as a client-construction error rather than as a hung read.
func TestNewClientRejectsANonPositiveTimeout(t *testing.T) {
	t.Parallel()

	logger, err := logForTest(t)
	if err != nil {
		t.Fatalf("logger: %v", err)
	}

	if _, err := newClientWithTimeout(t, logger, -time.Second); err == nil {
		t.Error("a negative timeout was accepted")
	}
}

// A Monitor or Sniffer radio is the case where the container pointer does not save the channel:
// the controller sends phy-ht-cfg and omits curr-freq inside it, because that leaf carries a "when"
// excluding those two modes and the invalid one. The width and the transmit power carry no such
// guard and are still reported, which is why only the channel may be suppressed.
func TestRadiosMonitorAndSnifferOmitOnlyTheChannel(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{
		"radio-oper-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"radio-type":"radio-80211bg",
		   "radio-mode":"radio-mode-sniffer","xor-radio-mode":"xor-radio-mode-sniffer",
		   "current-active-band":"dot11-2-dot-4-ghz-band","current-band-id":0,
		   "admin-state":"enabled","oper-state":"radio-up",
		   "phy-ht-cfg":{"cfg-data":{"chan-width":20}},
		   "radio-band-info":[{"band-id":0,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":20}}}]},
		  {"wtp-mac":"` + macAP2 + `","radio-slot-id":1,"radio-type":"radio-80211a",
		   "radio-mode":"radio-mode-monitor","xor-radio-mode":"xor-radio-mode-monitor",
		   "current-active-band":"dot11-5-ghz-band","current-band-id":1,
		   "admin-state":"enabled","oper-state":"radio-up",
		   "phy-ht-cfg":{"cfg-data":{"chan-width":40}},
		   "radio-band-info":[{"band-id":1,"phy-tx-pwr-lvl-cfg":{"cfg-data":{"curr-tx-power-in-dbm":17}}}]}
		]}`},
		"capwap-data": {body: `{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[
		  {"wtp-mac":"` + macAP1 + `","name":"TEST-AP01"},
		  {"wtp-mac":"` + macAP2 + `","name":"TEST-AP02"}
		]}`},
		"common-oper-data": {body: `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[]}`},
		"rrm-measurement": {body: `{"Cisco-IOS-XE-wireless-rrm-oper:rrm-measurement":[
		  {"wtp-mac":"` + macAP1 + `","radio-slot-id":0,"load":{"cca-util-percentage":0,"stations":0}},
		  {"wtp-mac":"` + macAP2 + `","radio-slot-id":1,"load":{"cca-util-percentage":0,"stations":0}}
		]}`},
		"rf-tags": {body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[]}}`},
	})

	radios, _, err := c.Radios(t.Context())
	if err != nil {
		t.Fatalf("Radios: %v", err)
	}

	if len(radios) != 2 {
		t.Fatalf("got %d rows, want 2", len(radios))
	}

	for _, r := range radios {
		// curr-freq carries a when guard on the radio mode and is genuinely absent here,
		// while chan-width beside it carries none and is sent. One guard cannot serve both.
		if r.Channel != nil {
			t.Errorf("%s: channel = %d, want the guarded leaf to arrive absent", r.Mode, *r.Channel)
		}

		if r.Width == nil {
			t.Errorf("%s: width is absent, but the leaf has no when guard and was sent", r.Mode)
		}

		if r.TxPowerDBm == nil {
			t.Errorf("%s: transmit power is absent, but the leaf has no when guard and was sent", r.Mode)
		}

		// A radio the client list could reach and that carries no client is a reported
		// zero, not an absence.
		if r.Clients == nil || *r.Clients != 0 {
			t.Errorf("%s: clients = %v, want a reported zero", r.Mode, r.Clients)
		}

		if r.ChUtil == nil || *r.ChUtil != 0 {
			t.Errorf("%s: utilization = %v, want a reported zero", r.Mode, r.ChUtil)
		}
	}

	if radios[0].Mode != "radio-mode-sniffer" || radios[1].Mode != "radio-mode-monitor" {
		t.Errorf("modes = %q / %q", radios[0].Mode, radios[1].Mode)
	}
}

// The join view's whole reason to exist is the access point capwap-data has dropped,
// so the fixture holds one joined and one not, and every never-happened instant is the
// epoch the controller actually sends.
func TestAPJoins(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"ap-join-stats": {body: `{"Cisco-IOS-XE-wireless-ap-global-oper:ap-join-stats":[
	  {"wtp-mac":"` + macAP1 + `","ap-disconnect-reason":"DTLS close alert from peer",
	   "reboot-reason":"ap-reboot-reason-none",
	   "ap-join-info":{"ap-name":"TEST-AP01","ap-ethernet-mac":"` + macCl1 + `","ap-ip-addr":"192.168.0.11",
	     "is-joined":true,"last-error-type":"ap-con-failure-run","last-join-failure-type":"jf-none",
	     "last-config-failure-type":"cf-none",
	     "last-succ-join-atmpt-time":"2026-08-24T10:00:00+00:00",
	     "last-succ-conf-atmpt-time":"2026-08-24T10:00:01+00:00",
	     "last-error-time":"1970-01-01T00:00:00+00:00"},
	   "ap-discovery-info":{"last-disc-failure-type":"disc-fail-none",
	     "last-success-disc-time":"2026-08-24T09:59:00+00:00"}},
	  {"wtp-mac":"` + macAP2 + `","ap-disconnect-reason":"Mode change to sniffer",
	   "ap-join-info":{"ap-name":"TEST-AP02","is-joined":false,
	     "last-error-type":"ap-con-failure-imgdwnld",
	     "last-succ-join-atmpt-time":"1970-01-01T00:00:00+00:00",
	     "last-succ-conf-atmpt-time":"1970-01-01T00:00:00+00:00"},
	   "ap-discovery-info":{}}
	]}`}})

	joins, err := c.APJoins(t.Context())
	if err != nil {
		t.Fatalf("APJoins: %v", err)
	}

	if len(joins) != 2 {
		t.Fatalf("got %d rows, want 2", len(joins))
	}

	if joins[0].Joined == nil || !*joins[0].Joined {
		t.Errorf("an explicit true decoded as %v", joins[0].Joined)
	}

	if joins[1].Joined == nil || *joins[1].Joined {
		t.Errorf("an explicit false decoded as %v", joins[1].Joined)
	}

	// An unjoined access point still carries its record, which is the whole point.
	if joins[1].Name != "TEST-AP02" || joins[1].DisconnectReason == "" {
		t.Errorf("the unjoined row lost its identity: %+v", joins[1])
	}

	// The epoch is the controller's "never", and it must not become an age.
	if joins[0].LastError.Unix() != 0 {
		t.Errorf("last-error-time = %v, want the epoch as sent", joins[0].LastError)
	}

	// An omitted container yields empty strings, not an absence: every container on
	// this path is a non-pointer struct in the SDK.
	if joins[1].LastDiscFailure != "" || joins[1].EthernetMAC != "" {
		t.Errorf("an omitted container produced values: %+v", joins[1])
	}
}

// The overview's access-point read is pruned to the three nodes it consumes, asserted on the
// request for the reason the tag view's own test gives, and the two reads share one constant.
// radioAPInfo is called directly because Radios makes five reads and the helper records the last.
func TestRadioAPInfoPrunesTheRequestToTheTagContainer(t *testing.T) {
	t.Parallel()

	var got string

	c := newClientWithQuery(t, &got,
		`{"Cisco-IOS-XE-wireless-access-point-oper:capwap-data":[{"wtp-mac":"`+macAP1+`","name":"TEST-AP01",
		  "tag-info":{"resolved-tag-info":{"resolved-rf-tag":"test-inside"}}}]}`)

	names, tags, err := c.radioAPInfo(t.Context())
	if err != nil {
		t.Fatalf("radioAPInfo: %v", err)
	}

	if want := "fields=wtp-mac;name;tag-info"; got != want {
		t.Errorf("query = %q, want %q", got, want)
	}

	// Both halves of the pruned body have to survive, or the prune saved bytes by
	// dropping a column: the name feeds the client tally's join and the RF tag feeds the
	// profile lookup.
	if names[macAP1] != "TEST-AP01" {
		t.Errorf("name = %q, want TEST-AP01", names[macAP1])
	}

	if tags[macAP1] != "test-inside" {
		t.Errorf("resolved RF tag = %q, want test-inside", tags[macAP1])
	}
}

// The three tag reads go through the SDK's own tag accessors, whose wrappers name the
// module-qualified key each keyed path returns. A wrapper naming the wrong key decodes to
// nothing and reports success, so each accessor is exercised against the envelope the
// controller sends rather than trusted.
func TestPolicyTags(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"policy-list-entries": {
		body: `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entries":{"policy-list-entry":[
		  {"tag-name":"test-wlan-flex","description":"test flex","wlan-policies":{"wlan-policy":[
		     {"wlan-profile-name":"test-wlan-profile01","policy-profile-name":"test-policy-profile01"},
		     {"wlan-profile-name":"test-wlan-profile02","policy-profile-name":"test-policy-profile01"}]}},
		  {"tag-name":"default-policy-tag"},
		  {"tag-name":"test-empty-desc","description":""}
		]}}`,
	}})

	tags, err := c.PolicyTags(t.Context())
	if err != nil {
		t.Fatalf("PolicyTags: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("got %d tags, want 3", len(tags))
	}

	if len(tags[0].Bindings) != 2 || tags[0].Bindings[1].WLANProfile != "test-wlan-profile02" {
		t.Errorf("bindings = %+v", tags[0].Bindings)
	}

	if tags[0].Description == nil || *tags[0].Description != "test flex" {
		t.Errorf("description = %v", tags[0].Description)
	}

	// A tag binding nothing sends no container at all, which is the built-in tag's
	// ordinary shape and must not become a binding to the empty string.
	if tags[1].Bindings != nil {
		t.Errorf("a tag with no container produced %+v", tags[1].Bindings)
	}

	if tags[1].Description != nil {
		t.Errorf("an omitted description decoded as %q", *tags[1].Description)
	}

	// A reported empty string is a pointer to "", which is what separates it from the
	// omission above and is why the row layer collapses it rather than the read.
	if tags[2].Description == nil || *tags[2].Description != "" {
		t.Errorf("a reported empty description decoded as %v", tags[2].Description)
	}
}

func TestSiteTags(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"site-tag-configs": {
		body: `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-configs":{"site-tag-config":[
		  {"site-tag-name":"test-site-flex","description":"test flex","ap-join-profile":"test-ap-profile01",
		   "flex-profile":"test-flex-profile01","is-local-site":false},
		  {"site-tag-name":"default-site-tag","is-local-site":true},
		  {"site-tag-name":"test-bare"}
		]}}`,
	}})

	tags, err := c.SiteTags(t.Context())
	if err != nil {
		t.Fatalf("SiteTags: %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("got %d tags, want 3", len(tags))
	}

	if tags[0].LocalSite == nil || *tags[0].LocalSite {
		t.Errorf("an explicit false decoded as %v, want a pointer to false", tags[0].LocalSite)
	}

	if tags[1].LocalSite == nil || !*tags[1].LocalSite {
		t.Errorf("an explicit true decoded as %v", tags[1].LocalSite)
	}

	// The leaf carries a schema default, so an omission is not "not local" and must stay
	// apart from the two readings above. This is what the read asks report-all for.
	if tags[2].LocalSite != nil {
		t.Errorf("an omitted flag decoded as %v, want nil", *tags[2].LocalSite)
	}

	if tags[2].APJoinProfile != nil || tags[2].FlexProfile != nil {
		t.Errorf("omitted profiles produced %+v", tags[2])
	}
}

func TestRFTags(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"rf-tags": {
		body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[
		  {"tag-name":"test-inside","description":"test inside","dot11b-rf-profile-name":"test-rf-profile01",
		   "dot11a-rf-profile-name":"test-rf-profile02","dot11-6ghz-rf-prof-name":"test-rf-profile05"},
		  {"tag-name":"default-rf-tag"}
		]}}`,
	}})

	tags, err := c.RFTags(t.Context())
	if err != nil {
		t.Fatalf("RFTags: %v", err)
	}

	if len(tags) != 2 {
		t.Fatalf("got %d tags, want 2", len(tags))
	}

	// The 2.4 GHz profile is the 802.11b leaf and the 5 GHz one the 802.11a leaf. Swapping
	// them is the one mistake this read can make and still look right.
	if deref(tags[0].Profile24GHz) != "test-rf-profile01" || deref(tags[0].Profile5GHz) != "test-rf-profile02" {
		t.Errorf("the band-to-leaf pairing is wrong: %+v", tags[0])
	}

	if deref(tags[0].Profile6GHz) != "test-rf-profile05" {
		t.Errorf("6 GHz profile = %q", deref(tags[0].Profile6GHz))
	}

	// The SDK declares each leaf as a pointer, so an omitted name arrives absent rather
	// than as an empty string the row layer cannot tell from a cleared one.
	if tags[1].Profile24GHz != nil || tags[1].Description != nil {
		t.Errorf("omitted names produced %+v", tags[1])
	}
}

// Every tag read asks for the values in force. These are configuration reads, and the
// parameter is what tells a leaf left at its default apart from one the controller
// withheld — measured as necessary on rf-tags, whose built-in tag omits all three per-band
// profile names on a plain read.
func TestTagReadsAskForTheDefaultsInForce(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		read func(*Client) error
	}{
		{
			name: "policy-tag",
			body: `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entries":{"policy-list-entry":[]}}`,
			read: func(c *Client) error { _, err := c.PolicyTags(t.Context()); return err },
		},
		{
			name: "site-tag",
			body: `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-configs":{"site-tag-config":[]}}`,
			read: func(c *Client) error { _, err := c.SiteTags(t.Context()); return err },
		},
		{
			name: "rf-tag",
			body: `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tags":{"rf-tag":[]}}`,
			read: func(c *Client) error { _, err := c.RFTags(t.Context()); return err },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got string

			c := newClientWithQuery(t, &got, tt.body)

			if err := tt.read(c); err != nil {
				t.Fatalf("read: %v", err)
			}

			// Spelt out rather than built from the option: comparing the call against
			// itself would pass whatever the call became.
			if want := "with-defaults=report-all"; got != want {
				t.Errorf("query = %q, want %q", got, want)
			}
		})
	}
}

// A failed tag read is fatal for its view: there is no secondary read to degrade, so the
// rows are dropped rather than a partial list being printed.
func TestTagReadsReportAFailedRead(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		route string
		read  func(*Client) error
	}{
		"policy-tag": {
			route: "policy-list-entries",
			read:  func(c *Client) error { _, err := c.PolicyTags(t.Context()); return err },
		},
		"site-tag": {
			route: "site-tag-configs",
			read:  func(c *Client) error { _, err := c.SiteTags(t.Context()); return err },
		},
		"rf-tag": {
			route: "rf-tags",
			read:  func(c *Client) error { _, err := c.RFTags(t.Context()); return err },
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := newClient(t, routes{tt.route: {status: http.StatusForbidden}})

			err := tt.read(c)
			if err == nil {
				t.Fatal("a 403 did not produce an error")
			}

			if cause, status := Classify(err); cause != CauseForbidden || status != http.StatusForbidden {
				t.Errorf("Classify = %q/%d, want forbidden/403", cause, status)
			}
		})
	}
}

// A node holding nothing answers with no body, which the SDK turns into a success. A
// controller with no tag of one kind is therefore an empty view rather than a failure.
func TestTagReadsTreatAnEmptyBodyAsAnEmptyView(t *testing.T) {
	t.Parallel()

	c := newClient(t, routes{"policy-list-entries": {}, "site-tag-configs": {}, "rf-tags": {}})

	policy, err := c.PolicyTags(t.Context())
	if err != nil || len(policy) != 0 {
		t.Errorf("PolicyTags = %+v, %v", policy, err)
	}

	site, err := c.SiteTags(t.Context())
	if err != nil || len(site) != 0 {
		t.Errorf("SiteTags = %+v, %v", site, err)
	}

	rf, err := c.RFTags(t.Context())
	if err != nil || len(rf) != 0 {
		t.Errorf("RFTags = %+v, %v", rf, err)
	}
}
