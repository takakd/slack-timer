// https://pkg.go.dev/github.com/golang-migrate/migrate/v4
// e.g CLI
// $ migrate -source ${MIGRATION_DIR} -database ${DATABASE_URL} up
// $ migrate -source ${MIGRATION_DIR} -database ${DATABASE_URL} down 2
// $ migrate -source ${MIGRATION_DIR} -database ${DATABASE_URL} force
// $ migrate -source ${MIGRATION_DIR} -database ${DATABASE_URL} create -ext sql -dir
//
package main

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"os"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/log"
	"strconv"
	"time"
)

// Create empty migration file.
// File name is {UNIX_TIMESTAMP}_name.sql
func createMigrationFile(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	dir := config.Get("MIGRATION_DIR", "")

	var err error
	if err != nil {
		_, err = os.Create(fmt.Sprintf("%s/%d_%s.up.sql", dir, time.Now().Unix(), name))
	}
	if err != nil {
		_, err = os.Create(fmt.Sprintf("%s/%d_%s.down.sql", dir, time.Now().Unix(), name))
	}
	return err
}

// Print usage.
func printUsage() {
	fmt.Print(`
Usage:	migrate COMMAND

Commands:
	up				Apply all up migration from current version
	down			Apply all down migration from current version
	new NAME		Create new migration file titled NAME
	force VERSION	Create new migration file titled NAME
`)
}

// Main
func main() {
	dir := config.Get("MIGRATION_DIR", "")
	if dir == "" {
		log.Error(fmt.Sprintf("migrate: must be set MIGRATION_DIR in environment variables MIGRATION_DIR=%s\n", dir))
		os.Exit(1)
	}

	if cmd := os.Args[1]; cmd == "create" {
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}
		err := createMigrationFile(os.Args[2])
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	url := config.Get("DATABASE_URL", "")
	if url == "" {
		log.Error(fmt.Sprintf("migrate: must be set DATABASE_URL in environment variables DATABASE_URL=%s\n", url))
		os.Exit(1)
	}

	m, err := migrate.New(dir, url)
	if err != nil {
		log.Error(fmt.Errorf("migrate: %w", err))
		os.Exit(1)
	}
	defer m.Close()

	if len(os.Args) < 1 {
		printUsage()
		os.Exit(0)
	}

	if cmd := os.Args[1]; cmd == "up" {
		if err := m.Up(); err != nil {
			log.Error(fmt.Errorf("migrate: %w", err))
		}
	} else if cmd == "down" {
		if err := m.Down(); err != nil {
			log.Error(fmt.Errorf("migrate: %w", err))
		}
	} else if cmd == "force" {
		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Error(fmt.Errorf("migrate: %w", err))
		}
		if err := m.Force(version); err != nil {
			log.Error(fmt.Errorf("migrate: %w", err))
		}
	}

	os.Exit(0)
}
