// Package version provides version information for this CLI application.
package version

// Version contains the version string injected at build time.
// This variable is set via -ldflags during the build process.
var Version = "dev"

// Get returns the current version of the application.
func Get() string {
	return Version
}
