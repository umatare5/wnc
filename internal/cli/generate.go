package cli

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
)

// generateTokenCommand prints the RESTCONF Basic auth token, and is the one command that contacts
// no controller. It is flat rather than a group of one, which would share no flags with its single
// child.
func generateTokenCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate-token",
		Aliases: []string{"g"},
		Usage:   "Print the Basic auth token for a controller account",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    config.FlagUsername,
				Aliases: []string{"u"},
				Usage:   "controller username",
				Sources: cli.EnvVars(config.EnvUsername),
			},
			&cli.StringFlag{
				Name:    config.FlagPassword,
				Aliases: []string{"p"},
				Usage:   "controller password; prefer " + config.EnvPassword + " or piped stdin",
				Sources: cli.EnvVars(config.EnvPassword),
			},
		},
		Action: tokenAction,
	}
}

// tokenAction assembles the token. RFC 7617 gives the colon to the separator, so a
// username holding one is rejected here rather than being silently truncated by
// the controller.
func tokenAction(ctx context.Context, cmd *cli.Command) error {
	username := cmd.String(config.FlagUsername)
	if username == "" {
		return fmt.Errorf("%w: username required: use --%s or %s",
			ErrUsage, config.FlagUsername, config.EnvUsername)
	}

	if strings.Contains(username, ":") {
		return fmt.Errorf("%w: username must not contain ':' (RFC 7617 reserves it as the separator)", ErrUsage)
	}

	password, err := readPassword(ctx, cmd)
	if err != nil {
		return err
	}

	token := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	_, err = fmt.Fprintln(cmd.Root().Writer, token)

	return err
}

// readPassword takes the password from the flag, the environment, or a piped
// stdin. There is no interactive prompt: suppressing the echo needs a dependency
// this CLI does not take, and a prompt that echoes would leave the credential on
// screen and in the scrollback.
func readPassword(ctx context.Context, cmd *cli.Command) (string, error) {
	// An environment variable that exists but is empty is not a password: urfave
	// reports it as set, and "export WNC_PASSWORD=" is how one is cleared.
	if p := cmd.String(config.FlagPassword); p != "" {
		return p, nil
	}

	st := runtimeFrom(ctx)
	if st == nil || !st.Streams.InPipe {
		return "", fmt.Errorf("%w: password required: use --%s, %s, or pipe it on stdin",
			ErrUsage, config.FlagPassword, config.EnvPassword)
	}

	b, err := io.ReadAll(st.Streams.In)
	if err != nil {
		return "", fmt.Errorf("reading password from stdin: %w", err)
	}

	// Only the line ending is stripped: a password may legitimately hold spaces.
	return strings.TrimRight(string(b), "\r\n"), nil
}
