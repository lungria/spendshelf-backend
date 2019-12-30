package main

import (
	"github.com/caarlos0/env"
)

// Config is struct for all configuration params of the project
type Config struct {
	HTTPAddr string `env:"WEB_HOOK_ADDR" envDefault:":80"`
	MongoURI string `env:"MONGO_URI" envDefault:"mongodb://root:toor@localhost:27017"`
	DBName   string `env:"SPEND_SHELF_DB" envDefault:"spendShelf"`
}

// NewConfig is parsing environment variables and returns filled Config
func NewConfig() (*Config, error) {
	c := Config{}
	err := env.Parse(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
