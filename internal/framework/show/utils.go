package cli

import "github.com/umatare5/wnc/internal/config"

// hasNoData checks if the data slice is empty
func hasNoData(data []any) bool {
	return len(data) == 0
}

// isJSONFormat checks if the format is JSON
func isJSONFormat(format string) bool {
	return format == config.PrintFormatJSON
}

// isAPMisconfigured checks if the AP is misconfigured
func isAPMisconfigured(isMisconfigured bool) bool {
	return isMisconfigured
}
