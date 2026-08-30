package config

import (
	"fmt"
	"strings"
)

// tagNameLimit is the controller's own cap on all three kinds, unreadable from the model: the key
// leaves declare a pattern and no length, yet a 33-character name answers 400 on every kind.
const tagNameLimit = 32

// NormalizeTagName checks a tag name against the pattern the key leaf declares,
// `[!-~]([ -~]*[!-~])?`: printable ASCII throughout, with no leading or trailing space. The SDK
// repeats the same four checks, but only once a client exists, which is one layer too late.
func NormalizeTagName(kind, s string) (string, error) {
	if s == "" {
		return "", fmt.Errorf("%s name: must not be empty", kind)
	}

	for i := range len(s) {
		if s[i] < 0x20 || s[i] > 0x7e {
			return "", fmt.Errorf("%s name %q: must be printable ASCII", kind, s)
		}
	}

	// A flag value reaches this untrimmed where a positional argument would not, so the space
	// check is a live guard. TestTagNameKeepsTheSpacesAFlagValueCarries drives it.
	if strings.HasPrefix(s, " ") || strings.HasSuffix(s, " ") {
		return "", fmt.Errorf("%s name %q: must not begin or end with a space", kind, s)
	}

	if len(s) > tagNameLimit {
		return "", fmt.Errorf("%s name %q: at most %d characters", kind, s, tagNameLimit)
	}

	return s, nil
}
