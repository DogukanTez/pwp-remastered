package store

import (
	"pwp-remastered/internal/database"
	"pwp-remastered/internal/domain"
	"time"
)

// EventStore handles event data operations
type EventStore interface {
	GetEvent(int) (*domain.Event, error)
	CreateEvent(*domain.Event) error
	UpdateEvent(*domain.Event) error
	DeleteEvent(int) error
	GetDatedUserEvents(int, time.Time, time.Time) ([]domain.Event, error)
}

// Example implementation using database layer
type eventDBStore struct {
	db database.Service // Assuming dependency injection
}

func NewEventStore(db database.Service) EventStore {
	return &eventDBStore{db: db}
}

func (s *eventDBStore) GetEvent(id int) (*domain.Event, error) {
	var event domain.Event
	var user domain.EventUser
	var eventType domain.EventType

	query := `
		SELECT 
			e.id, e.type_id, e.user_id, e.name, e.title, e.description, 
			e.start_date, e.end_date, e.road_price,
			u.id, u.username, u.first_name, u.last_name,
			et.id,et.type, et.language, et.color, et.is_pricable
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		LEFT JOIN event_types et ON e.type_id = et.id
		WHERE e.id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&event.ID, &event.TypeID, &event.UserID, &event.Name, &event.Title, &event.Description,
		&event.StartDate, &event.EndDate, &event.RoadPrice,
		&user.ID, &user.Username, &user.FirstName, &user.LastName,
		&eventType.ID, &eventType.Type, &eventType.Language, &eventType.Color, &eventType.IsPricable,
	)
	if err != nil {
		return nil, err
	}

	event.User = &user
	event.Type = &eventType
	return &event, nil
}

func (s *eventDBStore) CreateEvent(event *domain.Event) error {
	// Implementation here
	return nil
}

func (s *eventDBStore) UpdateEvent(event *domain.Event) error {
	// Implementation here
	return nil
}

func (s *eventDBStore) DeleteEvent(id int) error {
	// Implementation here
	return nil
}

func (s *eventDBStore) GetDatedUserEvents(id int, startdate time.Time, enddate time.Time) ([]domain.Event, error) {
	var events []domain.Event

	query := `
		SELECT
			e.id, e.type_id, e.user_id, e.name, e.title, e.description,
			e.start_date, e.end_date, e.road_price,
			u.id, u.username, u.first_name, u.last_name,
			et.id,et.type, et.language, et.color, et.is_pricable
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		LEFT JOIN event_types et ON e.type_id = et.id
		WHERE e.user_id = $1 AND e.start_date >= $2 AND e.end_date <= $3
		ORDER BY e.start_date`

	rows, err := s.db.Query(query, id, startdate, enddate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event domain.Event
		var user domain.EventUser
		var eventType domain.EventType

		err := rows.Scan(
			&event.ID, &event.TypeID, &event.UserID, &event.Name, &event.Title, &event.Description,
			&event.StartDate, &event.EndDate, &event.RoadPrice,
			&user.ID, &user.Username, &user.FirstName, &user.LastName,
			&eventType.ID, &eventType.Type, &eventType.Language, &eventType.Color, &eventType.IsPricable,
		)
		if err != nil {
			return nil, err
		}

		event.User = &user
		event.Type = &eventType
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil

}
