package show

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// The fan-out is exercised with a fetcher that never reaches the network. Client
// construction still happens for real, so the hosts are from the RFC 5737 test range:
// they are valid authorities the SDK accepts and nothing here connects to them.
const (
	hostA = "192.0.2.1"
	hostB = "192.0.2.2"
)

type fanRow struct {
	Name       *string `json:"name,omitzero"`
	Controller string  `json:"controller"`
}

func fanColumns() []render.Column[fanRow] {
	return []render.Column[fanRow]{
		{Key: "name", Header: "Name", Cell: func(r fanRow) string { return render.StrPtr(r.Name) }},
		{Key: keyController, Header: headController, Cell: func(r fanRow) string { return render.Str(r.Controller) }},
	}
}

// fanEnv builds an Env over the given controllers and returns it with the buffers it
// writes to.
func fanEnv(t *testing.T, format string, names ...string) (env Env, out, errOut *bytes.Buffer) {
	t.Helper()

	var stdout, stderr bytes.Buffer

	logger, err := log.NewWithOutput(&stderr, "info")
	if err != nil {
		t.Fatalf("building the logger: %v", err)
	}

	targets := make([]config.Target, 0, len(names))
	for i, name := range names {
		host := hostA
		if i%2 == 1 {
			host = hostB
		}

		targets = append(targets, config.Target{Name: name, Host: host, Token: fanToken})
	}

	return Env{
		Settings: config.Settings{
			Controllers: targets,
			Timeout:     time.Second,
			Insecure:    true,
			Format:      format,
			SortBy:      "name",
			SortOrder:   config.OrderAsc,
		},
		Logger:    logger,
		Out:       &stdout,
		UserAgent: "wnc/test",
	}, &stdout, &stderr
}

const fanToken = "TestToken0123456789ABCDEF=="

// rowsPerController answers with one row per controller, or an error for the names
// listed as failing.
func rowsPerController(failing, degraded map[string]error) Fetcher[fanRow] {
	return func(_ context.Context, _ *wnc.Client, t config.Target, rep *Reporter) ([]fanRow, error) {
		if err, ok := failing[t.Name]; ok {
			return nil, err
		}

		if err, ok := degraded[t.Name]; ok {
			rep.Degraded("secondary", err)
		}

		name := t.Name

		return []fanRow{{Name: &name, Controller: t.Name}}, nil
	}
}

func TestRunMergesEveryController(t *testing.T) {
	t.Parallel()

	env, out, errOut := fanEnv(t, config.FormatTable, "wlc-b", "wlc-a")

	if err := Run(t.Context(), env, fanColumns(), rowsPerController(nil, nil)); err != nil {
		t.Fatalf("Run: %v", err)
	}

	lines := strings.Split(strings.TrimRight(out.String(), "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("output = %q, want a heading and two rows", out.String())
	}

	// The rows are ordered by the sort key, not by the order the controllers answered.
	if !strings.HasPrefix(lines[1], "wlc-a") || !strings.HasPrefix(lines[2], "wlc-b") {
		t.Errorf("rows are not sorted: %q", out.String())
	}

	if errOut.Len() != 0 {
		t.Errorf("a clean run wrote to stderr: %q", errOut.String())
	}
}

func TestRunPartialFailure(t *testing.T) {
	t.Parallel()

	env, out, errOut := fanEnv(t, config.FormatTable, "wlc-good", "wlc-bad")

	err := Run(t.Context(), env, fanColumns(),
		rowsPerController(map[string]error{"wlc-bad": errors.New("boom")}, nil))

	if !errors.Is(err, ErrPartial) {
		t.Fatalf("Run = %v, want ErrPartial", err)
	}

	// The half of the fleet that answered is still printed: withholding it would hide
	// the healthy controllers.
	if !strings.Contains(out.String(), "wlc-good") {
		t.Errorf("the surviving rows were dropped: %q", out.String())
	}

	if strings.Contains(out.String(), "wlc-bad") {
		t.Errorf("a failed controller produced a row: %q", out.String())
	}

	// The failure is reported once: a sentence led by the controller, with the cause kept
	// in the field spelling the troubleshooting index is built on.
	if !strings.Contains(errOut.String(), "error: wlc-bad:") || !strings.Contains(errOut.String(), "cause=") {
		t.Errorf("stderr = %q", errOut.String())
	}
}

// A read that costs only some cells still makes the run a partial one, because the
// operator's view of those cells is not what the controller holds.
func TestRunDegradedReadIsPartial(t *testing.T) {
	t.Parallel()

	env, out, errOut := fanEnv(t, config.FormatTable, "wlc-a")

	err := Run(t.Context(), env, fanColumns(),
		rowsPerController(nil, map[string]error{"wlc-a": errors.New("secondary down")}))

	if !errors.Is(err, ErrPartial) {
		t.Fatalf("Run = %v, want ErrPartial", err)
	}

	if !strings.Contains(out.String(), "wlc-a") {
		t.Errorf("the row was dropped: %q", out.String())
	}

	if !strings.Contains(errOut.String(), "endpoint=secondary") {
		t.Errorf("stderr = %q", errOut.String())
	}
}

// With no controller answering there is nothing to print: a heading on its own says
// the fleet is empty, which is the opposite of what happened.
func TestRunTotalFailurePrintsNothing(t *testing.T) {
	t.Parallel()

	env, out, _ := fanEnv(t, config.FormatTable, "wlc-a", "wlc-b")

	err := Run(t.Context(), env, fanColumns(), rowsPerController(map[string]error{
		"wlc-a": errors.New("boom"),
		"wlc-b": errors.New("boom"),
	}, nil))

	if !errors.Is(err, ErrAllFailed) {
		t.Fatalf("Run = %v, want ErrAllFailed", err)
	}

	if out.Len() != 0 {
		t.Errorf("stdout = %q, want nothing", out.String())
	}
}

// An empty fleet is a normal answer: the heading prints and the run succeeds.
func TestRunEmptyFleet(t *testing.T) {
	t.Parallel()

	env, out, _ := fanEnv(t, config.FormatTable, "wlc-a")

	empty := func(context.Context, *wnc.Client, config.Target, *Reporter) ([]fanRow, error) {
		return nil, nil
	}

	if err := Run(t.Context(), env, fanColumns(), empty); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got := strings.TrimSpace(out.String()); got != "Name  Controller" {
		t.Errorf("output = %q, want the heading alone", got)
	}
}

func TestRunJSONFormat(t *testing.T) {
	t.Parallel()

	env, out, _ := fanEnv(t, config.FormatJSON, "wlc-a")

	if err := Run(t.Context(), env, fanColumns(), rowsPerController(nil, nil)); err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := `[{"name":"wlc-a","controller":"wlc-a"}]` + "\n"
	if out.String() != want {
		t.Errorf("output = %q, want %q", out.String(), want)
	}
}

// An interrupted run prints nothing: half a table read as a whole one is worse than no
// table at all.
func TestRunInterruptedPrintsNothing(t *testing.T) {
	t.Parallel()

	env, out, _ := fanEnv(t, config.FormatTable, "wlc-a")

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	err := Run(ctx, env, fanColumns(), rowsPerController(nil, nil))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run = %v, want context.Canceled", err)
	}

	if out.Len() != 0 {
		t.Errorf("stdout = %q, want nothing", out.String())
	}
}

// A sort key no column declares is a fault of the caller and is reported rather than
// producing an arbitrarily ordered table.
func TestRunRejectsAnUnknownSortKey(t *testing.T) {
	t.Parallel()

	env, _, _ := fanEnv(t, config.FormatTable, "wlc-a")
	env.Settings.SortBy = "nope"

	err := Run(t.Context(), env, fanColumns(), rowsPerController(nil, nil))
	if err == nil || !strings.Contains(err.Error(), "no such column") {
		t.Errorf("Run = %v, want a rejected sort key", err)
	}
}

// A host the SDK refuses costs that controller its rows and nothing else.
func TestRunReportsAClientThatCannotBeBuilt(t *testing.T) {
	t.Parallel()

	env, out, errOut := fanEnv(t, config.FormatTable, "wlc-good", "wlc-bad")
	env.Settings.Controllers[1].Token = ""

	err := Run(t.Context(), env, fanColumns(), rowsPerController(nil, nil))
	if !errors.Is(err, ErrPartial) {
		t.Fatalf("Run = %v, want ErrPartial", err)
	}

	if !strings.Contains(out.String(), "wlc-good") {
		t.Errorf("the healthy controller lost its row: %q", out.String())
	}

	if !strings.Contains(errOut.String(), "error: wlc-bad:") {
		t.Errorf("stderr = %q", errOut.String())
	}
}

// Reporting waits until every controller has finished and then walks them in the order
// they were given, so two runs of one command produce the same lines.
func TestRunReportsInControllerOrder(t *testing.T) {
	t.Parallel()

	first := captureFailures(t, "wlc-1", "wlc-2", "wlc-3")
	second := captureFailures(t, "wlc-1", "wlc-2", "wlc-3")

	if first != second {
		t.Errorf("two runs reported differently:\n%s\n%s", first, second)
	}

	if !strings.Contains(first, "wlc-1") || strings.Index(first, "wlc-1") > strings.Index(first, "wlc-3") {
		t.Errorf("failures are not in controller order: %s", first)
	}
}

func captureFailures(t *testing.T, names ...string) string {
	t.Helper()

	env, _, errOut := fanEnv(t, config.FormatTable, names...)

	failing := make(map[string]error, len(names))
	for _, n := range names {
		failing[n] = errors.New("boom")
	}

	_ = Run(t.Context(), env, fanColumns(), rowsPerController(failing, nil))

	return errOut.String()
}

// The SDK reports at Error and Debug only, and this layer reports every controller
// failure itself, so an info-level run must not carry a second copy from the SDK.
func TestSDKLoggerIsSilentAtInfo(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	logger, err := log.NewWithOutput(&buf, "info")
	if err != nil {
		t.Fatalf("building the logger: %v", err)
	}

	log.SDKLogger(logger).Error("GET failed")

	if buf.Len() != 0 {
		t.Errorf("the SDK logger wrote at info: %q", buf.String())
	}

	if logger.GetLevel() != logrus.InfoLevel {
		t.Errorf("level = %v", logger.GetLevel())
	}
}
