package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	CoreAPIURL string `envconfig:"CORE_API_URL" required:"true"`
	MongoDB    MongoDB
}

type MongoDB struct {
	URI string `envconfig:"MONGODB_URI" required:"true"`
}

// Get returns config, filled from environment variables
func Get() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
