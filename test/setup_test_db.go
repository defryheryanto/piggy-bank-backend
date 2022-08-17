package test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	postgres_migrate "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseInfos struct {
	Host     string
	Port     string
	Name     string
	Username string
	Password string
	SSLMode  string
}

func SetupTestDatabase(t *testing.T, envPath, migrationPath string) *gorm.DB {
	LoadEnv(t, envPath)
	dbInfo := DatabaseInfos{
		Host:     GetEnvOrDefault("PIGGY_BANK_TEST_DB_HOST", "piggybank"),
		Port:     GetEnvOrDefault("PIGGY_BANK_TEST_DB_PORT", "piggybank"),
		Name:     GetEnvOrDefault("PIGGY_BANK_TEST_DB_NAME", "piggybank"),
		Username: GetEnvOrDefault("PIGGY_BANK_TEST_DB_USERNAME", "piggybank"),
		Password: GetEnvOrDefault("PIGGY_BANK_TEST_DB_PASSWORD", "piggybank"),
		SSLMode:  GetEnvOrDefault("PIGGY_BANK_TEST_DB_SSL_MODE", "piggybank"),
	}
	migrateDB(t, dbInfo, migrationPath)
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbInfo.Host,
		dbInfo.Port,
		dbInfo.Name,
		dbInfo.Username,
		dbInfo.Password,
		dbInfo.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		t.Errorf("Error opening database")
	}

	return db
}

func TruncateTables(t *testing.T, db *gorm.DB, tables []string) {
	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", strings.Join(tables, ", "))
	db.Exec(query)
}

func migrateDB(t *testing.T, dbInfo DatabaseInfos, migrationPath string) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			dbInfo.Username,
			dbInfo.Password,
			dbInfo.Host,
			dbInfo.Port,
			dbInfo.Name,
			dbInfo.SSLMode,
		),
	)
	if err != nil {
		t.Errorf("error opening migration database %v", err)
	}

	driver, err := postgres_migrate.WithInstance(db, &postgres_migrate.Config{})
	if err != nil {
		t.Errorf("error generating postgres instance %v", err)
	}

	fsrc, err := (&file.File{}).Open(migrationPath)
	if err != nil {
		t.Errorf("error opening migration files %v", err)
	}

	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		t.Errorf("error generating migrate instance %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Errorf("error migrating db %v", err)
	}
}
