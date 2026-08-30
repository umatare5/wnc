package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// leafAP is spelt by more than one tree, declared once so show, reset, enable and disable
// cannot disagree about what an operator types.
const leafAP = "ap"

// yesFlag is declared on each leaf rather than on a parent, so the flag that decides whether a
// command acts is listed in that command's own OPTIONS and not as an inherited one.
func yesFlag() cli.Flag {
	return &cli.BoolFlag{
		Name:  config.FlagYes,
		Usage: "act without the confirmation prompt",
	}
}

// apNameFlag names the access point the six leaves of reset, enable and disable act on. It
// declares no Sources deliberately: urfave counts an empty environment variable as set and
// applies it through Set, so an exported variable would make Count report a target nobody
// named in the invocation.
func apNameFlag() cli.Flag {
	return &cli.StringFlag{
		Name:  config.FlagAPName,
		Usage: "access point name, as shown in the ap_name column of wnc show ap",
	}
}

// writeWording is every operator-facing string one action produces. The four templates take the
// same arguments on purpose: a leaf that worded its prompt about one target and its report about
// another would be undetectable from the exit code.
type writeWording struct {
	question string
	dryRun   string
	sent     string
	failed   string
}

// actionPrologue runs the guards every write shares and returns what a write needs. It reads
// no target: each tree reads and checks its own flag first, so the order that keeps exit 2
// meaning nothing was sent stays visible in the runner rather than buried in here.
func actionPrologue(ctx context.Context, cmd *cli.Command) (config.Target, *wnc.Client, error) {
	st := runtimeFrom(ctx)

	target, settings, err := execTarget(cmd, st)
	if err != nil {
		return config.Target{}, nil, err
	}

	if err := requireAnswerable(st, cmd.Bool(config.FlagYes), cmd.Bool(config.FlagDryRun)); err != nil {
		return config.Target{}, nil, err
	}

	if settings.Insecure {
		st.Logger.Warn("TLS certificate verification is disabled")
	}

	client, err := wnc.NewClient(target, settings, st.Logger, UserAgent())
	if err != nil {
		return config.Target{}, nil, err
	}

	return target, client, nil
}

// confirmAndAct is the whole of every write past its last read: the dry-run report, the prompt, the
// cancellation, the one call that changes something and the line that reports it. Fourteen leaves
// reach it, and a refused run changing nothing is the property all fourteen must have.
func confirmAndAct(
	ctx context.Context, cmd *cli.Command, w writeWording,
	act func(context.Context) error, args ...any,
) error {
	// Reporting the resolved target is the whole answer a dry run can give, because nothing
	// else is knowable without acting.
	if cmd.Bool(config.FlagDryRun) {
		_, err := fmt.Fprintf(cmd.Root().Writer, w.dryRun, args...)

		return err
	}

	ok, err := confirmAction(runtimeFrom(ctx), cmd.Bool(config.FlagYes), w.question, args...)
	if err != nil {
		return err
	}

	// A declined prompt is a completed run: the operator was asked and said no, which is the
	// command working rather than failing.
	if !ok {
		_, err = fmt.Fprintln(cmd.Root().Writer, "canceled")

		return err
	}

	if err := act(ctx); err != nil {
		return fmt.Errorf("%s: %s", fmt.Sprintf(w.failed, args...), wnc.Message(err))
	}

	// Whether "sent" may claim a completion is each wording's own business. save-config's RPC
	// declares an output container and the six tag leaves name the one node the controller answered
	// 201 or 204 for, so those seven say what happened; the other seven post an RPC declaring no
	// output, where a 204 establishes only that the instruction was accepted.
	_, err = fmt.Fprintf(cmd.Root().Writer, w.sent, args...)

	return err
}

// runAPAction is the whole of the four leaves that name one access point through --ap-name.
// Everything decidable without the controller is decided first so exit 2 keeps meaning nothing was
// sent, and the resolved address is discarded because all four RPCs take their ap-name arm.
func runAPAction(
	ctx context.Context, cmd *cli.Command, w writeWording,
	act func(ctx context.Context, c *wnc.Client, name string) error,
) error {
	name, err := requireAPName(cmd)
	if err != nil {
		return err
	}

	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	if _, err := resolveAP(ctx, client, target, name); err != nil {
		return err
	}

	return confirmAndAct(ctx, cmd, w, func(ctx context.Context) error {
		return act(ctx, client, name)
	}, name, target.Name)
}

// resolveAP settles that the controller holds the access point the operator named, and
// returns the base radio address the keyed radio read is keyed on. It is the one read every
// action makes: measured on 17.12, 17.15 and 17.18, ap-name-mac-map answers 404 for a name no
// access point holds, so a wrong name and a wrong --controller are told apart before a write.
func resolveAP(
	ctx context.Context, client *wnc.Client, target config.Target, name string,
) (string, error) {
	mac, found, err := client.APRadioMACByName(ctx, name)
	if err != nil {
		if cause, _ := wnc.Classify(err); cause == wnc.CauseNotFound {
			return "", absentAP(target, name)
		}

		return "", fmt.Errorf("reading ap-name-mac-map from %s: %s", target.Name, wnc.Message(err))
	}

	if !found {
		return "", absentAP(target, name)
	}

	return mac, nil
}

// absentAP words a name the controller does not hold. The 404 and the row that never arrived
// reach here differently and mean the same thing, so they are reported the same way.
func absentAP(target config.Target, name string) error {
	return fmt.Errorf("%s holds no access point named %s", target.Name, name)
}

// requireOne reads the flag that names the single thing a write acts on. Count and not IsSet:
// urfave keeps the last value of a repeated flag, so a second occurrence would move the target
// with nothing said.
func requireOne(cmd *cli.Command, flag, noun string) (string, error) {
	switch n := cmd.Count(flag); {
	case n == 0:
		return "", fmt.Errorf("%w: %s requires --%s: the %s's name", ErrUsage, cmd.Name, flag, noun)
	case n > 1:
		return "", fmt.Errorf("%w: one %s per invocation, --%s given %d times", ErrUsage, noun, flag, n)
	}

	return cmd.String(flag), nil
}

// requireAPName reads --ap-name. Only presence and emptiness are checked: measured on 17.18,
// ap-name-mac-map answers 404 for a 256-character key and for one holding a space, a slash or
// a multi-byte character alike, so the controller distinguishes no grammar this could enforce
// and one invented here would refuse a name it holds. The emptiness check stays because the
// SDK answers an empty key with its own not-found, which would reach the operator as a read
// failure at exit 1.
func requireAPName(cmd *cli.Command) (string, error) {
	name, err := requireOne(cmd, config.FlagAPName, "access point")
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("%w: --%s must not be empty", ErrUsage, config.FlagAPName)
	}

	return name, nil
}

// execTarget resolves the single controller an action runs against. An access point lives on
// one controller and this CLI must not guess which, so naming several is refused before
// anything is sent and exit 2 keeps meaning that nothing was.
func execTarget(cmd *cli.Command, st *runtimeState) (config.Target, config.Settings, error) {
	settings, err := config.ResolveExec(cmd, st.File)
	if err != nil {
		return config.Target{}, config.Settings{}, fmt.Errorf("%w: %w", ErrUsage, err)
	}

	if len(settings.Controllers) != 1 {
		return config.Target{}, config.Settings{}, fmt.Errorf(
			"%w: this command acts on one controller, %d given", ErrUsage, len(settings.Controllers))
	}

	return settings.Controllers[0], settings, nil
}

// requireAnswerable refuses a run that could never be confirmed. A piped stdin cannot answer
// a prompt, so this is settled before the client exists rather than after a read has already
// gone out on behalf of a run that was going to be refused.
func requireAnswerable(st *runtimeState, yes, dryRun bool) error {
	if yes || dryRun || !st.Streams.InPipe {
		return nil
	}

	return fmt.Errorf("%w: stdin is not a terminal: pass --%s to act without a prompt",
		ErrUsage, config.FlagYes)
}

// confirmAction asks before acting, and is shared by every command that writes. The target is
// what the operator typed, so what the prompt still catches is the controller beside it — and
// on the radio leaf the band, which is the controller's reading rather than the operator's.
func confirmAction(st *runtimeState, yes bool, question string, args ...any) (bool, error) {
	if yes {
		return true, nil
	}

	_, err := fmt.Fprintf(st.Streams.Out, question, args...)
	if err != nil {
		return false, err
	}

	line, err := bufio.NewReader(st.Streams.In).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return false, fmt.Errorf("reading the confirmation: %w", err)
	}

	answer := strings.ToLower(strings.TrimSpace(line))

	return answer == "y" || answer == "yes", nil
}
