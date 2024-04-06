package checker

import (
	"arithmetic_operations/internal/models"
	"errors"
)

func ValidateOperation(operation models.Operation) error {
	if operation.DurationInMilliSecond < 0 {
		return errors.New("operation duration should be more than 0")
	}

	if !models.IsAllowedOperation(operation.OperationKind) {
		return errors.New("it is not allowed operation")
	}

	return nil
}
