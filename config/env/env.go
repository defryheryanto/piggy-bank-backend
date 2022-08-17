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
	LISTEN_PORT            = "PORT"
	LISTEN_ADDRESS         = "ADDRESS"
	PIGGY_BANK_DB_HOST     = "PIGGY_BANK_DB_HOST"
	PIGGY_BANK_DB_PORT     = "PIGGY_BANK_DB_PORT"
	PIGGY_BANK_DB_NAME     = "PIGGY_BANK_DB_NAME"
	PIGGY_BANK_DB_USERNAME = "PIGGY_BANK_DB_USERNAME"
	PIGGY_BANK_DB_PASSWORD = "PIGGY_BANK_DB_PASSWORD"
	PIGGY_BANK_DB_SSL_MODE = "PIGGY_BANK_DB_SSL_MODE"
	JWT_SECRET_KEY         = "JWT_SECRET_KEY"
	AES_SECRET_KEY         = "AES_SECRET_KEY"
)

type env struct{}

func (e *env) ListenPort() string {
	return getStringOrDefault(LISTEN_PORT, "8080")
}

func (e *env) ListenAddress() string {
	return getStringOrDefault(LISTEN_ADDRESS, "localhost")
}

func (e *env) PiggyBankDBHost() string {
	return getStringOrDefault(PIGGY_BANK_DB_HOST, "localhost")
}

func (e *env) PiggyBankDBPort() string {
	return getStringOrDefault(PIGGY_BANK_DB_PORT, "5432")
}

func (e *env) PiggyBankDBName() string {
	return getStringOrDefault(PIGGY_BANK_DB_NAME, "piggybank")
}

func (e *env) PiggyBankDBUsername() string {
	return getStringOrDefault(PIGGY_BANK_DB_USERNAME, "piggybank")
}

func (e *env) PiggyBankDBPassword() string {
	return getStringOrDefault(PIGGY_BANK_DB_PASSWORD, "piggybank")
}

func (e *env) PiggyBankDBSSLMode() string {
	return getStringOrDefault(PIGGY_BANK_DB_SSL_MODE, "disable")
}

func (e *env) JWTSecretKey() string {
	return getStringOrDefault(JWT_SECRET_KEY, "secret")
}

func (e *env) AESSecretKey() string {
	return getStringOrDefault(AES_SECRET_KEY, "secret")
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
