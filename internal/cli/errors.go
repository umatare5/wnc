package cli

import "errors"

// ErrUsage marks a fault in what the operator asked for: a flag that does not
// parse, an unknown command, a rejected value, or a configuration file that cannot
// be read. Everything wrapping it exits 2.
var ErrUsage = errors.New("invalid usage")
