package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"errors"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(email, name, password string) (*models.User, error) {
	// Check if user exists
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, repository.ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           "user_" + time.Now().Format("20060102150405"),
		Email:        email,
		Name:         name,
		PasswordHash: string(hashedPassword),
		AuthProvider: "local",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	// Normalize email - trim whitespace and convert to lowercase
	email = strings.TrimSpace(strings.ToLower(email))
	// Trim any whitespace from password
	password = strings.TrimSpace(password)

	// Add debug logging
	log.Printf("Attempting login for email: %s", email)

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if err == repository.ErrUserNotFound {
			log.Printf("User not found: %s", email)
			return nil, errors.New("invalid credentials")
		}
		log.Printf("Database error when finding user: %v", err)
		return nil, err
	}

	log.Printf("User found, comparing password hash")

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		log.Printf("Password comparison failed: %v", err)
		return nil, errors.New("invalid credentials")
	}

	log.Printf("Login successful for user: %s", email)
	return user, nil
}

func (s *AuthService) GetGoogleAuthURL() string {
	// TODO: Implement Google OAuth URL generation
	return ""
}

func (s *AuthService) HandleGoogleCallback(code string) (string, *models.User, error) {
	// TODO: Implement Google OAuth callback handling
	return "", nil, errors.New("not implemented")
}
