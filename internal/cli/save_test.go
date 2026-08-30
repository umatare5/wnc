package cli

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

// newFailingSaveStub answers every request with the refusal a controller gives when it will not
// save, and returns its address. It is a stub of its own because controllerStub answers the
// save successfully, which is what every other case here needs.
func newFailingSaveStub(t *testing.T) string {
	t.Helper()

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	srv.StartTLS()

	return srv.Listener.Addr().String()
}

// Every guard is settled before a client exists, so none of these reaches a controller. The
// hit count is asserted alongside the exit code: this command names no target and reads
// nothing first, so the guards are the whole of what stands between an invocation and a write.
func TestSaveConfigRefusesBeforeContactingAnything(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		piped    bool
		mentions string
	}{
		{name: "no controller", args: nil, mentions: "no controller given"},
		{
			name:     "two controllers",
			args:     []string{"-c", "192.0.2.1", "-c", "192.0.2.2", "--access-token", fakeToken},
			mentions: "acts on one controller",
		},
		{
			name:     "a pipe that could not answer the prompt",
			args:     []string{"-c", "192.0.2.1", "--access-token", fakeToken},
			piped:    true,
			mentions: "stdin is not a terminal",
		},
		{
			name:     "a positional instead of a flag",
			args:     []string{"stray", "-c", "192.0.2.1", "--access-token", fakeToken},
			mentions: "takes no positional arguments",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			got := runCLI(t, "", tt.piped, append([]string{"save-config"}, tt.args...)...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}

			if n := stub.hit(saveConfigRPC); n != 0 {
				t.Errorf("the save was sent %d times by a refused run", n)
			}
		})
	}
}

// A dry run and a declined prompt are both completed runs that changed nothing, and the only
// proof of the second half is that the RPC never arrived.
func TestSaveConfigReportsWithoutSendingWhenItMust(t *testing.T) {
	tests := []struct {
		name string
		// root carries --dry-run, which is Local to the root and so must precede the command.
		root  []string
		stdin string
		says  string
	}{
		{name: "dry run", root: []string{"--dry-run"}, says: "would save the running configuration"},
		{name: "declined prompt", stdin: "n\n", says: "canceled"},
		{name: "an answer that is not yes", stdin: "later\n", says: "canceled"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			args := append(slices.Clone(tt.root),
				"save-config", "-c", stub.addr, "--access-token", fakeToken, "-k")

			got := runCLI(t, tt.stdin, false, args...)

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, tt.says) {
				t.Errorf("stdout %q does not say %q", got.stdout, tt.says)
			}

			if n := stub.hit(saveConfigRPC); n != 0 {
				t.Errorf("the save was sent %d times, want none", n)
			}
		})
	}
}

// The prompt has to state what the operator is agreeing to, because the RPC takes no target
// and persists whatever else is on the controller. That is the one thing about this command an
// operator cannot infer from what they typed.
func TestSaveConfigPromptNamesTheBlastRadius(t *testing.T) {
	stub := newControllerStub(t, http.StatusOK)

	got := runCLI(t, "n\n", false,
		"save-config", "-c", stub.addr, "--access-token", fakeToken, "-k")

	for _, want := range []string{stub.addr, "including changes this CLI did not make"} {
		if !strings.Contains(got.stdout, want) {
			t.Errorf("the prompt %q does not carry %q", got.stdout, want)
		}
	}
}

// --yes and a typed yes both reach the controller exactly once, and the reported line claims a
// completion rather than a dispatch: this is the one write whose RPC answers with the
// controller's own account of what it did.
func TestSaveConfigSends(t *testing.T) {
	tests := []struct {
		name  string
		extra []string
		stdin string
	}{
		{name: "yes flag", extra: []string{"--yes"}},
		{name: "typed yes", stdin: "y\n"},
		{name: "typed yes spelt out", stdin: "YES\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stub := newControllerStub(t, http.StatusOK)

			args := append([]string{"save-config", "-c", stub.addr, "--access-token", fakeToken, "-k"},
				tt.extra...)

			got := runCLI(t, tt.stdin, false, args...)

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, "running configuration saved") {
				t.Errorf("stdout %q does not report the save", got.stdout)
			}

			if n := stub.hit(saveConfigRPC); n != 1 {
				t.Errorf("the save was sent %d times, want once", n)
			}
		})
	}
}

// A controller that refuses the save must not be reported as having saved. Exit 1 and not 3:
// this command reaches one controller, so a failure is total rather than partial.
func TestSaveConfigReportsARefusal(t *testing.T) {
	stub := newFailingSaveStub(t)

	got := runCLI(t, "", false,
		"save-config", "-c", stub, "--access-token", fakeToken, "-k", "--yes")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
	}

	if strings.Contains(got.stdout, "saved") {
		t.Errorf("a refused save was reported as done: %q", got.stdout)
	}

	if !strings.Contains(got.stderr, "saving the running configuration") {
		t.Errorf("stderr %q does not name the operation", got.stderr)
	}
}
