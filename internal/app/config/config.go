package config

import (
	"net"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress        string `env:"RUN_ADDRESS" envDefault:":8081"`
	NumOfWorkers         int    `env:"NUM_OF_WORKERS" envDefault:"10"`
	BaseURL              string `env:"BASE_URL"`
	SigningKey           string `env:"SIGNING_KEY" envDefault:"some-secret-key"`
	RefreshKey           string `env:"REFRESH_KEY" envDefault:"some-refresh-secret-key"`
	TokenTTL             string `env:"TOKEN_TTL" envDefault:"24h"`
	RefreshTTL           string `env:"REFRESH_TTL" envDefault:"240h"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	DatabaseDSN          string `env:"DATABASE_URI" envDefault:"postgresql://user:password@localhost:5432/gophermart?sslmode=disable"`
}

func New() *Config {
	config := &Config{}
	env.Parse(config)
	config.SetBaseURL()
	return config
}

func (c *Config) SetBaseURL() {
	host, port, _ := net.SplitHostPort(c.ServerAddress)
	if host == "" {
		host = "localhost"
	}
	c.BaseURL = "http://" + host
	if port != "" && port != "80" {
		c.BaseURL += ":" + port
	}
}
