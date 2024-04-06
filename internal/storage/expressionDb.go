package storage

import (
	"arithmetic_operations/internal/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (s *PostgresqlDB) CreateExpression(expression *models.Expression) error {
	sql := `INSERT INTO expressions (expression, answer, status, created_at, completed_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := s.db.QueryRow(context.Background(), sql, expression.Expression, expression.Answer, expression.Status, expression.CreatedAt, expression.CompletedAt).Scan(&expression.Id)
	if err != nil {
		return fmt.Errorf("failed to create expression: %w", err)
	}
	return nil
}

func (s *PostgresqlDB) ReadAllExpressions() ([]*models.Expression, error) {
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var expressions []*models.Expression
	rows, err := conn.Query(context.Background(), `SELECT id, expression, answer, status, created_at, completed_at FROM expressions`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all expressions: %w", err)
	}

	for rows.Next() {
		expr := &models.Expression{}
		err := rows.Scan(&expr.Id, &expr.Expression, &expr.Answer, &expr.Status, &expr.CreatedAt, &expr.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into expression: %w", err)
		}
		expressions = append(expressions, expr)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	return expressions, nil
}

func (s *PostgresqlDB) ReadExpression(id int) (*models.Expression, error) {
	ctx := context.Background()
	row := s.db.QueryRow(ctx, `SELECT id, expression, answer, status, created_at, completed_at FROM expressions WHERE id = $1`, id)

	var expr models.Expression
	err := row.Scan(&expr.Id, &expr.Expression, &expr.Answer, &expr.Status, &expr.CreatedAt, &expr.CompletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no expression found with id %d: %w", id, err)
		}
		return nil, fmt.Errorf("failed to scan row into expression: %w", err)
	}

	return &expr, nil
}

func (s *PostgresqlDB) ReadAllExpressionsUndone() ([]*models.Expression, error) {
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var expressions []*models.Expression
	rows, err := conn.Query(context.Background(), `SELECT id, expression, answer, status, created_at, completed_at FROM expressions WHERE status = $1`, "in process")
	if err != nil {
		return nil, fmt.Errorf("failed to query all expressions: %w", err)
	}

	for rows.Next() {
		expr := &models.Expression{}
		err := rows.Scan(&expr.Id, &expr.Expression, &expr.Answer, &expr.Status, &expr.CreatedAt, &expr.CompletedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row into expression: %w", err)
		}
		expressions = append(expressions, expr)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	return expressions, nil
}

func (s *PostgresqlDB) UpdateExpression(e *models.Expression) error {
	_, err := s.db.Exec(context.Background(), `
		UPDATE expressions SET  
		answer = $1,
		status = $2,
		completed_at = $3
		WHERE id = $4`,
		e.Answer,
		e.Status,
		e.CompletedAt,
		e.Id,
	)

	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
