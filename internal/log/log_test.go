package log

import (
	"bytes"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewWithOutputLevels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		level   string
		wantErr bool
		want    logrus.Level
	}{
		{name: "info", level: "info", want: logrus.InfoLevel},
		{name: "debug", level: "debug", want: logrus.DebugLevel},
		{name: "warning", level: "warning", want: logrus.WarnLevel},
		{name: "error", level: "error", want: logrus.ErrorLevel},
		{name: "trace", level: "trace", want: logrus.TraceLevel},
		{name: "rejects unknown", level: "verbose", wantErr: true},
		{name: "rejects empty", level: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			l, err := NewWithOutput(&buf, tt.level)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewWithOutput(%q) = nil error, want one", tt.level)
				}

				return
			}

			if err != nil {
				t.Fatalf("NewWithOutput(%q) = %v", tt.level, err)
			}

			if l.GetLevel() != tt.want {
				t.Errorf("level = %v, want %v", l.GetLevel(), tt.want)
			}
		})
	}
}

// The debug form is logfmt, and it must stay reproducible: a timestamp or a color
// escape would make two runs of one command differ.
func TestDebugFormatterIsLogfmtAndDeterministic(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	l, err := NewWithOutput(&buf, "debug")
	if err != nil {
		t.Fatalf("NewWithOutput: %v", err)
	}

	l.WithField(fieldController, "WNC1").Error("read failed")
	got := buf.String()

	if strings.Contains(got, "\x1b[") {
		t.Errorf("output holds an ANSI escape: %q", got)
	}

	if strings.Contains(got, "time=") {
		t.Errorf("output holds a timestamp: %q", got)
	}

	if !strings.Contains(got, "controller=WNC1") {
		t.Errorf("field missing from %q", got)
	}
}

// One diagnostic renders as one sentence rather than as logfmt. The three middle
// field sets are internal/show/reporter.go's; the first is the fieldless TLS
// advisory internal/cli emits and the last a field nothing attaches.
func TestPlainFormatterRendersOneSentence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		fields logrus.Fields
		msg    string
		want   string
	}{
		{
			name: "no field at all",
			msg:  "TLS certificate verification is disabled",
			want: "error: TLS certificate verification is disabled\n",
		},
		{
			name:   "a note carries the controller only",
			fields: logrus.Fields{fieldController: "WNC1"},
			msg:    "2 rows excluded: no band reported",
			want:   "error: WNC1: 2 rows excluded: no band reported\n",
		},
		{
			// The status is left out on purpose: wnc.Message rebuilds an *APIError's
			// text from the status, so the sentence already carries it.
			name:   "a fatal read names the controller and the cause",
			fields: logrus.Fields{fieldController: "WNC1", fieldCause: "auth", fieldStatus: 401},
			msg:    "the controller answered 401 Unauthorized",
			want:   "error: WNC1: the controller answered 401 Unauthorized (cause=auth)\n",
		},
		{
			name: "a degraded read names the endpoint too",
			fields: logrus.Fields{
				fieldController: "WNC1", "endpoint": "oper-data", fieldCause: "not-found", fieldStatus: 404,
			},
			msg:  "the controller answered 404 Not Found",
			want: "error: WNC1: the controller answered 404 Not Found (cause=not-found, endpoint=oper-data)\n",
		},
		{
			// A field this rendering does not know is appended rather than dropped.
			name:   "an unknown field still reaches the operator",
			fields: logrus.Fields{fieldController: "WNC1", "attempt": 2},
			msg:    "retrying",
			want:   "error: WNC1: retrying (attempt=2)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			l, err := NewWithOutput(&buf, "info")
			if err != nil {
				t.Fatalf("NewWithOutput: %v", err)
			}

			l.WithFields(tt.fields).Error(tt.msg)

			if got := buf.String(); got != tt.want {
				t.Errorf("output = %q, want %q", got, tt.want)
			}
		})
	}
}

// The rendering is chosen by the verbosity that was asked for, so the boundary is
// pinned in both directions rather than at one level.
func TestFormatterFollowsTheLevel(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		level     string
		wantPlain bool
	}{
		{level: "error", wantPlain: true},
		{level: "warning", wantPlain: true},
		{level: "info", wantPlain: true},
		{level: "debug", wantPlain: false},
		{level: "trace", wantPlain: false},
	} {
		t.Run(tt.level, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			l, err := NewWithOutput(&buf, tt.level)
			if err != nil {
				t.Fatalf("NewWithOutput: %v", err)
			}

			l.WithField(fieldController, "WNC1").Error("boom")

			got := buf.String()
			plain := !strings.Contains(got, "level=error")

			if plain != tt.wantPlain {
				t.Errorf("--log-level %s produced %q; plain = %v, want %v", tt.level, got, plain, tt.wantPlain)
			}
		})
	}
}

func TestSDKLoggerMapsEverythingToDebug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level slog.Level
	}{
		{name: "error", level: slog.LevelError},
		{name: "warn", level: slog.LevelWarn},
		{name: "info", level: slog.LevelInfo},
		{name: "debug", level: slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			l, err := NewWithOutput(&buf, "debug")
			if err != nil {
				t.Fatalf("NewWithOutput: %v", err)
			}

			SDKLogger(l).Log(t.Context(), tt.level, "sdk said something")

			if got := buf.String(); !strings.Contains(got, "level=debug") {
				t.Errorf("record at %v did not land at debug: %q", tt.level, got)
			}
		})
	}
}

// The SDK reports every failure the fan-out also reports, so a run at the default
// level must not carry a second copy. Enabled consults the same mapper, which is what
// makes even an Error record silent here.
func TestSDKLoggerSilentAtTheDefaultLevel(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	l, err := NewWithOutput(&buf, "warning")
	if err != nil {
		t.Fatalf("NewWithOutput: %v", err)
	}

	sdk := SDKLogger(l)
	sdk.Error("GET failed")
	sdk.Info("connected")

	if got := buf.String(); got != "" {
		t.Errorf("a default-level run is not silent: %q", got)
	}
}

func TestLevelsOffersOnlyTheReachableThree(t *testing.T) {
	t.Parallel()

	got := Levels()
	if want := []string{"error", "warning", "debug"}; !slices.Equal(got, want) {
		t.Fatalf("Levels() = %v, want %v", got, want)
	}

	for _, name := range got {
		if _, err := logrus.ParseLevel(name); err != nil {
			t.Errorf("Levels() offers %q, which logrus rejects: %v", name, err)
		}
	}

	// The four this CLI never emits at stay out. panic and fatal answer a failed run
	// with a non-zero exit and an empty stderr, which defeats the failure report the
	// fan-out exists to produce; info and trace cannot be told from warning and debug.
	for _, name := range []string{"panic", "fatal", "info", "trace"} {
		if slices.Contains(got, name) {
			t.Errorf("Levels() offers %q, which no code path emits at", name)
		}
	}
}
