package storage

import (
	"arithmetic_operations/orchestrator/config"
	"arithmetic_operations/orchestrator/models"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PostgresqlDB struct {
	db *pgxpool.Pool
}

func PostgresqlOpen(cfg *config.Config) (*PostgresqlDB, error) {
	db, err := pgxpool.New(context.Background(), cfg.Storage.URL)
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	postgresql := &PostgresqlDB{db: db}
	err = postgresql.Init(cfg)

	if err != nil {
		return nil, err
	}

	return postgresql, nil
}

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

func (s *PostgresqlDB) CreateOperation(operation *models.Operation) error {
	const insertOperationSQL = `INSERT INTO operations (operation_kind, duration_in_millisec) VALUES ($1, $2)`

	// Use Exec to execute the INSERT statement
	_, err := s.db.Exec(context.Background(), insertOperationSQL, operation.OperationKind, operation.DurationInMilliSecond)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (s *PostgresqlDB) ReadAllOperations() ([]*models.Operation, error) {
	var operations []*models.Operation

	// Use pgxpool to run the query
	conn, err := s.db.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to acquire a connection: %w", err)
	}
	defer conn.Release()

	// Execute the query and iterate over the rows
	rows, err := conn.Query(context.Background(), `SELECT operation_kind, duration_in_millisec FROM operations`)
	if err != nil {
		return nil, fmt.Errorf("failed to query all operations: %w", err)
	}

	for rows.Next() {
		op := &models.Operation{}
		if err := rows.Scan(&op.OperationKind, &op.DurationInMilliSecond); err != nil {
			return nil, fmt.Errorf("failed to scan row into operation: %w", err)
		}
		operations = append(operations, op)
	}
	rows.Close()

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration over rows: %w", err)
	}

	return operations, nil
}

func (s *PostgresqlDB) UpdateOperation(operation *models.Operation) error {
	ctx := context.Background()
	const updateQuery = `UPDATE operations SET duration_in_millisec = $1 WHERE operation_kind = $2`

	// Execute the update query without preparation since pgx handles parameter substitution safely
	commandTag, err := s.db.Exec(ctx, updateQuery, operation.DurationInMilliSecond, operation.OperationKind)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	// Check if any rows were updated
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("no rows matched given operation kind")
	}

	return nil
}

func (s *PostgresqlDB) SeedOperation(cfg *config.Config) error {
	var databaseOperationsIsCreated = false
	operationsInDatabase, err := s.ReadAllOperations()

	if err != nil {
		return err
	}

	if len(operationsInDatabase) == cfg.Operation.CountOperation {
		databaseOperationsIsCreated = true
	}

	operations := []*models.Operation{
		{OperationKind: models.Addition, DurationInMilliSecond: cfg.Operation.DurationInMilliSecondAddition},
		{OperationKind: models.Subtraction, DurationInMilliSecond: cfg.Operation.DurationInMilliSecondSubtraction},
		{OperationKind: models.Multiplication, DurationInMilliSecond: cfg.Operation.DurationInMilliSecondMultiplication},
		{OperationKind: models.Division, DurationInMilliSecond: cfg.Operation.DurationInMilliSecondDivision},
	}
	if databaseOperationsIsCreated {
		for _, operation := range operations {
			err := s.UpdateOperation(operation)
			if err != nil {
				return fmt.Errorf("failed to update operation: %w", err)
			}
		}
	} else {
		for _, operation := range operations {
			err := s.CreateOperation(operation)
			if err != nil {
				return fmt.Errorf("failed to create operation: %w", err)
			}
		}
	}
	return nil
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
