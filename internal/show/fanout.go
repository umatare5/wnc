package show

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/wnc"
)

// Env is what a show command needs from the CLI layer.
type Env struct {
	Settings  config.Settings
	Logger    *logrus.Logger
	Out       io.Writer
	UserAgent string
}

// Fetcher reads one controller and builds its rows. The error it returns is fatal for that
// controller, so its rows are dropped; a read costing only some cells goes to the Reporter and
// the fetch continues.
type Fetcher[R any] func(ctx context.Context, c *wnc.Client, t config.Target, rep *Reporter) ([]R, error)

type result[R any] struct {
	rows     []R
	reporter *Reporter
	fatal    error
}

// Run executes one show command end to end. Every controller is read concurrently while the reads
// inside one are sequential, and reporting waits for all of them and then walks them in the order
// they were given, so two runs of the same command produce the same lines.
func Run[R any](ctx context.Context, env Env, cols []render.Column[R], fetch Fetcher[R]) error {
	results := make([]result[R], len(env.Settings.Controllers))

	var wg sync.WaitGroup

	for i, target := range env.Settings.Controllers {
		wg.Go(func() {
			results[i] = readOne(ctx, env, target, fetch)
		})
	}

	wg.Wait()

	// An interrupted run prints nothing: half a table read as a whole one is worse
	// than no table at all.
	if err := ctx.Err(); err != nil {
		return err
	}

	rows, failed, degraded := collect(env.Logger, results)

	out := outcome(len(results), failed, degraded)

	// With no controller answering there is nothing to print. A heading on its own
	// says the fleet is empty, which is the opposite of what happened.
	if errors.Is(out, ErrAllFailed) {
		return out
	}

	if err := renderRows(env, cols, rows); err != nil {
		return err
	}

	return out
}

func readOne[R any](
	ctx context.Context, env Env, target config.Target, fetch Fetcher[R],
) result[R] {
	rep := &Reporter{target: target}

	client, err := wnc.NewClient(target, env.Settings, env.Logger, env.UserAgent)
	if err != nil {
		return result[R]{reporter: rep, fatal: err}
	}

	rows, err := fetch(ctx, client, target, rep)
	if err != nil {
		return result[R]{reporter: rep, fatal: err}
	}

	return result[R]{rows: rows, reporter: rep}
}

// collect merges the rows and writes every controller's faults to the log, in the
// order the controllers were given.
func collect[R any](logger *logrus.Logger, results []result[R]) (rows []R, failed, degraded int) {
	for _, r := range results {
		if r.fatal != nil {
			failed++

			r.reporter.logFatal(logger, r.fatal)
		}

		if r.reporter.logDegraded(logger) {
			degraded++
		}

		rows = append(rows, r.rows...)
	}

	return rows, failed, degraded
}

func renderRows[R any](env Env, cols []render.Column[R], rows []R) error {
	if err := render.Sort(rows, cols, env.Settings.SortBy, env.Settings.Descending()); err != nil {
		return fmt.Errorf("--%s: %w", config.FlagSortBy, err)
	}

	if env.Settings.Format == config.FormatJSON {
		return render.JSON(env.Out, rows)
	}

	if env.Settings.Pretty {
		return render.PrettyTable(env.Out, cols, rows)
	}

	return render.Table(env.Out, cols, rows)
}

// outcome classifies the run.
func outcome(total, failed, degraded int) error {
	switch {
	case total > 0 && failed == total:
		return ErrAllFailed
	case failed > 0 || degraded > 0:
		return ErrPartial
	default:
		return nil
	}
}
