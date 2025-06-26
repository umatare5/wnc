package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// RegisterOverviewSubCommand registers a subcommand for showing overview of the wireless infrastructure.
func RegisterOverviewSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "overview",
			Usage:     "Show overview of the wireless infrastructure",
			UsageText: "wnc show overview [options...]",
			Aliases:   []string{"o"},
			Flags:     registerOverviewCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewShowCli(&c, &r, &u)

				c.SetShowCmdConfig(cmd)
				f.InvokeOverviewCli().ShowOverview()
				return nil
			},
		},
	}
}

// registerOverviewCmdFlags returns flags for the overview command.
func registerOverviewCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerControllersFlag()...)
	flags = append(flags, registerInsecureFlag()...)
	flags = append(flags, registerPrintFormatFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	flags = append(flags, registerRadioFlag()...)
	flags = append(flags, registerOverviewSortByFlag()...)
	flags = append(flags, registerSortOrderFlag()...)
	return flags
}
