package framework

import (
	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	cli "github.com/umatare5/wnc/internal/framework/show"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// ShowCli holds dependencies for show command operations
type ShowCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// New creates a new instance of the Cli struct
func NewShowCli(c *config.Config, r *infrastructure.Repository, u *application.Usecase) ShowCli {
	return ShowCli{
		Config:     c,
		Repository: r,
		Usecase:    u,
	}
}

// InvokeClientCli returns a new ClientCli struct
func (sc *ShowCli) InvokeClientCli() *cli.ClientCli {
	return &cli.ClientCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeApCli returns a new ApCli struct
func (sc *ShowCli) InvokeApCli() *cli.ApCli {
	return &cli.ApCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeApTagCli returns a new ApTagCli struct
func (sc *ShowCli) InvokeApTagCli() *cli.ApTagCli {
	return &cli.ApTagCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeWlanCli returns a new WlanCli struct
func (sc *ShowCli) InvokeWlanCli() *cli.WlanCli {
	return &cli.WlanCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}

// InvokeOverviewCli returns a new OverviewCli struct
func (sc *ShowCli) InvokeOverviewCli() *cli.OverviewCli {
	return &cli.OverviewCli{
		Config:     sc.Config,
		Repository: sc.Repository,
		Usecase:    sc.Usecase,
	}
}
