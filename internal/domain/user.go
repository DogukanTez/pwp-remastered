package domain

type User struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	IsAdmin        bool   `json:"is_admin"`
	IsUser         bool   `json:"is_user"`
	TenantID       int    `json:"tenant_id"`
	Status         int    `json:"status"`
}
