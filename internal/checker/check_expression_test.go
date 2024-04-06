package checker

import (
	"errors"
	"fmt"
	"testing"
)

func TestRemoveAllSpaces(t *testing.T) {
	str := RemoveAllSpaces("2 + 2")
	if str != "2+2" {
		t.Errorf("RemoveAllSpaces() error")
	}
}

func TestHasDoubleSymbol(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid expression",
			exp:  "1+2",
			want: nil,
		},
		{
			name: "Expression has double symbol",
			exp:  "4**2",
			want: errors.New("expression has doubled symbol"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HasDoubleSymbol(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("HasDoubleSymbol() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestIsValidParentheses(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Correct brackets",
			exp:  "(1+2)",
			want: nil,
		},
		{
			name: "Missing opening bracket",
			exp:  "1+2)",
			want: errors.New("expression has invalid parentheses"),
		},
		{
			name: "Missing closing bracket",
			exp:  "(1+2",
			want: errors.New("expression has invalid parentheses"),
		},
		{
			name: "Mismatched brackets",
			exp:  "(1+2]",
			want: errors.New("expression has invalid parentheses"),
		},
		{
			name: "Invalid brackets",
			exp:  "(1+2)())(",
			want: errors.New("expression has invalid parentheses"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidParentheses(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("IsValidParentheses() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestHasDivizionByZero(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid expression",
			exp:  "1+2",
			want: nil,
		},
		{
			name: "Expression with division by zero",
			exp:  "4/0",
			want: errors.New("expression has division by zero"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HasDivizionByZero(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("HasDivizionByZero() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestHasValidCharacters(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid expression",
			exp:  "(1+2)",
			want: nil,
		},
		{
			name: "Expression with extra characters",
			exp:  "(1+2)!",
			want: errors.New(fmt.Sprintf("expression has invalid character !")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HasValidCharacters(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("HasValidCharacters() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestContainsCorrectFloatPoint(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid expression",
			exp:  "(2.4+3.4)",
			want: nil,
		},
		{
			name: "Incorrect float point in end",
			exp:  "34+14.",
			want: errors.New("expression has a dot in a wrong place"),
		},
		{
			name: "Incorrect float point in center",
			exp:  "34.+14",
			want: errors.New("expression has a dot in a wrong place"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ContainsCorrectFloatPoint(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("ContainsCorrectFloatPoint() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestHasAtLeastOneExpression(t *testing.T) {
	tests := []struct {
		name string
		exp  string
		want error
	}{
		{
			name: "Valid expression",
			exp:  "1+3",
			want: nil,
		},
		{
			name: "Only numbers",
			exp:  "1234",
			want: errors.New("this string doesn't has at least one expression"),
		},
		{
			name: "Only chars",
			exp:  "invalid expression",
			want: errors.New("this string doesn't has at least one expression"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HasAtLeastOneExpression(tt.exp)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("HasAtLeastOneExpression() error = %v, wantErr %v", err, tt.want)
			}
		})
	}
}

func TestExpressionStartsWithNumber(t *testing.T) {
	testCases := []struct {
		name     string
		exp      string
		expected error
	}{
		{
			name:     "Starts with number",
			exp:      "1+2",
			expected: nil,
		},
		{
			name:     "Does not start with number",
			exp:      "a+2",
			expected: errors.New("this string doesn't start with number or bracket"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ExpressionStartsWithNumber(tc.exp)
			if (got == nil) != (tc.expected == nil) {
				t.Errorf("ExpressionStartsWithNumber() = %v, want %v", got, tc.expected)
			}
		})
	}
}
