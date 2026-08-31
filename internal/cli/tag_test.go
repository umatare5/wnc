package cli

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"sync"
	"testing"
)

// tagLeaves is every kind the set and delete trees carry, with the RESTCONF elements each
// one reaches. The three lists are keyed independently on the controller, so a leaf reaching
// another kind's list would be invisible to a per-kind test — every case below asserts that
// nothing else was touched.
var tagLeaves = []struct {
	leaf  string
	noun  string
	list  string
	field []string
}{
	{
		leaf: "policy-tag", noun: "policy tag",
		list:  "policy-list-entries",
		field: []string{"--description", "written by a test"},
	},
	{
		leaf: "site-tag", noun: "site tag",
		list:  "site-tag-configs",
		field: []string{"--ap-join-profile", "test-join"},
	},
	{
		leaf: "rf-tag", noun: "RF tag",
		list:  "rf-tags",
		field: []string{"--profile-5ghz", "test-rf-5"},
	},
}

const tagName = "test-tag-for-the-suite"

// tagStub answers the keyed read that decides create-from-update and counts every write by
// method and by the element it arrived on. The counts are the assertion that matters: a
// canceled or dry run must leave all of them at zero, and no status code can show that.
type tagStub struct {
	addr string

	mu   sync.Mutex
	hits map[string]int
}

// newTagStub serves a controller on which the tag either exists or does not. An absent tag
// answers 404 on the keyed URL, which is what the controller does — measured on 17.12.8 —
// rather than an empty record.
func newTagStub(t *testing.T, exists bool) *tagStub {
	t.Helper()

	s := &tagStub{hits: map[string]int{}}

	srv := httptest.NewTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		base := path.Base(req.URL.Path)

		s.mu.Lock()
		s.hits[req.Method+" "+base]++
		s.mu.Unlock()

		switch {
		case req.Method == http.MethodGet && strings.Contains(base, "="):
			if !exists {
				w.WriteHeader(http.StatusNotFound)

				return
			}

			w.Header().Set("Content-Type", "application/yang-data+json")
			_, _ = w.Write([]byte(tagRecord(base)))
		case req.Method == http.MethodPost:
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	srv.StartTLS()

	s.addr = srv.Listener.Addr().String()

	return s
}

// tagRecord answers a keyed read with a record the SDK's own payload type accepts. The key
// leaf differs per kind, so it is derived from the element the request arrived on.
func tagRecord(base string) string {
	switch {
	case strings.HasPrefix(base, "rf-tag="):
		return `{"Cisco-IOS-XE-wireless-rf-cfg:rf-tag":[{"tag-name":"` + tagName + `"}]}`
	case strings.HasPrefix(base, "site-tag-config="):
		return `{"Cisco-IOS-XE-wireless-site-cfg:site-tag-config":[{"site-tag-name":"` + tagName + `"}]}`
	default:
		return `{"Cisco-IOS-XE-wireless-wlan-cfg:policy-list-entry":[{"tag-name":"` + tagName + `"}]}`
	}
}

func (s *tagStub) hit(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.hits[key]
}

// requests counts everything that arrived, reads included. A guard that runs before the
// client exists must leave this at zero.
func (s *tagStub) requests() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	for _, v := range s.hits {
		n += v
	}

	return n
}

// writes counts every method that changes the controller, whatever element it arrived on.
func (s *tagStub) writes() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0

	for k, v := range s.hits {
		if !strings.HasPrefix(k, http.MethodGet+" ") {
			n += v
		}
	}

	return n
}

// Every guard below is settled before a client exists, so none of these reaches a
// controller. That is the property being asserted as much as the exit code.
func TestTagWritesRefuseBeforeContactingAnything(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		piped    bool
		mentions string
	}{
		{name: "no name", args: []string{}, mentions: "requires --name"},
		{name: "two names", args: []string{"--name", "a", "--name", "b"}, mentions: "--name given 2 times"},
		{name: "not printable ASCII", args: []string{"--name", "タグ"}, mentions: "printable ASCII"},
		{name: "no controller", args: []string{"--name", tagName}, mentions: "no controller given"},
		{
			name:     "two controllers",
			args:     []string{"--name", tagName, "-c", "h1", "-c", "h2", "--access-token", fakeToken},
			mentions: "one controller",
		},
		{
			name:     "piped stdin cannot answer the prompt",
			args:     []string{"--name", tagName, "-c", "h1", "--access-token", fakeToken},
			piped:    true,
			mentions: "--yes",
		},
	}

	for _, verb := range []string{"set", "delete"} {
		for _, leaf := range tagLeaves {
			for _, tt := range cases {
				t.Run(verb+"/"+leaf.leaf+"/"+tt.name, func(t *testing.T) {
					got := runCLI(t, "", tt.piped, append([]string{verb, leaf.leaf}, tt.args...)...)

					if got.code != ExitUsage {
						t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
					}

					if !strings.Contains(got.stderr, tt.mentions) {
						t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
					}

					if got.stdout != "" {
						t.Errorf("a refused write wrote to stdout: %q", got.stdout)
					}

					if strings.Contains(got.stderr, "cause=") {
						t.Errorf("a refused write contacted a controller: %q", got.stderr)
					}
				})
			}
		}
	}
}

// A combination the controller refuses is refused here, before a client exists. Each case
// below was measured against 17.12.8 as a 400 or as a write that never happened, and each
// message names both flags so the operator does not have to guess which half to drop.
func TestTagSetRefusesAContradictoryCombination(t *testing.T) {
	tests := []struct {
		name     string
		leaf     string
		extra    []string
		mentions string
	}{
		{
			// Both are keys of the wlan-policy list, so one alone binds nothing. Before the
			// guard this printed "updated" having sent no write at all.
			name: "wlan without a policy profile", leaf: "policy-tag",
			extra: []string{"--wlan", "test-corp"}, mentions: "--policy-profile",
		},
		{
			name: "policy profile without a wlan", leaf: "policy-tag",
			extra: []string{"--policy-profile", "test-pol"}, mentions: "--wlan",
		},
		{
			// flex-profile declares when "../is-local-site = 'false'", so the controller
			// answers 400 for the pair on every release measured.
			name: "a local site with a flex profile", leaf: "site-tag",
			extra: []string{"--local-site", "--flex-profile", "test-flex-profile01"}, mentions: "--flex-profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.leaf+"/"+tt.name, func(t *testing.T) {
			stub := newTagStub(t, true)

			args := append([]string{
				"set", tt.leaf, "--name", tagName,
				"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes",
			}, tt.extra...)

			got := runCLI(t, "", false, args...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d (stdout %q, stderr %q)", got.code, ExitUsage, got.stdout, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not name %q", got.stderr, tt.mentions)
			}

			if got.stdout != "" {
				t.Errorf("a refused combination wrote to stdout: %q", got.stdout)
			}

			// The point of running it before the client: nothing reaches the controller,
			// not even the existence read.
			if n := stub.requests(); n != 0 {
				t.Errorf("a refused combination made %d requests", n)
			}
		})
	}
}

// The pair that IS complete must still reach the wire, or the guard above has replaced one
// silent miss with another.
func TestTagSetSendsACompleteWLANBinding(t *testing.T) {
	stub := newTagStub(t, true)

	got := runCLI(t, "", false, "set", "policy-tag", "--name", tagName,
		"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes",
		"--wlan", "test-corp", "--policy-profile", "test-pol")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
	}

	if n := stub.writes(); n != 1 {
		t.Errorf("%d writes, want 1", n)
	}
}

// A name the controller does not hold is created and one it holds is updated, and the
// report says which. The verb that went out is asserted too: a create must POST to the
// list and an update must not.
func TestTagSetCreatesOrUpdates(t *testing.T) {
	for _, leaf := range tagLeaves {
		for _, tt := range []struct {
			name     string
			exists   bool
			wantSays string
			wantPost int
		}{
			{name: "absent", exists: false, wantSays: "created", wantPost: 1},
			{name: "present", exists: true, wantSays: "updated", wantPost: 0},
		} {
			t.Run(leaf.leaf+"/"+tt.name, func(t *testing.T) {
				stub := newTagStub(t, tt.exists)

				args := append([]string{
					"set", leaf.leaf, "--name", tagName,
					"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes",
				}, leaf.field...)

				got := runCLI(t, "", false, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if !strings.Contains(got.stdout, tt.wantSays) {
					t.Errorf("stdout %q does not say %q", got.stdout, tt.wantSays)
				}

				if !strings.Contains(got.stdout, leaf.noun) || !strings.Contains(got.stdout, tagName) {
					t.Errorf("stdout %q does not name the tag and its kind", got.stdout)
				}

				if n := stub.hit(http.MethodPost + " " + leaf.list); n != tt.wantPost {
					t.Errorf("POST %s arrived %d times, want %d", leaf.list, n, tt.wantPost)
				}

				// One list per kind: reaching another kind's list would bind the wrong thing.
				for _, other := range tagLeaves {
					if other.list == leaf.list {
						continue
					}

					if n := stub.hit(http.MethodPost + " " + other.list); n != 0 {
						t.Errorf("set %s posted to %s", leaf.leaf, other.list)
					}
				}
			})
		}
	}
}

// An update with no field named changes nothing, and says so rather than reporting a write
// it did not make.
func TestTagSetWithNoFieldChangesNothing(t *testing.T) {
	for _, leaf := range tagLeaves {
		t.Run(leaf.leaf, func(t *testing.T) {
			stub := newTagStub(t, true)

			got := runCLI(t, "", false, "set", leaf.leaf, "--name", tagName,
				"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if !strings.Contains(got.stdout, "no field given to change") {
				t.Errorf("stdout = %q", got.stdout)
			}

			if n := stub.writes(); n != 0 {
				t.Errorf("a run with no field made %d writes", n)
			}
		})
	}
}

// Deleting a name the controller does not hold is a failure rather than a silent success,
// and the DELETE is not sent.
func TestTagDeleteRefusesAnAbsentName(t *testing.T) {
	for _, leaf := range tagLeaves {
		t.Run(leaf.leaf, func(t *testing.T) {
			stub := newTagStub(t, false)

			got := runCLI(t, "", false, "delete", leaf.leaf, "--name", tagName,
				"-c", stub.addr, "--access-token", fakeToken, "-k", "--yes")

			if got.code != ExitFailure {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
			}

			if !strings.Contains(got.stderr, "holds no "+leaf.noun) {
				t.Errorf("stderr %q does not report the absence plainly", got.stderr)
			}

			if n := stub.writes(); n != 0 {
				t.Errorf("a refused delete made %d writes", n)
			}
		})
	}
}

// A dry run reads, reports and writes nothing. Reading is not changing, so it does contact
// the controller; the assertion is that no write does.
func TestTagDryRunWritesNothing(t *testing.T) {
	for _, verb := range []string{"set", "delete"} {
		for _, leaf := range tagLeaves {
			t.Run(verb+"/"+leaf.leaf, func(t *testing.T) {
				stub := newTagStub(t, true)

				// A set names a field: with none, an existing tag reports that there is
				// nothing to change, which is settled before the dry-run branch is reached.
				args := []string{
					"--dry-run", verb, leaf.leaf, "--name", tagName,
					"-c", stub.addr, "--access-token", fakeToken, "-k",
				}
				if verb == "set" {
					args = append(args, leaf.field...)
				}

				got := runCLI(t, "", true, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if !strings.Contains(got.stdout, "would ") {
					t.Errorf("stdout %q does not report what would happen", got.stdout)
				}

				if n := stub.writes(); n != 0 {
					t.Errorf("a dry run made %d writes", n)
				}
			})
		}
	}
}

// The prompt and the two ways past it. A declined write is the command working, so it exits
// zero and says so.
func TestTagConfirmation(t *testing.T) {
	for _, verb := range []string{"set", "delete"} {
		for _, tt := range []struct {
			name       string
			stdin      string
			extra      []string
			wantWrites int
		}{
			{name: "answered no", stdin: "n\n", wantWrites: 0},
			{name: "answered with nothing", stdin: "\n", wantWrites: 0},
			{name: "answered yes", stdin: "y\n", wantWrites: 1},
			{name: "flag instead of a prompt", extra: []string{"--yes"}, wantWrites: 1},
		} {
			t.Run(verb+"/"+tt.name, func(t *testing.T) {
				stub := newTagStub(t, true)

				args := append([]string{
					verb, "rf-tag", "--name", tagName,
					"-c", stub.addr, "--access-token", fakeToken, "-k",
				}, tt.extra...)
				if verb == "set" {
					args = append(args, "--description", "written by a test")
				}

				got := runCLI(t, tt.stdin, false, args...)

				if got.code != ExitOK {
					t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
				}

				if tt.wantWrites == 0 && !strings.Contains(got.stdout, "canceled") {
					t.Errorf("a declined write did not say so: %q", got.stdout)
				}

				if n := stub.writes(); n != tt.wantWrites {
					t.Errorf("%d writes, want %d", n, tt.wantWrites)
				}
			})
		}
	}
}

// urfave trims a positional argument and does not trim a flag value — measured against
// v3.11.0 — so a padded --name reaches the tag-name validator and is refused with nothing
// sent. An inner space is the key leaf's own business and passes.
func TestTagNameKeepsTheSpacesAFlagValueCarries(t *testing.T) {
	for _, tt := range []struct {
		name     string
		given    string
		wantCode int
		mentions string
	}{
		{
			name: "padded", given: "  padded  ",
			wantCode: ExitUsage, mentions: "must not begin or end with a space",
		},
		{name: "an inner space", given: "pad ded", wantCode: ExitOK, mentions: "pad ded"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stub := newTagStub(t, false)

			got := runCLI(t, "", false, "--dry-run", "set", "rf-tag", "--name", tt.given,
				"-c", stub.addr, "--access-token", fakeToken, "-k")

			if got.code != tt.wantCode {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, tt.wantCode, got.stderr)
			}

			if !strings.Contains(got.stdout+got.stderr, tt.mentions) {
				t.Errorf("neither stream mentions %q: %q %q", tt.mentions, got.stdout, got.stderr)
			}

			if tt.wantCode == ExitUsage && stub.requests() != 0 {
				t.Errorf("a refused name reached the controller %d times", stub.requests())
			}
		})
	}
}

// The combination guard runs ahead of the name, so a contradiction is named before a
// missing name is. Both faults are exit 2 and neither reaches a controller, so the message
// is the only thing that shows which one was decided first — and the order is the design.
func TestTagSetNamesTheContradictionBeforeTheMissingName(t *testing.T) {
	got := runCLI(t, "", false, "set", "policy-tag", "--wlan", "test-corp")

	if got.code != ExitUsage {
		t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
	}

	if !strings.Contains(got.stderr, "--policy-profile") {
		t.Errorf("stderr %q names the missing name rather than the contradiction", got.stderr)
	}
}
