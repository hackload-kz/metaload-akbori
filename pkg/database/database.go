package database

import (
	"biletter-service/internal/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(cfg config.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(0) // maximum amount of time a connection may be reused
	db.SetConnMaxIdleTime(0) // maximum amount of time a connection may be idle before being closed
	return db, nil
}
