package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/wnc"
)

// tagKind binds one tag kind to the writes that reach it. The three kinds live in three separate
// lists on the controller, keyed independently, so a name is unique within a kind and never across
// them.
type tagKind struct {
	leaf string
	noun string
	// flags the leaf declares, and verify the combinations among them the controller
	// refuses. verify runs before the client exists, so a contradiction is exit 2 with
	// nothing sent rather than a 400 after a socket was opened.
	flags  []cli.Flag
	verify func(wnc.TagFields) error
	exists func(context.Context, *wnc.Client, string) (bool, error)
	create func(context.Context, *wnc.Client, string, wnc.TagFields) error
	update func(context.Context, *wnc.Client, string, wnc.TagFields) error
	remove func(context.Context, *wnc.Client, string) error
}

// tagKinds is the whole write surface, declared once so the set and delete trees cannot
// disagree about which kinds exist.
func tagKinds() []tagKind {
	return []tagKind{
		{
			leaf:   "policy-tag",
			noun:   config.KindPolicyTag,
			flags:  policyTagFlags(),
			verify: verifyPolicyTagFields,
			exists: func(ctx context.Context, c *wnc.Client, n string) (bool, error) {
				return c.PolicyTagExists(ctx, n)
			},
			create: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.CreatePolicyTag(ctx, n, f)
			},
			update: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.UpdatePolicyTag(ctx, n, f)
			},
			remove: func(ctx context.Context, c *wnc.Client, n string) error {
				return c.DeletePolicyTag(ctx, n)
			},
		},
		{
			leaf:   "site-tag",
			noun:   config.KindSiteTag,
			flags:  siteTagFlags(),
			verify: verifySiteTagFields,
			exists: func(ctx context.Context, c *wnc.Client, n string) (bool, error) {
				return c.SiteTagExists(ctx, n)
			},
			create: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.CreateSiteTag(ctx, n, f)
			},
			update: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.UpdateSiteTag(ctx, n, f)
			},
			remove: func(ctx context.Context, c *wnc.Client, n string) error {
				return c.DeleteSiteTag(ctx, n)
			},
		},
		{
			leaf:  "rf-tag",
			noun:  config.KindRFTag,
			flags: rfTagFlags(),
			exists: func(ctx context.Context, c *wnc.Client, n string) (bool, error) {
				return c.RFTagExists(ctx, n)
			},
			create: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.CreateRFTag(ctx, n, f)
			},
			update: func(ctx context.Context, c *wnc.Client, n string, f wnc.TagFields) error {
				return c.UpdateRFTag(ctx, n, f)
			},
			remove: func(ctx context.Context, c *wnc.Client, n string) error {
				return c.DeleteRFTag(ctx, n)
			},
		},
	}
}

// setCommand creates or updates a tag, and is one of the three configuration surfaces AGENTS.md
// admits.
func setCommand() *cli.Command {
	leaves := make([]*cli.Command, 0, len(tagKinds()))

	for _, k := range tagKinds() {
		leaves = append(leaves, setLeaf(k))
	}

	return &cli.Command{
		Name:     "set",
		Usage:    "Create or update a tag on a controller",
		Flags:    execFlags(),
		Action:   parentAction,
		Commands: leaves,
	}
}

func deleteCommand() *cli.Command {
	leaves := make([]*cli.Command, 0, len(tagKinds()))

	for _, k := range tagKinds() {
		leaves = append(leaves, deleteLeaf(k))
	}

	return &cli.Command{
		Name:     "delete",
		Usage:    "Delete a tag from a controller",
		Flags:    execFlags(),
		Action:   parentAction,
		Commands: leaves,
	}
}

// setLeaf builds one kind's set command. A name that does not exist is created and one
// that does is updated in place, so the same invocation is safe to repeat.
func setLeaf(k tagKind) *cli.Command {
	return &cli.Command{
		Name:      k.leaf,
		Usage:     "Create or update one " + k.noun,
		UsageText: synopsis(config.FlagName),
		Description: "A name --name gives that the controller does not hold is created and one it\n" +
			"holds is updated, so the same command may be repeated. A field no flag names\n" +
			"is left as it is rather than cleared. Pass --dry-run to report and change\n" +
			"nothing.",
		Flags:  append([]cli.Flag{tagNameFlag(k.noun), yesFlag()}, k.flags...),
		Action: func(ctx context.Context, cmd *cli.Command) error { return runTagSet(ctx, cmd, k) },
	}
}

func deleteLeaf(k tagKind) *cli.Command {
	return &cli.Command{
		Name:      k.leaf,
		Usage:     "Delete one " + k.noun,
		UsageText: synopsis(config.FlagName),
		Description: "The name --name gives is read on the controller first, so a name it does not\n" +
			"hold is a failure rather than a silent success. Pass --dry-run to report and\n" +
			"change nothing.",
		Flags:  []cli.Flag{tagNameFlag(k.noun), yesFlag()},
		Action: func(ctx context.Context, cmd *cli.Command) error { return runTagDelete(ctx, cmd, k) },
	}
}

// runTagSet creates or updates one tag. The guard order is the reset tree's: everything
// decidable without the controller is decided first, so exit 2 keeps meaning nothing was
// sent, and the read that follows decides between a create and an update.
func runTagSet(ctx context.Context, cmd *cli.Command, k tagKind) error {
	fields := tagFields(cmd)

	// Ahead of the name, so a contradiction is named before a missing name is.
	if k.verify != nil {
		if err := k.verify(fields); err != nil {
			return fmt.Errorf("%w: %w", ErrUsage, err)
		}
	}

	name, err := requireTagName(cmd, k.noun)
	if err != nil {
		return err
	}

	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	present, err := k.exists(ctx, client, name)
	if err != nil {
		return fmt.Errorf("reading the %s %s from %s: %s", k.noun, name, target.Name, wnc.Message(err))
	}

	if present && fields.Empty() {
		_, err = fmt.Fprintf(cmd.Root().Writer,
			"%s %s on %s: exists, no field given to change\n", k.noun, name, target.Name)

		return err
	}

	verb, done, write := "Create", "created", k.create
	if present {
		verb, done, write = "Update", "updated", k.update
	}

	return confirmAndAct(ctx, cmd, writeWording{
		question: verb + " %s %s on %s? [y/N]: ",
		dryRun:   "%s %s on %s: would " + strings.ToLower(verb) + "\n",
		sent:     "%s %s on %s: " + done + "\n",
		failed:   "writing the %s %s on %s",
	}, func(ctx context.Context) error {
		return write(ctx, client, name, fields)
	}, k.noun, name, target.Name)
}

// runTagDelete removes one tag. The tag is read first so a name the controller does not
// hold is reported as such: RESTCONF answers a delete of an absent node with 404, which
// would otherwise reach the operator as an opaque status rather than as the plain fact.
func runTagDelete(ctx context.Context, cmd *cli.Command, k tagKind) error {
	name, err := requireTagName(cmd, k.noun)
	if err != nil {
		return err
	}

	target, client, err := actionPrologue(ctx, cmd)
	if err != nil {
		return err
	}

	present, err := k.exists(ctx, client, name)
	if err != nil {
		return fmt.Errorf("reading the %s %s from %s: %s", k.noun, name, target.Name, wnc.Message(err))
	}

	if !present {
		return fmt.Errorf("%s holds no %s named %s", target.Name, k.noun, name)
	}

	return confirmAndAct(ctx, cmd, writeWording{
		question: "Delete %s %s on %s? [y/N]: ",
		dryRun:   "%s %s on %s: would delete\n",
		sent:     "%s %s on %s: deleted\n",
		failed:   "deleting the %s %s on %s",
	}, func(ctx context.Context) error {
		return k.remove(ctx, client, name)
	}, k.noun, name, target.Name)
}

// requireTagName reads --name and applies the device's grammar. urfave trims a positional argument
// and does not trim a flag value, so the validator's leading- and trailing-space branch is the live
// guard here.
func requireTagName(cmd *cli.Command, kind string) (string, error) {
	name, err := requireOne(cmd, config.FlagName, kind)
	if err != nil {
		return "", err
	}

	normalized, err := config.NormalizeTagName(kind, name)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrUsage, err)
	}

	return normalized, nil
}

// verifyPolicyTagFields refuses half of the WLAN binding. Both leaves are keys of the wlan-policy
// list, so a lone half is a named field that reaches no payload, which made the run report a write
// it never made.
func verifyPolicyTagFields(f wnc.TagFields) error {
	if (f.WLAN == nil) != (f.PolicyProfile == nil) {
		return fmt.Errorf("--%s and --%s go together", config.FlagWLAN, config.FlagPolicyProfile)
	}

	return nil
}

// verifySiteTagFields refuses a local site carrying a flex profile. The flex-profile leaf declares
// when "../is-local-site = 'false'", so the pair is a when-violation, while a flex profile on its
// own is accepted and clears the flag itself.
func verifySiteTagFields(f wnc.TagFields) error {
	if f.FlexProfile != nil && f.LocalSite != nil && *f.LocalSite {
		return fmt.Errorf("--%s cannot be set with --%s: a flex profile is in force on a non-local site only",
			config.FlagLocalSite, config.FlagFlexProfile)
	}

	return nil
}

// tagFields reads the binding flags. A flag the operator did not give stays nil, which is
// what the write layer reads as "leave this alone".
func tagFields(cmd *cli.Command) wnc.TagFields {
	return wnc.TagFields{
		Description:   stringFlag(cmd, config.FlagDescription),
		Profile24GHz:  stringFlag(cmd, config.FlagProfile24GHz),
		Profile5GHz:   stringFlag(cmd, config.FlagProfile5GHz),
		Profile6GHz:   stringFlag(cmd, config.FlagProfile6GHz),
		APJoinProfile: stringFlag(cmd, config.FlagAPJoinProfile),
		FlexProfile:   stringFlag(cmd, config.FlagFlexProfile),
		LocalSite:     boolFlag(cmd, config.FlagLocalSite),
		WLAN:          stringFlag(cmd, config.FlagWLAN),
		PolicyProfile: stringFlag(cmd, config.FlagPolicyProfile),
	}
}

// stringFlag reports a flag the command declares and the operator set. A leaf declares
// only its own kind's flags, so IsSet is asked about a name that may not exist here.
func stringFlag(cmd *cli.Command, name string) *string {
	if !cmd.IsSet(name) {
		return nil
	}

	v := cmd.String(name)

	return &v
}

func boolFlag(cmd *cli.Command, name string) *bool {
	if !cmd.IsSet(name) {
		return nil
	}

	v := cmd.Bool(name)

	return &v
}
