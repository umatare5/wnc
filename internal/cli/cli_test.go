package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/show"
)

// Nothing in this file runs in parallel, and that is deliberate. urfave reads the
// WNC_* variables at parse time, so a developer's own shell would otherwise decide
// what these assertions see; clearing the environment needs t.Setenv, which forbids
// t.Parallel.

// fakeToken looks like a Basic auth token and decodes to nothing real.
const fakeToken = "TestToken0123456789ABCDEF=="

type outcome struct {
	stdout string
	stderr string
	code   int
}

// runCLI drives the real command tree with a neutralized environment and a default
// configuration path that does not exist.
func runCLI(t *testing.T, stdin string, piped bool, args ...string) outcome {
	t.Helper()

	return runCLIWithEnv(t, nil, stdin, piped, args...)
}

// runCLIWithEnv is runCLI with named variables restored after the neutralization. It
// exists because the neutralization has to happen inside the helper, which would
// otherwise overwrite anything a test set beforehand.
func runCLIWithEnv(t *testing.T, env map[string]string, stdin string, piped bool, args ...string) outcome {
	t.Helper()

	for _, k := range []string{
		"WNC_CONFIG", "WNC_CONTROLLER", "WNC_ACCESS_TOKEN", "WNC_USERNAME", "WNC_PASSWORD",
	} {
		t.Setenv(k, "")
	}

	for k, v := range env {
		t.Setenv(k, v)
	}

	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	var out, errOut bytes.Buffer

	code := RunWith(append([]string{"wnc"}, args...), Streams{
		Out:    &out,
		Err:    &errOut,
		In:     strings.NewReader(stdin),
		InPipe: piped,
	})

	return outcome{stdout: out.String(), stderr: errOut.String(), code: code}
}

func TestExitCodes(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		want       int
		wantStdout bool
		mentions   string
	}{
		{name: "bare root shows help", args: nil, want: ExitOK, wantStdout: true},
		{name: "version", args: []string{"--version"}, want: ExitOK, wantStdout: true},
		{name: "help flag", args: []string{"--help"}, want: ExitOK, wantStdout: true},
		{name: "show group shows help", args: []string{"show"}, want: ExitOK, wantStdout: true},
		// The nested path names no command, and the suggestion floor accepts it because both
		// words start with g, which keeps an operator who typed it from being left with nothing.
		{
			name: "the flattened command's old path", args: []string{"generate"},
			want: ExitUsage, mentions: "generate-token",
		},

		{name: "unknown flag", args: []string{"--bogus"}, want: ExitUsage, mentions: "not defined"},
		{
			name:     "a positional at a leaf",
			args:     []string{"show", "ap", "stray"},
			want:     ExitUsage,
			mentions: "takes no positional arguments",
		},
		{name: "unknown command", args: []string{"bogus"}, want: ExitUsage, mentions: "unknown command"},
		{
			name:     "unknown show subcommand",
			args:     []string{"show", "bogus"},
			want:     ExitUsage,
			mentions: "unknown command",
		},
		{
			// urfave builds its own cli.Exit for this path, which matches none of this
			// package's sentinels; the ExitCoder branch of the mapper catches it and
			// failureLine replaces its text, which would otherwise echo the word.
			name: "help for an unknown topic", args: []string{"--help", "bogus"},
			want: ExitUsage, mentions: "unknown help topic",
		},
		{
			// The help command is hidden, so the word reaches parentAction, which answers it.
			name: "help word", args: []string{"help"},
			want: ExitOK, wantStdout: true,
		},
		{
			name: "help for a leaf", args: []string{"help", "show", "ap"},
			want: ExitOK, wantStdout: true,
		},
		{
			name: "help for an unknown word", args: []string{"help", "bogus"},
			want: ExitUsage, mentions: "unknown command",
		},
		{
			// --dry-run is Local to the root, so a subcommand does not inherit it.
			name: "dry-run is not inherited", args: []string{"show", "ap", "--dry-run"},
			want: ExitUsage, mentions: "not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runCLI(t, "", false, tt.args...)

			if got.code != tt.want {
				t.Errorf("exit = %d, want %d\nstdout: %s\nstderr: %s", got.code, tt.want, got.stdout, got.stderr)
			}

			if tt.wantStdout && got.stdout == "" {
				t.Error("stdout is empty")
			}

			if !tt.wantStdout && got.stdout != "" {
				t.Errorf("a fault wrote to stdout: %q", got.stdout)
			}

			if tt.mentions != "" && !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}
		})
	}
}

// urfave evaluates the help flag before the usage hook, so an explicit --help placed
// first wins over a later unknown flag. Fixing that would mean reordering the
// library's own run loop; it is recorded here as the accepted behavior instead.
func TestHelpBeforeUnknownFlagWins(t *testing.T) {
	first := runCLI(t, "", false, "--help", "--bogus")
	if first.code != ExitOK || first.stdout == "" {
		t.Errorf("--help --bogus: exit %d, stdout %d bytes; want help and 0", first.code, len(first.stdout))
	}

	second := runCLI(t, "", false, "--bogus", "--help")
	if second.code != ExitUsage {
		t.Errorf("--bogus --help: exit %d, want %d", second.code, ExitUsage)
	}
}

// A parse fault must not put the whole help text into a pipe that is carrying data.
func TestUsageFaultsKeepStdoutClean(t *testing.T) {
	for _, args := range [][]string{
		{"--bogus"},
		{"bogus"},
		{"show", "bogus"},
		{"show", "ap", "--bogus"},
		{"generate-token", "--bogus"},
		{"reset", "ap", "stray"},
		{"generate-token", "stray"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			got := runCLI(t, "", false, args...)

			if got.stdout != "" {
				t.Errorf("stdout = %q, want empty", got.stdout)
			}

			if lines := strings.Count(strings.TrimRight(got.stderr, "\n"), "\n"); lines != 0 {
				t.Errorf("stderr has %d extra lines: %q", lines, got.stderr)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stdin string
		piped bool
		want  int
		out   string
	}{
		{
			name: "flags", args: []string{"generate-token", "-u", "admin", "-p", "secret"},
			want: ExitOK, out: "YWRtaW46c2VjcmV0",
		},
		{
			name: "piped stdin", args: []string{"generate-token", "-u", "admin"},
			stdin: "secret\n", piped: true, want: ExitOK, out: "YWRtaW46c2VjcmV0",
		},
		{
			// A password may hold spaces, so only the line ending is stripped.
			name: "piped stdin keeps inner spaces", args: []string{"generate-token", "-u", "admin"},
			stdin: "two words\r\n", piped: true, want: ExitOK, out: "YWRtaW46dHdvIHdvcmRz",
		},
		{name: "no username", args: []string{"generate-token"}, want: ExitUsage},
		{name: "no password and no pipe", args: []string{"generate-token", "-u", "admin"}, want: ExitUsage},
		{
			name: "colon in the username", args: []string{"generate-token", "-u", "a:b", "-p", "x"},
			want: ExitUsage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runCLI(t, tt.stdin, tt.piped, tt.args...)

			if got.code != tt.want {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, tt.want, got.stderr)
			}

			if tt.out != "" && strings.TrimSpace(got.stdout) != tt.out {
				t.Errorf("stdout = %q, want %q", strings.TrimSpace(got.stdout), tt.out)
			}
		})
	}
}

// urfave applies a slice flag's environment value at the level that declares it and
// then appends what a deeper command supplies, so a controller flag sourced from the
// environment would query the host the operator named plus every exported one. The
// flag therefore declares no source and config.Resolve reads the variable itself.
func TestControllerFlagIsNotMergedWithItsEnvironmentVariable(t *testing.T) {
	const named = "192.0.2.90"

	env := []string{"192.0.2.91", "192.0.2.92"}
	got := runCLI(t, "", false, "show", "ap-tag",
		"-c", named, "--access-token", fakeToken, "-t", "1ms")

	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d", got.code, ExitFailure)
	}

	if !strings.Contains(got.stderr, named) {
		t.Errorf("the named controller is missing from stderr: %q", got.stderr)
	}

	for _, h := range env {
		if strings.Contains(got.stderr, h) {
			t.Errorf("an environment host leaked into the run: %q", h)
		}
	}

	// Same invocation with the variable exported: the flag must still win outright.
	got = runCLIWithEnv(t, map[string]string{"WNC_CONTROLLER": strings.Join(env, ",")},
		"", false, "show", "ap-tag", "-c", named, "--access-token", fakeToken, "-t", "1ms")

	if !strings.Contains(got.stderr, named) {
		t.Errorf("the named controller is missing from stderr: %q", got.stderr)
	}

	for _, h := range env {
		if strings.Contains(got.stderr, h) {
			t.Errorf("%q was queried although only %q was named", h, named)
		}
	}
}

// The variable stands in for the flag when the flag is absent, and its comma form is
// the only way one variable can carry several hosts.
func TestControllerEnvironmentVariableAppliesWhenTheFlagIsAbsent(t *testing.T) {
	got := runCLIWithEnv(t, map[string]string{"WNC_CONTROLLER": "192.0.2.91,192.0.2.92"},
		"", false, "show", "ap-tag", "--access-token", fakeToken, "-t", "1ms")
	if got.code != ExitFailure {
		t.Fatalf("exit = %d, want %d", got.code, ExitFailure)
	}

	for _, h := range []string{"192.0.2.91", "192.0.2.92"} {
		if !strings.Contains(got.stderr, h) {
			t.Errorf("%q from the environment was not queried: %q", h, got.stderr)
		}
	}
}

// A rejected password must never appear in a diagnostic, and neither must a rejected
// controller entry.
func TestFaultsNeverEchoACredential(t *testing.T) {
	cases := [][]string{
		{"generate-token", "-u", "a:b", "-p", fakeToken},
		{"show", "ap-tag", "-c", "https://h", "--access-token", fakeToken},
		{"show", "ap-tag", "-c", "2001:db8::1", "--access-token", fakeToken},
		// The pre-split grammar packed the token after the port. Fed to the current
		// one it is rejected as a port, and the token must not reach the message.
		{"show", "ap-tag", "-c", "h:" + fakeToken, "--access-token", fakeToken},
		// The userinfo position carries no colon, so the port check cannot see it.
		{"show", "ap-tag", "-c", fakeToken + "@h", "--access-token", fakeToken},
		{"show", "ap-tag", "-c", "h"},
		// A wrapper that drops its quotes hands the flag and its value over as one argument,
		// which urfave's parser appends verbatim to its own message.
		{"show", "ap-tag", "-c " + fakeToken},
		{"generate-token", "-u", "admin", "-p " + fakeToken},
		// The mysql idiom: a value written attached to its short flag. The flag package
		// appends the whole word, and trimArg cuts at ':' and ' ' alone, so only the
		// declaration lookup separates the name from the value.
		{"generate-token", "-u", "admin", "-p" + fakeToken},
		{"show", "ap-tag", "-t" + fakeToken},
		// The same shape spelt entirely in the characters a flag name uses. The name is
		// withheld here rather than cut, because -p is declared and reporting it would
		// leave the rest of the word in the message.
		{"generate-token", "-u", "admin", "-psecret-password"},
		// An attached value behind a letter no command declares. Nothing here can be
		// named, so nothing is.
		{"show", "ap-tag", "-z" + fakeToken},
		// A long flag whose name is a stray colon leaves trimArg nothing but the dash.
		// Without the empty check the first-character lookup would index past the end of
		// the string. One dash instead of two makes it a positional argument, refused by
		// count before any of this.
		{"show", "ap-tag", "--:" + fakeToken},
		// -t is --timeout and --access-token has no short name, so this is the slip that
		// reaches a parse the flag itself performs. urfave quotes the value and the
		// duration parser quotes it again, which printed it twice.
		{"show", "-t", fakeToken, "ap"},
		// A leftover word is refused by count and never repeated, which is what makes a
		// wrapper's misplaced password on generate-token safe to report.
		{"generate-token", "-u", "admin", fakeToken},
		{"show", "ap-tag", "-c", "h", "--access-token", fakeToken, fakeToken},
		// A token pasted where a subcommand belongs, which a group reports without the word.
		{"show", fakeToken},
		{fakeToken},
		// A word after --help is a help topic, and urfave answers an unknown one with
		// its own ExitCoder whose text carries the word — the one error this package
		// does not build, replaced wholesale by failureLine.
		{"--help", fakeToken},
		{"show", "-h", fakeToken},
		{"show", "ap", "--help", fakeToken},
		// The enum-valued slots: -b, -o, -f and -r sit beside -c and -t on the same
		// commands, so each is a mis-paste target and its fault names no value.
		{"show", "ap", "-c", "h", "--access-token", fakeToken, "--sort-by", fakeToken},
		{"show", "ap", "-c", "h", "--access-token", fakeToken, "--sort-order", fakeToken},
		{"show", "ap", "-c", "h", "--access-token", fakeToken, "--format", fakeToken},
		{"--log-level", fakeToken, "show", "ap"},
		{"show", "client", "-c", "h", "--access-token", fakeToken, "--radio", fakeToken},
	}

	for _, args := range cases {
		t.Run(strings.Join(args[:min(2, len(args))], " "), func(t *testing.T) {
			got := runCLI(t, "", false, args...)

			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
			}

			for _, secret := range []string{fakeToken, "secret-password"} {
				if strings.Contains(got.stderr, secret) || strings.Contains(got.stdout, secret) {
					t.Errorf("output echoed %q:\nstdout %q\nstderr %q", secret, got.stdout, got.stderr)
				}
			}
		})
	}
}

// TestFaultsNeverEchoACredential pins one direction only: nothing secret comes out.
// These pin the other — which name does come out — one case per branch of
// declaredFlagName plus the help-topic replacement, so deleting a branch fails a test
// rather than only widening the silence.
func TestUndefinedFlagFaultsNameOnlyDeclaredFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		mentions string
		rejects  string
	}{
		{
			// trimArg cuts the argument at the space, leaving a whole declared name.
			name: "whole declared name", args: []string{"--dry-run " + fakeToken},
			mentions: "not defined: -dry-run",
		},
		{
			// The declared name survives as the first character of the attached form.
			name:     "attached value behind a declared letter",
			args:     []string{"generate-token", "-u", "admin", "-p" + fakeToken},
			mentions: "not defined: -p",
		},
		{
			// A flag-spelt remainder reads as a mistyped long name, so no name is safe.
			name:     "flag-shaped remainder withholds the name",
			args:     []string{"generate-token", "-u", "admin", "-psecret-password"},
			mentions: "not defined", rejects: "not defined: -",
		},
		{
			name: "undeclared letter names nothing", args: []string{"show", "ap-tag", "-z" + fakeToken},
			mentions: "not defined", rejects: "not defined: -",
		},
		{
			// urfave's own ExitCoder for an unknown help topic carries the word, so
			// failureLine replaces its text wholesale.
			name: "unknown help topic is replaced", args: []string{"--help", fakeToken},
			mentions: "unknown help topic", rejects: "No help topic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runCLI(t, "", false, tt.args...)

			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}

			if tt.rejects != "" && strings.Contains(got.stderr, tt.rejects) {
				t.Errorf("stderr %q must not mention %q", got.stderr, tt.rejects)
			}
		})
	}
}

func TestSettingsFaults(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		mentions string
	}{
		{name: "no controller", args: []string{"show", "ap-tag"}, mentions: "no controller given"},
		{
			name:     "zero timeout",
			args:     []string{"show", "ap-tag", "-c", "h", "--access-token", fakeToken, "-t", "0"},
			mentions: "must be positive",
		},
		{
			name:     "unknown format",
			args:     []string{"show", "ap-tag", "-c", "h", "--access-token", fakeToken, "-f", "yaml"},
			mentions: "accepted values",
		},
		{
			name:     "unknown sort key",
			args:     []string{"show", "ap-tag", "-c", "h", "--access-token", fakeToken, "-b", "bogus"},
			mentions: "accepted keys",
		},
		{
			name:     "unknown sort order",
			args:     []string{"show", "ap-tag", "-c", "h", "--access-token", fakeToken, "--sort-order", "sideways"},
			mentions: "accepted values",
		},
		{
			name:     "unknown log level",
			args:     []string{"--log-level", "loud", "show", "ap-tag", "-c", "h", "--access-token", fakeToken},
			mentions: "accepted values",
		},
		{
			name:     "unknown band",
			args:     []string{"show", "client", "-c", "h", "--access-token", fakeToken, "-r", "60"},
			mentions: "accepted values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runCLI(t, "", false, tt.args...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}
		})
	}
}

func TestConfigFile(t *testing.T) {
	t.Run("a named file that is not there is a fault", func(t *testing.T) {
		got := runCLI(t, "", false, "--config", filepath.Join(t.TempDir(), "absent.json"), "show", "ap-tag")
		if got.code != ExitUsage {
			t.Errorf("exit = %d, want %d", got.code, ExitUsage)
		}
	})

	t.Run("dry-run reports a valid file", func(t *testing.T) {
		path := writeFile(t, `{
		  "note": "lab",
		  "timeout": "30s",
		  "token": "`+fakeToken+`",
		  "controllers": [{"name": "lab", "host": "192.168.0.231"}]
		}`)

		got := runCLI(t, "", false, "--config", path, "--dry-run")
		if got.code != ExitOK {
			t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
		}

		if !strings.Contains(got.stdout, "1 controller") {
			t.Errorf("stdout = %q", got.stdout)
		}
	})

	t.Run("dry-run rejects an unknown key", func(t *testing.T) {
		got := runCLI(t, "", false, "--config", writeFile(t, `{"timeuot":"30s"}`), "--dry-run")
		if got.code != ExitUsage {
			t.Errorf("exit = %d, want %d", got.code, ExitUsage)
		}

		if !strings.Contains(got.stderr, "/timeuot") {
			t.Errorf("stderr %q does not carry the JSON pointer", got.stderr)
		}
	})

	t.Run("dry-run rejects a bad controller entry", func(t *testing.T) {
		got := runCLI(t, "", false,
			"--config", writeFile(t, `{"token":"`+fakeToken+`","controllers":[{"host":"https://h"}]}`), "--dry-run")

		if got.code != ExitUsage {
			t.Errorf("exit = %d, want %d", got.code, ExitUsage)
		}

		if strings.Contains(got.stderr, fakeToken) {
			t.Errorf("stderr echoed the token: %q", got.stderr)
		}
	})

	t.Run("dry-run contacts nothing", func(t *testing.T) {
		// 240.0.0.1 is unroutable, so a run that tried to reach it would hang or fail
		// rather than return at once.
		path := writeFile(t, `{"token":"`+fakeToken+`","controllers":[{"host":"240.0.0.1"}]}`)

		got := runCLI(t, "", false, "--config", path, "--dry-run")
		if got.code != ExitOK {
			t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
		}
	})

	t.Run("a flag beats the file", func(t *testing.T) {
		path := writeFile(t, `{"format":"json","token":"`+fakeToken+`","controllers":[{"host":"h"}]}`)

		// The file's format is valid and the flag's is not, so the fault itself is
		// the proof the flag won; the rejected value is withheld from the message.
		got := runCLI(t, "", false, "--config", path, "show", "ap-tag", "-f", "yaml")
		if got.code != ExitUsage || !strings.Contains(got.stderr, "--format") {
			t.Errorf("the flag did not reach validation: exit %d, stderr %q", got.code, got.stderr)
		}
	})
}

// The hook is consulted on the running command alone, never on a parent, so a node
// without one lets urfave print the whole help text and exit through its own path.
func TestEveryCommandHasAUsageHook(t *testing.T) {
	root := newRootCommand(Streams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})

	var walk func(*cli.Command, string)

	count := 0

	walk = func(cmd *cli.Command, path string) {
		count++

		if cmd.OnUsageError == nil {
			t.Errorf("command %q has no OnUsageError", path)
		}

		for _, sub := range cmd.Commands {
			walk(sub, path+" "+sub.Name)
		}
	}

	walk(root, root.Name)

	// root, deauth, delete and three, disable and two, enable and two, generate-token, reset
	// and two, save-config, set and three, show and nine.
	if want := 31; count != want {
		t.Errorf("walked %d commands, want %d", count, want)
	}
}

func writeFile(t *testing.T, body string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	return path
}

// A listing of a command's own vocabulary needs no controller and no token, so it must
// not fail on their absence and must not reach the network.
func TestSortKeysNeedsNoController(t *testing.T) {
	for _, tt := range []struct {
		noun string
		keys []string
	}{
		{"overview", show.OverviewKeys()},
		{"ap", show.APKeys()},
		{"ap-join", show.APJoinKeys()},
		{"ap-tag", show.APTagKeys()},
		{"client", show.ClientKeys()},
		{"wlan", show.WLANKeys()},
		{"policy-tag", show.PolicyTagKeys()},
		{"site-tag", show.SiteTagKeys()},
		{"rf-tag", show.RFTagKeys()},
	} {
		t.Run(tt.noun, func(t *testing.T) {
			got := runCLI(t, "", false, "show", tt.noun, "--sort-keys")

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if want := strings.Join(tt.keys, "\n") + "\n"; got.stdout != want {
				t.Errorf("stdout = %q, want %q", got.stdout, want)
			}

			if got.stderr != "" {
				t.Errorf("stderr = %q, want empty", got.stderr)
			}
		})
	}
}

// The listing answers before the settings are resolved, so a value Resolve would have
// rejected does not turn a listing into a fault.
func TestSortKeysOutranksARejectedValue(t *testing.T) {
	for _, args := range [][]string{
		{"show", "ap", "--sort-keys", "-b", "bogus"},
		{"show", "ap", "--sort-keys", "-t", "0"},
		{"show", "ap", "--sort-keys", "-c", "192.0.2.1"},
	} {
		t.Run(strings.Join(args[2:], " "), func(t *testing.T) {
			got := runCLI(t, "", false, args...)

			if got.code != ExitOK {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if got.stdout == "" {
				t.Error("stdout is empty")
			}
		})
	}
}

// Both flags are declared per leaf because urfave's GLOBAL OPTIONS section lists the
// root's flags only: declared on the show parent they would work and be invisible.
func TestEveryShowSubcommandCarriesTheSortFlags(t *testing.T) {
	root := newRootCommand(Streams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})

	var group *cli.Command

	for _, c := range root.Commands {
		if c.Name == "show" {
			group = c
		}
	}

	if group == nil {
		t.Fatal("the show group is missing")
	}

	for _, leaf := range group.Commands {
		names := map[string]bool{}
		for _, f := range leaf.Flags {
			for _, n := range f.Names() {
				names[n] = true
			}
		}

		for _, want := range []string{"sort-by", "sort-keys"} {
			if !names[want] {
				t.Errorf("show %s does not declare --%s", leaf.Name, want)
			}
		}
	}
}

// The whole point of the listing is that no OPTIONS line has to enumerate the keys.
func TestShowHelpLinesStayNarrow(t *testing.T) {
	const limit = 90

	for _, noun := range []string{
		"overview", "ap", "ap-join", "ap-tag", "client", "wlan",
		"policy-tag", "site-tag", "rf-tag",
	} {
		t.Run(noun, func(t *testing.T) {
			got := runCLI(t, "", false, "show", noun, "--help")

			inOptions := false
			for _, line := range strings.Split(got.stdout, "\n") {
				switch {
				case strings.HasPrefix(line, "OPTIONS:"):
					inOptions = true
				case strings.HasPrefix(line, "GLOBAL OPTIONS:"):
					inOptions = false
				case inOptions && len(line) > limit:
					t.Errorf("%d columns exceeds %d: %q", len(line), limit, line)
				}
			}
		})
	}
}

// -o is the conventional output-format shorthand, so --format owns it here and --sort-order is
// long-only. Each fault names the flag that consumed the value and withholds the value, which is
// what makes the swap visible in both directions.
func TestFormatOwnsTheOutputShorthand(t *testing.T) {
	for _, tt := range []struct {
		name     string
		args     []string
		mentions string
	}{
		{name: "-o names the format", args: []string{"-o", "yaml"}, mentions: "--format: accepted values"},
		{name: "-f still names the format", args: []string{"-f", "yaml"}, mentions: "--format: accepted values"},
		{name: "-o no longer names the order", args: []string{"-o", "desc"}, mentions: "--format: accepted values"},
		{
			name:     "--sort-order kept its own vocabulary",
			args:     []string{"--sort-order", "sideways"},
			mentions: "--sort-order: accepted values",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			args := append([]string{"show", "ap-tag", "-c", "h", "--access-token", fakeToken}, tt.args...)

			got := runCLI(t, "", false, args...)
			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if !strings.Contains(got.stderr, tt.mentions) {
				t.Errorf("stderr %q does not mention %q", got.stderr, tt.mentions)
			}
		})
	}
}

// The other half of the swap: -o json has to pass resolution rather than merely parse,
// so the run is taken as far as the read and must fail on the host and not on the value.
func TestOutputShorthandReachesTheRead(t *testing.T) {
	for _, short := range []string{"-o", "-f"} {
		t.Run(short, func(t *testing.T) {
			got := runCLI(t, "", false, "show", "ap-tag",
				"-c", "192.0.2.99", "--access-token", fakeToken, "-t", "1ms", short, "json")

			if got.code != ExitFailure {
				t.Errorf("exit = %d, want %d (stderr %q)", got.code, ExitFailure, got.stderr)
			}
		})
	}
}

// urfave appends the completion subtree during its own setup pass, after newRootCommand has
// returned, so TestEveryCommandHasAUsageHook cannot see those five nodes;
// ConfigureShellCompletionCommand is what carries the hook to them.
func TestCompletionSubtreeCarriesTheUsageHook(t *testing.T) {
	for _, args := range [][]string{
		{"completion", "--bogus"},
		{"completion", "bash", "--bogus"},
		{"completion", "zsh", "--bogus"},
		{"completion", "fish", "--bogus"},
		{"completion", "pwsh", "--bogus"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			got := runCLI(t, "", false, args...)

			if got.code != ExitUsage {
				t.Errorf("exit = %d, want %d", got.code, ExitUsage)
			}

			if got.stdout != "" {
				t.Errorf("stdout = %q, want empty", got.stdout)
			}

			want := "see 'wnc " + strings.Join(args[:len(args)-1], " ") + " --help'"
			if !strings.Contains(got.stderr, want) {
				t.Errorf("stderr %q does not name %q", got.stderr, want)
			}

			if lines := strings.Count(strings.TrimRight(got.stderr, "\n"), "\n"); lines != 0 {
				t.Errorf("stderr has %d extra lines: %q", lines, got.stderr)
			}
		})
	}
}

// Each shell writes its script to stdout and nothing else, and the command urfave
// generates is hidden, so enabling it adds no line to the root's own listing.
func TestCompletionWritesAScriptAndStaysHidden(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish", "pwsh"} {
		t.Run(shell, func(t *testing.T) {
			got := runCLI(t, "", false, "completion", shell)

			if got.code != ExitOK {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
			}

			if got.stdout == "" {
				t.Error("no script on stdout")
			}

			if got.stderr != "" {
				t.Errorf("stderr = %q, want empty", got.stderr)
			}
		})
	}

	if root := runCLI(t, "", false, "--help"); strings.Contains(root.stdout, "completion") {
		t.Errorf("the root listing names the hidden command: %q", root.stdout)
	}
}

// The flag the generated scripts call back with. The reply is one token per line on
// stdout, and it names flags without ever repeating the value one was given.
func TestGenerateShellCompletionAnswersOnStdout(t *testing.T) {
	got := runCLI(t, "", false, "show", "--generate-shell-completion")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
	}

	if !strings.Contains(got.stdout, "ap-tag:") {
		t.Errorf("stdout = %q, want the show leaves", got.stdout)
	}

	got = runCLI(t, "", false, "show", "ap", "--access-token", fakeToken, "--generate-shell-completion")
	if strings.Contains(got.stdout, fakeToken) || strings.Contains(got.stderr, fakeToken) {
		t.Errorf("completion echoed the token:\nstdout %q\nstderr %q", got.stdout, got.stderr)
	}
}

// leafPaths is every childless node's argv path, walked from the tree rather than listed, so
// a leaf added later is covered without an edit here.
func leafPaths(t *testing.T) [][]string {
	t.Helper()

	root := newRootCommand(Streams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})

	var (
		paths [][]string
		walk  func(*cli.Command, []string)
	)

	walk = func(cmd *cli.Command, path []string) {
		if len(cmd.Commands) == 0 {
			paths = append(paths, path)

			return
		}

		for _, sub := range cmd.Commands {
			walk(sub, append(slices.Clone(path), sub.Name))
		}
	}

	walk(root, nil)

	return paths
}

// Every value a subcommand takes is named by a flag, so a leftover word is a fault wherever
// it lands. The word is a credential here on purpose: a leftover on generate-token can be the
// password a wrapper misplaced, and the message reports the count rather than the word.
func TestEveryLeafRefusesAPositional(t *testing.T) {
	paths := leafPaths(t)

	// deauth, three delete, two disable, two enable, generate-token, two reset, save-config,
	// three set, nine show.
	if want := 24; len(paths) != want {
		t.Fatalf("walked %d leaves, want %d", len(paths), want)
	}

	for _, path := range paths {
		if slices.Contains(path, "completion") {
			t.Fatalf("the walk reached urfave's completion subtree: %v", path)
		}

		t.Run(strings.Join(path, "/"), func(t *testing.T) {
			got := runCLI(t, "", false, append(slices.Clone(path), fakeToken)...)

			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitUsage, got.stderr)
			}

			if got.stdout != "" {
				t.Errorf("a refused run wrote to stdout: %q", got.stdout)
			}

			for _, want := range []string{
				"takes no positional arguments, 1 given",
				"see 'wnc " + strings.Join(path, " ") + " --help'",
			} {
				if !strings.Contains(got.stderr, want) {
					t.Errorf("stderr %q does not carry %q", got.stderr, want)
				}
			}

			for _, secret := range []string{fakeToken, "secret-password"} {
				if strings.Contains(got.stderr, secret) {
					t.Errorf("stderr %q echoes the word it refused", got.stderr)
				}
			}

			if n := len(strings.Split(strings.TrimRight(got.stderr, "\n"), "\n")); n != 1 {
				t.Errorf("stderr holds %d lines, want 1: %q", n, got.stderr)
			}
		})
	}
}

// The hint is derived from the leaf's own flags, so it names the target where there is one
// and stays silent where there is not. generate-token is the case that matters: nothing may
// invite an operator to put a password on a flag.
func TestArgumentRefusalHintsTheTargetFlag(t *testing.T) {
	for _, tt := range []struct {
		path []string
		want string
	}{
		{path: []string{"reset", "ap"}, want: ": use --ap-name"},
		{path: []string{"reset", "capwap"}, want: ": use --ap-name"},
		{path: []string{"enable", "radio"}, want: ": use --ap-name"},
		{path: []string{"disable", "ap"}, want: ": use --ap-name"},
		{path: []string{"show", "client"}, want: ": use --ap-name"},
		{path: []string{"set", "rf-tag"}, want: ": use --name"},
		{path: []string{"delete", "policy-tag"}, want: ": use --name"},
		{path: []string{"show", "ap"}, want: ""},
		{path: []string{"show", "wlan"}, want: ""},
		{path: []string{"generate-token"}, want: ""},
	} {
		t.Run(strings.Join(tt.path, "/"), func(t *testing.T) {
			got := runCLI(t, "", false, append(slices.Clone(tt.path), "stray")...)

			switch {
			case tt.want == "" && strings.Contains(got.stderr, ": use "):
				t.Errorf("stderr %q hints a flag this leaf does not take", got.stderr)
			case tt.want != "" && !strings.Contains(got.stderr, tt.want):
				t.Errorf("stderr %q does not hint %q", got.stderr, tt.want)
			}
		})
	}
}

// attachLeafRules skips a leaf with no action, because urfave installs its own help action
// for one during setup and the closure would call nil. Such a leaf would run unguarded, so
// the tree is held to declaring one everywhere.
func TestEveryLeafDeclaresAnAction(t *testing.T) {
	root := newRootCommand(Streams{Out: &bytes.Buffer{}, Err: &bytes.Buffer{}})

	var walk func(*cli.Command, string)

	count := 0

	walk = func(cmd *cli.Command, path string) {
		if len(cmd.Commands) == 0 {
			count++

			if cmd.Action == nil {
				t.Errorf("leaf %q declares no action", path)
			}

			return
		}

		for _, sub := range cmd.Commands {
			walk(sub, path+" "+sub.Name)
		}
	}

	walk(root, root.Name)

	if want := 24; count != want {
		t.Errorf("walked %d leaves, want %d", count, want)
	}
}

// The completion subtree is urfave's own and answers for its own arguments. It is appended during
// the library's setup pass, after every walk in newRootCommand has run, so the refusal cannot
// reach it and must not: that subtree took the shell name as a positional argument on earlier
// releases.
func TestCompletionSubtreeKeepsItsOwnArguments(t *testing.T) {
	got := runCLI(t, "", false, "completion", "bash", "extra")

	if got.code != ExitOK {
		t.Fatalf("exit = %d, want %d (stderr %q)", got.code, ExitOK, got.stderr)
	}

	if !strings.Contains(got.stdout, "complete") {
		t.Errorf("stdout carries no completion script: %q", got.stdout)
	}

	if got.stderr != "" {
		t.Errorf("stderr = %q, want empty", got.stderr)
	}
}

// The same bound TestShowHelpLinesStayNarrow holds the show leaves to. --ap-name widens the names
// column on the six leaves that declare it, which takes --slot's own line over the bound unless
// its default text is hidden, and generate-token is out because it already prints wider.
//
// save-config and deauth are out because, being flat, each declares the connection flags itself,
// so what a nested acting command shows under INHERITED OPTIONS lands in its OPTIONS instead and
// this scan reads only one of the two headings.
func TestExecHelpLinesStayNarrow(t *testing.T) {
	const limit = 90

	flat := []string{"show", "generate-token", "save-config", "deauth"}

	for _, path := range leafPaths(t) {
		if slices.Contains(flat, path[0]) {
			continue
		}

		t.Run(strings.Join(path, "/"), func(t *testing.T) {
			got := runCLI(t, "", false, append(slices.Clone(path), "--help")...)

			inOptions := false
			for _, line := range strings.Split(got.stdout, "\n") {
				switch {
				case strings.HasPrefix(line, "OPTIONS:"):
					inOptions = true
				case strings.HasPrefix(line, "GLOBAL OPTIONS:"):
					inOptions = false
				case inOptions && len(line) > limit:
					t.Errorf("%d columns exceeds %d: %q", len(line), limit, line)
				}
			}
		})
	}
}

// Deleting ArgsUsage left urfave's USAGE line saying only "[options]", so every acting leaf
// stages the fragment naming what it refuses to run without and attachLeafRules prefixes the
// path. The path is the part a leaf cannot build for itself, and getting it wrong is silent
// in the help rather than in a run.
func TestActingLeavesNameTheirTargetInTheSynopsis(t *testing.T) {
	for _, tt := range []struct {
		path []string
		want string
	}{
		{path: []string{"reset", "ap"}, want: "wnc reset ap --ap-name <ap-name> [options]"},
		{path: []string{"reset", "capwap"}, want: "wnc reset capwap --ap-name <ap-name> [options]"},
		{path: []string{"enable", "ap"}, want: "wnc enable ap --ap-name <ap-name> [options]"},
		{
			path: []string{"disable", "radio"},
			want: "wnc disable radio --ap-name <ap-name> --slot <n> [options]",
		},
		{path: []string{"set", "rf-tag"}, want: "wnc set rf-tag --name <name> [options]"},
		{path: []string{"delete", "site-tag"}, want: "wnc delete site-tag --name <name> [options]"},
		// The one leaf whose mandatory flags are a choice rather than a conjunction. Without
		// this row the second arm can vanish from the line an operator reads first, silently.
		{
			path: []string{"deauth"},
			want: "wnc deauth (--mac <mac> | --username <username>) [options]",
		},
		// A show leaf takes no mandatory value, so it keeps urfave's own line.
		{path: []string{"show", "ap"}, want: "wnc show ap [options]"},
	} {
		t.Run(strings.Join(tt.path, "/"), func(t *testing.T) {
			got := runCLI(t, "", false, append(slices.Clone(tt.path), "--help")...)

			if !strings.Contains(got.stdout, "USAGE:\n   "+tt.want+"\n") {
				t.Errorf("help does not carry USAGE %q:\n%s", tt.want, got.stdout)
			}
		})
	}
}

// The counterpart to TestFaultsNeverEchoACredential: a word spelt like a command name is
// still repeated, and still carries the near miss. Without this the redaction could swallow
// every word and the credential test would go on passing.
func TestUnknownCommandStillNamesACommandShapedWord(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want string
	}{
		{args: []string{"bogus"}, want: `unknown command "bogus"`},
		{args: []string{"shwo"}, want: `unknown command "shwo" (did you mean show?)`},
		{args: []string{"show", "ap-tags"}, want: `unknown command "ap-tags" (did you mean ap-tag?)`},
		{args: []string{"show", "BOGUS"}, want: "unknown command:"},
	} {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			got := runCLI(t, "", false, tt.args...)

			if got.code != ExitUsage {
				t.Fatalf("exit = %d, want %d", got.code, ExitUsage)
			}

			if !strings.Contains(got.stderr, tt.want) {
				t.Errorf("stderr %q does not carry %q", got.stderr, tt.want)
			}
		})
	}
}
