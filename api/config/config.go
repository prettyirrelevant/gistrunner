package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	RedisURL        string `env:"REDIS_URL,notEmpty"`
	Environment     string `env:"ENVIRONMENT,notEmpty"`
	DatabaseURL     string `env:"DATABASE_URL,notEmpty"`
	Port            int    `env:"PORT" envDefault:"4567"`
	DockerEngineURL string `env:"DOCKER_ENGINE_URL,notEmpty"`
}

func New() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
