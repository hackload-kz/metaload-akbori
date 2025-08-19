package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
)

type UserRepository interface {
	GetByID(userID int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByID(userID int) (*models.User, error) {
	query := `
		SELECT user_id, email, password_hash, password_plain, first_name, surname, 
		birthday, registered_at, is_active, last_logged_in
		FROM users WHERE user_id = $1`
	
	var user models.User
	err := r.db.QueryRow(query, userID).Scan(&user.UserID, &user.Email, &user.PasswordHash,
		&user.PasswordPlain, &user.FirstName, &user.Surname, &user.Birthday,
		&user.RegisteredAt, &user.IsActive, &user.LastLoggedIn)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT user_id, email, password_hash, password_plain, first_name, surname,
		birthday, registered_at, is_active, last_logged_in
		FROM users WHERE email = $1`
	
	var user models.User
	err := r.db.QueryRow(query, email).Scan(&user.UserID, &user.Email, &user.PasswordHash,
		&user.PasswordPlain, &user.FirstName, &user.Surname, &user.Birthday,
		&user.RegisteredAt, &user.IsActive, &user.LastLoggedIn)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &user, nil
}