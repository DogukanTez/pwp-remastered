package domain

import (
	"time"
)

type Event struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	RoadPrice   float64   `json:"road_price"`
}

type EventList struct {
	Events []Event `json:"events"`
}
