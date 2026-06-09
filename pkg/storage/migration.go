package storage

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func RunMigrations(db *pgxpool.Pool) error {
	sqlDB := stdlib.OpenDB(*db.Config().ConnConfig)
	defer sqlDB.Close()

	dbDriver, err := pgx.WithInstance(sqlDB, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("error creating database driver: %w", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("error creating embed migration files: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"pgx5",
		dbDriver,
	)
	if err != nil {
		return fmt.Errorf("error initializing migrate: %w", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("error applying migrations: %w", err)
	}

	return nil
}
