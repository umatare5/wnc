package cli

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/render"
	"github.com/umatare5/wnc/internal/show"
)

// showCommand groups the read-only views.
func showCommand() *cli.Command {
	return &cli.Command{
		Name:    "show",
		Aliases: []string{"s"},
		Usage:   "Display controller state",
		Flags:   showFlags(),
		Action:  parentAction,
		Commands: []*cli.Command{
			overviewCommand(),
			apCommand(),
			apJoinCommand(),
			apTagCommand(),
			clientCommand(),
			wlanCommand(),
			policyTagCommand(),
			siteTagCommand(),
			rfTagCommand(),
		},
	}
}

// absenceNote is the one sentence every show leaf's DESCRIPTION carries. A dash is not a
// zero and not a No, and that is the single thing a reader most needs from these views.
const absenceNote = `A cell reading "-" is a value the controller did not send.`

func overviewCommand() *cli.Command {
	return &cli.Command{
		Name:    "overview",
		Aliases: []string{"o"},
		Usage:   "Per-radio RF summary across 2.4, 5 and 6 GHz",
		Description: "One row per access point radio, sorted by ap_name.\n" +
			absenceNote + "\n" +
			"Admin is the radio's own state: an access-point-level disable leaves it\n" +
			"Enabled with Oper reading Down, and wnc show ap is the authority instead.",
		Flags: append(sortFlags(show.DefaultSortAPName),
			radioFlag(),
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			band, err := config.ResolveRadio(cmd)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrUsage, err)
			}

			return runShow(ctx, cmd, show.OverviewKeys(), show.DefaultSortAPName,
				show.OverviewColumns(), show.FetchOverview(band))
		},
	}
}

func apCommand() *cli.Command {
	return &cli.Command{
		Name:    leafAP,
		Aliases: []string{"a"},
		Usage:   "Associated access points",
		Description: "One row per access point in capwap-data, sorted by ap_name.\n" +
			absenceNote + "\n" +
			"Admin is the access point's own state, which an access-point-level disable\n" +
			"changes and the Admin column of wnc show overview does not.",
		Flags: sortFlags(show.DefaultSortAPName),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.APKeys(), show.DefaultSortAPName,
				show.APColumns(), show.FetchAPs)
		},
	}
}

// apJoinCommand takes the "join" alias rather than a letter that would read as a filter.
func apJoinCommand() *cli.Command {
	return &cli.Command{
		Name:    "ap-join",
		Aliases: []string{"join"},
		Usage:   "Join, discovery and DTLS outcome per access point, joined or not",
		Description: "One row per access point the controller remembers, joined or not, sorted\n" +
			"by ap_name.\n" +
			absenceNote + "\n" +
			"It is the only view that reports an access point capwap-data has dropped.",
		Flags: sortFlags(show.DefaultSortAPName),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.APJoinKeys(), show.DefaultSortAPName,
				show.APJoinColumns(), show.FetchAPJoins)
		},
	}
}

// apTagCommand takes "tag" rather than a single letter: "wnc show t -t 30s" reads
// as neither the command nor the timeout.
func apTagCommand() *cli.Command {
	return &cli.Command{
		Name:    "ap-tag",
		Aliases: []string{"tag"},
		Usage:   "Tag assignment and its resolved values, per access point",
		Description: "One row per access point, sorted by ap_name.\n" +
			absenceNote + "\n" +
			"The tag columns are the resolved tags in force; the two profile columns\n" +
			"come from the configured site tag and agree only while Tag Source is Static.",
		Flags: sortFlags(show.DefaultSortAPName),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.APTagKeys(), show.DefaultSortAPName,
				show.APTagColumns(), show.FetchAPTags)
		},
	}
}

func clientCommand() *cli.Command {
	return &cli.Command{
		Name:    "client",
		Aliases: []string{"c"},
		Usage:   "Associated wireless clients",
		Description: "One row per associated client, sorted by mac.\n" +
			absenceNote + "\n" +
			"--radio, --ssid and --ap-name narrow the list. A client whose band the\n" +
			"controller did not report is excluded by --radio, and the count is logged.",
		Flags: append(sortFlags(show.DefaultSortMAC),
			radioFlag(),
			&cli.StringFlag{
				Name:    config.FlagSSID,
				Aliases: []string{"s"},
				Usage:   "keep only clients on this SSID",
			},
			&cli.StringFlag{
				Name:  config.FlagAPName,
				Usage: "keep only clients on this AP",
			},
		),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			band, err := config.ResolveRadio(cmd)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrUsage, err)
			}

			filter := show.ClientFilter{
				Band:   band,
				SSID:   cmd.String(config.FlagSSID),
				APName: cmd.String(config.FlagAPName),
			}

			return runShow(ctx, cmd, show.ClientKeys(), show.DefaultSortMAC,
				show.ClientColumns(), show.FetchClients(filter))
		},
	}
}

func wlanCommand() *cli.Command {
	return &cli.Command{
		Name:    "wlan",
		Aliases: []string{"w"},
		Usage:   "Configured WLANs and their bound policy profiles",
		Description: "One row per WLAN and each policy profile bound to it, sorted by wlan_id,\n" +
			"so a WLAN bound under two tags appears twice and an unbound one appears\n" +
			"once with its policy columns empty.\n" +
			absenceNote,
		Flags: sortFlags(show.DefaultSortWLANID),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.WLANKeys(), show.DefaultSortWLANID,
				show.WLANColumns(), show.FetchWLANs)
		},
	}
}

// The three tag views take no alias: "s" is the show group's own and "-r" is --radio. "tag" stays
// on ap-tag, which is the one tag view an operator reaches during an incident.
func policyTagCommand() *cli.Command {
	return &cli.Command{
		Name:  "policy-tag",
		Usage: "Configured policy tags and the WLANs they bind",
		Description: "One row per WLAN binding the tag carries, sorted by policy_tag, so a tag\n" +
			"binding three WLANs appears three times and one binding none appears once.\n" +
			absenceNote + "\n" +
			"WLAN is the WLAN profile name the binding keys on, not always the SSID.",
		Flags: sortFlags(show.DefaultSortPolicyTag),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.PolicyTagKeys(), show.DefaultSortPolicyTag,
				show.PolicyTagColumns(), show.FetchPolicyTags)
		},
	}
}

func siteTagCommand() *cli.Command {
	return &cli.Command{
		Name:  "site-tag",
		Usage: "Configured site tags and the profiles they name",
		Description: "One row per site tag, sorted by site_tag.\n" +
			absenceNote + "\n" +
			"The read asks for the values in force, so a leaf a tag left at its default\n" +
			"is reported rather than arriving as an absence.",
		Flags: sortFlags(show.DefaultSortSiteTag),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.SiteTagKeys(), show.DefaultSortSiteTag,
				show.SiteTagColumns(), show.FetchSiteTags)
		},
	}
}

func rfTagCommand() *cli.Command {
	return &cli.Command{
		Name:  "rf-tag",
		Usage: "Configured RF tags and their per-band RF profiles",
		Description: "One row per RF tag, sorted by rf_tag.\n" +
			absenceNote + "\n" +
			"The read asks for the values in force, because a plain read omits the\n" +
			"built-in tag's three per-band profile names.",
		Flags: sortFlags(show.DefaultSortRFTag),
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return runShow(ctx, cmd, show.RFTagKeys(), show.DefaultSortRFTag,
				show.RFTagColumns(), show.FetchRFTags)
		},
	}
}

// runShow is the body every show subcommand shares: resolve the settings against the
// configuration file, then hand the fan-out its columns and its read. Resolution
// happens here rather than in the root's Before because a subcommand's own flags are
// not visible from the root, and IsSet is what gives a flag precedence over the file.
func runShow[R any](
	ctx context.Context,
	cmd *cli.Command,
	sortKeys []string,
	defaultSortBy string,
	cols []render.Column[R],
	fetch show.Fetcher[R],
) error {
	st := runtimeFrom(ctx)

	// Answered before anything is resolved: listing a command's own vocabulary needs
	// no controller and no token, so it must not fail on their absence.
	if cmd.Bool(config.FlagSortKeys) {
		return printSortKeys(cmd, st.Logger, sortKeys)
	}

	settings, err := config.Resolve(cmd, st.File, sortKeys, defaultSortBy)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUsage, err)
	}

	if settings.Insecure {
		st.Logger.Warn("TLS certificate verification is disabled")
	}

	// IsSet and not settings.Pretty: a configuration file carrying pretty next to
	// format json would otherwise warn on every run, when nobody asked for both in
	// this invocation.
	if cmd.IsSet(config.FlagPretty) && settings.Format == config.FormatJSON {
		st.Logger.Warn("--" + config.FlagPretty + " styles the table only; --" +
			config.FlagFormat + " " + config.FormatJSON + " output is unchanged")
	}

	return show.Run(ctx, show.Env{
		Settings:  settings,
		Logger:    st.Logger,
		Out:       st.Streams.Out,
		UserAgent: UserAgent(),
	}, cols, fetch)
}

// ignoredBySortKeys are the flags a listing parses and does not act on. Naming them
// is cheaper than the alternative: --sort-keys short-circuits before the settings are
// resolved, so a value that would have been rejected there passes unremarked.
// --access-token is deliberately absent: it is the one flag here with an environment
// source, and IsSet is true for an exported variable, so listing it would warn on
// every invocation of anyone who exports the token.
var ignoredBySortKeys = []string{
	config.FlagController, config.FlagInsecure,
	config.FlagFormat, config.FlagTimeout, config.FlagPretty,
	config.FlagSortBy, config.FlagSortOrder,
	config.FlagRadio, config.FlagSSID, config.FlagAPName,
}

// printSortKeys answers --sort-keys with one key per line. The slice is the same one --sort-by is
// validated against, so the listing and the accepted set are one source.
func printSortKeys(cmd *cli.Command, logger *logrus.Logger, keys []string) error {
	for _, name := range ignoredBySortKeys {
		if cmd.IsSet(name) {
			logger.Warnf("--%s does not affect --%s", name, config.FlagSortKeys)
		}
	}

	for _, k := range keys {
		if _, err := fmt.Fprintln(cmd.Root().Writer, k); err != nil {
			return err
		}
	}

	return nil
}
