package server

import (
	"encoding/json"
	"net/http"
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/services"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/matthewhartstonge/argon2"
)

type UserHandlers struct {
	userService *services.UserService
}

func NewUserHandlers(userService *services.UserService) *UserHandlers {
	return &UserHandlers{
		userService: userService,
	}
}

func (h *UserHandlers) RegisterRoutes(r chi.Router) {
	r.Route("/users", func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Get("/", h.ListUsers)
		r.Post("/", h.CreateUser)
		r.Get("/{id}", h.GetUser)
		r.Put("/{id}", h.UpdateUser)
		r.Put("/me", h.UpdateSelfUser)
		r.Put("/me/password", h.UpdateSelfPassword)
		// r.Delete("/{id}", h.DeleteUser)
		r.With(AdminMiddleware).Post("/{id}/status", h.ChangeUserStatus)
	})
	r.Post("/login", h.Login)
}

func (h *UserHandlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	caller, err := ExtractUserFromRequest(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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

	users, err := h.userService.ListUsers(&caller)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string][]domain.User{"users": users}
	jsonResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (h *UserHandlers) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.userService.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResp)
}

func (h *UserHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = id

	// Extract caller from JWT (for now, only ID and is_admin fields are extracted)
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

	if err := h.userService.UpdateUser(&caller, &user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (h *UserHandlers) UpdateSelfUser(w http.ResponseWriter, r *http.Request) {

	var caller domain.User

	if err := json.NewDecoder(r.Body).Decode(&caller); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	if err := h.userService.UpdateSelfUser(&caller); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(caller)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func (h *UserHandlers) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandlers) ChangeUserStatus(w http.ResponseWriter, r *http.Request) {
	var caller domain.User

	if err := json.NewDecoder(r.Body).Decode(&caller); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.userService.ChangeUserStatus(&caller, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AuthMiddleware checks for JWT in Authorization header and enforces authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		token, err := ParseJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if !token.Claims.(jwt.MapClaims)["is_admin"].(bool) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Login returns a JWT token
func (h *UserHandlers) Login(w http.ResponseWriter, r *http.Request) {
	type loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		http.Error(w, "Giriş bilgileri hatalı.", http.StatusUnauthorized)
		return
	}

	if user.Status == 0 {
		http.Error(w, "Kulanıcı hesabı inaktif.", http.StatusForbidden)
		return
	}
	// Use argon2 to verify password
	if ok, _ := argon2.VerifyEncoded([]byte(req.Password), []byte(user.HashedPassword)); !ok {
		http.Error(w, "Giriş bilgileri hatalı.", http.StatusUnauthorized)
		return
	}
	token, err := GenerateJWT(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "token",
	// 	Value:    token,
	// 	HttpOnly: true,
	// 	Secure:   false, // https kullanıyorsan true yap
	// 	SameSite: http.SameSiteStrictMode,
	// 	Path:     "/",
	// 	Expires:  time.Now().Add(1 * time.Hour),
	// })
	// w.WriteHeader(http.StatusOK)

}

func (h *UserHandlers) UpdateSelfPassword(w http.ResponseWriter, r *http.Request) {
	var caller domain.User
	var passwordRequest struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&passwordRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
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

	if err := h.userService.UpdateSelfPassword(&caller, passwordRequest.Password); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
