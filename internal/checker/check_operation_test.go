package checker

import (
	"arithmetic_operations/internal/models"
	"errors"
	"testing"
)

func TestValidateOperation(t *testing.T) {
	tests := []struct {
		name string
		exp  models.Operation
		want error
	}{
		{
			name: "Valid operation",
			exp: models.Operation{
				OperationKind:         "addition",
				DurationInMilliSecond: 200,
			},
			want: nil,
		},
		{
			name: "Operation below zero",
			exp: models.Operation{
				OperationKind:         "addition",
				DurationInMilliSecond: -200,
			},
			want: errors.New("operation duration should be more than 0"),
		},
		{
			name: "Operation not allowed",
			exp: models.Operation{
				OperationKind:         "add",
				DurationInMilliSecond: 200,
			},
			want: errors.New("it is not allowed operation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOperation(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("ValidateOperation() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
