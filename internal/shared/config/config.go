package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port    string `env:"PORT" envDefault:"8080"`
	GinMode string `env:"GIN_MODE" envDefault:"release"` // "debug" | "release" | "test"
	AppEnv  string `env:"APP_ENV" envDefault:"local"`    // "local" | "prod"

	// Data
	DataFile string `env:"DATA_FILE" envDefault:"data/items.json"`

	// Cache
	CacheTTL  time.Duration `env:"CACHE_TTL" envDefault:"60s"`
	CacheSize int64         `env:"CACHE_SIZE" envDefault:"1000"`

	// Idempotency
	IdempotencyTTL  time.Duration `env:"IDEMPOTENCY_TTL" envDefault:"5m"`
	IdempotencySize int64         `env:"IDEMPOTENCY_SIZE" envDefault:"5000"`

	// Server timeouts
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func Load() Config {
	appEnv := getenv("APP_ENV", "local")

	if appEnv == "prod" {
		_ = godotenv.Load(".env.prod") // Load .env.prod if it exists; if not, use the environment
	} else {
		_ = godotenv.Load(".env.local") // Load .env.local if it exists; if not, use the environment
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
