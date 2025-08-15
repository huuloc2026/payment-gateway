package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName string

	API struct {
		Port        string
		MetricsPort string
	}

	Processor struct {
		MetricsPort string
	}

	DB struct {
		DSN string
	}

	Redis struct {
		Addr string
		DB   int
	}

	NATS struct {
		URL     string
		Subject string
	}

	Security struct {
		HMACSecret string
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func New() *Config {
	cfg := &Config{}
	cfg.AppName = getenv("APP_NAME", "payment-gateway")

	cfg.API.Port = getenv("API_PORT", "3000")
	cfg.API.MetricsPort = getenv("API_METRICS_PORT", "2112")

	cfg.Processor.MetricsPort = getenv("PROCESSOR_METRICS_PORT", "2113")

	cfg.DB.DSN = getenv("POSTGRES_DSN", "postgres://pguser:pgpass@localhost:5432/pgw?sslmode=disable")

	cfg.Redis.Addr = getenv("REDIS_ADDR", "localhost:6379")
	cfg.Redis.DB = 0

	cfg.NATS.URL = getenv("NATS_URL", "nats://localhost:4222")
	cfg.NATS.Subject = getenv("NATS_SUBJECT", "payments.created")

	cfg.Security.HMACSecret = getenv("HMAC_SECRET", "supersecret")
	return cfg
}

func (c *Config) String() string {
	return fmt.Sprintf("API:%s DB:%s Redis:%s NATS:%s", c.API.Port, c.DB.DSN, c.Redis.Addr, c.NATS.URL)
}
