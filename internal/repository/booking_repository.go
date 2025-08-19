package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type BookingRepository interface {
	Create(booking *models.Booking) (*models.Booking, error)
	GetByID(id int64) (*models.Booking, error)
	GetByIDForUpdate(id int64) (*models.Booking, error)
	Update(booking *models.Booking) error
	GetByUserID(userID int) ([]models.Booking, error)
	GetAll() ([]models.Booking, error)
	GetByOrderID(orderID string) (*models.Booking, error)
	WithTx(tx *sql.Tx) BookingRepository
}

type bookingRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewBookingRepository(db *sql.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) WithTx(tx *sql.Tx) BookingRepository {
	return &bookingRepository{db: r.db, tx: tx}
}

func (r *bookingRepository) getExecutor() interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *bookingRepository) Create(booking *models.Booking) (*models.Booking, error) {
	query := `
		INSERT INTO bookings (event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	now := time.Now()
	booking.CreatedAt = now
	booking.UpdatedAt = now

	executor := r.getExecutor()
	err := executor.QueryRow(query, booking.EventID, booking.UserID, booking.Status,
		booking.TotalAmount, booking.PaymentID, booking.OrderID, booking.CreatedAt, booking.UpdatedAt).Scan(&booking.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	return booking, nil
}

func (r *bookingRepository) GetByID(id int64) (*models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at
		FROM bookings WHERE id = $1`

	var booking models.Booking
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&booking.ID, &booking.EventID, &booking.UserID,
		&booking.Status, &booking.TotalAmount, &booking.PaymentID, &booking.OrderID,
		&booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &booking, nil
}

func (r *bookingRepository) GetByIDForUpdate(id int64) (*models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at
		FROM bookings WHERE id = $1 FOR UPDATE`

	var booking models.Booking
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&booking.ID, &booking.EventID, &booking.UserID,
		&booking.Status, &booking.TotalAmount, &booking.PaymentID, &booking.OrderID,
		&booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking for update: %w", err)
	}

	return &booking, nil
}

func (r *bookingRepository) Update(booking *models.Booking) error {
	query := `
		UPDATE bookings 
		SET status = $1, total_amount = $2, payment_id = $3, order_id = $4, updated_at = $5
		WHERE id = $6`

	booking.UpdatedAt = time.Now()

	executor := r.getExecutor()
	_, err := executor.Exec(query, booking.Status, booking.TotalAmount, booking.PaymentID,
		booking.OrderID, booking.UpdatedAt, booking.ID)

	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	return nil
}

func (r *bookingRepository) GetByUserID(userID int) ([]models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at
		FROM bookings WHERE user_id = $1 ORDER BY created_at DESC`

	executor := r.getExecutor()
	rows, err := executor.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings by user: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.EventID, &booking.UserID,
			&booking.Status, &booking.TotalAmount, &booking.PaymentID,
			&booking.OrderID, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *bookingRepository) GetAll() ([]models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at
		FROM bookings ORDER BY created_at DESC`

	executor := r.getExecutor()
	rows, err := executor.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all bookings: %w", err)
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		err := rows.Scan(&booking.ID, &booking.EventID, &booking.UserID,
			&booking.Status, &booking.TotalAmount, &booking.PaymentID,
			&booking.OrderID, &booking.CreatedAt, &booking.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *bookingRepository) GetByOrderID(orderID string) (*models.Booking, error) {
	query := `
		SELECT id, event_id, user_id, status, total_amount, payment_id, order_id, created_at, updated_at
		FROM bookings WHERE order_id = $1`

	var booking models.Booking
	executor := r.getExecutor()
	err := executor.QueryRow(query, orderID).Scan(&booking.ID, &booking.EventID, &booking.UserID,
		&booking.Status, &booking.TotalAmount, &booking.PaymentID, &booking.OrderID,
		&booking.CreatedAt, &booking.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking by order ID: %w", err)
	}

	return &booking, nil
}
