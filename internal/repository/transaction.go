package repository

import (
	"database/sql"
	"fmt"
)

// Transaction wraps a database transaction
type Transaction struct {
	tx *sql.Tx
}

// TransactionRepository provides repository operations within a transaction
type TransactionRepository struct {
	tx      *sql.Tx
	Event   EventRepository
	Booking BookingRepository
	Seat    SeatRepository
	User    UserRepository
}

// TransactionFunc is a function that executes within a transaction
type TransactionFunc func(repo *TransactionRepository) error

// TransactionManager provides transaction management
type TransactionManager struct {
	db *sql.DB
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(db *sql.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
// If the function returns an error, the transaction is rolled back
// Otherwise, the transaction is committed
func (tm *TransactionManager) WithTransaction(fn TransactionFunc) error {
	tx, err := tm.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create repositories that use the transaction
	txRepo := &TransactionRepository{
		tx:      tx,
		Event:   NewEventRepository(tm.db).WithTx(tx),
		Booking: NewBookingRepository(tm.db).WithTx(tx),
		Seat:    NewSeatRepository(tm.db).WithTx(tx),
		User:    NewUserRepository(tm.db).WithTx(tx),
	}

	// Execute the function
	err = fn(txRepo)
	if err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("transaction failed: %w, rollback failed: %v", err, rollbackErr)
		}
		return err
	}

	// Commit if no errors
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// TransactionalRepository interface that all repositories should implement
// to support both regular database connections and transactions
type TransactionalRepository interface {
	WithTx(tx *sql.Tx) TransactionalRepository
}
