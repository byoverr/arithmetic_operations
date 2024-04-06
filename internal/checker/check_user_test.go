package checker

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsValidUsername(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid username",
			exp:  "serejkaaa",
			want: nil,
		},
		{
			name: "Invalid username",
			exp:  "ser",
			want: errors.New("expression has division by zero"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidUsername(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("isValidUsername() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestIsValidPassword(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid password",
			exp:  "Qwerty_123",
			want: nil,
		},
		{
			name: "Missing upper case character",
			exp:  "qwerty_123",
			want: fmt.Errorf("password must have at least one upper case character"),
		},
		{
			name: "Missing lower case character",
			exp:  "QWERTY_123",
			want: fmt.Errorf("password must have at least one lower case character"),
		},
		{
			name: "Missing numeric character",
			exp:  "qwertY_",
			want: fmt.Errorf("password must have at least one numeric character"),
		},
		{
			name: "Missing special character",
			exp:  "qwertY123",
			want: fmt.Errorf("password must have at least one special character"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidPassword(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("HasValidCharacters() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}
