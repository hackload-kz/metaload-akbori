package models

import (
	"time"
)

type Event struct {
	ID           int64     `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Type         string    `json:"type" db:"type"`
	DatetimeStart time.Time `json:"datetime_start" db:"datetime_start"`
	Provider     string    `json:"provider" db:"provider"`
}