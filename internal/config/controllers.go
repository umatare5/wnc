package config

import (
	"errors"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"
)

// Environment variables the flags fall back to.
const (
	// EnvController supplies --controller. The sibling exporter reads the same name
	// for its single host, so one pair of variables drives both tools.
	EnvController = "WNC_CONTROLLER"
	// EnvAccessToken supplies --access-token, the one token applied to every host
	// --controller names.
	EnvAccessToken = "WNC_ACCESS_TOKEN"
	// EnvUsername supplies --username of generate-token.
	EnvUsername = "WNC_USERNAME"
	// EnvPassword supplies --password of generate-token. Preferring it to the flag
	// keeps the credential out of the process arguments.
	EnvPassword = "WNC_PASSWORD"
)

// maxColonsInHostPort is the number of colons a host:port authority may hold. More
// than that and the value is a bare IPv6 address, which RFC 3986 requires bracketed.
const maxColonsInHostPort = 1

// Target is one controller a run queries. Host is the authority alone: the SDK forces https and
// checks the same forms later, so validating here is what makes a fault a usage fault.
type Target struct {
	Name  string
	Host  string
	Token string
}

// Faults a --controller element can carry. They are sentinels so the caller can
// prefix the element's index without ever quoting the element.
var (
	errScheme     = errors.New("scheme prefix is not accepted")
	errPathInHost = errors.New("host must not contain '/', '?' or '#'")
	errUserinfo   = errors.New("userinfo is not accepted")
	errBareIPv6   = errors.New("an IPv6 address must be bracketed")
	errBrackets   = errors.New("malformed bracketed IPv6 address")
	errPort       = errors.New("port is not a number in 0-65535")
)

// errNoFileToken is the one fault a file-wide token produces, reported once for the whole list.
var errNoFileToken = errors.New(`controllers are listed with no token: add "token" to the file`)

// syntaxHint is appended to every element fault so the operator sees the accepted
// form without the CLI having to echo what they typed.
const syntaxHint = "expected host[:port] (IPv6: [2001:db8::1]:443)"

// ParseControllers validates host[:port] elements and pairs each with the one token the caller
// resolved. An element carries no token of its own, so a forgotten one cannot leave the port
// occupying its place.
func ParseControllers(hosts []string, token string) ([]Target, error) {
	targets := make([]Target, 0, len(hosts))

	for i, elem := range hosts {
		elem = strings.TrimSpace(elem)
		if elem == "" {
			return nil, fmt.Errorf("--%s[%d]: empty element; %s", FlagController, i+1, syntaxHint)
		}

		if err := ValidateAuthority(elem); err != nil {
			return nil, fmt.Errorf("--%s[%d]: %w; %s", FlagController, i+1, err, syntaxHint)
		}

		targets = append(targets, Target{Name: elem, Host: elem, Token: token})
	}

	return targets, nil
}

// ControllersFromEnv splits EnvController on commas, which is the only list form one variable can
// carry. An unset or empty variable yields no hosts.
func ControllersFromEnv() []string {
	spec := strings.TrimSpace(os.Getenv(EnvController))
	if spec == "" {
		return nil
	}

	return strings.Split(spec, ",")
}

// ValidateAuthority runs at parse time, so a malformed host is a usage fault at exit 2 before a
// client exists. '@' is refused here because the SDK's own check quotes the element it rejects, and
// neither a hostname nor a port may hold one, so nothing legal is turned away.
func ValidateAuthority(a string) error {
	switch {
	case strings.Contains(a, "://"):
		return errScheme
	case strings.ContainsAny(a, "/?#"):
		return errPathInHost
	case strings.Contains(a, "@"):
		return errUserinfo
	case strings.HasPrefix(a, "["):
		return validateBracketed(a)
	case strings.Count(a, ":") > maxColonsInHostPort:
		return errBareIPv6
	case strings.Contains(a, ":"):
		return validatePort(a[strings.IndexByte(a, ':')+1:])
	}

	return nil
}

func validateBracketed(a string) error {
	end := strings.IndexByte(a, ']')
	if end < 0 {
		return errBrackets
	}

	if _, err := netip.ParseAddr(a[1:end]); err != nil {
		return errBrackets
	}

	rest := a[end+1:]
	if rest == "" {
		return nil
	}

	if !strings.HasPrefix(rest, ":") {
		return errBrackets
	}

	return validatePort(rest[1:])
}

func validatePort(p string) error {
	if _, err := strconv.ParseUint(p, 10, 16); err != nil {
		return errPort
	}

	return nil
}

// TargetsFromFile converts the file's controllers array and pairs every entry with the
// file's one token. A file entry may carry a display name; without one the authority
// labels the rows, matching what the flag produces.
func TargetsFromFile(entries []Controller, token string) ([]Target, error) {
	if len(entries) > 0 && token == "" {
		return nil, errNoFileToken
	}

	targets := make([]Target, 0, len(entries))

	for i, e := range entries {
		host := deref(e.Host)
		if host == "" {
			return nil, fmt.Errorf("controllers[%d]: empty host", i)
		}

		if err := ValidateAuthority(host); err != nil {
			return nil, fmt.Errorf("controllers[%d]: %w; %s", i, err, syntaxHint)
		}

		name := deref(e.Name)
		if name == "" {
			name = host
		}

		targets = append(targets, Target{Name: name, Host: host, Token: token})
	}

	return targets, nil
}

func deref(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}
