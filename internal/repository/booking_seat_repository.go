package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

type BookingSeatRepository interface {
	Create(booingSeat *models.BookingSeat) (*models.BookingSeat, error)
	GetByID(id int64) (*models.BookingSeat, error)
	GetBySeatID(seatID int64) ([]*models.BookingSeat, error)
	GetByBookingID(bookingID int64) ([]*models.BookingSeat, error)
	GetByBookingIDs(bookingIDs []int64) ([]*models.BookingSeat, error)
	Update(bookingSeat *models.BookingSeat) (*models.BookingSeat, error)
	Delete(id int64) error
	WithTx(txt *sql.Tx) BookingSeatRepository
}

type bookingSeatRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewBookingSeatRepository(db *sql.DB) BookingSeatRepository {
	return &bookingSeatRepository{db: db}
}

func (r *bookingSeatRepository) WithTx(tx *sql.Tx) BookingSeatRepository {
	return &bookingSeatRepository{db: r.db, tx: tx}
}

func (r *bookingSeatRepository) getExecutor() interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *bookingSeatRepository) Create(bookingSeat *models.BookingSeat) (*models.BookingSeat, error) {
	query := `
		INSERT INTO booking_seats(booking_id, seat_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	now := time.Now()
	bookingSeat.CreatedAt = now

	executor := r.getExecutor()
	err := executor.QueryRow(query, bookingSeat.BookingID, bookingSeat.SeatID, now).Scan(&bookingSeat.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to create bookingSeat: %w", err)
	}

	return bookingSeat, nil
}

func (r *bookingSeatRepository) GetByID(id int64) (*models.BookingSeat, error) {
	query := `
		SELECT id, booking_id, seat_id, created_at
		FROM booking_seats WHERE id = $1`

	var bookingSeat models.BookingSeat
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&bookingSeat.ID, &bookingSeat.BookingID, &bookingSeat.SeatID, &bookingSeat.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get bookingSeat: %w", err)
	}

	return &bookingSeat, nil
}

func (r *bookingSeatRepository) GetBySeatID(seatID int64) ([]*models.BookingSeat, error) {
	query := `
		SELECT id, booking_id, seat_id, created_at
		FROM booking_seats WHERE seat_id = $1`

	executor := r.getExecutor()
	rows, err := executor.Query(query, seatID)

	if err != nil {
		return nil, fmt.Errorf("failed to get bookingSeats by seatID: %w", err)
	}

	var bookingSeats []*models.BookingSeat
	for rows.Next() {
		var bookingSeat models.BookingSeat
		err := rows.Scan(&bookingSeat.ID, &bookingSeat.BookingID, &bookingSeat.SeatID, &bookingSeat.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan bookingSeat: %w", err)
		}
		bookingSeats = append(bookingSeats, &bookingSeat)
	}

	return bookingSeats, nil
}

func (r *bookingSeatRepository) GetByBookingID(bookingID int64) ([]*models.BookingSeat, error) {
	query := `
		SELECT id, booking_id, seat_id, created_at
		FROM booking_seats WHERE booking_id = $1`

	executor := r.getExecutor()
	rows, err := executor.Query(query, bookingID)

	if err != nil {
		return nil, fmt.Errorf("failed to get bookingSeats by bookingID: %w", err)
	}

	var bookingSeats []*models.BookingSeat
	for rows.Next() {
		var bookingSeat models.BookingSeat
		err = rows.Scan(&bookingSeat.ID, &bookingSeat.BookingID, &bookingSeat.SeatID, &bookingSeat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookingSeat: %w", err)
		}
		bookingSeats = append(bookingSeats, &bookingSeat)
	}

	return bookingSeats, nil
}

func (r *bookingSeatRepository) GetByBookingIDs(bookingIDs []int64) ([]*models.BookingSeat, error) {
	if len(bookingIDs) == 0 {
		return []*models.BookingSeat{}, nil
	}

	query := `
		SELECT id, booking_id, seat_id, created_at
		FROM booking_seats WHERE booking_id = ANY($1)`

	executor := r.getExecutor()
	rows, err := executor.Query(query, pq.Array(bookingIDs))

	if err != nil {
		return nil, fmt.Errorf("failed to get bookingSeats by bookingID: %w", err)
	}

	var bookingSeats []*models.BookingSeat
	for rows.Next() {
		var bookingSeat models.BookingSeat
		err = rows.Scan(&bookingSeat.ID, &bookingSeat.BookingID, &bookingSeat.SeatID, &bookingSeat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan bookingSeat: %w", err)
		}
		bookingSeats = append(bookingSeats, &bookingSeat)
	}

	return bookingSeats, nil
}

func (r *bookingSeatRepository) Update(bookingSeat *models.BookingSeat) (*models.BookingSeat, error) {
	query := `
		UPDATE booking_seats SET booking_id = $1, seat_id = $2
		WHERE id = $3`

	executor := r.getExecutor()
	err := executor.QueryRow(query, bookingSeat.ID, bookingSeat.SeatID, bookingSeat.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to update bookingSeat: %w", err)
	}

	return bookingSeat, nil
}

func (r *bookingSeatRepository) Delete(id int64) error {
	query := `DELETE FROM booking_seats WHERE id = $1`

	executor := r.getExecutor()
	_, err := executor.Exec(query, id)

	if err != nil {
		return fmt.Errorf("failed to delete bookingSeat: %w", err)
	}

	return nil
}
