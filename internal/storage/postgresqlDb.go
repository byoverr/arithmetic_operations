package storage

import (
	"arithmetic_operations/internal/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PgxIface interface {
	Begin(context.Context, *config.Config) (*PostgresqlDB, error)
	Close(context.Context) error
}

func (s *PostgresqlDB) Begin(ctx context.Context, cfg *config.Config) (*PostgresqlDB, error) {
	db, err := pgxpool.New(context.Background(), cfg.Storage.URL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	postgresql := &PostgresqlDB{db: db}
	return postgresql, err
}
func (s *PostgresqlDB) Close(ctx context.Context) error {
	// Закрываем пул соединений с базой данных
	s.db.Close()
	return nil
}

type PostgresqlDB struct {
	db *pgxpool.Pool
}

func PostgresqlOpen(cfg *config.Config) (*PostgresqlDB, error) {
	db := &PostgresqlDB{}
	postgresql, err := db.Begin(context.Background(), cfg)

	err = postgresql.Init(cfg)

	if err != nil {
		return nil, err
	}

	return postgresql, nil
}

func (s *PostgresqlDB) Init(cfg *config.Config) error {
	q := `
CREATE TABLE IF NOT EXISTS expressions (
    id SERIAL PRIMARY KEY,
    expression TEXT,
    answer VARCHAR,
    status VARCHAR,
    created_at timestamp,
    completed_at timestamp NULL 
);

CREATE TABLE IF NOT EXISTS operations (
    id SERIAL PRIMARY KEY,
    operation_kind VARCHAR UNIQUE,
    duration_in_millisec INT
);
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE,
    hash_password TEXT,
    created_at timestamp
);
`

	if _, err := s.db.Exec(context.Background(), q); err != nil {
		return err
	}

	err := s.SeedOperation(cfg)
	if err != nil {
		return err
	}
	return err
}
