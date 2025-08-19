package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type SeatRepository interface {
	GetByEventID(eventID int64) ([]models.Seat, error)
	GetByID(id int64) (*models.Seat, error)
	GetByIDForUpdate(id int64) (*models.Seat, error)
	GetByIDs(ids []int64) ([]models.Seat, error)
	UpdateStatus(seatID int64, status models.SeatStatus) error
	Update(seat *models.Seat) error
	ReserveSeats(seatIDs []int64, userID int) error
	ReleaseSeats(seatIDs []int64) error
	WithTx(tx *sql.Tx) SeatRepository
}

type seatRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewSeatRepository(db *sql.DB) SeatRepository {
	return &seatRepository{db: db}
}

func (r *seatRepository) WithTx(tx *sql.Tx) SeatRepository {
	return &seatRepository{db: r.db, tx: tx}
}

func (r *seatRepository) getExecutor() interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *seatRepository) GetByEventID(eventID int64) ([]models.Seat, error) {
	query := `
		SELECT id, event_id, row_number, seat_number, status, price, created_at, updated_at, version
		FROM seats WHERE event_id = $1 ORDER BY row_number, seat_number`

	executor := r.getExecutor()
	rows, err := executor.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to query seats: %w", err)
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(&seat.ID, &seat.EventID, &seat.RowNumber, &seat.SeatNumber,
			&seat.Status, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt, &seat.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat: %w", err)
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

func (r *seatRepository) GetByID(id int64) (*models.Seat, error) {
	query := `
		SELECT id, event_id, row_number, seat_number, status, price, created_at, updated_at, version
		FROM seats WHERE id = $1`

	var seat models.Seat
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&seat.ID, &seat.EventID, &seat.RowNumber,
		&seat.SeatNumber, &seat.Status, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt, &seat.Version)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get seat: %w", err)
	}

	return &seat, nil
}

func (r *seatRepository) GetByIDForUpdate(id int64) (*models.Seat, error) {
	query := `
		SELECT id, event_id, row_number, seat_number, status, price, created_at, updated_at, version
		FROM seats WHERE id = $1 FOR UPDATE`

	var seat models.Seat
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&seat.ID, &seat.EventID, &seat.RowNumber,
		&seat.SeatNumber, &seat.Status, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt, &seat.Version)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get seat for update: %w", err)
	}

	return &seat, nil
}

func (r *seatRepository) GetByIDs(ids []int64) ([]models.Seat, error) {
	if len(ids) == 0 {
		return []models.Seat{}, nil
	}

	query := `
		SELECT id, event_id, row_number, seat_number, status, price, created_at, updated_at, version
		FROM seats WHERE id = ANY($1)`

	executor := r.getExecutor()
	rows, err := executor.Query(query, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to query seats by IDs: %w", err)
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		err := rows.Scan(&seat.ID, &seat.EventID, &seat.RowNumber, &seat.SeatNumber,
			&seat.Status, &seat.Price, &seat.CreatedAt, &seat.UpdatedAt, &seat.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to scan seat: %w", err)
		}
		seats = append(seats, seat)
	}

	return seats, nil
}

func (r *seatRepository) UpdateStatus(seatID int64, status models.SeatStatus) error {
	query := `UPDATE seats SET status = $1, updated_at = $2 WHERE id = $3`

	executor := r.getExecutor()
	_, err := executor.Exec(query, status, time.Now(), seatID)
	if err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}

	return nil
}

func (r *seatRepository) Update(seat *models.Seat) error {
	query := `
		UPDATE seats 
		SET status = $1, updated_at = $2 
		WHERE id = $3`

	seat.UpdatedAt = time.Now()
	executor := r.getExecutor()
	_, err := executor.Exec(query, seat.Status, seat.UpdatedAt, seat.ID)
	if err != nil {
		return fmt.Errorf("failed to update seat: %w", err)
	}

	return nil
}

func (r *seatRepository) ReserveSeats(seatIDs []int64, userID int) error {
	if len(seatIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE seats 
		SET status = $1, updated_at = $2 
		WHERE id = ANY($3) AND status = $4`

	result, err := tx.Exec(query, models.SeatStatusReserved, time.Now(), seatIDs, models.SeatStatusFree)
	if err != nil {
		return fmt.Errorf("failed to reserve seats: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if int(rowsAffected) != len(seatIDs) {
		return fmt.Errorf("some seats are not available for reservation")
	}

	return tx.Commit()
}

func (r *seatRepository) ReleaseSeats(seatIDs []int64) error {
	if len(seatIDs) == 0 {
		return nil
	}

	query := `UPDATE seats SET status = $1, updated_at = $2 WHERE id = ANY($3)`

	executor := r.getExecutor()
	_, err := executor.Exec(query, models.SeatStatusFree, time.Now(), seatIDs)
	if err != nil {
		return fmt.Errorf("failed to release seats: %w", err)
	}

	return nil
}
