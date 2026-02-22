// Package config contens the configuration for the application.
package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/charmbracelet/log"
)

// Config the actual configuration.
type Config struct {
	AppID      string `env:"APP_ID,required"`
	AppSecret  string `env:"APP_SECRET,required"`
	Email      string `env:"EMAIL,required"`
	Password   string `env:"PASSWORD,required"`
	InverterSN string `env:"SN,required"`
}

// Must returns the config or exit 1.
func Must() Config {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		log.Fatal("failed to parse config", "err", err)
	}
	return cfg
}
