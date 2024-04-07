package storage

import (
	"arithmetic_operations/internal/models"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:admin@localhost:5432/arithmetic_operations?sslmode=disable") // Replace with your str
	if err != nil {
		t.Fatalf("an error %s", err)
	}
	r := &PostgresqlDB{db: db}
	tests := []struct {
		name    string
		input   *models.User
		want    int64
		wantErr bool
	}{
		{
			name: "Ok",
			input: &models.User{
				Username:     "tester_user",
				HashPassword: "Qwerty_123",
				CreatedAt:    time.Now(),
			},
			want:    int64(1), // Replace with needed id
			wantErr: false,
		},
		{
			name: "Same User",
			input: &models.User{
				Username:     "tester_user",
				HashPassword: "Qwerty_123",
				CreatedAt:    time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.CreateUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	db, errDb := pgxpool.New(context.Background(), "postgres://postgres:admin@localhost:5432/arithmetic_operations?sslmode=disable") // Replace with your str
	if errDb != nil {
		t.Fatalf("an error %s", errDb)
	}
	r := &PostgresqlDB{db: db}

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Ok",
			input:   "tester_user",
			wantErr: false,
		},
		{
			name:    "Unknown User",
			input:   "some_user",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.GetUser(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got.Username, tt.input)
			}
		})
	}
}
