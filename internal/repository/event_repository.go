package repository

import (
	"biletter-service/internal/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type EventRepository interface {
	FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.Event, error)
	GetByID(id int64) (*models.Event, error)
}

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.Event, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `SELECT id, title, description, type, datetime_start, provider FROM events`
	
	if query != nil && *query != "" {
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+*query+"%")
		argIndex++
	}

	if date != nil {
		conditions = append(conditions, fmt.Sprintf("DATE(datetime_start) = $%d", argIndex))
		args = append(args, date.Format("2006-01-02"))
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY datetime_start"
	
	offset := (page - 1) * pageSize
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		err := rows.Scan(&event.ID, &event.Title, &event.Description, 
			&event.Type, &event.DatetimeStart, &event.Provider)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}

func (r *eventRepository) GetByID(id int64) (*models.Event, error) {
	query := `SELECT id, title, description, type, datetime_start, provider FROM events WHERE id = $1`
	
	var event models.Event
	err := r.db.QueryRow(query, id).Scan(&event.ID, &event.Title, &event.Description,
		&event.Type, &event.DatetimeStart, &event.Provider)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}

	return &event, nil
}