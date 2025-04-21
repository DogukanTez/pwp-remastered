package store

import (
	"pwp-remastered/internal/database"
	"pwp-remastered/internal/domain"
)

// EventStore handles event data operations
type EventStore interface {
	GetEvent(int) (*domain.Event, error)
	CreateEvent(*domain.Event) error
	UpdateEvent(*domain.Event) error
	DeleteEvent(int) error
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
	var user domain.User
	var eventType domain.EventType

	query := `
		SELECT 
			e.id, e.type_id, e.user_id, e.name, e.title, e.description, 
			e.start_date, e.end_date, e.road_price,
			u.username, u.first_name, u.last_name,
			et.type, et.language, et.color, et.is_pricable
		FROM events e
		LEFT JOIN users u ON e.user_id = u.id
		LEFT JOIN event_types et ON e.type_id = et.id
		WHERE e.id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&event.ID, &event.TypeID, &event.UserID, &event.Name, &event.Title, &event.Description,
		&event.StartDate, &event.EndDate, &event.RoadPrice,
		&user.Username, &user.FirstName, &user.LastName,
		&eventType.Type, &eventType.Language, &eventType.Color, &eventType.IsPricable,
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
