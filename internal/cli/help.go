package cli

import (
	"context"
	"strings"

	"github.com/urfave/cli/v3"
)

// helpWord is the one leftover argument a group answers rather than refuses.
// HideHelpCommand on the root removes urfave's help command from the whole tree, so
// without this the word reaches parentAction as an unknown command.
const helpWord = "help"

// helpFor answers "help" with the help of whatever the remaining arguments name.
// urfave's ShowCommandHelp descends one level per call, so the path is walked here, and
// an argument that names no command stays a usage fault.
func helpFor(ctx context.Context, cmd *cli.Command, path []string) error {
	parent := cmd

	for _, name := range path {
		sub := cmd.Command(name)
		if sub == nil {
			return unknownCommand(cmd, name)
		}

		parent, cmd = cmd, sub
	}

	return showHelp(ctx, cmd, parent)
}

// showHelp prints one command's own help. A leaf is printed through its parent because
// that is urfave's only path that honors CustomHelpTemplate, which is where the
// inherited-options section lives.
func showHelp(ctx context.Context, cmd, parent *cli.Command) error {
	switch {
	case cmd.Root() == cmd:
		return cli.ShowRootCommandHelp(cmd)
	case parent != nil && parent != cmd && len(cmd.Commands) == 0:
		return cli.ShowCommandHelp(ctx, parent, cmd.Name)
	default:
		return cli.ShowSubcommandHelp(cmd)
	}
}

// attachInheritedOptions gives every leaf a help template that also lists the flags its
// parent declares: urfave builds GLOBAL OPTIONS from the root's persistent flags alone,
// so a connection flag a leaf parses through its parent appears in no leaf's help. Only a
// group's flags are copied down, because the root's are that GLOBAL OPTIONS section.
func attachInheritedOptions(root *cli.Command) {
	for _, group := range root.Commands {
		section := inheritedOptions(group.Flags)
		if section == "" {
			continue
		}

		for _, leaf := range group.Commands {
			if len(leaf.Commands) == 0 {
				leaf.CustomHelpTemplate = cli.CommandHelpTemplate + section
			}
		}
	}
}

// inheritedOptions renders the section. It is built in Go rather than by the template because
// urfave's template data is the leaf, which reaches the root's flags and never its parent's, and a
// Local flag is left out because urfave does not apply one to a subcommand.
func inheritedOptions(flags []cli.Flag) string {
	var section strings.Builder

	for _, fl := range flags {
		if local, ok := fl.(cli.LocalFlag); ok && local.IsLocal() {
			continue
		}

		section.WriteString("\n   " + fl.String())
	}

	if section.Len() == 0 {
		return ""
	}

	return "\nINHERITED OPTIONS:" + section.String() + "\n"
}
