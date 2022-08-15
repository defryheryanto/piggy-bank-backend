package main

import (
	"fmt"
	"log"

	"github.com/defryheryanto/piggy-bank-backend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupDatabase() *gorm.DB {
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		config.PiggyBankDBHost(),
		config.PiggyBankDBPort(),
		config.PiggyBankDBName(),
		config.PiggyBankDBUsername(),
		config.PiggyBankDBPassword(),
		config.PiggyBankDBSSLMode(),
	)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening database")
	}
	fmt.Println("Success connecting to piggy bank database")

	return db
}
