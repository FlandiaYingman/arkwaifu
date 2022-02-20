package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Address string `env:"ADDRESS" envDefault:"0.0.0.0:7080"`
	// ResourceLocation is the location where Arkwaifu puts its resources.
	// The "resources" include AVG images and backgrounds.
	//
	// There will be "{config.ResourceLocation}/{resVersion}" subdirectories under ResourceLocation, which are the
	// resources of the specified version.
	ResourceLocation string `env:"RESOURCE_LOCATION" envDefault:"./arkwaifu_resource"`

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
