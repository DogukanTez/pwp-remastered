package services

import (
	"fmt"
	"pwp-remastered/internal/domain"
	"pwp-remastered/internal/store"

	"github.com/matthewhartstonge/argon2"
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
func (s *UserService) UpdateUser(caller *domain.User, user *domain.User) error {
	argon := argon2.DefaultConfig()

	// fmt.Println("caller:", caller)
	// if caller.IsAdmin == false || caller.ID == user.ID {
	// 	return errors.New("Unauthorized")
	// }

	existingUser, err := s.store.GetUser(user.ID)
	if err != nil {
		fmt.Println("Error retrieving user:", err)
		return err
	}
	if existingUser == nil {
		fmt.Println("User not found")
		return err
	} else {
		hashedPassword, err := argon.HashEncoded([]byte(user.HashedPassword))
		if err != nil {
			return err
		}
		user.HashedPassword = string(hashedPassword)
		// Check if the user exists

		return s.store.UpdateUser(caller, user)
	}

}

// DeleteUser removes a user by ID
func (s *UserService) DeleteUser(id int) error {
	return s.store.DeleteUser(id)
}

// ListUsers retrieves all users
func (s *UserService) ListUsers() ([]domain.User, error) {
	return s.store.ListUsers()
}

func (s *UserService) ChangeUserStatus(caller *domain.User, id int) error {
	return s.store.ChangeUserStatus(caller, id)
}

func (s *UserService) UpdateSelfUser(caller *domain.User) error {

	return s.store.UpdateSelfUser(caller)
}
