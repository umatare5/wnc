package application

import (
	"encoding/base64"
	"fmt"

	"github.com/umatare5/wnc/internal/config"
	"github.com/umatare5/wnc/internal/infrastructure"
)

// TokenUsecase handles token generation operations
type TokenUsecase struct {
	Config     *config.Config
	Repository *infrastructure.Repository
}

// GenerateBasicAuthToken generates a basic authentication token using the username and password from the configuration
func (u *TokenUsecase) GenerateBasicAuthToken() string {
	creds := fmt.Sprintf("%s:%s", u.Config.GenerateCmdConfig.Username, u.Config.GenerateCmdConfig.Password)
	return base64.StdEncoding.EncodeToString([]byte(creds))
}
