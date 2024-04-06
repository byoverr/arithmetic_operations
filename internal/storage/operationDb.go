package storage

import (
	"arithmetic_operations/internal/config"
	"arithmetic_operations/internal/models"
	"context"
	"fmt"
)

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
