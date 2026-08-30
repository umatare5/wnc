package cli

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/show"
)

// Exit codes. A run that reached every controller and rendered its rows exits 0
// even when there were no rows: an empty fleet is a normal answer.
const (
	ExitOK      = 0
	ExitFailure = 1
	ExitUsage   = 2
	ExitPartial = 3
	ExitSignal  = 130
)

// exitCode maps the run's outcome, and the order is load-bearing. A cli.ExitCoder comes first
// because urfave builds one itself and it matches none of this package's sentinels, and the signal
// check comes next so an interrupted fan-out is not reported as a partial success.
func exitCode(ctx context.Context, err error) int {
	var coder cli.ExitCoder

	switch {
	case errors.As(err, &coder):
		return ExitUsage
	case ctx.Err() != nil:
		return ExitSignal
	case errors.Is(err, ErrUsage):
		return ExitUsage
	case errors.Is(err, show.ErrPartial):
		return ExitPartial
	case err != nil:
		return ExitFailure
	default:
		return ExitOK
	}
}

// reported tells whether the fan-out has already written this failure to the log,
// which is the case for every outcome it classifies per controller. Printing a
// summary line on top of those would repeat what the operator just read.
func reported(err error) bool {
	return errors.Is(err, show.ErrPartial) || errors.Is(err, show.ErrAllFailed)
}
