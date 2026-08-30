// Package log builds the process logger and the slog bridge the WNC SDK takes.
package log

import (
	"fmt"
	"io"
	"log/slog"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/sirupsen/logrus"
	logrusslog "github.com/sirupsen/logrus/hooks/slog"
)

// New builds the process logger. Diagnostics go to stderr so they never mix with
// the data on stdout, and color is disabled unconditionally because a TTY check
// would make the output differ between a terminal and a pipe.
func New(level string) (*logrus.Logger, error) {
	return NewWithOutput(os.Stderr, level)
}

func NewWithOutput(w io.Writer, level string) (*logrus.Logger, error) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %q: %w", level, err)
	}

	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(lv)
	l.SetFormatter(formatterFor(lv))

	return l, nil
}

// formatterFor picks the rendering from the verbosity asked for. logfmt is the debug form because
// a bug report wants every field separately; below that the fields become a sentence.
func formatterFor(lv logrus.Level) logrus.Formatter {
	if lv >= logrus.DebugLevel {
		return &logrus.TextFormatter{DisableTimestamp: true, DisableColors: true}
	}

	return plainFormatter{}
}

// The three fields this package handles by name; every other field internal/show/reporter.go
// attaches is rendered by the generic path in tailFields.
const (
	fieldController = "controller"
	fieldCause      = "cause"
	fieldStatus     = "status"
)

// plainFormatter renders one diagnostic as a sentence led by the controller it
// happened on, with the cause token kept in its field spelling so what an operator
// reads is what an operator searches docs/TROUBLESHOOTING.md for.
type plainFormatter struct{}

func (plainFormatter) Format(e *logrus.Entry) ([]byte, error) {
	var b strings.Builder

	b.WriteString(e.Level.String())
	b.WriteString(": ")

	if v, ok := e.Data[fieldController]; ok {
		fmt.Fprintf(&b, "%v: ", v)
	}

	b.WriteString(e.Message)

	if tail := tailFields(e.Data); tail != "" {
		b.WriteString(" (" + tail + ")")
	}

	b.WriteByte('\n')

	return []byte(b.String()), nil
}

// tailFields renders the fields that follow the sentence. cause leads because that spelling is the
// heading docs/TROUBLESHOOTING.md is indexed by, an unknown field is appended rather than dropped,
// and the two that never appear are controller, which led the line, and status, which wnc.Message
// has already rebuilt the text from.
func tailFields(data logrus.Fields) string {
	pairs := make([]string, 0, len(data))

	if v, ok := data[fieldCause]; ok {
		pairs = append(pairs, fmt.Sprintf("%s=%v", fieldCause, v))
	}

	for _, name := range slices.Sorted(maps.Keys(data)) {
		if name == fieldController || name == fieldCause || name == fieldStatus {
			continue
		}

		pairs = append(pairs, fmt.Sprintf("%s=%v", name, data[name]))
	}

	return strings.Join(pairs, ", ")
}

// SDKLogger adapts the logger to the *slog.Logger the SDK accepts. Every SDK record is mapped to
// Debug because internal/show reports each controller failure itself, and Enabled consults the
// same mapper, so the SDK is silent at warning and above.
func SDKLogger(l *logrus.Logger) *slog.Logger {
	return slog.New(logrusslog.NewHandler(l, &logrusslog.HandlerOptions{
		LevelMapper: func(slog.Level) logrus.Level { return logrus.DebugLevel },
	}))
}

// Levels lists the accepted --log-level values, for the flag's usage string and for the error a
// rejected value produces. Three of logrus's seven: nothing emits at Info, Trace, Fatal or Panic,
// and the spellings come from logrus, where WarnLevel is "warning".
func Levels() []string {
	return []string{
		logrus.ErrorLevel.String(),
		logrus.WarnLevel.String(),
		logrus.DebugLevel.String(),
	}
}
