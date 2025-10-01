package config

import (
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port    string `env:"PORT" envDefault:"8080"`
	GinMode string `env:"GIN_MODE" envDefault:"release"` // "debug" | "release" | "test"
	AppEnv  string `env:"APP_ENV" envDefault:"local"`    // "local" | "prod"
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
