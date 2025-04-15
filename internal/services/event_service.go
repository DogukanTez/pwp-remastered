package services

import (
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/store"
)

// EventService handles business logic for events
type EventService struct {
	store store.EventStore
}

// NewEventService creates a new event service
func NewEventService(eventStore store.EventStore) *EventService {
	return &EventService{
		store: eventStore,
	}
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(id int) (*domain.Event, error) {
	return s.store.GetEvent(id)
}

// CreateEvent persists a new event
func (s *EventService) CreateEvent(event *domain.Event) error {
	return s.store.CreateEvent(event)
}

// UpdateEvent modifies an existing event
func (s *EventService) UpdateEvent(event *domain.Event) error {
	return s.store.UpdateEvent(event)
}

// DeleteEvent removes an event by ID
func (s *EventService) DeleteEvent(id int) error {
	return s.store.DeleteEvent(id)
}
