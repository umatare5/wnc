// Package humanize provides human-readable formatting functions
package humanize

import (
	"github.com/dustin/go-humanize"
)

// FormatComma formats a number with comma separators for thousands
func FormatComma(n int64) string {
	return humanize.Comma(n)
}

// FormatBytes formats byte values with appropriate units and comma-separated numbers
func FormatBytes(bytes int64) string {
	if bytes < 1024 {
		return FormatComma(bytes) + " B"
	}
	kb := bytes / 1024
	return FormatComma(kb) + " KB"
}

// FormatTimeoutSeconds formats timeout values in seconds with comma-separated numbers
func FormatTimeoutSeconds(seconds int64) string {
	return FormatComma(seconds) + "s"
}
