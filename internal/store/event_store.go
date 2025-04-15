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
	// Implementation here
	return &domain.Event{}, nil
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
