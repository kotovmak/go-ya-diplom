package config

import (
	"net"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	NumOfWorkers  int    `env:"NUM_OR_WORKERS" envDefault:"10"`
	BaseURL       string `env:"BASE_URL"`
	DatabaseDSN   string `env:"DATABASE_DSN" envDefault:"postgresql://user:password@localhost:5432/gophermart?sslmode=disable"`
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
