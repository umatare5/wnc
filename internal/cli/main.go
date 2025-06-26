package cli

import (
	"context"
	"log"
	"os"

	generateCmd "github.com/umatare5/wnc/internal/cli/generate"
	showCmd "github.com/umatare5/wnc/internal/cli/show"
	"github.com/urfave/cli/v3"
)

// Run initializes and starts the CLI application.
func Run() {
	cmd := &cli.Command{
		Name:      "wnc",
		Usage:     "Client for Cisco C9800 Wireless Network Controller API",
		UsageText: "wnc [command] [options...]",
		Version:   getVersion(),
		Commands:  registerSubCommands(),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			_ = cli.ShowAppHelp(cmd)
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// registerSubCommands registers the commands for the CLI application.
func registerSubCommands() []*cli.Command {
	cmds := []*cli.Command{}
	cmds = append(cmds, generateCmd.RegisterGenerateCommand()...)
	cmds = append(cmds, showCmd.RegisterShowCommand()...)
	return cmds
}
