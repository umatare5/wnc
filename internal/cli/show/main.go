package subcommand

import (
	"context"

	"github.com/urfave/cli/v3"
)

// RegisterShowCommand registers the main show command.
func RegisterShowCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "show",
			Usage:     "Show information about the wireless infrastructure",
			UsageText: "wnc show [subcommand] [options...]",
			Aliases:   []string{"s"},
			Commands:  registerShowSubCommands(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				_ = cli.ShowSubcommandHelp(cmd)
				return nil
			},
		},
	}
}

// registerShowSubCommands returns subcommands for the show command.
func registerShowSubCommands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, RegisterApSubCommand()...)
	cmds = append(cmds, RegisterApTagSubCommand()...)
	cmds = append(cmds, RegisterClientSubCommand()...)
	cmds = append(cmds, RegisterOverviewSubCommand()...)
	cmds = append(cmds, RegisterWlanSubCommand()...)
	return cmds
}
