package infra

import "github.com/caarlos0/env/v6"

type Config struct {
	Address     string `env:"ADDRESS" envDefault:"0.0.0.0:7080"`
	Root        string `env:"ROOT" envDefault:"./arkwaifu-root"`
	PostgresDSN string `env:"POSTGRES_DSN"`
}

func ProvideConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
