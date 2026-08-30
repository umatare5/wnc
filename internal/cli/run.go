package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v3"
)

// Run executes the CLI and returns the process exit code, and owns the signal context because main
// exits through os.Exit.
func Run(args []string) int {
	piped, err := stdinIsPiped()
	if err != nil {
		piped = false
	}

	return RunWith(args, Streams{Out: os.Stdout, Err: os.Stderr, In: os.Stdin, InPipe: piped})
}

// RunWith is Run against caller-supplied streams, which is how the tests drive the
// real command tree.
func RunWith(args []string, st Streams) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	root := newRootCommand(st)
	err := root.Run(ctx, args)
	code := exitCode(ctx, err)

	if err != nil && code != ExitSignal && !reported(err) {
		// This is the last thing written; a failure here has nowhere left to go.
		_, _ = fmt.Fprintf(st.Err, "wnc: %s\n", failureLine(err))
	}

	return code
}

// failureLine is the one line a failure prints. Nothing in this package builds a cli.ExitCoder, so
// one can only be urfave's own, and its "No help topic" text repeats the word verbatim past every
// redaction in usage.go — "-h <token>" is one keystroke from "-h", so a fixed line replaces it.
func failureLine(err error) string {
	var coder cli.ExitCoder
	if errors.As(err, &coder) {
		return "invalid usage: unknown help topic: see 'wnc --help'"
	}

	return err.Error()
}

// stdinIsPiped reports whether stdin is anything other than a terminal, which is
// what decides whether generate-token may read a password from it.
func stdinIsPiped() (bool, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}

	return info.Mode()&os.ModeCharDevice == 0, nil
}
