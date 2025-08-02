package framework

import (
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/framework/generate"
	"github.com/umatare5/wnc/internal/framework/show"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// GenerateCli holds dependencies for generate command operations
type GenerateCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// ShowCli holds dependencies for show command operations
type ShowCli struct {
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

// New creates a new instance of the Cli struct
func NewShowCli(c *config.Config, r *infrastructure.Repository, u *application.Usecase) ShowCli {
	return ShowCli{
		Config:     c,
		Repository: r,
		Usecase:    u,
	}
}

// InvokeTokenCli returns a new TokenCli struct
func (gc *GenerateCli) InvokeTokenCli() *generate.TokenCli {
	return &generate.TokenCli{
		Config:     gc.Config,
		Repository: gc.Repository,
		Usecase:    gc.Usecase,
	}
}

// InvokeClientCli returns a new ClientCli struct
func (sc *ShowCli) InvokeClientCli() *show.ClientCli {
	return &show.ClientCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeApCli returns a new ApCli struct
func (sc *ShowCli) InvokeApCli() *show.ApCli {
	return &show.ApCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeApTagCli returns a new ApTagCli struct
func (sc *ShowCli) InvokeApTagCli() *show.ApTagCli {
	return &show.ApTagCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeWlanCli returns a new WlanCli struct
func (sc *ShowCli) InvokeWlanCli() *show.WlanCli {
	return &show.WlanCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeOverviewCli returns a new OverviewCli struct
func (sc *ShowCli) InvokeOverviewCli() *show.OverviewCli {
	return &show.OverviewCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}
