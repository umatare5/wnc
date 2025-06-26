package subcommand

import (
	"github.com/umatare5/wnc/internal/config"
	"github.com/urfave/cli/v3"
)

// registerUsernameFlag returns the username flag for authentication.
func registerUsernameFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.UsernameFlagName,
			Usage:   "Username to generate Basic Authentication header",
			Aliases: []string{"u"},
			Value:   "",
		},
	}
}

// registerPasswordFlag returns the password flag for authentication.
func registerPasswordFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.PasswordFlagName,
			Usage:   "Password to generate Basic Authentication header",
			Aliases: []string{"p"},
			Value:   "",
		},
	}
}
