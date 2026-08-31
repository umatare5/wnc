package cli

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"sync"
	"testing"
)

// Documentation addresses from the range IANA reserves. docMAC is what the resolve answers with
// and nothing an operator sees carries it, while testClientMAC is separate because deauth does
// print a client's address; testClientRowMAC differs from it so a leaf echoing the operator's
// argument instead of the controller's spelling cannot pass.
const (
	docMAC           = "00:00:5e:00:53:01"
	testAPName       = "TEST-AP01"
	testClientMAC    = "00:00:5e:00:53:a1"
	testClientRowMAC = "00:00:5E:00:53:A1"
)

// The usernames the collection read answers with. One is carried by a single session and one by
// two, because the count is what the username arm's prompt reports and a fixture holding one
// session could not tell the singular wording from the plural one. A username no session carries
// is not a constant: any other value is one.
const (
	testClientUsername = "test-user-1"
	testSharedUsername = "test-user-2"
)

// clientCollection is the unkeyed read the username resolve makes. The fourth row carries the
// empty username most clients do, so a resolve that matched loosely would count it.
const clientCollection = `{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[` +
	`{"client-mac":"` + testClientRowMAC + `","username":"` + testClientUsername + `"},` +
	`{"client-mac":"00:00:5E:00:53:A2","username":"` + testSharedUsername + `"},` +
	`{"client-mac":"00:00:5E:00:53:A3","username":"` + testSharedUsername + `"},` +
	`{"client-mac":"00:00:5E:00:53:A4","username":""}]}`

// The last path element each RPC arrives on.
const (
	apResetRPC     = "Cisco-IOS-XE-wireless-access-point-cmd-rpc:ap-reset"
	capwapResetRPC = "Cisco-IOS-XE-wireless-access-point-cmd-rpc:set-rad-capwap-reset"
	saveConfigRPC  = "cisco-ia:save-config"
	deauthRPC      = "Cisco-IOS-XE-wireless-client-rpc:apf-ms-delete-all"
)

// saveReply is the shape a controller answered the save with. save-config is the one RPC this
// stub answers with a body, because it is the one that declares an output container.
const saveReply = `{"cisco-ia:output":{"result":"Save running-config successful"}}`

// controllerStub answers the reads and RPCs of reset, save-config and deauth, and counts what
// arrived. The count is the assertion that matters: a canceled or dry run must leave the RPC
// untouched, and nothing else can prove that.
type controllerStub struct {
	addr string

	mu   sync.Mutex
	hits map[string]int
}

// newControllerStub serves one access point under testAPName. A nameStatus other than 200
// stands for a controller holding no access point under that name.
func newControllerStub(t *testing.T, nameStatus int) *controllerStub {
	t.Helper()

	s := &controllerStub{hits: map[string]int{}}

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		base := path.Base(req.URL.Path)

		s.mu.Lock()
		s.hits[base]++
		s.mu.Unlock()

		switch {
		case base == "ap-name-mac-map="+testAPName:
			if nameStatus != http.StatusOK {
				w.WriteHeader(nameStatus)

				return
			}

			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(`{"Cisco-IOS-XE-wireless-access-point-oper:ap-name-mac-map":[` +
				`{"wtp-name":"` + testAPName + `","wtp-mac":"` + docMAC + `","eth-mac":"00:00:5e:00:53:11"}]}`))
		case base == apResetRPC, base == capwapResetRPC:
			w.WriteHeader(http.StatusNoContent)
		case base == saveConfigRPC:
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(saveReply))
		case base == "common-oper-data="+testClientMAC:
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(`{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[` +
				`{"client-mac":"` + testClientRowMAC + `","ap-name":"` + testAPName +
				`","co-state":"client-status-run"}]}`))
		case base == deauthRPC:
			w.WriteHeader(http.StatusNoContent)
		case base == "common-oper-data":
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(clientCollection))
		case strings.HasPrefix(base, "common-oper-data="):
			// Any other address is one the controller holds no client at, which is what
			// TestDeauthRefusesAnAddressTheControllerDoesNotHold reads.
			w.WriteHeader(http.StatusNotFound)
		case strings.HasPrefix(base, "ap-name-mac-map="):
			// Any other name is one the controller does not hold, which is what
			// TestAPNameSwallowsTheFlagBehindIt reads.
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Errorf("unexpected request for %s", req.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	srv.StartTLS()

	s.addr = srv.Listener.Addr().String()

	return s
}

func (s *controllerStub) hit(base string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.hits[base]
}

// Every guard below is settled before a client exists, so none of these reaches a
// controller. That is the property being asserted as much as the exit code: an action
// that is going to be refused must be refused while exit 2 still means nothing was sent.
func TestResetRefusesBeforeContactingAnything(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		piped    bool
		mentions string
	}{
		{
			name:     "no name",
			args:     []string{},
			mentions: "requires --ap-name",
		},
		{
			name:     "two names",
			args:     []string{"--ap-name", testAPName, "--ap-name", "TEST-AP02"},
			mentions: "one access point per invocation, --ap-name given 2 times",
		},
		{
			// The SDK answers an empty key with its own not-found, which would reach the
			// operator as a read failure at exit 1 rather than as what they typed.
			name:     "an empty name",
			args:     []string{"--ap-name", ""},
			mentions: "must not be empty",
		},
		{
			// TrimSpace and not a bare comparison: a name of spaces is as unusable a key as
			// an empty one, and the SDK would answer both with its own not-found.
			name:     "a name of spaces",
			args:     []string{"--ap-name", "   "},
			mentions: "must not be empty",
		},
		{
			name:     "no controller",
			args:     []string{"--ap-name", testAPName},
			mentions: "no controller given",
		},
		{
			name:     "two controllers",
			args:     []string{"--ap-name", testAPName, "-c", "h1", "-c", "h2", "--access-token", fakeToken},
			mentions: "one controller",
		},
		{
			name:     "piped stdin cannot answer the prompt",
			args:     []string{"--ap-name", testAPName, "-c", "h1", "--access-token", fakeToken},
			piped:    true,
			mentions: "--yes",
		},
	}

	for _, leaf := range resetLeaves {
		for _, tt := range tests {
			t.Run(leaf.leaf+"/"+tt.name, func(t *testing.T) {
				got := runCLI(t, "", tt.piped, append([]string{"reset", leaf.leaf}, tt.args...)...)

				if got.code != ExitUsage {
					t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
				}

				if !strings.Contains(got.stderr, tt.mentions) {
					t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
				}

				if got.stdout != "" {
					t.Errorf("a refused action wrote to stdout: %q", got.stdout)
				}

				// A read would have failed against these hosts and logged a cause; none is
				// expected, because nothing was sent.
				if strings.Contains(got.stderr, "cause=") {
					t.Errorf("a refused action contacted a controller: %q", got.stderr)
				}
			})
		}
	}
}

// The name reaches the controller stage as typed. Nothing local judges its characters or its
// length, for the reason requireAPName states.
func TestResetAPAcceptsTheNameAsTyped(t *testing.T) {
	for _, name := range []string{"TEST-AP01", "test-ap.floor-3", "TEST-AP01-北"} {
		t.Run(name, func(t *testing.T) {
			got := runCLI(t, "", false, "reset", "ap", "--ap-name", name)

			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
			}

			// It got past the positional and failed on the missing controller, which is the
			// next guard in order.
			if !strings.Contains(got.stderr, "no controller given") {
				t.Errorf("stderr %q, want the controller fault rather than a name fault", got.stderr)
			}
		})
	}
}

// The group with no leaf prints help rather than acting, like every other parent here.
func TestResetGroupShowsHelp(t *testing.T) {
	got := runCLI(t, "", false, "reset")

	if got.code != ExitOK {
		t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
	}

	for _, leaf := range resetLeaves {
		if !strings.Contains(got.stdout, leaf.leaf) {
			t.Errorf("help does not list the %s leaf: %q", leaf.leaf, got.stdout)
		}
	}
}

func TestResetAPUnknownLeafIsAUsageFault(t *testing.T) {
	got := runCLI(t, "", false, "reset", "bogus")

	if got.code != ExitUsage {
		t.Errorf("exit = %d, want %d", got.code, ExitUsage)
	}
}

// The two flags that decide whether this command acts have to be readable from the
// leaf, because urfave lists a parent's flags in the parent's help alone.
func TestResetHelpNamesTheFlagsThatDecideWhetherItActs(t *testing.T) {
	for _, leaf := range resetLeaves {
		t.Run(leaf.leaf, func(t *testing.T) {
			got := runCLI(t, "", false, "reset", leaf.leaf, "--help")

			for _, want := range []string{"--ap-name", "--yes", "--dry-run", "ap_name"} {
				if !strings.Contains(got.stdout, want) {
					t.Errorf("help does not mention %q:\n%s", want, got.stdout)
				}
			}
		})
	}
}

// resetLeaves is every leaf of the reset tree with the RPC each one must reach and the
// word it reports. Both leaves are driven through the same cases on purpose: the guards
// are shared code, and a test per leaf is what would let one of them quietly lose a guard.
var resetLeaves = []struct {
	leaf    string
	rpc     string
	sent    string
	wouldDo string
}{
	{leaf: "ap", rpc: apResetRPC, sent: "reset sent", wouldDo: "would reset"},
	{leaf: "capwap", rpc: capwapResetRPC, sent: "capwap reset sent", wouldDo: "would reset the CAPWAP session"},
}

// The prompt, the cancellation and the two ways of getting past it, against a stub
// controller. Each case asserts what reached the RPC, because the exit code alone
// cannot tell a canceled run from one that acted and said so.
func TestResetConfirmation(t *testing.T) {
	answers := []struct {
		name     string
		stdin    string
		extra    []string
		wantSent int
	}{
		{name: "answered no", stdin: "n\n", wantSent: 0},
		{name: "answered with nothing", stdin: "\n", wantSent: 0},
		{name: "answered yes", stdin: "y\n", wantSent: 1},
		{name: "answered the whole word", stdin: "YES\n", wantSent: 1},
		{name: "flag instead of a prompt", stdin: "", extra: []string{"--yes"}, wantSent: 1},
	}

	for _, leaf := range resetLeaves {
		for _, tt := range answers {
			t.Run(leaf.leaf+"/"+tt.name, func(t *testing.T) {
				stub := newControllerStub(t, http.StatusOK)

				args := []string{
					"reset", leaf.leaf, "--ap-name", testAPName,
					"-c", stub.addr, "--access-token", fakeToken, "-k",
				}
				args = append(args, tt.extra...)

				got := runCLI(t, tt.stdin, false, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				want := leaf.sent
				if tt.wantSent == 0 {
					want = "canceled"
				}

				if !strings.Contains(got.stdout, want) {
					t.Errorf("stdout %q does not contain %q", got.stdout, want)
				}

				if n := stub.hit(leaf.rpc); n != tt.wantSent {
					t.Errorf("%s arrived %d times, want %d", leaf.rpc, n, tt.wantSent)
				}

				// The other leaf's RPC must never be the one that went out.
				for _, other := range resetLeaves {
					if other.rpc == leaf.rpc {
						continue
					}

					if n := stub.hit(other.rpc); n != 0 {
						t.Errorf("reset %s reached %s %d times", leaf.leaf, other.rpc, n)
					}
				}

				if !strings.Contains(got.stdout, testAPName) {
					t.Errorf("stdout %q does not name the access point", got.stdout)
				}

				// An access point is identified by name and never by address. The resolve
				// answered with one, so nothing but a format string keeps it off the stream.
				if strings.Contains(got.stdout, docMAC) {
					t.Errorf("stdout %q carries an address", got.stdout)
				}

				// Neither RPC declares an output container, so a 204 establishes acceptance
				// and nothing more. Claiming the access point restarted, or that its session
				// is back, would state what was not read.
				for _, overclaim := range []string{"rebooted", "restarted", "reloaded", "rejoined"} {
					if strings.Contains(got.stdout, overclaim) {
						t.Errorf("stdout %q claims %q, which the exchange did not establish",
							got.stdout, overclaim)
					}
				}
			})
		}
	}
}

// A dry run resolves the target and stops. Reading is not changing, so it contacts the
// controller; the assertion is that the RPC does not.
func TestResetDryRunSendsNothing(t *testing.T) {
	for _, leaf := range resetLeaves {
		t.Run(leaf.leaf, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			got := runCLI(t, "", true, "--dry-run", "reset", leaf.leaf, "--ap-name", testAPName,
				"-c", stub.addr, "--access-token", fakeToken, "-k")

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, leaf.wouldDo) || !strings.Contains(got.stdout, testAPName) {
				t.Errorf("stdout = %q, want %q and the access point's name", got.stdout, leaf.wouldDo)
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

// The two shapes an absent name arrives in. A controller holding no access point under the
// name answers 404, and a 200 carrying no row is not a measured shape for this list, so both
// are reported the same way. Neither reaches the RPC, which is the whole reason the resolve
// survives a name-arm write.
func TestResetRefusesANameTheControllerDoesNotHold(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		mentions string
	}{
		{name: "a 404", status: http.StatusNotFound, mentions: "holds no access point named"},
		{name: "an answer with no row", status: http.StatusNoContent, mentions: "holds no access point named"},
	}

	for _, leaf := range resetLeaves {
		for _, tt := range tests {
			t.Run(leaf.leaf+"/"+tt.name, func(t *testing.T) {
				stub := newControllerStub(t, tt.status)

				got := runCLI(t, "", false, "reset", leaf.leaf, "--ap-name", testAPName,
					"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

				if got.code != ExitFailure {
					t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
				}

				if !strings.Contains(got.stderr, tt.mentions) {
					t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
				}

				if n := stub.hit(leaf.rpc); n != 0 {
					t.Errorf("%s was sent anyway", leaf.rpc)
				}
			})
		}
	}
}

// urfave consumes the next argv element as a flag's value with no hyphen test, so a name
// left out makes --yes the target. The cost is bounded and observable: one keyed read, the
// name reported as one the controller does not hold, and the consent token spent on the
// value rather than on the prompt — so a run that was going to be confirmed is not.
func TestAPNameSwallowsTheFlagBehindIt(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLI(t, "", false, "reset", "ap", "--ap-name", "--yes",
		"-c", stub.addr, "--access-token", fakeToken, "-k")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
	}

	if !strings.Contains(got.stderr, "holds no access point named") {
		t.Errorf("stderr %q does not report the name as absent", got.stderr)
	}

	if n := stub.hit("ap-name-mac-map=--yes"); n != 1 {
		t.Errorf("the name was read %d times, want 1", n)
	}

	if n := stub.hit(apResetRPC); n != 0 {
		t.Errorf("%d RPCs were sent for a target that does not exist", n)
	}
}
