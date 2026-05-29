// Package db handles database connection and migration runner.
package db

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// New opens a postgres connection pool and pings the DB.
func New() (*sqlx.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

// RunMigrations applies all pending up migrations from the migrations directory.
func RunMigrations(db *sqlx.DB) error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL is not set")
	}

	m, err := migrate.New("file://db/migrations", dsn)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			slog.Warn("migrate source close error", "err", srcErr)
		}
		if dbErr != nil {
			slog.Warn("migrate db close error", "err", dbErr)
		}
	}()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		// If the DB is dirty from a previously failed migration, force-reset to
		// the last known-good version and retry once so the server can start.
		version, dirty, vErr := m.Version()
		if vErr == nil && dirty {
			slog.Warn("migrations: dirty state detected, forcing recovery", "version", version)
			if fErr := m.Force(int(version) - 1); fErr != nil {
				return fmt.Errorf("force migration version: %w", fErr)
			}
			if rErr := m.Up(); rErr != nil && rErr != migrate.ErrNoChange {
				return fmt.Errorf("run migrations after dirty-state recovery: %w", rErr)
			}
			slog.Info("migrations: recovered from dirty state, up to date")

			return nil
		}
		return fmt.Errorf("run migrations: %w", err)
	}

	slog.Info("migrations: up to date")

	return nil
}
