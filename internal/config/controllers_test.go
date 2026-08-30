package config

import (
	"strings"
	"testing"
)

func TestParseControllersAccepts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		hosts []string
		want  []Target
	}{
		{
			name:  "host only",
			hosts: []string{"192.168.0.231"},
			want:  []Target{{Name: "192.168.0.231", Host: "192.168.0.231", Token: fakeToken}},
		},
		{
			name:  "host and port",
			hosts: []string{"192.168.0.231:443"},
			want:  []Target{{Name: "192.168.0.231:443", Host: "192.168.0.231:443", Token: fakeToken}},
		},
		{
			name:  "dns name",
			hosts: []string{"wnc1.example.internal"},
			want:  []Target{{Name: "wnc1.example.internal", Host: "wnc1.example.internal", Token: fakeToken}},
		},
		{
			name:  "bracketed ipv6",
			hosts: []string{"[2001:db8::1]"},
			want:  []Target{{Name: "[2001:db8::1]", Host: "[2001:db8::1]", Token: fakeToken}},
		},
		{
			name:  "bracketed ipv6 with port",
			hosts: []string{"[2001:db8::1]:443"},
			want:  []Target{{Name: "[2001:db8::1]:443", Host: "[2001:db8::1]:443", Token: fakeToken}},
		},
		{
			name:  "repeated flag with surrounding space",
			hosts: []string{" a.example ", " b.example "},
			want: []Target{
				{Name: "a.example", Host: "a.example", Token: fakeToken},
				{Name: "b.example", Host: "b.example", Token: fakeToken},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseControllers(tt.hosts, fakeToken)
			if err != nil {
				t.Fatalf("ParseControllers: %v", err)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d targets, want %d", len(got), len(tt.want))
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("target[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseControllersRejects(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		hosts []string
		want  string
	}{
		{name: "empty element", hosts: []string{"a.example", ""}, want: "empty element"},
		{name: "https scheme", hosts: []string{"https://a.example"}, want: "scheme prefix"},
		{name: "path", hosts: []string{"a.example/restconf"}, want: "must not contain"},
		{name: "query", hosts: []string{"a.example?x=1"}, want: "must not contain"},
		{name: "fragment", hosts: []string{"a.example#f"}, want: "must not contain"},
		{name: "bare ipv6", hosts: []string{"2001:db8::1"}, want: "must be bracketed"},
		{name: "unclosed bracket", hosts: []string{"[2001:db8::1"}, want: "malformed bracketed"},
		{name: "not an address in brackets", hosts: []string{"[nope]"}, want: "malformed bracketed"},
		{name: "junk after brackets", hosts: []string{"[2001:db8::1]x"}, want: "malformed bracketed"},
		{name: "port out of range", hosts: []string{"a.example:70000"}, want: "0-65535"},
		{name: "port not a number", hosts: []string{"a.example:https"}, want: "0-65535"},
		// The old grammar packed the token after the port. Fed to the new one it is
		// rejected as a port and never quoted, which is how a stale export fails.
		{name: "old host:token habit", hosts: []string{"a.example:" + fakeToken}, want: "0-65535"},
		// A token written where userinfo goes carries no colon, so the port check cannot see it.
		// errUserinfo refuses these before the SDK's own refusal could quote the element.
		{name: "token as userinfo", hosts: []string{fakeToken + "@a.example"}, want: "userinfo"},
		{name: "userinfo pair", hosts: []string{"u:" + fakeToken + "@a.example"}, want: "userinfo"},
		{name: "userinfo before brackets", hosts: []string{fakeToken + "@[2001:db8::1]:443"}, want: "userinfo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := ParseControllers(tt.hosts, fakeToken)
			if err == nil {
				t.Fatalf("ParseControllers(%q) = nil error, want one", tt.hosts)
			}

			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not mention %q", err, tt.want)
			}
		})
	}
}

// The element body may be a credential, so a fault names the element by index and
// states the accepted syntax without ever quoting what was typed.
func TestParseControllersNeverEchoesTheElement(t *testing.T) {
	t.Parallel()

	specs := []string{
		"https://a.example",
		"2001:db8::1",
		"a.example/x",
		"a.example:70000",
		// The stale form: the token now lands where a port is expected.
		"a.example:" + fakeToken,
		// The userinfo position carries no colon, so only the '@' case sees it and
		// the non-echo property of errUserinfo has to be pinned on its own.
		fakeToken + "@a.example",
		"u:" + fakeToken + "@a.example",
		fakeToken + "@[2001:db8::1]:443",
	}

	for i, spec := range specs {
		t.Run(spec[:min(len(spec), 24)], func(t *testing.T) {
			t.Parallel()

			_, err := ParseControllers([]string{"ok.example", spec}, fakeToken)
			if err == nil {
				t.Fatalf("ParseControllers(%q) = nil error, want one", spec)
			}

			if strings.Contains(err.Error(), fakeToken) {
				t.Errorf("case %d echoed the token: %q", i, err)
			}

			if !strings.Contains(err.Error(), "--controller[2]") {
				t.Errorf("case %d does not name the element by its 1-based index: %q", i, err)
			}
		})
	}
}

func TestTargetsFromFile(t *testing.T) {
	t.Parallel()

	name, host := "lab", "192.168.0.231"

	got, err := TargetsFromFile([]Controller{{Name: &name, Host: &host}}, fakeToken)
	if err != nil {
		t.Fatalf("TargetsFromFile: %v", err)
	}

	if len(got) != 1 || got[0].Name != "lab" || got[0].Host != host {
		t.Fatalf("targets = %+v", got)
	}

	// Without a name the authority labels the rows, matching what the repeated flag produces.
	unnamed, err := TargetsFromFile([]Controller{{Host: &host}}, fakeToken)
	if err != nil {
		t.Fatalf("TargetsFromFile: %v", err)
	}

	if unnamed[0].Name != host {
		t.Errorf("name = %q, want the host %q", unnamed[0].Name, host)
	}
}

func TestTargetsFromFileRejects(t *testing.T) {
	t.Parallel()

	empty, host := "", "192.168.0.231"
	scheme := "https://192.168.0.231"

	tests := []struct {
		name  string
		entry Controller
		token string
		want  string
	}{
		{name: "no host", token: fakeToken, want: "empty host"},
		{name: "empty host", entry: Controller{Host: &empty}, token: fakeToken, want: "empty host"},
		// The token is the file's rather than the entry's, so its absence is one fault
		// for the whole list instead of one per entry.
		{name: "no token", entry: Controller{Host: &host}, want: "no token"},
		{name: "scheme in host", entry: Controller{Host: &scheme}, token: fakeToken, want: "scheme prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := TargetsFromFile([]Controller{tt.entry}, tt.token)
			if err == nil {
				t.Fatalf("TargetsFromFile = nil error, want one")
			}

			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not mention %q", err, tt.want)
			}

			if strings.Contains(err.Error(), fakeToken) {
				t.Errorf("error echoed the token: %q", err)
			}
		})
	}
}
