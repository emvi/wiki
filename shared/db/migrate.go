package db

import (
	"emviwiki/shared/config"
	"fmt"
	"github.com/emvi/logbuch"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	migrateConnectionString = `postgres://%s:%s/%s?user=%s&password=%s&sslmode=%s&sslcert=%s&sslkey=%s&sslrootcert=%s&connect_timeout=60`
)

// Migrate runs database migrations scripts and panics in case of an error.
func Migrate() {
	dir := config.Get().Migrate.Dir
	logbuch.Info("Migrating database schema...", logbuch.Fields{"dir": dir})
	m, err := migrate.New(
		"file://"+dir,
		postgresConnectionString())

	if err != nil {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"err": err})
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"err": err})
		return
	}

	if sourceErr, dbErr := m.Close(); sourceErr != nil || dbErr != nil {
		logbuch.Fatal("Error migrating database schema", logbuch.Fields{"source_err": sourceErr, "db_err": dbErr})
	}

	logbuch.Info("Successfully migrated database schema")
}

func postgresConnectionString() string {
	cfg := config.Get().Migrate
	host := cfg.Database.Host
	port := cfg.Database.Port
	user := cfg.Database.User
	password := cfg.Database.Password
	schema := cfg.Database.Schema
	sslmode := cfg.Database.SSLMode

	if host == "" {
		host = defaultHost
	}

	if port == "" {
		port = defaultPort
	}

	if user == "" {
		user = defaultUser
	}

	if password == "" {
		password = defaultPassword
	}

	if schema == "" {
		schema = defaultSchema
	}

	if sslmode == "" {
		sslmode = defaultSSLMode
	}

	return fmt.Sprintf(migrateConnectionString, host, port, schema, user, password, sslmode,
		cfg.Database.SSLCert, cfg.Database.SSLKey, cfg.Database.SSLRootCert)
}
