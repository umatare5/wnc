package cli

import (
	"github.com/umatare5/wnc/pkg/version"
)

func getVersion() string {
	return version.Get()
}
