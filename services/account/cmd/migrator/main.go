package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// for local run
// go run ./cmd/migrator --postgres-user=postgres --postgres-password=postgres --postgres-host=localhost --migrations-path=./migrations --migrations-table=account-migrations-table
func main() {
	var postgresUser, postgresPassword, postgresHost, migrationsPath, migrationsTable string

	flag.StringVar(&postgresUser, "postgres-user", "", "username for postgresql database")
	flag.StringVar(&postgresPassword, "postgres-password", "", "password for postgresql database")
	flag.StringVar(&postgresHost, "postgres-host", "", "hostname of postgresql database")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "account-migrations-table", "name of migrations table")
	flag.Parse()

	if postgresUser == "" {
		panic("postgres-user is required")
	}
	if postgresPassword == "" {
		panic("postgres-password is required")
	}
	if postgresHost == "" {
		panic("postgres-host is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:5432/bank?sslmode=disable&x-migrations-table=%s",
		postgresUser,
		postgresPassword,
		postgresHost,
		migrationsTable,
	)

	m, err := migrate.New(
		"file://"+migrationsPath,
		postgresURL,
	)

	if err != nil {
		panic(err)
	}
	//Up method will apply all migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply ")
			return
		}

		panic(err)
	}

	fmt.Println("all migrations applied")
}
