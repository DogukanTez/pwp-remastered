package server

import (
	"encoding/json"
	"net/http"
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/services"
	"pwp-remastered/internal/store"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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
		r.Use(AuthMiddleware)
		r.Get("/dated/{id}", h.GetDatedUserEvents)
		r.Get("/dated", h.GetAllDatedEvents)
		r.Get("/dated/me", h.GetSelfDatedEvents)
		r.Get("/{id}", h.GetEvent)
		r.Post("/", h.CreateEvent)
		r.Put("/{id}", h.UpdateEvent)
		r.Delete("/{id}", h.DeleteEvent)
		r.Get("/types", h.GetEventTypes)
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
	var caller domain.User
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString != "" {
		token, err := ParseJWT(tokenString)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if idVal, ok := claims["user_id"].(float64); ok {
					caller.ID = int(idVal)
				}
				if isAdmin, ok := claims["is_admin"].(bool); ok {
					caller.IsAdmin = isAdmin
				}
			}
		}
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.eventService.CreateEvent(&event, &caller); err != nil {
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(event)
}

// UpdateEvent updates an existing event
func (h *EventHandlers) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var caller domain.User
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString != "" {
		token, err := ParseJWT(tokenString)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if idVal, ok := claims["user_id"].(float64); ok {
					caller.ID = int(idVal)
				}
				if isAdmin, ok := claims["is_admin"].(bool); ok {
					caller.IsAdmin = isAdmin
				}
			}
		}
	}

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
	if err := h.eventService.UpdateEvent(&event, &caller); err != nil {
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

	if err := h.eventService.DeleteEvent(eventID); err != nil {
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandlers) GetDatedUserEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	startDateStr := r.URL.Query().Get("startdate")
	endDateStr := r.URL.Query().Get("enddate")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing startdate or enddate", http.StatusBadRequest)
		return
	}

	// ISO 8601 formatını parse et (örnek: "2025-03-01")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid startdate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid enddate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var caller domain.User
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString != "" {
		token, err := ParseJWT(tokenString)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if idVal, ok := claims["user_id"].(float64); ok {
					caller.ID = int(idVal)
				}
				if isAdmin, ok := claims["is_admin"].(bool); ok {
					caller.IsAdmin = isAdmin
				}
			}
		}
	}

	events, err := h.eventService.GetDatedUserEvents(&caller, userID, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandlers) GetAllDatedEvents(w http.ResponseWriter, r *http.Request) {
	caller, err := ExtractUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	startDateStr := r.URL.Query().Get("startdate")
	endDateStr := r.URL.Query().Get("enddate")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing startdate or enddate", http.StatusBadRequest)
		return
	}

	// ISO 8601 formatını parse et (örnek: "2025-03-01")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid startdate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid enddate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	events, err := h.eventService.GetAllDatedEvents(&caller, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandlers) GetSelfDatedEvents(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startdate")
	endDateStr := r.URL.Query().Get("enddate")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing startdate or enddate", http.StatusBadRequest)
		return
	}

	// ISO 8601 formatını parse et (örnek: "2025-03-01")
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid startdate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid enddate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var caller domain.User
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString != "" {
		token, err := ParseJWT(tokenString)
		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if idVal, ok := claims["user_id"].(float64); ok {
					caller.ID = int(idVal)
				}
				if isAdmin, ok := claims["is_admin"].(bool); ok {
					caller.IsAdmin = isAdmin
				}
			}
		}
	}

	events, err := h.eventService.GetSelfDatedEvents(&caller, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (h *EventHandlers) GetEventTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.eventService.GetEventTypes()
	if err != nil {
		http.Error(w, "Failed to retrieve event types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(types)
}

/*func (h *EventHandlers) GetSelfDatedEvents(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("startdate")
	endDateStr := r.URL.Query().Get("enddate")

	if startDateStr == "" || endDateStr == "" {
		http.Error(w, "Missing startdate or enddate", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(w, "Invalid startdate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(w, "Invalid enddate format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	var caller domain.User

	// ⬇️ JWT'yi cookie'den oku
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	tokenString := tokenCookie.Value

	token, err := ParseJWT(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if idVal, ok := claims["user_id"].(float64); ok {
			caller.ID = int(idVal)
		}
		if isAdmin, ok := claims["is_admin"].(bool); ok {
			caller.IsAdmin = isAdmin
		}
	}

	events, err := h.eventService.GetSelfDatedEvents(&caller, startDate, endDate)
	if err != nil {
		http.Error(w, "Failed to retrieve events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}*/
