package test

import (
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/source"
	"github.com/joho/godotenv"
	_ "github.com/mattes/migrate/source/file"
)

func LoadEnv(t *testing.T, envPath string) {
	godotenv.Load(envPath)
}

func GetEnvOrDefault(key, def string) string {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	return val
}
