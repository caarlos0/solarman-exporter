package config

import (
	"github.com/caarlos0/env/v9"
	"github.com/charmbracelet/log"
)

type Config struct {
	AppID      string `env:"APP_ID,required"`
	AppSecret  string `env:"APP_SECRET,required"`
	Email      string `env:"EMAIL,required"`
	Password   string `env:"PASSWORD,required"`
	InverterSN string `env:"SN,required"`
}

func Must() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("failed to parse config", "err", err)
	}
	return cfg
}
