package domain

import (
	"time"
)

type Event struct {
	ID     int `json:"id"`
	TypeID int `json:"type_id"`
	UserID int `json:"user_id"`
	// Username    string    `json:"username"`
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	RoadPrice   float64    `json:"road_price"`
	User        *EventUser `json:"user,omitempty"`
	Type        *EventType `json:"type,omitempty"`
}

type EventList struct {
	Events []Event `json:"events"`
}

type EventType struct {
	ID         int     `json:"id"`
	Type       string  `json:"type"`
	Language   string  `json:"language"`
	Color      *string `json:"color"`
	IsPricable bool    `json:"is_pricable"`
}

type EventTypeList struct {
	EventTypes []EventType `json:"event_types"`
}

type EventUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
