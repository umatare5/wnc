package cli

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// resetCommand groups the actions that restart something. Every leaf names one target, resolves
// it against the controller before acting, and asks first unless --yes was given.
func resetCommand() *cli.Command {
	return &cli.Command{
		Name:     "reset",
		Usage:    "Restart an access point or its controller session",
		Flags:    execFlags(),
		Action:   parentAction,
		Commands: []*cli.Command{resetAPCommand(), resetCAPWAPCommand()},
	}
}

// resetAPCommand names its target with --ap-name. Both RPCs in this tree declare a mandatory choice
// offering an ap-name or a mac-addr arm, so the name an operator reads out of show ap is what goes
// on the wire and no address is spelt anywhere.
func resetAPCommand() *cli.Command {
	return &cli.Command{
		Name:      leafAP,
		Usage:     "Restart one access point",
		UsageText: synopsis(config.FlagAPName),
		Description: "--ap-name is the ap_name column of wnc show ap. The controller resolves it\n" +
			"first, so a name it holds no access point under is refused before the RPC.\n" +
			"Pass --dry-run to name the target and change nothing.",
		Flags:  []cli.Flag{apNameFlag(), yesFlag()},
		Action: runResetAP,
	}
}

// resetCAPWAPCommand restarts the controller session rather than the access point. The
// access point's own uptime is unchanged across it, which is what makes this the smaller
// remedy of the two and why it is a separate leaf rather than a flag on reset ap.
func resetCAPWAPCommand() *cli.Command {
	return &cli.Command{
		Name:      "capwap",
		Usage:     "Reset one access point's controller session",
		UsageText: synopsis(config.FlagAPName),
		Description: "--ap-name is the ap_name column of wnc show ap. The controller resolves it\n" +
			"first, so a name it holds no access point under is refused before the RPC.\n" +
			"The access point does not reboot: only its CAPWAP session is re-established.\n" +
			"Pass --dry-run to name the target and change nothing.",
		Flags:  []cli.Flag{apNameFlag(), yesFlag()},
		Action: runResetCAPWAP,
	}
}

func runResetAP(ctx context.Context, cmd *cli.Command) error {
	return runAPAction(ctx, cmd, writeWording{
		question: "Reset %s on %s? Its clients disconnect for about four minutes. [y/N]: ",
		dryRun:   "%s on %s: would reset\n",
		sent:     "%s on %s: reset sent\n",
		failed:   "resetting %s on %s",
	}, func(ctx context.Context, c *wnc.Client, name string) error {
		return c.ResetAP(ctx, name)
	})
}

// runResetCAPWAP resets one access point's controller session.
func runResetCAPWAP(ctx context.Context, cmd *cli.Command) error {
	return runAPAction(ctx, cmd, writeWording{
		question: "Reset the CAPWAP session of %s on %s? " +
			"It rejoins within about ten seconds and does not reboot. [y/N]: ",
		dryRun: "%s on %s: would reset the CAPWAP session\n",
		sent:   "%s on %s: capwap reset sent\n",
		failed: "resetting the CAPWAP session of %s on %s",
	}, func(ctx context.Context, c *wnc.Client, name string) error {
		return c.ResetCAPWAP(ctx, name)
	})
}
