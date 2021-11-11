package config

import (
	"flag"
	"net"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress        string        `env:"RUN_ADDRESS" envDefault:":8081"`
	NumOfWorkers         int           `env:"NUM_OF_WORKERS" envDefault:"10"`
	BaseURL              string        `env:"BASE_URL"`
	SigningKey           string        `env:"SIGNING_KEY" envDefault:"some-secret-key"`
	RefreshKey           string        `env:"REFRESH_KEY" envDefault:"some-refresh-secret-key"`
	TokenTTL             time.Duration `env:"TOKEN_TTL" envDefault:"24h"`
	RefreshTTL           time.Duration `env:"REFRESH_TTL" envDefault:"240h"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
	DatabaseDSN          string        `env:"DATABASE_URI" envDefault:"postgresql://user:password@localhost:5432/gophermart?sslmode=disable"`
}

func New() (*Config, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}
	config.SetBaseURL()
	return config, nil
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

func (c *Config) InitFlags() {
	flag.Func("a", "Server start address string", func(flagValue string) error {
		if flagValue != "" {
			c.ServerAddress = flagValue
		}
		c.SetBaseURL()
		return nil
	})
	flag.Func("b", "Base URL string for generated short link", func(flagValue string) error {
		if flagValue != "" {
			c.BaseURL = flagValue
		}
		return nil
	})
	flag.Func("r", "Accrual system address", func(flagValue string) error {
		if flagValue != "" {
			c.AccrualSystemAddress = flagValue
		}
		return nil
	})
	flag.Func("d", "Database DSN string", func(flagValue string) error {
		if flagValue != "" {
			c.DatabaseDSN = flagValue
		}
		return nil
	})
	flag.Parse()
}
