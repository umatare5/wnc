package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// deauthCommand drops a client's session. It is a flat tree rather than a reset leaf: that tree
// restarts an access point and keys every leaf on a name the controller holds, and a client has
// neither.
func deauthCommand() *cli.Command {
	return &cli.Command{
		Name:      "deauth",
		Usage:     "Deauthenticate a client on a controller",
		UsageText: synopsisChoice(config.FlagMAC, config.FlagUsername),
		Description: "--mac and --username are the mac and username columns of wnc show client,\n" +
			"and one invocation gives one of them. The controller resolves it first, so a\n" +
			"value it holds no client at is refused before the RPC, which answers the same\n" +
			"whether or not a client was there. A username may hold more than one session,\n" +
			"and the prompt says how many. The client is dropped and reconnects on its own\n" +
			"within about four minutes. The operation is absent before 17.15. Pass --dry-run\n" +
			"to name the target and change nothing.",
		Flags:  append(execFlags(), macFlag(), clientUsernameFlag(), yesFlag()),
		Action: runDeauth,
	}
}

// macFlag names the client this leaf acts on. It declares no Sources, for the reason
// apNameFlag gives.
func macFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  config.FlagMAC,
		Usage: "client MAC address, as shown in the mac column of wnc show client",
	}
}

// clientUsernameFlag names the client's own username, which is not the controller account
// generate-token asks for under the same flag name. Declaring no Sources is what keeps the two
// apart: WNC_USERNAME supplies that one, and reading it here would let an exported controller
// login select whichever clients happen to carry it.
func clientUsernameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  config.FlagUsername,
		Usage: "client username, as shown in the username column of wnc show client",
	}
}

// runDeauth drops one arm's worth of sessions. The read before the post is the guard and not a
// courtesy: the RPC answers 204 for an identifier associated to nothing exactly as it does for a
// session it dropped, so without it a reported deauthentication and a wrong target are the same
// output.
func runDeauth(ctx context.Context, cmd *cli.Command) error {
	flag, value, err := requireDeauthTarget(cmd)
	if err != nil {
		return err
	}

	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	if flag == config.FlagUsername {
		return deauthByUsername(ctx, cmd, client, target, value)
	}

	return deauthByMAC(ctx, cmd, client, target, value)
}

// deauthByMAC drops the one session at an address. The address the row carries is what the prompt,
// the report and the wire all name, because the controller already serves it in the form the SDK
// normalizes to.
func deauthByMAC(
	ctx context.Context, cmd *cli.Command, client *wnc.Client, target config.Target, mac string,
) error {
	wire, found, err := client.ClientByMAC(ctx, mac)
	if err != nil {
		if cause, _ := wnc.Classify(err); cause == wnc.CauseNotFound {
			return absentClient(target, mac)
		}

		return readFailure(target, err)
	}

	if !found {
		return absentClient(target, mac)
	}

	return confirmAndAct(ctx, cmd, writeWording{
		question: "Deauthenticate %s on %s? It is dropped and reconnects on its own. [y/N]: ",
		dryRun:   "%s on %s: would deauthenticate\n",
		sent:     "%s on %s: deauthenticate sent\n",
		failed:   "deauthenticating %s on %s",
	}, func(ctx context.Context) error {
		return absentOperation(client.DeauthenticateClientByMAC(ctx, wire))
	}, wire, target.Name)
}

// deauthByUsername drops every session under a username. The RPC's leaf states no cardinality, so
// the prompt names the number the controller holds; the collection read answers no 404 the way the
// keyed one does, so a failure here is a failure to read and not an absent client.
func deauthByUsername(
	ctx context.Context, cmd *cli.Command, client *wnc.Client, target config.Target, username string,
) error {
	sessions, err := client.ClientsByUsername(ctx, username)
	if err != nil {
		return readFailure(target, err)
	}

	if sessions == 0 {
		return fmt.Errorf("%s holds no client authenticated as %s", target.Name, username)
	}

	return confirmAndAct(ctx, cmd, writeWording{
		question: "Deauthenticate %s on %s? Each is dropped and reconnects on its own. [y/N]: ",
		dryRun:   "%s on %s: would deauthenticate\n",
		sent:     "%s on %s: deauthenticate sent\n",
		failed:   "deauthenticating %s on %s",
	}, func(ctx context.Context) error {
		return absentOperation(client.DeauthenticateClientByUsername(ctx, username))
	}, underUsername(sessions, username), target.Name)
}

// underUsername words the resolved target of the username arm. It is a count and not an address
// because the controller canonicalizes no username, so the number of sessions is the only thing
// the read adds to what the operator typed.
func underUsername(sessions int, username string) string {
	if sessions == 1 {
		return "1 client authenticated as " + username
	}

	return fmt.Sprintf("%d clients authenticated as %s", sessions, username)
}

// readFailure words a resolve that did not answer. Both arms read the same collection, so both
// name it, and Message is what keeps the response body an APIError carries out of the line.
func readFailure(target config.Target, err error) error {
	return fmt.Errorf("reading common-oper-data from %s: %s", target.Name, wnc.Message(err))
}

// absentClient words an address the controller holds no client at. The 404 the keyed read
// answers and a 200 carrying no row reach here differently and mean the same thing, so they are
// reported the same way — the shape absentAP takes for a name.
func absentClient(target config.Target, mac string) error {
	return fmt.Errorf("%s holds no client at %s", target.Name, mac)
}

// requireDeauthTarget reads the arm of the RPC's mandatory choice this invocation fills, and not
// through requireOne, whose message calls the target a name. Only presence and emptiness are
// checked: the SDK normalizes the address spellings, the username leaf is a bare string the schema
// restricts in no way, and an empty username is the value most clients carry.
func requireDeauthTarget(cmd *cli.Command) (flag, value string, err error) {
	macs, usernames := cmd.Count(config.FlagMAC), cmd.Count(config.FlagUsername)

	switch {
	case macs == 0 && usernames == 0:
		return "", "", fmt.Errorf("%w: %s requires --%s or --%s: the client's address or its username",
			ErrUsage, cmd.Name, config.FlagMAC, config.FlagUsername)
	case macs > 0 && usernames > 0:
		return "", "", fmt.Errorf("%w: --%s and --%s select differently, so one invocation gives one of them",
			ErrUsage, config.FlagMAC, config.FlagUsername)
	}

	flag, noun := config.FlagMAC, "client"
	if usernames > 0 {
		flag, noun = config.FlagUsername, "username"
	}

	if n := cmd.Count(flag); n > 1 {
		return "", "", fmt.Errorf("%w: one %s per invocation, --%s given %d times",
			ErrUsage, noun, flag, n)
	}

	value = cmd.String(flag)
	if strings.TrimSpace(value) == "" {
		return "", "", fmt.Errorf("%w: --%s must not be empty", ErrUsage, flag)
	}

	return flag, value, nil
}

// absentOperation re-words the answer a release that does not serve the operation gives. The target
// was resolved a moment earlier on the same controller, so a rejected path is the reading left, and
// Message would otherwise report the status alone.
func absentOperation(err error) error {
	if cause, status := wnc.Classify(err); cause == wnc.CauseHTTP && status == http.StatusBadRequest {
		return errors.New("the controller rejected the operation, which is how a release before 17.15 answers it")
	}

	return err
}
