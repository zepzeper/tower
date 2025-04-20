package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/zepzeper/tower/internal/database"
	"github.com/zepzeper/tower/internal/database/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	jwtSecret         string
	tokenExpiry       time.Duration
  databaseManager   *database.Manager
}

func NewAuthService(jwtSecret string, tokenExpiry time.Duration, databaseManager *database.Manager) *AuthService {
	return &AuthService{
		jwtSecret:    jwtSecret,
		tokenExpiry:  tokenExpiry,
		databaseManager: databaseManager,
	}
}

func (s *AuthService) CreateUser(email, password string) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return s.databaseManager.Repos.Auth().CreateUser(email, string(hashedPassword))
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.databaseManager.Repos.Auth().GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(s.tokenExpiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) UserExists(email string) (bool, error) {
	return s.databaseManager.Repos.Auth().UserExists(email)
}
