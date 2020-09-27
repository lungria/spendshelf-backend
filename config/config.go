package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	MonoAccountID string `env:"SHELF_MONO_ACCOUNT_ID,required"`
	MonoAPIKey    string `env:"SHELF_MONO_API_KEY,required"`
	DBConnString  string `env:"SHELF_DB_CONN,required"`
	MonoBaseURL   string `env:"SHELF_MONO_BASE_URL,required"`
}

func FromEnv() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
