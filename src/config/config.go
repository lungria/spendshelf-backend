package config

import (
	"github.com/caarlos0/env"
)

// EnvironmentConfiguration is struct for all configuration params of the project
type EnvironmentConfiguration struct {
	HTTPAddr   string `env:"SPENDSHELF_HTTP_ADDR" envDefault:":8081"`
	DBName     string `env:"SPENDSHELF_DB_NAME" envDefault:"spendshelf.db"`
	MonoAPIKey string `env:"SPENDSHELF_MONO_API_KEY" envDefault:"MONO_API"`
}

// NewConfig is parsing environment variables and returns filled EnvironmentConfiguration
func NewConfig() (*EnvironmentConfiguration, error) {
	c := EnvironmentConfiguration{}
	err := env.Parse(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *EnvironmentConfiguration) GetHTTPAddr() string {
	return c.HTTPAddr
}

func (c *EnvironmentConfiguration) GetDBName() string {
	return c.DBName
}
