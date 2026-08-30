// Package cli wires the command tree, owns the process exit code, and is the only
// place that decides what reaches stdout and what reaches stderr.
package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// Streams is injected rather than read from os so a test can drive the whole tree, the
// piped-stdin branch of generate-token included.
type Streams struct {
	Out    io.Writer
	Err    io.Writer
	In     io.Reader
	InPipe bool
}

// runtimeKey keys the per-run state urfave's Before installs on the context, which is what lets
// the root resolve the file and the logger once for every action.
type runtimeKey struct{}

type runtimeState struct {
	Logger   *logrus.Logger
	Streams  Streams
	File     config.File
	Path     string
	Explicit bool
}

// runtimeFrom reads the state back. The root's Before is the only place that installs it and it
// runs first, so every command finds it here.
func runtimeFrom(ctx context.Context) *runtimeState {
	st, ok := ctx.Value(runtimeKey{}).(*runtimeState)
	if !ok {
		return nil
	}

	return st
}

// newRootCommand builds the whole tree. ExitErrHandler displaces urfave's HandleExitCoder call,
// which reaches os.Exit, and HideHelpCommand removes the help subcommand and is inherited by every
// child; the three walks that follow are not inherited and so are applied per node.
//
// The completion subtree is urfave's own and is appended during the library's setup pass, after
// every walk below has run. ConfigureShellCompletionCommand is the seam that reaches it, set on the
// root alone because urfave's condition tests that field without also testing for the root, and it
// carries the usage hook and not attachLeafRules.
func newRootCommand(st Streams) *cli.Command {
	root := &cli.Command{
		Name:                            "wnc",
		Usage:                           "Operate Cisco Catalyst 9800 Wireless Network Controllers",
		Version:                         Version(),
		Writer:                          st.Out,
		ErrWriter:                       st.Err,
		Reader:                          st.In,
		HideHelpCommand:                 true,
		EnableShellCompletion:           true,
		ConfigureShellCompletionCommand: attachUsageHandler,
		ExitErrHandler:                  func(context.Context, *cli.Command, error) {},
		Flags:                           rootFlags(),
		Before:                          rootBefore(st),
		Action:                          rootAction,
		Commands: []*cli.Command{
			deauthCommand(),
			deleteCommand(),
			disableCommand(),
			enableCommand(),
			generateTokenCommand(),
			resetCommand(),
			saveConfigCommand(),
			setCommand(),
			showCommand(),
		},
	}

	attachUsageHandler(root)
	attachInheritedOptions(root)
	attachLeafRules(root, "")

	return root
}

// rootBefore reads the configuration file once and builds the logger. Doing it at
// the root rather than per command is what lets the file's log_level govern the
// reporting of the run that follows.
func rootBefore(streams Streams) cli.BeforeFunc {
	return func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
		path, explicit := config.Path(cmd.String(config.FlagConfig))

		// The logger does not exist yet, so an advisory from the read is buffered
		// and replayed once the level is known.
		var pending []string

		file, err := config.Load(path, explicit, func(m string) { pending = append(pending, m) })
		if err != nil {
			return ctx, fmt.Errorf("%w: %w", ErrUsage, err)
		}

		level, err := config.ResolveLogLevel(cmd, file, log.Levels())
		if err != nil {
			return ctx, fmt.Errorf("%w: %w", ErrUsage, err)
		}

		logger, err := log.NewWithOutput(streams.Err, level)
		if err != nil {
			return ctx, fmt.Errorf("%w: %w", ErrUsage, err)
		}

		for _, m := range pending {
			logger.Warn(m)
		}

		return context.WithValue(ctx, runtimeKey{}, &runtimeState{
			Logger:   logger,
			Streams:  streams,
			File:     file,
			Path:     path,
			Explicit: explicit,
		}), nil
	}
}

// rootAction handles the two ways the root is reached with no subcommand.
func rootAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.Bool(config.FlagDryRun) {
		return dryRun(ctx, cmd)
	}

	return parentAction(ctx, cmd)
}

// dryRun validates the configuration file and contacts nothing. The syntax, the
// unknown-key check and the duration decode already happened in Before, so what is
// left is the parts of a file that only a show command would otherwise reach.
func dryRun(ctx context.Context, cmd *cli.Command) error {
	st := runtimeFrom(ctx)

	// A default path that does not exist was skipped by Before. Reporting a file
	// that is not there as valid would be the wrong answer, so it is read again as
	// if it had been named.
	if !st.Explicit {
		if _, err := config.Load(st.Path, true, nil); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
	}

	targets, err := config.ValidateFile(st.File, log.Levels())
	if err != nil {
		return fmt.Errorf("%w: %s: %w", ErrUsage, st.Path, err)
	}

	_, err = fmt.Fprintf(cmd.Root().Writer, "%s: valid, %d controller(s)\n", st.Path, len(targets))

	return err
}
