package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/defryheryanto/piggy-bank-backend/config"
	_ "github.com/defryheryanto/piggy-bank-backend/config/env"
	_ "github.com/golang-migrate/migrate/source"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
)

func main() {
	downFlag := flag.Bool("down", false, "database migration down")
	flag.Parse()
	fmt.Println("Opening piggy bank database connection...")
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			config.PiggyBankDBUsername(),
			config.PiggyBankDBPassword(),
			config.PiggyBankDBHost(),
			config.PiggyBankDBPort(),
			config.PiggyBankDBName(),
			config.PiggyBankDBSSLMode(),
		),
	)
	fmt.Println("Piggy bank database connected.")
	if err != nil {
		log.Fatalf("error opening migration database %v", err)
	}

	fmt.Println("Generating postgres instance...")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("error generating postgres instance %v", err)
	}
	fmt.Println("Postgres instance generated.")

	fmt.Println("Opening migration files...")
	fsrc, err := (&file.File{}).Open("file://db/migrations")
	if err != nil {
		log.Fatalf("error opening migration files %v", err)
	}
	fmt.Println("Migration files opened.")

	fmt.Println("Creating migration instance...")
	m, err := migrate.NewWithInstance("file", fsrc, "postgres", driver)
	if err != nil {
		log.Fatalf("error generating migrate instance %v", err)
	}
	fmt.Println("Migration instance created.")

	if *downFlag {
		fmt.Println("Rollback migration..")
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		version, _, _ := m.Version()
		fmt.Printf("Rollback complete to version %d.\n", version)
	} else {
		fmt.Println("Migrating migration..")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
		fmt.Println("Migrate complete.")
	}
}
