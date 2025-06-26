package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// RegisterWlanSubCommand registers a subcommand for listing WLANs.
func RegisterWlanSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "wlan",
			Usage:     "Show the ESSIDs in the wireless infrastructure",
			UsageText: "wnc show wlan [options...]",
			Aliases:   []string{"w"},
			Flags:     registerWlanCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewShowCli(&c, &r, &u)

				c.SetShowCmdConfig(cmd)
				f.InvokeWlanCli().ShowWlan()
				return nil
			},
		},
	}
}

// registerWlanCmdFlags returns flags for the wlan command.
func registerWlanCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerControllersFlag()...)
	flags = append(flags, registerInsecureFlag()...)
	flags = append(flags, registerPrintFormatFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	return flags
}
