package config

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

// Flag names, shared by the definitions in internal/cli and the lookups here so
// the two cannot drift.
const (
	FlagConfig      = "config"
	FlagLogLevel    = "log-level"
	FlagDryRun      = "dry-run"
	FlagController  = "controller"
	FlagAccessToken = "access-token"
	FlagInsecure    = "insecure"
	FlagFormat      = "format"
	FlagPretty      = "pretty"
	FlagTimeout     = "timeout"
	FlagSortBy      = "sort-by"
	FlagSortKeys    = "sort-keys"
	FlagSortOrder   = "sort-order"
	FlagRadio       = "radio"
	FlagSSID        = "ssid"
	FlagAPName      = "ap-name"
	FlagUsername    = "username"
	FlagPassword    = "password"
	FlagYes         = "yes"
	FlagSlot        = "slot"

	// FlagMAC names a client by address. It is its own flag rather than FlagAPName reused: a
	// client carries no name on the wire, so the address is what a row and the wire share.
	FlagMAC = "mac"

	// FlagName is the tag a set or delete leaf writes. The access-point target reuses
	// FlagAPName rather than this: show client already spells that one, so the target and
	// the filter cannot drift apart.
	FlagName = "name"

	FlagDescription   = "description"
	FlagProfile24GHz  = "profile-24ghz"
	FlagProfile5GHz   = "profile-5ghz"
	FlagProfile6GHz   = "profile-6ghz"
	FlagAPJoinProfile = "ap-join-profile"
	FlagFlexProfile   = "flex-profile"
	FlagLocalSite     = "local-site"
	FlagWLAN          = "wlan"
	FlagPolicyProfile = "policy-profile"
)

// The three tag kinds, spelt as the noun every message about them uses rather than as the set
// and delete leaf names, which are hyphenated. A controller keys each kind in its own list, so
// a name is only unique within one kind.
const (
	KindPolicyTag = "policy tag"
	KindSiteTag   = "site tag"
	KindRFTag     = "RF tag"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
)

const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

// Flag defaults, declared once so the flag definition in internal/cli and the merge
// below cannot disagree about what an unset value means.
const (
	DefaultFormat    = FormatTable
	DefaultSortOrder = OrderAsc
)

// Radio band selectors accepted by --radio. They are the display values of the
// band column, so a filter and a rendered row read the same.
const (
	Band24 = "2.4"
	Band5  = "5"
	Band6  = "6"
)

// DefaultTimeout bounds one request: internal/wnc passes it to sdk.WithTimeout,
// which sets it as http.Client.Timeout. The wall-clock ceiling of a run is this
// times the number of sequential reads a command makes per controller.
const DefaultTimeout = 60 * time.Second

// DefaultLogLevel shows a run's warnings and its failures; nothing emits at Info.
const DefaultLogLevel = "warning"

// Settings is the resolved configuration one command acts on.
type Settings struct {
	Controllers []Target
	Timeout     time.Duration
	Insecure    bool
	Format      string
	Pretty      bool
	SortBy      string
	SortOrder   string
}

func (s Settings) Descending() bool {
	return s.SortOrder == OrderDesc
}

// stringIsSet reports whether a string flag carries a usable value. urfave counts an
// environment variable that exists but is empty as set, and "export WNC_X=" is the
// normal way to clear one, so an empty value is treated as absent rather than as an
// empty host list or an empty password.
func stringIsSet(cmd *cli.Command, name string) bool {
	return cmd.IsSet(name) && cmd.String(name) != ""
}

// Resolve merges one show command's settings. Precedence is the flag or its
// environment variable, then the file, then the flag's own default: urfave marks an
// environment-sourced flag as set, so IsSet covers both of the first two.
func Resolve(cmd *cli.Command, file File, sortKeys []string, defaultSortBy string) (Settings, error) {
	targets, err := resolveTargets(cmd, file)
	if err != nil {
		return Settings{}, err
	}

	timeout, err := resolveTimeout(cmd, file)
	if err != nil {
		return Settings{}, err
	}

	format, err := resolveChoice(cmd, FlagFormat, DefaultFormat, file.Format, []string{FormatTable, FormatJSON})
	if err != nil {
		return Settings{}, err
	}

	order, err := resolveChoice(cmd, FlagSortOrder, DefaultSortOrder, nil, []string{OrderAsc, OrderDesc})
	if err != nil {
		return Settings{}, err
	}

	sortBy := defaultSortBy
	if stringIsSet(cmd, FlagSortBy) {
		sortBy = cmd.String(FlagSortBy)
		if !slices.Contains(sortKeys, sortBy) {
			return Settings{}, fmt.Errorf("--%s: accepted keys are %s",
				FlagSortBy, strings.Join(sortKeys, ", "))
		}
	}

	return Settings{
		Controllers: targets,
		Timeout:     timeout,
		Insecure:    resolveBool(cmd, FlagInsecure, file.Insecure),
		Format:      format,
		Pretty:      resolveBool(cmd, FlagPretty, file.Pretty),
		SortBy:      sortBy,
		SortOrder:   order,
	}, nil
}

// ResolveExec merges the settings an action needs. It is separate from Resolve rather
// than a call into it: an action declares none of the output or sort flags, so reading
// them through a command that never defined them would tie the exec tree to defaults
// only a show command validates.
func ResolveExec(cmd *cli.Command, file File) (Settings, error) {
	targets, err := resolveTargets(cmd, file)
	if err != nil {
		return Settings{}, err
	}

	timeout, err := resolveTimeout(cmd, file)
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		Controllers: targets,
		Timeout:     timeout,
		Insecure:    resolveBool(cmd, FlagInsecure, file.Insecure),
	}, nil
}

// resolveTargets picks the controller list, naming all three sources in the fault.
func resolveTargets(cmd *cli.Command, file File) ([]Target, error) {
	if hosts := resolveHosts(cmd); len(hosts) > 0 {
		token := resolveToken(cmd, file)
		if token == "" {
			return nil, fmt.Errorf("no token given: use --%s, %s, or a configuration file",
				FlagAccessToken, EnvAccessToken)
		}

		return ParseControllers(hosts, token)
	}

	if len(file.Controllers) > 0 {
		return TargetsFromFile(file.Controllers, resolveToken(cmd, file))
	}

	return nil, fmt.Errorf("no controller given: use --%s, %s, or a configuration file",
		FlagController, EnvController)
}

// resolveHosts picks the controller hosts. The flag wins outright when it carries
// anything, and only its absence falls through to the environment: the two are never
// merged, which is exactly the concatenation urfave would have performed had the flag
// declared the variable as a source.
func resolveHosts(cmd *cli.Command) []string {
	if hosts := cmd.StringSlice(FlagController); len(hosts) > 0 {
		return hosts
	}

	return ControllersFromEnv()
}

// resolveToken picks the one token every controller is read with. The flag carries the
// environment variable as its own source, so what is left here is the file, which is the
// safest of the three: a token on the command line is visible to every process on the host.
func resolveToken(cmd *cli.Command, file File) string {
	if token := cmd.String(FlagAccessToken); token != "" {
		return token
	}

	return deref(file.Token)
}

// resolveTimeout merges the timeout and rejects a non-positive value, which would
// otherwise make every read fail with a deadline already past.
func resolveTimeout(cmd *cli.Command, file File) (time.Duration, error) {
	d := cmd.Duration(FlagTimeout)

	if !cmd.IsSet(FlagTimeout) && file.Timeout != nil {
		d = file.Timeout.Duration()
	}

	if d <= 0 {
		return 0, fmt.Errorf("--%s %s: must be positive", FlagTimeout, d)
	}

	return d, nil
}

// resolveChoice merges a string flag restricted to a fixed set. An empty value falls back to
// the declared default rather than being rejected: writing "--format=" produces one, and
// neither caller declares an environment source that could.
func resolveChoice(
	cmd *cli.Command, name, defaultValue string, fromFile *string, allowed []string,
) (string, error) {
	v := cmd.String(name)
	if v == "" {
		v = defaultValue
	}

	if !stringIsSet(cmd, name) && fromFile != nil {
		v = *fromFile
	}

	// The rejected value is deliberately absent from the fault: -b, -o, -f and -r sit
	// beside -c and -t on the same commands, so an enum slot is a mis-paste target for
	// a credential and its message may name only the flag and the accepted set.
	if !slices.Contains(allowed, v) {
		return "", fmt.Errorf("--%s: accepted values are %s", name, strings.Join(allowed, ", "))
	}

	return v, nil
}

func resolveBool(cmd *cli.Command, name string, fromFile *bool) bool {
	if !cmd.IsSet(name) && fromFile != nil {
		return *fromFile
	}

	return cmd.Bool(name)
}

// ResolveLogLevel merges the log level. The root reads it before any subcommand
// runs, so the file's value reaches the logger that reports the rest of the run.
func ResolveLogLevel(cmd *cli.Command, file File, allowed []string) (string, error) {
	level := cmd.String(FlagLogLevel)
	if level == "" {
		level = DefaultLogLevel
	}

	if !stringIsSet(cmd, FlagLogLevel) && file.LogLevel != nil {
		level = *file.LogLevel
	}

	if !slices.Contains(allowed, level) {
		return "", fmt.Errorf("--%s: accepted values are %s",
			FlagLogLevel, strings.Join(allowed, ", "))
	}

	return level, nil
}

// ValidateFile checks the parts of a configuration file that only a show command would otherwise
// reach, and returns the controllers it declares. This is what --dry-run reports on.
func ValidateFile(file File, logLevels []string) ([]Target, error) {
	targets, err := TargetsFromFile(file.Controllers, deref(file.Token))
	if err != nil {
		return nil, err
	}

	// The file's rejected values are withheld for the same reason resolveChoice
	// withholds a flag's: the file holds a token, so a value pasted one key away
	// from it can be that token.
	if file.Format != nil && !slices.Contains([]string{FormatTable, FormatJSON}, *file.Format) {
		return nil, fmt.Errorf("format: accepted values are %s, %s", FormatTable, FormatJSON)
	}

	if file.LogLevel != nil && !slices.Contains(logLevels, *file.LogLevel) {
		return nil, fmt.Errorf("log_level: accepted values are %s", strings.Join(logLevels, ", "))
	}

	if file.Timeout != nil && file.Timeout.Duration() <= 0 {
		return nil, fmt.Errorf("timeout %s: must be positive", file.Timeout.Duration())
	}

	return targets, nil
}

// ResolveRadio validates the band filter. An unset filter is the empty string,
// which selects every band.
func ResolveRadio(cmd *cli.Command) (string, error) {
	if !stringIsSet(cmd, FlagRadio) {
		return "", nil
	}

	band := cmd.String(FlagRadio)
	if !slices.Contains([]string{Band24, Band5, Band6}, band) {
		return "", fmt.Errorf("--%s: accepted values are %s, %s, %s",
			FlagRadio, Band24, Band5, Band6)
	}

	return band, nil
}
