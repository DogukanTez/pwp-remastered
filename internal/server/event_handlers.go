package server

import (
	"encoding/json"
	"net/http"
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/services"
	"pwp-remastered/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandlers struct {
	eventService services.EventService
	eventStore   store.EventStore
}

// NewEventHandlers creates a new event handlers
func NewEventHandlers(eventService services.EventService, eventStore store.EventStore) *EventHandlers {
	return &EventHandlers{
		eventService: eventService,
		eventStore:   eventStore,
	}
}

func (h *EventHandlers) RegisterRoutes(r chi.Router) {
	r.Route("/events", func(r chi.Router) {
		r.Get("/{id}", h.GetEvent)
		r.Post("/", h.CreateEvent)
		r.Put("/{id}", h.UpdateEvent)
		r.Delete("/{id}", h.DeleteEvent)
	})
}

// GetEvent returns a single event by ID
func (h *EventHandlers) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := h.eventStore.GetEvent(eventID)
	if err != nil {
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// CreateEvent creates a new event
func (h *EventHandlers) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.eventStore.CreateEvent(&event); err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// UpdateEvent updates an existing event
func (h *EventHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	event.ID = eventID
	if err := h.eventStore.UpdateEvent(&event); err != nil {
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// DeleteEvent deletes an event by ID
func (h *EventHandlers) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid event ID", http.StatusBadRequest)
		return
	}

	if err := h.eventStore.DeleteEvent(eventID); err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
