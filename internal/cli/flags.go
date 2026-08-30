package cli

import (
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/log"
)

// rootFlags are the flags every command sees, except --dry-run: that one is Local
// so "wnc show ap --dry-run" is rejected as an unknown flag rather than silently
// validating nothing. Local governs parsing and not reading, so "wnc --dry-run reset
// ap" still reaches the leaf, which is what stops an action running under it.
func rootFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.FlagConfig,
			Usage:   "path to the JSON configuration file",
			Sources: cli.EnvVars(config.EnvConfig),
		},
		&cli.StringFlag{
			Name:  config.FlagLogLevel,
			Usage: "log verbosity (" + strings.Join(log.Levels(), "|") + ")",
			Value: config.DefaultLogLevel,
		},
		&cli.BoolFlag{
			Name:  config.FlagDryRun,
			Usage: "report what would happen and change nothing",
			Local: true,
		},
	}
}

// showFlags are shared by every show subcommand. urfave makes a parent's flag persistent unless
// Local is set, so declaring them here is what puts them on all nine.
func showFlags() []cli.Flag {
	return []cli.Flag{
		// Deliberately no Sources: urfave applies a slice flag's environment value at
		// the level that declares it, then appends what a deeper command supplies, so
		// "show ap -c host" with WNC_CONTROLLER exported would query the named host
		// and the exported ones. config.Resolve reads the variable itself instead.
		&cli.StringSliceFlag{
			Name:    config.FlagController,
			Aliases: []string{"c"},
			// The variable is spelled out here because Sources is absent, and urfave
			// only appends the bracketed name for a flag that declares one.
			Usage: "controller host[:port], repeatable [$" + config.EnvController + "]",
		},
		&cli.StringFlag{
			Name:    config.FlagAccessToken,
			Usage:   "Basic auth token for every controller",
			Sources: cli.EnvVars(config.EnvAccessToken),
		},
		&cli.BoolFlag{
			Name:    config.FlagInsecure,
			Aliases: []string{"k"},
			Usage:   "skip TLS certificate verification",
		},
		// -o is the conventional short name for an output format, and -f is kept beside it
		// because this CLI's own examples are written with it.
		&cli.StringFlag{
			Name:    config.FlagFormat,
			Aliases: []string{"o", "f"},
			Usage:   "output format (" + config.FormatTable + "|" + config.FormatJSON + ")",
			Value:   config.DefaultFormat,
		},
		// Pretty is a table style rather than a format: it is inert under --format
		// json, where the run still succeeds and a warning says the flag did nothing.
		&cli.BoolFlag{
			Name:  config.FlagPretty,
			Usage: "draw the table with borders and status glyphs",
		},
		&cli.DurationFlag{
			Name:        config.FlagTimeout,
			Aliases:     []string{"t"},
			Usage:       "request timeout",
			Value:       config.DefaultTimeout,
			DefaultText: config.DefaultTimeout.String(),
		},
		// Long-only on purpose: -o is --format's, and the two take disjoint values.
		&cli.StringFlag{
			Name:  config.FlagSortOrder,
			Usage: "sort direction (" + config.OrderAsc + "|" + config.OrderDesc + ")",
			Value: config.DefaultSortOrder,
		},
	}
}

// execFlags are shared by the commands that act on a controller. The four connection
// flags are spelt again rather than extracted from showFlags: extracting them would
// move --timeout from sixth to fourth in every show command's help, and the OPTIONS
// transcripts in docs would all have to move with it.
func execFlags() []cli.Flag {
	return []cli.Flag{
		// No Sources, for the reason showFlags gives: urfave appends a slice flag's
		// environment value to what a deeper command supplies.
		&cli.StringSliceFlag{
			Name:    config.FlagController,
			Aliases: []string{"c"},
			Usage:   "controller host[:port] [$" + config.EnvController + "]",
		},
		&cli.StringFlag{
			Name:    config.FlagAccessToken,
			Usage:   "Basic auth token for the controller",
			Sources: cli.EnvVars(config.EnvAccessToken),
		},
		&cli.BoolFlag{
			Name:    config.FlagInsecure,
			Aliases: []string{"k"},
			Usage:   "skip TLS certificate verification",
		},
		&cli.DurationFlag{
			Name:        config.FlagTimeout,
			Aliases:     []string{"t"},
			Usage:       "request timeout",
			Value:       config.DefaultTimeout,
			DefaultText: config.DefaultTimeout.String(),
		},
	}
}

// The binding flags, one set per tag kind. Each leaf declares only its own: a flag on the
// set parent would be persistent and would let "set rf-tag --ap-join-profile" parse into a
// field an RF tag has no room for.
func policyTagFlags() []cli.Flag {
	return []cli.Flag{
		descriptionFlag(config.KindPolicyTag),
		// Both are keys of the wlan-policy list, so one without the other binds nothing and
		// is refused rather than dropped.
		&cli.StringFlag{
			Name:  config.FlagWLAN,
			Usage: "WLAN profile to bind, required with --" + config.FlagPolicyProfile,
		},
		&cli.StringFlag{
			Name:  config.FlagPolicyProfile,
			Usage: "policy profile the WLAN is bound to, required with --" + config.FlagWLAN,
		},
	}
}

func siteTagFlags() []cli.Flag {
	return []cli.Flag{
		descriptionFlag(config.KindSiteTag),
		&cli.StringFlag{
			Name:  config.FlagAPJoinProfile,
			Usage: "AP join profile to bind",
		},
		// The leaf declares when "../is-local-site = 'false'", so a flex profile clears the
		// flag on its own and naming both is refused before anything is sent.
		&cli.StringFlag{
			Name:  config.FlagFlexProfile,
			Usage: "flex profile to bind, which clears --" + config.FlagLocalSite,
		},
		&cli.BoolFlag{
			Name:  config.FlagLocalSite,
			Usage: "mark the site local",
		},
	}
}

func rfTagFlags() []cli.Flag {
	return []cli.Flag{
		descriptionFlag(config.KindRFTag),
		&cli.StringFlag{
			Name:  config.FlagProfile24GHz,
			Usage: "2.4 GHz RF profile to bind",
		},
		&cli.StringFlag{
			Name:  config.FlagProfile5GHz,
			Usage: "5 GHz RF profile to bind",
		},
		&cli.StringFlag{
			Name:  config.FlagProfile6GHz,
			Usage: "6 GHz RF profile to bind",
		},
	}
}

// tagNameFlag names the tag a set or delete leaf writes, worded per kind as descriptionFlag
// is. The 32-character cap is the controller's own and unreadable from the model, so the help
// states it rather than leaving a 400 to teach it.
func tagNameFlag(kind string) cli.Flag {
	return &cli.StringFlag{
		Name:  config.FlagName,
		Usage: kind + " name, at most 32 characters",
	}
}

func descriptionFlag(kind string) cli.Flag {
	return &cli.StringFlag{
		Name:  config.FlagDescription,
		Usage: "description for the " + kind,
	}
}

// sortFlags builds the per-command sort pair. Neither can be declared once on the
// parent: the default differs per command, and urfave's GLOBAL OPTIONS section lists
// the root's persistent flags only, so a flag the show parent declares would work from
// a leaf while appearing in no leaf's help.
func sortFlags(defaultKey string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    config.FlagSortBy,
			Aliases: []string{"b"},
			Usage:   "sort `key` (see --" + config.FlagSortKeys + ")",
			Value:   defaultKey,
		},
		&cli.BoolFlag{
			Name:  config.FlagSortKeys,
			Usage: "print the keys --" + config.FlagSortBy + " accepts, then exit",
		},
	}
}

// radioFlag filters rows by band. The values are the band column's own display
// strings so a filter and a rendered row read alike.
func radioFlag() cli.Flag {
	return &cli.StringFlag{
		Name:    config.FlagRadio,
		Aliases: []string{"r"},
		Usage:   "band filter (" + config.Band24 + "|" + config.Band5 + "|" + config.Band6 + ")",
	}
}
