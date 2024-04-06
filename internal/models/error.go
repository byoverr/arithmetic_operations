package models

type Error struct {
	Error string `json:"error"`
}

func NewError(s string) *Error {
	return &Error{Error: s}
}
