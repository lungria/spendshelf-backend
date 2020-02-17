package config

import (
	"github.com/caarlos0/env"
)

// EnvironmentConfiguration is struct for all configuration params of the project
type EnvironmentConfiguration struct {
	HTTPAddr   string `env:"HTTP_ADDR" envDefault:":8081"`
	MongoURI   string `env:"SPENDSHELF_MONGO_URI" envDefault:"mongodb://root:toor@localhost:27017"`
	DBName     string `env:"SPEND_SHELF_DB" envDefault:"spendShelf"`
	MonoAPIKey string `env:"SPEND_SHELF_MONO_APIKEY" envDefault:"MONO_API"`
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
