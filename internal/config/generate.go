package config

import (
	"github.com/jinzhu/configor"
	"github.com/umatare5/wnc/pkg/log"
	"github.com/urfave/cli/v3"
)

const (
	UsernameFlagName = "username"
	PasswordFlagName = "password"
)

// GenerateCmdConfig holds generate command configuration
type GenerateCmdConfig struct {
	Username string
	Password string
}

// SetGenerateCmdConfig initializes the configuration
func (c *Config) SetGenerateCmdConfig(cli *cli.Command) {
	err := c.validateGenerateCliFlags(cli)
	if err != nil {
		log.Fatal(err)
	}

	cfg := GenerateCmdConfig{
		Username: cli.String(UsernameFlagName),
		Password: cli.String(PasswordFlagName),
	}

	err = configor.New(&configor.Config{}).Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	c.GenerateCmdConfig = cfg
}

// validateGenerateCliFlags checks if the flags are valid
func (c *Config) validateGenerateCliFlags(cli *cli.Command) error {
	if cli.String(UsernameFlagName) == "" {
		log.Fatal("error: username is empty")
	}
	if cli.String(PasswordFlagName) == "" {
		log.Fatal("error: password is empty")
	}

	return nil
}
