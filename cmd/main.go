// Package main is the wnc entry point. It holds nothing but the exit call: a
// deferred function here would not run under os.Exit, so the signal context and
// every other cleanup lives inside internal/cli.Run.
package main

import (
	"os"

	"github.com/umatare5/wnc/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
