package config

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/urfave/cli/v3"
)

// resolveWith runs a command carrying the real show flags and returns what Resolve
// made of them. Driving the actual flag definitions is the point: a test that built a
// Settings by hand would not notice a flag renamed on one side only.
func resolveWith(t *testing.T, file File, args ...string) (Settings, error) {
	t.Helper()

	// resolveHosts reads EnvController itself rather than through a flag source, so a
	// developer's own shell would otherwise supply the controllers these tests assert
	// on. Clearing it needs t.Setenv, which is why nothing reaching this helper is parallel.
	t.Setenv(EnvController, "")

	var (
		got     Settings
		gotErr  error
		sortKey = "ap_name"
	)

	cmd := &cli.Command{
		Name: "show",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: FlagController, Aliases: []string{"c"}},
			&cli.StringFlag{Name: FlagAccessToken},
			&cli.BoolFlag{Name: FlagInsecure, Aliases: []string{"k"}},
			&cli.StringFlag{Name: FlagFormat, Aliases: []string{"f"}, Value: DefaultFormat},
			&cli.DurationFlag{Name: FlagTimeout, Aliases: []string{"t"}, Value: DefaultTimeout},
			&cli.StringFlag{Name: FlagSortOrder, Aliases: []string{"o"}, Value: DefaultSortOrder},
			&cli.StringFlag{Name: FlagSortBy, Aliases: []string{"b"}, Value: sortKey},
			&cli.StringFlag{Name: FlagRadio, Aliases: []string{"r"}},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			got, gotErr = Resolve(c, file, []string{sortKey, "mac"}, sortKey)

			return nil
		},
	}

	if err := cmd.Run(t.Context(), append([]string{"show"}, args...)); err != nil {
		t.Fatalf("running the command: %v", err)
	}

	return got, gotErr
}

// resolveExecWith drives ResolveExec through a command declaring only the flags the
// exec tree does, which is the point of the assertion: the output and sort flags are
// absent, so anything ResolveExec read from them would fail here.
func resolveExecWith(t *testing.T, file File, args ...string) (Settings, error) {
	t.Helper()

	t.Setenv(EnvController, "")

	var (
		got    Settings
		gotErr error
	)

	cmd := &cli.Command{
		Name: "reset",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: FlagController, Aliases: []string{"c"}},
			&cli.StringFlag{Name: FlagAccessToken},
			&cli.BoolFlag{Name: FlagInsecure, Aliases: []string{"k"}},
			&cli.DurationFlag{Name: FlagTimeout, Aliases: []string{"t"}, Value: DefaultTimeout},
			&cli.BoolFlag{Name: FlagYes},
		},
		Action: func(_ context.Context, c *cli.Command) error {
			got, gotErr = ResolveExec(c, file)

			return nil
		},
	}

	if err := cmd.Run(t.Context(), append([]string{"reset"}, args...)); err != nil {
		t.Fatalf("running the command: %v", err)
	}

	return got, gotErr
}

func TestResolveExec(t *testing.T) {
	t.Run("the flag supplies the controller and the timeout", func(t *testing.T) {
		got, err := resolveExecWith(t, File{},
			"-c", "192.0.2.1", "--access-token", fakeToken, "-t", "5s", "-k")
		if err != nil {
			t.Fatalf("ResolveExec: %v", err)
		}

		if len(got.Controllers) != 1 || got.Controllers[0].Host != "192.0.2.1" {
			t.Errorf("controllers = %+v", got.Controllers)
		}

		if got.Timeout != 5*time.Second || !got.Insecure {
			t.Errorf("timeout = %s, insecure = %v", got.Timeout, got.Insecure)
		}
	})

	t.Run("the file supplies what no flag set", func(t *testing.T) {
		fileTimeout := Dur(15 * time.Second)
		insecure := true

		got, err := resolveExecWith(t, File{
			Timeout:     &fileTimeout,
			Insecure:    &insecure,
			Token:       strp(fakeToken),
			Controllers: []Controller{{Host: strp("192.0.2.10")}},
		})
		if err != nil {
			t.Fatalf("ResolveExec: %v", err)
		}

		if got.Timeout != 15*time.Second || !got.Insecure {
			t.Errorf("timeout = %s, insecure = %v", got.Timeout, got.Insecure)
		}

		if len(got.Controllers) != 1 || got.Controllers[0].Host != "192.0.2.10" {
			t.Errorf("controllers = %+v", got.Controllers)
		}
	})

	t.Run("a controller list is refused when nothing names one", func(t *testing.T) {
		if _, err := resolveExecWith(t, File{}); err == nil {
			t.Error("no controller was accepted")
		}
	})

	t.Run("a non-positive timeout is refused", func(t *testing.T) {
		_, err := resolveExecWith(t, File{}, "-c", "192.0.2.1", "--access-token", fakeToken, "-t", "0s")
		if err == nil {
			t.Error("a zero timeout was accepted")
		}
	})
}

func TestResolvePrecedence(t *testing.T) {
	table := "table"
	fileTimeout := Dur(15 * time.Second)
	insecure := true

	file := File{
		Format:      &table,
		Timeout:     &fileTimeout,
		Insecure:    &insecure,
		Token:       strp(fakeToken),
		Controllers: []Controller{{Host: strp("192.0.2.10")}},
	}

	t.Run("the file supplies what no flag set", func(t *testing.T) {
		got, err := resolveWith(t, file)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if len(got.Controllers) != 1 || got.Controllers[0].Host != "192.0.2.10" {
			t.Errorf("controllers = %+v", got.Controllers)
		}

		if got.Timeout != 15*time.Second || got.Format != table || !got.Insecure {
			t.Errorf("settings = %+v", got)
		}
	})

	t.Run("a flag beats the file", func(t *testing.T) {
		got, err := resolveWith(t, file, "-c", "192.0.2.20", "--"+FlagAccessToken, fakeToken,
			"-t", "5s", "-f", "json")
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if got.Controllers[0].Host != "192.0.2.20" {
			t.Errorf("controllers = %+v", got.Controllers)
		}

		if got.Timeout != 5*time.Second || got.Format != FormatJSON {
			t.Errorf("settings = %+v", got)
		}
	})

	t.Run("the flag default applies where neither speaks", func(t *testing.T) {
		got, err := resolveWith(t, File{Token: file.Token, Controllers: file.Controllers})
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if got.Timeout != DefaultTimeout || got.Format != DefaultFormat || got.SortOrder != DefaultSortOrder {
			t.Errorf("settings = %+v", got)
		}

		if got.Insecure {
			t.Error("insecure defaulted to true")
		}
	})
}

// The flag is repeated rather than comma separated, and one token covers every host it
// names: the flag selects, the 0600 file supplies the credential.
func TestResolveRepeatedControllerAndFileToken(t *testing.T) {
	file := File{
		Token: strp(fakeToken),
		Controllers: []Controller{
			{Host: strp("192.0.2.10")},
			{Host: strp("192.0.2.11")},
		},
	}

	t.Run("every host named is queried, in the order given", func(t *testing.T) {
		got, err := resolveWith(t, file, "-c", "192.0.2.11", "-c", "192.0.2.10",
			"--"+FlagAccessToken, fakeToken)
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if len(got.Controllers) != 2 ||
			got.Controllers[0].Host != "192.0.2.11" || got.Controllers[1].Host != "192.0.2.10" {
			t.Fatalf("controllers = %+v", got.Controllers)
		}
	})

	t.Run("a host with no flag token takes the file's", func(t *testing.T) {
		got, err := resolveWith(t, file, "-c", "192.0.2.11")
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if len(got.Controllers) != 1 || got.Controllers[0].Token != fakeToken {
			t.Fatalf("controllers = %+v", got.Controllers)
		}
	})

	// A host the file does not list is still queried: the file supplies the credential and does
	// not decide which hosts a flag may name.
	t.Run("a host the file does not name takes the file's token too", func(t *testing.T) {
		got, err := resolveWith(t, file, "-c", "192.0.2.99")
		if err != nil {
			t.Fatalf("Resolve: %v", err)
		}

		if len(got.Controllers) != 1 || got.Controllers[0].Token != fakeToken {
			t.Fatalf("controllers = %+v", got.Controllers)
		}
	})

	t.Run("no token anywhere is a fault", func(t *testing.T) {
		_, err := resolveWith(t, File{}, "-c", "192.0.2.99")
		if err == nil {
			t.Fatal("Resolve accepted a controller with no token")
		}

		if !strings.Contains(err.Error(), "no token given") {
			t.Errorf("error %q does not say a token is missing", err)
		}

		if strings.Contains(err.Error(), fakeToken) {
			t.Errorf("error echoed a token: %q", err)
		}
	})
}

func TestResolveRejections(t *testing.T) {
	ctrl := []Controller{{Host: strp("192.0.2.10")}}

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no controller anywhere", args: nil, want: "no controller given"},
		{name: "zero timeout", args: []string{"-t", "0"}, want: "must be positive"},
		{name: "negative timeout", args: []string{"-t", "-1s"}, want: "must be positive"},
		{name: "unknown format", args: []string{"-f", "yaml"}, want: "accepted values"},
		{name: "unknown order", args: []string{"-o", "sideways"}, want: "accepted values"},
		{name: "unknown sort key", args: []string{"-b", "nope"}, want: "accepted keys"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := File{}
			if tt.name != "no controller anywhere" {
				file.Token, file.Controllers = strp(fakeToken), ctrl
			}

			_, err := resolveWith(t, file, tt.args...)
			if err == nil {
				t.Fatalf("Resolve accepted %v", tt.args)
			}

			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error %q does not mention %q", err, tt.want)
			}
		})
	}
}

// A timeout in the file must be rejected the same way a flag would be, or a bad file
// value reaches the transport as an already-expired deadline.
func TestResolveRejectsANonPositiveFileTimeout(t *testing.T) {
	zero := Dur(0)
	file := File{
		Timeout:     &zero,
		Token:       strp(fakeToken),
		Controllers: []Controller{{Host: strp("192.0.2.10")}},
	}

	if _, err := resolveWith(t, file); err == nil {
		t.Error("a zero timeout in the file was accepted")
	}
}

func TestResolveSortDirection(t *testing.T) {
	file := File{Token: strp(fakeToken), Controllers: []Controller{{Host: strp("192.0.2.10")}}}

	asc, err := resolveWith(t, file)
	if err != nil || asc.Descending() {
		t.Errorf("default order = %q, descending = %v (err %v)", asc.SortOrder, asc.Descending(), err)
	}

	desc, err := resolveWith(t, file, "-o", OrderDesc)
	if err != nil || !desc.Descending() {
		t.Errorf("explicit desc = %q, descending = %v (err %v)", desc.SortOrder, desc.Descending(), err)
	}
}

func TestResolveRadioFilter(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "unset selects every band", args: nil, want: ""},
		{name: "2.4", args: []string{"-r", Band24}, want: Band24},
		{name: "5", args: []string{"-r", Band5}, want: Band5},
		{name: "6", args: []string{"-r", Band6}, want: Band6},
		{name: "a band that does not exist", args: []string{"-r", "60"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got    string
				gotErr error
			)

			cmd := &cli.Command{
				Name:  "show",
				Flags: []cli.Flag{&cli.StringFlag{Name: FlagRadio, Aliases: []string{"r"}}},
				Action: func(_ context.Context, c *cli.Command) error {
					got, gotErr = ResolveRadio(c)

					return nil
				},
			}

			if err := cmd.Run(t.Context(), append([]string{"show"}, tt.args...)); err != nil {
				t.Fatalf("running the command: %v", err)
			}

			if tt.wantErr {
				if gotErr == nil {
					t.Error("an unknown band was accepted")
				}

				return
			}

			if gotErr != nil {
				t.Fatalf("ResolveRadio: %v", gotErr)
			}

			if got != tt.want {
				t.Errorf("band = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveLogLevel(t *testing.T) {
	levels := []string{"error", "warning", "debug"}
	debug := "debug"

	tests := []struct {
		name    string
		args    []string
		file    File
		want    string
		wantErr bool
	}{
		{name: "the default", want: DefaultLogLevel},
		{name: "from the file", file: File{LogLevel: &debug}, want: "debug"},
		{
			name: "the flag beats the file",
			args: []string{"--log-level", "error"},
			file: File{LogLevel: &debug},
			want: "error",
		},
		{name: "a level that does not exist", args: []string{"--log-level", "loud"}, wantErr: true},
		{name: "a level the file invented", file: File{LogLevel: strp("loud")}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got    string
				gotErr error
			)

			cmd := &cli.Command{
				Name: "wnc",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: FlagLogLevel, Value: DefaultLogLevel},
				},
				Action: func(_ context.Context, c *cli.Command) error {
					got, gotErr = ResolveLogLevel(c, tt.file, levels)

					return nil
				},
			}

			if err := cmd.Run(t.Context(), append([]string{"wnc"}, tt.args...)); err != nil {
				t.Fatalf("running the command: %v", err)
			}

			if tt.wantErr {
				if gotErr == nil {
					t.Error("an unknown level was accepted")
				}

				return
			}

			if gotErr != nil || got != tt.want {
				t.Errorf("level = %q (err %v), want %q", got, gotErr, tt.want)
			}
		})
	}
}

func TestValidateFile(t *testing.T) {
	levels := []string{"error", "warning", "debug"}
	bad := "yaml"
	badLevel := "loud"
	zero := Dur(0)

	tests := []struct {
		name    string
		file    File
		want    string
		wantErr bool
	}{
		{
			name: "a full file",
			file: File{
				Format:      strp(FormatJSON),
				LogLevel:    strp("debug"),
				Token:       strp(fakeToken),
				Controllers: []Controller{{Host: strp("192.0.2.10")}},
			},
		},
		{name: "an empty file", file: File{}},
		{name: "unknown format", file: File{Format: &bad}, want: "format", wantErr: true},
		{name: "unknown log level", file: File{LogLevel: &badLevel}, want: "log_level", wantErr: true},
		{name: "zero timeout", file: File{Timeout: &zero}, want: "timeout", wantErr: true},
		{
			// One fault for the whole list rather than one per entry: the token is the
			// file's, so --dry-run reports it once however many hosts are listed.
			name:    "controllers with no token",
			file:    File{Controllers: []Controller{{Host: strp("192.0.2.10")}}},
			want:    "no token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateFile(tt.file, levels)

			if tt.wantErr {
				if err == nil {
					t.Fatal("ValidateFile accepted the file")
				}

				if !strings.Contains(err.Error(), tt.want) {
					t.Errorf("error %q does not mention %q", err, tt.want)
				}

				return
			}

			if err != nil {
				t.Errorf("ValidateFile: %v", err)
			}
		})
	}
}

// The file's enum slots sit two keys from its token, so a fault on either may name
// only the key and the accepted set — the rejected value can be that token.
func TestValidateFileNeverEchoesAValue(t *testing.T) {
	t.Parallel()

	host := "192.0.2.10"

	for name, file := range map[string]File{
		"format":    {Format: strp(fakeToken), Token: strp(fakeToken), Controllers: []Controller{{Host: &host}}},
		"log_level": {LogLevel: strp(fakeToken), Token: strp(fakeToken), Controllers: []Controller{{Host: &host}}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			_, err := ValidateFile(file, []string{"error", "warning", "debug"})
			if err == nil {
				t.Fatal("ValidateFile accepted the file")
			}

			if strings.Contains(err.Error(), fakeToken) {
				t.Errorf("error echoed the value: %q", err)
			}
		})
	}
}

func strp(s string) *string { return &s }
