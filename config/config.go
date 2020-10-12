package config

import (
	"github.com/caarlos0/env/v6"
)

// Config describes app configuration parameters.
type Config struct {
	MonoAccountID   string `env:"SHELF_MONO_ACCOUNT_ID,required"`
	MonoAPIKey      string `env:"SHELF_MONO_API_KEY,required"`
	DBConnString    string `env:"SHELF_DB_CONN,required"`
	MonoBaseURL     string `env:"SHELF_MONO_BASE_URL,required"`
	WebAPIPort      string `env:"SHELF_API_PORT" envDefault:"8080"`
	EnableImportJob bool   `env:"SHELF_MONO_ENABLE_IMPORT" envDefault:"false"`
}

// FromEnv parses configuration based on declared environment variables.
func FromEnv() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
