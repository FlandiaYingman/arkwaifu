package infra

import (
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"os"
	"path/filepath"
)

type Config struct {
	Address     string `env:"ADDRESS" envDefault:"0.0.0.0:7080"`
	Root        string `env:"ROOT" envDefault:"./arkwaifu-root"`
	PostgresDSN string `env:"POSTGRES_DSN"`
	Users       []User
}

type User struct {
	ID   uuid.UUID
	Name string
}

func ProvideConfig() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.Users, err = provideUsers(cfg.Root)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func provideUsers(root string) ([]User, error) {
	path := filepath.Join(root, "users.yaml")
	file, err := os.ReadFile(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	users := make([]User, 0)
	err = yaml.Unmarshal(file, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}
