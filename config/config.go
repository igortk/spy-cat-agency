package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type (
	Config struct {
		DataBaseConfig         DataBaseConfig         `envPrefix:"DB_"`
		HttpConfig             HttpConfig             `envPrefix:"HTTP_"`
		ExternalServicesConfig ExternalServicesConfig `envPrefix:"API_"`
	}

	DataBaseConfig struct {
		Host string `env:"URL" envDefault:"postgres://postgres:password@localhost:5432/SpyCatAgency?sslmode=disable"`
	}

	HttpConfig struct {
		Port int `env:"PORT"  envDefault:"8081"`
	}

	ExternalServicesConfig struct {
		TheCatAPIURL string `env:"THE_CAT_API_URL" envDefault:"https://api.thecatapi.com/v1/breeds"`
	}
)

func GetConfig() (*Config, error) {
	log.Info(fmt.Sprintf("Getting config..."))
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.New("can't parse config")
	}
	return &cfg, nil
}
