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
	WithTx(tx *sql.Tx) EventRepository
}

type eventRepository struct {
	db *sql.DB
	tx *sql.Tx
}

func NewEventRepository(db *sql.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) WithTx(tx *sql.Tx) EventRepository {
	return &eventRepository{db: r.db, tx: tx}
}

func (r *eventRepository) getExecutor() interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
} {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *eventRepository) FindEvents(query *string, date *time.Time, page, pageSize int) ([]models.Event, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	baseQuery := `SELECT id, title, description, type, datetime_start, provider FROM events`

	// Используем полнотекстовый поиск с индексом gin для лучшей производительности
	if query != nil && *query != "" {
		// Сначала пробуем полнотекстовый поиск
		//conditions = append(conditions, fmt.Sprintf("to_tsvector('russian', title || ' ' || description) @@ plainto_tsquery('russian', $%d)", argIndex))
		//args = append(args, *query)
		//argIndex++

		// Добавляем резервный ILIKE поиск для случаев, когда полнотекстовый поиск не даст результатов
		// Используем lower() функцию для оптимизации с индексом
		//conditions[len(conditions)-1] += fmt.Sprintf("(lower(title) LIKE lower($%d) OR lower(description) LIKE lower($%d))", argIndex, argIndex)
		conditions = append(conditions, fmt.Sprintf("(lower(title) LIKE lower($%d) OR lower(description) LIKE lower($%d))", argIndex, argIndex))
		args = append(args, "%"+*query+"%")
		argIndex++
	}

	// Оптимизированный поиск по дате с использованием индекса
	if date != nil {
		conditions = append(conditions, fmt.Sprintf("DATE(datetime_start) = $%d", argIndex))
		args = append(args, date.Format("2006-01-02"))
		argIndex++
	}

	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY id"

	// Добавляем пагинацию с защитой от больших offset'ов
	offset := (page - 1) * pageSize
	if offset > 10000 { // Защита от глубокой пагинации
		offset = 10000
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, pageSize, offset)

	//executor := r.getExecutor() // в целях оптимизации убрал вызов через executor
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
	executor := r.getExecutor()
	err := executor.QueryRow(query, id).Scan(&event.ID, &event.Title, &event.Description,
		&event.Type, &event.DatetimeStart, &event.Provider)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}

	return &event, nil
}
