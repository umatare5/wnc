// Package config resolves one run's settings from the flags, the environment and
// an optional JSON file, in that order of precedence.
package config

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
)

// EnvConfig names the environment variable that selects the configuration file.
const EnvConfig = "WNC_CONFIG"

// groupOtherPerm are the group and other permission bits. A configuration file
// carrying tokens should not have them set.
const groupOtherPerm fs.FileMode = 0o077

// File mirrors the on-disk configuration. Every scalar field is a pointer so a key the
// file never set stays distinguishable from one set to a zero value, which is what
// lets a flag default win only where the file is genuinely silent.
type File struct {
	Note     *string `json:"note"`
	Timeout  *Dur    `json:"timeout"`
	Insecure *bool   `json:"insecure"`
	Format   *string `json:"format"`
	Pretty   *bool   `json:"pretty"`
	LogLevel *string `json:"log_level"`
	// Token is the one Basic auth token every controller in the file is read with. It is
	// file-wide rather than per entry, so the file holds one secret however many hosts it lists.
	Token       *string      `json:"token"`
	Controllers []Controller `json:"controllers"`
}

// Controller is one entry of the file's controllers array. Name labels the rows a
// controller produced; note carries an operational remark, because JSON has no
// comments and a rejected unknown key leaves nowhere else to put one.
type Controller struct {
	Name *string `json:"name"`
	Host *string `json:"host"`
	Note *string `json:"note"`
}

// DefaultPath is the file consulted when neither the flag nor the environment names one.
func DefaultPath() string {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return filepath.Join(dir, "wnc", "config.json")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, ".config", "wnc", "config.json")
}

// Path picks the configuration file. The value passed in is the --config flag, which urfave has
// already filled from WNC_CONFIG, so the environment is not consulted twice; the bool reports
// whether the choice was explicit, which decides whether a missing file is an error.
func Path(flagValue string) (path string, explicit bool) {
	if flagValue != "" {
		return flagValue, true
	}

	return DefaultPath(), false
}

// Load reads the configuration file, and warn receives an advisory that does not stop the run. A
// missing file is an error only when the path was chosen explicitly, and the strict decode rejects
// an unknown key, a duplicate, a case-differing key, a comment and a trailing comma.
func Load(path string, explicit bool, warn func(string)) (File, error) {
	var file File

	if path == "" {
		if explicit {
			return file, errors.New("no configuration file path given")
		}

		return file, nil
	}

	f, err := os.Open(path) //nolint:gosec // the path is the operator's own --config value
	if err != nil {
		if !explicit && errors.Is(err, fs.ErrNotExist) {
			return file, nil
		}

		return file, fmt.Errorf("configuration file: %w", err)
	}
	defer func() { _ = f.Close() }()

	checkPerm(f, path, warn)

	if err := json.UnmarshalRead(f, &file, json.RejectUnknownMembers(true)); err != nil {
		return File{}, fmt.Errorf("%s: %w", path, redactDecodeError(err))
	}

	return file, nil
}

// checkPerm warns when the file is readable beyond its owner. The mode is taken from the open
// handle rather than a second stat of the path, which closes the window between the two.
func checkPerm(f *os.File, path string, warn func(string)) {
	if runtime.GOOS == "windows" || warn == nil {
		return
	}

	info, err := f.Stat()
	if err != nil {
		return
	}

	if mode := info.Mode().Perm(); mode&groupOtherPerm != 0 {
		warn(fmt.Sprintf("configuration file %s is mode %#o; it holds tokens, so 0600 is expected", path, mode))
	}
}

// redactDecodeError rebuilds a json/v2 decode failure from its JSON pointer and its cause.
// SemanticError.Error() must not be called: it appends the offending JSON value, and a nested
// TextUnmarshaler failure quotes the input it could not parse.
func redactDecodeError(err error) error {
	var serr *json.SemanticError
	if !errors.As(err, &serr) || serr.JSONPointer == "" {
		return err
	}

	// The cause is what may quote the input: a TextUnmarshaler failure carries the text it could
	// not parse, and timeout is the one member of this file that has one. Only the two causes
	// naming a member rather than a value survive.
	if serr.Err != nil && safeCause(serr.Err) {
		return fmt.Errorf("%s: %w", serr.JSONPointer, serr.Err)
	}

	if serr.GoType != nil {
		return fmt.Errorf("%s: cannot decode into %s", serr.JSONPointer, serr.GoType)
	}

	return fmt.Errorf("%s: invalid value", serr.JSONPointer)
}

// safeCause reports whether a decode cause carries a member name and never a value. Both
// sentinels are about keys, which the pointer already exposes.
func safeCause(err error) bool {
	return errors.Is(err, json.ErrUnknownName) || errors.Is(err, jsontext.ErrDuplicateName)
}
