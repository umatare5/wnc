package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// RegisterApSubCommand registers a subcommand for listing access points.
func RegisterApSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "ap",
			Usage:     "Show the access points",
			UsageText: "wnc show ap [options...]",
			Aliases:   []string{"a"},
			Flags:     registerApCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewShowCli(&c, &r, &u)

				c.SetShowCmdConfig(cmd)
				f.InvokeApCli().ShowAp()
				return nil
			},
		},
	}
}

// registerApCmdFlags returns flags for the ap command.
func registerApCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerControllersFlag()...)
	flags = append(flags, registerInsecureFlag()...)
	flags = append(flags, registerPrintFormatFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	return flags
}
