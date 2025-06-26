package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// RegisterApTagSubCommand registers a subcommand for listing access point tags.
func RegisterApTagSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "ap-tag",
			Usage:     "Show the access point tags",
			UsageText: "wnc show ap-tag [options...]",
			Aliases:   []string{"t"},
			Flags:     registerApTagCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewShowCli(&c, &r, &u)

				c.SetShowCmdConfig(cmd)
				f.InvokeApTagCli().ShowApTag()
				return nil
			},
		},
	}
}

// registerApTagCmdFlags returns flags for the ap-tag command.
func registerApTagCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerControllersFlag()...)
	flags = append(flags, registerInsecureFlag()...)
	flags = append(flags, registerPrintFormatFlag()...)
	flags = append(flags, registerTimeoutFlag()...)
	return flags
}
