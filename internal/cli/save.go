package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

// saveConfigCommand persists the controller's running configuration. It is a flat tree rather than
// a leaf of reset, which admits an RPC only if the RPC persists nothing: this one exists to
// persist, and its schema declares the output container that test refuses.
func saveConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "save-config",
		Usage: "Save the running configuration to the startup configuration",
		Description: "The startup configuration is the only destination: no file may be named, and\n" +
			"every change on the controller is persisted rather than only what this CLI\n" +
			"wrote. An access point's admin state is unaffected, being no part of the\n" +
			"configuration. The save took two to four seconds on every release measured, so\n" +
			"a --timeout a read survives can still refuse it. Pass --dry-run to name the\n" +
			"controller and change nothing.",
		Flags:  append(execFlags(), yesFlag()),
		Action: runSaveConfig,
	}
}

// runSaveConfig saves one controller's configuration. Nothing is resolved first: the RPC names
// no target, so the guard sequence is the shared prologue and then the prompt, and there is no
// read whose absence could tell the operator they named the wrong controller.
func runSaveConfig(ctx context.Context, cmd *cli.Command) error {
	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	logger := runtimeFrom(ctx).Logger

	return confirmAndAct(ctx, cmd, writeWording{
		question: "Save the running configuration of %s? Every change on the controller is " +
			"persisted, including changes this CLI did not make. [y/N]: ",
		dryRun: "%s: would save the running configuration\n",
		sent:   "%s: running configuration saved\n",
		failed: "saving the running configuration on %s",
	}, func(ctx context.Context) error {
		result, err := client.SaveConfig(ctx)
		if err != nil {
			return err
		}

		// The controller's own account, kept for a bug report rather than printed: it adds
		// nothing the reported line does not.
		logger.WithField("controller", target.Name).Debug(result)

		return nil
	}, target.Name)
}
