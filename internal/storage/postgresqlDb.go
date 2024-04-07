package storage

import (
	"arithmetic_operations/internal/config"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PostgresqlDB struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func PostgresqlOpen(cfg *config.Config, ctx context.Context) (*PostgresqlDB, error) {
	db, err := pgxpool.New(ctx, cfg.Storage.URL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	postgresql := &PostgresqlDB{db: db}

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
