package subcommand

import (
	"context"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework"
	"github.com/umatare5/wnc/internal/infrastructure"

	"github.com/urfave/cli/v3"
)

// registerTokenSubCommand registers a subcommand for generating tokens.
func registerTokenSubCommand() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "token",
			Usage:     "Generate a basic authentication token to connect to the WNC",
			UsageText: "wnc generate token [options...]",
			Aliases:   []string{"t"},
			Flags:     registerTokenSubCmdFlags(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				c := config.New()
				r := infrastructure.New(&c)
				u := application.New(&c, &r)
				f := framework.NewGenerateCli(&c, &r, &u)

				c.SetGenerateCmdConfig(cmd)
				f.InvokeTokenCli().GenerateToken()
				return nil
			},
		},
	}
}

// registerTokenSubCmdFlags returns flags for the token command.
func registerTokenSubCmdFlags() []cli.Flag {
	flags := []cli.Flag{}
	flags = append(flags, registerUsernameFlag()...)
	flags = append(flags, registerPasswordFlag()...)
	return flags
}
