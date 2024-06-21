package main

import (
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

var (
	migrationsPath string
	dbURL          string
)

func init() {
	flag.StringVar(&migrationsPath, "migrations-path", "file://migrations", "Path to migration files")
	flag.StringVar(&dbURL, "db-url", "", "URL to PostgreSQL database")
	flag.Parse()
}

func main() {
	args := flag.Args()

	switch {
	case len(args) < 1:
		log.Fatal("You must specify a migrate command (up, down, goto, etc)")
	}

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}
	defer m.Close()

	// Parse the migrate command
	command := args[0]
	switch command {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down()

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		log.Printf("Current version: %d, Dirty: %v", version, dirty)
		return
	default:
		log.Fatalf("Unknown migrate command: %s", command)
	}

	if err != nil {
		log.Fatalf("migrate %s: %v", command, err)
	}

	log.Printf("migrate %s completed", command)
}
