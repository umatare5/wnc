package subcommand

import (
	"context"

	"github.com/urfave/cli/v3"
)

// RegisterGenerateCommand registers the main generate command.
func RegisterGenerateCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "generate",
			Usage:     "Generate values and files for the WNC",
			UsageText: "wnc generate [subcommand] [options...]",
			Aliases:   []string{"g"},
			Commands:  registerGenerateSubCommands(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				_ = cli.ShowSubcommandHelp(cmd)
				return nil
			},
		},
	}
}

// registerGenerateSubCommands returns subcommands for the generate command.
func registerGenerateSubCommands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, registerTokenSubCommand()...)
	return cmds
}
