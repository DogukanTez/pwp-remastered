package services

import (
	"errors"
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/store"
	"time"
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
func (s *EventService) CreateEvent(event *domain.Event, caller *domain.User) error {
	return s.store.CreateEvent(event, caller)
}

// UpdateEvent modifies an existing event
func (s *EventService) UpdateEvent(event *domain.Event, caller *domain.User) error {
	userID := caller.ID

	if event.UserID != userID {
		return errors.New("Caller is not the owner of the event")
	}

	return s.store.UpdateEvent(event, caller)
}

// DeleteEvent removes an event by ID
func (s *EventService) DeleteEvent(id int) error {

	return s.store.DeleteEvent(id)
}

// GetDatedUserEvents retrieves events for a user within a date range
func (s *EventService) GetDatedUserEvents(caller *domain.User, userID int, startDate time.Time, endDate time.Time) ([]domain.Event, error) {
	if !caller.IsAdmin {
		return nil, errors.New("Caller is not admin")
		// return s.store.GetSelfDatedEvents(caller, startDate, endDate)
	}
	return s.store.GetDatedUserEvents(userID, startDate, endDate)
}

// GetAllDatedEvents retrieves all events within a date range
func (s *EventService) GetAllDatedEvents(caller *domain.User, startDate time.Time, endDate time.Time) ([]domain.Event, error) {
	if !caller.IsAdmin {
		return nil, errors.New("Caller is not authorized")
	}
	return s.store.GetAllDatedEvents(startDate, endDate)
}

func (s *EventService) GetSelfDatedEvents(caller *domain.User, startDate time.Time, endDate time.Time) ([]domain.Event, error) {
	return s.store.GetSelfDatedEvents(caller, startDate, endDate)
}

func (s *EventService) GetEventTypes() ([]domain.EventType, error) {
	return s.store.GetEventTypes()
}
