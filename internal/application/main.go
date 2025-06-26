package application

import (
	"github.com/umatare5/wnc/internal/config"

	"github.com/umatare5/wnc/internal/infrastructure"
)

// Usecase holds configuration and repository dependencies
type Usecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// New creates a new instance of the Usecase struct
func New(c *config.Config, r *infrastructure.Repository) Usecase {
	return Usecase{
		Config:     c,
		Repository: r,
	}
}

// InvokeTokenUsecase returns a new TokenUsecase struct
func (u *Usecase) InvokeTokenUsecase() *TokenUsecase {
	return &TokenUsecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeClientUsecase returns a new ClientUsecase struct
func (u *Usecase) InvokeClientUsecase() *ClientUsecase {
	return &ClientUsecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeApUsecase returns a new ApUsecase struct
func (u *Usecase) InvokeApUsecase() *ApUsecase {
	return &ApUsecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeWlanUsecase returns a new WlanUsecase struct
func (u *Usecase) InvokeWlanUsecase() *WlanUsecase {
	return &WlanUsecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}

// InvokeOverviewUsecase returns a new OverviewUsecase struct
func (u *Usecase) InvokeOverviewUsecase() *OverviewUsecase {
	return &OverviewUsecase{
		Config:     u.Config,
		Repository: u.Repository,
	}
}
