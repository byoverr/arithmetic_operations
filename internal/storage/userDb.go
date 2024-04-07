package storage

import (
	"arithmetic_operations/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (s *PostgresqlDB) CreateUser(user *models.User) (int64, error) {
	var id int
	sql := "INSERT INTO users (username, hash_password, created_at) values ($1, $2, $3) RETURNING id"

	row := s.db.QueryRow(context.Background(), sql, user.Username, user.HashPassword, user.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}
	id64 := int64(id)
	return id64, nil
}

func (s *PostgresqlDB) GetUser(username string) (*models.User, error) {
	var user models.User
	ctx := context.Background()
	row := s.db.QueryRow(ctx, "SELECT id, username, hash_password, created_at FROM users WHERE username=$1", username)
	err := row.Scan(&user.Id, &user.Username, &user.HashPassword, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no user found with username %s: %w", username, err)
		}
		return nil, fmt.Errorf("failed to scan row into user: %w", err)
	}

	return &user, nil
}
