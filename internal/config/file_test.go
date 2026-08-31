package config

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// fakeToken looks like a Basic auth token but decodes to nothing real. It exists so
// a test can assert the value never reaches a diagnostic.
const fakeToken = "TestToken0123456789ABCDEF=="

func writeConfig(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	return path
}

func TestLoadAcceptsAFullFile(t *testing.T) {
	t.Parallel()

	path := writeConfig(t, `{
	  "note": "WNC1",
	  "timeout": "30s",
	  "insecure": true,
	  "format": "json",
	  "log_level": "debug",
	  "token": "`+fakeToken+`",
	  "controllers": [
	    {"name": "WNC1", "host": "192.168.0.1", "note": "retiring"}
	  ]
	}`)

	file, err := Load(path, true, nil)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if file.Timeout == nil || file.Timeout.Duration() != 30*time.Second {
		t.Errorf("timeout = %v, want 30s", file.Timeout)
	}

	if file.Insecure == nil || !*file.Insecure {
		t.Errorf("insecure = %v, want true", file.Insecure)
	}

	if len(file.Controllers) != 1 || file.Token == nil {
		t.Fatalf("controllers = %#v", file.Controllers)
	}

	if got := *file.Controllers[0].Name; got != "WNC1" {
		t.Errorf("name = %q", got)
	}
}

// A key the file never set must stay distinguishable from one set to a zero value,
// which is what lets a flag default win only where the file is genuinely silent.
func TestLoadKeepsOmissionDistinctFromZero(t *testing.T) {
	t.Parallel()

	omitted, err := Load(writeConfig(t, `{"note":"n"}`), true, nil)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	explicit, err := Load(writeConfig(t, `{"insecure":false,"format":"table"}`), true, nil)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if omitted.Insecure != nil {
		t.Errorf("omitted insecure = %v, want nil", omitted.Insecure)
	}

	if explicit.Insecure == nil || *explicit.Insecure {
		t.Errorf("explicit insecure = %v, want a pointer to false", explicit.Insecure)
	}
}

func TestLoadIsStrict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		is   error
		want string
	}{
		{
			name: "unknown key",
			body: `{"timeuot":"30s"}`,
			is:   json.ErrUnknownName,
			want: "/timeuot",
		},
		{
			name: "duplicate key",
			body: `{"note":"a","note":"b"}`,
			is:   jsontext.ErrDuplicateName,
		},
		{
			// json/v2 matches names case-sensitively where v1 did not, so a
			// capitalized key is an unknown one rather than a silent match.
			name: "key differing only in case",
			body: `{"Timeout":"30s"}`,
			is:   json.ErrUnknownName,
			want: "/Timeout",
		},
		{
			name: "comment",
			body: "{\"note\":\"n\"} // why",
			want: "invalid character",
		},
		{
			name: "trailing comma",
			body: `{"note":"n",}`,
			want: "invalid character",
		},
		{
			name: "trailing value",
			body: `{"note":"a"} {"note":"b"}`,
			want: "invalid character",
		},
		{
			// The parser's own message quotes the text it could not parse, so the
			// fault carries the pointer and the Go type instead of the cause.
			name: "duration without a unit",
			body: `{"timeout":"30"}`,
			want: "cannot decode into",
		},
		{
			// A bare integer is rejected before UnmarshalText is reached. The
			// rejection still matters: a numeric decode would have read 30 as
			// nanoseconds.
			name: "duration as a bare integer",
			body: `{"timeout":30}`,
			want: "/timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := Load(writeConfig(t, tt.body), true, nil)
			if err == nil {
				t.Fatalf("Load(%s) = nil error, want one", tt.body)
			}

			if tt.is != nil && !errors.Is(err, tt.is) {
				t.Errorf("error %v does not match sentinel %v", err, tt.is)
			}

			if tt.want != "" && !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not mention %q", err, tt.want)
			}
		})
	}
}

// json/v2's SemanticError.Error appends the offending scalar for some decodes and
// the cause quotes it for others, so the message is rebuilt from the pointer rather
// than taken as given.
func TestLoadNeverEchoesATokenValue(t *testing.T) {
	t.Parallel()

	bodies := []string{
		`{"token":` + `12345` + `}`,
		`{"token":["` + fakeToken + `"]}`,
		`{"token":{"v":"` + fakeToken + `"}}`,
		`{"timeout":"` + fakeToken + `"}`,
	}

	for _, body := range bodies {
		t.Run(body[:min(len(body), 40)], func(t *testing.T) {
			t.Parallel()

			_, err := Load(writeConfig(t, body), true, nil)
			if err == nil {
				t.Fatalf("Load(%s) = nil error, want one", body)
			}

			// The token is the only credential a configuration file carries, so the
			// literal written into the body above is the whole of what an echo prints.
			if strings.Contains(err.Error(), fakeToken) {
				t.Errorf("error echoed the credential: %q", err)
			}
		})
	}
}

func TestLoadMissingFile(t *testing.T) {
	t.Parallel()

	missing := filepath.Join(t.TempDir(), "absent.json")

	if _, err := Load(missing, true, nil); err == nil {
		t.Error("an explicitly named missing file must be an error")
	}

	file, err := Load(missing, false, nil)
	if err != nil {
		t.Errorf("a missing default file must be ignored, got %v", err)
	}

	if file.Controllers != nil {
		t.Errorf("controllers = %#v, want nil", file.Controllers)
	}
}

func TestLoadWarnsOnPermissiveMode(t *testing.T) {
	t.Parallel()

	if runtime.GOOS == "windows" {
		t.Skip("permission bits are not meaningful on windows")
	}

	path := writeConfig(t, `{"note":"n"}`)
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("Chmod: %v", err)
	}

	var warnings []string
	if _, err := Load(path, true, func(m string) { warnings = append(warnings, m) }); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(warnings) != 1 || !strings.Contains(warnings[0], "0600") {
		t.Errorf("warnings = %#v", warnings)
	}
}

func TestLoadDoesNotWarnOnOwnerOnlyMode(t *testing.T) {
	t.Parallel()

	var warnings []string
	if _, err := Load(
		writeConfig(t, `{"note":"n"}`),
		true,
		func(m string) { warnings = append(warnings, m) },
	); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(warnings) != 0 {
		t.Errorf("warnings = %#v, want none", warnings)
	}
}

func TestPathPrecedence(t *testing.T) {
	t.Parallel()

	got, explicit := Path("/etc/wnc.json")
	if got != "/etc/wnc.json" || !explicit {
		t.Errorf("Path(flag) = %q, %v", got, explicit)
	}

	got, explicit = Path("")
	if explicit {
		t.Errorf("Path(\"\") reported an explicit choice: %q", got)
	}
}
