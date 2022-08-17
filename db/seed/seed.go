package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/defryheryanto/piggy-bank-backend/db/seed/seeder"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	environment := flag.String("environment", "test", "Define the environment 'development' or 'test', default to 'test'")
	accountType := flag.Bool("account_types", false, "Seed account_types table")
	flag.Parse()
	db := setupDatabase(*environment)
	if *accountType {
		fmt.Println("seeding account_types")
		seeder := seeder.NewAccountTypeSeeder(db)
		seeder.Seed()
	}

	fmt.Println("seed complete")
}

func setupDatabase(environment string) *gorm.DB {
	godotenv.Load()
	hostKey := "PIGGY_BANK_TEST_DB_HOST"
	portKey := "PIGGY_BANK_TEST_DB_PORT"
	nameKey := "PIGGY_BANK_TEST_DB_NAME"
	userKey := "PIGGY_BANK_TEST_DB_USERNAME"
	passwordKey := "PIGGY_BANK_TEST_DB_PASSWORD"
	sslModeKey := "PIGGY_BANK_TEST_DB_SSL_MODE"
	if environment == "development" {
		hostKey = "PIGGY_BANK_DB_HOST"
		portKey = "PIGGY_BANK_DB_PORT"
		nameKey = "PIGGY_BANK_DB_NAME"
		userKey = "PIGGY_BANK_DB_USERNAME"
		passwordKey = "PIGGY_BANK_DB_PASSWORD"
		sslModeKey = "PIGGY_BANK_DB_SSL_MODE"
	}
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		os.Getenv(hostKey),
		os.Getenv(portKey),
		os.Getenv(nameKey),
		os.Getenv(userKey),
		os.Getenv(passwordKey),
		os.Getenv(sslModeKey),
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database")
	}
	fmt.Printf("connected to %s\n", os.Getenv(nameKey))

	return db
}
