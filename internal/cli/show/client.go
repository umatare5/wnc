package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// RegisterClientSubCommand registers a subcommand for listing wireless clients.
func RegisterClientSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "client",
			Usage:     "Show the wireless clients",
			UsageText: "wnc show client [options...]",
			Aliases:   []string{"c"},
			Flags:     registerClientCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewShowCli(&c, &r, &u)

				c.SetShowCmdConfig(cmd)
				f.InvokeClientCli().ShowClient()
				return nil
			},
		},
	}
}

// registerClientCmdFlags returns flags for the client command.
func registerClientCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerControllersFlag()...)
	flags = append(flags, registerInsecureFlag()...)
	flags = append(flags, registerPrintFormatFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	flags = append(flags, registerRadioFlag()...)
	flags = append(flags, registerSSIDFlag()...)
	flags = append(flags, registerClientSortByFlag()...)
	flags = append(flags, registerSortOrderFlag()...)
	return flags
}
