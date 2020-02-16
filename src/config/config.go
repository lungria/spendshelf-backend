package config

import (
	"github.com/caarlos0/env"
)

// EnvironmentConfiguration is struct for all configuration params of the project
type EnvironmentConfiguration struct {
	HTTPAddr   string `env:"SPENDSHELF_HTTP_ADDR" envDefault:":8081"`
	DBName     string `env:"SPENDSHELF_DB_NAME" envDefault:"spendshelf.db"`
	Topic      string `env:"SPENDSHELF_MQTT_TOPIC" envDefault:"spendshelf/transactions"`
	BrokerHost string `env:"SPENDSHELF_MQTT_HOST" envDefault:"mqtt:1883"`
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

func (c *EnvironmentConfiguration) GetTopic() string {
	return c.Topic
}

func (c *EnvironmentConfiguration) GetBrokerHost() string {
	return c.BrokerHost
}
