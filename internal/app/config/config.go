package config

import "github.com/caarlos0/env/v6"

type Config struct {
	Address string `env:"ADDRESS" envDefault:"0.0.0.0:7080"`

	// ResourceLocation is where Arkwaifu puts resource files.
	//
	// The resource files are the raw resources from the game. Including raw assets, raw gamedata, etc.
	// The resource files are safe to remove, as long as Arkwaifu is not running.
	// "{ResourceLocation}/{resVersion}" subdirectories are the resources of the specified version.
	ResourceLocation string `env:"RESOURCE_LOCATION" envDefault:"./arkwaifu_resource"`
	// StaticLocation is where Arkwaifu puts static files.
	//
	// The static files are produced by the resources. Including webp resources, super-resolution resources, etc.
	// Normally, Arkwaifu will not change the content of existing static files, unless a breaking change is made.
	StaticLocation string `env:"STATIC_LOCATION" envDefault:"./arkwaifu_static"`

	PostgresDSN string `env:"POSTGRES_DSN"`

	// ForceUpdate indicates whether updateloop should force an update (only once).
	// This environment variable is useful in development. Don't use in production unless you really find it fits.
	ForceUpdate bool `env:"FORCE_UPDATE" envDefault:"false"`
}

func ProvideConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
