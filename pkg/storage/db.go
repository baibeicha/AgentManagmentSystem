package storage

import (
	"AgentManagmentSystem/pkg/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
	datasource string
}

func (db *DB) RunMigrations() error {
	return RunMigrations(db.Pool)
}

func Connect(url, dbname, user, password, sslmode string) (*DB, error) {
	if sslmode == "" {
		sslmode = "disable"
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, password, url, dbname, sslmode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, datasource)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect database: %w", err)
	}

	return &DB{
		Pool:       pool,
		datasource: datasource,
	}, nil
}

func ConnectSource(source *config.Datasource) (*DB, error) {
	db, err := Connect(
		source.Url,
		source.DB,
		source.User,
		source.Password,
		source.SslMode,
	)

	if err != nil {
		return nil, fmt.Errorf("error connecting to datasource: %w", err)
	}

	return db, nil
}

func (db *DB) Close() {
	db.Pool.Close()
	slog.Info("database connection pool closed successfully")
}
