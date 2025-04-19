package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-pkgz/auth/v2/token"
	"github.com/tower/internal/database/models"
)

type DBUserStore struct {
	db *sql.DB
}

func NewDBUserStore(db *sql.DB) *DBUserStore {
	return &DBUserStore{db: db}
}

func (s *DBUserStore) Get(username string) (token.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, email, picture FROM users WHERE email = $1`
	row := s.db.QueryRowContext(ctx, query, username)

	var dbUser models.DBUser
	err := row.Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Picture)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return token.User{}, errors.New("user not found")
		}
		return token.User{}, err
	}

	return dbUser.ToTokenUser(), nil
}

func (s *DBUserStore) Add(username string, user token.User) error {
	// This method shouldn't be used for password storage
	// Use CreateUser instead (see below)
	_, err := s.Get(username)
	return err
}

// CreateUserWithPassword handles user creation with password
func (s *DBUserStore) CreateUserWithPassword(user models.DBUser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (id, name, email, picture, password) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Picture,
		user.Password,
	)
	return err
}

// GetWithPassword retrieves user with password for authentication
func (s *DBUserStore) GetWithPassword(email string) (DBUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, email, picture, password FROM users WHERE email = $1`
	row := s.db.QueryRowContext(ctx, query, email)

	var dbUser DBUser
	err := row.Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.Picture, &dbUser.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return DBUser{}, errors.New("user not found")
		}
		return DBUser{}, err
	}

	return dbUser, nil
}
