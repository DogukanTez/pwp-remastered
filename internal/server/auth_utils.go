package server

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"pwp-remastered/internal/domain"
)

// ExtractUserFromRequest extracts and validates JWT from Authorization header and returns the user info
func ExtractUserFromRequest(r *http.Request) (domain.User, error) {
	var caller domain.User
	tokenString := r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	if tokenString == "" {
		return caller, http.ErrNoCookie // Use as unauthorized error
	}
	token, err := ParseJWT(tokenString)
	if err != nil || !token.Valid {
		return caller, http.ErrNoCookie // Use as unauthorized error
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if idVal, ok := claims["user_id"].(float64); ok {
			caller.ID = int(idVal)
		}
		if username, ok := claims["username"].(string); ok {
			caller.Username = username
		}
		if isAdmin, ok := claims["is_admin"].(bool); ok {
			caller.IsAdmin = isAdmin
		}
	}
	return caller, nil
}
