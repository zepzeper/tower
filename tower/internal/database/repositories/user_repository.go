package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/zepzeper/tower/internal/database/models"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll() ([]models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at, last_login
		FROM users
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.LastLogin,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`
	
	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found: %s", id)
		}
		return user, err
	}
	
	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	query := `
		SELECT id, email, password, name, role, created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`
	
	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLogin,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found with email: %s", email)
		}
		return user, err
	}
	
	return user, nil
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user models.User) error {
	query := `
		INSERT INTO users (id, email, password, name, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)
	
	return err
}

// Update updates an existing user
func (r *UserRepository) Update(user models.User) error {
	query := `
		UPDATE users
		SET email = $1, name = $2, role = $3, updated_at = $4
		WHERE id = $5
	`
	
	user.UpdatedAt = time.Now()
	
	result, err := r.db.Exec(
		query,
		user.Email,
		user.Name,
		user.Role,
		user.UpdatedAt,
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
		return fmt.Errorf("user not found: %s", user.ID)
	}
	
	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id, hashedPassword string) error {
	query := `
		UPDATE users
		SET password = $1, updated_at = $2
		WHERE id = $3
	`
	
	now := time.Now()
	
	result, err := r.db.Exec(
		query,
		hashedPassword,
		now,
		id,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}
	
	return nil
}

// Delete removes a user from the database
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}
	
	return nil
}

// UpdateLastLogin updates a user's last login timestamp
func (r *UserRepository) UpdateLastLogin(id string) error {
	query := `
		UPDATE users
		SET last_login = $1
		WHERE id = $2
	`
	
	now := time.Now()
	
	result, err := r.db.Exec(
		query,
		now,
		id,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}
	
	return nil
}
