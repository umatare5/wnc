package cli

import (
	"net/http"
	"net/http/httptest"
	"path"
	"slices"
	"strings"
	"testing"
)

// newRejectingDeauthStub answers the read and rejects the RPC the way a release that does not
// serve it does. It is a stub of its own because controllerStub answers the RPC successfully,
// and because this is the one case where the read has to succeed for the rejection to mean
// what the CLI says it means.
func newRejectingDeauthStub(t *testing.T) string {
	t.Helper()

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch path.Base(req.URL.Path) {
		case "common-oper-data=" + testClientMAC:
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(`{"Cisco-IOS-XE-wireless-client-oper:common-oper-data":[` +
				`{"client-mac":"` + testClientRowMAC + `","co-state":"client-status-run"}]}`))
		case "common-oper-data":
			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(clientCollection))
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	srv.StartTLS()

	return srv.Listener.Addr().String()
}

// Every guard is settled before a client exists, so none of these reaches a controller. The hit
// count is asserted beside the exit code, because a guard that ran too late would still produce
// exit 2 while having already sent the RPC.
func TestDeauthRefusesBeforeContactingAnything(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		piped    bool
		mentions string
	}{
		{name: "no target", args: nil, mentions: "requires --mac or --username"},
		{
			name:     "two addresses",
			args:     []string{"--mac", testClientMAC, "--mac", "00:00:5e:00:53:a2"},
			mentions: "one client per invocation, --mac given 2 times",
		},
		{
			name:     "two usernames",
			args:     []string{"--username", testClientUsername, "--username", testSharedUsername},
			mentions: "one username per invocation, --username given 2 times",
		},
		{
			// The RPC's choice is mandatory and the controller resolves the first arm it finds,
			// so sending both would let the controller pick which one the operator meant.
			name:     "both arms",
			args:     []string{"--mac", testClientMAC, "--username", testClientUsername},
			mentions: "one invocation gives one of them",
		},
		{
			name:     "an empty address",
			args:     []string{"--mac", ""},
			mentions: "must not be empty",
		},
		{
			// An empty username is the value most clients carry, so it would select nearly the
			// whole fleet rather than nothing.
			name:     "an empty username",
			args:     []string{"--username", ""},
			mentions: "must not be empty",
		},
		{
			name:     "no controller",
			args:     []string{"--mac", testClientMAC},
			mentions: "no controller given",
		},
		{
			name: "two controllers",
			args: []string{
				"--mac", testClientMAC, "-c", "192.0.2.1", "-c", "192.0.2.2",
				"--access-token", fakeToken,
			},
			mentions: "acts on one controller",
		},
		{
			name:     "a pipe that could not answer the prompt",
			args:     []string{"--mac", testClientMAC, "-c", "192.0.2.1", "--access-token", fakeToken},
			piped:    true,
			mentions: "stdin is not a terminal",
		},
		{
			name:     "a positional instead of a flag",
			args:     []string{testClientMAC, "-c", "192.0.2.1", "--access-token", fakeToken},
			mentions: "takes no positional arguments, 1 given: use --mac",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			got := runCLI(t, "", tt.piped, append([]string{"deauth"}, tt.args...)...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}

			if n := stub.hit(deauthRPC); n != 0 {
				t.Errorf("the RPC was sent %d times by a refused run", n)
			}
		})
	}
}

// The read before the post is the whole reason this command can report truthfully: the RPC
// answers 204 for an address associated to nothing exactly as it does for a client it dropped.
// So an address the controller does not hold must be refused here, and the RPC must not arrive.
func TestDeauthRefusesAnAddressTheControllerDoesNotHold(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLI(t, "", false, "deauth", "--mac", "00:00:5e:00:53:a3",
		"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
	}

	if !strings.Contains(got.stderr, "holds no client at") {
		t.Errorf("stderr %q does not report the absence", got.stderr)
	}

	if n := stub.hit(deauthRPC); n != 0 {
		t.Errorf("the RPC was sent %d times for an address nothing holds", n)
	}
}

// A dry run and a declined prompt are both completed runs that changed nothing. The dry run
// still reads, so it doubles as an existence probe — which is the only thing it can report,
// since the RPC's answer says nothing either way.
func TestDeauthReportsWithoutSendingWhenItMust(t *testing.T) {
	tests := []struct {
		name string
		// root carries --dry-run, which is Local to the root and so must precede the command.
		root []string
		// target defaults to the address arm, so a row naming none exercises that one.
		target []string
		stdin  string
		says   string
	}{
		{name: "dry run", root: []string{"--dry-run"}, says: "would deauthenticate"},
		{name: "declined prompt", stdin: "n\n", says: "canceled"},
		{name: "an answer that is not yes", stdin: "later\n", says: "canceled"},
		{
			name: "dry run on the username arm",
			root: []string{"--dry-run"}, target: []string{"--username", testSharedUsername},
			says: "2 clients authenticated as " + testSharedUsername,
		},
		{
			name:   "declined prompt on the username arm",
			target: []string{"--username", testClientUsername}, stdin: "n\n", says: "canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			target := tt.target
			if target == nil {
				target = []string{"--mac", testClientMAC}
			}

			args := append(slices.Clone(tt.root), "deauth")
			args = append(args, target...)
			args = append(args, "-c", stub.addr, "--access-token", fakeToken, "-k")

			got := runCLI(t, tt.stdin, false, args...)

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, tt.says) {
				t.Errorf("stdout %q does not say %q", got.stdout, tt.says)
			}

			if n := stub.hit(deauthRPC); n != 0 {
				t.Errorf("the RPC was sent %d times, want none", n)
			}
		})
	}
}

// The reported line says "sent" and not a past participle: this RPC declares no output
// container, so a 204 establishes that the instruction was accepted and nothing more.
func TestDeauthSends(t *testing.T) {
	tests := []struct {
		name string
		// target defaults to the address arm, so a row naming none exercises that one.
		target []string
		extra  []string
		stdin  string
	}{
		{name: "yes flag", extra: []string{"--yes"}},
		{name: "typed yes", stdin: "y\n"},
		{name: "typed yes spelt out", stdin: "YES\n"},
		{
			// One post per invocation on this arm too: the RPC takes the username and the
			// controller decides how many sessions that is, so the CLI must not loop over them.
			name:   "the username arm posts once for several sessions",
			target: []string{"--username", testSharedUsername}, extra: []string{"--yes"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			target := tt.target
			if target == nil {
				target = []string{"--mac", testClientMAC}
			}

			args := append([]string{"deauth"}, target...)
			args = append(args, "-c", stub.addr, "--access-token", fakeToken, "-k")
			args = append(args, tt.extra...)

			got := runCLI(t, tt.stdin, false, args...)

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, "deauthenticate sent") {
				t.Errorf("stdout %q does not report the dispatch", got.stdout)
			}

			if n := stub.hit(deauthRPC); n != 1 {
				t.Errorf("the RPC was sent %d times, want once", n)
			}
		})
	}
}

// A release that does not serve the operation answers 400, not 404, so the status alone reaches
// the operator as an opaque refusal. The target was read on the same controller a moment
// earlier, which is what licenses naming the release instead. Both arms are covered because the
// re-wording is one call each makes, and an arm that dropped it would fail on neither the other's
// test nor the unit test of the classification.
func TestDeauthWordsARejectedPathAsAnAbsentOperation(t *testing.T) {
	tests := map[string][]string{
		"the address arm":  {"--mac", testClientMAC},
		"the username arm": {"--username", testClientUsername},
	}

	for name, target := range tests {
		t.Run(name, func(t *testing.T) {
			addr := newRejectingDeauthStub(t)

			args := append([]string{"deauth"}, target...)
			args = append(args, "-c", addr, "--access-token", fakeToken, "-k", "--yes")

			got := runCLI(t, "", false, args...)

			if got.code != ExitFailure {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
			}

			if !strings.Contains(got.stderr, "before 17.15") {
				t.Errorf("stderr %q does not name the release the operation is absent on", got.stderr)
			}

			if strings.Contains(got.stdout, "sent") {
				t.Errorf("a rejected post was reported as sent: %q", got.stdout)
			}
		})
	}
}

// The prompt names the client and the controller, and says the client comes back on its own.
// Overstating a single-client deauthentication would be a defect of its own.
func TestDeauthPromptNamesTheTargetAndTheRecovery(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLI(t, "n\n", false, "deauth", "--mac", testClientMAC,
		"-c", stub.addr, "--access-token", fakeToken, "-k")

	for _, want := range []string{testClientRowMAC, stub.addr, "reconnects on its own"} {
		if !strings.Contains(got.stdout, want) {
			t.Errorf("the prompt %q does not carry %q", got.stdout, want)
		}
	}
}

// The username arm's resolve adds one thing to what the operator typed, and the prompt is where
// it belongs: the RPC's user-name leaf states no cardinality, so how many sessions the controller
// holds is the operator's only warning of the blast radius. The singular row is what keeps the
// count from being rendered as a bare plural on the common case.
func TestDeauthPromptNamesHowManySessionsAUsernameHolds(t *testing.T) {
	tests := []struct {
		name     string
		username string
		says     string
	}{
		{name: "one session", username: testClientUsername, says: "1 client authenticated as "},
		{name: "two sessions", username: testSharedUsername, says: "2 clients authenticated as "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			got := runCLI(t, "n\n", false, "deauth", "--username", tt.username,
				"-c", stub.addr, "--access-token", fakeToken, "-k")

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			for _, want := range []string{tt.says + tt.username, stub.addr, "reconnects on its own"} {
				if !strings.Contains(got.stdout, want) {
					t.Errorf("the prompt %q does not carry %q", got.stdout, want)
				}
			}
		})
	}
}

// A username the controller holds no session under is refused before the post, for the reason the
// address arm is: the RPC answers 204 either way, so a reported deauthentication and a typo would
// otherwise be the same output.
func TestDeauthRefusesAUsernameTheControllerDoesNotHold(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLI(t, "", false, "deauth", "--username", "test-nobody",
		"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
	}

	if !strings.Contains(got.stderr, "holds no client authenticated as") {
		t.Errorf("stderr %q does not report the absence", got.stderr)
	}

	if n := stub.hit(deauthRPC); n != 0 {
		t.Errorf("the RPC was sent %d times for a username nothing holds", n)
	}
}

// WNC_USERNAME is the controller login generate-token reads, and this flag shares its name with
// that one. So the variable must not reach here: an operator with it exported would otherwise
// have a deauth target chosen for them, and the sessions it selected would be whichever clients
// happened to authenticate under the controller account's name.
func TestDeauthDoesNotTakeItsTargetFromTheControllerAccountVariable(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLIWithEnv(t, map[string]string{"WNC_USERNAME": testClientUsername}, "", false,
		"deauth", "-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
	}

	if !strings.Contains(got.stderr, "requires --mac or --username") {
		t.Errorf("stderr %q does not report the missing target", got.stderr)
	}

	if n := stub.hit(deauthRPC); n != 0 {
		t.Errorf("the RPC was sent %d times for a target nobody named", n)
	}
}
