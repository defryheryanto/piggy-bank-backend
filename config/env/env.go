package env

import (
	"os"

	"github.com/defryheryanto/piggy-bank-backend/config"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	config.SetConfig(&env{})
}

const (
	LISTEN_PORT    = "PORT"
	LISTEN_ADDRESS = "ADDRESS"
)

type env struct{}

func (e *env) ListenPort() string {
	return getStringOrDefault(LISTEN_PORT, "8080")
}

func (e *env) ListenAddress() string {
	return getStringOrDefault(LISTEN_ADDRESS, "localhost")
}

func getStringOrDefault(key, def string) string {
	return getEnvOrDefault(key, def)
}

func getEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}
