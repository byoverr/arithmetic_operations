package models

import "time"

type User struct {
	Id           int       `json:"-" db:"id"`
	Username     string    `json:"username" validate:"username, required"`
	HashPassword string    `json:"password" validate:"password, required"`
	CreatedAt    time.Time `json:"createdAt"`
}

func NewUser(username, password string) *User {
	return &User{
		Username:     username,
		HashPassword: password,
		CreatedAt:    time.Now(),
	}
}

type RegisterUser struct {
	Username string `json:"username" validate:"username, required"`
	Password string `json:"password" validate:"password, required"`
}
