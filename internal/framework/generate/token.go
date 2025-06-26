package cli

import (
	"fmt"

	"github.com/umatare5/wnc/internal/application"
	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TokenCli handles token generation CLI operations
type TokenCli struct {
	Config     *config.Config
	Repository *infrastructure.Repository
	Usecase    *application.Usecase
}

// GenerateToken generates a token using the TokenUsecase
func (tc *TokenCli) GenerateToken() {
	ret := tc.Usecase.InvokeTokenUsecase().GenerateBasicAuthToken()

	fmt.Println(ret)
}
