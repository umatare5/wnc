package framework

import (
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	cli "github.com/umatare5/wnc/internal/framework/generate"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// GenerateCli holds dependencies for generate command operations
type GenerateCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// New creates a new instance of the Cli struct
func NewGenerateCli(c *config.Config, r *infrastructure.Repository, u *application.Usecase) GenerateCli {
	return GenerateCli{
		Config:     c,
		Repository: r,
		Usecase:    u,
	}
}

// InvokeTokenCli returns a new TokenCli struct
func (gc *GenerateCli) InvokeTokenCli() *cli.TokenCli {
	return &cli.TokenCli{
		Config:     gc.Config,
		Repository: gc.Repository,
		Usecase:    gc.Usecase,
	}
}
