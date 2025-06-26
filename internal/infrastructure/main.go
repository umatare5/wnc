package infrastructure

import (
	"github.com/umatare5/wnc/internal/config"
)

// Repository holds configuration and provides access to different repositories.
type Repository struct {
	Config *config.Config
}

// New creates a new Repository instance with the provided configuration.
func New(c *config.Config) Repository {
	return Repository{
		Config: c,
	}
}

// InvokeClientRepository returns a new instance of the ClientRepository struct.
func (r *Repository) InvokeClientRepository() *ClientRepository {
	return &ClientRepository{
		Config: r.Config,
	}
}

// InvokeApRepository returns a new instance of the ApRepository struct.
func (r *Repository) InvokeApRepository() *ApRepository {
	return &ApRepository{
		Config: r.Config,
	}
}

// InvokeWlanRepository returns a new instance of the WlanRepository struct.
func (r *Repository) InvokeWlanRepository() *WlanRepository {
	return &WlanRepository{
		Config: r.Config,
	}
}

// InvokeRadioRepository returns a new instance of the RadioRepository struct.
func (r *Repository) InvokeRadioRepository() *RadioRepository {
	return &RadioRepository{
		Config: r.Config,
	}
}

// InvokeRrmRepository returns a new instance of the RrmRepository struct.
func (r *Repository) InvokeRrmRepository() *RrmRepository {
	return &RrmRepository{
		Config: r.Config,
	}
}

// InvokeRfRepository returns a new instance of the RfRepository struct.
func (r *Repository) InvokeRfRepository() *RfRepository {
	return &RfRepository{
		Config: r.Config,
	}
}

// InvokeDot11Repository returns a new instance of the Dot11Repository struct.
func (r *Repository) InvokeDot11Repository() *Dot11Repository {
	return &Dot11Repository{
		Config: r.Config,
	}
}
