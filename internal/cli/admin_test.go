package cli

import (
	"io"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"sync"
	"testing"
)

// The last element each request arrives on. The access point's own admin state and one
// radio's are separate RPCs, so the element is what says which of the two a run reached.
const (
	apAdminRPC    = "Cisco-IOS-XE-wireless-access-point-cfg-rpc:set-ap-admin-state"
	radioAdminRPC = "Cisco-IOS-XE-wireless-access-point-cfg-rpc:set-ap-slot-admin-state"
	radioRead     = "radio-oper-data="
)

// The mode spellings both RPCs take. They are the write family's, not the adminstate-* pair
// the read side reports, so a value read back cannot stand in for either.
const (
	modeEnabled  = "admin-state-enabled"
	modeDisabled = "admin-state-disabled"
)

// The rows the keyed radio read answers with. Every band, radio-type and state spelling is a
// declared member of the schema the controller serves; the three pairs the RPC's must
// statement forbids are composed here, because no controller reports them.
const (
	radio24InSlot0 = `{"radio-slot-id":0,"radio-type":"radio-80211bg",` +
		`"current-active-band":"dot11-2-dot-4-ghz-band","admin-state":"enabled"}`
	radio5InSlot1 = `{"radio-slot-id":1,"radio-type":"radio-80211a",` +
		`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`
	radio24InSlot1 = `{"radio-slot-id":1,"radio-type":"radio-80211bg",` +
		`"current-active-band":"dot11-2-dot-4-ghz-band","admin-state":"enabled"}`
	radio5InSlot0 = `{"radio-slot-id":0,"radio-type":"radio-80211a",` +
		`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`
	radioNoBand      = `{"radio-slot-id":1,"radio-type":"radio-80211a","admin-state":"enabled"}`
	remoteLANInSlot2 = `{"radio-slot-id":2,"radio-type":"radio-remote-lan"}`
	radioInvalidBand = `{"radio-slot-id":1,"radio-type":"radio-80211a",` +
		`"current-active-band":"dot11-invalid-band","admin-state":"enabled"}`
	radioXORInSlot2 = `{"radio-slot-id":2,"radio-type":"radio-80211-xor-5-6ghz",` +
		`"current-active-band":"dot11-6-ghz-band","admin-state":"enabled"}`
	radioXORInSlot1 = `{"radio-slot-id":1,"radio-type":"radio-80211-xor-5-6ghz",` +
		`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`
	radioABGNInSlot0 = `{"radio-slot-id":0,"radio-type":"radio-80211abgn",` +
		`"current-active-band":"dot11-2-dot-4-ghz-band","admin-state":"enabled"}`
	radio6GHzInSlot2 = `{"radio-slot-id":2,"radio-type":"radio-80211-6ghz",` +
		`"current-active-band":"dot11-6-ghz-band","admin-state":"enabled"}`
	radioUWBInSlot1 = `{"radio-slot-id":1,"radio-type":"radio-uwb",` +
		`"current-active-band":"dot11-5-ghz-band","admin-state":"enabled"}`
	radioNoType = `{"radio-slot-id":1,"current-active-band":"dot11-5-ghz-band",` +
		`"admin-state":"enabled"}`
)

// radioAnswer is what the keyed radio read returns. An empty row is a list the controller
// returned with no member and a status is the read failing outright; the refusals below have
// to keep those two apart from each other and from a row that is present.
type radioAnswer struct {
	status int
	row    string
}

// adminStub answers the name read, the keyed radio read and both RPCs, and keeps each
// request's body. The body is the assertion no count can make: an RPC that arrived once says
// nothing about which of the two modes it carried.
type adminStub struct {
	addr string

	mu     sync.Mutex
	hits   map[string]int
	bodies map[string]string
}

// newAdminStub serves one access point under testAPName, whose row carries docMAC, and answers
// every keyed radio read with the same answer, whichever slot it arrives for.
func newAdminStub(t *testing.T, radio radioAnswer) *adminStub {
	t.Helper()

	return newAdminStubWithAddress(t, radio, docMAC)
}

// newAdminStubWithAddress is newAdminStub with the address the map row carries. A row
// carrying none is not a measured shape, and the guard that refuses it is the reason this
// is a parameter: the keyed radio read is keyed on that address.
func newAdminStubWithAddress(t *testing.T, radio radioAnswer, mapMAC string) *adminStub {
	t.Helper()

	s := &adminStub{hits: map[string]int{}, bodies: map[string]string{}}

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		base := path.Base(req.URL.Path)

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Errorf("reading the request body for %s: %v", base, err)
		}

		s.mu.Lock()
		s.hits[base]++
		s.bodies[base] = string(body)
		s.mu.Unlock()

		switch {
		case base == "ap-name-mac-map="+testAPName:
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(`{"Cisco-IOS-XE-wireless-access-point-oper:ap-name-mac-map":[` +
				`{"wtp-name":"` + testAPName + `","wtp-mac":"` + mapMAC + `","eth-mac":"00:00:5e:00:53:11"}]}`))
		case strings.HasPrefix(base, radioRead):
			if radio.status != 0 {
				w.WriteHeader(radio.status)

				return
			}

			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(
				`{"Cisco-IOS-XE-wireless-access-point-oper:radio-oper-data":[` + radio.row + `]}`))
		case base == apAdminRPC, base == radioAdminRPC:
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request for %s", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	srv.StartTLS()

	s.addr = srv.Listener.Addr().String()

	return s
}

func (s *adminStub) hit(base string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.hits[base]
}

func (s *adminStub) body(base string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.bodies[base]
}

// requests counts everything that arrived, the name read included. A guard that runs before
// the client exists must leave this at zero.
func (s *adminStub) requests() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	for _, v := range s.hits {
		n += v
	}

	return n
}

// adminLeaves is every leaf both verbs carry, with the RPC each one must reach, what its
// report has to name, and the part of its prompt that says what the write covers.
var adminLeaves = []struct {
	leaf   string
	rpc    string
	extra  []string
	names  []string
	prompt string
}{
	{
		leaf: leafAP, rpc: apAdminRPC,
		names: []string{testAPName},
		// Measured on 17.15.6: after an AP-level disable both radios still report their own
		// admin-state as enabled, so which of the two states is being set is not something
		// the operator can infer from the radio rows afterwards.
		prompt: "not one radio's",
	},
	{
		leaf: "radio", rpc: radioAdminRPC, extra: []string{"--slot", "1"},
		names:  []string{testAPName, "5 GHz"},
		prompt: "slot 1 (5 GHz)",
	},
}

// adminModes pairs each verb with the mode it must put on the wire and the words its run must
// not produce.
var adminModes = []struct {
	verb      string
	otherVerb string
	mode      string
	otherMode string
}{
	{verb: "enable", otherVerb: "disable", mode: modeEnabled, otherMode: modeDisabled},
	{verb: "disable", otherVerb: "enable", mode: modeDisabled, otherMode: modeEnabled},
}

// Both trees come out of adminVerbs, and an inverted on field there leaves the exit code, the
// request count and the path unchanged, so the mode field of the body is the only place the
// difference appears. The report is checked beside it because a format string carrying the other
// verb's word would lie about a write that was correct.
func TestAdminSendsTheModeItsVerbNames(t *testing.T) {
	for _, leaf := range adminLeaves {
		for _, verb := range adminModes {
			t.Run(verb.verb+"/"+leaf.leaf, func(t *testing.T) {
				stub := newAdminStub(t, radioAnswer{row: radio5InSlot1})

				args := append([]string{
					verb.verb, leaf.leaf, "--ap-name", testAPName,
					"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes",
				}, leaf.extra...)

				got := runCLI(t, "", false, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if n := stub.hit(leaf.rpc); n != 1 {
					t.Fatalf("%s arrived %d times, want 1", leaf.rpc, n)
				}

				if body := stub.body(leaf.rpc); !strings.Contains(body, verb.mode) ||
					strings.Contains(body, verb.otherMode) {
					t.Errorf("%s sent %s, want the %s mode", verb.verb, body, verb.mode)
				}

				// The two RPCs cover different scopes on the controller, so reaching the
				// other leaf's would act on more or less than the operator named.
				for _, other := range adminLeaves {
					if other.rpc == leaf.rpc {
						continue
					}

					if n := stub.hit(other.rpc); n != 0 {
						t.Errorf("%s %s reached %s %d times", verb.verb, leaf.leaf, other.rpc, n)
					}
				}

				if !strings.Contains(got.stdout, verb.verb+" sent") {
					t.Errorf("stdout %q does not report %q", got.stdout, verb.verb+" sent")
				}

				// Neither RPC declares an output container, so the report may name the
				// instruction and never the resulting state.
				for _, wrong := range []string{verb.otherVerb, "enabled", "disabled"} {
					if strings.Contains(got.stdout, wrong) {
						t.Errorf("stdout %q carries %q", got.stdout, wrong)
					}
				}

				if !strings.Contains(got.stdout, testAPName) {
					t.Errorf("stdout %q does not name the access point", got.stdout)
				}

				// An access point is identified by name and never by address. The resolve
				// answered with one, and the radio write puts it on the wire, so nothing but
				// a format string keeps it off the stream.
				if strings.Contains(got.stdout, docMAC) {
					t.Errorf("stdout %q carries an address", got.stdout)
				}
			})
		}
	}
}

// The slot is the operator's and the band number is the radio type's, so a dual-band radio takes 3
// whichever band it is on and the served band's own number answers 400. Slot 0 is here for the flag
// as much as for the band: an explicit zero has to act.
func TestAdminRadioSendsTheSlotAndTheBandTheTypeTakes(t *testing.T) {
	tests := []struct {
		name string
		slot string
		row  string
		wire string
	}{
		{name: "2.4 GHz on slot 0", slot: "0", row: radio24InSlot0, wire: `"slot-id":0,"band":"1"`},
		{name: "5 GHz on slot 1", slot: "1", row: radio5InSlot1, wire: `"slot-id":1,"band":"2"`},
		{name: "a dual-band radio on slot 0", slot: "0", row: radioABGNInSlot0, wire: `"slot-id":0,"band":"3"`},
		{name: "an XOR radio on slot 2", slot: "2", row: radioXORInSlot2, wire: `"slot-id":2,"band":"3"`},
		{name: "a dedicated 6 GHz radio", slot: "2", row: radio6GHzInSlot2, wire: `"slot-id":2,"band":"4"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newAdminStub(t, radioAnswer{row: tt.row})

			got := runCLI(t, "", false, "disable", "radio", "--ap-name", testAPName, "--slot", tt.slot,
				"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			// The read is keyed by the slot, so the element it arrived on is the slot the
			// write that follows is about.
			if n := stub.hit(radioRead + docMAC + "," + tt.slot); n != 1 {
				t.Errorf("the keyed read for slot %s arrived %d times, want 1", tt.slot, n)
			}

			if body := stub.body(radioAdminRPC); !strings.Contains(body, tt.wire) {
				t.Errorf("%s sent %s, want %s", radioAdminRPC, body, tt.wire)
			}
		})
	}
}

// Every guard settled before a client exists, so none of these reaches the controller it names:
// exit 2 has to keep meaning nothing was sent, the name read included.
//
// The radio leaf alone is driven here, because runRadioAdmin calls the guards itself and in its own
// order where the ap leaf shares runAPAction, which reset_test.go already drives.
func TestAdminRadioRefusesBeforeContactingAnything(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		controllers int
		piped       bool
		mentions    string
	}{
		{
			name: "no name", args: []string{"--slot", "1"},
			controllers: 1, mentions: "requires --ap-name",
		},
		{
			name:        "two names",
			args:        []string{"--slot", "1", "--ap-name", testAPName, "--ap-name", "TEST-AP02"},
			controllers: 1, mentions: "one access point per invocation, --ap-name given 2 times",
		},
		{
			name: "an empty name", args: []string{"--slot", "1", "--ap-name", ""},
			controllers: 1, mentions: "must not be empty",
		},
		{name: "no slot", args: []string{"--ap-name", testAPName}, controllers: 1, mentions: "requires --slot"},
		{
			name: "a slot above the declared range", args: []string{"--ap-name", testAPName, "--slot", "4"},
			controllers: 1, mentions: "accepted values are 0 to 3",
		},
		{
			name: "a slot below the declared range", args: []string{"--ap-name", testAPName, "--slot", "-1"},
			controllers: 1, mentions: "accepted values are 0 to 3",
		},
		{
			name: "no controller", args: []string{"--ap-name", testAPName, "--slot", "1"},
			controllers: 0, mentions: "no controller given",
		},
		{
			name: "two controllers", args: []string{"--ap-name", testAPName, "--slot", "1"},
			controllers: 2, mentions: "one controller",
		},
		{
			name: "piped stdin cannot answer the prompt", args: []string{"--ap-name", testAPName, "--slot", "1"},
			controllers: 1, piped: true, mentions: "--yes",
		},
		// The slot is decidable without the controller, so it is decided first.
		{
			name: "no slot outranks two controllers", args: []string{"--ap-name", testAPName},
			controllers: 2, mentions: "requires --slot",
		},
		// And the name is decided before the slot, for the same reason.
		{
			name: "no name outranks no slot", args: []string{},
			controllers: 1, mentions: "requires --ap-name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newAdminStub(t, radioAnswer{row: radio5InSlot1})

			args := append([]string{"disable", "radio"}, tt.args...)
			for range tt.controllers {
				args = append(args, "-c", stub.addr)
			}

			args = append(args, "--access-token", fakeToken, "-k")
			if !tt.piped {
				args = append(args, "--yes")
			}

			got := runCLI(t, "", tt.piped, args...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}

			if got.stdout != "" {
				t.Errorf("a refused action wrote to stdout: %q", got.stdout)
			}

			if n := stub.requests(); n != 0 {
				t.Errorf("a refused action made %d requests", n)
			}
		})
	}
}

// Every state radioForAdmin refuses, each after the read and none with a write, so these exit 1
// rather than 2: the controller was reached and answered.
//
// The must-statement rows are here so a pair the controller would refuse is refused before a write:
// band 1 with slot 1 answers 400 naming the must, and band 2 is accepted on slot 1 or 2 alone. One
// verb runs the table, because radioForAdmin is reached before the verb is used.
func TestAdminRadioRefusesAStateTheRPCCannotName(t *testing.T) {
	tests := []struct {
		name     string
		slot     string
		radio    radioAnswer
		mentions string
	}{
		{
			name: "a list with no row", slot: "3", radio: radioAnswer{},
			mentions: "holds no radio in slot 3",
		},
		{
			name: "the read answered 404", slot: "3", radio: radioAnswer{status: http.StatusNotFound},
			mentions: "holds no radio in slot 3",
		},
		{
			name: "a remote-LAN port", slot: "2", radio: radioAnswer{row: remoteLANInSlot2},
			mentions: "remote-LAN port",
		},
		{
			name: "no band for the slot", slot: "1", radio: radioAnswer{row: radioNoBand},
			mentions: "reports no band",
		},
		{
			name: "a served band the prompt has no label for", slot: "1",
			radio: radioAnswer{row: radioInvalidBand}, mentions: "unknown band",
		},
		{
			name: "a radio type the RPC has no band number for", slot: "1",
			radio: radioAnswer{row: radioUWBInSlot1}, mentions: "has no band number for",
		},
		{
			name: "a record carrying no radio type", slot: "1",
			radio: radioAnswer{row: radioNoType}, mentions: "reports no radio type",
		},
		{
			name: "2.4 GHz on a slot the must statement forbids", slot: "1",
			radio: radioAnswer{row: radio24InSlot1}, mentions: "accepts that radio on slot 0 only",
		},
		{
			name: "5 GHz on a slot the must statement forbids", slot: "0",
			radio: radioAnswer{row: radio5InSlot0}, mentions: "accepts that radio on slot 1 or 2 only",
		},
		{
			name: "a dual-band radio on a slot the must statement forbids", slot: "1",
			radio: radioAnswer{row: radioXORInSlot1}, mentions: "accepts that radio on slot 0 or 2 only",
		},
		{
			name: "the read failed", slot: "1", radio: radioAnswer{status: http.StatusInternalServerError},
			mentions: "reading radio-oper-data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newAdminStub(t, tt.radio)

			got := runCLI(t, "", false, "disable", "radio", "--ap-name", testAPName, "--slot", tt.slot,
				"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

			if got.code != ExitFailure {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}

			if got.stdout != "" {
				t.Errorf("a refused action wrote to stdout: %q", got.stdout)
			}

			if n := stub.hit(radioRead + docMAC + "," + tt.slot); n != 1 {
				t.Errorf("the keyed read arrived %d times, want 1", n)
			}

			if n := stub.hit(radioAdminRPC); n != 0 {
				t.Errorf("%s was sent anyway, %d times", radioAdminRPC, n)
			}
		})
	}
}

// A dry run resolves the target and stops. Reading is not changing, so it does contact the
// controller; the assertion is that the RPC does not. The radio leaf's line names the band as
// well, which is knowable only after the read — the reason runRadioAdmin cannot reuse
// runAPAction.
func TestAdminDryRunReadsAndSendsNothing(t *testing.T) {
	for _, leaf := range adminLeaves {
		for _, verb := range adminModes {
			t.Run(verb.verb+"/"+leaf.leaf, func(t *testing.T) {
				stub := newAdminStub(t, radioAnswer{row: radio5InSlot1})

				args := append([]string{
					"--dry-run", verb.verb, leaf.leaf, "--ap-name", testAPName,
					"-c", stub.addr, "--access-token", fakeToken, "-k",
				}, leaf.extra...)

				got := runCLI(t, "", true, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if !strings.Contains(got.stdout, "would "+verb.verb) {
					t.Errorf("stdout %q does not report %q", got.stdout, "would "+verb.verb)
				}

				for _, want := range leaf.names {
					if !strings.Contains(got.stdout, want) {
						t.Errorf("stdout %q does not name %q", got.stdout, want)
					}
				}

				if n := stub.hit(leaf.rpc); n != 0 {
					t.Errorf("a dry run sent %s %d times", leaf.rpc, n)
				}

				// A dry run needs no answer, so a piped stdin must not refuse it.
				if strings.Contains(got.stderr, "--yes") {
					t.Errorf("a dry run demanded a confirmation: %q", got.stderr)
				}
			})
		}
	}
}

// The prompt, the cancellation and the flag past it. runRadioAdmin does not reuse runAPAction —
// the radio read has to land before the prompt so the band is checked rather than consented
// to — so this tree holds a second copy of the sequence and both leaves are driven through
// it. One verb runs it, because what differs per verb is the wording and not the guard.
func TestAdminConfirmation(t *testing.T) {
	answers := []struct {
		name     string
		stdin    string
		extra    []string
		prompted bool
		wantSent int
	}{
		{name: "answered no", stdin: "n\n", prompted: true, wantSent: 0},
		{name: "answered yes", stdin: "y\n", prompted: true, wantSent: 1},
		{name: "flag instead of a prompt", extra: []string{"--yes"}, wantSent: 1},
	}

	for _, leaf := range adminLeaves {
		for _, tt := range answers {
			t.Run(leaf.leaf+"/"+tt.name, func(t *testing.T) {
				stub := newAdminStub(t, radioAnswer{row: radio5InSlot1})

				args := append([]string{
					"disable", leaf.leaf, "--ap-name", testAPName,
					"-c", stub.addr, "--access-token", fakeToken, "-k",
				}, leaf.extra...)
				args = append(args, tt.extra...)

				got := runCLI(t, tt.stdin, false, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if n := stub.hit(leaf.rpc); n != tt.wantSent {
					t.Errorf("%s arrived %d times, want %d", leaf.rpc, n, tt.wantSent)
				}

				if tt.wantSent == 0 && !strings.Contains(got.stdout, "canceled") {
					t.Errorf("a declined action did not say so: %q", got.stdout)
				}

				// The prompt names what the write covers, which the two leaves differ on.
				if tt.prompted && !strings.Contains(got.stdout, leaf.prompt) {
					t.Errorf("the prompt %q does not carry %q", got.stdout, leaf.prompt)
				}
			})
		}
	}
}

// The one place this file pins the suffix every usage fault carries. It belongs on a leaf two
// levels down: FullName is what makes the suffix a command an operator can retype, and the
// suggestion is drawn from the flags this leaf parses, its own --slot included.
func TestAdminRadioUsageFaultNamesTheLeafsOwnHelp(t *testing.T) {
	got := runCLI(t, "", false, "enable", "radio", "--slotx", "1", "--ap-name", testAPName)

	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}

	for _, want := range []string{"did you mean --slot?", "see 'wnc enable radio --help'"} {
		if !strings.Contains(got.stderr, want) {
			t.Errorf("stderr %q does not carry %q", got.stderr, want)
		}
	}
}

// Each constructor must build the verb it is named for: adminVerb.name becomes the
// command's Name, so a swapped index renames the command an operator types.
func TestAdminCommandsAreNamedForTheVerbTheyBuild(t *testing.T) {
	if got := enableCommand().Name; got != "enable" {
		t.Errorf("enableCommand builds %q, want enable", got)
	}

	if got := disableCommand().Name; got != "disable" {
		t.Errorf("disableCommand builds %q, want disable", got)
	}
}

// radio-oper-data is keyed on wtp-mac, so a map row carrying no address leaves the read with
// no key. Asking anyway would key on the empty string, and the controller's answer would be
// reported as a slot the access point does not hold.
func TestAdminRadioRefusesAMapRowWithNoAddress(t *testing.T) {
	stub := newAdminStubWithAddress(t, radioAnswer{row: radio5InSlot1}, "")

	got := runCLI(t, "", false, "disable", "radio", "--ap-name", testAPName, "--slot", "1",
		"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
	}

	if !strings.Contains(got.stderr, "reports no radio address") {
		t.Errorf("stderr %q does not name the missing address", got.stderr)
	}

	if n := stub.hit(radioRead); n != 0 {
		t.Errorf("the radio read went out on an empty key %d times", n)
	}

	if n := stub.hit(radioAdminRPC); n != 0 {
		t.Errorf("%s was sent anyway", radioAdminRPC)
	}
}

// The one run that observes both halves of the split at once: the prompt names the band the radio is
// serving, and the wire carries the number its type takes. Merging the two tables again would make
// this line and this body agree, which is exactly the defect it exists to catch.
func TestAdminRadioNamesTheServedBandAndSendsTheTypesBand(t *testing.T) {
	stub := newAdminStub(t, radioAnswer{row: radioXORInSlot2})

	got := runCLI(t, "", false, "disable", "radio", "--ap-name", testAPName, "--slot", "2",
		"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
	}

	if !strings.Contains(got.stdout, "slot 2 (6 GHz)") {
		t.Errorf("stdout %q does not name the band the radio is serving", got.stdout)
	}

	body := stub.body(radioAdminRPC)

	if !strings.Contains(body, `"band":"3"`) {
		t.Errorf("%s sent %s, want the dual-band number", radioAdminRPC, body)
	}

	for _, wrong := range []string{`"band":"4"`, `"band":"6"`} {
		if strings.Contains(body, wrong) {
			t.Errorf("%s sent %s, which carries %s", radioAdminRPC, body, wrong)
		}
	}

	// An empty label would render as "( GHz)" and a bad verb as a %! directive: both would
	// leave the operator confirming a radio the line does not name.
	for _, broken := range []string{"( GHz)", "%!"} {
		if strings.Contains(got.stdout, broken) {
			t.Errorf("stdout %q carries %q", got.stdout, broken)
		}
	}
}
