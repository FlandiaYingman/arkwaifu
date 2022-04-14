package updateloop

import "github.com/caarlos0/env/v6"

type Config struct {
	// ForceUpdate indicates whether updateloop should force an update operation (only once).
	// This environment variable is useful in development. Don't use in production unless you find it fits.
	ForceUpdate bool `env:"FORCE_UPDATE" envDefault:"false"`
	// ForceSubmit indicates whether updateloop should force a submit operation (only once).
	// This environment variable is useful in development. Don't use in production unless you find it fits.
	ForceSubmit bool `env:"FORCE_SUBMIT" envDefault:"false"`
	// ForceResVersion indicates whether updateloop should force the remote resVersion to a specific value (only once).
	// This environment variable is useful in development. Don't use in production unless you find it fits.
	ForceResVersion string `env:"FORCE_RES_VERSION" envDefault:""`
}

func ProvideConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
