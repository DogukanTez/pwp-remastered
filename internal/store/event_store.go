package store

import (
	"database/sql"
	"errors"
	"pwp-remastered/internal/database"
	"pwp-remastered/internal/domain"

	"time"
)

// EventStore handles event data operations
type EventStore interface {
	GetEvent(int) (*domain.Event, error)
	CreateEvent(*domain.Event, *domain.User) error
	UpdateEvent(*domain.Event, *domain.User) error
	DeleteEvent(int) error
	GetDatedUserEvents(int, time.Time, time.Time) ([]domain.Event, error)
	GetAllDatedEvents(time.Time, time.Time) ([]domain.Event, error)
	GetSelfDatedEvents(*domain.User, time.Time, time.Time) ([]domain.Event, error)
	GetEventType(int) (*domain.EventType, error)
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

func (s *eventDBStore) CreateEvent(event *domain.Event, caller *domain.User) error {
	userID := caller.ID
	typeID := event.TypeID

	EventType, err := s.GetEventType(typeID)

	if err != nil {
		return err
	}

	event.Type = EventType
	event.UserID = userID

	query := `
		INSERT INTO events (type_id, user_id, name, title, description, start_date, end_date, road_price)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	err = s.db.QueryRow(query, typeID, userID, event.Name, event.Title, event.Description, event.StartDate, event.EndDate, event.RoadPrice).Scan(&event.ID)
	if err != nil {
		return err
	}

	userQuery := `
		SELECT u.id, u.username, u.first_name, u.last_name
		FROM users u
		WHERE u.id = $1`
	var user domain.EventUser
	err = s.db.QueryRow(userQuery, userID).Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName)
	if err != nil {
		return err
	}
	event.User = &user

	return nil
}

func (s *eventDBStore) UpdateEvent(event *domain.Event, caller *domain.User) error {
	//TODO: type cannot be manually changed add it to query and remove user_id

	query := `
		UPDATE events
		SET type_id = $1, user_id = $2, name = $3, title = $4, description = $5, start_date = $6, end_date = $7, road_price = $8
		WHERE id = $9
		RETURNING id`

	_, err := s.db.Exec(query, event.TypeID, event.UserID, event.Name, event.Title, event.Description, event.StartDate, event.EndDate, event.RoadPrice, event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *eventDBStore) DeleteEvent(id int) error {
	query := `
		DELETE FROM events
		WHERE id = $1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
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

func (s *eventDBStore) GetAllDatedEvents(startdate time.Time, enddate time.Time) ([]domain.Event, error) {
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
		WHERE e.start_date >= $1 AND e.end_date <= $2
		ORDER BY e.start_date`

	rows, err := s.db.Query(query, startdate, enddate)
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

func (s *eventDBStore) GetSelfDatedEvents(caller *domain.User, startdate time.Time, enddate time.Time) ([]domain.Event, error) {

	events, err := s.GetDatedUserEvents(caller.ID, startdate, enddate)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *eventDBStore) GetEventType(id int) (*domain.EventType, error) {
	var eventType domain.EventType
	query := `
		SELECT id, type, language, color, is_pricable
		FROM event_types
		WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&eventType.ID, &eventType.Type, &eventType.Language, &eventType.Color, &eventType.IsPricable,
	)

	if err == sql.ErrNoRows {
		var ErrEventTypeNotFound = errors.New("event type not found")
		return nil, ErrEventTypeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &eventType, nil
}
