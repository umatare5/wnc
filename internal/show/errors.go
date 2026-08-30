// Package show fans one read out across the configured controllers, joins what
// each answered into rows, and hands them to a renderer.
package show

import "errors"

// Outcomes the fan-out reports. Both are raised only after every failure has been logged per
// controller, so the caller adds no message of its own.
var (
	// ErrPartial means at least one read failed while at least one succeeded, whether the failure
	// cost a controller its rows or only some cells. The rows that exist are still printed.
	ErrPartial = errors.New("partial failure")

	// ErrAllFailed means no controller answered.
	ErrAllFailed = errors.New("all controllers failed")
)
