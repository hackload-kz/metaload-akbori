package models

import (
	"time"
)

type User struct {
	UserID       int       `json:"user_id" db:"user_id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	PasswordPlain *string  `json:"-" db:"password_plain"`
	FirstName    string    `json:"first_name" db:"first_name"`
	Surname      string    `json:"surname" db:"surname"`
	Birthday     *time.Time `json:"birthday" db:"birthday"`
	RegisteredAt time.Time `json:"registered_at" db:"registered_at"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	LastLoggedIn time.Time `json:"last_logged_in" db:"last_logged_in"`
}