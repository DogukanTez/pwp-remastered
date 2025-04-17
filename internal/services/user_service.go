package services

import (
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/store"
)

// UserService handles business logic for users
type UserService struct {
	store store.UserStore
}

// NewUserService creates a new user service
func NewUserService(userStore store.UserStore) *UserService {
	return &UserService{
		store: userStore,
	}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*domain.User, error) {
	return s.store.GetUser(id)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.store.GetUserByUsername(username)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(user *domain.User) error {
	return s.store.CreateUser(user)
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(user *domain.User) error {
	return s.store.UpdateUser(user)
}

// DeleteUser removes a user by ID
func (s *UserService) DeleteUser(id int) error {
	return s.store.DeleteUser(id)
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.store.ListUsers()
}
