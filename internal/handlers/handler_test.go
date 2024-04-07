package handlers

import (
	"arithmetic_operations/internal/models"
	"arithmetic_operations/internal/prettylogger"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandlerGetAllOperations(t *testing.T) {
	// Mock logger
	opts := prettylogger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}
	handler := prettylogger.NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Mock operation reader function
	mockOperations := []*models.Operation{
		{OperationKind: "addition", DurationInMilliSecond: 100},
		{OperationKind: "subtraction", DurationInMilliSecond: 200},
	}
	operationReader := func() ([]*models.Operation, error) {
		return mockOperations, nil
	}

	// Create a request to pass to the handler
	req := httptest.NewRequest("GET", "/operations", nil)
	rec := httptest.NewRecorder()

	// Call the handler function with the mocked dependencies
	HandlerGetAllOperations(logger, operationReader).ServeHTTP(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Decode the response body
	var response []*models.Operation
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response body
	assert.Equal(t, mockOperations, response)
}

func TestHandlerGetAllOperations_Error(t *testing.T) {
	// Mock logger
	opts := prettylogger.PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	}
	handler := prettylogger.NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)

	// Mock operation reader function to return an error
	operationReader := func() ([]*models.Operation, error) {
		return nil, errors.New("mock error")
	}

	// Create a request to pass to the handler
	req := httptest.NewRequest("GET", "/operations", nil)
	rec := httptest.NewRecorder()

	// Call the handler function with the mocked dependencies
	HandlerGetAllOperations(logger, operationReader).ServeHTTP(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	// Decode the response body
	var response models.Error
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Check the response body
	assert.Equal(t, models.Error(models.Error{Error: "no operations"}), response)
}
