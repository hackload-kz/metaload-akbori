package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
	"sync"
)

type UserRepository interface {
	GetByID(userID int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	PreloadCache() error
	WithTx(tx *sql.Tx) UserRepository
}

// UserCache provides thread-safe caching for users
type UserCache struct {
	mu           sync.RWMutex
	usersByID    map[int]*models.User
	usersByEmail map[string]*models.User
}

func NewUserCache() *UserCache {
	return &UserCache{
		usersByID:    make(map[int]*models.User),
		usersByEmail: make(map[string]*models.User),
	}
}

func (c *UserCache) GetByID(userID int) *models.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.usersByID[userID]
}

func (c *UserCache) GetByEmail(email string) *models.User {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.usersByEmail[email]
}

func (c *UserCache) Set(user *models.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.usersByID[user.UserID] = user
	c.usersByEmail[user.Email] = user
}

type userRepository struct {
	db    *sql.DB
	tx    *sql.Tx
	cache *UserCache
}

var globalUserCache = NewUserCache()

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db:    db,
		cache: globalUserCache,
	}
}

func (r *userRepository) WithTx(tx *sql.Tx) UserRepository {
	return &userRepository{db: r.db, tx: tx, cache: r.cache}
}

func (r *userRepository) getExecutor() interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *userRepository) GetByID(userID int) (*models.User, error) {
	// Проверяем кэш сначала
	if cachedUser := r.cache.GetByID(userID); cachedUser != nil {
		return cachedUser, nil
	}

	query := `
		SELECT user_id, email, password_hash, password_plain, first_name, surname, 
		birthday, registered_at, is_active, last_logged_in
		FROM users WHERE user_id = $1`

	var user models.User
	executor := r.getExecutor()
	err := executor.QueryRow(query, userID).Scan(&user.UserID, &user.Email, &user.PasswordHash,
		&user.PasswordPlain, &user.FirstName, &user.Surname, &user.Birthday,
		&user.RegisteredAt, &user.IsActive, &user.LastLoggedIn)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Сохраняем в кэш
	r.cache.Set(&user)

	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	// Проверяем кэш сначала
	if cachedUser := r.cache.GetByEmail(email); cachedUser != nil {
		return cachedUser, nil
	}

	query := `
		SELECT user_id, email, password_hash, password_plain, first_name, surname,
		birthday, registered_at, is_active, last_logged_in
		FROM users WHERE email = $1`

	var user models.User
	executor := r.getExecutor()
	err := executor.QueryRow(query, email).Scan(&user.UserID, &user.Email, &user.PasswordHash,
		&user.PasswordPlain, &user.FirstName, &user.Surname, &user.Birthday,
		&user.RegisteredAt, &user.IsActive, &user.LastLoggedIn)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	// Сохраняем в кэш
	r.cache.Set(&user)

	return &user, nil
}

// PreloadCache загружает всех пользователей в кэш при старте приложения
func (r *userRepository) PreloadCache() error {
	query := `
		SELECT user_id, email, password_hash, password_plain, first_name, surname,
		birthday, registered_at, is_active, last_logged_in
		FROM users WHERE is_active = true`

	rows, err := r.db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to preload users cache: %w", err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.UserID, &user.Email, &user.PasswordHash,
			&user.PasswordPlain, &user.FirstName, &user.Surname, &user.Birthday,
			&user.RegisteredAt, &user.IsActive, &user.LastLoggedIn)
		if err != nil {
			return fmt.Errorf("failed to scan user during cache preload: %w", err)
		}

		r.cache.Set(&user)
		count++
	}

	fmt.Printf("✅ Preloaded %d users into cache\n", count)
	return nil
}
