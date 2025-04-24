package store

import (
	"database/sql"
	"errors"
	"pwp-remastered/internal/database"
	"pwp-remastered/internal/domain"
)

// UserStore defines the interface for user data operations
type UserStore interface {
	GetUser(id int) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	CreateUser(user *domain.User) error
	UpdateUser(caller *domain.User, user *domain.User) error
	DeleteUser(id int) error
	ListUsers() ([]domain.User, error)
	ChangeUserStatus(id int) error
}

type userDBStore struct {
	db database.Service
}

// NewUserStore creates a new UserStore instance
func NewUserStore(db database.Service) UserStore {
	return &userDBStore{db: db}
}

func (s *userDBStore) GetUser(id int) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, hashed_password, email, first_name, last_name, 
		       is_admin, is_user, tenant_id, status
		FROM users WHERE id = $1`

	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.Email,
		&user.FirstName, &user.LastName, &user.IsAdmin, &user.IsUser,
		&user.TenantID, &user.Status,
	)

	if err == sql.ErrNoRows {
		var ErrUserNotFound = errors.New("user not found")
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userDBStore) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	query := `
		SELECT id, username, hashed_password, email, first_name, last_name, 
		       is_admin, is_user, tenant_id, status
		FROM users WHERE username = $1`

	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.HashedPassword, &user.Email,
		&user.FirstName, &user.LastName, &user.IsAdmin, &user.IsUser,
		&user.TenantID, &user.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userDBStore) CreateUser(user *domain.User) error {
	query := `
		INSERT INTO users (username, hashed_password, email, first_name, last_name, 
		                  is_admin, is_user, tenant_id, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	err := s.db.QueryRow(
		query,
		user.Username, user.HashedPassword, user.Email,
		user.FirstName, user.LastName, user.IsAdmin,
		user.IsUser, user.TenantID, user.Status,
	).Scan(&user.ID)

	return err
}

func (s *userDBStore) UpdateUser(caller *domain.User, user *domain.User) error {
	// Admin değilse, is_admin alanını değiştirmesin
	if !caller.IsAdmin {
		var currentIsAdmin bool
		err := s.db.QueryRow("SELECT is_admin FROM users WHERE id = $1", user.ID).Scan(&currentIsAdmin)
		if err != nil {
			return err
		}
		user.IsAdmin = currentIsAdmin
	}

	query := `
		UPDATE users 
		SET username = $1, hashed_password = $2, email = $3,
		    first_name = $4, last_name = $5, is_admin = $6,
		    is_user = $7, tenant_id = $8, status = $9
		WHERE id = $10`

	result, err := s.db.Exec(
		query,
		user.Username, user.HashedPassword, user.Email,
		user.FirstName, user.LastName, user.IsAdmin,
		user.IsUser, user.TenantID, user.Status,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *userDBStore) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *userDBStore) ListUsers() ([]domain.User, error) {
	query := `
		SELECT id, username, hashed_password, email, first_name, last_name, 
		       is_admin, is_user, tenant_id, status
		FROM users`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.HashedPassword, &user.Email,
			&user.FirstName, &user.LastName, &user.IsAdmin, &user.IsUser,
			&user.TenantID, &user.Status,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userDBStore) ChangeUserStatus(id int) error {
	query := `UPDATE users SET status = 1 - status WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
