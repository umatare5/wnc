package cli

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
)

// onUsageError turns a parse fault into an ErrUsage chain. Returning early here
// also suppresses urfave's own behavior of dumping the entire help text to stdout
// on a parse error, which would put usage text in a pipe carrying data.
func onUsageError(_ context.Context, cmd *cli.Command, err error, _ bool) error {
	raw := err.Error()
	msg := redactInvalidValue(cmd, redactUndefinedFlag(cmd, raw))

	// The near miss is scored on the argument as it arrived rather than on what survived the
	// redaction, because the redaction now withholds a name it cannot vouch for and a typo is
	// exactly the case that needs one matched. It cannot echo either way: SuggestFlag returns
	// one of this tree's own names and suggestion floors the rest.
	if arg, found := strings.CutPrefix(raw, undefinedFlag); found {
		typed := strings.TrimLeft(trimArg(arg), "-")
		msg += suggestion(typed, cli.SuggestFlag(parsedFlags(cmd), typed, false))
	}

	return usageFault(cmd, msg)
}

// undefinedFlag is urfave's own unknown-flag message up to, but not including, the single '-' its
// parser prepends to the argument, which arrives with its dashes stripped and its '=' tail cut.
const undefinedFlag = "flag provided but not defined: "

// redactUndefinedFlag keeps only a name this command declares, which is the guard
// redactInvalidValue already applies to a value. urfave splits an argument at '=' and so never
// echoes "--flag=token", but its parser appends a single word verbatim and "-p<password>" arrives
// whole, which trimArg's ':' and ' ' cuts do not touch.
func redactUndefinedFlag(cmd *cli.Command, msg string) string {
	arg, found := strings.CutPrefix(msg, undefinedFlag)
	if !found {
		return msg
	}

	if name := declaredFlagName(cmd, trimArg(arg)); name != "" {
		return undefinedFlag + "-" + name
	}

	return strings.TrimSuffix(undefinedFlag, ": ")
}

// declaredFlagName names the flag an undefined argument began with, and only where the command
// declares that name. A remainder that is command-shaped is read as a mistyped long name and gives
// no name at all, because reporting its first character would name a flag that is in fact declared
// — and an all-lower-case password takes that branch too.
func declaredFlagName(cmd *cli.Command, arg string) string {
	typed := strings.TrimLeft(arg, "-")

	switch {
	case typed == "":
		return ""
	case declaresFlag(cmd, typed):
		return typed
	case commandShaped(typed[1:]):
		return ""
	case declaresFlag(cmd, typed[:1]):
		return typed[:1]
	}

	return ""
}

// declaresFlag reports whether the command parses a flag under this exact name. Both
// redactions ask it, so neither can drift from the set the parse itself used.
func declaresFlag(cmd *cli.Command, name string) bool {
	for _, fl := range parsedFlags(cmd) {
		if slices.Contains(fl.Names(), name) {
			return true
		}
	}

	return false
}

// trimArg cuts an argument down to the name a message may repeat. Neither a flag name
// nor a command name holds ':' or ' ', so what is dropped is never part of either.
func trimArg(arg string) string {
	name, _, _ := strings.Cut(arg, ":")
	name, _, _ = strings.Cut(name, " ")

	return name
}

// invalidValue is urfave's message for a value its flag cannot parse. It quotes the
// argument as it was typed, and the type's own parser quotes it a second time.
const invalidValue = "invalid value "

// redactInvalidValue keeps the flag and drops the value. -t is --timeout while
// --access-token has no short name, so "show -t <token>" is a plausible slip and the one
// parse fault that prints an argument twice. The name is taken from the message only
// when the command declares it, so a value carrying the separator cannot reach the result.
func redactInvalidValue(cmd *cli.Command, msg string) string {
	if !strings.HasPrefix(msg, invalidValue) {
		return msg
	}

	_, rest, _ := strings.Cut(msg, " for flag -")

	if name, _, _ := strings.Cut(rest, ":"); declaresFlag(cmd, name) {
		return invalidValue + "for flag -" + name
	}

	return invalidValue + "for a flag"
}

// usageFault is the shape every parse and dispatch fault takes. FullName is the path
// from the root, so a fault at a leaf names the leaf's own help, which after the
// inherited section is the one listing carrying both its flags and its parent's.
func usageFault(cmd *cli.Command, msg string) error {
	return fmt.Errorf("%w: %s: see '%s --help'", ErrUsage, msg, cmd.FullName())
}

// suggestion offers urfave's nearest match, which is only ever one of this CLI's own names and so
// cannot echo what was typed. Three floors urfave has none of: a match must start as the typed word
// does, must be longer than an alias, and must differ from it.
func suggestion(typed, match string) string {
	name := strings.TrimLeft(match, "-")
	if len(name) < 2 || typed == "" || typed == name || typed[0] != name[0] {
		return ""
	}

	return " (did you mean " + match + "?)"
}

// parsedFlags is every flag the command parses: each ancestor's persistent ones, then
// its own, which is the set urfave assembles for the parse itself. Building it is what
// lets a fault at a leaf name a flag the parent declares.
func parsedFlags(cmd *cli.Command) []cli.Flag {
	var flags []cli.Flag

	for _, ancestor := range cmd.Lineage()[1:] {
		for _, fl := range ancestor.Flags {
			if local, ok := fl.(cli.LocalFlag); ok && local.IsLocal() {
				continue
			}

			flags = append(flags, fl)
		}
	}

	return append(flags, cmd.Flags...)
}

// attachUsageHandler sets the hook on every node of the tree. urfave consults the
// running command's own field and never a parent's, so a node left without one
// falls back to the library's path.
func attachUsageHandler(cmd *cli.Command) {
	cmd.OnUsageError = onUsageError

	for _, sub := range cmd.Commands {
		attachUsageHandler(sub)
	}
}

// refuseArgs makes a leftover positional argument a usage fault. urfave raises none of its own — a
// word that names no subcommand is handed to the action through Args() — and the count is reported
// without the word, because a leftover on generate-token can be the password a wrapper misplaced.
func refuseArgs(cmd *cli.Command) error {
	n := cmd.Args().Len()
	if n == 0 {
		return nil
	}

	msg := fmt.Sprintf("this command takes no positional arguments, %d given", n)
	if flag := targetFlag(cmd); flag != "" {
		msg += ": use --" + flag
	}

	return usageFault(cmd, msg)
}

// targetFlag names the flag a leftover word most likely belonged to, read from the leaf's own
// declarations so it cannot drift from them. Every other leaf, generate-token included, gets
// no hint rather than an invitation to put a credential on a flag.
func targetFlag(cmd *cli.Command) string {
	for _, name := range []string{config.FlagAPName, config.FlagName, config.FlagMAC} {
		for _, fl := range cmd.Flags {
			if slices.Contains(fl.Names(), name) {
				return name
			}
		}
	}

	return ""
}

// requiredPlaceholder is the value spelling each mandatory flag takes in a synopsis, keyed on
// the flag constant so a synopsis cannot name a flag this tree does not declare.
var requiredPlaceholder = map[string]string{
	config.FlagAPName:   "<ap-name>",
	config.FlagName:     "<name>",
	config.FlagSlot:     "<n>",
	config.FlagMAC:      "<mac>",
	config.FlagUsername: "<username>",
}

// synopsis stages the fragment naming a leaf's mandatory flags, which attachLeafRules completes
// with the path from the root. What a leaf refuses to run without belongs in the one line an
// operator reads first.
func synopsis(flags ...string) string {
	parts := make([]string, 0, len(flags))

	for _, name := range flags {
		parts = append(parts, "--"+name+" "+requiredPlaceholder[name])
	}

	return strings.Join(parts, " ")
}

// synopsisChoice stages the same fragment for two flags a leaf takes one of. deauth is the only
// such leaf, so the brackets are spelt here rather than in synopsis, which joins its flags as a
// conjunction.
func synopsisChoice(a, b string) string {
	return "(" + synopsis(a) + " | " + synopsis(b) + ")"
}

// attachLeafRules gives every leaf the refusal of a leftover positional argument and a USAGE line
// carrying its path, neither of which a leaf can state from inside its own literal. It is a walk of
// its own because attachUsageHandler is also passed to ConfigureShellCompletionCommand, and a
// leftover word inside that subtree is urfave's to ignore rather than this tree's to refuse.
//
// A childless node with no action is skipped, because urfave installs its own help action for one
// during setup and the closure would call nil; TestEveryLeafDeclaresAnAction keeps that from being
// a silent hole.
func attachLeafRules(cmd *cli.Command, path string) {
	path = strings.TrimSpace(path + " " + cmd.Name)

	for _, sub := range cmd.Commands {
		attachLeafRules(sub, path)
	}

	if len(cmd.Commands) > 0 || cmd.Action == nil {
		return
	}

	if cmd.UsageText != "" {
		cmd.UsageText = path + " " + cmd.UsageText + " [options]"
	}

	action := cmd.Action

	cmd.Action = func(ctx context.Context, leaf *cli.Command) error {
		if err := refuseArgs(leaf); err != nil {
			return err
		}

		return action(ctx, leaf)
	}
}

// parentAction runs when a command that only groups subcommands is invoked on its
// own. A leftover argument is an unknown command: urfave would otherwise pass it
// through as a positional argument and run the group silently. The one word answered
// rather than refused is "help", which the hidden help command would have answered.
func parentAction(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()

	if len(args) > 0 && args[0] == helpWord {
		return helpFor(ctx, cmd, args[1:])
	}

	if len(args) > 0 {
		return unknownCommand(cmd, args[0])
	}

	return showHelp(ctx, cmd, nil)
}

// unknownCommand reports a word no command answers. A match urfave scored against its own help
// alias is dropped: SuggestCommand appends "help" and "h" to every candidate list, and this tree
// removed the command that carried them. A word that is not command-shaped is not repeated at all.
func unknownCommand(cmd *cli.Command, arg string) error {
	name := trimArg(arg)
	if !commandShaped(name) {
		return usageFault(cmd, "unknown command")
	}

	match := cli.SuggestCommand(cmd.Commands, name)
	if match != helpWord && cmd.Command(match) == nil {
		match = ""
	}

	return usageFault(cmd, fmt.Sprintf("unknown command %q%s", name, suggestion(name, match)))
}

// commandShaped reports whether a word is spelt the way this tree spells its command names:
// lower-case letters, digits and hyphens. A Basic auth token carries upper case, so the shape
// separates a mistyped command from a pasted credential, imperfectly.
func commandShaped(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
		default:
			return false
		}
	}

	return true
}
